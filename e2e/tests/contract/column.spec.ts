// Contract for handson-column (#109, Python). The column proto keys columns by board_id; the
// contract treats a project id as the board id. project -> column listing needs a new RPC
// (proposed/project_column.proto), so that test is fixme.
import { randomUUID } from "node:crypto";
import { expect, test } from "@playwright/test";
import { cfg } from "./lib/config";
import { client, USER_ID } from "./lib/grpc";

const column = () => client("column/column.proto", "column", "ColumnService", cfg.columnGrpc);

test("create, get and list columns", async () => {
  const call = column();
  const board_id = randomUUID();
  const { column: created } = await call("CreateColumn", { name: "todo", board_id });
  expect(created.id).toBeTruthy();
  expect((await call("GetColumn", { id: created.id })).column).toMatchObject({ id: created.id, name: "todo", board_id });
  const { columns } = await call("ListColumns", { board_id, page: 1, page_size: 100 });
  expect(columns.map((c: { id: string }) => c.id)).toEqual([created.id]);
});

// fixme: needs ProjectService.ListProjectColumns (proposed/project_column.proto), see #72/#109.
test.fixme("list columns of a project via project -> column", async () => {
  const p = client("project/project.proto", "project", "ProjectService", cfg.projectGrpc);
  const list = client("project_column.proto", "project", "ProjectService", cfg.projectGrpc, true);
  const { project } = await p("CreateProject", { name: "with columns", owner_id: USER_ID });
  const { column: created } = await column()("CreateColumn", { name: "todo", board_id: project.id });
  const { columns } = await list("ListProjectColumns", { project_id: project.id });
  expect(columns.map((c: { id: string }) => c.id)).toEqual([created.id]);
});
