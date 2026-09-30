from concurrent import futures
import logging
import grpc
import os
import argparse  # argparseをインポート
from grpc_health.v1 import health
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

from internal.role.model.role import RoleModel
from internal.role.handler.role import RoleHandler
from internal.role.service.role import RoleService
from internal.role.repository.role import RoleRepository
from pkg.db.db import Database
from role import role_pb2_grpc


def configure_health_server(server: grpc.Server):
    # ヘルスサーバーの作成（非ブロッキング等のオプションは不要）
    health_servicer = health.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)

    # サーバー全体のデフォルトの状態も設定可能（必要なら）
    health_servicer.set("", health_pb2.HealthCheckResponse.SERVING)


def server():
    """gRPCサーバーの起動
    Args:
        None

    """

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    db = Database(db_url, [RoleModel])
    # DI
    with db.session_scope() as session:
        role_repository = RoleRepository(session)
        role_service = RoleService(role_repository)
        role_handler = RoleHandler(role_service)
        role_pb2_grpc.add_RoleServiceServicer_to_server(role_handler, server)
    # gRPCサーバーにヘルスチェックを追加
    configure_health_server(server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


def migrate():
    """データベースのマイグレーション
    Args:
        None
    """
    db = Database(db_url, [RoleModel])  # SQLiteのURLを指定
    db.migrate()
    db.close()


def reset():
    """データベースのリセット
    Args:
        None
    """
    db = Database(db_url, [RoleModel])  # SQLiteのURLを指定
    db.drop_all()
    db.migrate()
    db.close()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)

    port = os.getenv("PORT", "50051")
    db_url = os.getenv(
        "DATABASE_URL",
        "postgresql+psycopg2://your_db_user:your_db_password@db:5432/your_db_name",
    )

    # argparseでコマンドライン引数を処理
    parser = argparse.ArgumentParser(description="Manage the gRPC server and database.")
    parser.add_argument(
        "command", choices=["server", "migrate", "reset"], help="Command to execute"
    )

    args = parser.parse_args()

    if args.command == "server":
        server()
    elif args.command == "migrate":
        migrate()
    elif args.command == "reset":
        reset()
    else:
        print("Invalid command. Use 'server', 'migrate', or 'reset'.")
        exit(1)
