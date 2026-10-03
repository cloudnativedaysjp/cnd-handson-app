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
  const res = await fetch(path, { ...init, headers });
  // 待っている間に別のトークンでログインし直していたら、そのセッションは捨てない
  if (res.status === 401 && store.getState().session.token === token) {
    store.dispatch(logout());
  }
  if (!res.ok) {
    throw new ApiError(path, res.status);
  }
  return res.json();
}
