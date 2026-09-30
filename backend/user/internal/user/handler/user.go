package handler

import (
	"context"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/user/internal/user/service"
	userpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/user"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
}

// CreateUserはユーザーを作成するgRPCメソッド
func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	// ユーザー作成処理（サービス層に委譲）
	// ここでユーザーをデータベースに保存などの処理を行う
	roleId, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid role_id: %v", err)
	}
	userModel, err := service.CreateUser(req.GetEmail(), req.GetPassword(), req.GetName(), roleId)
	if err != nil {
		return nil, err
	}

	// ユーザー情報をgRPCレスポンスに変換
	return &userpb.UserResponse{
		User: &userpb.User{
			Id:        userModel.ID.String(),
			Name:      userModel.Name,
			Email:     userModel.Email,
			RoleId:    userModel.RoleID.String(),
			CreatedAt: timestamppb.New(userModel.CreatedAt),
			UpdatedAt: timestamppb.New(userModel.UpdatedAt),
		},
	}, nil
}

// UpdateUserはユーザー情報を更新するgRPCメソッド
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	name := req.GetName()
	email := req.GetEmail()
	password := req.GetPassword()
	roleId, err := uuid.Parse(req.GetRoleId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid role_id: %v", err)
	}

	userId, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	// stringは空文字チェック、uuidはnilチェック
	var namePtr, emailPtr, passwordPtr *string
	var roleIDPtr *uuid.UUID

	if name != "" {
		namePtr = &name
	}
	if email != "" {
		emailPtr = &email
	}
	if password != "" {
		passwordPtr = &password
	}
	roleIDPtr = &roleId // 空文字は事前にエラーにしているのでOK

	userModel, err := service.UpdateUser(userId, namePtr, emailPtr, passwordPtr, roleIDPtr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &userpb.UserResponse{
		User: &userpb.User{
			Id:        userModel.ID.String(),
			Name:      userModel.Name,
			Email:     userModel.Email,
			RoleId:    userModel.RoleID.String(),
			CreatedAt: timestamppb.New(userModel.CreatedAt),
			UpdatedAt: timestamppb.New(userModel.UpdatedAt),
		},
	}, nil
}

// GetUserはIDに基づいてユーザーを取得するgRPCメソッド
func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	// ユーザー取得処理（サービス層に委譲）
	userId, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}
	userModel, err := service.GetUserByID(userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// ユーザーを返却
	return &userpb.UserResponse{
		User: &userpb.User{
			Id:        userModel.ID.String(),
			Name:      userModel.Name,
			Email:     userModel.Email,
			RoleId:    userModel.RoleID.String(),
			CreatedAt: timestamppb.New(userModel.CreatedAt),
			UpdatedAt: timestamppb.New(userModel.UpdatedAt),
		},
	}, nil
}

// DeleteUserはユーザーを削除するgRPCメソッド
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	// ユーザー削除処理（サービス層に委譲）
	userId, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}
	err = service.DeleteUser(userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	// 削除成功のレスポンスを返す
	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}

// vaはユーザーの認証を行うgRPCメソッド
func (s *UserServiceServer) VerifyPassword(ctx context.Context, req *userpb.VerifyPasswordRequest) (*userpb.UserResponse, error) {

	// ユーザー認証処理（サービス層に委譲）
	userModel, err := service.VerifyPassword(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify password: %v", err)
	}

	// ユーザーを返却
	return &userpb.UserResponse{
		User: &userpb.User{
			Id:        userModel.ID.String(),
			Name:      userModel.Name,
			Email:     userModel.Email,
			RoleId:    userModel.RoleID.String(),
			CreatedAt: timestamppb.New(userModel.CreatedAt),
			UpdatedAt: timestamppb.New(userModel.UpdatedAt),
		},
	}, nil
}
