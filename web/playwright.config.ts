import { defineConfig, devices } from "@playwright/test";

/**
 * E2E against a running stack:
 *   WEB_URL  (default http://localhost:3000)  the Next.js admin
 *   API_URL  (default http://localhost:8080)  the Habarchy API (for seeding)
 *   E2E_EMAIL / E2E_PASSWORD                  an existing admin user
 */
export default defineConfig({
  testDir: "./e2e",
  timeout: 60_000,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "github" : "list",
  use: {
    baseURL: process.env.WEB_URL ?? "http://localhost:3000",
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    locale: "en-US",
  },
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        // PW_CHROMIUM lets CI / sandboxes reuse a preinstalled Chromium
        // instead of downloading one.
        launchOptions: process.env.PW_CHROMIUM ? { executablePath: process.env.PW_CHROMIUM } : {},
      },
    },
  ],
  outputDir: "test-results",
});
