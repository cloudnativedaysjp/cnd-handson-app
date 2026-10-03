package telemetry

import (
	"context"
	"log/slog"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// healthcheck は数秒おきに来るので、スパンとログから外す
const healthPrefix = "/grpc.health.v1.Health/"

// ServerOptions は otelgrpc の stats handler と、リクエストごとに rpc.method / code を 1 行出す interceptor
func ServerOptions(log *slog.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithFilter(filters.Not(filters.HealthCheck())))),
		grpc.ChainUnaryInterceptor(logUnary(log)),
	}
}

func ClientOptions() []grpc.DialOption {
	return []grpc.DialOption{grpc.WithStatsHandler(otelgrpc.NewClientHandler())}
}

func logUnary(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if !strings.HasPrefix(info.FullMethod, healthPrefix) {
			log.InfoContext(ctx, "rpc", "rpc.method", info.FullMethod, "code", status.Code(err).String())
		}
		return resp, err
	}
}
