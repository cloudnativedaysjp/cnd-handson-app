// Contract for handson-column (#109, Python). The column proto keys columns by board_id; the
// contract treats a project id as the board id, and only the project's owner may use it (#184).
import { expect, test } from "@playwright/test";
import { cfg } from "./lib/config";
import { client, grpc, OTHER_USER_ID, USER_ID } from "./lib/grpc";

const column = () => client("column/column.proto", "column", "ColumnService", cfg.columnGrpc);
const project = () => client("project/project.proto", "project", "ProjectService", cfg.projectGrpc);

test("create, get and list columns", async () => {
  const call = column();
  const { project: board } = await project()("CreateProject", { name: "board", owner_id: USER_ID });
  const board_id = board.id;
  const { column: created } = await call("CreateColumn", { name: "todo", board_id });
  expect(created.id).toBeTruthy();
  expect((await call("GetColumn", { id: created.id })).column).toMatchObject({ id: created.id, name: "todo", board_id });
  const { columns } = await call("ListColumns", { board_id, page: 1, page_size: 100 });
  expect(columns.map((c: { id: string }) => c.id)).toEqual([created.id]);
});

test("list columns of a project via project -> column", async () => {
  const p = client("project/project.proto", "project", "ProjectService", cfg.projectGrpc);
  const { project } = await p("CreateProject", { name: "with columns", owner_id: USER_ID });
  const { column: created } = await column()("CreateColumn", { name: "todo", board_id: project.id });
  const { columns } = await p("ListProjectColumns", { project_id: project.id });
  expect(columns.map((c: { id: string }) => c.id)).toEqual([created.id]);
});

test("only the project owner can see its columns", async () => {
  const call = column();
  const { project: board } = await project()("CreateProject", { name: "private", owner_id: USER_ID });
  const { column: created } = await call("CreateColumn", { name: "todo", board_id: board.id });
  const notFound = { code: grpc.status.NOT_FOUND };
  await expect(call("GetColumn", { id: created.id }, OTHER_USER_ID)).rejects.toMatchObject(notFound);
  await expect(call("ListColumns", { board_id: board.id }, OTHER_USER_ID)).rejects.toMatchObject(notFound);
});
