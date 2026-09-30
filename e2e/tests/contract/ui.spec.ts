// UI contract for both entry variants (#73). Selectors are the contract: data-testid values below
// must exist in both variants. The idp seeds the user (CONTRACT_USER_EMAIL/PASSWORD) with id
// USER_ID; the test creates its own project owned by that id.
import { expect, test } from "@playwright/test";
import { cfg, entries } from "./lib/config";
import { client, USER_ID } from "./lib/grpc";

let projectName: string;
test.beforeAll(async () => {
  projectName = `ui-${Date.now()}`;
  await client("project/project.proto", "project", "ProjectService", cfg.projectGrpc)("CreateProject", {
    name: projectName,
    owner_id: USER_ID,
  });
});

for (const e of entries) {
  test(`ui ${e.variant}: login -> projects -> tasks -> create task`, async ({ page }) => {
    await page.goto(`${e.url}/`);
    await page.getByTestId("login-email").fill(cfg.userEmail);
    await page.getByTestId("login-password").fill(cfg.userPassword);
    await page.getByTestId("login-submit").click();

    await expect(page.getByTestId("project-list")).toBeVisible();
    await page.getByTestId("project-item").filter({ hasText: projectName }).click();
    await expect(page.getByTestId("task-list")).toBeVisible();

    const title = `task-${Date.now()}`;
    await page.getByTestId("new-task-title").fill(title);
    await page.getByTestId("new-task-submit").click();
    await expect(page.getByTestId("task-list").getByText(title)).toBeVisible();
  });
}
