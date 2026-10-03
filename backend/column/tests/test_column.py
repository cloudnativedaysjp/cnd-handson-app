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


class FakeContext:
    def __init__(self, metadata):
        self.metadata = metadata

    def invocation_metadata(self):
        return self.metadata

    def abort(self, code, details):
        raise PermissionError(code)


class Details:
    def __init__(self, method):
        self.method = method


def test_require_user_id():
    import grpc

    from internal.column.handler.column import RequireUserID

    def call(method, metadata):
        h = grpc.unary_unary_rpc_method_handler(lambda req, ctx: "ok")
        wrapped = RequireUserID().intercept_service(lambda _: h, Details(method))
        return wrapped.unary_unary(None, FakeContext(metadata))

    m = "/column.ColumnService/ListColumns"
    assert call(m, [("x-user-id", "00000000-0000-4000-8000-000000000071")]) == "ok"
    for bad in ([], [("x-user-id", "not-a-uuid")]):
        with pytest.raises(PermissionError) as e:
            call(m, bad)
        assert e.value.args[0] == grpc.StatusCode.UNAUTHENTICATED
    assert call("/grpc.health.v1.Health/Check", []) == "ok"
