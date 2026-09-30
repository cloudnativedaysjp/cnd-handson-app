import { defineConfig } from "@playwright/test";

// Runs against an already-running stack (`make up`); it does not start servers itself.
export default defineConfig({
  testDir: "./tests",
  testIgnore: "**/contract/**", // run separately: `pnpm test:contract`
  outputDir: "./test-results",
  reporter: [["list"]],
  use: { baseURL: process.env.BASE_URL ?? "http://localhost:5173" },
});
