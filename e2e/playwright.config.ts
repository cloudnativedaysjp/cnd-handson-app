import { defineConfig } from "@playwright/test";

// Runs against an already-running stack (`make up`); it does not start servers itself.
export default defineConfig({
  testDir: "./tests",
  outputDir: "./test-results",
  reporter: [["list"]],
  use: { baseURL: process.env.BASE_URL ?? "http://localhost:5173" },
});
