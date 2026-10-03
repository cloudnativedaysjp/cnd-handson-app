package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/repository"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Project, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]*model.Project, error)
	Create(ctx context.Context, p *model.Project) error
	Update(ctx context.Context, p *model.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectService struct {
	repo Repository
}

func NewProjectService(repo Repository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, req *projectpb.CreateProjectRequest) (*model.Project, error) {
	if req.GetName() == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidArgument)
	}
	ownerID, err := ParseID(req.GetOwnerId())
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

func (s *ProjectService) Get(ctx context.Context, id string) (*model.Project, error) {
	pid, err := ParseID(id)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.Get(ctx, pid)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return p, err
}

func (s *ProjectService) List(ctx context.Context, ownerID string) ([]*model.Project, error) {
	oid := uuid.Nil
	if ownerID != "" {
		var err error
		if oid, err = ParseID(ownerID); err != nil {
			return nil, err
		}
	}
	return s.repo.List(ctx, oid)
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
	pid, err := ParseID(id)
	if err != nil {
		return err
	}
	err = s.repo.Delete(ctx, pid)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

func ParseID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid id %q", ErrInvalidArgument, s)
	}
	return id, nil
}
