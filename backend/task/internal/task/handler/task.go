package handler

import (
	"context"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/service"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServiceServer struct {
	taskpb.UnimplementedTaskServiceServer
}

// GetTaskはIDに基づいてタスクを取得するgRPCメソッド
func (s *TaskServiceServer) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.TaskResponse, error) {
	// タスク取得処理（サービス層に委譲）
	task_id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid task_id: %v", err)
	}
	taskModel, err := service.GetTaskByID(task_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}

	return &taskpb.TaskResponse{
		Task: &taskpb.Task{
			Id:          taskModel.ID.String(),
			Title:       taskModel.Title,
			Description: taskModel.Description,
			Status:      taskModel.Status,
			StartTime:   timestamppb.New(taskModel.Start_time),
			EndTime:     timestamppb.New(taskModel.End_time),
			ColumnId:    taskModel.Column_id.String(),
			AssigneeId:  taskModel.Assignee_id.String(),
		},
	}, nil
}

// ListTasksはタスクのリストを取得するgRPCメソッド
func (s *TaskServiceServer) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	// タスクリスト取得処理（サービス層に委譲）
	columnId, err := uuid.Parse(req.GetColumnId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid column_id: %v", err)
	}
	assigneeId, err := uuid.Parse(req.GetAssigneeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid assignee_id: %v", err)
	}
	tasks, totalCount, err := service.ListTasks(columnId, assigneeId, req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	var taskResponses []*taskpb.Task
	for _, taskModel := range tasks {
		taskResponses = append(taskResponses, &taskpb.Task{
			Id:          taskModel.ID.String(),
			Title:       taskModel.Title,
			Description: taskModel.Description,
			Status:      taskModel.Status,
			StartTime:   timestamppb.New(taskModel.Start_time),
			EndTime:     timestamppb.New(taskModel.End_time),
			ColumnId:    taskModel.Column_id.String(),
			AssigneeId:  taskModel.Assignee_id.String(),
		})
	}

	return &taskpb.ListTasksResponse{
		Tasks:      taskResponses,
		TotalCount: totalCount,
	}, nil
}

// CreateTaskは新しいタスクを作成するgRPCメソッド
func (s *TaskServiceServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.TaskResponse, error) {
	column_id := uuid.Nil
	if req.GetColumnId() != "" {
		parsedColumn_id, err := uuid.Parse(req.GetColumnId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid column_id: %v", err)
		}
		column_id = parsedColumn_id
	}
	assignee_id := uuid.Nil
	if req.GetAssigneeId() != "" {
		parsedAssignee_id, err := uuid.Parse(req.GetAssigneeId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid assignee_id: %v", err)
		}
		assignee_id = parsedAssignee_id
	}
	// タスク作成処理（サービス層に委譲）
	taskModel, err := service.CreateTask(req.GetTitle(), req.GetDescription(), req.GetStatus(), column_id, assignee_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	return &taskpb.TaskResponse{
		Task: &taskpb.Task{
			Id:          taskModel.ID.String(),
			Title:       taskModel.Title,
			Description: taskModel.Description,
			Status:      taskModel.Status,
			StartTime:   timestamppb.New(taskModel.Start_time),
			EndTime:     timestamppb.New(taskModel.End_time),
			ColumnId:    taskModel.Column_id.String(),
			AssigneeId:  taskModel.Assignee_id.String(),
		},
	}, nil
}

// UpdateTaskは既存のタスクを更新するgRPCメソッド
func (s *TaskServiceServer) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.TaskResponse, error) {
	task_id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid task_id: %v", err)
	}

	// タスク更新処理（サービス層に委譲）
	taskModel, err := service.UpdateTask(task_id, req.GetTask(), req.GetUpdateMask())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}
	return &taskpb.TaskResponse{Task: &taskpb.Task{
		Id:          taskModel.ID.String(),
		Title:       taskModel.Title,
		Description: taskModel.Description,
		Status:      taskModel.Status,
		StartTime:   timestamppb.New(taskModel.Start_time),
		EndTime:     timestamppb.New(taskModel.End_time),
		ColumnId:    taskModel.Column_id.String(),
		AssigneeId:  taskModel.Assignee_id.String(),
	}}, nil

}

// DeleteTaskはタスクを削除するgRPCメソッド
func (s *TaskServiceServer) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	task_id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid task_id: %v", err)
	}
	// タスク削除処理（サービス層に委譲）
	err = service.DeleteTask(task_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &taskpb.DeleteTaskResponse{
		Success: true,
	}, nil
}
