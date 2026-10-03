package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [command] [args]")
		fmt.Println("Commands: create-task, update-task, get-task, delete-task")
		os.Exit(1)
	}

	command := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// task は x-user-id を必須にしている。既定は契約テストの demo ユーザー
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", cmp.Or(os.Getenv("TASK_USER_ID"), "00000000-0000-4000-8000-000000000071"))
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
	client := taskpb.NewTaskServiceClient(conn)

	switch command {
	case "create-task":
		if len(os.Args) < 7 {
			log.Fatalf("Usage: create-task <title> <description> <status> <column_id> <assignee_id>")
		}
		runCreateTask(ctx, client, os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])

	case "update-task":
		runUpdateTask(ctx, client, os.Args[2:]) // フラグで処理

	case "get-task":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: get-task <id>")
		}
		runGetTask(ctx, client, os.Args[2])
	case "list-tasks":
		if len(os.Args) < 4 {
			log.Fatalf("Usage: list-tasks <column_id>")
		}
		runListTasks(ctx, client, os.Args[2], os.Args[3], 1, 10) // デフォルトのページとページサイズを指定
	case "delete-task":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: delete-task <id>")
		}
		runDeleteTask(ctx, client, os.Args[2])

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func runCreateTask(ctx context.Context, client taskpb.TaskServiceClient, title, description, status, columnID, assigneeID string) {
	resp, err := client.CreateTask(ctx, &taskpb.CreateTaskRequest{
		Title:       title,
		Description: description,
		Status:      status,
		ColumnId:    columnID,
		AssigneeId:  assigneeID,
	})
	if err != nil {
		log.Fatalf("タスク作成失敗: %v", err)
	}
	fmt.Printf("Task created: %+v\n", resp.Task)
}

func runUpdateTask(ctx context.Context, client taskpb.TaskServiceClient, args []string) {
	fs := flag.NewFlagSet("update-task", flag.ExitOnError)

	id := fs.String("id", "", "Task ID to update (required)")
	title := fs.String("title", "", "New title")
	description := fs.String("description", "", "New description")
	status := fs.String("status", "", "New status")
	columnID := fs.String("column_id", "", "New column ID")
	assigneeID := fs.String("assignee_id", "", "New assignee ID")

	err := fs.Parse(args)
	if err != nil {
		log.Fatalf("update-task: %v", err)
	}

	if *id == "" {
		log.Fatal("update-task: -id is required")
	}

	task := &taskpb.Task{Id: *id}
	var updatePaths []string

	if *title != "" {
		task.Title = *title
		updatePaths = append(updatePaths, "title")
	}
	if *description != "" {
		task.Description = *description
		updatePaths = append(updatePaths, "description")
	}
	if *status != "" {
		task.Status = *status
		updatePaths = append(updatePaths, "status")
	}
	if *columnID != "" {
		task.ColumnId = *columnID
		updatePaths = append(updatePaths, "column_id")
	}
	if *assigneeID != "" {
		task.AssigneeId = *assigneeID
		updatePaths = append(updatePaths, "assignee_id")
	}

	if len(updatePaths) == 0 {
		log.Fatal("update-task: at least one field to update must be specified")
	}

	req := &taskpb.UpdateTaskRequest{
		Id:         *id,
		Task:       task,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: updatePaths},
	}

	resp, err := client.UpdateTask(ctx, req)
	if err != nil {
		log.Fatalf("UpdateTask failed: %v", err)
	}

	fmt.Printf("Task updated: %+v\n", resp.Task)
}

func runGetTask(ctx context.Context, client taskpb.TaskServiceClient, id string) {
	resp, err := client.GetTask(ctx, &taskpb.GetTaskRequest{Id: id})
	if err != nil {
		log.Fatalf("タスク取得失敗: %v", err)
	}
	fmt.Printf("Task: %+v\n", resp.Task)
}

func runListTasks(ctx context.Context, client taskpb.TaskServiceClient, columnID string, assigneeID string, page int32, pageSize int32) {
	resp, err := client.ListTasks(ctx, &taskpb.ListTasksRequest{
		ColumnId:   columnID,
		AssigneeId: assigneeID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		log.Fatalf("タスク一覧取得失敗: %v", err)
	}
	fmt.Printf("タスク一覧: \n")
	for _, task := range resp.Tasks {
		fmt.Printf("ID: %s, Title: %s, Description: %s, Status: %s\n", task.Id, task.Title, task.Description, task.Status)
	}
	fmt.Printf("Total Count: %d\n", resp.TotalCount)

}
func runDeleteTask(ctx context.Context, client taskpb.TaskServiceClient, id string) {
	resp, err := client.DeleteTask(ctx, &taskpb.DeleteTaskRequest{Id: id})
	if err != nil {
		log.Fatalf("タスク削除失敗: %v", err)
	}
	fmt.Printf("削除成功: %v\n", resp.Success)
}
