package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/repository"
	columnpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/column"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Project, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]*model.Project, error)
	Create(ctx context.Context, p *model.Project) error
	Update(ctx context.Context, p *model.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Tasks interface {
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*taskpb.Task, error)
}

type Columns interface {
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*columnpb.Column, error)
}

type ProjectService struct {
	repo    Repository
	tasks   Tasks
	columns Columns
}

func NewProjectService(repo Repository, tasks Tasks, columns Columns) *ProjectService {
	return &ProjectService{repo: repo, tasks: tasks, columns: columns}
}

func (s *ProjectService) Create(ctx context.Context, req *projectpb.CreateProjectRequest) (*model.Project, error) {
	if req.GetName() == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidArgument)
	}
	// リクエストの owner_id は使わず、呼び出し元を所有者にする（#184）
	ownerID, err := caller(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	p := &model.Project{
		ID:          uuid.New(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Get は呼び出し元が所有者のときだけプロジェクトを返す。所有者でなければ、存在を隠すため NotFound にする
func (s *ProjectService) Get(ctx context.Context, id string) (*model.Project, error) {
	pid, err := ParseID(id)
	if err != nil {
		return nil, err
	}
	user, err := caller(ctx)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.Get(ctx, pid)
	if errors.Is(err, repository.ErrNotFound) || (err == nil && p.OwnerID != user) {
		return nil, ErrNotFound
	}
	return p, err
}

// List は呼び出し元が所有するプロジェクトだけを返す
func (s *ProjectService) List(ctx context.Context) ([]*model.Project, error) {
	user, err := caller(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.List(ctx, user)
}

// Update は空でないフィールドだけを書き換える
func (s *ProjectService) Update(ctx context.Context, req *projectpb.UpdateProjectRequest) (*model.Project, error) {
	p, err := s.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if req.GetName() != "" {
		p.Name = req.GetName()
	}
	if req.GetDescription() != "" {
		p.Description = req.GetDescription()
	}
	p.UpdatedAt = time.Now()
	err = s.repo.Update(ctx, p)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) Delete(ctx context.Context, id string) error {
	p, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, p.ID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// ListTasks はプロジェクトがあることを確かめてから、task サービスに問い合わせる
func (s *ProjectService) ListTasks(ctx context.Context, projectID string) ([]*taskpb.Task, error) {
	p, err := s.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return s.tasks.ListByProject(ctx, p.ID)
}

// ListColumns はプロジェクトがあることを確かめてから、column サービスに問い合わせる
func (s *ProjectService) ListColumns(ctx context.Context, projectID string) ([]*columnpb.Column, error) {
	p, err := s.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return s.columns.ListByProject(ctx, p.ID)
}

// CheckAccess は呼び出し元がプロジェクトの所有者かを確かめる。task / column が操作の前に使う
func (s *ProjectService) CheckAccess(ctx context.Context, projectID string) error {
	_, err := s.Get(ctx, projectID)
	return err
}

func caller(ctx context.Context) (uuid.UUID, error) {
	id, ok := userid.FromContext(ctx)
	if !ok {
		return uuid.Nil, ErrUnauthenticated
	}
	return uuid.Parse(id)
}

func ParseID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid id %q", ErrInvalidArgument, s)
	}
	return id, nil
}
