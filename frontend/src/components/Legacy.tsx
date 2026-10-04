import type React from "react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import type { Project } from "../features/projects/types";
import { api } from "../utils/api";
import "./legacy.css";

// legacy 版（VARIANT=legacy）の画面。2000 年代の社内 Web 業務システムの作法で作る。データと API は modern と共通

// 入口のイメージはビルド時の VARIANT で legacy / modern を切り替える。bff/Dockerfile が VITE_VARIANT に渡す
export const legacy = import.meta.env.VITE_VARIANT === "legacy";

type Task = {
  id: string;
  title: string;
  description: string;
  columnId: string;
  startTime?: string;
  endTime?: string;
};
type Column = { id: string; name: string };

const STANDARD_COLUMNS = ["未着手", "作業中", "完了"];

const formatDate = (iso?: string) =>
  iso
    ? new Date(iso).toLocaleString("ja-JP", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
      })
    : "";

// パンくずとメニューに出すため、プロジェクト名を一覧から引く
const useProjectName = (id: string) => {
  const [name, setName] = useState("");
  useEffect(() => {
    api<{ projects: Project[] }>("/api/projects")
      .then((r) => setName(r.projects.find((p) => p.id === id)?.name ?? ""))
      .catch(() => {});
  }, [id]);
  return name;
};

const useProjectData = (id: string) => {
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
        setError("タスクと列を読み込めませんでした。再読み込みしてください"),
      )
      .finally(() => setLoaded(true));
  }, [id]);
  return { tasks, setTasks, columns, setColumns, loaded, error, setError };
};

const createColumn = (projectId: string, name: string) =>
  api<{ column: Column }>(`/api/projects/${projectId}/columns`, {
    method: "POST",
    body: JSON.stringify({ name }),
  }).then((r) => r.column);

// 作った順に並ぶので、1 つずつ順に作る
const createStandardColumns = async (projectId: string) => {
  const created: Column[] = [];
  for (const name of STANDARD_COLUMNS) {
    created.push(await createColumn(projectId, name));
  }
  return created;
};

type Crumb = { label: string; to?: string };

const Frame: React.FC<{
  screenId: string;
  title: string;
  project?: { id: string; name: string };
  crumbs: Crumb[];
  error?: string;
  children: React.ReactNode;
}> = ({ screenId, title, project, crumbs, error, children }) => {
  const navigate = useNavigate();
  useEffect(() => {
    document.title = `${title} - タスク管理システム`;
  }, [title]);
  const current = (id: string) => (screenId === id ? "page" : undefined);
  return (
    <div className="legacy">
      <header className="legacy-band">
        <strong>タスク管理システム</strong>
        <span>
          <span className="legacy-screen-id">
            {screenId} {title}
          </span>
          <button type="button" onClick={() => navigate("/logout")}>
            ログアウト
          </button>
        </span>
      </header>
      <nav className="legacy-menu" aria-label="メニュー">
        <Link to="/" aria-current={current("PRJ-010")}>
          プロジェクト一覧
        </Link>
        {project && (
          <>
            <Link
              to={`/projects/${project.id}`}
              aria-current={current("TSK-010")}
            >
              タスク一覧
            </Link>
            <Link
              to={`/projects/${project.id}/columns`}
              aria-current={current("COL-010")}
            >
              列の設定
            </Link>
          </>
        )}
      </nav>
      <main className="legacy-main">
        <p className="legacy-crumbs">
          {crumbs.map((c, i) => (
            <span key={c.label}>
              {i > 0 && " > "}
              {c.to ? <Link to={c.to}>{c.label}</Link> : c.label}
            </span>
          ))}
        </p>
        <h1>{title}</h1>
        {error && (
          <p className="legacy-error" role="alert">
            {error}
          </p>
        )}
        {children}
      </main>
    </div>
  );
};

