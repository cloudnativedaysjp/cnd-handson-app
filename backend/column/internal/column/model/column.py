import uuid
from datetime import datetime, timezone

from sqlalchemy import VARCHAR, Column, DateTime, Uuid
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class ColumnModel(Base):
    __tablename__ = "columns"

    id = Column(Uuid, primary_key=True, default=uuid.uuid4)
    name = Column(VARCHAR(255), nullable=False)
    # project の ID を入れる（proto のフィールド名は互換のため board_id のまま）
    board_id = Column(Uuid, index=True)
    # 列は作った順に並べる
    created_at = Column(
        DateTime(timezone=True),
        nullable=False,
        default=lambda: datetime.now(timezone.utc),
    )
