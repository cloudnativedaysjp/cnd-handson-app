"""docs/conventions.md の計測の約束を column（Python）で満たす。Go の pkg/telemetry と同じ形にする"""

import json
import logging
import os
import sys
import threading
from collections.abc import Callable
from datetime import datetime, timezone

import grpc
from opentelemetry import metrics, propagate, trace
from opentelemetry.baggage.propagation import W3CBaggagePropagator
from opentelemetry.exporter.otlp.proto.http.metric_exporter import OTLPMetricExporter
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.exporter.prometheus import PrometheusMetricReader
from opentelemetry.instrumentation.grpc import (
    client_interceptor,
    filters,
    server_interceptor,
)
from opentelemetry.instrumentation.grpc.grpcext import intercept_channel
from opentelemetry.propagators.composite import CompositePropagator
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from prometheus_client import start_http_server

METRICS_PORT = 9464
HEALTH_PREFIX = "/grpc.health.v1.Health/"


def setup() -> Callable[..., None]:
    """provider と propagator を設定し、/metrics を公開する。
    endpoint・service name は OTEL_* の env から SDK が解決する。返す関数で終了処理をする"""
    resource = Resource.create()
    # 終了処理は shutdown() の期限付きのものだけにする（atexit で再び詰まらせない）
    tracer_provider = TracerProvider(resource=resource, shutdown_on_exit=False)
    tracer_provider.add_span_processor(BatchSpanProcessor(OTLPSpanExporter()))

    readers = [PrometheusMetricReader()]
    if os.getenv("OTEL_METRICS_EXPORTER") == "otlp":
        readers.append(PeriodicExportingMetricReader(OTLPMetricExporter()))
    meter_provider = MeterProvider(
        resource=resource, metric_readers=readers, shutdown_on_exit=False
    )

    trace.set_tracer_provider(tracer_provider)
    metrics.set_meter_provider(meter_provider)
    propagate.set_global_textmap(
        CompositePropagator([TraceContextTextMapPropagator(), W3CBaggagePropagator()])
    )
    start_http_server(METRICS_PORT)

    def shutdown(timeout: float = 3.0) -> None:
        # collector が落ちているとエクスポートの再試行で止まらないので、期限で打ち切る
        def run() -> None:
            tracer_provider.shutdown()
            meter_provider.shutdown()

        t = threading.Thread(target=run, daemon=True)
        t.start()
        t.join(timeout)

    return shutdown


class JSONFormatter(logging.Formatter):
    """Go と同じキー（time level msg service trace_id span_id）で 1 行の JSON を書く"""

    service = os.getenv("OTEL_SERVICE_NAME", "")

    def format(self, record: logging.LogRecord) -> str:
        ctx = trace.get_current_span().get_span_context()
        out = {
            "time": datetime.fromtimestamp(record.created, timezone.utc).isoformat(),
            "level": record.levelname.lower(),
            "msg": record.getMessage(),
            "service": self.service,
            "trace_id": f"{ctx.trace_id:032x}" if ctx.is_valid else "",
            "span_id": f"{ctx.span_id:016x}" if ctx.is_valid else "",
        }
        out.update(getattr(record, "fields", {}))
        if record.exc_info:
            out["err"] = self.formatException(record.exc_info)
        return json.dumps(out, separators=(",", ":"))


def setup_logging() -> None:
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(JSONFormatter())
    logging.basicConfig(level=logging.INFO, handlers=[handler], force=True)


def go_code(code: grpc.StatusCode) -> str:
    """Go の codes.Code.String() と同じ表記にする（NOT_FOUND → NotFound）"""
    special = {grpc.StatusCode.OK: "OK", grpc.StatusCode.CANCELLED: "Canceled"}
    return special.get(code) or "".join(w.capitalize() for w in code.name.split("_"))


def client_channel(addr: str) -> grpc.Channel:
    """trace context を引き継ぐ client の channel。平文なのは、クラスタ内の暗号化を Istio の mTLS に任せるため"""
    return intercept_channel(grpc.insecure_channel(addr), client_interceptor())


def server_interceptors(log: logging.Logger) -> list[grpc.ServerInterceptor]:
    """計装（スパン）を外側、リクエストログを内側に置く。ログに trace_id を入れるため"""
    tracing = server_interceptor(filter_=filters.negate(filters.health_check()))
    return [tracing, RequestLog(log)]


class RequestLog(grpc.ServerInterceptor):
    """リクエストごとに rpc.method と code を 1 行出す（health check は除く）"""

    def __init__(self, log: logging.Logger):
        self.log = log

    def intercept_service(self, continuation, details):
        handler = continuation(details)
        if details.method.startswith(HEALTH_PREFIX):
            return handler
        if handler is None:
            fields = {"rpc.method": details.method, "code": "Unimplemented"}
            self.log.info("rpc", extra={"fields": fields})
            return None
        if handler.unary_unary is None:
            return handler
        inner = handler.unary_unary

        def logged(request, context):
            code = grpc.StatusCode.UNKNOWN
            try:
                resp = inner(request, context)
                code = context.code() or grpc.StatusCode.OK
                return resp
            except Exception:
                code = context.code() or grpc.StatusCode.UNKNOWN
                raise
            finally:
                fields = {"rpc.method": details.method, "code": go_code(code)}
                self.log.info("rpc", extra={"fields": fields})

        return grpc.unary_unary_rpc_method_handler(
            logged,
            request_deserializer=handler.request_deserializer,
            response_serializer=handler.response_serializer,
        )
