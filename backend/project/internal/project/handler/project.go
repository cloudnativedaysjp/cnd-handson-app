package handler

import (
	"context"
	"errors"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/service"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProjectServiceServer struct {
	projectpb.UnimplementedProjectServiceServer
	svc *service.ProjectService
}

func NewProjectServiceServer(svc *service.ProjectService) *ProjectServiceServer {
	return &ProjectServiceServer{svc: svc}
}

func (s *ProjectServiceServer) CreateProject(ctx context.Context, req *projectpb.CreateProjectRequest) (*projectpb.ProjectResponse, error) {
	p, err := s.svc.Create(ctx, req)
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &projectpb.ProjectResponse{Project: toProto(p)}, nil
}

func (s *ProjectServiceServer) GetProject(ctx context.Context, req *projectpb.GetProjectRequest) (*projectpb.ProjectResponse, error) {
	p, err := s.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &projectpb.ProjectResponse{Project: toProto(p)}, nil
}

func (s *ProjectServiceServer) ListProjects(ctx context.Context, req *projectpb.ListProjectsRequest) (*projectpb.ListProjectsResponse, error) {
	projects, err := s.svc.List(ctx, req.GetOwnerId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	res := &projectpb.ListProjectsResponse{}
	for _, p := range projects {
		res.Projects = append(res.Projects, toProto(p))
	}
	return res, nil
}

func (s *ProjectServiceServer) UpdateProject(ctx context.Context, req *projectpb.UpdateProjectRequest) (*projectpb.ProjectResponse, error) {
	p, err := s.svc.Update(ctx, req)
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	return &projectpb.ProjectResponse{Project: toProto(p)}, nil
}

func (s *ProjectServiceServer) DeleteProject(ctx context.Context, req *projectpb.DeleteProjectRequest) (*projectpb.DeleteProjectResponse, error) {
	if err := s.svc.Delete(ctx, req.GetId()); err != nil {
		return nil, toStatus(ctx, err)
	}
	return &projectpb.DeleteProjectResponse{Success: true}, nil
}

func (s *ProjectServiceServer) ListProjectTasks(ctx context.Context, req *projectpb.ListProjectTasksRequest) (*projectpb.ListProjectTasksResponse, error) {
	tasks, err := s.svc.ListTasks(ctx, req.GetProjectId())
	if err != nil {
		return nil, toStatus(ctx, err)
	}
	res := &projectpb.ListProjectTasksResponse{}
	for _, t := range tasks {
		res.Tasks = append(res.Tasks, &projectpb.ProjectTask{
			Id:        t.GetId(),
			Title:     t.GetTitle(),
			Status:    t.GetStatus(),
			ColumnId:  t.GetColumnId(),
			ProjectId: t.GetProjectId(),
		})
	}
	return res, nil
}

func toProto(p *model.Project) *projectpb.Project {
	return &projectpb.Project{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		OwnerId:     p.OwnerID.String(),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

func toStatus(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case status.Code(err) == codes.Unavailable || status.Code(err) == codes.DeadlineExceeded:
		return err // 下流（task）の障害はそのまま伝える
	default:
		trace.SpanFromContext(ctx).RecordError(err) // ログは interceptor の 1 行に任せ、原因はトレースで追う
		return status.Error(codes.Internal, "internal error")
	}
}
