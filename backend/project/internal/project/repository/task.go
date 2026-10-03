package repository

import (
	"context"

	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/google/uuid"
)

// TaskRepository は task サービスからプロジェクトのタスクを取る
type TaskRepository struct {
	client taskpb.TaskServiceClient
}

func NewTaskRepository(client taskpb.TaskServiceClient) *TaskRepository {
	return &TaskRepository{client: client}
}

// ListByProject は page を指定せず、プロジェクトのタスクを全件返す
func (r *TaskRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*taskpb.Task, error) {
	res, err := r.client.ListTasks(ctx, &taskpb.ListTasksRequest{ProjectId: projectID.String()})
	if err != nil {
		return nil, err
	}
	return res.GetTasks(), nil
}
