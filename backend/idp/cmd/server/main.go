package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/handler"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/service"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/token"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/pkg/db"
	idppb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/idp"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: idp-service [server|migrate]")
		os.Exit(1)
	}
	slog.SetDefault(telemetry.NewLogger(os.Stdout))

	var err error
	switch os.Args[1] {
	case "server":
		err = runServer()
	case "migrate":
		err = runMigrate()
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}
	if err != nil {
		slog.Error("exit", "err", err)
		os.Exit(1)
	}
}

func runServer() error {
	shutdownTelemetry, err := telemetry.Setup(context.Background())
	if err != nil {
		return fmt.Errorf("telemetry: %w", err)
	}
	conn, err := db.Open()
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	key, err := token.ParsePrivateKey(os.Getenv("IDP_SIGNING_KEY"))
	if err != nil {
		return fmt.Errorf("IDP_SIGNING_KEY: %w", err)
	}
	issuer, err := token.NewRS256Issuer(key, os.Getenv("IDP_ISS"), os.Getenv("IDP_AUD"))
	if err != nil {
		return fmt.Errorf("token issuer: %w", err)
	}

	grpcLis, err := net.Listen("tcp", ":"+envOr("PORT", "50051"))
	if err != nil {
		return err
	}
	httpLis, err := net.Listen("tcp", ":"+envOr("HTTP_PORT", "8080"))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(telemetry.ServerOptions(slog.Default())...)
	svc := service.NewIdpService(
		repository.NewUserRepository(conn),
		repository.NewRoleRepository(conn),
		repository.NewRefreshTokenRepository(conn),
		issuer,
	)
	idppb.RegisterIdpServiceServer(grpcServer, handler.NewIdpServiceServer(svc))
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	httpServer := &http.Server{
		Handler:           handler.NewHTTPHandler(issuer.JWK(), os.Getenv("IDP_ISS")),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := grpcServer.Serve(grpcLis); err != nil {
			slog.Error("gRPC server stopped", "err", err)
		}
	}()
	go func() {
		if err := httpServer.Serve(httpLis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server stopped", "err", err)
		}
	}()
	slog.Info("listening", "grpc", grpcLis.Addr().String(), "http", httpLis.Addr().String())

	return telemetry.WaitAndStop(10*time.Second,
		func(context.Context) error { grpcServer.GracefulStop(); return nil },
		httpServer.Shutdown,
		shutdownTelemetry,
	)
}

func runMigrate() error {
	conn, err := db.Open()
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if err := conn.AutoMigrate(&model.Role{}, &model.User{}, &model.RefreshToken{}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := seed(context.Background(), conn); err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	slog.Info("migration completed")
	return nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
