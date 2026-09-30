from concurrent import futures
import logging
import grpc
import os
import argparse  # argparseをインポート
from grpc_health.v1 import health
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc


from column import column_pb2_grpc
from internal.column.model.column import ColumnModel
from internal.column.repository.column import ColumnRepository
from internal.column.service.column import ColumnService
from internal.column.handler.column import ColumnHandler
from pkg.db.db import Database


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

    db = Database(db_url, [ColumnModel])
    # DI
    with db.session_scope() as session:
        column_repository = ColumnRepository(session)
        column_service = ColumnService(column_repository)
        column_handler = ColumnHandler(column_service)
        column_pb2_grpc.add_ColumnServiceServicer_to_server(column_handler, server)

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
    db = Database(db_url, [ColumnModel])  # SQLiteのURLを指定
    db.migrate()
    db.close()


def reset():
    """データベースのリセット
    Args:
        None
    """
    db = Database(db_url, [ColumnModel])  # SQLiteのURLを指定
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
