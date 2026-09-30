package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/user"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/client/user_client.go [command] [args]")
		fmt.Println("Available commands:")
		fmt.Println("  create-user <name> <email> <password> <role_id>")
		fmt.Println("  update-user <id> <name> <email> <password> <role_id>")
		fmt.Println("  get-user <id>")
		fmt.Println("  delete-user <id>")
		fmt.Println("  verify-password <user_id> <password>")
		os.Exit(1)
	}

	command := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}()

	client := userpb.NewUserServiceClient(conn)

	switch command {
	case "create-user":
		if len(os.Args) < 6 {
			log.Fatalln("Usage: create-user <name> <email> <password> <role_id>")
		}
		runCreateUser(ctx, client, os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	case "update-user":
		if len(os.Args) < 3 {
			log.Fatalln("Usage: update-user <id> [name] [email] [password] [role_id]")
		}
		runUpdateUser(ctx, conn, os.Args[2:])
	case "get-user":
		if len(os.Args) < 3 {
			log.Fatalln("Usage: get-user <id>")
		}
		runGetUser(ctx, client, os.Args[2])
	case "delete-user":
		if len(os.Args) < 3 {
			log.Fatalln("Usage: delete-user <id>")
		}
		runDeleteUser(ctx, client, os.Args[2])
	case "verify-password":
		if len(os.Args) < 4 {
			log.Fatalln("Usage: verify-password <user_id> <password>")
		}
		runVerifyPassword(ctx, client, os.Args[2], os.Args[3])
	default:
		fmt.Println("Unknown command:", command)
	}
}

func runCreateUser(ctx context.Context, client userpb.UserServiceClient, name, email, password, roleID string) {
	resp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		RoleId:   roleID,
	})
	if err != nil {
		log.Fatalf("ユーザー作成失敗: %v", err)
	}
	fmt.Printf("User created: %+v\n", resp.User)
}

func runUpdateUser(ctx context.Context, conn *grpc.ClientConn, args []string) {
	if len(args) < 3 {
		log.Fatalln("Usage: update-user <id> [name] [email] [password] [role_id]")
	}

	client := userpb.NewUserServiceClient(conn)
	id := args[0]

	req := &userpb.UpdateUserRequest{
		Id: id,
	}

	if len(args) > 1 && args[1] != "" {
		req.Name = args[1]
	}
	if len(args) > 2 && args[2] != "" {
		req.Email = args[2]
	}
	if len(args) > 3 && args[3] != "" {
		req.Password = args[3]
	}
	if len(args) > 4 && args[4] != "" {
		req.RoleId = args[4]
	}

	resp, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("UpdateUser失敗: %v", err)
	}

	fmt.Printf("User updated: %v\n", resp.User)
}

func runGetUser(ctx context.Context, client userpb.UserServiceClient, id string) {
	resp, err := client.GetUser(ctx, &userpb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("ユーザー取得失敗: %v", err)
	}
	fmt.Printf("User: %+v\n", resp.User)
}

func runDeleteUser(ctx context.Context, client userpb.UserServiceClient, id string) {
	resp, err := client.DeleteUser(ctx, &userpb.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("ユーザー削除失敗: %v", err)
	}
	fmt.Printf("削除成功: %v\n", resp.Success)
}

func runVerifyPassword(ctx context.Context, client userpb.UserServiceClient, email string, password string) {
	resp, err := client.VerifyPassword(ctx, &userpb.VerifyPasswordRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("パスワード検証失敗: %v", err)
	}
	fmt.Printf("Password valid: %v\n", resp.User)
}
