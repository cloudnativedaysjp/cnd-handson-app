import {
  DndContext,
  type DragEndEvent,
  KeyboardSensor,
  PointerSensor,
  useDraggable,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import type React from "react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  type Column,
  type Project,
  STANDARD_COLUMNS,
  type Task,
  createColumn,
  createProject,
  createStandardColumns,
  createTask,
  deleteColumn,
  deleteTask,
  formatDate,
  getTask,
  listProjects,
  renameColumn,
  updateTask,
  useProjectData,
  useProjectName,
} from "../board";
import "./modern.css";

// modern 版（VARIANT=modern）の画面。状態ごとのレーンにカードを並べるカンバン。データと API は src/board.ts で legacy と共有する

const Shell: React.FC<{
  project?: { id: string; name: string };
  current: "projects" | "board" | "statuses";
  children: React.ReactNode;
}> = ({ project, current, children }) => {
  const navigate = useNavigate();
  const here = (name: string) => (current === name ? "page" : undefined);
  return (
    <div className="modern">
      <header className="modern-bar">
        <Link to="/" className="modern-brand">
          Kanban
        </Link>
        <nav className="modern-nav" aria-label="メニュー">
          <Link to="/" aria-current={here("projects")}>
            プロジェクト
          </Link>
          {project && (
            <>
              <Link to={`/projects/${project.id}`} aria-current={here("board")}>
                ボード
              </Link>
              <Link
                to={`/projects/${project.id}/columns`}
                aria-current={here("statuses")}
              >
                状態の設定
              </Link>
            </>
          )}
        </nav>
        <button
          type="button"
          className="modern-button is-quiet"
          onClick={() => navigate("/logout")}
        >
          ログアウト
        </button>
      </header>
      <main className="modern-main">{children}</main>
    </div>
  );
};

const ErrorNote: React.FC<{ message: string }> = ({ message }) =>
  message ? (
    <p className="modern-error" role="alert">
      {message}
    </p>
  ) : null;

const StandardColumnsButton: React.FC<{
  projectId: string;
  onCreated: (columns: Column[]) => void;
  onError: (message: string) => void;
}> = ({ projectId, onCreated, onError }) => (
  <button
    type="button"
    className="modern-button is-primary"
    onClick={() =>
      createStandardColumns(projectId)
        .then(onCreated)
        .catch(() => onError("状態を作成できませんでした"))
    }
  >
    標準の状態（{STANDARD_COLUMNS.join("・")}）を作成
  </button>
);

export const ModernProjectList: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loaded, setLoaded] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    document.title = "プロジェクト - Kanban";
    listProjects()
      .then(setProjects)
      .catch(() =>
        setError("プロジェクトを読み込めませんでした。再読み込みしてください"),
      )
      .finally(() => setLoaded(true));
  }, []);

  const create = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError("プロジェクト名を入力してください");
      return;
    }
    try {
      const project = await createProject(name, description);
      setProjects((ps) => [project, ...ps]);
      setName("");
      setDescription("");
      setError("");
    } catch {
      setError("プロジェクトを作成できませんでした");
    }
  };

  return (
    <Shell current="projects">
      <h1>プロジェクト</h1>
      <p className="modern-lead">開くとボードが表示されます。</p>
      <ErrorNote message={error} />
      <div className="modern-projects">
        <div>
          <ul className="modern-list" data-testid="project-list">
            {projects.map((p) => (
              <li key={p.id}>
                <Link to={`/projects/${p.id}`} data-testid="project-item">
                  <strong>{p.name}</strong>
                  {p.description && <small>{p.description}</small>}
                  <time dateTime={p.createdAt}>{formatDate(p.createdAt)}</time>
                </Link>
              </li>
            ))}
          </ul>
          {loaded && projects.length === 0 && (
            <div className="modern-empty">
              <p>
                プロジェクトはまだありません。新しいプロジェクトを作成してください。
              </p>
            </div>
          )}
        </div>
        <form className="modern-panel" onSubmit={create}>
          <h2>新しいプロジェクト</h2>
          <label className="modern-field">
            <span>プロジェクト名</span>
            <input
              className="modern-input"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </label>
          <label className="modern-field">
            <span>説明</span>
            <input
              className="modern-input"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </label>
          <button type="submit" className="modern-button is-primary">
            作成
          </button>
        </form>
      </div>
    </Shell>
  );
};

