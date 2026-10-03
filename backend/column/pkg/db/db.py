import os

from sqlalchemy import Engine, create_engine
from sqlalchemy.engine import URL
from sqlalchemy.orm import sessionmaker


def build_db_url() -> str:
    """環境変数 DB_* から接続URLを組み立てる。未設定があれば起動時にエラー"""
    keys = ["DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_DB"]
    missing = [k for k in keys if not os.environ.get(k)]
    if missing:
        raise SystemExit(
            f"required environment variables are not set: {', '.join(missing)}"
        )
    return URL.create(
        "postgresql+psycopg2",
        username=os.environ["DB_USER"],
        password=os.environ["DB_PASSWORD"],
        host=os.environ["DB_HOST"],
        port=int(os.environ["DB_PORT"]),
        database=os.environ["DB_DB"],
    ).render_as_string(hide_password=False)


def make_sessions(engine: Engine) -> sessionmaker:
    # 返したモデルをセッションの外でも読めるよう、commit 後に期限切れにしない
    return sessionmaker(bind=engine, expire_on_commit=False)


def make_engine(url: str) -> Engine:
    return create_engine(url, pool_pre_ping=True)
