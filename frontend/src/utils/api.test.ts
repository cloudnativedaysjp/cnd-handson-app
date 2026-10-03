import { afterEach, describe, expect, test, vi } from "vitest";
import { login } from "../features/sessions/slice";
import { store } from "../store";
import { api } from "./api";

describe("api", () => {
  afterEach(() => vi.unstubAllGlobals());

  test("sends the token and logs out on 401", async () => {
    store.dispatch(login({ token: "tok" }));
    const fetchMock = vi
      .fn()
      .mockResolvedValue(new Response("", { status: 401 }));
    vi.stubGlobal("fetch", fetchMock);

    await expect(api("/api/projects")).rejects.toThrow("401");
    expect(fetchMock.mock.calls[0][1].headers.Authorization).toBe("Bearer tok");
    expect(store.getState().session.isLoggedIn).toBe(false);
  });
});