const Card: React.FC<{ task: Task; projectId: string }> = ({
  task,
  projectId,
}) => {
  const { attributes, listeners, setNodeRef, transform, isDragging } =
    useDraggable({ id: task.id });
  // dnd-kit は role="button" を付ける。リンクのままにして、Enter で詳細を開けるようにする。持ち上げるのは Space
  const { role: _role, ...linkAttributes } = attributes;
  return (
    <li>
      <Link
        ref={setNodeRef}
        to={`/projects/${projectId}/tasks/${task.id}`}
        className={`modern-card${isDragging ? " is-dragging" : ""}`}
        style={
          transform
            ? { transform: `translate(${transform.x}px, ${transform.y}px)` }
            : undefined
        }
        {...listeners}
        {...linkAttributes}
      >
        <strong>{task.title}</strong>
        {task.description && <small>{task.description}</small>}
      </Link>
    </li>
  );
};

const Lane: React.FC<{
  id: string;
  name: string;
  tasks: Task[];
  projectId: string;
}> = ({ id, name, tasks, projectId }) => {
  const { setNodeRef, isOver } = useDroppable({ id });
  return (
    <section
      ref={setNodeRef}
      className={`modern-lane${isOver ? " is-over" : ""}`}
      aria-label={`${name} ${tasks.length}件`}
    >
      <header>
        {name}
        <span className="modern-count">{tasks.length}</span>
      </header>
      <ul className="modern-cards">
        {tasks.map((t) => (
          <Card key={t.id} task={t} projectId={projectId} />
        ))}
      </ul>
      {tasks.length === 0 && (
        <p className="modern-lane-empty">ここにカードをドラッグ</p>
      )}
    </section>
  );
};

// 列に入っていないタスクをまとめるレーン。ドロップ先にはしない
const UNSET = "__unset__";

export const ModernBoard: React.FC = () => {
  const { id = "", taskId } = useParams();
  const projectName = useProjectName(id);
  const { tasks, setTasks, columns, setColumns, loaded, error, setError } =
    useProjectData(id);
  const [title, setTitle] = useState("");
  const [query, setQuery] = useState("");
  // クリックでカードを開けるよう、少し動かしてからドラッグを始める
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(KeyboardSensor, {
      keyboardCodes: {
        start: ["Space"],
        cancel: ["Escape"],
        end: ["Space"],
      },
    }),
  );

  useEffect(() => {
    document.title = `${projectName || "ボード"} - Kanban`;
  }, [projectName]);

  const inColumn = (t: Task) => columns.some((c) => c.id === t.columnId);
  const shown = tasks.filter((t) => t.title.includes(query));
  const unset = shown.filter((t) => !inColumn(t));

  const create = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    try {
      // 新しいタスクは先頭の状態に入れる
      const task = await createTask(id, title, columns[0]?.id ?? "");
      setTasks((ts) => [...ts, task]);
      setTitle("");
      setError("");
    } catch {
      setError("タスクを追加できませんでした");
    }
  };

  const move = async ({ active, over }: DragEndEvent) => {
    const columnId = over?.id;
    const task = tasks.find((t) => t.id === active.id);
    if (!task || typeof columnId !== "string" || columnId === UNSET) return;
    if (task.columnId === columnId) return;
    // 先に動かしてから保存する。失敗したら戻す
    setTasks((ts) =>
      ts.map((t) => (t.id === task.id ? { ...t, columnId } : t)),
    );
    try {
      await updateTask(task.id, { columnId });
    } catch {
      setTasks((ts) => ts.map((t) => (t.id === task.id ? task : t)));
      setError("状態を変更できませんでした");
    }
  };

  return (
    <Shell project={{ id, name: projectName }} current="board">
      <h1>{projectName || "ボード"}</h1>
      <p className="modern-lead">
        カードをドラッグすると状態が変わります。カードを開くと詳細を編集できます。
      </p>
      <ErrorNote message={error} />
      <div className="modern-toolbar">
        <form onSubmit={create}>
          <input
            className="modern-input"
            placeholder="新しいタスクの件名"
            aria-label="新しいタスクの件名"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            data-testid="new-task-title"
          />
          <button
            type="submit"
            className="modern-button is-primary"
            data-testid="new-task-submit"
          >
            追加
          </button>
        </form>
        <input
          type="search"
          className="modern-input modern-search"
          placeholder="件名で絞り込む"
          aria-label="件名で絞り込む"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      {loaded && columns.length === 0 && (
        <div className="modern-empty">
          <p>状態がまだありません。状態ごとにレーンができます。</p>
          <div className="modern-row" style={{ justifyContent: "center" }}>
            <StandardColumnsButton
              projectId={id}
              onCreated={setColumns}
              onError={setError}
            />
            <Link to={`/projects/${id}/columns`} className="modern-button">
              自分で設定する
            </Link>
          </div>
        </div>
      )}

      <DndContext sensors={sensors} onDragEnd={move}>
        <div className="modern-board" data-testid="task-list">
          {columns.map((c) => (
            <Lane
              key={c.id}
              id={c.id}
              name={c.name}
              tasks={shown.filter((t) => t.columnId === c.id)}
              projectId={id}
            />
          ))}
          {unset.length > 0 && (
            <Lane id={UNSET} name="未設定" tasks={unset} projectId={id} />
          )}
        </div>
      </DndContext>

      {taskId && (
        <TaskDrawer
          projectId={id}
          taskId={taskId}
          columns={columns}
          onSaved={(task) =>
            setTasks((ts) => ts.map((t) => (t.id === task.id ? task : t)))
          }
          onDeleted={() => setTasks((ts) => ts.filter((t) => t.id !== taskId))}
        />
      )}
    </Shell>
  );
};

