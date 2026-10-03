import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { expect, test, vi } from "vitest";
import { api } from "../utils/api";
import { LegacyProject } from "./Legacy";

vi.mock("../utils/api", () => ({ api: vi.fn() }));

test("legacy project: lists, creates and updates tasks", async () => {
  const mocked = vi.mocked(api);
  mocked.mockImplementation(async (path: string, init?: RequestInit) => {
    if (path === "/api/projects/p1/tasks" && !init) {
      return {
        tasks: [{ id: "t1", title: "first", status: "todo", columnId: "c1" }],
      };
    }
    if (path === "/api/projects/p1/columns") {
      return { columns: [{ id: "c1", name: "backlog" }] };
    }
    if (init?.method === "POST") {
      return {
        task: { id: "t2", title: "second", status: "todo", columnId: "" },
      };
    }
    return {
      task: { id: "t1", title: "first", status: "done", columnId: "c1" },
    };
  });
  render(
    <MemoryRouter initialEntries={["/projects/p1"]}>
      <Routes>
        <Route path="/projects/:id" element={<LegacyProject />} />
      </Routes>
    </MemoryRouter>,
  );

  expect(await screen.findByText("first")).toBeInTheDocument();
  expect(screen.getAllByText("backlog")).toHaveLength(2);

  fireEvent.change(screen.getByTestId("new-task-title"), {
    target: { value: "second" },
  });
  fireEvent.click(screen.getByTestId("new-task-submit"));
  expect(await screen.findByText("second")).toBeInTheDocument();
  expect(mocked).toHaveBeenCalledWith("/api/projects/p1/tasks", {
    method: "POST",
    body: JSON.stringify({ title: "second", status: "todo" }),
  });

  fireEvent.change(screen.getAllByRole("combobox")[0], {
    target: { value: "done" },
  });
  expect(await screen.findByDisplayValue("done")).toBeInTheDocument();
  expect(mocked).toHaveBeenCalledWith("/api/tasks/t1", {
    method: "PATCH",
    body: JSON.stringify({ status: "done" }),
  });
});
