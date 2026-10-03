from typing import Optional
from uuid import UUID

from internal.column.model.column import ColumnModel
from internal.column.repository.column import ColumnRepository


class InvalidArgument(ValueError):
    pass


class NotFound(LookupError):
    pass


def parse_id(value: str) -> UUID:
    try:
        return UUID(value)
    except ValueError:
        raise InvalidArgument(f"invalid id: {value!r}") from None


def parse_optional_id(value: str) -> Optional[UUID]:
    """空文字は未指定（None）として扱う"""
    return parse_id(value) if value else None


class ColumnService:
    def __init__(self, repository: ColumnRepository) -> None:
        self.repository = repository

    def create(self, name: str, board_id: str) -> ColumnModel:
        if not name:
            raise InvalidArgument("name is required")
        return self.repository.create(name, parse_optional_id(board_id))

    def get(self, column_id: str) -> ColumnModel:
        column = self.repository.get(parse_id(column_id))
        if column is None:
            raise NotFound("column not found")
        return column

    def update(self, column_id: str, name: str, board_id: str) -> ColumnModel:
        column = self.repository.update(
            parse_id(column_id), name, parse_optional_id(board_id)
        )
        if column is None:
            raise NotFound("column not found")
        return column

    def delete(self, column_id: str) -> None:
        if not self.repository.delete(parse_id(column_id)):
            raise NotFound("column not found")

    def list(
        self, board_id: str, page: int, page_size: int
    ) -> tuple[list[ColumnModel], int]:
        return self.repository.list(parse_optional_id(board_id), page, page_size)
