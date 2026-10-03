package service_test

import (
	"context"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/service"
	columnpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/column"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	projects map[uuid.UUID]*model.Project
}

type fakeTasks struct {
	byProject map[uuid.UUID][]*taskpb.Task
}

func (f *fakeTasks) ListByProject(_ context.Context, id uuid.UUID) ([]*taskpb.Task, error) {
	return f.byProject[id], nil
}

func newService() *service.ProjectService {
	return newServiceWithTasks(&fakeTasks{})
}

type fakeColumns struct {
	byProject map[uuid.UUID][]*columnpb.Column
}

func (f *fakeColumns) ListByProject(_ context.Context, id uuid.UUID) ([]*columnpb.Column, error) {
	return f.byProject[id], nil
}

func newServiceWithTasks(tasks service.Tasks) *service.ProjectService {
	return service.NewProjectService(&fakeRepo{projects: map[uuid.UUID]*model.Project{}}, tasks, &fakeColumns{})
}

func (f *fakeRepo) Get(_ context.Context, id uuid.UUID) (*model.Project, error) {
	p, ok := f.projects[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (f *fakeRepo) List(_ context.Context, ownerID uuid.UUID) ([]*model.Project, error) {
	var out []*model.Project
	for _, p := range f.projects {
		if ownerID == uuid.Nil || p.OwnerID == ownerID {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *fakeRepo) Create(_ context.Context, p *model.Project) error {
	f.projects[p.ID] = p
	return nil
}

func (f *fakeRepo) Update(_ context.Context, p *model.Project) error {
	f.projects[p.ID] = p
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := f.projects[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.projects, id)
	return nil
}

// as は user が呼び出し元になる ctx を返す（本番では userid.Require が入れる）
func as(user string) context.Context {
	return userid.NewContext(context.Background(), user)
}

func TestCreateGetListUpdateDelete(t *testing.T) {
	owner := uuid.NewString()
	ctx := as(owner)
	svc := newService()

	p, err := svc.Create(ctx, &projectpb.CreateProjectRequest{Name: "p1", Description: "d", OwnerId: uuid.NewString()})
	require.NoError(t, err)
	assert.Equal(t, owner, p.OwnerID.String(), "the caller owns the project, not the requested owner_id")
	_, err = svc.Create(as(uuid.NewString()), &projectpb.CreateProjectRequest{Name: "other"})
	require.NoError(t, err)

	got, err := svc.Get(ctx, p.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "p1", got.Name)

	mine, err := svc.List(ctx)
	require.NoError(t, err)
	assert.Len(t, mine, 1, "only the caller's projects")

	updated, err := svc.Update(ctx, &projectpb.UpdateProjectRequest{Id: p.ID.String(), Name: "p2"})
	require.NoError(t, err)
	assert.Equal(t, "p2", updated.Name)
	assert.Equal(t, "d", updated.Description, "empty fields are left as is")

	require.NoError(t, svc.Delete(ctx, p.ID.String()))
	_, err = svc.Get(ctx, p.ID.String())
	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.ErrorIs(t, svc.Delete(ctx, p.ID.String()), service.ErrNotFound)
}

func TestOthersSeeNotFound(t *testing.T) {
	tasks := &fakeTasks{byProject: map[uuid.UUID][]*taskpb.Task{}}
	columns := &fakeColumns{byProject: map[uuid.UUID][]*columnpb.Column{}}
	svc := service.NewProjectService(&fakeRepo{projects: map[uuid.UUID]*model.Project{}}, tasks, columns)
	p, err := svc.Create(as(uuid.NewString()), &projectpb.CreateProjectRequest{Name: "p"})
	require.NoError(t, err)
	id := p.ID.String()
	other := as(uuid.NewString())

	_, err = svc.Get(other, id)
	assert.ErrorIs(t, err, service.ErrNotFound)
	_, err = svc.Update(other, &projectpb.UpdateProjectRequest{Id: id, Name: "x"})
	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.ErrorIs(t, svc.Delete(other, id), service.ErrNotFound)
	_, err = svc.ListTasks(other, id)
	assert.ErrorIs(t, err, service.ErrNotFound)
	_, err = svc.ListColumns(other, id)
	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.ErrorIs(t, svc.CheckAccess(other, id), service.ErrNotFound)
	mine, err := svc.List(other)
	require.NoError(t, err)
	assert.Empty(t, mine)
}

func TestRequiresCaller(t *testing.T) {
	svc := newService()
	_, err := svc.Create(context.Background(), &projectpb.CreateProjectRequest{Name: "p"})
	assert.ErrorIs(t, err, service.ErrUnauthenticated)
	_, err = svc.List(context.Background())
	assert.ErrorIs(t, err, service.ErrUnauthenticated)
}

func TestRejectsBadInput(t *testing.T) {
	ctx := as(uuid.NewString())
	svc := newService()
	_, err := svc.Create(ctx, &projectpb.CreateProjectRequest{})
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.Get(ctx, "bad")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}

func TestListTasksAndColumnsForOwner(t *testing.T) {
	ctx := as(uuid.NewString())
	tasks := &fakeTasks{byProject: map[uuid.UUID][]*taskpb.Task{}}
	columns := &fakeColumns{byProject: map[uuid.UUID][]*columnpb.Column{}}
	svc := service.NewProjectService(&fakeRepo{projects: map[uuid.UUID]*model.Project{}}, tasks, columns)
	p, err := svc.Create(ctx, &projectpb.CreateProjectRequest{Name: "p"})
	require.NoError(t, err)
	tasks.byProject[p.ID] = []*taskpb.Task{{Id: "t1"}}
	columns.byProject[p.ID] = []*columnpb.Column{{Id: "c1"}}

	gotTasks, err := svc.ListTasks(ctx, p.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "t1", gotTasks[0].GetId())
	gotColumns, err := svc.ListColumns(ctx, p.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "c1", gotColumns[0].GetId())
	assert.NoError(t, svc.CheckAccess(ctx, p.ID.String()))
	_, err = svc.ListTasks(ctx, uuid.NewString())
	assert.ErrorIs(t, err, service.ErrNotFound)
}
