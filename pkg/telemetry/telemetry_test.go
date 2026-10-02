package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func lastLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[len(lines)-1]), &m))
	return m
}

func TestLoggerWritesConventionKeys(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "handson-test")
	var buf bytes.Buffer
	log := NewLogger(&buf)
	ctx, span := sdktrace.NewTracerProvider().Tracer("t").Start(context.Background(), "op")
	defer span.End()

	log.With("extra", 1).InfoContext(ctx, "hello")

	m := lastLine(t, &buf)
	assert.Equal(t, "info", m["level"])
	assert.Equal(t, "hello", m["msg"])
	assert.Equal(t, "handson-test", m["service"])
	assert.Equal(t, span.SpanContext().TraceID().String(), m["trace_id"])
	assert.Equal(t, span.SpanContext().SpanID().String(), m["span_id"])
	_, err := time.Parse(time.RFC3339, m["time"].(string))
	assert.NoError(t, err)
}

func TestLogUnaryLogsMethodAndCode(t *testing.T) {
	var buf bytes.Buffer
	intercept := logUnary(NewLogger(&buf))
	handler := func(context.Context, any) (any, error) { return nil, status.Error(codes.NotFound, "x") }

	_, err := intercept(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/idp.IdpService/Login"}, handler)
	assert.Equal(t, codes.NotFound, status.Code(err))
	m := lastLine(t, &buf)
	assert.Equal(t, "/idp.IdpService/Login", m["rpc.method"])
	assert.Equal(t, "NotFound", m["code"])

	buf.Reset()
	_, _ = intercept(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: healthPrefix + "Check"}, handler)
	assert.Empty(t, buf.String(), "health checks are not logged")
}

// collector が落ちていても起動・計測・終了が止まらず、/metrics は常に返る
func TestSetupServesMetricsWithoutCollector(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:1")
	t.Setenv("OTEL_METRICS_EXPORTER", "otlp")
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	start := time.Now()
	shutdown, err := Setup(context.Background())
	require.NoError(t, err)
	assert.Less(t, time.Since(start), time.Second)

	counter, err := otel.Meter("test").Int64Counter("http.server.request.count")
	require.NoError(t, err)
	counter.Add(context.Background(), 1)
	_, span := otel.Tracer("test").Start(context.Background(), "op")
	span.End()

	res, err := http.Get("http://127.0.0.1" + MetricsAddr + "/metrics")
	require.NoError(t, err)
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Contains(t, string(body), "http_server_request_count")

	// エクスポートの失敗は返ってよいが、timeout を超えて止まってはいけない
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stopStart := time.Now()
	_ = shutdown(ctx)
	assert.Less(t, time.Since(stopStart), 3*time.Second)
}
