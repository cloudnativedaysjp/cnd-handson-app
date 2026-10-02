// Package telemetry は docs/conventions.md の計測の約束を Go サービスで共通化する
package telemetry

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const MetricsAddr = ":9464"

// Setup は TracerProvider / MeterProvider / propagator を設定し、/metrics を MetricsAddr で公開する。
// endpoint・service name などは OTEL_* の env から SDK が解決する。返す関数で終了処理をする
func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	// エクスポートの失敗（collector 停止など）はリクエストを止めず、JSON ログに残す
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { slog.Warn("otel export failed", "err", err) }))

	res := resource.Default()
	traceExp, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExp), sdktrace.WithResource(res))

	reg := prometheus.NewRegistry()
	promReader, err := otelprom.New(otelprom.WithRegisterer(reg))
	if err != nil {
		return nil, err
	}
	opts := []sdkmetric.Option{sdkmetric.WithResource(res), sdkmetric.WithReader(promReader)}
	if os.Getenv("OTEL_METRICS_EXPORTER") == "otlp" {
		metricExp, err := otlpmetrichttp.New(ctx)
		if err != nil {
			return nil, err
		}
		opts = append(opts, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)))
	}
	mp := sdkmetric.NewMeterProvider(opts...)

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ln, err := net.Listen("tcp", MetricsAddr)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("metrics server stopped", "err", err)
		}
	}()

	return func(ctx context.Context) error {
		return errors.Join(srv.Shutdown(ctx), tp.Shutdown(ctx), mp.Shutdown(ctx))
	}, nil
}
