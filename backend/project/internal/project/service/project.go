package service

import (
	"strings"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// projectService はプロジェクトサービスの実装
type projectService struct {
	repo repository.ProjectRepository
}

// DefaultProjectService はデフォルトのプロジェクトサービス
var DefaultProjectService = NewProjectService(repository.DefaultProjectRepository)

// NewProjectService は新しいプロジェクトサービスを作成する
func NewProjectService(repo repository.ProjectRepository) *projectService {
	return &projectService{repo: repo}
}

// CreateProject プロジェクトを新規作成する
func (s *projectService) CreateProject(name string, description string, ownerID uuid.UUID) (*model.Project, error) {
	// 入力値の検証
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "project name cannot be empty")
	}
	if ownerID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "owner ID cannot be empty")
	}

	// プロジェクトオブジェクト作成
	project := model.Project{
		ID:          uuid.New(), // 新規UUIDを生成
		Name:        strings.ToLower(name),
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// データベースに保存
	return s.repo.CreateProject(&project)
}

// GetProject プロジェクトをIDで取得する
func (s *projectService) GetProject(projectID uuid.UUID) (*model.Project, error) {
	// IDの検証
	if projectID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "project ID cannot be empty")
	}

	// リポジトリからプロジェクトを取得
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "project not found: %v", err)
	}

	return project, nil
}

// ListProjects すべてのプロジェクトを取得する
func (s *projectService) ListProjects() ([]*model.Project, error) {
	return s.repo.ListProjects()
}

// ListProjectsByOwner 特定のオーナーのプロジェクトを取得する
func (s *projectService) ListProjectsByOwner(ownerID uuid.UUID) ([]*model.Project, error) {
	if ownerID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "owner ID cannot be empty")
	}

	return s.repo.GetProjectsByOwnerID(ownerID)
}

// UpdateProject プロジェクト情報を更新する
func (s *projectService) UpdateProject(projectID uuid.UUID, name *string, description *string) (*model.Project, error) {
	// IDの検証
	if projectID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "project ID cannot be empty")
	}

	// 現在のプロジェクト情報を取得
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "project not found: %v", err)
	}

	// 更新するフィールドの処理
	if name != nil {
		if *name == "" {
			return nil, status.Errorf(codes.InvalidArgument, "project name cannot be empty")
		}
		project.Name = *name
	}

	if description != nil {
		project.Description = *description
	}

	// 更新日時を設定
	project.UpdatedAt = time.Now()

	// データベースを更新
	return s.repo.UpdateProject(project)
}

// DeleteProject プロジェクトを削除する
func (s *projectService) DeleteProject(projectID uuid.UUID) error {
	// IDの検証
	if projectID == uuid.Nil {
		return status.Errorf(codes.InvalidArgument, "project ID cannot be empty")
	}

	// プロジェクトの存在確認
	_, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return status.Errorf(codes.NotFound, "project not found: %v", err)
	}

	// プロジェクト削除
	return s.repo.DeleteProject(projectID)
}

// 以下は後方互換性のための関数群
// CreateProject プロジェクトを新規作成する（後方互換性用）
func CreateProject(name string, description string, ownerID uuid.UUID) (*model.Project, error) {
	return DefaultProjectService.CreateProject(name, description, ownerID)
}

// GetProject プロジェクトをIDで取得する（後方互換性用）
func GetProject(projectID uuid.UUID) (*model.Project, error) {
	return DefaultProjectService.GetProject(projectID)
}

// ListProjects すべてのプロジェクトを取得する（後方互換性用）
func ListProjects() ([]*model.Project, error) {
	return DefaultProjectService.ListProjects()
}

// ListProjectsByOwner 特定のオーナーのプロジェクトを取得する（後方互換性用）
func ListProjectsByOwner(ownerID uuid.UUID) ([]*model.Project, error) {
	return DefaultProjectService.ListProjectsByOwner(ownerID)
}

// UpdateProject プロジェクト情報を更新する（後方互換性用）
func UpdateProject(projectID uuid.UUID, name *string, description *string) (*model.Project, error) {
	return DefaultProjectService.UpdateProject(projectID, name, description)
}

// DeleteProject プロジェクトを削除する（後方互換性用）
func DeleteProject(projectID uuid.UUID) error {
	return DefaultProjectService.DeleteProject(projectID)
}
