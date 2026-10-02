# Habarchy — build prompt

> Paste this whole file into Claude Code (or another coding agent) opened in the `habarchy` repo.

You are building **Habarchy** (Turkmen "Habarçy" = messenger) from scratch in this repository. It is a self-hosted, multi-tenant **notification gateway**, similar to Twilio + SendGrid + OneSignal combined: client apps call ONE API and Habarchy delivers the message through **SMS, Email, Push (Firebase FCM) or Telegram**, with queues, retries, templates, delivery status, webhooks and an admin panel.

The repo is a monorepo with three apps:

```
habarchy/
├── backend/   # Go 1.23 + Fiber v2 — REST API, workers, providers
├── web/       # Next.js 15 (App Router) + TypeScript + Tailwind + shadcn/ui — admin panel
├── mobile/    # Flutter 3.x — admin/monitoring app + reference client for push tokens
├── sdk/       # go/, ts/, dart/, php/ thin client libraries
├── deploy/    # docker-compose, nginx, env examples
└── docs/      # OpenAPI spec, architecture (Mermaid), Postman collection
```

---

## 1. Backend — Go + Fiber

### Stack
- Go 1.23, Fiber v2, PostgreSQL 16 (pgx + **sqlc** for typed queries), Redis 7 (queues, rate limits, idempotency, OTP store).
- Queue: **asynq** (Redis-based) with separate queues `sms`, `email`, `push`, `telegram`, `webhooks` and priorities high/normal/low. Asynqmon mounted at `/admin/queues` behind admin auth.
- Migrations: **goose** (SQL files in `backend/migrations`).
- Config: env vars via `envconfig`, documented in `.env.example`.
- Logging: `zerolog` structured JSON; request-id middleware; OpenTelemetry traces optional behind a flag.
- Validation: `go-playground/validator`; phone normalization to E.164 with `nyaruka/phonenumbers`.
- OpenAPI 3 spec generated with `swaggo/swag` annotations, served at `/api/docs`.
- Auth:
  - Admin panel users: email + password (argon2id), JWT access (15 min) + refresh (30 d, rotated, stored hashed in DB), optional TOTP 2FA.
  - Client projects: API key in header `X-Api-Key` (prefix `hb_live_` / `hb_test_`, SHA-256 hashed at rest, shown once), optional HMAC-SHA256 request signature `X-Signature` + `X-Timestamp` (5 min tolerance), per-key IP allowlist and scopes.
- Layout: clean architecture
  ```
  backend/
    cmd/api/          # HTTP server
    cmd/worker/       # asynq workers
    cmd/scheduler/    # scheduled messages, quota reset, cleanup
    internal/domain/      # entities, enums, errors
    internal/app/         # use cases (SendMessage, RenderTemplate, VerifyOTP ...)
    internal/ports/       # interfaces: SmsProvider, EmailProvider, PushProvider, Repo, Queue
    internal/adapters/
        http/             # Fiber handlers, middleware, DTOs
        postgres/         # sqlc repos
        redis/
        providers/sms/{httpgeneric,smpp}
        providers/email/smtp
        providers/push/fcm
        providers/telegram
    migrations/
    pkg/                  # reusable helpers
  ```

### Domain model (PostgreSQL)
1. **users** — admin users; **project_members** (role: owner, admin, developer, viewer).
2. **projects** (tenants) — name, slug, status, daily_quota, monthly_quota, webhook_url, webhook_secret (encrypted), default_locale, allowed_ips.
3. **api_keys** — project_id, prefix, key_hash, scopes[], ip_allowlist[], last_used_at, revoked_at.
4. **providers** — project_id, channel (sms/email/push/telegram), type (http_sms, smpp, smtp, fcm, telegram_bot), priority, is_active, credentials JSONB **encrypted with AES-256-GCM** (master key from env/KMS), rate_limit_per_sec.
   - `http_sms`: fully configurable URL, method, headers, body template (Go template with `{{.To}} {{.Text}} {{.Sender}}`), success matcher (status code / JSON path / regex), message-id extractor.
   - `smpp`: real SMPP 3.4 client (bind_transceiver, submit_sm with UCS2 for Cyrillic/Turkmen, long-message concatenation via UDH, deliver_sm for DLR, enquire_link keep-alive, auto-reconnect). Use `linxGnu/gosmpp` or `fiorix/go-smpp`.
   - `smtp`: host, port, TLS mode, user, pass, from name/address, reply-to; HTML + plain text; attachments (base64 or URL).
   - `fcm`: Firebase HTTP v1 with a service-account JSON per project; tokens, topics, data payload, Android/APNs options, batch up to 500; handle `UNREGISTERED` → disable the device token.
   - `telegram_bot`: bot token, chat_id, parse_mode Markdown/HTML, inline buttons optional.
   - **Fallback**: ordered by priority; on retryable error try the next provider.
