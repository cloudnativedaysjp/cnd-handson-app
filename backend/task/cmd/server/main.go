package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/handler"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/service"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/pkg/db"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-service [server|migrate]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "server":
		runServer()
	case "migrate":
		runMigrate()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func runServer() {
	conn, err := db.Open()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	svc := service.NewTaskService(repository.NewTaskRepository(conn))
	taskpb.RegisterTaskServiceServer(grpcServer, handler.NewTaskServiceServer(svc))
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runMigrate() {
	conn, err := db.Open()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := conn.AutoMigrate(&model.Task{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed")
}
