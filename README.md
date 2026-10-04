# Habarçy (Habarchy)

**TK:** Habarçy — bir API bilen SMS, e-mail, push (Firebase FCM) we Telegram habarlaryny iberýän, köp proýektli (multi-tenant), nobatly (queue), gaýtadan synanyşýan (retry/fallback), şablonly we loglary bolan merkezi habar şlýuzy. Twilio + SendGrid + OneSignal ýaly, ýöne öz serweriňizde.

**EN:** Habarchy is a self-hosted, multi-tenant notification gateway: one API to send SMS, email, push (FCM) and Telegram messages with queues, retries, provider fallback, templates, delivery status, webhooks and an admin panel.

## Stack

| Part | Tech |
|------|------|
| `backend/` | Go 1.23 + Fiber v2, PostgreSQL, Redis, asynq |
| `web/` | Next.js 15 + TypeScript + Tailwind + shadcn/ui |
| `mobile/` | Flutter 3 + Riverpod + firebase_messaging |
| `sdk/` | Go, TypeScript, Dart, PHP clients |

## Status

In progress — built step by step from [PROMPT.md](PROMPT.md).

| Step | Scope | Status |
|------|-------|--------|
| 1 | Monorepo skeleton, backend config, migrations, domain models, sqlc | ✅ done |
| 2 | Auth (admin JWT + 2FA, API keys + HMAC), projects, members, templates | ✅ done |
| 3 | Messages API, asynq workers, providers, fallback, webhooks, OTP | ⏳ |
| 4 | Admin API, SSE, usage aggregation | ⏳ |
| 5 | Next.js admin panel | ⏳ |
| 6 | Flutter app | ⏳ |
| 7 | SDKs, docs, docker, CI | ⏳ |

## Development

```sh
# dependencies
docker compose -f deploy/docker-compose.dev.yml up -d   # Postgres 16 + Redis 7
make -C backend tools                                     # sqlc, goose

# backend
cp backend/.env.example backend/.env                      # set JWT_SECRET and MASTER_KEY
cd backend
go run ./cmd/api -genkey                                  # prints a HABARCHY_MASTER_KEY
make migrate                                              # goose up
HABARCHY_ADMIN_EMAIL=you@example.com HABARCHY_ADMIN_PASSWORD='min 8 chars' go run ./cmd/api -create-admin
make run                                                  # API on :8080  (/healthz, /readyz)
make test                                                 # unit tests; set HABARCHY_TEST_DATABASE_URL for integration tests
make lint
```

See [docs/architecture.md](docs/architecture.md) for the component and data-model diagrams and
[docs/auth.md](docs/auth.md) for login, API keys and request signing.

### API surface so far

| Area | Endpoints |
|------|-----------|
| Admin auth | `POST /api/admin/auth/login`, `/refresh`, `/logout`, `/logout-all`, `GET /me`, `PUT /me/password`, `POST /me/totp/{setup,confirm,disable}` |
| Projects | `GET/POST /api/admin/projects`, `GET/PATCH/DELETE /projects/{id}`, `GET/PUT /members`, `DELETE /members/{user_id}` |
| API keys | `GET/POST /projects/{id}/api-keys`, `DELETE /api-keys/{key_id}` |
| Templates (admin) | `GET/POST /projects/{id}/templates`, `GET/PUT/DELETE /{tid}`, `GET /{tid}/versions`, `POST /{tid}/versions/{v}/restore`, `POST /preview`, `POST /{tid}/preview` |
| Templates (public) | `GET/POST /api/v1/templates`, `GET/PUT/DELETE /{tid}`, `POST /preview`, `GET /api/v1/me` (key auth) |

## Quick example (planned API)

```http
POST /api/v1/messages
X-Api-Key: hb_live_xxx
Content-Type: application/json

{ "channel": "sms", "to": "+99365123456", "template": "otp", "data": { "code": "4821" } }
```
