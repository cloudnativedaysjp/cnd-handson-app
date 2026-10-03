import functools
import uuid

import grpc
from opentelemetry import trace

from column import column_pb2, column_pb2_grpc
from internal.column.service.column import ColumnService, InvalidArgument, NotFound


def to_proto(c) -> column_pb2.Column:
    return column_pb2.Column(
        id=str(c.id), name=c.name, board_id=str(c.board_id) if c.board_id else ""
    )


def grpc_errors(method):
    """service の例外を gRPC のステータスに変える"""

    @functools.wraps(method)
    def wrapper(self, request, context):
        try:
            return method(self, request, context)
        except InvalidArgument as e:
            context.abort(grpc.StatusCode.INVALID_ARGUMENT, str(e))
        except NotFound as e:
            context.abort(grpc.StatusCode.NOT_FOUND, str(e))

    return wrapper


class ColumnHandler(column_pb2_grpc.ColumnServiceServicer):
    def __init__(self, service: ColumnService):
        self.service = service

    @grpc_errors
    def CreateColumn(self, request, context):
        column = self.service.create(request.name, request.board_id)
        return column_pb2.ColumnResponse(column=to_proto(column))

    @grpc_errors
    def GetColumn(self, request, context):
        return column_pb2.ColumnResponse(column=to_proto(self.service.get(request.id)))

    @grpc_errors
    def UpdateColumn(self, request, context):
        column = self.service.update(request.id, request.name, request.board_id)
        return column_pb2.ColumnResponse(column=to_proto(column))

    @grpc_errors
    def DeleteColumn(self, request, context):
        self.service.delete(request.id)
        return column_pb2.DeleteColumnResponse(success=True)

    @grpc_errors
    def ListColumns(self, request, context):
        columns, total = self.service.list(
            request.board_id, request.page, request.page_size
        )
        return column_pb2.ListColumnsResponse(
            columns=[to_proto(c) for c in columns], total_count=total
        )


class RequireUserID(grpc.ServerInterceptor):
    """ColumnService の呼び出しに x-user-id（JWT は入口が検証済み）を求め、スパンに user.id を付ける。
    health check は metadata を送らないので対象外"""

    def intercept_service(self, continuation, details):
        handler = continuation(details)
        if handler is None or not details.method.startswith("/column.ColumnService/"):
            return handler
        inner = handler.unary_unary

        def checked(request, context):
            values = [v for k, v in context.invocation_metadata() if k == "x-user-id"]
            if len(values) != 1:
                context.abort(grpc.StatusCode.UNAUTHENTICATED, "x-user-id is required")
            try:
                user_id = str(uuid.UUID(values[0]))
            except ValueError:
                context.abort(
                    grpc.StatusCode.UNAUTHENTICATED, "x-user-id must be a UUID"
                )
            trace.get_current_span().set_attribute("user.id", user_id)
            return inner(request, context)

        return grpc.unary_unary_rpc_method_handler(
            checked,
            request_deserializer=handler.request_deserializer,
            response_serializer=handler.response_serializer,
        )
