import { test, expect, type Page } from "@playwright/test";

const API = process.env.API_URL ?? "http://localhost:8080";
const EMAIL = process.env.E2E_EMAIL ?? "admin@habarchy.tm";
const PASSWORD = process.env.E2E_PASSWORD ?? "Admin-pass-123";
const shots = process.env.E2E_SHOTS; // directory for documentation screenshots

async function shot(page: Page, name: string) {
  if (shots) await page.screenshot({ path: `${shots}/${name}.png`, fullPage: false });
}

test.describe.configure({ mode: "serial" });

test("login, create project and key, send a test message, see it delivered", async ({ page, request }) => {
  // ---- login ----
  await page.goto("/login");
  await expect(page.getByText(/sign in to habarchy|habarçy-a giriň|вход в habarchy/i)).toBeVisible();
  await shot(page, "01-login");
  await page.getByLabel(/email|e-mail/i).fill(EMAIL);
  await page.getByLabel(/password|parol|пароль/i).fill(PASSWORD);
  await page.getByRole("button", { name: /sign in|gir|войти/i }).click();
  await page.waitForURL((u) => !u.pathname.startsWith("/login"));

  // ---- project ----
  const name = `E2E ${Date.now().toString(36)}`;
  await page.goto("/projects");
  await page.getByTestId("new-project").click();
  await page.getByTestId("project-name").fill(name);
  await page.getByTestId("create-project").click();
  await page.waitForURL("/");
  await expect(page.getByRole("heading", { name })).toBeVisible();
  await shot(page, "02-dashboard");

  // ---- template ----
  await page.goto("/templates");
  await page.getByTestId("new-template").click();
  await page.locator("#tpl-key").fill("welcome");
  await page.locator("#tpl-body").fill("Salam {{.name}}! Habarchy işleýär.");
  await page.getByRole("button", { name: /^create$|^döret$|^создать$/i }).click();
  await page.waitForURL(/\/templates\/[0-9a-f-]+/);
  await expect(page.locator("#body")).toHaveValue(/Salam/);
  await page.locator("#sample").fill('{"name":"Aman"}');
  await expect(page.getByText("Salam Aman! Habarchy işleýär.")).toBeVisible({ timeout: 10_000 });
  await shot(page, "03-template-editor");

  // ---- api key (test key -> sandbox provider) ----
  await page.goto("/api-keys");
  await page.getByTestId("new-api-key").click();
  await page.locator("#k-name").fill("e2e");
  await page.getByRole("button", { name: /^create$|^döret$|^создать$/i }).click();
  const created = page.getByTestId("created-key");
  await expect(created).toBeVisible();
  const key = (await created.locator("code").innerText()).trim();
  expect(key.startsWith("hb_test_")).toBeTruthy();
  await shot(page, "04-api-keys");

  // ---- send through the public API with that key ----
  const send = await request.post(`${API}/api/v1/messages`, {
    headers: { "X-Api-Key": key, "Content-Type": "application/json" },
    data: { channel: "sms", to: "+99365123456", template: "welcome", data: { name: "Aman" }, metadata: { source: "e2e" } },
  });
  expect(send.status()).toBe(202);
  const { data } = (await send.json()) as { data: { id: string } };

  // ---- it shows up in the log, delivered by the sandbox worker ----
  await page.goto("/messages");
  await expect(page.getByTestId("message-row").first()).toBeVisible();
  await page.goto(`/messages/${data.id}`);
  await expect(page.getByText("Salam Aman! Habarchy işleýär.")).toBeVisible();
  await expect(page.getByText(/^delivered$|gowşuryldy|доставлено/i).first()).toBeVisible({ timeout: 20_000 });
  await shot(page, "05-message-detail");

  await page.goto("/messages");
  await shot(page, "06-messages");
  await page.goto("/providers");
  await shot(page, "07-providers");
  await page.goto("/health");
  await shot(page, "08-health");

  // ---- language switch (default is Turkmen) to English and back ----
  await page.getByRole("button", { name: /language|dil|язык/i }).click();
  await page.getByRole("menuitem", { name: /English/ }).click();
  await expect(page.getByRole("heading", { name: /System health/ })).toBeVisible();
  await shot(page, "09-health-en");
  await page.getByRole("button", { name: /language|dil|язык/i }).click();
  await page.getByRole("menuitem", { name: /Türkmençe/ }).click();
  await expect(page.getByRole("heading", { name: /Ulgam ýagdaýy/ })).toBeVisible();
});

test("unauthenticated users are redirected to login", async ({ page }) => {
  await page.context().clearCookies();
  await page.goto("/messages");
  await page.waitForURL(/\/login/);
  await expect(page.getByLabel(/email|e-mail/i)).toBeVisible();
});
