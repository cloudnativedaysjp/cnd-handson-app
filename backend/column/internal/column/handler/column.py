from column import column_pb2_grpc, column_pb2
import grpc


class ColumnHandler(column_pb2_grpc.ColumnServiceServicer):
    """ColumnHandlerの実装"""

    def __init__(self, column_service):
        """ColumnHandlerの初期化
        Args:
            column_service: ColumnServiceのインスタンス
        """
        self.column_service = column_service

    def CreateColumn(self, request, context):
        """Columnの生成
        Args:
        request: CreateColumnRequest
            Columnの生成に必要な情報を含むリクエスト
        context: grpc.ServicerContext
        """
        try:
            # ColumnService から生成された ColumnModel を取得
            column_model = self.column_service.create(request.name, request.board_id)

            # ColumnResponse を構築
            response = column_pb2.ColumnResponse(
                column=column_pb2.Column(
                    id=str(column_model.id),
                    name=column_model.name,
                    board_id=str(column_model.board_id),
                )
            )
            return response
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error creating column: {str(e)}")
            return column_pb2.ColumnResponse()

    def GetColumn(self, request, context):
        """Columnの取得"""
        try:
            # ColumnService から ColumnModel を取得
            column_model = self.column_service.get(column_id=request.id)
            if column_model is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Column not found!")
                return column_pb2.ColumnResponse()

            # ColumnResponse を構築
            response = column_pb2.ColumnResponse(
                column=column_pb2.Column(
                    id=str(column_model.id),
                    name=column_model.name,
                    board_id=str(column_model.board_id),
                )
            )
            return response
        except Exception as e:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Column not found: {str(e)}")
            return column_pb2.ColumnResponse()

    def UpdateColumn(self, request, context):
        """Columnの更新"""
        try:

            # ここでColumnを更新する処理を実装する
            column_model = self.column_service.update(
                request.id, request.name, request.board_id
            )
            if column_model is None:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("Column not found!")
                return column_pb2.ColumnResponse()

            # ColumnResponse を構築
            response = column_pb2.ColumnResponse(
                column=column_pb2.Column(
                    id=str(column_model.id),
                    name=column_model.name,
                    board_id=str(column_model.board_id),
                )
            )
            return response
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error updating column: {str(e)}")
            return column_pb2.ColumnResponse()

    def DeleteColumn(self, request, context):
        """Columnの削除"""
        # ここでColumnを削除する処理を実装する
        success = self.column_service.delete(request.id)
        if not success:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("Column not found!")
            return column_pb2.DeleteColumnResponse(success=False)
        response = column_pb2.DeleteColumnResponse(success=True)
        return response

    def ListColumns(self, request, context):
        """Columnの一覧取得"""
        # ここでColumnの一覧を取得する処理を実装する
        try:
            ColumnModels = self.column_service.list(
                request.board_id, request.page, request.page_size
            )
            columns = []
            for column_model in ColumnModels:
                columns.append(
                    column_pb2.Column(
                        id=str(column_model.id),
                        name=column_model.name,
                        board_id=str(column_model.board_id),
                    )
                )
            total_count = self.column_service.total_count(request.board_id)
            return column_pb2.ListColumnsResponse(
                columns=columns, total_count=total_count
            )
        except Exception as e:
            context.abort(grpc.StatusCode.INTERNAL, f"Error listing columns: {e}")
