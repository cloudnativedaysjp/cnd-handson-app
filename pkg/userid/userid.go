// Package userid は gRPC metadata の x-user-id を受け取り、下流へ引き継ぐ（JWT の検証は入口が行う）
package userid

import (
	"context"
	"strings"

	"github.com/google/uuid"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const Key = "x-user-id"

type ctxKey struct{}

// Require は servicePrefix（例: "/task.TaskService/"）の呼び出しに UUID の x-user-id を求め、
// ctx に入れてスパンに user.id を付ける。health check は metadata を送らないので prefix で外す
func Require(servicePrefix string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !strings.HasPrefix(info.FullMethod, servicePrefix) {
			return handler(ctx, req)
		}
		md, _ := metadata.FromIncomingContext(ctx)
		vals := md.Get(Key)
		if len(vals) != 1 {
			return nil, status.Error(codes.Unauthenticated, "x-user-id is required")
		}
		id, err := uuid.Parse(vals[0])
		// ゼロの UUID は「未指定」と区別できないので受け付けない
		if err != nil || id == uuid.Nil {
			return nil, status.Error(codes.Unauthenticated, "x-user-id must be a non-zero UUID")
		}
		trace.SpanFromContext(ctx).SetAttributes(semconv.UserID(id.String()))
		return handler(NewContext(ctx, id.String()), req)
	}
}

// NewContext は ctx に呼び出し元のユーザー ID を入れる
func NewContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext は Require が ctx に入れたユーザー ID を返す
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKey{}).(string)
	return id, ok
}

// Forward は Require で受けた x-user-id を下流の呼び出しの metadata に付ける
func Forward() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if id, ok := FromContext(ctx); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, Key, id)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
