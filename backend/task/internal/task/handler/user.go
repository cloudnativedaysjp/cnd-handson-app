package handler

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

// RequireUserID は TaskService の呼び出しに x-user-id（JWT は入口が検証済み）を求め、スパンに user.id を付ける。
// health check は metadata を送らないので対象外
func RequireUserID(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if !strings.HasPrefix(info.FullMethod, "/task.TaskService/") {
		return handler(ctx, req)
	}
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get("x-user-id")
	if len(vals) != 1 {
		return nil, status.Error(codes.Unauthenticated, "x-user-id is required")
	}
	id, err := uuid.Parse(vals[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "x-user-id must be a UUID")
	}
	trace.SpanFromContext(ctx).SetAttributes(semconv.UserID(id.String()))
	return handler(ctx, req)
}