5. **templates** — project_id, key (`otp`, `welcome`…), channel, locale (tk/ru/en), subject, body with `{{.var}}` placeholders, required_vars[], version, is_active. Keep old versions (**template_versions**).
6. **contacts** — project address book: external_id, phone, email, telegram_chat_id, locale, tags[], attributes JSONB; **devices** — contact_id, platform (android/ios/web), fcm_token (unique), app_version, last_seen_at, is_active.
7. **messages** — the central table: id (UUID v7), project_id, channel, to_address (normalized), contact_id, template_key, template_version, rendered_subject, rendered_body, status enum (`queued, processing, sent, delivered, failed, cancelled`), priority, provider_id, provider_message_id, error_code, error_message, attempts, scheduled_at, sent_at, delivered_at, cost_micros, currency, metadata JSONB, idempotency_key, batch_id. Indexes on (project_id, created_at desc), (status, scheduled_at), unique (project_id, idempotency_key). **Partition by month** on created_at.
8. **message_events** — timeline per message (queued, attempt, provider_response, delivered, failed, webhook_sent).
9. **batches**, **webhook_deliveries** (event, payload, signature, response_code, attempts, next_retry_at), **otp_codes** (in Redis, hashed), **audit_logs**, **usage_daily** (aggregated counts/cost per project/channel/day).

### Public API `/api/v1` (API-key auth, JSON envelope `{ "data", "meta", "error": {"code","message","details"} }`)
- `POST /messages` — body: `channel` (`sms|email|push|telegram|auto`), `to` (string or `{contact_id}`), `template` + `data` **or** `body`/`subject`, `locale`, `scheduled_at`, `priority`, `idempotency_key`, `metadata`. Returns **202** with `{id, status}`.
  - `auto` picks the best available channel for the contact in the project's configured order (default push → sms → email).
- `POST /messages/batch` — up to 1000 recipients, same template, per-recipient `data`. Returns batch id; `GET /batches/{id}` for progress.
- `GET /messages/{id}` — status + events timeline. `GET /messages?status=&channel=&from=&to=&cursor=&limit=` cursor pagination.
- `POST /messages/{id}/cancel` — only while `queued`/scheduled.
- `POST /otp/send` `{channel, to, length?, ttl?, template?}` and `POST /otp/verify` `{to, code}` — numeric codes, configurable length/TTL/max attempts, rate-limited per address and per IP, codes stored hashed in Redis, never logged.
- `GET/POST/PUT/DELETE /templates`, `/contacts`, `POST /devices` (register FCM token), `DELETE /devices/{token}`.
- `GET /usage?from=&to=&group_by=channel|day`.
- Inbound: `POST /callbacks/sms/{provider_id}` (HTTP DLR), FCM feedback handling, Telegram webhook optional.
- Errors: `validation_failed`, `unauthorized`, `forbidden_scope`, `quota_exceeded` (429), `rate_limited` (429), `duplicate_request` (returns original message), `provider_unavailable`, `invalid_recipient`.

### Admin API `/api/admin` (JWT) — used by web + mobile
- auth (login, refresh, logout, 2FA), me.
- projects CRUD + members, api-keys (create → plaintext once, revoke), providers CRUD + **test-send**, templates CRUD + **preview with sample data**, contacts/devices, messages log + timeline + raw provider response + **resend**, webhooks log + resend, usage/cost reports, dashboard stats (sent/delivered/failed per day per channel, p95 delivery latency), system health (queue depth, failed tasks, provider status), audit log.
- Server-Sent Events `GET /api/admin/stream` for live message status on the dashboard.

### Reliability rules
- Every send is an asynq task; retry 3× with backoff 10 s / 60 s / 300 s; non-retryable provider errors fail immediately with a clear `error_code`.
- Per-project and per-provider rate limiting with Redis token bucket (SMPP throughput N msg/s).
- Quotas enforced at enqueue time; scheduler resets daily counters and aggregates `usage_daily`.
- Idempotency key valid 24 h.
- Webhooks: HMAC-SHA256 signature header `X-Habarchy-Signature`, events `message.sent|delivered|failed`, `batch.completed`; retry with exponential backoff up to 8 attempts.
- Graceful shutdown, health endpoints `/healthz` and `/readyz`, Prometheus metrics at `/metrics`.

