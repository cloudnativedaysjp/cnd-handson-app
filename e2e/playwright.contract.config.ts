import { defineConfig } from "@playwright/test";

// Contract tests: expected to fail until #65 services exist. Run via `make contract`.
export default defineConfig({
  testDir: "./tests/contract",
  outputDir: "./test-results",
  reporter: [["list"]],
  timeout: 15_000,
});