const TaskDrawer: React.FC<{
  projectId: string;
  taskId: string;
  columns: Column[];
  onSaved: (task: Task) => void;
  onDeleted: () => void;
}> = ({ projectId, taskId, columns, onSaved, onDeleted }) => {
  const navigate = useNavigate();
  const [task, setTask] = useState<Task>();
  const [form, setForm] = useState({
    title: "",
    description: "",
    columnId: "",
  });
  const [error, setError] = useState("");
  const [saved, setSaved] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const close = () => navigate(`/projects/${projectId}`);

  useEffect(() => {
    getTask(taskId)
      .then((t) => {
        setTask(t);
        setForm({
          title: t.title,
          description: t.description,
          columnId: t.columnId,
        });
      })
      .catch(() =>
        setError("このタスクは見つかりません。削除された可能性があります"),
      );
  }, [taskId]);

  // Esc で閉じる
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === "Escape" && close();
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  });

  const save = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.title.trim()) {
      setError("件名を入力してください");
      return;
    }
    try {
      const t = await updateTask(taskId, form);
      setTask(t);
      onSaved(t);
      setError("");
      setSaved(true);
    } catch {
      setError("タスクを保存できませんでした");
    }
  };

  const remove = async () => {
    try {
      await deleteTask(taskId);
      onDeleted();
      close();
    } catch {
      setError("タスクを削除できませんでした");
    }
  };

  const edit = (patch: Partial<typeof form>) => {
    setForm({ ...form, ...patch });
    setSaved(false);
  };

  return (
    <>
      <button
        type="button"
        className="modern-scrim"
        aria-label="詳細を閉じる"
        tabIndex={-1}
        onClick={close}
      />
      <aside className="modern-drawer" aria-label="タスクの詳細">
        <header>
          <strong>タスクの詳細</strong>
          <button
            type="button"
            className="modern-button is-quiet"
            onClick={close}
          >
            閉じる
          </button>
        </header>
        <ErrorNote message={error} />
        {task && (
          <form onSubmit={save}>
            <label className="modern-field">
              <span>件名</span>
              <input
                className="modern-input"
                value={form.title}
                onChange={(e) => edit({ title: e.target.value })}
              />
            </label>
            <label className="modern-field">
              <span>状態</span>
              <select
                className="modern-select"
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
            </label>
            <label className="modern-field">
              <span>詳細</span>
              <textarea
                className="modern-textarea"
                value={form.description}
                onChange={(e) => edit({ description: e.target.value })}
              />
            </label>
            <div className="modern-row">
              <button type="submit" className="modern-button is-primary">
                保存
              </button>
              {saved && <output>保存しました</output>}
            </div>
            <div className="modern-meta">
              <div>
                <span>作成</span>
                <time dateTime={task.startTime}>
                  {formatDate(task.startTime)}
                </time>
              </div>
              <div>
                <span>更新</span>
                <time dateTime={task.endTime}>{formatDate(task.endTime)}</time>
              </div>
            </div>
          </form>
        )}
        {task && (
          <div className="modern-confirm">
            {confirming ? (
              <>
                <p>削除すると元に戻せません。</p>
                <div className="modern-row">
                  <button
                    type="button"
                    className="modern-button is-danger"
                    onClick={remove}
                  >
                    削除する
                  </button>
                  <button
                    type="button"
                    className="modern-button"
                    onClick={() => setConfirming(false)}
                  >
                    やめる
                  </button>
                </div>
              </>
            ) : (
              <button
                type="button"
                className="modern-button is-quiet"
                onClick={() => setConfirming(true)}
              >
                このタスクを削除
              </button>
            )}
          </div>
        )}
      </aside>
    </>
  );
};

