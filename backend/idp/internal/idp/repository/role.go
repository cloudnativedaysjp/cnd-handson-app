package repository

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	return r.first(ctx, "id = ?", id)
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*model.Role, error) {
	return r.first(ctx, "name = ?", name)
}

func (r *roleRepository) first(ctx context.Context, query string, arg any) (*model.Role, error) {
	var role model.Role
	err := r.db.WithContext(ctx).Where(query, arg).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}
