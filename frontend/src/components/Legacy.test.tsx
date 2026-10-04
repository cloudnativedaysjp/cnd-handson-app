import { fireEvent, render, screen, within } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, expect, test, vi } from "vitest";
import { api } from "../utils/api";
import { LegacyColumns, LegacyProject, LegacyTask } from "./Legacy";

vi.mock("../utils/api", () => ({ api: vi.fn() }));
const mocked = vi.mocked(api);

// method とパスで返す値を決める。書き込みは呼ばれた内容で確かめる
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
          description: "下書きあり",
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
        <Route path="/projects/:id" element={<LegacyProject />} />
        <Route path="/projects/:id/tasks/:taskId" element={<LegacyTask />} />
        <Route path="/projects/:id/columns" element={<LegacyColumns />} />
      </Routes>
    </MemoryRouter>,
  );

const called = (path: string, method: string, body?: object) =>
  expect(mocked).toHaveBeenCalledWith(path, {
    method,
    ...(body && { body: JSON.stringify(body) }),
  });

test("task list: filters by column, registers into the first column, moves a task", async () => {
  show("/projects/p1");
  const list = await screen.findByTestId("task-list");
  expect(await within(list).findByText("資料を作る")).toBeInTheDocument();

  fireEvent.click(screen.getByRole("button", { name: "完了 1件" }));
  expect(within(list).queryByText("資料を作る")).not.toBeInTheDocument();
  expect(screen.getByText("全 2 件中 1 件を表示")).toBeInTheDocument();
  fireEvent.click(screen.getByRole("button", { name: "すべて 2件" }));

  fireEvent.change(screen.getByTestId("new-task-title"), {
    target: { value: "会議" },
  });
  fireEvent.click(screen.getByTestId("new-task-submit"));
  expect(await within(list).findByText("会議")).toBeInTheDocument();
  called("/api/projects/p1/tasks", "POST", { title: "会議", columnId: "c1" });

  fireEvent.change(screen.getByLabelText("資料を作る の状態"), {
    target: { value: "c2" },
  });
  expect(
    await screen.findByRole("button", { name: "完了 2件" }),
  ).toBeInTheDocument();
  called("/api/tasks/t1", "PATCH", { columnId: "c2" });
});

test("task list: offers the standard columns when the project has none", async () => {
  columns = [];
  show("/projects/p1");
  fireEvent.click(await screen.findByRole("button", { name: /標準の列/ }));
  expect(
    await screen.findByRole("button", { name: "作業中 0件" }),
  ).toBeInTheDocument();
  const created = mocked.mock.calls
    .filter(([p, i]) => p === "/api/projects/p1/columns" && i)
    .map(([, i]) => i?.body);
  expect(created).toEqual(
    ["未着手", "作業中", "完了"].map((name) => JSON.stringify({ name })),
  );
});

test("task detail: saves the form and deletes after confirming", async () => {
  show("/projects/p1/tasks/t1");
  expect(await screen.findByDisplayValue("下書きあり")).toBeInTheDocument();
  fireEvent.change(await screen.findByLabelText(/件名/), {
    target: { value: "資料を直す" },
  });
  fireEvent.change(screen.getByLabelText("詳細"), {
    target: { value: "金曜まで" },
  });
  fireEvent.click(screen.getByRole("button", { name: "保存" }));
  expect(await screen.findByText("保存しました")).toBeInTheDocument();
  called("/api/tasks/t1", "PATCH", {
    title: "資料を直す",
    description: "金曜まで",
    columnId: "c1",
  });

  fireEvent.click(screen.getByRole("button", { name: "このタスクを削除" }));
  expect(mocked).not.toHaveBeenCalledWith("/api/tasks/t1", {
    method: "DELETE",
  });
  fireEvent.click(screen.getByRole("button", { name: "削除する" }));
  await vi.waitFor(() => called("/api/tasks/t1", "DELETE"));
});

test("column settings: renames and deletes after confirming", async () => {
  show("/projects/p1/columns");
  fireEvent.change(await screen.findByLabelText("未着手 の列名"), {
    target: { value: "待ち" },
  });
  fireEvent.click(screen.getAllByRole("button", { name: "変更" })[0]);
  await vi.waitFor(() => called("/api/columns/c1", "PATCH", { name: "待ち" }));

  fireEvent.click(screen.getAllByRole("button", { name: "削除" })[1]);
  expect(screen.getByText(/未設定」になります/)).toBeInTheDocument();
  fireEvent.click(screen.getByRole("button", { name: "削除する" }));
  await vi.waitFor(() =>
    expect(screen.queryByLabelText("完了 の列名")).not.toBeInTheDocument(),
  );
  called("/api/columns/c2", "DELETE");
});
