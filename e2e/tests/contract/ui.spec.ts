// UI contract for both entry variants (#73). Selectors are the contract: data-testid values below
// must exist in both variants. The user is seeded by the idp (CONTRACT_USER_EMAIL/PASSWORD).
import { expect, test } from "@playwright/test";
import { cfg, entries } from "./lib/config";

for (const e of entries) {
  test(`ui ${e.variant}: login -> projects -> tasks -> create task`, async ({ page }) => {
    await page.goto(`${e.url}/`);
    await page.getByTestId("login-email").fill(cfg.userEmail);
    await page.getByTestId("login-password").fill(cfg.userPassword);
    await page.getByTestId("login-submit").click();

    await expect(page.getByTestId("project-list")).toBeVisible();
    await page.getByTestId("project-item").first().click();
    await expect(page.getByTestId("task-list")).toBeVisible();

    const title = `task-${Date.now()}`;
    await page.getByTestId("new-task-title").fill(title);
    await page.getByTestId("new-task-submit").click();
    await expect(page.getByTestId("task-list").getByText(title)).toBeVisible();
  });
}
