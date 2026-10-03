import { logout } from "../features/sessions/slice";
import { store } from "../store";

// 入口（bff）の REST API を呼ぶ。401 ならセッションを捨て、PrivateLayout がログイン画面に戻す
export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const { token } = store.getState().session;
  const res = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
      ...init.headers,
    },
  });
  if (res.status === 401) {
    store.dispatch(logout());
  }
  if (!res.ok) {
    throw new Error(`${path}: ${res.status}`);
  }
  return res.json();
}
