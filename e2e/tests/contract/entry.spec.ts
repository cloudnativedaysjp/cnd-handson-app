// Contract for handson-legacy / handson-modern (#73). Traces: no OTel collector in compose yet,
// so span assertions are skipped until one exists.
import { expect, test } from "@playwright/test";
import { entries } from "./lib/config";
import { jsonLogs } from "./lib/logs";

for (const e of entries) {
  test.describe(`entry ${e.variant}`, () => {
    test("/ returns HTML", async ({ request }) => {
      const res = await request.get(`${e.url}/`);
      expect(res.status()).toBe(200);
      expect(res.headers()["content-type"]).toContain("text/html");
    });

    test(`/color returns ${e.color}`, async ({ request }) => {
      const res = await request.get(`${e.url}/color`);
      expect(res.status()).toBe(200);
      expect((await res.text()).trim()).toBe(e.color);
    });

    test("/healthz returns 200", async ({ request }) => {
      expect((await request.get(`${e.url}/healthz`)).status()).toBe(200);
    });

    test("logs one JSON line per request", async ({ request }) => {
      const since = new Date().toISOString();
      await request.get(`${e.url}/color`);
      await expect(async () => {
        const hit = jsonLogs(e.service, since).filter((l) => l.path === "/color");
        expect(hit).toHaveLength(1);
        expect(hit[0]).toMatchObject({
          variant: e.variant,
          color: e.color,
          method: "GET",
          path: "/color",
          status: 200,
        });
        expect(typeof hit[0].duration_ms).toBe("number");
        expect(hit[0].trace_id).toBeTruthy();
      }).toPass({ timeout: 10_000 });
    });
  });
}
