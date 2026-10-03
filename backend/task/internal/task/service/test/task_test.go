package service_test

import (
	"context"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/repository"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/service"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type fakeRepo struct {
	tasks map[uuid.UUID]*model.Task
}

func newFake() *fakeRepo { return &fakeRepo{tasks: map[uuid.UUID]*model.Task{}} }

func (f *fakeRepo) Get(_ context.Context, id uuid.UUID) (*model.Task, error) {
	t, ok := f.tasks[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (f *fakeRepo) List(_ context.Context, flt repository.Filter, _, _ int32) ([]*model.Task, int32, error) {
	var out []*model.Task
	for _, t := range f.tasks {
		if flt.ColumnID != uuid.Nil && t.Column_id != flt.ColumnID {
			continue
		}
		out = append(out, t)
	}
	return out, int32(len(out)), nil
}

func (f *fakeRepo) Create(_ context.Context, t *model.Task) error { f.tasks[t.ID] = t; return nil }

func (f *fakeRepo) Update(_ context.Context, t *model.Task) error { f.tasks[t.ID] = t; return nil }

func (f *fakeRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := f.tasks[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.tasks, id)
	return nil
}

func TestCreateGetListDelete(t *testing.T) {
	ctx := context.Background()
	svc := service.NewTaskService(newFake())

	created, err := svc.Create(ctx, "t1", "d", "todo", uuid.Nil, uuid.Nil)
	require.NoError(t, err)

	got, err := svc.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, "t1", got.Title)

	tasks, total, err := svc.List(ctx, repository.Filter{}, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, int32(1), total)
	assert.Equal(t, created.ID, tasks[0].ID)

	require.NoError(t, svc.Delete(ctx, created.ID))
	_, err = svc.Get(ctx, created.ID)
	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.ErrorIs(t, svc.Delete(ctx, created.ID), service.ErrNotFound)
}

func TestCreateRequiresTitle(t *testing.T) {
	_, err := service.NewTaskService(newFake()).Create(context.Background(), "", "", "", uuid.Nil, uuid.Nil)
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}

func TestUpdateAppliesMaskedFields(t *testing.T) {
	ctx := context.Background()
	svc := service.NewTaskService(newFake())
	created, err := svc.Create(ctx, "t1", "d", "todo", uuid.Nil, uuid.Nil)
	require.NoError(t, err)
	col := uuid.New()

	updated, err := svc.Update(ctx, created.ID,
		&taskpb.Task{Title: "ignored", Status: "done", ColumnId: col.String()},
		&fieldmaskpb.FieldMask{Paths: []string{"status", "column_id"}})
	require.NoError(t, err)
	assert.Equal(t, "t1", updated.Title)
	assert.Equal(t, "done", updated.Status)
	assert.Equal(t, col, updated.Column_id)

	_, err = svc.Update(ctx, created.ID, &taskpb.Task{}, &fieldmaskpb.FieldMask{Paths: []string{"id"}})
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
	_, err = svc.Update(ctx, created.ID, &taskpb.Task{ColumnId: "bad"}, &fieldmaskpb.FieldMask{Paths: []string{"column_id"}})
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}

func TestParseOptionalUUID(t *testing.T) {
	id, err := service.ParseOptionalUUID("")
	require.NoError(t, err)
	assert.Equal(t, uuid.Nil, id)
	_, err = service.ParseOptionalUUID("bad")
	assert.ErrorIs(t, err, service.ErrInvalidArgument)
}
