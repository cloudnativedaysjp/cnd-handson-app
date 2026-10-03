from typing import Optional
from uuid import UUID

from internal.column.model.column import ColumnModel
from internal.column.repository.column import ColumnRepository
from internal.column.repository.project import ProjectRepository


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
    def __init__(
        self, repository: ColumnRepository, projects: ProjectRepository
    ) -> None:
        self.repository = repository
        self.projects = projects

    def can_use(self, board_id: Optional[UUID]) -> bool:
        """board はプロジェクトで、所有者だけが使える。所有者でなければ存在を明かさないよう NotFound にする"""
        return board_id is None or self.projects.check_access(board_id)

    def check_board(self, board_id: Optional[UUID]) -> None:
        if not self.can_use(board_id):
            raise NotFound("board not found")

    def create(self, name: str, board_id: str) -> ColumnModel:
        if not name:
            raise InvalidArgument("name is required")
        board = parse_optional_id(board_id)
        self.check_board(board)
        return self.repository.create(name, board)

    def get(self, column_id: str) -> ColumnModel:
        column = self.repository.get(parse_id(column_id))
        if column is None or not self.can_use(column.board_id):
            raise NotFound("column not found")
        return column

    def update(self, column_id: str, name: str, board_id: str) -> ColumnModel:
        current = self.get(column_id)
        board = parse_optional_id(board_id)
        if board != current.board_id:
            self.check_board(board)
        column = self.repository.update(current.id, name, board)
        if column is None:
            raise NotFound("column not found")
        return column

    def delete(self, column_id: str) -> None:
        if not self.repository.delete(self.get(column_id).id):
            raise NotFound("column not found")

    def list(
        self, board_id: str, page: int, page_size: int
    ) -> tuple[list[ColumnModel], int]:
        board = parse_optional_id(board_id)
        self.check_board(board)
        return self.repository.list(board, page, page_size)