const StandardColumnsButton: React.FC<{
  projectId: string;
  onCreated: (columns: Column[]) => void;
  onError: (message: string) => void;
}> = ({ projectId, onCreated, onError }) => (
  <button
    type="button"
    className="legacy-primary"
    onClick={() =>
      createStandardColumns(projectId)
        .then(onCreated)
        .catch(() => onError("列を作成できませんでした"))
    }
  >
    標準の列（{STANDARD_COLUMNS.join("・")}）を作成
  </button>
);

export const LegacyProjectList: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    api<{ projects: Project[] }>("/api/projects")
      .then((r) => setProjects(r.projects))
      .catch(() =>
        setError("プロジェクトを読み込めませんでした。再読み込みしてください"),
      )
      .finally(() => setLoaded(true));
  }, []);

  const create = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError("名称を入力してください");
      return;
    }
    try {
      const { project } = await api<{ project: Project }>("/api/projects", {
        method: "POST",
        body: JSON.stringify({ name, description }),
      });
      setProjects((ps) => [...ps, project]);
      setName("");
      setDescription("");
      setError("");
    } catch {
      setError("プロジェクトを作成できませんでした");
    }
  };

  return (
    <Frame
      screenId="PRJ-010"
      title="プロジェクト一覧"
      crumbs={[{ label: "プロジェクト一覧" }]}
      error={error}
    >
      <div className="legacy-scroll">
        <table data-testid="project-list">
          <thead>
            <tr>
              <th>プロジェクト名</th>
              <th>説明</th>
              <th className="legacy-date">作成日時</th>
            </tr>
          </thead>
          <tbody>
            {projects.map((p) => (
              <tr key={p.id}>
                <td>
                  <Link to={`/projects/${p.id}`} data-testid="project-item">
                    {p.name}
                  </Link>
                </td>
                <td>{p.description}</td>
                <td>{formatDate(p.createdAt)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {loaded && projects.length === 0 && (
        <p className="legacy-note">
          プロジェクトはまだありません。下のフォームで作成してください。
        </p>
      )}
      <h2>プロジェクトの新規作成</h2>
      <form onSubmit={create}>
        <table className="legacy-form">
          <tbody>
            <tr>
              <th>
                <label htmlFor="project-name">名称</label>
                <span className="legacy-required">必須</span>
              </th>
              <td>
                <input
                  id="project-name"
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </td>
            </tr>
            <tr>
              <th>
                <label htmlFor="project-description">説明</label>
              </th>
              <td>
                <input
                  id="project-description"
                  type="text"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </td>
            </tr>
          </tbody>
        </table>
        <div className="legacy-actions">
          <button type="submit" className="legacy-primary">
            作成
          </button>
        </div>
      </form>
    </Frame>
  );
};

// 列で絞り込むときの値。"" はすべて、NONE は列に入っていないタスク
const NONE = "__none__";

export const LegacyProject: React.FC = () => {
  const { id = "" } = useParams();
  const projectName = useProjectName(id);
  const { tasks, setTasks, columns, setColumns, loaded, error, setError } =
    useProjectData(id);
  const [title, setTitle] = useState("");
  const [keyword, setKeyword] = useState("");
  const [search, setSearch] = useState({ keyword: "", column: "" });

  const inColumn = (t: Task) => columns.some((c) => c.id === t.columnId);
  const unassigned = tasks.filter((t) => !inColumn(t));
  const shown = tasks.filter(
    (t) =>
      t.title.includes(search.keyword) &&
      (search.column === "" ||
        (search.column === NONE ? !inColumn(t) : t.columnId === search.column)),
  );

  const create = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    try {
      const { task } = await api<{ task: Task }>(`/api/projects/${id}/tasks`, {
        method: "POST",
        // 件名だけで登録したタスクは、先頭の列に入れる
        body: JSON.stringify({ title, columnId: columns[0]?.id ?? "" }),
      });
      setTasks((ts) => [...ts, task]);
      setTitle("");
      setError("");
    } catch {
      setError("タスクを登録できませんでした");
    }
  };

  const move = async (taskId: string, columnId: string) => {
    try {
      const { task } = await api<{ task: Task }>(`/api/tasks/${taskId}`, {
        method: "PATCH",
        body: JSON.stringify({ columnId }),
      });
      setTasks((ts) => ts.map((t) => (t.id === taskId ? task : t)));
    } catch {
      setError("状態を変更できませんでした");
    }
  };

  const filterBy = (column: string) => setSearch({ ...search, column });

  return (
    <Frame
      screenId="TSK-010"
      title="タスク一覧"
      project={{ id, name: projectName }}
      crumbs={[
        { label: "プロジェクト一覧", to: "/" },
        { label: projectName || "プロジェクト" },
        { label: "タスク一覧" },
      ]}
      error={error}
    >
      {loaded && columns.length === 0 && (
        <div className="legacy-panel">
          <p className="legacy-note">
            このプロジェクトには列がありません。列がタスクの状態になります。
          </p>
          <div className="legacy-actions">
            <StandardColumnsButton
              projectId={id}
              onCreated={setColumns}
              onError={setError}
            />
            <Link to={`/projects/${id}/columns`}>列を自分で設定する</Link>
          </div>
        </div>
      )}

      <ul className="legacy-counts" aria-label="状態ごとの件数">
        <li>
          <button
            type="button"
            aria-pressed={search.column === ""}
            onClick={() => filterBy("")}
          >
            すべて {tasks.length}件
          </button>
        </li>
        {columns.map((c) => (
          <li key={c.id}>
            <button
              type="button"
              aria-pressed={search.column === c.id}
              onClick={() => filterBy(c.id)}
            >
              {c.name} {tasks.filter((t) => t.columnId === c.id).length}件
            </button>
          </li>
        ))}
        {unassigned.length > 0 && (
          <li>
            <button
              type="button"
              aria-pressed={search.column === NONE}
              onClick={() => filterBy(NONE)}
            >
              未設定 {unassigned.length}件
            </button>
          </li>
        )}
      </ul>

      <form
        className="legacy-panel legacy-row"
        onSubmit={(e) => {
          e.preventDefault();
          setSearch({ ...search, keyword });
        }}
      >
        <label>
          件名{" "}
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
          />
        </label>
        <label>
          状態{" "}
          <select
            value={search.column}
            onChange={(e) => filterBy(e.target.value)}
          >
            <option value="">すべて</option>
            {columns.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
            <option value={NONE}>未設定</option>
          </select>
        </label>
        <button type="submit">検索</button>
        <button
          type="button"
          onClick={() => {
            setKeyword("");
            setSearch({ keyword: "", column: "" });
          }}
        >
          クリア
        </button>
      </form>

      <form className="legacy-panel legacy-row" onSubmit={create}>
        <label htmlFor="new-task-title">新しいタスクの件名</label>
        <input
          id="new-task-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          data-testid="new-task-title"
        />
        <button
          type="submit"
          className="legacy-primary"
          data-testid="new-task-submit"
        >
          登録
        </button>
        {columns[0] && <span>「{columns[0].name}」に入ります</span>}
      </form>

      <div className="legacy-scroll">
        <table data-testid="task-list">
          <thead>
            <tr>
              <th className="legacy-num">No.</th>
              <th>件名</th>
              <th>状態</th>
            </tr>
          </thead>
          <tbody>
            {shown.map((t, i) => (
              <tr key={t.id}>
                <td className="legacy-num">{i + 1}</td>
                <td>
                  <Link to={`/projects/${id}/tasks/${t.id}`}>{t.title}</Link>
                </td>
                <td>
                  <select
                    aria-label={`${t.title} の状態`}
                    value={inColumn(t) ? t.columnId : ""}
                    onChange={(e) => move(t.id, e.target.value)}
                  >
                    {!inColumn(t) && <option value="">未設定</option>}
                    {columns.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="legacy-note">
        {loaded && tasks.length === 0
          ? "タスクはまだありません。上の欄に件名を入れて登録してください。"
          : `全 ${tasks.length} 件中 ${shown.length} 件を表示`}
      </p>
    </Frame>
  );
};

export const LegacyTask: React.FC = () => {
  const { id = "", taskId = "" } = useParams();
  const navigate = useNavigate();
  const projectName = useProjectName(id);
  const { columns, error, setError } = useProjectData(id);
  const [task, setTask] = useState<Task>();
  const [missing, setMissing] = useState(false);
  const [form, setForm] = useState({
    title: "",
    description: "",
    columnId: "",
  });
  const [saved, setSaved] = useState(false);
  const [confirming, setConfirming] = useState(false);

  useEffect(() => {
    api<{ task: Task }>(`/api/tasks/${taskId}`)
      .then(({ task }) => {
        setTask(task);
        setForm({
          title: task.title,
          description: task.description,
          columnId: task.columnId,
        });
      })
      .catch(() => setMissing(true));
  }, [taskId]);

  const save = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title.trim()) {
      setError("件名を入力してください");
      return;
    }
    try {
      const { task } = await api<{ task: Task }>(`/api/tasks/${taskId}`, {
        method: "PATCH",
        body: JSON.stringify(form),
      });
      setTask(task);
      setError("");
      setSaved(true);
    } catch {
      setError("タスクを保存できませんでした");
    }
  };

  const remove = async () => {
    try {
      await api(`/api/tasks/${taskId}`, { method: "DELETE" });
      navigate(`/projects/${id}`);
    } catch {
      setError("タスクを削除できませんでした");
    }
  };

  const edit = (patch: Partial<typeof form>) => {
    setForm({ ...form, ...patch });
    setSaved(false);
  };

  return (
    <Frame
      screenId="TSK-020"
      title="タスク詳細"
      project={{ id, name: projectName }}
      crumbs={[
        { label: "プロジェクト一覧", to: "/" },
        { label: projectName || "プロジェクト", to: `/projects/${id}` },
        { label: "タスク詳細" },
      ]}
      error={error}
    >
      {missing && (
        <p className="legacy-note">
          このタスクは見つかりません。削除された可能性があります。{" "}
          <Link to={`/projects/${id}`}>タスク一覧に戻る</Link>
        </p>
      )}
      {task && (
        <>
          <form onSubmit={save}>
            <table className="legacy-form">
              <tbody>
                <tr>
                  <th>
                    <label htmlFor="task-title">件名</label>
                    <span className="legacy-required">必須</span>
                  </th>
                  <td>
                    <input
                      id="task-title"
                      type="text"
                      value={form.title}
                      onChange={(e) => edit({ title: e.target.value })}
                    />
                  </td>
                </tr>
                <tr>
                  <th>
                    <label htmlFor="task-description">詳細</label>
                  </th>
                  <td>
                    <textarea
                      id="task-description"
                      rows={5}
                      value={form.description}
                      onChange={(e) => edit({ description: e.target.value })}
                    />
                  </td>
                </tr>
                <tr>
                  <th>
                    <label htmlFor="task-column">状態</label>
                  </th>
                  <td>
                    <select
                      id="task-column"
                      value={form.columnId}
                      onChange={(e) => edit({ columnId: e.target.value })}
                    >
                      {!columns.some((c) => c.id === form.columnId) && (
                        <option value={form.columnId}>未設定</option>
                      )}
                      {columns.map((c) => (
                        <option key={c.id} value={c.id}>
                          {c.name}
                        </option>
                      ))}
                    </select>
                  </td>
                </tr>
                <tr>
                  <th>登録日時</th>
                  <td>{formatDate(task.startTime)}</td>
                </tr>
                <tr>
                  <th>更新日時</th>
                  <td>{formatDate(task.endTime)}</td>
                </tr>
              </tbody>
            </table>
            <div className="legacy-actions">
              <button type="submit" className="legacy-primary">
                保存
              </button>
              <Link to={`/projects/${id}`}>タスク一覧に戻る</Link>
              {saved && <span role="status">保存しました</span>}
            </div>
          </form>

          <h2>タスクの削除</h2>
          <div className="legacy-actions">
            {confirming ? (
              <>
                <span className="legacy-confirm">
                  削除すると元に戻せません。削除しますか？
                </span>
                <button type="button" onClick={remove}>
                  削除する
                </button>
                <button type="button" onClick={() => setConfirming(false)}>
                  やめる
                </button>
              </>
            ) : (
              <button type="button" onClick={() => setConfirming(true)}>
                このタスクを削除
              </button>
            )}
          </div>
        </>
      )}
    </Frame>
  );
};

export const LegacyColumns: React.FC = () => {
  const { id = "" } = useParams();
  const projectName = useProjectName(id);
  const { columns, setColumns, loaded, error, setError } = useProjectData(id);
  const [names, setNames] = useState<Record<string, string>>({});
  const [newName, setNewName] = useState("");
  const [confirming, setConfirming] = useState("");

  const add = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim()) {
      setError("列名を入力してください");
      return;
    }
    try {
      const column = await createColumn(id, newName);
      setColumns((cs) => [...cs, column]);
      setNewName("");
      setError("");
    } catch {
      setError("列を追加できませんでした");
    }
  };

  const rename = async (columnId: string) => {
    const name = names[columnId]?.trim();
    if (!name) return;
    try {
      const { column } = await api<{ column: Column }>(
        `/api/columns/${columnId}`,
        { method: "PATCH", body: JSON.stringify({ name }) },
      );
      setColumns((cs) => cs.map((c) => (c.id === columnId ? column : c)));
      setError("");
    } catch {
      setError("列名を変更できませんでした");
    }
  };

  const remove = async (columnId: string) => {
    try {
      await api(`/api/columns/${columnId}`, { method: "DELETE" });
      setColumns((cs) => cs.filter((c) => c.id !== columnId));
      setConfirming("");
      setError("");
    } catch {
      setError("列を削除できませんでした");
    }
  };

  return (
    <Frame
      screenId="COL-010"
      title="列の設定"
      project={{ id, name: projectName }}
      crumbs={[
        { label: "プロジェクト一覧", to: "/" },
        { label: projectName || "プロジェクト", to: `/projects/${id}` },
        { label: "列の設定" },
      ]}
      error={error}
    >
      <p className="legacy-note">
        列はタスクの状態です。タスク一覧では、ここに並んだ順に表示します。
      </p>
      {loaded && columns.length === 0 && (
        <div className="legacy-panel">
          <p className="legacy-note">列がありません。</p>
          <StandardColumnsButton
            projectId={id}
            onCreated={setColumns}
            onError={setError}
          />
        </div>
      )}
      {columns.length > 0 && (
        <div className="legacy-scroll">
          <table>
            <thead>
              <tr>
                <th className="legacy-num">順</th>
                <th>列名</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {columns.map((c, i) => (
                <tr key={c.id}>
                  <td className="legacy-num">{i + 1}</td>
                  <td>
                    <div className="legacy-row">
                      <input
                        type="text"
                        aria-label={`${c.name} の列名`}
                        value={names[c.id] ?? c.name}
                        onChange={(e) =>
                          setNames({ ...names, [c.id]: e.target.value })
                        }
                      />
                      <button type="button" onClick={() => rename(c.id)}>
                        変更
                      </button>
                    </div>
                  </td>
                  <td>
                    {confirming === c.id ? (
                      <div className="legacy-row">
                        <span className="legacy-confirm">
                          この列のタスクは「未設定」になります。
                        </span>
                        <button type="button" onClick={() => remove(c.id)}>
                          削除する
                        </button>
                        <button type="button" onClick={() => setConfirming("")}>
                          やめる
                        </button>
                      </div>
                    ) : (
                      <button type="button" onClick={() => setConfirming(c.id)}>
                        削除
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <h2>列の追加</h2>
      <form className="legacy-row" onSubmit={add}>
        <label htmlFor="column-name">列名</label>
        <input
          id="column-name"
          type="text"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
        />
        <button type="submit" className="legacy-primary">
          追加
        </button>
      </form>
    </Frame>
  );
};
