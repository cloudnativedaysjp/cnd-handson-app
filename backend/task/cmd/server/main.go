package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/handler"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/service"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/pkg/db"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-service [server|migrate]")
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
	lis, err := net.Listen("tcp", ":"+cmp.Or(os.Getenv("PORT"), "50051"))
	if err != nil {
		return err
	}

	opts := append(telemetry.ServerOptions(slog.Default()), grpc.ChainUnaryInterceptor(userid.Require("/task.TaskService/")))
	grpcServer := grpc.NewServer(opts...)
	svc := service.NewTaskService(repository.NewTaskRepository(conn))
	taskpb.RegisterTaskServiceServer(grpcServer, handler.NewTaskServiceServer(svc))
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	serveErr := make(chan error, 1)
	go func() { serveErr <- grpcServer.Serve(lis) }()
	slog.Info("listening", "grpc", lis.Addr().String())

	return telemetry.WaitAndStop(5*time.Second, serveErr, telemetry.GRPCStop(grpcServer), shutdownTelemetry)
}

func runMigrate() error {
	conn, err := db.Open()
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if err := conn.AutoMigrate(&model.Task{}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	slog.Info("migration completed")
	return nil
}
