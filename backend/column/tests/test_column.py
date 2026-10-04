import uuid

import pytest
from sqlalchemy import create_engine
from sqlalchemy.pool import StaticPool

from internal.column.model.column import Base
from internal.column.repository.column import ColumnRepository
from internal.column.service.column import ColumnService, InvalidArgument, NotFound
from pkg.db.db import make_sessions


class FakeProjects:
    def __init__(self):
        self.allowed = set()

    def check_access(self, project_id):
        return project_id in self.allowed


@pytest.fixture
def projects():
    return FakeProjects()


@pytest.fixture
def service(projects):
    engine = create_engine(
        "sqlite://", connect_args={"check_same_thread": False}, poolclass=StaticPool
    )
    Base.metadata.create_all(engine)
    return ColumnService(ColumnRepository(make_sessions(engine)), projects)


def test_create_get_list_update_delete(service, projects):
    board, other = uuid.uuid4(), uuid.uuid4()
    projects.allowed |= {board, other}
    board = str(board)
    created = service.create("todo", board)
    assert service.get(str(created.id)).name == "todo"

    service.create("other board", str(other))
    service.create("no board", "")
    columns, total = service.list(board, 1, 100)
    assert total == 1 and [c.id for c in columns] == [created.id]

    updated = service.update(str(created.id), "doing", "")
    assert updated.name == "doing" and str(updated.board_id) == board

    columns, total = service.list("", 1, 100)
    assert total == 1 and columns[0].name == "no board"

    service.delete(str(created.id))
    with pytest.raises(NotFound):
        service.get(str(created.id))
    with pytest.raises(NotFound):
        service.delete(str(created.id))


def test_hides_columns_of_others_boards(service, projects):
    mine, theirs = uuid.uuid4(), uuid.uuid4()
    projects.allowed |= {mine, theirs}
    column = service.create("todo", str(theirs))
    projects.allowed.discard(theirs)

    for call in (
        lambda: service.create("x", str(theirs)),
        lambda: service.get(str(column.id)),
        lambda: service.update(str(column.id), "x", ""),
        lambda: service.delete(str(column.id)),
        lambda: service.list(str(theirs), 1, 100),
    ):
        with pytest.raises(NotFound):
            call()
    mine_column = service.create("todo", str(mine))
    with pytest.raises(NotFound):
        service.update(str(mine_column.id), "", str(theirs))


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

    from pkg.userid.userid import RequireUserID, current

    def call(method, metadata):
        h = grpc.unary_unary_rpc_method_handler(lambda req, ctx: current.get("ok"))
        interceptor = RequireUserID("/column.ColumnService/")
        wrapped = interceptor.intercept_service(lambda _: h, Details(method))
        return wrapped.unary_unary(None, FakeContext(metadata))

    m = "/column.ColumnService/ListColumns"
    user = "00000000-0000-4000-8000-000000000071"
    assert call(m, [("x-user-id", user)]) == user
    zero = [("x-user-id", "00000000-0000-0000-0000-000000000000")]
    for bad in ([], [("x-user-id", "not-a-uuid")], zero):
        with pytest.raises(PermissionError) as e:
            call(m, bad)
        assert e.value.args[0] == grpc.StatusCode.UNAUTHENTICATED
    assert call("/grpc.health.v1.Health/Check", []) == "ok"


def test_handler_keeps_status_of_project_errors():
    import grpc

    from internal.column.handler.column import grpc_errors

    class Unavailable(grpc.RpcError):
        def code(self):
            return grpc.StatusCode.UNAVAILABLE

        def details(self):
            return "project is down"

    @grpc_errors
    def call(self, request, context):
        raise Unavailable()

    with pytest.raises(PermissionError) as e:
        call(None, None, FakeContext([]))
    assert e.value.args[0] == grpc.StatusCode.UNAVAILABLE


def test_lists_columns_in_creation_order(service, projects):
    board = uuid.uuid4()
    projects.allowed.add(board)
    for name in ("todo", "doing", "done"):
        service.create(name, str(board))
    columns, _ = service.list(str(board), 0, 0)
    assert [c.name for c in columns] == ["todo", "doing", "done"]
