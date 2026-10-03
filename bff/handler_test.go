package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestColorLogsOneLineWithTrace(t *testing.T) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider())
	var buf bytes.Buffer
	h := newHandler(telemetry.NewLogger(&buf).With("variant", "legacy", "color", "blue"), "blue", t.TempDir(), nil, nil, nil, nil)

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

func TestWebFallsBackToIndex(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("index"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "assets"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("js"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := newHandler(telemetry.NewLogger(io.Discard), "blue", dir, nil, nil, nil, nil)

	for path, want := range map[string]struct {
		code int
		body string
	}{
		"/":                  {http.StatusOK, "index"},
		"/projects/1":        {http.StatusOK, "index"},
		"/assets/":           {http.StatusOK, "index"},
		"/assets/app.js":     {http.StatusOK, "js"},
		"/assets/missing.js": {http.StatusNotFound, "404 page not found\n"},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != want.code || rec.Body.String() != want.body {
			t.Errorf("%s: got %d %q, want %d %q", path, rec.Code, rec.Body.String(), want.code, want.body)
		}
	}
}
