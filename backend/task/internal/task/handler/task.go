package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/service"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServiceServer struct {
	taskpb.UnimplementedTaskServiceServer
	svc *service.TaskService
}

func NewTaskServiceServer(svc *service.TaskService) *TaskServiceServer {
	return &TaskServiceServer{svc: svc}
}

func (s *TaskServiceServer) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.TaskResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	task, err := s.svc.Get(ctx, id)
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &taskpb.TaskResponse{Task: toProto(task)}, nil
}

func (s *TaskServiceServer) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	columnID, err := service.ParseOptionalUUID(req.GetColumnId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	assigneeID, err := service.ParseOptionalUUID(req.GetAssigneeId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	tasks, total, err := s.svc.List(ctx, repository.Filter{ColumnID: columnID, AssigneeID: assigneeID}, req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	res := &taskpb.ListTasksResponse{TotalCount: total}
	for _, t := range tasks {
		res.Tasks = append(res.Tasks, toProto(t))
	}
	return res, nil
}

func (s *TaskServiceServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.TaskResponse, error) {
	columnID, err := service.ParseOptionalUUID(req.GetColumnId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	assigneeID, err := service.ParseOptionalUUID(req.GetAssigneeId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	task, err := s.svc.Create(ctx, req.GetTitle(), req.GetDescription(), req.GetStatus(), columnID, assigneeID)
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &taskpb.TaskResponse{Task: toProto(task)}, nil
}

func (s *TaskServiceServer) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.TaskResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	task, err := s.svc.Update(ctx, id, req.GetTask(), req.GetUpdateMask())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &taskpb.TaskResponse{Task: toProto(task)}, nil
}

func (s *TaskServiceServer) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	if err := s.svc.Delete(ctx, id); err != nil {
		return nil, toStatus(ctx, err)
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

// 未指定の ID（uuid.Nil）は空文字で返す
func optionalID(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}

func toProto(t *model.Task) *taskpb.Task {
	return &taskpb.Task{
		Id:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		StartTime:   timestamppb.New(t.Start_time),
		EndTime:     timestamppb.New(t.End_time),
		ColumnId:    optionalID(t.Column_id),
		AssigneeId:  optionalID(t.Assignee_id),
	}
}

func toStatus(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		slog.ErrorContext(ctx, "internal error", "err", err)
		return status.Error(codes.Internal, "internal error")
	}
}
