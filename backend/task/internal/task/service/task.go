package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/repository"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	ErrNotFound        = errors.New("task not found")
	ErrInvalidArgument = errors.New("invalid argument")
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Task, error)
	List(ctx context.Context, f repository.Filter, page, pageSize int32) ([]*model.Task, int32, error)
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	task, err := s.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return task, err
}

func (s *TaskService) List(ctx context.Context, f repository.Filter, page, pageSize int32) ([]*model.Task, int32, error) {
	return s.repo.List(ctx, f, page, pageSize)
}

func (s *TaskService) Create(ctx context.Context, title, description, status string, columnID, assigneeID uuid.UUID) (*model.Task, error) {
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidArgument)
	}
	now := time.Now()
	task := &model.Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      status,
		Start_time:  now,
		End_time:    now,
		Column_id:   columnID,
		Assignee_id: assigneeID,
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, req *taskpb.Task, mask *fieldmaskpb.FieldMask) (*model.Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, path := range mask.GetPaths() {
		switch path {
		case "title":
			task.Title = req.GetTitle()
		case "description":
			task.Description = req.GetDescription()
		case "status":
			task.Status = req.GetStatus()
		case "column_id":
			if task.Column_id, err = ParseOptionalUUID(req.GetColumnId()); err != nil {
				return nil, err
			}
		case "assignee_id":
			if task.Assignee_id, err = ParseOptionalUUID(req.GetAssigneeId()); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("%w: unsupported field %q", ErrInvalidArgument, path)
		}
	}
	task.End_time = time.Now()
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// ParseOptionalUUID は空文字を uuid.Nil（未指定）として扱う
func ParseOptionalUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	return id, nil
}
