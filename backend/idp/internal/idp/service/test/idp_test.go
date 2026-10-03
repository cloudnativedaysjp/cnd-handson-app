package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUsers struct {
	byEmail map[string]*model.User
}

func newFakeUsers() *fakeUsers { return &fakeUsers{byEmail: map[string]*model.User{}} }

func (f *fakeUsers) Create(_ context.Context, u *model.User) error {
	if _, ok := f.byEmail[u.Email]; ok {
		return repository.ErrAlreadyExists
	}
	f.byEmail[u.Email] = u
	return nil
}

func (f *fakeUsers) GetByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return u, nil
}

type fakeIssuer struct{}

func (fakeIssuer) Issue(u *model.User) (string, time.Time, error) {
	return "token-for-" + u.ID.String(), time.Unix(1700000000, 0), nil
}

func newService() *service.IdpService {
	return service.NewIdpService(newFakeUsers(), fakeIssuer{})
}

func TestRegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	id, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)

	tokens, err := svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)
	assert.Equal(t, "token-for-"+id.String(), tokens.AccessToken)
	assert.Equal(t, int64(1700000000), tokens.ExpiresAt.Unix())
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	_, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)
	_, err = svc.Register(ctx, "alice2", "alice@example.com", "other")
	assert.ErrorIs(t, err, service.ErrEmailTaken)
}

func TestRegisterRequiresEmailAndPassword(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	_, err := svc.Register(ctx, "alice", "", "secret")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.Register(ctx, "alice", "alice@example.com", "")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}

func TestLoginRejectsBadCredentials(t *testing.T) {
	ctx := context.Background()
	svc := newService()
	_, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)

	_, err = svc.Login(ctx, "alice@example.com", "wrong")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	_, err = svc.Login(ctx, "nobody@example.com", "secret")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}
