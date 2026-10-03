package repository

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type Filter struct {
	ColumnID   uuid.UUID
	AssigneeID uuid.UUID
	ProjectID  uuid.UUID
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Get(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List は uuid.Nil の条件を無視し、新しい順に返す。page / pageSize が 0 以下なら全件
func (r *TaskRepository) List(ctx context.Context, f Filter, page, pageSize int32) ([]*model.Task, int32, error) {
	query := r.db.WithContext(ctx).Model(&model.Task{})
	if f.ColumnID != uuid.Nil {
		query = query.Where("column_id = ?", f.ColumnID)
	}
	if f.AssigneeID != uuid.Nil {
		query = query.Where("assignee_id = ?", f.AssigneeID)
	}
	// project_id を指定しないときは、プロジェクトなしのタスクだけを返す（他人のプロジェクトのタスクを出さない）
	if f.ProjectID != uuid.Nil {
		query = query.Where("project_id = ?", f.ProjectID)
	} else {
		query = query.Where("project_id = ? OR project_id IS NULL", uuid.Nil)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page > 0 && pageSize > 0 {
		query = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}
	var tasks []*model.Task
	if err := query.Order("start_time DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, int32(total), nil
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Task{}, id)
	if res.Error == nil && res.RowsAffected == 0 {
		return ErrNotFound
	}
	return res.Error
}
