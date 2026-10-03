package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/handler"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/service"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/token"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/pkg/db"
	idppb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/idp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: idp-service [server|migrate]")
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
	key, err := token.ParsePrivateKey(os.Getenv("IDP_SIGNING_KEY"))
	if err != nil {
		log.Fatalf("Invalid IDP_SIGNING_KEY: %v", err)
	}
	issuer, err := token.NewRS256Issuer(key, os.Getenv("IDP_ISS"), os.Getenv("IDP_AUD"))
	if err != nil {
		log.Fatalf("Failed to set up token issuer: %v", err)
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
	svc := service.NewIdpService(repository.NewUserRepository(conn), repository.NewRefreshTokenRepository(conn), issuer)
	idppb.RegisterIdpServiceServer(grpcServer, handler.NewIdpServiceServer(svc))

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
	if err := db.Migrate(conn, &model.User{}, &model.RefreshToken{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed")
}
