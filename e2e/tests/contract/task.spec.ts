// Contract for handson-task (#105), using the existing task.proto.
import { expect, test } from "@playwright/test";
import { cfg } from "./lib/config";
import { client, grpc, OTHER_USER_ID, USER_ID } from "./lib/grpc";

test("create, get and list a task", async () => {
  const call = client("task/task.proto", "task", "TaskService", cfg.taskGrpc);
  const { task } = await call("CreateTask", { title: "contract task", description: "d", status: "todo" });
  expect(task.id).toBeTruthy();
  expect((await call("GetTask", { id: task.id })).task).toMatchObject({ id: task.id, title: "contract task" });
  const { tasks } = await call("ListTasks", { page: 1, page_size: 100 });
  expect(tasks.map((t: { id: string }) => t.id)).toContain(task.id);
});

test("only the project owner can see its tasks", async () => {
  const call = client("task/task.proto", "task", "TaskService", cfg.taskGrpc);
  const p = client("project/project.proto", "project", "ProjectService", cfg.projectGrpc);
  const { project } = await p("CreateProject", { name: "private", owner_id: USER_ID });
  const { task } = await call("CreateTask", { title: "private", status: "todo", project_id: project.id });
  const notFound = { code: grpc.status.NOT_FOUND };
  await expect(call("GetTask", { id: task.id }, OTHER_USER_ID)).rejects.toMatchObject(notFound);
  await expect(call("ListTasks", { project_id: project.id }, OTHER_USER_ID)).rejects.toMatchObject(notFound);
});
