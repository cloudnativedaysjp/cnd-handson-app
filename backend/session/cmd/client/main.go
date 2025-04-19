package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	sessionpb "github.com/cloudnativedaysjp/cnd-handson-app/backend/session/proto"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/client/main.go [command] [args]")
		fmt.Println("Available commands:")
		fmt.Println("  generate-access-token <user-id>")
		fmt.Println("  validate-access-token <access-token>")
		fmt.Println("  generate-refresh-token <user-id>")
		fmt.Println("  revoke-refresh-token <user-id><refresh-token>")
		fmt.Println("  validate-refresh-token <refresh-token> <user-id>")
		os.Exit(1)
	}

	command := os.Args[1]

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch command {
	case "generate-access-token":
		if len(os.Args) < 3 {
			log.Fatalln("User IDを指定してください")
		}
		userID := os.Args[2]
		runGenerateAccessToken(ctx, conn, userID)

	case "validate-access-token":
		if len(os.Args) < 3 {
			log.Fatalln("Access Tokenを指定してください")
		}
		token := os.Args[2]
		runValidateAccessToken(ctx, conn, token)

	case "generate-refresh-token":
		if len(os.Args) < 3 {
			log.Fatalln("User IDを指定してください")
		}
		userID := os.Args[2]
		runGenerateRefreshToken(ctx, conn, userID)
	case "revoke-refresh-token":
		if len(os.Args) < 4 {
			log.Fatalln("Refresh TokenとUser IDを指定してください")
		}
		token := os.Args[2]
		userID := os.Args[3]
		runRevokeRefreshToken(ctx, conn, token, userID)
	case "validate-refresh-token":
		if len(os.Args) < 4 {
			log.Fatalln("Refresh TokenとUser IDを指定してください")
		}
		token := os.Args[2]
		userID := os.Args[3]
		runValidateRefreshToken(ctx, conn, token, userID)
	default:
		fmt.Println("Unknown command:", command)
	}
}

// AccessToken
func runGenerateAccessToken(ctx context.Context, conn *grpc.ClientConn, userID string) {
	client := sessionpb.NewAccessTokenServiceClient(conn)

	resp, err := client.GenerateAccessToken(ctx, &sessionpb.GenerateAccessTokenRequest{
		UserId: userID,
	})
	if err != nil {
		log.Fatalf("AccessToken生成失敗: %v", err)
	}

	fmt.Printf("AccessToken: %s\nExpiresAt: %d\n", resp.AccessToken, resp.ExpiresAt)
}

func runValidateAccessToken(ctx context.Context, conn *grpc.ClientConn, token string) {
	client := sessionpb.NewAccessTokenServiceClient(conn)

	resp, err := client.ValidateAccessToken(ctx, &sessionpb.ValidateAccessTokenRequest{
		AccessToken: token,
	})
	if err != nil {
		log.Fatalf("トークン検証失敗: %v", err)
	}

	fmt.Printf("Valid: %v\nUserID: %s\nError: %s\n", resp.Valid, resp.UserId, resp.Error)
}

// RefreshToken
func runGenerateRefreshToken(ctx context.Context, conn *grpc.ClientConn, userID string) {
	client := sessionpb.NewRefreshTokenServiceClient(conn)

	resp, err := client.GenerateRefreshToken(ctx, &sessionpb.GenerateRefreshTokenRequest{
		UserId: userID,
	})

	if err != nil {
		log.Fatalf("RefreshToken生成失敗: %v", err)
	}

	fmt.Printf("RefreshToken: %s\nExpiresAt: %d\n", resp.RefreshToken, resp.ExpiresAt)
}

func runRevokeRefreshToken(ctx context.Context, conn *grpc.ClientConn, token string, userID string) {
	client := sessionpb.NewRefreshTokenServiceClient(conn)

	resp, err := client.RevokeRefreshToken(ctx, &sessionpb.RevokeRefreshTokenRequest{
		RefreshToken: token,
		UserId:       userID,
	})
	if err != nil {
		log.Fatalf("RefreshToken無効化失敗: %v", err)
	}

	fmt.Printf("Success: %v\n", resp.Success)
}

func runValidateRefreshToken(ctx context.Context, conn *grpc.ClientConn, token string, userid string) {
	client := sessionpb.NewRefreshTokenServiceClient(conn)

	resp, err := client.ValidateRefreshToken(ctx, &sessionpb.ValidateRefreshTokenRequest{
		RefreshToken: token,
		UserId:       userid,
	})
	if err != nil {
		log.Fatalf("RefreshToken検証失敗: %v", err)
	}

	fmt.Printf("Valid: %v\n", resp.Valid)
}
