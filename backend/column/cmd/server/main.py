import argparse
import logging
import os
from concurrent import futures

import grpc
from grpc_health.v1 import health, health_pb2, health_pb2_grpc

from column import column_pb2_grpc
from internal.column.handler.column import ColumnHandler
from internal.column.model.column import Base
from internal.column.repository.column import ColumnRepository
from internal.column.service.column import ColumnService
from pkg.db.db import build_db_url, make_engine, make_sessions


def serve() -> None:
    engine = make_engine(build_db_url())
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service = ColumnService(ColumnRepository(make_sessions(engine)))
    column_pb2_grpc.add_ColumnServiceServicer_to_server(ColumnHandler(service), server)

    health_servicer = health.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)
    health_servicer.set("", health_pb2.HealthCheckResponse.SERVING)

    port = os.getenv("PORT", "50051")
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    logging.info("gRPC server listening on port %s", port)
    server.wait_for_termination()


def migrate() -> None:
    engine = make_engine(build_db_url())
    Base.metadata.create_all(engine)
    engine.dispose()
    logging.info("migration completed")


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    parser = argparse.ArgumentParser(description="column service")
    parser.add_argument("command", choices=["server", "migrate"])
    {"server": serve, "migrate": migrate}[parser.parse_args().command]()
