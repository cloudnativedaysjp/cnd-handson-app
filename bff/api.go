package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// registerAPI は frontend 向けの REST API。タスクと列も project 経由で取り、入口からの gRPC の呼び先を project に絞る
func registerAPI(route func(string, http.HandlerFunc), projects projectpb.ProjectServiceClient, verify verifyFunc) {
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
