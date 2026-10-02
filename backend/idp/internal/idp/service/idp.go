package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type TokenIssuer interface {
	Issue(user *model.User) (token string, expiresAt time.Time, err error)
}

type IdpService struct {
	users  repository.UserRepository
	issuer TokenIssuer
	now    func() time.Time
}

func NewIdpService(users repository.UserRepository, issuer TokenIssuer) *IdpService {
	return &IdpService{users: users, issuer: issuer, now: time.Now}
}

func (s *IdpService) Register(ctx context.Context, name, email, password string) (uuid.UUID, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return uuid.Nil, fmt.Errorf("%w: email and password are required", ErrInvalidArgument)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	now := s.now()
	user := &model.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return uuid.Nil, ErrEmailTaken
		}
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (s *IdpService) Login(ctx context.Context, email, password string) (*Tokens, error) {
	user, err := s.users.GetByEmail(ctx, strings.TrimSpace(email))
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	token, exp, err := s.issuer.Issue(user)
	if err != nil {
		return nil, err
	}
	return &Tokens{AccessToken: token, ExpiresAt: exp}, nil
}
