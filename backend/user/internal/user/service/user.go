package service

import (
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/user/internal/user/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/user/internal/user/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/user/pkg/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateUser(Email string, Password string, Name string, RoleID uuid.UUID) (*model.User, error) {
	HashPassword, err := auth.HashPassword(Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:           uuid.New(), // UUIDを生成
		Name:         Name,       // デフォルト値
		Email:        Email,
		PasswordHash: HashPassword,
		RoleID:       uuid.Nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return repository.CreateUser(&user)
}

func GetUserByID(userID uuid.UUID) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	if user == nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return user, nil
}

func GetUserByName(name string) (*model.User, error) {
	user, err := repository.GetUserByName(name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	return user, nil
}

func UpdateUser(userID uuid.UUID, Name *string, Email *string, Password *string, RoleID *uuid.UUID) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if Password != nil {
		HashPassword, err := auth.HashPassword(*Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = HashPassword
	}

	if Name != nil {
		user.Name = *Name
	}

	if Email != nil {
		user.Email = *Email
	}

	if RoleID != nil {
		user.RoleID = *RoleID
	}

	user.UpdatedAt = time.Now()

	return repository.UpdateUser(user)
}

func DeleteUser(userID uuid.UUID) error {

	return repository.DeleteUser(userID)
}

func VerifyPassword(email string, password string) (bool, error) {
	user, err := repository.GetUserByEmail(email)
	// ユーザーが存在しない場合はエラーを返す
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	// パスワードを検証
	if !auth.ComparePassword(user.PasswordHash, password) {
		return false, nil
	}
	return true, nil
}
