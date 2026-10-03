package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	idppb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/idp"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	taskpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/task"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// registerAPI は frontend 向けの REST API。一覧は project 経由で取り、タスクの書き込みだけ task を直接呼ぶ。
// task も呼び出し元がプロジェクトの所有者かを確かめる
func registerAPI(route func(string, http.HandlerFunc), idp idppb.IdpServiceClient, projects projectpb.ProjectServiceClient, tasks taskpb.TaskServiceClient, verify verifyFunc) {
	// ブラウザは gRPC を呼べないので、入口が idp の Login を中継する。refresh token は使わないので返さない
	route("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		var body struct{ Email, Password string }
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<10)).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest)
			return
		}
		resp, err := idp.Login(r.Context(), &idppb.LoginRequest{Email: body.Email, Password: body.Password})
		if err != nil {
			writeError(w, httpStatus(status.Code(err)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"accessToken": resp.GetAccessToken()})
	})
	api := func(pattern string, call func(ctx context.Context, r *http.Request) (proto.Message, error)) {
		route(pattern, func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok {
				writeError(w, http.StatusUnauthorized)
				return
			}
			user, err := verify(r.Context(), token)
			if err != nil {
				writeError(w, http.StatusUnauthorized)
				return
			}
			resp, err := call(userid.NewContext(r.Context(), user), r)
			if err != nil {
				writeError(w, httpStatus(status.Code(err)))
				return
			}
			b, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(resp)
			if err != nil {
				writeError(w, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(b)
		})
	}
	api("GET /api/projects", func(ctx context.Context, _ *http.Request) (proto.Message, error) {
		return projects.ListProjects(ctx, &projectpb.ListProjectsRequest{})
	})
	api("GET /api/projects/{id}/tasks", func(ctx context.Context, r *http.Request) (proto.Message, error) {
		return projects.ListProjectTasks(ctx, &projectpb.ListProjectTasksRequest{ProjectId: r.PathValue("id")})
	})
	api("GET /api/projects/{id}/columns", func(ctx context.Context, r *http.Request) (proto.Message, error) {
		return projects.ListProjectColumns(ctx, &projectpb.ListProjectColumnsRequest{ProjectId: r.PathValue("id")})
	})
	api("POST /api/projects/{id}/tasks", func(ctx context.Context, r *http.Request) (proto.Message, error) {
		var body struct{ Title, Status string }
		if err := decode(r, &body); err != nil {
			return nil, err
		}
		return tasks.CreateTask(ctx, &taskpb.CreateTaskRequest{Title: body.Title, Status: body.Status, ProjectId: r.PathValue("id")})
	})
	api("PATCH /api/tasks/{id}", func(ctx context.Context, r *http.Request) (proto.Message, error) {
		var body struct{ Status string }
		if err := decode(r, &body); err != nil {
			return nil, err
		}
		return tasks.UpdateTask(ctx, &taskpb.UpdateTaskRequest{
			Id: r.PathValue("id"), Task: &taskpb.Task{Status: body.Status},
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"status"}},
		})
	})
}

// decode は JSON の本文を読む。読めなければ 400 になるよう InvalidArgument を返す
func decode(r *http.Request, v any) error {
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<10)).Decode(v); err != nil {
		return status.Error(codes.InvalidArgument, "invalid body")
	}
	return nil
}

func httpStatus(c codes.Code) int {
	switch c {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	default:
		return http.StatusBadGateway
	}
}

// 下流のエラーメッセージは内部の情報を含みうるので、ステータスの文言だけ返す
func writeError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(code)})
}
