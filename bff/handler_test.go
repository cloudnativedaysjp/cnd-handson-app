package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestColorLogsOneLineWithTrace(t *testing.T) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	var buf bytes.Buffer
	h := newHandler(telemetry.NewLogger(&buf).With("variant", "legacy", "color", "blue"), "blue")

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/color", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "blue" {
		t.Fatalf("got %d %q", rec.Code, rec.Body.String())
	}

	var line map[string]any
	if err := json.Unmarshal(buf.Bytes(), &line); err != nil {
		t.Fatalf("want one JSON line, got %q: %v", buf.String(), err)
	}
	for k, want := range map[string]any{"variant": "legacy", "color": "blue", "method": "GET", "path": "/color", "status": float64(200)} {
		if line[k] != want {
			t.Errorf("%s = %v, want %v", k, line[k], want)
		}
	}
	if _, ok := line["duration_ms"].(float64); !ok {
		t.Errorf("duration_ms = %v, want number", line["duration_ms"])
	}
	if line["trace_id"] == "" {
		t.Error("trace_id is empty")
	}
}
