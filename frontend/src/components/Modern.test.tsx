import { fireEvent, render, screen, within } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, expect, test, vi } from "vitest";
import { api } from "../utils/api";
import { ModernBoard, ModernColumns } from "./Modern";

vi.mock("../utils/api", () => ({ api: vi.fn() }));
const mocked = vi.mocked(api);

let columns: { id: string; name: string }[];
beforeEach(() => {
  mocked.mockReset();
  columns = [
    { id: "c1", name: "未着手" },
    { id: "c2", name: "完了" },
  ];
  mocked.mockImplementation(async (path: string, init?: RequestInit) => {
    const method = init?.method ?? "GET";
    const body = init?.body ? JSON.parse(String(init.body)) : {};
    if (path === "/api/projects")
      return { projects: [{ id: "p1", name: "Alpha" }] };
    if (path === "/api/projects/p1/tasks" && method === "GET") {
      return {
        tasks: [
          { id: "t1", title: "資料を作る", description: "", columnId: "c1" },
          { id: "t2", title: "見積を出す", description: "", columnId: "c2" },
        ],
      };
    }
    if (path === "/api/projects/p1/columns" && method === "GET")
      return { columns };
    if (path === "/api/projects/p1/columns")
      return { column: { id: `n-${body.name}`, name: body.name } };
    if (path === "/api/projects/p1/tasks")
      return { task: { id: "t3", title: body.title, columnId: body.columnId } };
    if (path === "/api/tasks/t1" && method === "GET") {
      return {
        task: {
          id: "t1",
          title: "資料を作る",
          description: "下書き",
          columnId: "c1",
        },
      };
    }
    if (path.startsWith("/api/tasks/") && method === "PATCH") {
      return {
        task: { id: "t1", title: "資料を作る", description: "", ...body },
      };
    }
    if (path.startsWith("/api/columns/") && method === "PATCH")
      return { column: { id: "c1", name: body.name } };
    return {};
  });
});

const show = (path: string) =>
  render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route path="/projects/:id" element={<ModernBoard />} />
        <Route path="/projects/:id/tasks/:taskId" element={<ModernBoard />} />
        <Route path="/projects/:id/columns" element={<ModernColumns />} />
      </Routes>
    </MemoryRouter>,
  );

const called = (path: string, method: string, body?: object) =>
  expect(mocked).toHaveBeenCalledWith(path, {
    method,
    ...(body && { body: JSON.stringify(body) }),
  });

test("board: shows a lane per status, filters by title, adds to the first lane", async () => {
  show("/projects/p1");
  const todo = await screen.findByRole("region", { name: "未着手 1件" });
  expect(within(todo).getByText("資料を作る")).toBeInTheDocument();
  expect(screen.getByRole("region", { name: "完了 1件" })).toBeInTheDocument();

  fireEvent.change(screen.getByLabelText("件名で絞り込む"), {
    target: { value: "見積" },
  });
  expect(screen.queryByText("資料を作る")).not.toBeInTheDocument();
  fireEvent.change(screen.getByLabelText("件名で絞り込む"), {
    target: { value: "" },
  });

  fireEvent.change(screen.getByTestId("new-task-title"), {
    target: { value: "会議" },
  });
  fireEvent.click(screen.getByTestId("new-task-submit"));
  expect(
    await screen.findByRole("region", { name: "未着手 2件" }),
  ).toBeInTheDocument();
  called("/api/projects/p1/tasks", "POST", { title: "会議", columnId: "c1" });
});

test("board: offers the standard statuses when the project has none", async () => {
  columns = [];
  show("/projects/p1");
  fireEvent.click(await screen.findByRole("button", { name: /標準の状態/ }));
  expect(
    await screen.findByRole("region", { name: "作業中 0件" }),
  ).toBeInTheDocument();
});

test("task panel: saves the status change and deletes after confirming", async () => {
  show("/projects/p1/tasks/t1");
  const panel = await screen.findByRole("complementary", {
    name: "タスクの詳細",
  });
  expect(await within(panel).findByDisplayValue("下書き")).toBeInTheDocument();
  fireEvent.change(within(panel).getByLabelText("状態"), {
    target: { value: "c2" },
  });
  fireEvent.click(within(panel).getByRole("button", { name: "保存" }));
  expect(await within(panel).findByText("保存しました")).toBeInTheDocument();
  called("/api/tasks/t1", "PATCH", {
    title: "資料を作る",
    description: "下書き",
    columnId: "c2",
  });

  fireEvent.click(
    within(panel).getByRole("button", { name: "このタスクを削除" }),
  );
  fireEvent.click(within(panel).getByRole("button", { name: "削除する" }));
  await vi.waitFor(() => called("/api/tasks/t1", "DELETE"));
});

test("statuses: renames and deletes after confirming", async () => {
  show("/projects/p1/columns");
  fireEvent.change(await screen.findByLabelText("未着手 の名前"), {
    target: { value: "待ち" },
  });
  fireEvent.click(screen.getAllByRole("button", { name: "名前を保存" })[0]);
  await vi.waitFor(() => called("/api/columns/c1", "PATCH", { name: "待ち" }));

  fireEvent.click(screen.getAllByRole("button", { name: "削除" })[1]);
  fireEvent.click(screen.getByRole("button", { name: "削除する" }));
  await vi.waitFor(() =>
    expect(screen.queryByLabelText("完了 の名前")).not.toBeInTheDocument(),
  );
  called("/api/columns/c2", "DELETE");
});
