import net from "node:net";
import { expect, test } from "@playwright/test";

test("frontend page loads", async ({ page }) => {
  const res = await page.goto("/");
  expect(res?.ok()).toBe(true);
  await expect(page.locator("#root")).not.toBeEmpty();
  await page.screenshot({ path: "screenshots/smoke.png", fullPage: true });
});

// Only the services implemented in docker-compose.yaml; contract tests come in #71.
const grpcPorts = { user: 50051, session: 50052, project: 50053 };

for (const [name, port] of Object.entries(grpcPorts)) {
  test(`${name} gRPC port is reachable`, async () => {
    const connect = () =>
      new Promise<void>((resolve, reject) => {
        const sock = net.connect(port, "127.0.0.1", () => {
          sock.destroy();
          resolve();
        });
        sock.setTimeout(3000, () => sock.destroy(new Error("timeout")));
        sock.on("error", reject);
      });
    await expect(connect).toPass({ timeout: 30_000 });
  });
}
