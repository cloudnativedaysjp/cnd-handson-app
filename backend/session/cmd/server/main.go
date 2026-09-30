package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/session/internal/auth/handler"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/session/internal/auth/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/session/pkg/db"
	sessionpb "github.com/cloudnativedaysjp/cnd-handson-app/backend/session/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myapp [server|migrate]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "server":
		runServer()
	case "migrate":
		runMigrate()
	case "reset":
		resetDB()
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}

func runServer() {
	// init処理
	err := loadEnv()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	if k := os.Getenv("JWT_SECRET_KEY"); k == "" || k == "changeme" {
		log.Fatalf("JWT_SECRET_KEY must be set to a non-default value")
	}
	_, err = db.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// gRPCサーバーの設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPCサーバーの作成
	grpcServer := grpc.NewServer()

	// gRPCサービスの登録
	accessTokenService := &handler.AccessTokenServiceHandler{}
	sessionpb.RegisterAccessTokenServiceServer(grpcServer, accessTokenService)

	refreshTokenService := &handler.RefreshTokenServiceHandler{}
	sessionpb.RegisterRefreshTokenServiceServer(grpcServer, refreshTokenService)

	// Health checkの登録
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthSrv)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrate() {
	conn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = db.MigrateDB(conn, model.RefreshToken{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully.")
}
func resetDB() {
	conn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = db.ResetDB(conn, model.RefreshToken{})
	if err != nil {
		log.Fatalf("Database reset failed: %v", err)
	}

	log.Println("Database reset completed successfully.")
}

func loadEnv() error {
	return nil // 環境変数読み込みの実装（.env）
}
