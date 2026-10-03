package repository

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var p model.Project
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List は ownerID のプロジェクトだけを返す。全件を返す経路は作らない
func (r *ProjectRepository) List(ctx context.Context, ownerID uuid.UUID) ([]*model.Project, error) {
	var projects []*model.Project
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Order("created_at DESC").Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// Update は UPDATE だけを実行する（Save だと消えた行を INSERT で作り直してしまう）
func (r *ProjectRepository) Update(ctx context.Context, p *model.Project) error {
	res := r.db.WithContext(ctx).Model(p).Select("*").Updates(p)
	if res.Error == nil && res.RowsAffected == 0 {
		return ErrNotFound
	}
	return res.Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Project{}, id)
	if res.Error == nil && res.RowsAffected == 0 {
		return ErrNotFound
	}
	return res.Error
}
