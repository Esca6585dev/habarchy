# Habarçy SMS — build prompt (step by step)

> **TK:** Bu faýl "Habarçy SMS" önümini ädim-ädim gurmak üçin doly ýumuş
> sanawy. Her faza aýratyn gurulýar we commit edilýär. Esas: **habarçy**
> (eýýäm multi-tenant, OTP, SMPP 0796, kontakt import, QR gateway, admin panel
> bar) — täze zatlar şonuň üstüne goşulýar, hiç zat ýokdan ýazylmaýar.
>
> **EN:** Step-by-step build list for the "Habarçy SMS" product. Each phase is
> built and committed on its own. Foundation is the existing **habarchy**
> monorepo; extend it, do not rewrite. Reuse the existing projects/contacts/
> groups/providers/messages/gateway code.

## Product in one paragraph

A self-hosted SMS platform for Turkmenistan. The **owner (admin)** runs the
server and sees everything — who is registered, who sent what, to whom. End
users (**clients**) install a client app, sign in with **their phone number +
OTP** (the OTP goes out through the `0796` operator SMPP short code), import
**their own** contacts (Excel / CSV / sqlite / phone), organise them into
groups, and send SMS. Each client's numbers are **isolated** — never merged
with another client's. Sending goes out either through the operator SMPP
(`0796`) or through **phones the client paired** via QR + OTP (an Android
"sender" app with a SIM). A client may be allowed to send from phones shared
by **other** accounts too; sending fans out **100 at a time** and can be
managed/queued across several phones. The server is always listening; every
message is logged with sender, recipient, time and status.

## Apps to ship (4–5)

| # | App | Who | Base |
|---|-----|-----|------|
| 1 | **Admin APK** (owner) | the owner (65656585) | `mobile/` (extend) |
| 2 | **Client APK** | end users | new flavor/app from `mobile/` |
| 3 | **Sender APK** (SMS gateway) | runs on a phone with a SIM | `gateway/` (extend) |
| 4 | **Web panel** | owner + clients in a browser | `web/` (extend) |
| 5 | **Web OTP/landing** (optional) | public sign-in / download page | `web/` route |

## Branding / UI

- Name: **Habarçy SMS**. Logo: speech-bubble + SIM/signal motif; deliver as
  `assets/` (Flutter) + `public/` (web), light & dark. See Phase 6.
- Visual style: **green "chat" look** (WhatsApp/Telegram-like bubbles for the
  compose + message log), rounded cards, tk/ru/en, dark mode. Reuse the
  existing shadcn/Tailwind (web) and theme (mobile); do not restyle from zero.

---

## Phase 0 — 0796 SMPP + OTP  ✅ DONE

Goal: SMS/OTP deliver through the operator SMPP short code (`0796`).

- [x] `cmd/smpptest` — bind/submit verifier reading `SMPP_*` env.
- [x] `-seed` creates an active `smpp` provider from `SMPP_*` (idempotent).
- [x] `otp` templates already seeded → `POST /api/v1/otp/send` goes via 0796.
- [x] Docs: `docs/sms-0796.md`. Secrets only in env (AES-256-GCM at rest).

Acceptance: `make smpptest` prints `✓ bind OK`; `-to <num>` delivers one SMS.

---

## Phase 1 — Client auth: phone + OTP login, owner super-admin, isolation

**TK:** Clientler öz nomeri + OTP bilen girýär; eýe (65656585) super-admin;
nomerler ulanyjylar arasynda garyşmaýar.

What exists: admin users (email+password+TOTP, JWT), API keys + HMAC for
service clients, OTP send/verify as a messaging feature, projects = tenants.

New work:
1. **Phone+OTP end-user identity.** Add an `end_users` concept (or reuse
   `users` with an auth method flag): a client is identified by E.164 phone.
   - `POST /api/auth/phone/start` `{phone}` → sends an OTP via the 0796 SMPP
     provider (reuse `otp.Service`), rate-limited per phone + IP.
   - `POST /api/auth/phone/verify` `{phone, code}` → on success, create the
     user on first login (auto-provision a **personal project/tenant** named
     after the phone) and return JWT access+refresh (reuse `auth` tokens).
2. **Owner super-admin.** Seed/promote the phone **+99365656585** (bare
   `65656585`) as `owner` with a platform-wide role: can list all tenants,
   users, providers, messages. Add a `platform_admin` boolean on the user and
   gate new `/api/admin/platform/*` read endpoints with it.
