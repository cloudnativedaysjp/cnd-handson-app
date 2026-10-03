package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/service"
	"github.com/google/uuid"
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

func (f *fakeUsers) GetByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	for _, u := range f.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeUsers) GetByEmail(_ context.Context, email string) (*model.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return u, nil
}

type fakeRefresh struct {
	byUser       map[uuid.UUID]*model.RefreshToken
	beforeRotate func()
}

func (f *fakeRefresh) Save(_ context.Context, t *model.RefreshToken) error {
	f.byUser[t.UserID] = t
	return nil
}

func (f *fakeRefresh) Get(_ context.Context, id uuid.UUID) (*model.RefreshToken, error) {
	t, ok := f.byUser[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (f *fakeRefresh) Rotate(_ context.Context, prevHash string, next *model.RefreshToken) (bool, error) {
	if f.beforeRotate != nil {
		f.beforeRotate()
	}
	if cur, ok := f.byUser[next.UserID]; !ok || cur.Token != prevHash {
		return false, nil
	}
	f.byUser[next.UserID] = next
	return true, nil
}

type fakeRoles struct {
	roles []model.Role
}

func (f *fakeRoles) GetByID(_ context.Context, id uuid.UUID) (*model.Role, error) {
	for i := range f.roles {
		if f.roles[i].ID == id {
			return &f.roles[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

func (f *fakeRoles) GetByName(_ context.Context, name string) (*model.Role, error) {
	for i := range f.roles {
		if f.roles[i].Name == name {
			return &f.roles[i], nil
		}
	}
	return nil, repository.ErrNotFound
}

// fakeIssuer は渡された roles を記録する
type fakeIssuer struct {
	lastRoles []string
}

func (f *fakeIssuer) Issue(u *model.User, roles []string) (string, time.Time, error) {
	f.lastRoles = roles
	return "token-for-" + u.ID.String(), time.Unix(1700000000, 0), nil
}

type fixture struct {
	svc     *service.IdpService
	users   *fakeUsers
	refresh *fakeRefresh
	issuer  *fakeIssuer
}

func newFixture() *fixture {
	f := &fixture{
		users:   newFakeUsers(),
		refresh: &fakeRefresh{byUser: map[uuid.UUID]*model.RefreshToken{}},
		issuer:  &fakeIssuer{},
	}
	roles := &fakeRoles{roles: []model.Role{{ID: uuid.New(), Name: service.DefaultRole}}}
	f.svc = service.NewIdpService(f.users, roles, f.refresh, f.issuer)
	return f
}

func newService() (*service.IdpService, *fakeRefresh) {
	f := newFixture()
	return f.svc, f.refresh
}

func TestRegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService()

	id, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)

	tokens, err := svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)
	assert.Equal(t, "token-for-"+id.String(), tokens.AccessToken)
	assert.Equal(t, int64(1700000000), tokens.ExpiresAt.Unix())
	assert.True(t, strings.HasPrefix(tokens.RefreshToken, id.String()+"."))
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService()

	_, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)
	_, err = svc.Register(ctx, "alice2", "alice@example.com", "other")
	assert.ErrorIs(t, err, service.ErrEmailTaken)
}

func TestRegisterRequiresEmailAndPassword(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService()

	_, err := svc.Register(ctx, "alice", "", "secret")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.Register(ctx, "alice", "alice@example.com", "")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}

func TestLoginRejectsBadCredentials(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService()
	_, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)

	_, err = svc.Login(ctx, "alice@example.com", "wrong")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	_, err = svc.Login(ctx, "nobody@example.com", "secret")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestRefreshRotatesToken(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService()
	id, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)
	first, err := svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)

	second, err := svc.Refresh(ctx, first.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, "token-for-"+id.String(), second.AccessToken)
	assert.NotEqual(t, first.RefreshToken, second.RefreshToken)

	_, err = svc.Refresh(ctx, first.RefreshToken)
	assert.ErrorIs(t, err, service.ErrInvalidCredentials, "used token must not be reusable")
	_, err = svc.Refresh(ctx, second.RefreshToken)
	assert.NoError(t, err)
}

func TestRefreshRejectsInvalidTokens(t *testing.T) {
	ctx := context.Background()
	svc, refresh := newService()
	id, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)
	tokens, err := svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)

	for _, tok := range []string{"", "no-dot", "not-a-uuid.x", uuid.NewString() + ".x", id.String() + ".wrong"} {
		_, err := svc.Refresh(ctx, tok)
		assert.ErrorIs(t, err, service.ErrInvalidCredentials, tok)
	}

	refresh.byUser[id].Exp = time.Now().Add(-time.Minute).Unix()
	_, err = svc.Refresh(ctx, tokens.RefreshToken)
	assert.ErrorIs(t, err, service.ErrInvalidCredentials, "expired")
}

func TestLoginPutsRoleInToken(t *testing.T) {
	ctx := context.Background()
	f := newFixture()
	_, err := f.svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)

	_, err = f.svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)
	assert.Equal(t, []string{service.DefaultRole}, f.issuer.lastRoles)
}

func TestLoginWithUnknownRoleGivesEmptyRoles(t *testing.T) {
	ctx := context.Background()
	f := newFixture()
	_, err := f.svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)
	f.users.byEmail["alice@example.com"].RoleID = uuid.Nil // user サービスから引き継いだユーザー

	_, err = f.svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)
	assert.NotNil(t, f.issuer.lastRoles)
	assert.Empty(t, f.issuer.lastRoles)
}

// 同じトークンでの Refresh が並行したとき、照合後に先を越されたら失敗する
func TestRefreshFailsWhenTokenRotatedConcurrently(t *testing.T) {
	ctx := context.Background()
	svc, refresh := newService()
	_, err := svc.Register(ctx, "alice", "alice@example.com", "secret")
	require.NoError(t, err)
	tokens, err := svc.Login(ctx, "alice@example.com", "secret")
	require.NoError(t, err)

	refresh.beforeRotate = func() {
		refresh.beforeRotate = nil
		_, err := svc.Refresh(ctx, tokens.RefreshToken) // 先に完了する側
		require.NoError(t, err)
	}
	_, err = svc.Refresh(ctx, tokens.RefreshToken)
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}
