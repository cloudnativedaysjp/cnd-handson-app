package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	Issue(user *model.User, roles []string) (token string, expiresAt time.Time, err error)
}

const (
	refreshTokenTTL = 30 * 24 * time.Hour
	// DefaultRole は Register したユーザーに付ける role。migrate で作る
	DefaultRole = "member"
)

type IdpService struct {
	users   repository.UserRepository
	roles   repository.RoleRepository
	refresh repository.RefreshTokenRepository
	issuer  TokenIssuer
}

func NewIdpService(users repository.UserRepository, roles repository.RoleRepository, refresh repository.RefreshTokenRepository, issuer TokenIssuer) *IdpService {
	return &IdpService{users: users, roles: roles, refresh: refresh, issuer: issuer}
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
	role, err := s.roles.GetByName(ctx, DefaultRole)
	if err != nil {
		return uuid.Nil, fmt.Errorf("default role %q: %w", DefaultRole, err)
	}
	now := time.Now()
	user := &model.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		RoleID:       role.ID,
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
	return s.issueTokens(ctx, user, "")
}

// Refresh はリフレッシュトークンを使い捨てにし、新しい組を返す
func (s *IdpService) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	userID, secret, ok := strings.Cut(refreshToken, ".")
	if !ok {
		return nil, ErrInvalidCredentials
	}
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	stored, err := s.refresh.Get(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if stored.Exp < time.Now().Unix() || bcrypt.CompareHashAndPassword([]byte(stored.Token), []byte(secret)) != nil {
		return nil, ErrInvalidCredentials
	}
	user, err := s.users.GetByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user, stored.Token)
}

// prevHash が空でなければ、その行がまだ残っているときだけ置き換える（同じトークンの並行利用を 1 回に絞る）
func (s *IdpService) issueTokens(ctx context.Context, user *model.User, prevHash string) (*Tokens, error) {
	roles, err := s.roleNames(ctx, user)
	if err != nil {
		return nil, err
	}
	access, exp, err := s.issuer.Issue(user, roles)
	if err != nil {
		return nil, err
	}
	refresh, row, err := newRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	if prevHash == "" {
		err = s.refresh.Save(ctx, row)
	} else {
		var ok bool
		if ok, err = s.refresh.Rotate(ctx, prevHash, row); err == nil && !ok {
			err = ErrInvalidCredentials
		}
	}
	if err != nil {
		return nil, err
	}
	return &Tokens{AccessToken: access, RefreshToken: refresh, ExpiresAt: exp}, nil
}

// user サービスから引き継いだユーザーは role_id が空のことがあるため、見つからなければ roles は空にする
func (s *IdpService) roleNames(ctx context.Context, user *model.User) ([]string, error) {
	role, err := s.roles.GetByID(ctx, user.RoleID)
	if errors.Is(err, repository.ErrNotFound) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	return []string{role.Name}, nil
}

// トークンは "<user_id>.<secret>"。user_id で行を引き、secret を bcrypt で照合する
func newRefreshToken(userID uuid.UUID) (string, *model.RefreshToken, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", nil, err
	}
	secret := base64.RawURLEncoding.EncodeToString(buf)
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	row := &model.RefreshToken{UserID: userID, Token: string(hash), Exp: time.Now().Add(refreshTokenTTL).Unix()}
	return userID.String() + "." + secret, row, nil
}
