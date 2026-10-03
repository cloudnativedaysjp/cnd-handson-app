"""gRPC metadata の x-user-id を受け取り、下流へ引き継ぐ（JWT の検証は入口が行う）。Go の pkg/userid と同じ約束"""

import uuid
from contextvars import ContextVar

import grpc
from opentelemetry import trace

KEY = "x-user-id"
current: ContextVar[str] = ContextVar("user_id")


class RequireUserID(grpc.ServerInterceptor):
    """service_prefix の呼び出しに UUID の x-user-id を求め、current に入れてスパンに user.id を付ける。
    health check は metadata を送らないので prefix で外す"""

    def __init__(self, service_prefix: str):
        self.service_prefix = service_prefix

    def intercept_service(self, continuation, details):
        handler = continuation(details)
        if handler is None or not details.method.startswith(self.service_prefix):
            return handler
        inner = handler.unary_unary

        def checked(request, context):
            values = [v for k, v in context.invocation_metadata() if k == KEY]
            if len(values) != 1:
                context.abort(grpc.StatusCode.UNAUTHENTICATED, "x-user-id is required")
            try:
                user_id = uuid.UUID(values[0])
            except ValueError:
                user_id = uuid.UUID(int=0)
            # ゼロの UUID は「未指定」と区別できないので受け付けない
            if user_id.int == 0:
                context.abort(
                    grpc.StatusCode.UNAUTHENTICATED,
                    "x-user-id must be a non-zero UUID",
                )
            trace.get_current_span().set_attribute("user.id", str(user_id))
            token = current.set(str(user_id))
            try:
                return inner(request, context)
            finally:
                current.reset(token)

        return grpc.unary_unary_rpc_method_handler(
            checked,
            request_deserializer=handler.request_deserializer,
            response_serializer=handler.response_serializer,
        )


def forward_metadata() -> list[tuple[str, str]]:
    """RequireUserID で受けた x-user-id を、下流の呼び出しに付ける metadata にする"""
    return [(KEY, current.get())]
