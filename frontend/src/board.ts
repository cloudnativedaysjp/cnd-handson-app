import { useEffect, useState } from "react";
import { api } from "./utils/api";

// legacy と modern の画面が共有するデータと API。違うのは見せ方だけにする

// フロントエンドのイメージはビルド時の VARIANT で legacy / modern を切り替える。bff/Dockerfile が VITE_VARIANT に渡す
export const legacy = import.meta.env.VITE_VARIANT === "legacy";

export type Project = {
  id: string;
  name: string;
  description: string;
  createdAt?: string;
};
export type Task = {
  id: string;
  title: string;
  description: string;
  columnId: string;
  startTime?: string;
  endTime?: string;
};
export type Column = { id: string; name: string };

export const STANDARD_COLUMNS = ["未着手", "作業中", "完了"];

export const formatDate = (iso?: string) =>
  iso
    ? new Date(iso).toLocaleString("ja-JP", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
      })
    : "";

export const listProjects = () =>
  api<{ projects: Project[] }>("/api/projects").then((r) => r.projects);

export const createProject = (name: string, description: string) =>
  api<{ project: Project }>("/api/projects", {
    method: "POST",
    body: JSON.stringify({ name, description }),
  }).then((r) => r.project);

export const getTask = (taskId: string) =>
  api<{ task: Task }>(`/api/tasks/${taskId}`).then((r) => r.task);

export const createTask = (
  projectId: string,
  title: string,
  columnId: string,
) =>
  api<{ task: Task }>(`/api/projects/${projectId}/tasks`, {
    method: "POST",
    body: JSON.stringify({ title, columnId }),
  }).then((r) => r.task);

// 送った項目だけが変わる
export const updateTask = (
  taskId: string,
  patch: Partial<Pick<Task, "title" | "description" | "columnId">>,
) =>
  api<{ task: Task }>(`/api/tasks/${taskId}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  }).then((r) => r.task);

export const deleteTask = (taskId: string) =>
  api(`/api/tasks/${taskId}`, { method: "DELETE" });

export const createColumn = (projectId: string, name: string) =>
  api<{ column: Column }>(`/api/projects/${projectId}/columns`, {
    method: "POST",
    body: JSON.stringify({ name }),
  }).then((r) => r.column);

export const renameColumn = (columnId: string, name: string) =>
  api<{ column: Column }>(`/api/columns/${columnId}`, {
    method: "PATCH",
    body: JSON.stringify({ name }),
  }).then((r) => r.column);

export const deleteColumn = (columnId: string) =>
  api(`/api/columns/${columnId}`, { method: "DELETE" });

// 作った順に並ぶので、1 つずつ順に作る
export const createStandardColumns = async (projectId: string) => {
  const created: Column[] = [];
  for (const name of STANDARD_COLUMNS) {
    created.push(await createColumn(projectId, name));
  }
  return created;
};

// パンくずやヘッダーに出すため、プロジェクト名を一覧から引く
export const useProjectName = (id: string) => {
  const [name, setName] = useState("");
  useEffect(() => {
    listProjects()
      .then((ps) => setName(ps.find((p) => p.id === id)?.name ?? ""))
      .catch(() => {});
  }, [id]);
  return name;
};

export const useProjectData = (id: string) => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [columns, setColumns] = useState<Column[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => {
    Promise.all([
      api<{ tasks: Task[] }>(`/api/projects/${id}/tasks`),
      api<{ columns: Column[] }>(`/api/projects/${id}/columns`),
    ])
      .then(([t, c]) => {
        setTasks(t.tasks);
        setColumns(c.columns);
      })
      .catch(() =>
        setError("タスクと状態を読み込めませんでした。再読み込みしてください"),
      )
      .finally(() => setLoaded(true));
  }, [id]);
  return { tasks, setTasks, columns, setColumns, loaded, error, setError };
};
