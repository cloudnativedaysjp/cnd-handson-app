// Contract for handson-project (#72).
import { expect, test } from "@playwright/test";
import { cfg } from "./lib/config";
import { client, USER_ID } from "./lib/grpc";

const project = () => client("project/project.proto", "project", "ProjectService", cfg.projectGrpc);

test("create and get a project", async () => {
  const call = project();
  const { project: created } = await call("CreateProject", {
    name: "contract project",
    description: "d",
    owner_id: USER_ID,
  });
  expect(created.id).toBeTruthy();
  const { project: got } = await call("GetProject", { id: created.id });
  expect(got).toMatchObject({ id: created.id, name: "contract project" });
  const { projects } = await call("ListProjects", { owner_id: USER_ID });
  expect(projects.map((p: { id: string }) => p.id)).toContain(created.id);
});

test("list tasks of a project via project -> task", async () => {
  const p = project();
  const t = client("task/task.proto", "task", "TaskService", cfg.taskGrpc);

  const { project: created } = await p("CreateProject", { name: "with tasks", owner_id: USER_ID });
  const { task } = await t("CreateTask", { title: "in project", status: "todo", project_id: created.id });
  const { tasks } = await p("ListProjectTasks", { project_id: created.id });
  expect(tasks.map((x: { id: string }) => x.id)).toEqual([task.id]);
});
