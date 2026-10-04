package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	columnpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/column"
	idppb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/idp"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	jwt "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testUser = "0b6f6d6e-4b1a-4c3e-9a35-3f1f6c1d2e01"

type fakeProjects struct {
	projectpb.ProjectServiceClient
	gotUser string
	created *projectpb.CreateProjectRequest
}

func (f *fakeProjects) CreateProject(_ context.Context, req *projectpb.CreateProjectRequest, _ ...grpc.CallOption) (*projectpb.ProjectResponse, error) {
	f.created = req
	return &projectpb.ProjectResponse{Project: &projectpb.Project{Id: "p3", Name: req.GetName()}}, nil
}

func (f *fakeProjects) ListProjects(ctx context.Context, _ *projectpb.ListProjectsRequest, _ ...grpc.CallOption) (*projectpb.ListProjectsResponse, error) {
	f.gotUser, _ = userid.FromContext(ctx)
	return &projectpb.ListProjectsResponse{Projects: []*projectpb.Project{{Id: "p1", Name: "demo"}}}, nil
}

func (f *fakeProjects) ListProjectTasks(context.Context, *projectpb.ListProjectTasksRequest, ...grpc.CallOption) (*projectpb.ListProjectTasksResponse, error) {
	return nil, status.Error(codes.NotFound, "project p2 owned by someone else")
}

type fakeTasks struct {
	taskpb.TaskServiceClient
	created *taskpb.CreateTaskRequest
	updated *taskpb.UpdateTaskRequest
	deleted string
}

func (f *fakeTasks) GetTask(_ context.Context, req *taskpb.GetTaskRequest, _ ...grpc.CallOption) (*taskpb.TaskResponse, error) {
	return &taskpb.TaskResponse{Task: &taskpb.Task{Id: req.GetId(), Description: "details"}}, nil
}

