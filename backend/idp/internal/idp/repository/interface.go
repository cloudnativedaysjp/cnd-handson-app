package repository

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type RefreshTokenRepository interface {
	// Save はユーザーの既存のトークンを置き換える
	Save(ctx context.Context, token *model.RefreshToken) error
	Get(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error)
	// Rotate はハッシュが prevHash のままの行だけを next に置き換え、置き換えたかを返す
	Rotate(ctx context.Context, prevHash string, next *model.RefreshToken) (bool, error)
}

type RoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error)
	GetByName(ctx context.Context, name string) (*model.Role, error)
}
