// Contract for handson-project (#72). Listing tasks of a project needs a new RPC
// (proposed/project_task.proto), so that test is fixme.
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

// fixme: needs ProjectService.ListProjectTasks and task.project_id (proposed/*.proto), see #72/#105.
test.fixme("list tasks of a project via project -> task", async () => {
  const p = client("project/project.proto", "project", "ProjectService", cfg.projectGrpc);
  const t = client("task_project.proto", "task", "TaskService", cfg.taskGrpc, true);
  const list = client("project_task.proto", "project", "ProjectService", cfg.projectGrpc, true);

  const { project: created } = await p("CreateProject", { name: "with tasks", owner_id: USER_ID });
  const { task } = await t("CreateTask", { title: "in project", status: "todo", project_id: created.id });
  const { tasks } = await list("ListProjectTasks", { project_id: created.id });
  expect(tasks.map((x: { id: string }) => x.id)).toEqual([task.id]);
});
