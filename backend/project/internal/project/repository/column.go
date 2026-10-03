package repository

import (
	"context"

	columnpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/column"
	"github.com/google/uuid"
)

// ColumnRepository は column サービスからプロジェクトの列を取る。column の board_id にはプロジェクトの ID が入る
type ColumnRepository struct {
	client columnpb.ColumnServiceClient
}

func NewColumnRepository(client columnpb.ColumnServiceClient) *ColumnRepository {
	return &ColumnRepository{client: client}
}

// ListByProject は page を指定せず、プロジェクトの列を全件返す
func (r *ColumnRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*columnpb.Column, error) {
	res, err := r.client.ListColumns(ctx, &columnpb.ListColumnsRequest{BoardId: projectID.String()})
	if err != nil {
		return nil, err
	}
	return res.GetColumns(), nil
}
