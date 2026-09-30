package handler

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/session/internal/auth/service"
	sessionpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/session" // your generated proto package
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccessTokenServiceHandler アクセストークンサービスのハンドラー
type AccessTokenServiceHandler struct {
	// Dependency injection for storage or other services
	sessionpb.UnimplementedAccessTokenServiceServer
}

// GenerateAccessToken アクセストークン生成
func (h *AccessTokenServiceHandler) GenerateAccessToken(ctx context.Context, req *sessionpb.GenerateAccessTokenRequest) (*sessionpb.AccessTokenResponse, error) {
	// ユーザーIDを取得し、アクセストークンを生成（例: ランダムなトークン生成）
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}
	accessToken, expiresAt, err := service.GenerateAccessToken(userId)
	if err != nil {
		return nil, err
	}

	// アクセストークンと有効期限をレスポンスに設定
	return &sessionpb.AccessTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

// ValidateAccessToken アクセストークン検証
func (h *AccessTokenServiceHandler) ValidateAccessToken(ctx context.Context, req *sessionpb.ValidateAccessTokenRequest) (*sessionpb.ValidateAccessTokenResponse, error) {
	// ここでアクセストークンの検証ロジックを追加（例: データベースやキャッシュを使用）
	if req.AccessToken == "" {
		return nil, errors.New("access token is empty")
	}
	// アクセストークンを検証
	valid, userID, err := service.ValidateAccessToken(req.AccessToken)
	if err != nil {
		return &sessionpb.ValidateAccessTokenResponse{
			Valid:  false,
			UserId: "",
			Error:  err.Error(),
		}, err
	}

	// 仮に検証成功とする（実際には検証ロジックを追加）
	return &sessionpb.ValidateAccessTokenResponse{
		Valid:  valid,
		UserId: userID.String(), // ユーザーIDを返す
		Error: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}, nil
}
