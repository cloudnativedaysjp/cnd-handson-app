import grpc
import argparse  # argparseをインポート
from typing import Any

from role import role_pb2_grpc, role_pb2


def create_Request(
    stub: role_pb2_grpc.RoleServiceStub, command: str, arguments: list[str]
) -> Any:
    """コマンドに応じたリクエストを生成する関数
    Args:
        command: コマンド名
    """
    if command == "create":
        response = stub.CreateRole(
            role_pb2.CreateRoleRequest(name=arguments[0], description=arguments[1])
        )
        return (
            f"id={response.role.id}, name={response.role.name}, "
            f"description={response.role.description}"
        )
    elif command == "get":
        response = stub.GetRole(role_pb2.GetRoleRequest(id=arguments[0]))
        return (
            f"id={response.role.id}, name={response.role.name}, "
            f"description={response.role.description}"
        )

    elif command == "update":
        response = stub.UpdateRole(
            role_pb2.UpdateRoleRequest(
                id=arguments[0], name=arguments[1], description=arguments[2]
            )
        )
        return (
            f"id={response.role.id}, name={response.role.name}, "
            f"description={response.role.description}"
        )

    elif command == "delete":
        response = stub.DeleteRole(role_pb2.DeleteRoleRequest(id=arguments[0]))
        if response.success:
            return f"Role with id {arguments[0]} deleted successfully"
        else:
            return f"Failed to delete role with id {arguments[0]}"
    elif command == "list":
        return stub.ListRoles()
    else:
        raise ValueError("Invalid command")


def run(server_address: str, command: str, arguments: list[str]):
    # サーバーのアドレスとポートを指定
    # gRPCチャンネルを作成
    with grpc.insecure_channel(server_address) as channel:
        # スタブを作成
        stub = role_pb2_grpc.RoleServiceStub(channel)
        # コマンドに応じたリクエストを生成
        response = create_Request(stub, command, arguments)

        # レスポンスを表示
        print("Response from server:", response)


if __name__ == "__main__":
    # コマンドライン引数のパーサーを作成
    parser = argparse.ArgumentParser(description="gRPC Client for Role Service")
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
