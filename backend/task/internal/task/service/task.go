package service

import (
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/repository"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func GetTaskByID(task_id uuid.UUID) (*model.Task, error) {

	task, err := repository.GetTaskByID(task_id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func ListTasks(column_id uuid.UUID, assignee_id uuid.UUID, page int32, page_size int32) ([]*model.Task, int32, error) {
	tasks, totalCount, err := repository.List(column_id, assignee_id, page, page_size)
	if err != nil {
		return nil, 0, err
	}
	if tasks == nil {
		return nil, 0, fmt.Errorf("task not found")
	}
	return tasks, totalCount, nil
}

func CreateTask(title string, description string, status string, column_id uuid.UUID, assignee_id uuid.UUID) (*model.Task, error) {
	// タスクを作成
	task := &model.Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      status,
		Start_time:  time.Now(),
		End_time:    time.Now(),
		Column_id:   column_id,
		Assignee_id: assignee_id}
	err := repository.Create(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func UpdateTask(task_id uuid.UUID, requestTask *taskpb.Task, updateMask *fieldmaskpb.FieldMask) (*model.Task, error) {
	taskModel, err := GetTaskByID(task_id)
	if err != nil {
		return nil, err
	}

	for _, path := range updateMask.Paths {
		switch path {
		case "title":
			taskModel.Title = requestTask.GetTitle()
		case "description":
			taskModel.Description = requestTask.GetDescription()
		case "status":
			taskModel.Status = requestTask.GetStatus()
		case "column_id":
			column_id, err := uuid.Parse(requestTask.GetColumnId())
			if err != nil {
				return nil, fmt.Errorf("invalid column_id: %v", err)
			}
			taskModel.Column_id = column_id
		case "assignee_id":
			assignee_id, err := uuid.Parse(requestTask.GetAssigneeId())
			if err != nil {
				return nil, fmt.Errorf("invalid assignee_id: %v", err)
			}
			taskModel.Assignee_id = assignee_id
		default:
			return nil, fmt.Errorf("unsupported field: %s", path)
		}
	}

	// タスクを更新
	taskModel.End_time = time.Now()
	err = repository.Update(taskModel)
	if err != nil {
		return nil, err
	}

	return taskModel, nil
}

func DeleteTask(task_id uuid.UUID) error {
	err := repository.Delete(task_id)
	if err != nil {
		return err
	}
	return nil
}
