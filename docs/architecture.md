# Architecture

## Components

```mermaid
flowchart LR
    subgraph Clients
        L[Laravel / tds.gov.tm]
        F[Flutter apps]
        O[Other services]
    end
    subgraph Habarchy
        API[cmd/api<br/>Fiber REST + SSE]
        W[cmd/worker<br/>asynq consumers]
        S[cmd/scheduler<br/>partitions, usage, cleanup]
        PG[(PostgreSQL 16<br/>messages partitioned by month)]
        R[(Redis 7<br/>queues · rate limits · idempotency · OTP)]
    end
    subgraph Providers
        SMS[HTTP SMS / SMPP 3.4]
        MAIL[SMTP]
        FCM[Firebase FCM v1]
        TG[Telegram Bot API]
    end
    ADMIN[Next.js admin] --> API
    MOB[Flutter admin] --> API
    L & F & O -- X-Api-Key --> API
    API --> PG
    API -- enqueue --> R
    R -- dequeue --> W
    W --> PG
    W --> SMS & MAIL & FCM & TG
    W -- webhooks --> L
    S --> PG
```

## Backend layout (clean architecture)

```
backend/
  cmd/{api,worker,scheduler}      entry points
  internal/domain                 entities, enums, error codes (no deps)
  internal/app                    use cases (step 2+)
  internal/ports                  interfaces: providers, queue, cipher
  internal/adapters/http          Fiber server, middleware, JSON envelope
  internal/adapters/postgres      pgx pool, goose migrations, sqlc queries
  internal/adapters/redis         shared Redis client
  internal/adapters/providers     sms/{httpgeneric,smpp} email/smtp push/fcm telegram
  internal/queue                  asynq queue & task names
  internal/scheduler              interval job runner
  migrations/                     goose SQL (embedded)
  pkg/{phone,crypto,ids,logger}   reusable helpers
```

## Data model

```mermaid
erDiagram
    users ||--o{ project_members : has
    projects ||--o{ project_members : has
    projects ||--o{ api_keys : owns
    projects ||--o{ providers : configures
    projects ||--o{ templates : owns
    templates ||--o{ template_versions : history
    projects ||--o{ contacts : "address book"
    contacts ||--o{ devices : "push tokens"
    projects ||--o{ batches : creates
    projects ||--o{ messages : sends
    batches ||--o{ messages : groups
    messages ||--o{ message_events : timeline
    projects ||--o{ webhook_deliveries : notifies
    projects ||--o{ usage_daily : aggregates
    projects ||--o{ audit_logs : records
```

Notes:

- `messages` is range-partitioned by month on `created_at`; the primary key is
  `(id, created_at)`. `ensure_messages_partition(date)` creates partitions and the
  scheduler calls it ahead of time. A `messages_default` partition catches strays.
- Idempotency keys are enforced in `message_idempotency_keys` (true unique
  constraint) plus Redis for 24 h; a unique index on a partitioned table would
  have to include the partition key.
- Provider credentials, webhook secrets and TOTP seeds are AES-256-GCM encrypted
  with `HABARCHY_MASTER_KEY`; API keys are stored as SHA-256 hashes.

## Request authentication

```mermaid
sequenceDiagram
    participant App as Client app
    participant API as cmd/api
    participant DB as PostgreSQL
    App->>API: POST /api/v1/... X-Api-Key (+ X-Timestamp, X-Signature)
    API->>DB: SELECT api_keys WHERE key_hash = sha256(key)
    DB-->>API: key + project
    API->>API: revoked? expired? ip allowed? project active?
    API->>API: verify HMAC(key, ts\nMETHOD\npath\nsha256(body)) if present/required
    API->>API: scope check for the route
    API-->>App: 2xx or {error:{code}}
```

Details in [auth.md](auth.md).

## Send sequence

```mermaid
sequenceDiagram
    participant App as Client app
    participant API as cmd/api
    participant R as Redis
    participant PG as PostgreSQL
    participant W as cmd/worker
    participant P as Provider
    participant H as Webhook URL
    App->>API: POST /api/v1/messages {channel, to, template, data, idempotency_key}
    API->>API: resolve contact · pick channel (auto) · normalize address · render template
    API->>R: INCR quota:{project}:d:{day} (seed from DB on first use)
    API->>R: SETNX idem:{project}:{key}
    API->>PG: INSERT messages (queued) + message_events(queued) + idempotency key
    API->>R: asynq enqueue send:sms (task id = message id, ProcessAt = scheduled_at)
    API-->>App: 202 {id, status: queued}
    R-->>W: dequeue
    W->>PG: UPDATE status = processing, attempts + 1
    W->>R: token bucket wait (provider rate_limit_per_sec)
    W->>P: send
    P-->>W: provider_message_id
    W->>PG: status = sent, events(sent)
    W->>PG: INSERT webhook_deliveries(message.sent) → enqueue webhook:deliver
    W->>H: POST signed JSON (X-Habarchy-Signature)
    P-->>API: DLR callback POST /callbacks/sms/{provider_id}
    API->>PG: status = delivered, events(delivered) → webhook message.delivered
```

## Provider fallback and retries

```mermaid
flowchart TD
    A[task send:*] --> B{providers for channel<br/>ordered by priority}
    B -- none --> F1[failed: provider_unavailable]
    B --> C[try provider i]
    C -- success --> S[sent → webhook message.sent]
    C -- permanent error<br/>invalid_recipient, auth... --> F2[failed immediately<br/>webhook message.failed]
    C -- retryable error --> D{more providers?}
    D -- yes --> C
    D -- no --> E{attempt < max retry?}
    E -- yes --> Q[status back to queued<br/>asynq retry 10s / 60s / 300s]
    E -- no --> F3[failed: retries exhausted]
    Q --> A
```

Test API keys (`hb_test_`) skip the provider list and use the sandbox, which accepts
every message and marks it delivered (addresses containing `FAIL` / `RETRY` force errors).

## Webhook delivery and retry

```mermaid
sequenceDiagram
    participant W as worker
    participant PG as PostgreSQL
    participant H as Webhook URL
    W->>PG: INSERT webhook_deliveries (event, payload, url)
    loop up to 8 attempts
        W->>H: POST payload · X-Habarchy-Signature: t=ts,v1=hmac(secret, ts.body)
        alt 2xx
            H-->>W: ok
            W->>PG: delivered_at = now()
        else error / timeout
            W->>PG: attempts + 1, response_code, next_retry_at
            Note over W: backoff 10s · 30s · 1m · 5m · 15m · 1h · 3h
        end
    end
```

Events: `message.sent`, `message.delivered`, `message.failed`, `batch.completed`.
Receivers verify with HMAC-SHA256 over `"<timestamp>.<raw body>"` using the project's
webhook secret and reject timestamps older than a few minutes (see [auth.md](auth.md)).
