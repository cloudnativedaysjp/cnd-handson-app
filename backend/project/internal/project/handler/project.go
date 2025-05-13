package handler

import (
	"context"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/service"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/backend/project/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProjectServiceServer struct {
	projectpb.UnimplementedProjectServiceServer
}

// modelからprotoへの変換ヘルパー関数
func convertToProtoProject(project *model.Project) *projectpb.Project {
	return &projectpb.Project{
		Id:          project.ID.String(),
		Name:        project.Name,
	Description: project.Description,
		OwnerId:     project.OwnerID.String(),
		CreatedAt:   timestamppb.New(project.CreatedAt),
	UpdatedAt:   timestamppb.New(project.UpdatedAt),
	}
}

// CreateProject プロジェクトを作成するgRPCハンドラー
func (s *ProjectServiceServer) CreateProject(ctx context.Context, req *projectpb.CreateProjectRequest) (*projectpb.ProjectResponse, error) {
	// オーナーIDをUUIDに変換
	ownerID, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner_id: %v", err)
	}

	// サービス層に処理を委譲
	projectModel, err := service.CreateProject(req.GetName(), req.GetDescription(), ownerID)
	if err != nil {
		return nil, err
	}

	// レスポンスを返す
	return &projectpb.ProjectResponse{
		Project: convertToProtoProject(projectModel),
	}, nil
}

// GetProject プロジェクトを取得するgRPCハンドラー
func (s *ProjectServiceServer) GetProject(ctx context.Context, req *projectpb.GetProjectRequest) (*projectpb.ProjectResponse, error) {
	// プロジェクトIDをUUIDに変換
	projectID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid project_id: %v", err)
	}

	// サービス層に処理を委譲
	projectModel, err := service.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	// レスポンスを返す
	return &projectpb.ProjectResponse{
		Project: convertToProtoProject(projectModel),
	}, nil
}

// ListProjects プロジェクト一覧を取得するgRPCハンドラー
func (s *ProjectServiceServer) ListProjects(ctx context.Context, req *projectpb.ListProjectsRequest) (*projectpb.ListProjectsResponse, error) {
	var projects []*model.Project
	var err error

	// オーナーIDが指定されている場合は、そのオーナーのプロジェクトのみを取得
	if req.GetOwnerId() != "" {
		ownerID, err := uuid.Parse(req.GetOwnerId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner_id: %v", err)
		}
		projects, err = service.ListProjectsByOwner(ownerID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to list projects by owner: %v", err)
		}
	} else {
		// オーナーIDが指定されていない場合は、すべてのプロジェクトを取得
		projects, err = service.ListProjects()
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list projects: %v", err)
	}

	// protoメッセージに変換
	protoProjects := make([]*projectpb.Project, 0, len(projects))
	for _, project := range projects {
		protoProjects = append(protoProjects, convertToProtoProject(project))
	}

	// レスポンスを返す
	return &projectpb.ListProjectsResponse{
		Projects: protoProjects,
	}, nil
}

// UpdateProject プロジェクト情報を更新するgRPCハンドラー
func (s *ProjectServiceServer) UpdateProject(ctx context.Context, req *projectpb.UpdateProjectRequest) (*projectpb.ProjectResponse, error) {
	// プロジェクトIDをUUIDに変換
	projectID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid project_id: %v", err)
	}

	// 更新するフィールドのポインタを準備
	var namePtr, descPtr *string
	name := req.GetName()
	desc := req.GetDescription()

	// 空でないフィールドのみ更新
	if name != "" {
		namePtr = &name
	}
	if desc != "" {
		descPtr = &desc
	}

	// サービス層に処理を委譲
	projectModel, err := service.UpdateProject(projectID, namePtr, descPtr)
	if err != nil {
		return nil, err
	}

	// レスポンスを返す
	return &projectpb.ProjectResponse{
		Project: convertToProtoProject(projectModel),
	}, nil
}

// DeleteProject プロジェクトを削除するgRPCハンドラー
func (s *ProjectServiceServer) DeleteProject(ctx context.Context, req *projectpb.DeleteProjectRequest) (*projectpb.DeleteProjectResponse, error) {
	// プロジェクトIDをUUIDに変換
	projectID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid project_id: %v", err)
	}

	// サービス層に処理を委譲
	err = service.DeleteProject(projectID)
	if err != nil {
		return nil, err
	}

	// 削除成功のレスポンスを返す
	return &projectpb.DeleteProjectResponse{
		Success: true,
	}, nil
}
