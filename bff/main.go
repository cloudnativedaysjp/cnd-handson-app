// Command bff は入口（handson-legacy / handson-modern）。VARIANT で color と画面を切り替える
package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	projectpb "github.com/cloudnativedaysjp/cnd-handson-app/gen/go/project"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/userid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var colors = map[string]string{"legacy": "blue", "modern": "green"}

// frontend のビルド成果物。Dockerfile がここに置く
const webDir = "/web"

func main() {
	slog.SetDefault(telemetry.NewLogger(os.Stdout))
	var err error
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		err = healthcheck()
	} else {
		err = run()
	}
	if err != nil {
		slog.Error("exit", "err", err)
		os.Exit(1)
	}
}

func run() error {
	variant := os.Getenv("VARIANT")
	color, ok := colors[variant]
	if !ok {
		return fmt.Errorf("VARIANT must be legacy or modern, got %q", variant)
	}
	env, err := requireEnv("PROJECT_ADDR", "IDP_JWKS_URL", "IDP_ISS", "IDP_AUD")
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(env["PROJECT_ADDR"], append(telemetry.ClientOptions(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(userid.Forward()))...)
	if err != nil {
		return fmt.Errorf("project client: %w", err)
	}
	shutdownTelemetry, err := telemetry.Setup(context.Background())
	if err != nil {
		return fmt.Errorf("telemetry: %w", err)
	}
	lis, err := net.Listen("tcp", ":"+cmp.Or(os.Getenv("PORT"), "8080"))
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler: newHandler(slog.Default().With("variant", variant, "color", color), color, webDir,
			projectpb.NewProjectServiceClient(conn), newVerifier(env["IDP_JWKS_URL"], env["IDP_ISS"], env["IDP_AUD"])),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(lis) }()
	slog.Info("listening", "http", lis.Addr().String(), "variant", variant)

	return telemetry.WaitAndStop(5*time.Second, serveErr, srv.Shutdown,
		func(context.Context) error { return conn.Close() }, shutdownTelemetry)
}

func requireEnv(keys ...string) (map[string]string, error) {
	env := map[string]string{}
	var missing []string
	for _, k := range keys {
		if env[k] = os.Getenv(k); env[k] == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("required env: %s", strings.Join(missing, ", "))
	}
	return env, nil
}

// scratch イメージには wget が無いため、compose の healthcheck は同じバイナリで /healthz を叩く
func healthcheck() error {
	res, err := http.Get("http://127.0.0.1:" + cmp.Or(os.Getenv("PORT"), "8080") + "/healthz")
	if err != nil {
		return err
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("healthz: status %d", res.StatusCode)
	}
	return nil
}
