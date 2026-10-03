import json
import logging

import grpc
from opentelemetry.sdk.trace import TracerProvider

from pkg.telemetry import telemetry


def test_json_log_has_convention_keys():
    record = logging.LogRecord("t", logging.INFO, __file__, 1, "hello", None, None)
    record.fields = {"rpc.method": "/x"}
    with TracerProvider().get_tracer("t").start_as_current_span("op") as span:
        out = json.loads(telemetry.JSONFormatter().format(record))
    assert (
        out["level"] == "info" and out["msg"] == "hello" and out["rpc.method"] == "/x"
    )
    assert out["trace_id"] == f"{span.get_span_context().trace_id:032x}"
    assert {"time", "service", "span_id"} <= out.keys()


def test_go_code_matches_go_names():
    assert telemetry.go_code(grpc.StatusCode.OK) == "OK"
    assert telemetry.go_code(grpc.StatusCode.NOT_FOUND) == "NotFound"
    assert telemetry.go_code(grpc.StatusCode.CANCELLED) == "Canceled"


class FakeContext:
    def __init__(self, code=None):
        self._code = code

    def code(self):
        return self._code


class Details:
    def __init__(self, method):
        self.method = method


def test_request_log_writes_one_line_per_rpc(caplog):
    caplog.set_level(logging.INFO)
    interceptor = telemetry.RequestLog(logging.getLogger("t"))

    def call(method, ctx):
        h = grpc.unary_unary_rpc_method_handler(lambda req, c: "resp")
        return interceptor.intercept_service(lambda _: h, Details(method)).unary_unary(
            None, ctx
        )

    assert call("/column.ColumnService/GetColumn", FakeContext()) == "resp"
    call("/column.ColumnService/GetColumn", FakeContext(grpc.StatusCode.NOT_FOUND))
    call("/grpc.health.v1.Health/Check", FakeContext())
    fields = [r.fields for r in caplog.records]
    assert fields == [
        {"rpc.method": "/column.ColumnService/GetColumn", "code": "OK"},
        {"rpc.method": "/column.ColumnService/GetColumn", "code": "NotFound"},
    ]


def test_json_log_is_compact():
    record = logging.LogRecord("t", logging.INFO, __file__, 1, "m", None, None)
    assert '": "' not in telemetry.JSONFormatter().format(record)
