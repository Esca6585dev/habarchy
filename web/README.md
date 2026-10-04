# Habarchy web admin

Next.js 15 (App Router, React 19) + TypeScript + Tailwind v4 + Radix-based UI (shadcn/ui
style, components vendored in `src/components/ui`), TanStack Query, react-hook-form + zod,
Recharts, next-intl (**tk** default / ru / en), next-themes.

## How it talks to the API

- The browser never sees tokens. `POST /api/auth/login` (a Next route handler) calls the
  backend and stores the JWT pair in **httpOnly cookies** (`habarchy_access`, `habarchy_refresh`).
- Client code calls `/api/backend/<backend path>`; the proxy route attaches
  `Authorization: Bearer`, refreshes an expired access token once (rotating the refresh
  cookie) and streams responses back — including the SSE live feed.
- `src/middleware.ts` redirects unauthenticated visitors to `/login`.
- The API client is generated from the OpenAPI spec: `npm run openapi` →
  `src/lib/api/schema.d.ts`, used through `openapi-fetch` (`src/lib/api/client.ts`). No hand-written fetchers.
- `/api-docs` (Swagger UI) and `/admin/queues` (asynqmon) are rewritten to the backend so
  the admin cookie authenticates them.

## Pages

login · dashboard (counters, per-day chart, by channel, live SSE feed, recent failures) ·
projects (list / create / switch) · settings (quotas, webhook URL + secret, IP allowlist, auto
channel order, members, delete) · API keys (create shows the key once, revoke) · providers
(per-type JSON credentials with samples, test-send, masked settings, DLR URL) · templates
(editor with variables sidebar, live preview, version history, restore) · contacts & devices ·
messages (filters, infinite scroll, detail with timeline / raw provider response / webhooks,
resend, cancel) · webhooks (log, payload, resend) · usage (group by, CSV export) · system
health · audit log · profile (password, TOTP 2FA, sign out everywhere).

Roles drive what is shown: viewers read, developers edit templates and resend, admins manage
providers / keys / members, owners delete projects.

## Develop

```sh
cp .env.example .env.local            # HABARCHY_API_URL=http://localhost:8080
npm install
npm run openapi                       # regenerate types after backend spec changes
npm run dev                           # http://localhost:3000
npm run lint && npm run typecheck && npm run build
```

## E2E (Playwright)

Needs a running API + worker and an admin user (`go run ./cmd/api -create-admin`):

```sh
API_URL=http://localhost:8080 WEB_URL=http://localhost:3000 \
E2E_EMAIL=admin@habarchy.tm E2E_PASSWORD='Admin-pass-123' npm run e2e
# PW_CHROMIUM=/path/to/chrome reuses a preinstalled browser; E2E_SHOTS=dir saves screenshots
```

The flow: login → create project → create template (live preview) → create a test API key →
send a message through the public API with that key → the message appears in the log and is
delivered by the sandbox worker → language switch.
