package test

import (
	"os"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/project/internal/project/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	// テスト実行中はDB接続をスキップする
	if err := os.Setenv("DB_SKIP_CONNECT", "true"); err != nil {
		panic("Failed to set environment variable: " + err.Error())
	}
}

// インメモリプロジェクトリポジトリの実装
type InMemoryProjectRepository struct {
	projects map[uuid.UUID]*model.Project
	byOwner  map[uuid.UUID][]*model.Project
}

// NewInMemoryProjectRepository は新しいインメモリリポジトリを作成する
func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		projects: make(map[uuid.UUID]*model.Project),
		byOwner:  make(map[uuid.UUID][]*model.Project),
	}
}

// GetProjectByID はIDによるプロジェクトの取得
func (r *InMemoryProjectRepository) GetProjectByID(projectID uuid.UUID) (*model.Project, error) {
	project, exists := r.projects[projectID]
	if !exists {
		return nil, service.ErrProjectNotFound
	}
	return project, nil
}

// GetProjectsByOwnerID はオーナーIDによるプロジェクト一覧の取得
func (r *InMemoryProjectRepository) GetProjectsByOwnerID(ownerID uuid.UUID) ([]*model.Project, error) {
	projects, exists := r.byOwner[ownerID]
	if !exists {
		return []*model.Project{}, nil
	}
	return projects, nil
}

// ListProjects は全プロジェクトの取得
func (r *InMemoryProjectRepository) ListProjects() ([]*model.Project, error) {
	projects := make([]*model.Project, 0, len(r.projects))
	for _, p := range r.projects {
		projects = append(projects, p)
	}
	return projects, nil
}

// CreateProject は新規プロジェクトの作成
func (r *InMemoryProjectRepository) CreateProject(project *model.Project) (*model.Project, error) {
	// 既存のIDを持つプロジェクトがないことを確認
	if _, exists := r.projects[project.ID]; exists {
		return nil, service.ErrProjectAlreadyExists
	}

	// プロジェクトを保存
	r.projects[project.ID] = project

	// オーナーごとのプロジェクト一覧を更新
	ownerProjects := r.byOwner[project.OwnerID]
	r.byOwner[project.OwnerID] = append(ownerProjects, project)

	return project, nil
}

// UpdateProject はプロジェクト情報の更新
func (r *InMemoryProjectRepository) UpdateProject(project *model.Project) (*model.Project, error) {
	if _, exists := r.projects[project.ID]; !exists {
		return nil, service.ErrProjectNotFound
	}

	// プロジェクトを更新
	r.projects[project.ID] = project
	return project, nil
}

// DeleteProject はプロジェクトの削除
func (r *InMemoryProjectRepository) DeleteProject(projectID uuid.UUID) error {
	project, exists := r.projects[projectID]
	if !exists {
		return service.ErrProjectNotFound
	}

	// プロジェクトを削除
	delete(r.projects, projectID)

	// オーナーごとのプロジェクト一覧から削除
	ownerProjects := r.byOwner[project.OwnerID]
	var newOwnerProjects []*model.Project
	for _, p := range ownerProjects {
		if p.ID != projectID {
			newOwnerProjects = append(newOwnerProjects, p)
		}
	}
	r.byOwner[project.OwnerID] = newOwnerProjects

	return nil
}

// TestCreateProject はCreateProjectのテスト
func TestCreateProject(t *testing.T) {
	// テスト用リポジトリとサービスの作成
	repo := NewInMemoryProjectRepository()
	projectService := service.NewProjectService(repo)

	// テストデータ
	name := "CloudNative Days"
	description := "CloudNative Days hands-on project"
	ownerID := uuid.New()

	// テスト実行
	result, err := projectService.CreateProject(name, description, ownerID)

	// 検証
	assert.NoError(t, err, "Should not return an error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, name, result.Name, "Project name should match input")
	assert.Equal(t, description, result.Description, "Project description should match input")
	assert.Equal(t, ownerID, result.OwnerID, "Project owner ID should match input")

	// // リポジトリに正しく保存されていることを確認
	// savedProject, err := repo.GetProjectByID(result.ID)
	// assert.NoError(t, err, "Should be able to retrieve the project")
	// assert.Equal(t, name, savedProject.Name, "Saved project name should match")
	// assert.Equal(t, description, savedProject.Description, "Saved project description should match")
	// assert.Equal(t, ownerID, savedProject.OwnerID, "Saved project owner ID should match")
}