3. **Tenant isolation (hard guarantee).** Every contact/group/message/device
   query is already scoped by `project_id`; add tests proving a client cannot
   read or send to another tenant's contacts, and that imports land only in
   the caller's tenant. Numbers never merge across tenants.
4. **Permissions.** Reuse `project_members` roles (owner/admin/developer/
   viewer). "Everyone with permission can upload/manage contacts" = members
   with ≥ developer on a tenant. Document the matrix.

Deliverables: migrations (`platform_admin`, phone-login fields), handlers,
`auth` extensions, isolation tests, `docs/auth-phone.md` (tk). Owner login
seeded.

Acceptance: a new phone can request an OTP, verify, and get a token scoped to
its own fresh tenant; owner phone sees the platform list; isolation tests pass.

---

## Phase 2 — Contacts: import (Excel/CSV/sqlite), format guide, groups

**TK:** Excel/CSV/sqlite import; upload etmezden öň "excel nähili durmaly"
nusga görkezmek; create/update grupba; sms-sender-app nomerlerini import.

What exists: contacts + groups CRUD, import for Excel/CSV/Word/vCard/Google/
iCloud (`POST /api/v1/contacts/import`), contact trash/restore, Turkmen phone
normalization, groups in `POST /messages/batch`.

New work:
1. **Format guidance before upload.** Add `GET /api/v1/contacts/import/template`
   returning (a) the required column layout and (b) a downloadable sample
   `.xlsx`/`.csv`. Columns: `phone` (required, +993/national ok), `name`,
   `email`, `tags`, plus free attribute columns. Web + mobile show this
   **before** the file picker, with a "download sample" button and a preview
   of valid vs. invalid/duplicate rows (reuse the existing importer's report).
2. **sqlite import.** Add a sqlite source to the importer: accept a `.db`/
   `.sqlite` upload, let the user pick the table + map the phone/name columns,
   import into the caller's tenant.
3. **Import from sms-sender-app.** One-off migration tool
   (`cmd/import-smssender`) that reads the sms-sender-app backend's
   contacts/history (see that repo's `docs/DEVICES_API.md`) or an exported
   sqlite/Excel and loads numbers into a chosen Habarçy tenant, de-duplicating
   and reporting invalids. Never merges numbers across tenants.
4. **Groups.** Already create/update/delete + members; batch into groups of
   100 is a send concern (Phase 3). Ensure rename + multi-select exist in UI.

Deliverables: template endpoint + sample files, sqlite importer, `cmd/import-
smssender`, web+mobile import screens with the format guide, tests.
`docs/contacts-import.md` updated (tk section + the required format).

Acceptance: uploading the sample imports cleanly; a bad row is reported, not
silently dropped; sqlite import works; sms-sender numbers land in one tenant.

---

## Phase 3 — Phones & sending: QR+OTP pairing, multi-phone, 100/100 queue

**TK:** Client app bilen sender app QR + OTP arkaly birleşýär; client öz
telefonlaryny dolandyrýar; öz + birikdirilen beýleki akauntlardan ugradýar;
100-100 nobat; serwer hemişe diňläp durýar.

What exists: `android_sms` provider + `gateway/` APK, QR pairing payload
(`habarchy://gateway?url=&key=`), outbox long-poll, sent/delivered reports,
heartbeat, inbound forward, provider fallback by priority, asynq queues with
rate limits, `POST /messages/batch` up to 1000.

New work:
1. **QR + OTP pairing from the client app.** The client app creates an
   `android_sms` provider in its tenant and shows the pairing QR; the **sender
   APK** scans it. Add an **OTP confirmation** step on pairing (server issues a
   short code shown in the client app, entered on the phone, or vice versa) so
   a scanned QR alone cannot pair. Reuse `otp.Service` / gateway key hashing.
2. **Manage own phones.** Client screen listing their paired phones (name,
   battery, network, operator, online/offline from heartbeat), with rename,
   disable, re-pair, delete. Reuse `GET /providers/{pid}/pairing`, health.
3. **Send from other accounts' phones (sharing).** A phone owner can **share**
   a sender phone with another account; the grantee can enqueue through it
   (scoped, revocable, audited). New `provider_shares` table + endpoints;
   enforce in the send path and in fallback selection.
4. **100-at-a-time fan-out + management.** Group sends split into batches of
   **100** (like sms-sender-app's groups), round-robin across the selected/
   available phones, with per-batch send/continue/reset, live counters
   (sent / remaining / failed) and a global stop. Reuse asynq + rate limits;
   expose batch progress (`GET /batches/{id}`), pause/resume, and per-phone
   throughput caps.
5. **Always-on server.** Already the worker + long-poll; document the deploy
   (compose) and health so "server listening" is guaranteed; add alerting when
   a phone goes offline mid-batch (fallback to 0796 SMPP or another phone).

Deliverables: pairing-OTP flow, `provider_shares` + endpoints, 100-batch
fan-out + progress/controls, client "my phones" screen, sender-APK QR scan +
OTP, tests. `docs/gateway.md` updated.

Acceptance: client pairs a phone via QR+OTP; a 250-number group sends as 3
batches of 100 across phones with live counters; a shared phone works for the
grantee and can be revoked; a phone dropping offline falls back.

---

## Phase 4 — The apps (admin / client / sender / web)

**TK:** 4–5 app: Admin APK (eýe), Client APK, Sender APK, Web panel, (Web OTP).

1. **Client APK** — new Flutter flavor from `mobile/`: phone+OTP login,
   contacts + import (Phase 2), groups, compose (green chat UI), my phones
   (Phase 3), send + live progress, history, settings (tk/ru/en, dark).
2. **Admin APK (owner)** — extend `mobile/`: platform dashboard (who's
   registered, totals), all-tenants message log (who→whom, when, status),
   providers/health, users. Gated by `platform_admin`.
3. **Sender APK** — extend `gateway/`: QR **scan** (camera) + OTP confirm,
   SIM slot pick, counters, foreground service (already), nicer single-screen
   UI, autostart tips.
4. **Web panel** — extend `web/`: client login (phone+OTP) and owner login,
   same features as the apps, compose + message log in the green chat style.
5. **Web OTP/landing (optional)** — public page: app download links + phone
   sign-in.

Deliverables: build flavors, `--dart-define`/env wiring, store-free APK build
docs, Playwright smoke for web login+send.

Acceptance: each app builds (`flutter build apk`, `npm run build`) and runs
against the API; a client can do login→import→send end to end.

---

## Phase 5 — "Who sent what to whom" (reporting & control)

**TK:** Name ugratdy, kime ugratdy, haçan, näme ýagdaý — görünmeli we
dolandyrylmaly.

What exists: message log + timeline + raw provider response + resend/cancel,
dashboard stats, usage/cost, SSE live feed, audit log, Prometheus.

New work: per-client and platform-wide views that answer "who sent what to
whom" at a glance — sender (account/phone), recipient, group, channel (SMPP
vs. which phone), time, status, cost; filters + CSV export; live feed; owner
sees all tenants, client sees only its own. Reuse the admin message-log and
dashboard endpoints; add the sender/phone dimension and a platform scope.

Acceptance: owner filters the platform log by tenant/phone/status/date and
exports CSV; a client sees only its own sends; every send is attributable.

---

## Phase 6 — Branding: "Habarçy SMS" logo + green chat theme

1. Logo set (SVG + PNG, light/dark, app icons) under `assets/` and
   `web/public/`; wire `flutter_launcher_icons` and web favicons/manifest.
2. Green "chat" theme tokens applied to compose + message log in web + mobile
   (bubbles, status ticks, rounded cards), tk/ru/en, dark mode, accessible.

Acceptance: all apps show the Habarçy SMS logo/icon and the green chat UI,
consistent in light & dark.

---

## Cross-cutting rules (every phase)

- **Isolation first:** every new read/write is scoped by tenant; add a test
  that a cross-tenant access is denied.
- **No secrets in git:** operator/phone credentials only in env or the
  AES-256-GCM provider store; never printed, never logged.
- **Reuse, don't rewrite:** extend existing services/handlers/screens.
- **Each phase:** migrations + backend + web + mobile + tests + a short tk
  doc, built green (`make build`, `go vet`, `flutter analyze`, `npm run
  build`) before commit. One phase = one focused PR/commit set.
- **OTP safety:** codes hashed in Redis, rate-limited, never logged (already).

## Build order

Phase 0 ✅ → 1 → 2 → 3 → 4 → 5 → 6. Phases 1–3 are backend-heavy and unlock
the apps; 4–6 are app/UX. Start each with its migrations + tests, finish with
its tk doc.
