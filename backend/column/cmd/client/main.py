import os

import grpc
import argparse  # argparseをインポート
from typing import Any

from column import column_pb2_grpc, column_pb2


def create_Request(
    stub: column_pb2_grpc.ColumnServiceStub, command: str, arguments: list[str]
) -> Any:
    """コマンドに応じたリクエストを生成する関数
    Args:
        command: コマンド名
    """
    if command == "create":
        response = stub.CreateColumn(
            column_pb2.CreateColumnRequest(name=arguments[0], board_id=arguments[1])
        )
        return (
            f"id={response.column.id}, name={response.column.name}, "
            f"board_id={response.column.board_id}"
        )
    elif command == "get":
        response = stub.GetColumn(column_pb2.GetColumnRequest(id=arguments[0]))
        return (
            f"id={response.column.id}, name={response.column.name}, "
            f"board_id={response.column.board_id}"
        )

    elif command == "update":
        response = stub.UpdateColumn(
            column_pb2.UpdateColumnRequest(
                id=arguments[0], name=arguments[1], board_id=arguments[2]
            )
        )
        return (
            f"id={response.column.id}, name={response.column.name}, "
            f"board_id={response.column.board_id}"
        )

    elif command == "delete":
        response = stub.DeleteColumn(column_pb2.DeleteColumnRequest(id=arguments[0]))
        if response.success:
            return f"Column with id {arguments[0]} deleted successfully"
        else:
            return f"Failed to delete column with id {arguments[0]}"
    elif command == "list":
        return stub.ListColumns(
            column_pb2.ListColumnsRequest(
                board_id=arguments[0],
                page=int(arguments[1]),
                page_size=int(arguments[2]),
            )
        )
    else:
        raise ValueError("Invalid command")


class UserID(grpc.UnaryUnaryClientInterceptor):
    """column は x-user-id を必須にしている。既定は契約テストの demo ユーザー"""

    user_id = os.getenv("COLUMN_USER_ID", "00000000-0000-4000-8000-000000000071")

    def intercept_unary_unary(self, continuation, details, request):
        metadata = [*(details.metadata or []), ("x-user-id", self.user_id)]
        return continuation(details._replace(metadata=metadata), request)


def run(server_address: str, command: str, arguments: list[str]):
    # サーバーのアドレスとポートを指定
    # gRPCチャンネルを作成
    with grpc.insecure_channel(server_address) as channel:
        # スタブを作成
        stub = column_pb2_grpc.ColumnServiceStub(
            grpc.intercept_channel(channel, UserID())
        )
        # コマンドに応じたリクエストを生成
        response = create_Request(stub, command, arguments)

        # レスポンスを表示
        print("Response from server:", response)


if __name__ == "__main__":
    # コマンドライン引数のパーサーを作成
    parser = argparse.ArgumentParser(description="gRPC Client for Column Service")
    parser.add_argument(
        "--port", type=int, default=50051, help="Port number to connect to the server"
    )
    parser.add_argument(
        "command",
        choices=["create", "get", "update", "delete", "list"],
        help="Command to execute",
    )
    parser.add_argument("arguments", nargs="+", help="Arguments for the command")

    args = parser.parse_args()

    # コマンドに応じた引数を取得
    server_address = f"localhost:{args.port}"
    command = args.command
    arguments = args.arguments

    # コマンドを実行
    run(server_address, command, arguments)
