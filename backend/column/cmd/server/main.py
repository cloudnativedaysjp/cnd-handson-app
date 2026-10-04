import argparse
import logging
import os
import signal
import threading
from concurrent import futures

import grpc
from grpc_health.v1 import health, health_pb2, health_pb2_grpc
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor
from sqlalchemy import text

from column import column_pb2_grpc
from internal.column.handler.column import ColumnHandler
from internal.column.model.column import Base
from internal.column.repository.column import ColumnRepository
from internal.column.repository.project import ProjectRepository
from internal.column.service.column import ColumnService
from pkg.db.db import build_db_url, make_engine, make_sessions
from pkg.telemetry import telemetry
from pkg.userid.userid import RequireUserID

log = logging.getLogger("column")


def serve() -> None:
    shutdown_telemetry = telemetry.setup()
    engine = make_engine(build_db_url())
    SQLAlchemyInstrumentor().instrument(engine=engine)

    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10),
        interceptors=[
            *telemetry.server_interceptors(log),
            RequireUserID("/column.ColumnService/"),
        ],
    )
    projects = ProjectRepository(telemetry.client_channel(os.environ["PROJECT_ADDR"]))
    service = ColumnService(ColumnRepository(make_sessions(engine)), projects)
    column_pb2_grpc.add_ColumnServiceServicer_to_server(ColumnHandler(service), server)
    health_servicer = health.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)
    health_servicer.set("", health_pb2.HealthCheckResponse.SERVING)

    port = os.getenv("PORT", "50051")
    server.add_insecure_port(f"[::]:{port}")
    stop = threading.Event()
    for sig in (signal.SIGTERM, signal.SIGINT):
        signal.signal(sig, lambda *_: stop.set())
    server.start()
    log.info("listening", extra={"fields": {"grpc": f":{port}"}})

    stop.wait()
    # grace を過ぎても戻らないハンドラーで止まらないよう、待ちにも上限を付ける
    server.stop(grace=5).wait(timeout=6)
    shutdown_telemetry()


def migrate() -> None:
    engine = make_engine(build_db_url())
    Base.metadata.create_all(engine)
    # create_all は既存のテーブルに列を足さないので、前からあるボリュームには自分で足す
    with engine.begin() as conn:
        conn.execute(
            text(
                "ALTER TABLE columns ADD COLUMN IF NOT EXISTS created_at"
                " TIMESTAMPTZ NOT NULL DEFAULT now()"
            )
        )
    engine.dispose()
    log.info("migration completed")


if __name__ == "__main__":
    telemetry.setup_logging()
    parser = argparse.ArgumentParser(description="column service")
    parser.add_argument("command", choices=["server", "migrate"])
    {"server": serve, "migrate": migrate}[parser.parse_args().command]()
