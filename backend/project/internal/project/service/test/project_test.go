package service_test

import (
	"context"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/service"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
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

func newServiceWithTasks(tasks service.Tasks) *service.ProjectService {
	return service.NewProjectService(&fakeRepo{projects: map[uuid.UUID]*model.Project{}}, tasks)
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

func TestCreateGetListUpdateDelete(t *testing.T) {
	ctx := context.Background()
	svc := newService()
	owner := uuid.NewString()

	p, err := svc.Create(ctx, &projectpb.CreateProjectRequest{Name: "p1", Description: "d", OwnerId: owner})
	require.NoError(t, err)
	_, err = svc.Create(ctx, &projectpb.CreateProjectRequest{Name: "other", OwnerId: uuid.NewString()})
	require.NoError(t, err)

	got, err := svc.Get(ctx, p.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "p1", got.Name)

	mine, err := svc.List(ctx, owner)
	require.NoError(t, err)
	assert.Len(t, mine, 1)
	all, err := svc.List(ctx, "")
	require.NoError(t, err)
	assert.Len(t, all, 2)

	updated, err := svc.Update(ctx, &projectpb.UpdateProjectRequest{Id: p.ID.String(), Name: "p2"})
	require.NoError(t, err)
	assert.Equal(t, "p2", updated.Name)
	assert.Equal(t, "d", updated.Description, "empty fields are left as is")

	require.NoError(t, svc.Delete(ctx, p.ID.String()))
	_, err = svc.Get(ctx, p.ID.String())
	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.ErrorIs(t, svc.Delete(ctx, p.ID.String()), service.ErrNotFound)
}

func TestRejectsBadInput(t *testing.T) {
	ctx := context.Background()
	svc := newService()
	_, err := svc.Create(ctx, &projectpb.CreateProjectRequest{OwnerId: uuid.NewString()})
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.Create(ctx, &projectpb.CreateProjectRequest{Name: "p", OwnerId: "bad"})
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.Get(ctx, "bad")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.List(ctx, "bad")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}

func TestListTasksRequiresExistingProject(t *testing.T) {
	ctx := context.Background()
	tasks := &fakeTasks{byProject: map[uuid.UUID][]*taskpb.Task{}}
	svc := newServiceWithTasks(tasks)
	p, err := svc.Create(ctx, &projectpb.CreateProjectRequest{Name: "p", OwnerId: uuid.NewString()})
	require.NoError(t, err)
	tasks.byProject[p.ID] = []*taskpb.Task{{Id: "t1", ProjectId: p.ID.String()}}

	got, err := svc.ListTasks(ctx, p.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "t1", got[0].GetId())

	_, err = svc.ListTasks(ctx, uuid.NewString())
	assert.ErrorIs(t, err, service.ErrNotFound)
}
