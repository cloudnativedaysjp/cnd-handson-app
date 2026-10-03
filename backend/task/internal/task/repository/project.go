package repository

import (
	"context"

	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProjectRepository は project に、呼び出し元がプロジェクトの所有者かを問い合わせる
type ProjectRepository struct {
	client projectpb.ProjectServiceClient
}

func NewProjectRepository(client projectpb.ProjectServiceClient) *ProjectRepository {
	return &ProjectRepository{client: client}
}

// CheckAccess は所有者でなければ ErrNotFound を返す
func (r *ProjectRepository) CheckAccess(ctx context.Context, projectID uuid.UUID) error {
	_, err := r.client.CheckProjectAccess(ctx, &projectpb.CheckProjectAccessRequest{ProjectId: projectID.String()})
	if status.Code(err) == codes.NotFound {
		return ErrNotFound
	}
	return err
}
