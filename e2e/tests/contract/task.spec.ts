// Contract for handson-task (#105), using the existing task.proto.
import { expect, test } from "@playwright/test";
import { cfg } from "./lib/config";
import { client } from "./lib/grpc";

test("create, get and list a task", async () => {
  const call = client("task/task.proto", "task", "TaskService", cfg.taskGrpc);
  const { task } = await call("CreateTask", { title: "contract task", description: "d", status: "todo" });
  expect(task.id).toBeTruthy();
  expect((await call("GetTask", { id: task.id })).task).toMatchObject({ id: task.id, title: "contract task" });
  const { tasks } = await call("ListTasks", { page: 1, page_size: 100 });
  expect(tasks.map((t: { id: string }) => t.id)).toContain(task.id);
});
