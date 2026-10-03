import uuid

import pytest
from sqlalchemy import create_engine
from sqlalchemy.pool import StaticPool

from internal.column.model.column import Base
from internal.column.repository.column import ColumnRepository
from internal.column.service.column import ColumnService, InvalidArgument, NotFound
from pkg.db.db import make_sessions


@pytest.fixture
def service():
    engine = create_engine(
        "sqlite://", connect_args={"check_same_thread": False}, poolclass=StaticPool
    )
    Base.metadata.create_all(engine)
    return ColumnService(ColumnRepository(make_sessions(engine)))


def test_create_get_list_update_delete(service):
    board = str(uuid.uuid4())
    created = service.create("todo", board)
    assert service.get(str(created.id)).name == "todo"

    service.create("other board", str(uuid.uuid4()))
    columns, total = service.list(board, 1, 100)
    assert total == 1 and [c.id for c in columns] == [created.id]

    updated = service.update(str(created.id), "doing", "")
    assert updated.name == "doing" and str(updated.board_id) == board

    service.delete(str(created.id))
    with pytest.raises(NotFound):
        service.get(str(created.id))
    with pytest.raises(NotFound):
        service.delete(str(created.id))


def test_rejects_bad_input(service):
    with pytest.raises(InvalidArgument):
        service.create("", "")
    with pytest.raises(InvalidArgument):
        service.create("todo", "not-a-uuid")
    with pytest.raises(InvalidArgument):
        service.get("not-a-uuid")