export const ModernColumns: React.FC = () => {
  const { id = "" } = useParams();
  const projectName = useProjectName(id);
  const { columns, setColumns, loaded, error, setError } = useProjectData(id);
  const [names, setNames] = useState<Record<string, string>>({});
  const [newName, setNewName] = useState("");
  const [confirming, setConfirming] = useState("");

  useEffect(() => {
    document.title = "状態の設定 - Kanban";
  }, []);

  const add = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim()) {
      setError("状態の名前を入力してください");
      return;
    }
    try {
      const column = await createColumn(id, newName);
      setColumns((cs) => [...cs, column]);
      setNewName("");
      setError("");
    } catch {
      setError("状態を追加できませんでした");
    }
  };

  const rename = async (columnId: string) => {
    const name = names[columnId]?.trim();
    if (!name) return;
    try {
      const column = await renameColumn(columnId, name);
      setColumns((cs) => cs.map((c) => (c.id === columnId ? column : c)));
      setError("");
    } catch {
      setError("名前を変更できませんでした");
    }
  };

  const remove = async (columnId: string) => {
    try {
      await deleteColumn(columnId);
      setColumns((cs) => cs.filter((c) => c.id !== columnId));
      setConfirming("");
      setError("");
    } catch {
      setError("状態を削除できませんでした");
    }
  };

  return (
    <Shell project={{ id, name: projectName }} current="statuses">
      <div className="modern-narrow">
        <h1>状態の設定</h1>
        <p className="modern-lead">
          {projectName ? `${projectName} の` : ""}
          ボードには、ここに並んだ順にレーンが表示されます。
        </p>
        <ErrorNote message={error} />
        {loaded && columns.length === 0 && (
          <div className="modern-empty">
            <p>状態がまだありません。</p>
            <StandardColumnsButton
              projectId={id}
              onCreated={setColumns}
              onError={setError}
            />
          </div>
        )}
        {columns.length > 0 && (
          <ul className="modern-statuses">
            {columns.map((c) => (
              <li key={c.id}>
                <input
                  className="modern-input"
                  aria-label={`${c.name} の名前`}
                  value={names[c.id] ?? c.name}
                  onChange={(e) =>
                    setNames({ ...names, [c.id]: e.target.value })
                  }
                />
                <button
                  type="button"
                  className="modern-button"
                  onClick={() => rename(c.id)}
                >
                  名前を保存
                </button>
                {confirming === c.id ? (
                  <>
                    <button
                      type="button"
                      className="modern-button is-danger"
                      onClick={() => remove(c.id)}
                    >
                      削除する
                    </button>
                    <button
                      type="button"
                      className="modern-button"
                      onClick={() => setConfirming("")}
                    >
                      やめる
                    </button>
                  </>
                ) : (
                  <button
                    type="button"
                    className="modern-button is-quiet"
                    onClick={() => setConfirming(c.id)}
                  >
                    削除
                  </button>
                )}
                {confirming === c.id && (
                  <small style={{ color: "var(--danger)", flexBasis: "100%" }}>
                    この状態のタスクは「未設定」になります。
                  </small>
                )}
              </li>
            ))}
          </ul>
        )}
        <h2>状態を追加</h2>
        <form className="modern-row" onSubmit={add}>
          <input
            className="modern-input"
            aria-label="追加する状態の名前"
            placeholder="例: レビュー待ち"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
          <button type="submit" className="modern-button is-primary">
            追加
          </button>
        </form>
      </div>
    </Shell>
  );
};
