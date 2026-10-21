import { logout } from "../features/sessions/slice";
import { store } from "../store";

export class ApiError extends Error {
  constructor(
    path: string,
    readonly status: number,
  ) {
    super(`${path}: ${status}`);
  }
}

// 入口（bff）の REST API を呼ぶ。401 ならセッションを捨て、PrivateLayout がログイン画面に戻す
export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const { token } = store.getState().session;
  const headers = new Headers(init.headers);
  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  const start = performance.now();
  const res = await fetch(path, { ...init, headers });
  record({
    method: init.method ?? "GET",
    path,
    status: res.status,
    ms: Math.round(performance.now() - start),
    version: res.headers.get("X-App-Version") ?? "",
  });
  // 待っている間に別のトークンでログインし直していたら、そのセッションは捨てない
  if (res.status === 401 && store.getState().session.token === token) {
    store.dispatch(logout());
  }
  if (!res.ok) {
    throw new ApiError(path, res.status);
  }
  return res.json();
}

export type RequestLog = {
  method: string;
  path: string;
  status: number;
  ms: number;
  version: string;
};

// argo-rollouts 章で、どのバージョンが応答したかを LivePanel に出す。直近のものだけ持つ。
// Canary の途中で再読み込みしても続きから見られるよう、sessionStorage に残す
let recent: RequestLog[] = JSON.parse(
  sessionStorage.getItem("recentRequests") ?? "[]",
);
const listeners = new Set<() => void>();

function record(log: RequestLog) {
  recent = [log, ...recent].slice(0, 10);
  sessionStorage.setItem("recentRequests", JSON.stringify(recent));
  for (const l of listeners) l();
}

export const recentRequests = {
  subscribe(l: () => void) {
    listeners.add(l);
    return () => listeners.delete(l);
  },
  get: () => recent,
};
