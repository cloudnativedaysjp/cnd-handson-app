package handler

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/service"
	idppb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/idp"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdpServiceServer struct {
	idppb.UnimplementedIdpServiceServer
	svc *service.IdpService
}

func NewIdpServiceServer(svc *service.IdpService) *IdpServiceServer {
	return &IdpServiceServer{svc: svc}
}

func (s *IdpServiceServer) Register(ctx context.Context, req *idppb.RegisterRequest) (*idppb.RegisterResponse, error) {
	id, err := s.svc.Register(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &idppb.RegisterResponse{UserId: id.String()}, nil
}

func (s *IdpServiceServer) Login(ctx context.Context, req *idppb.LoginRequest) (*idppb.TokenResponse, error) {
	t, err := s.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return toTokenResponse(t), nil
}

func (s *IdpServiceServer) Refresh(ctx context.Context, req *idppb.RefreshRequest) (*idppb.TokenResponse, error) {
	t, err := s.svc.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return toTokenResponse(t), nil
}

func toTokenResponse(t *service.Tokens) *idppb.TokenResponse {
	return &idppb.TokenResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.ExpiresAt.Unix(),
	}
}

func toStatus(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		trace.SpanFromContext(ctx).RecordError(err) // ログは interceptor の 1 行に任せ、原因はトレースで追う
		return status.Error(codes.Internal, "internal error")
	}
}