### Tests
- Unit tests for template rendering, phone normalization, HMAC, provider adapters (with httptest / fake SMPP server).
- Integration tests with **testcontainers-go** (Postgres + Redis) for the send flow: enqueue → worker → status → webhook.
- `make test`, `make lint` (golangci-lint), `make run`, `make migrate`.

---

## 2. Web admin — Next.js

- Next.js 15 App Router, TypeScript strict, Tailwind, **shadcn/ui**, TanStack Query, react-hook-form + zod, Recharts for charts, next-intl for **tk (default) / ru / en**.
- Auth: JWT from admin API stored in httpOnly cookie via a Next route handler; refresh handled in middleware; protected `(dashboard)` route group.
- Pages: login, dashboard (charts + live SSE feed), projects (list/create/settings/members), API keys (create dialog shows key once with copy button), providers (per-channel forms with dynamic fields per provider type, test-send), templates (editor with variables sidebar, live preview, version history), contacts & devices, messages log (filters, infinite scroll, detail drawer with timeline and raw response, resend), webhooks log, usage & cost report (export CSV), system health, audit log, profile/2FA.
- Generate the API client from the backend OpenAPI spec (`openapi-typescript` + `openapi-fetch`); do not hand-write fetchers.
- Dark/light theme, responsive, accessible (keyboard, aria). Clean, professional dashboard look; no placeholder lorem ipsum.
- Playwright e2e for login + send-test-message flow.

---

## 3. Mobile — Flutter

- Flutter 3.x, Dart 3, **Riverpod** for state, **go_router**, **dio** with interceptors (JWT refresh), **freezed + json_serializable** models (generate from the OpenAPI spec with `openapi-generator` dart-dio, or hand-write DTOs if generation is unstable).
- Localization tk/ru/en with `flutter_localizations` + ARB files.
- Screens: login (+2FA), project switcher, dashboard (today's counts, failed alerts, mini charts with fl_chart), messages log with filters + detail timeline + resend, templates list + preview, providers status + test-send, API keys (view/revoke), settings.
- **Push demo mode**: the app registers its own FCM token with `POST /api/v1/devices` of a chosen project and shows incoming pushes, so it doubles as a reference client for `firebase_messaging` integration (foreground, background, tap handling, Android notification channel, iOS APNs setup notes).
- Material 3, adaptive layout for tablets, offline-friendly caching of the last dashboard.
- Build scripts: `flutter build apk --release` with flavors `dev`/`prod` and `--dart-define=API_URL=...`; document signing and FCM `google-services.json` setup.

---

## 4. SDKs, docs, deploy

- `sdk/go`, `sdk/ts`, `sdk/dart`, `sdk/php` — minimal clients: `SendMessage`, `SendBatch`, `GetMessage`, `SendOTP`, `VerifyOTP`, `RegisterDevice`, with HMAC signing support. Each has a README with a 5-line example. The PHP one must be usable from a Laravel 12 project (tds.gov.tm) via a Composer path repository.
- `docs/architecture.md` with Mermaid diagrams (components, send sequence, fallback, webhook retry), `docs/openapi.yaml`, Postman collection, `docs/providers.md` explaining how to add a new SMS provider in one Go file.
- `deploy/docker-compose.yml`: api, worker, scheduler, web, postgres, redis, nginx (TLS via certbot optional), asynqmon. One command: `docker compose up -d`. Seed command creates a demo project, prints an API key, and installs sample templates (`otp`, `welcome`, `password_reset`) in tk/ru/en.
- GitHub Actions: backend lint+test, web lint+build, flutter analyze+test, docker image build on tags.
- Root `README.md` in **Turkmen and English**: what Habarchy is, architecture picture, quick start, curl examples for every channel, how tds.gov.tm / Flutter apps integrate.

---

## Working rules
1. Work in steps and **run tests after each step**, showing me the result:
   1. monorepo skeleton, backend config, migrations, domain models, sqlc;
   2. auth (admin JWT + API keys), projects, templates;
   3. messages API + asynq workers + providers (http_sms, smtp, fcm first; smpp and telegram next) + fallback + webhooks + OTP;
   4. admin API + SSE + usage aggregation;
   5. Next.js admin panel;
   6. Flutter app;
   7. SDKs, docs, docker, CI, README.
2. Commit after every completed step with clear English messages; remote is `https://github.com/Esca6585dev/habarchy`, branch `main`.
3. Never log secrets, OTP codes, API keys or provider credentials. Encrypt credentials at rest.
4. Do not invent provider APIs: implement real SMPP 3.4 and real FCM HTTP v1; keep the generic HTTP SMS provider fully configurable from the admin panel.
5. Prefer small, typed, tested code over frameworks-on-frameworks. Explain briefly what you are about to create before creating it.

Start with step 1 now.
