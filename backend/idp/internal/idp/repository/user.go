package repository

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrAlreadyExists
	}
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return r.first(ctx, "id = ?", id)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.first(ctx, "email = ?", email)
}

func (r *userRepository) first(ctx context.Context, query string, arg any) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where(query, arg).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
