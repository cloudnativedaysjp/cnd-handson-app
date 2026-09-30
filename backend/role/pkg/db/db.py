import os

from sqlalchemy import create_engine, Engine
from sqlalchemy.engine import URL
from sqlalchemy.orm import sessionmaker, Session
from contextlib import contextmanager


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


class Database:
    def __init__(self, db_url: str, models: list) -> None:
        """Databaseの初期化
        Args:
            db_url: データベースのURL
            models: SQLAlchemyのモデルクラスのリスト
        """
        self.engine = create_engine(db_url, echo=True)  # echo=TrueはSQLログ出したいなら
        self.SessionLocal = sessionmaker(
            bind=self.engine, autoflush=False, autocommit=False
        )
        self.models = models

    def get_engine(self) -> Engine:
        """SQLAlchemyエンジンの取得"""
        return self.engine

    def get_session(self) -> Session:
        """SQLAlchemyセッションの取得"""
        return self.SessionLocal()

    def migrate(self) -> None:
        """データベースのマイグレーション"""
        for model in self.models:
            model.__table__.create(self.engine, checkfirst=True)

        return None

    def drop_all(self) -> None:
        """データベースの全テーブル削除"""
        for model in self.models:
            model.__table__.drop(self.engine, checkfirst=True)

        return None

    def close(self) -> None:
        """データベース接続のクローズ"""
        self.engine.dispose()
        return None

    @contextmanager
    def session_scope(self) -> Session:
        """with構文でセッション管理する"""
        session = self.get_session()
        try:
            yield session
            session.commit()
        except Exception:
            session.rollback()
            raise
        finally:
            session.close()
