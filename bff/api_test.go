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
}

func (f *fakeTasks) CreateTask(_ context.Context, req *taskpb.CreateTaskRequest, _ ...grpc.CallOption) (*taskpb.TaskResponse, error) {
	f.created = req
	return &taskpb.TaskResponse{Task: &taskpb.Task{Id: "t1", Title: req.GetTitle()}}, nil
}

func (f *fakeTasks) UpdateTask(_ context.Context, req *taskpb.UpdateTaskRequest, _ ...grpc.CallOption) (*taskpb.TaskResponse, error) {
	f.updated = req
	return &taskpb.TaskResponse{Task: &taskpb.Task{Id: req.GetId(), Status: req.GetTask().GetStatus()}}, nil
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
	h := newHandler(telemetry.NewLogger(io.Discard), "blue", t.TempDir(), nil, projects, tasks, newVerifier(jwksSrv.URL, "iss", "aud"))

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
	if rec := call(http.MethodPost, "/api/projects/p1/tasks", `{"title":"t","status":"todo"}`); rec.Code != http.StatusOK ||
		tasks.created.GetProjectId() != "p1" || tasks.created.GetTitle() != "t" || tasks.created.GetStatus() != "todo" {
		t.Errorf("create: %d %s %v", rec.Code, rec.Body.String(), tasks.created)
	}
	if rec := call(http.MethodPatch, "/api/tasks/t1", `{"status":"done"}`); rec.Code != http.StatusOK ||
		tasks.updated.GetId() != "t1" || tasks.updated.GetTask().GetStatus() != "done" || tasks.updated.GetUpdateMask().GetPaths()[0] != "status" {
		t.Errorf("update: %d %s %v", rec.Code, rec.Body.String(), tasks.updated)
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
	h := newHandler(telemetry.NewLogger(io.Discard), "blue", t.TempDir(), fakeIdp{}, nil, nil, nil)
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
