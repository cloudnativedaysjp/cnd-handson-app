import functools

import grpc

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
