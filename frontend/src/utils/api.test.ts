import { afterEach, describe, expect, test, vi } from "vitest";
import { login } from "../features/sessions/slice";
import { store } from "../store";
import { ApiError, api, recentRequests } from "./api";

describe("api", () => {
  afterEach(() => vi.unstubAllGlobals());

  test("sends the token and logs out on 401", async () => {
    store.dispatch(login({ token: "tok" }));
    const fetchMock = vi
      .fn()
      .mockResolvedValue(new Response("", { status: 401 }));
    vi.stubGlobal("fetch", fetchMock);

    await expect(api("/api/projects")).rejects.toBeInstanceOf(ApiError);
    expect(fetchMock.mock.calls[0][1].headers.get("Authorization")).toBe(
      "Bearer tok",
    );
    expect(store.getState().session.isLoggedIn).toBe(false);
  });

  test("keeps a newer session when an older token gets 401", async () => {
    store.dispatch(login({ token: "old" }));
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => {
        store.dispatch(login({ token: "new" }));
        return new Response("", { status: 401 });
      }),
    );

    await expect(api("/api/projects")).rejects.toThrow("401");
    expect(store.getState().session.token).toBe("new");
  });

  test("records the version that responded", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response("", {
          status: 502,
          headers: { "X-App-Version": "modern-v2" },
        }),
      ),
    );

    await expect(
      api("/api/projects/p1/tasks", { method: "POST" }),
    ).rejects.toThrow("502");
    expect(recentRequests.get()[0]).toMatchObject({
      method: "POST",
      status: 502,
      version: "modern-v2",
    });
  });
});
