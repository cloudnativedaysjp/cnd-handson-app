package handler_test

import (
	"context"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/task/internal/task/handler"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func call(ctx context.Context, method string) error {
	_, err := handler.RequireUserID(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method},
		func(context.Context, any) (any, error) { return "ok", nil })
	return err
}

func withUser(id string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-user-id", id))
}

func TestRequireUserID(t *testing.T) {
	const m = "/task.TaskService/ListTasks"
	assert.NoError(t, call(withUser("00000000-0000-4000-8000-000000000071"), m))
	assert.Equal(t, codes.Unauthenticated, status.Code(call(context.Background(), m)))
	assert.Equal(t, codes.Unauthenticated, status.Code(call(withUser("not-a-uuid"), m)))
	assert.NoError(t, call(context.Background(), "/grpc.health.v1.Health/Check"), "health checks carry no user")
}
