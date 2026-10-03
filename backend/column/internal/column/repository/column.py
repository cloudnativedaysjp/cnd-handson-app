from typing import Optional
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.orm import sessionmaker

from internal.column.model.column import ColumnModel


class ColumnRepository:
    """リクエストごとにセッションを開く（gRPC はスレッドプールで並行に呼ぶため共有しない）"""

    def __init__(self, sessions: sessionmaker):
        self.sessions = sessions

    def create(self, name: str, board_id: Optional[UUID]) -> ColumnModel:
        with self.sessions.begin() as s:
            column = ColumnModel(name=name, board_id=board_id)
            s.add(column)
        return column

    def get(self, column_id: UUID) -> Optional[ColumnModel]:
        with self.sessions() as s:
            return s.get(ColumnModel, column_id)

    def update(
        self, column_id: UUID, name: str, board_id: Optional[UUID]
    ) -> Optional[ColumnModel]:
        with self.sessions.begin() as s:
            column = s.get(ColumnModel, column_id)
            if column is None:
                return None
            if name:
                column.name = name
            if board_id is not None:
                column.board_id = board_id
        return column

    def delete(self, column_id: UUID) -> bool:
        with self.sessions.begin() as s:
            column = s.get(ColumnModel, column_id)
            if column is None:
                return False
            s.delete(column)
        return True

    def list(
        self, board_id: Optional[UUID], page: int, page_size: int
    ) -> tuple[list[ColumnModel], int]:
        """board_id が None なら board の無い column だけ。page / page_size が 0 以下ならページングしない"""
        query = select(ColumnModel).where(ColumnModel.board_id == board_id)
        with self.sessions() as s:
            total = s.scalar(select(func.count()).select_from(query.subquery()))
            if page > 0 and page_size > 0:
                query = query.offset((page - 1) * page_size).limit(page_size)
            return list(s.scalars(query.order_by(ColumnModel.name))), total