func (f *fakeTasks) DeleteTask(_ context.Context, req *taskpb.DeleteTaskRequest, _ ...grpc.CallOption) (*taskpb.DeleteTaskResponse, error) {
	f.deleted = req.GetId()
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (f *fakeTasks) CreateTask(_ context.Context, req *taskpb.CreateTaskRequest, _ ...grpc.CallOption) (*taskpb.TaskResponse, error) {
	f.created = req
	return &taskpb.TaskResponse{Task: &taskpb.Task{Id: "t1", Title: req.GetTitle()}}, nil
}

func (f *fakeTasks) UpdateTask(_ context.Context, req *taskpb.UpdateTaskRequest, _ ...grpc.CallOption) (*taskpb.TaskResponse, error) {
	f.updated = req
	return &taskpb.TaskResponse{Task: &taskpb.Task{Id: req.GetId(), ColumnId: req.GetTask().GetColumnId()}}, nil
}

type fakeColumns struct {
	columnpb.ColumnServiceClient
	created *columnpb.CreateColumnRequest
	updated *columnpb.UpdateColumnRequest
	deleted string
}

func (f *fakeColumns) UpdateColumn(_ context.Context, req *columnpb.UpdateColumnRequest, _ ...grpc.CallOption) (*columnpb.ColumnResponse, error) {
	f.updated = req
	return &columnpb.ColumnResponse{Column: &columnpb.Column{Id: req.GetId(), Name: req.GetName()}}, nil
}

func (f *fakeColumns) DeleteColumn(_ context.Context, req *columnpb.DeleteColumnRequest, _ ...grpc.CallOption) (*columnpb.DeleteColumnResponse, error) {
	f.deleted = req.GetId()
	return &columnpb.DeleteColumnResponse{Success: true}, nil
}

func (f *fakeColumns) CreateColumn(_ context.Context, req *columnpb.CreateColumnRequest, _ ...grpc.CallOption) (*columnpb.ColumnResponse, error) {
	f.created = req
	return &columnpb.ColumnResponse{Column: &columnpb.Column{Id: "c1", Name: req.GetName(), BoardId: req.GetBoardId()}}, nil
}

func TestAPI(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	b64 := base64.RawURLEncoding.EncodeToString
	jwksSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kty": "RSA", "kid": "k1", "n": b64(key.N.Bytes()), "e": b64(big.NewInt(int64(key.E)).Bytes()),
		}}})
	}))
	defer jwksSrv.Close()

	sign := func(kid string, mod func(*jwt.RegisteredClaims)) string {
		c := jwt.RegisteredClaims{
			Issuer: "iss", Audience: jwt.ClaimStrings{"aud"}, Subject: testUser,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		}
		if mod != nil {
			mod(&c)
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
		tok.Header["kid"] = kid
		s, err := tok.SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	hs256, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Subject: testUser}).SignedString([]byte("secret"))

	projects := &fakeProjects{}
	tasks := &fakeTasks{}
	columns := &fakeColumns{}
	h := newHandler(telemetry.NewLogger(io.Discard), "blue", t.TempDir(), nil, projects, tasks, columns, newVerifier(jwksSrv.URL, "iss", "aud"))

	for name, tc := range map[string]struct {
		path, auth string
		code       int
	}{
		"valid":         {"/api/projects", "Bearer " + sign("k1", nil), http.StatusOK},
		"no token":      {"/api/projects", "", http.StatusUnauthorized},
		"wrong aud":     {"/api/projects", "Bearer " + sign("k1", func(c *jwt.RegisteredClaims) { c.Audience = jwt.ClaimStrings{"other"} }), http.StatusUnauthorized},
		"wrong iss":     {"/api/projects", "Bearer " + sign("k1", func(c *jwt.RegisteredClaims) { c.Issuer = "other" }), http.StatusUnauthorized},
		"expired":       {"/api/projects", "Bearer " + sign("k1", func(c *jwt.RegisteredClaims) { c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute)) }), http.StatusUnauthorized},
		"no exp":        {"/api/projects", "Bearer " + sign("k1", func(c *jwt.RegisteredClaims) { c.ExpiresAt = nil }), http.StatusUnauthorized},
		"unknown kid":   {"/api/projects", "Bearer " + sign("k2", nil), http.StatusUnauthorized},
		"hs256":         {"/api/projects", "Bearer " + hs256, http.StatusUnauthorized},
		"downstream NF": {"/api/projects/p2/tasks", "Bearer " + sign("k1", nil), http.StatusNotFound},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		if tc.auth != "" {
			req.Header.Set("Authorization", tc.auth)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != tc.code {
			t.Errorf("%s: got %d %s, want %d", name, rec.Code, rec.Body.String(), tc.code)
		}
		if strings.Contains(rec.Body.String(), "owned by") {
			t.Errorf("%s: downstream message leaked: %s", name, rec.Body.String())
		}
	}

	if projects.gotUser != testUser {
		t.Errorf("x-user-id = %q, want %q", projects.gotUser, testUser)
	}
	call := func(method, path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+sign("k1", nil))
		h.ServeHTTP(rec, req)
		return rec
	}
	if rec := call(http.MethodGet, "/api/projects", ""); !strings.Contains(rec.Body.String(), `"name":"demo"`) {
		t.Errorf("body = %s", rec.Body.String())
	}
	if rec := call(http.MethodPost, "/api/projects/p1/tasks", `{"title":"t","columnId":"c1"}`); rec.Code != http.StatusOK ||
		tasks.created.GetProjectId() != "p1" || tasks.created.GetTitle() != "t" || tasks.created.GetColumnId() != "c1" {
		t.Errorf("create: %d %s %v", rec.Code, rec.Body.String(), tasks.created)
	}
	if rec := call(http.MethodPatch, "/api/tasks/t1", `{"columnId":"c2"}`); rec.Code != http.StatusOK ||
		tasks.updated.GetId() != "t1" || tasks.updated.GetTask().GetColumnId() != "c2" || tasks.updated.GetUpdateMask().GetPaths()[0] != "column_id" {
		t.Errorf("update: %d %s %v", rec.Code, rec.Body.String(), tasks.updated)
	}
	if rec := call(http.MethodPost, "/api/projects/p1/columns", `{"name":"doing"}`); rec.Code != http.StatusOK ||
		columns.created.GetBoardId() != "p1" || columns.created.GetName() != "doing" {
		t.Errorf("create column: %d %s %v", rec.Code, rec.Body.String(), columns.created)
	}
	if rec := call(http.MethodPatch, "/api/tasks/t1", `{"title":"t2","description":""}`); rec.Code != http.StatusOK ||
		strings.Join(tasks.updated.GetUpdateMask().GetPaths(), ",") != "title,description" || tasks.updated.GetTask().GetTitle() != "t2" {
		t.Errorf("update fields: %d %v", rec.Code, tasks.updated)
	}
	if rec := call(http.MethodGet, "/api/tasks/t1", ""); !strings.Contains(rec.Body.String(), `"description":"details"`) {
		t.Errorf("get task: %d %s", rec.Code, rec.Body.String())
	}
	long := strings.Repeat("あ", 2000)
	if rec := call(http.MethodPatch, "/api/tasks/t1", `{"description":"`+long+`"}`); rec.Code != http.StatusOK || tasks.updated.GetTask().GetDescription() != long {
		t.Errorf("long description: %d", rec.Code)
	}
	if rec := call(http.MethodDelete, "/api/tasks/t1", ""); rec.Code != http.StatusOK || tasks.deleted != "t1" {
		t.Errorf("delete task: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodPost, "/api/projects", `{"name":"new","description":"d"}`); rec.Code != http.StatusOK ||
		projects.created.GetName() != "new" || projects.created.GetOwnerId() != "" {
		t.Errorf("create project: %d %v", rec.Code, projects.created)
	}
	if rec := call(http.MethodPatch, "/api/columns/c1", `{"name":"renamed"}`); rec.Code != http.StatusOK ||
		columns.updated.GetId() != "c1" || columns.updated.GetName() != "renamed" {
		t.Errorf("update column: %d %v", rec.Code, columns.updated)
	}
	if rec := call(http.MethodDelete, "/api/columns/c1", ""); rec.Code != http.StatusOK || columns.deleted != "c1" {
		t.Errorf("delete column: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodPatch, "/api/tasks/t1", "not json"); rec.Code != http.StatusBadRequest {
		t.Errorf("bad body: got %d", rec.Code)
	}
}

type fakeIdp struct{ idppb.IdpServiceClient }

func (fakeIdp) Login(_ context.Context, req *idppb.LoginRequest, _ ...grpc.CallOption) (*idppb.TokenResponse, error) {
	if req.GetPassword() != "pw" {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &idppb.TokenResponse{AccessToken: "tok", RefreshToken: "refresh"}, nil
}

func TestLogin(t *testing.T) {
	h := newHandler(telemetry.NewLogger(io.Discard), "blue", t.TempDir(), fakeIdp{}, nil, nil, nil, nil)
	for body, want := range map[string]struct {
		code int
		body string
	}{
		`{"email":"a@example.com","password":"pw"}`:    {http.StatusOK, `{"accessToken":"tok"}`},
		`{"email":"a@example.com","password":"wrong"}`: {http.StatusUnauthorized, `{"error":"Unauthorized"}`},
		`not json`: {http.StatusBadRequest, `{"error":"Bad Request"}`},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body)))
		if rec.Code != want.code || strings.TrimSpace(rec.Body.String()) != want.body {
			t.Errorf("%s: got %d %s, want %d %s", body, rec.Code, rec.Body.String(), want.code, want.body)
		}
	}
}
