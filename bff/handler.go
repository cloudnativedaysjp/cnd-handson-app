package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	idppb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/idp"
	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func newHandler(log *slog.Logger, color, webDir string, idp idppb.IdpServiceClient, projects projectpb.ProjectServiceClient, verify verifyFunc) http.Handler {
	mux := http.NewServeMux()
	// パターン（"GET /color"）をそのままスパン名にする
	route := func(pattern string, h http.HandlerFunc) {
		mux.Handle(pattern, otelhttp.NewHandler(accessLog(log, h), pattern))
	}
	web := http.FileServer(http.Dir(webDir))
	route("GET /", func(w http.ResponseWriter, r *http.Request) {
		fi, err := os.Stat(filepath.Join(webDir, filepath.Clean(r.URL.Path)))
		switch {
		case err == nil && !fi.IsDir():
			web.ServeHTTP(w, r)
		// /projects/1 などの画面の URL は react-router が解決する。ディレクトリも一覧を出さず index.html にする
		case filepath.Ext(r.URL.Path) == "":
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
		// 無いアセットは 404 にする。HTML を返すとブラウザが JS として読んで分かりにくく壊れる
		default:
			http.NotFound(w, r)
		}
	})
	route("GET /color", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, color)
	})
	registerAPI(route, idp, projects, verify)
	// probe は数秒おきに来るため、スパンとログを出さない
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

// trace_id をログに載せるため、otelhttp の内側（span の中）で書く
func accessLog(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
