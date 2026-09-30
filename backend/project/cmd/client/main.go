package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/backend/project/proto"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: project-client [command] [args]")
		fmt.Println("Available commands:")
		fmt.Println("  create-project <name> <description> <owner_id>")
		fmt.Println("  update-project <id> <name> <description>")
		fmt.Println("  get-project <id>")
		fmt.Println("  list-projects [owner_id]")
		fmt.Println("  delete-project <id>")
		os.Exit(1)
	}

	command := os.Args[1]

	// gRPC接続設定
	target := os.Getenv("PROJECT_ADDR")
	if target == "" {
		target = "localhost:50051"
	}

	// gRPCクライアント接続作成
	ctx := context.Background()
	// staticcheckの警告を回避するために@latestアノテーションを使用
	//nolint:staticcheck // SA1019: grpc.DialContext is still required by this version
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close connection: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := projectpb.NewProjectServiceClient(conn)

	switch command {
	case "create-project":
		if len(os.Args) < 5 {
			log.Fatalln("Usage: create-project <name> <description> <owner_id>")
		}
		runCreateProject(ctx, client, os.Args[2], os.Args[3], os.Args[4])
	case "update-project":
		if len(os.Args) < 4 {
			log.Fatalln("Usage: update-project <id> <name> <description>")
		}
		name := ""
		desc := ""
		if len(os.Args) > 3 {
			name = os.Args[3]
		}
		if len(os.Args) > 4 {
			desc = os.Args[4]
		}
		runUpdateProject(ctx, client, os.Args[2], name, desc)
	case "get-project":
		if len(os.Args) < 3 {
			log.Fatalln("Usage: get-project <id>")
		}
		runGetProject(ctx, client, os.Args[2])
	case "list-projects":
		ownerID := ""
		if len(os.Args) > 2 {
			ownerID = os.Args[2]
		}
		runListProjects(ctx, client, ownerID)
	case "delete-project":
		if len(os.Args) < 3 {
			log.Fatalln("Usage: delete-project <id>")
		}
		runDeleteProject(ctx, client, os.Args[2])
	default:
		fmt.Println("Unknown command:", command)
	}
}

func runCreateProject(ctx context.Context, client projectpb.ProjectServiceClient, name, description, ownerID string) {
	resp, err := client.CreateProject(ctx, &projectpb.CreateProjectRequest{
		Name:        name,
		Description: description,
		OwnerId:     ownerID,
	})
	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}
	fmt.Printf("Project created: ID=%s, Name=%s\n", resp.Project.Id, resp.Project.Name)
	printProjectDetails(resp.Project)
}

func runUpdateProject(ctx context.Context, client projectpb.ProjectServiceClient, id, name, description string) {
	req := &projectpb.UpdateProjectRequest{
		Id:          id,
		Name:        name,
		Description: description,
	}

	resp, err := client.UpdateProject(ctx, req)
	if err != nil {
		log.Fatalf("Failed to update project: %v", err)
	}
	fmt.Printf("Project updated: ID=%s\n", resp.Project.Id)
	printProjectDetails(resp.Project)
}

func runGetProject(ctx context.Context, client projectpb.ProjectServiceClient, id string) {
	resp, err := client.GetProject(ctx, &projectpb.GetProjectRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}
	fmt.Printf("Project details:\n")
	printProjectDetails(resp.Project)
}

func runListProjects(ctx context.Context, client projectpb.ProjectServiceClient, ownerID string) {
	req := &projectpb.ListProjectsRequest{}
	if ownerID != "" {
		req.OwnerId = ownerID
	}

	resp, err := client.ListProjects(ctx, req)
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	fmt.Printf("Total %d projects found\n", len(resp.Projects))
	for i, p := range resp.Projects {
		fmt.Printf("--- Project %d ---\n", i+1)
		printProjectDetails(p)
	}
}

func runDeleteProject(ctx context.Context, client projectpb.ProjectServiceClient, id string) {
	resp, err := client.DeleteProject(ctx, &projectpb.DeleteProjectRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("Failed to delete project: %v", err)
	}
	fmt.Printf("Project deletion: success=%v\n", resp.Success)
}

func printProjectDetails(project *projectpb.Project) {
	fmt.Printf("ID: %s\n", project.Id)
	fmt.Printf("Name: %s\n", project.Name)
	fmt.Printf("Description: %s\n", project.Description)
	fmt.Printf("Owner ID: %s\n", project.OwnerId)
	fmt.Printf("Created at: %v\n", project.CreatedAt.AsTime())
	fmt.Printf("Updated at: %v\n", project.UpdatedAt.AsTime())
	fmt.Println("-----------------------")
}
