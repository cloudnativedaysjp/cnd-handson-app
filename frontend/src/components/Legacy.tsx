import type React from "react";
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import type { Project } from "../features/projects/types";
import { api } from "../utils/api";

// legacy 版（VARIANT=legacy）の画面。昔の業務システム風に、グレーの枠線の表で見せる。データと API は modern と共通
type Task = { id: string; title: string; status: string; columnId: string };
type Column = { id: string; name: string };

const STATUSES = ["todo", "doing", "done"];

const page: React.CSSProperties = {
  padding: 16,
  fontFamily: '"MS PGothic", Osaka, sans-serif',
  fontSize: 13,
  color: "#000",
  background: "#f0f0f0",
  minHeight: "100vh",
};
const table: React.CSSProperties = {
  borderCollapse: "collapse",
  background: "#fff",
  marginBottom: 16,
  minWidth: 480,
};
const cell: React.CSSProperties = {
  border: "1px solid #808080",
  padding: "2px 8px",
};
const head: React.CSSProperties = {
  ...cell,
  background: "#d4d0c8",
  textAlign: "left",
};
const h1: React.CSSProperties = {
  fontSize: 16,
  borderBottom: "2px solid #808080",
  margin: "0 0 12px",
};

export const LegacyProjectList: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        setProjects(
          (await api<{ projects: Project[] }>("/api/projects")).projects,
        );
      } catch {
        setFailed(true);
      }
    };
    load();
  }, []);

  return (
    <div style={page}>
      <h1 style={h1}>プロジェクト一覧</h1>
      {failed && (
        <p style={{ color: "red" }}>プロジェクトを読み込めませんでした</p>
      )}
      <table style={table} data-testid="project-list">
        <thead>
          <tr>
            <th style={head}>プロジェクト名</th>
            <th style={head}>説明</th>
          </tr>
        </thead>
        <tbody>
          {projects.map((p) => (
            <tr key={p.id}>
              <td style={cell}>
                <Link to={`/projects/${p.id}`} data-testid="project-item">
                  {p.name}
                </Link>
              </td>
              <td style={cell}>{p.description}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export const LegacyProject: React.FC = () => {
  const { id = "" } = useParams();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [columns, setColumns] = useState<Column[]>([]);
  const [title, setTitle] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      try {
        const [t, c] = await Promise.all([
          api<{ tasks: Task[] }>(`/api/projects/${id}/tasks`),
          api<{ columns: Column[] }>(`/api/projects/${id}/columns`),
        ]);
        setTasks(t.tasks);
        setColumns(c.columns);
      } catch {
        setError("読み込めませんでした");
      }
    };
    load();
  }, [id]);

  const create = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    try {
      const { task } = await api<{ task: Task }>(`/api/projects/${id}/tasks`, {
        method: "POST",
        body: JSON.stringify({ title, status: STATUSES[0] }),
      });
      setTasks([...tasks, task]);
      setTitle("");
    } catch {
      setError("登録できませんでした");
    }
  };

  const changeStatus = async (taskId: string, status: string) => {
    try {
      const { task } = await api<{ task: Task }>(`/api/tasks/${taskId}`, {
        method: "PATCH",
        body: JSON.stringify({ status }),
      });
      setTasks(tasks.map((t) => (t.id === taskId ? task : t)));
    } catch {
      setError("更新できませんでした");
    }
  };

  const columnName = (columnId: string) =>
    columns.find((c) => c.id === columnId)?.name ?? "";

  return (
    <div style={page}>
      <p>
        <Link to="/">&lt;&lt; プロジェクト一覧へ戻る</Link>
      </p>
      <h1 style={h1}>タスク一覧</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={create} style={{ marginBottom: 12 }}>
        件名：
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          data-testid="new-task-title"
        />{" "}
        <button type="submit" data-testid="new-task-submit">
          登録
        </button>
      </form>
      <table style={table} data-testid="task-list">
        <thead>
          <tr>
            <th style={head}>件名</th>
            <th style={head}>列</th>
            <th style={head}>状態</th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((t) => (
            <tr key={t.id}>
              <td style={cell}>{t.title}</td>
              <td style={cell}>{columnName(t.columnId)}</td>
              <td style={cell}>
                <select
                  value={t.status}
                  onChange={(e) => changeStatus(t.id, e.target.value)}
                >
                  {/* 一覧に無い状態のタスクも、今の値を表示できるようにする */}
                  {[...new Set([t.status, ...STATUSES])].map((s) => (
                    <option key={s}>{s}</option>
                  ))}
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <h1 style={h1}>列一覧</h1>
      <table style={table}>
        <thead>
          <tr>
            <th style={head}>列名</th>
          </tr>
        </thead>
        <tbody>
          {columns.map((c) => (
            <tr key={c.id}>
              <td style={cell}>{c.name}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
