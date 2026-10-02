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

Planning stage. The full build specification lives in [PROMPT.md](PROMPT.md).

## Quick example (planned API)

```http
POST /api/v1/messages
X-Api-Key: hb_live_xxx
Content-Type: application/json

{ "channel": "sms", "to": "+99365123456", "template": "otp", "data": { "code": "4821" } }
```
