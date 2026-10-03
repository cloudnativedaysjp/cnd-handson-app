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

type Projects interface {
	CheckAccess(ctx context.Context, projectID uuid.UUID) error
}

type TaskService struct {
	repo     Repository
	projects Projects
}

func NewTaskService(repo Repository, projects Projects) *TaskService {
	return &TaskService{repo: repo, projects: projects}
}

// checkProject はプロジェクトにひも付くときだけ、呼び出し元が所有者かを project に確かめる（#184）
func (s *TaskService) checkProject(ctx context.Context, projectID uuid.UUID) error {
	if projectID == uuid.Nil {
		return nil
	}
	err := s.projects.CheckAccess(ctx, projectID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	task, err := s.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := s.checkProject(ctx, task.Project_id); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) List(ctx context.Context, req *taskpb.ListTasksRequest) ([]*model.Task, int32, error) {
	ids, err := parseOptionalIDs(req.GetColumnId(), req.GetAssigneeId(), req.GetProjectId())
	if err != nil {
		return nil, 0, err
	}
	if err := s.checkProject(ctx, ids[2]); err != nil {
		return nil, 0, err
	}
	f := repository.Filter{ColumnID: ids[0], AssigneeID: ids[1], ProjectID: ids[2]}
	return s.repo.List(ctx, f, req.GetPage(), req.GetPageSize())
}

func (s *TaskService) Create(ctx context.Context, req *taskpb.CreateTaskRequest) (*model.Task, error) {
	if req.GetTitle() == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidArgument)
	}
	ids, err := parseOptionalIDs(req.GetColumnId(), req.GetAssigneeId(), req.GetProjectId())
	if err != nil {
		return nil, err
	}
	if err := s.checkProject(ctx, ids[2]); err != nil {
		return nil, err
	}
	now := time.Now()
	task := &model.Task{
		ID:          uuid.New(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Status:      req.GetStatus(),
		Start_time:  now,
		End_time:    now,
		Column_id:   ids[0],
		Assignee_id: ids[1],
		Project_id:  ids[2],
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
			if task.Column_id, err = parseOptionalUUID(req.GetColumnId()); err != nil {
				return nil, err
			}
		case "assignee_id":
			if task.Assignee_id, err = parseOptionalUUID(req.GetAssigneeId()); err != nil {
				return nil, err
			}
		case "project_id":
			if task.Project_id, err = parseOptionalUUID(req.GetProjectId()); err != nil {
				return nil, err
			}
			// 移動先のプロジェクトも、呼び出し元が所有者でなければならない
			if err := s.checkProject(ctx, task.Project_id); err != nil {
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
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func parseOptionalIDs(ss ...string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(ss))
	for i, s := range ss {
		var err error
		if ids[i], err = parseOptionalUUID(s); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// parseOptionalUUID は空文字を uuid.Nil（未指定）として扱う
func parseOptionalUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	return id, nil
}
