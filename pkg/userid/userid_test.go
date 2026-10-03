package userid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const demo = "00000000-0000-4000-8000-000000000071"

func call(ctx context.Context, method string) (context.Context, error) {
	var got context.Context
	_, err := Require("/task.TaskService/")(ctx, nil, &grpc.UnaryServerInfo{FullMethod: method},
		func(c context.Context, _ any) (any, error) { got = c; return nil, nil })
	return got, err
}

func withUser(id string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(Key, id))
}

func TestRequire(t *testing.T) {
	const m = "/task.TaskService/ListTasks"
	_, err := call(withUser(demo), m)
	assert.NoError(t, err)
	_, err = call(context.Background(), m)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
	_, err = call(withUser("not-a-uuid"), m)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
	_, err = call(context.Background(), "/grpc.health.v1.Health/Check")
	assert.NoError(t, err, "health checks carry no user")
}

func TestForwardCopiesUserToOutgoingMetadata(t *testing.T) {
	ctx, err := call(withUser(demo), "/task.TaskService/ListTasks")
	assert.NoError(t, err)
	var sent metadata.MD
	err = Forward()(ctx, "/task.TaskService/ListTasks", nil, nil, nil,
		func(c context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			sent, _ = metadata.FromOutgoingContext(c)
			return nil
		})
	assert.NoError(t, err)
	assert.Equal(t, []string{demo}, sent.Get(Key))
}
