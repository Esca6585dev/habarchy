# Authentication

Habarchy has two audiences with two different schemes.

| Audience | Scheme | Header |
|----------|--------|--------|
| Admin panel (web, mobile) | JWT access token + rotated refresh token | `Authorization: Bearer <access>` |
| Client applications | API key, optional HMAC request signature | `X-Api-Key`, `X-Signature`, `X-Timestamp` |

## Admin users

- Passwords are hashed with **argon2id** (64 MiB, t=3, p=2) in PHC format.
- `POST /api/admin/auth/login` → `{ tokens: { access_token, refresh_token, expires_in }, user }`.
  Access tokens are HS256 JWTs valid 15 minutes; refresh tokens are random 256-bit
  values valid 30 days and stored only as SHA-256 hashes.
- `POST /api/admin/auth/refresh` rotates the refresh token. Presenting an already
  rotated token is treated as theft: **every** session of that user is revoked.
- `POST /api/admin/auth/logout` revokes one refresh token; `/auth/logout-all` revokes all.
- Changing the password (`PUT /api/admin/me/password`) revokes all sessions.
- **2FA (TOTP):** `POST /me/totp/setup` returns a secret and `otpauth://` URL (show as QR),
  `POST /me/totp/confirm {code}` enables it, `POST /me/totp/disable {password}` turns it off.
  When enabled, login without a code answers `401` with `error.details.totp_required = true`;
  retry with `totp_code`.

Project access is role based: `owner` > `admin` > `developer` > `viewer`.

| Action | Minimum role |
|--------|--------------|
| View project, members, templates | viewer |
| Create / edit templates, list API keys | developer |
| Edit project settings, members, create / revoke API keys, providers | admin |
| Grant owner role, delete project | owner |

A project a user is not a member of answers `404`, never `403`, so project ids cannot be probed.

## API keys

Created in the admin panel (`POST /api/admin/projects/{id}/api-keys`). The plaintext
is returned **once**; only its SHA-256 hash is stored.

```
hb_live_<43 url-safe base64 chars>    reaches real providers
hb_test_<43 url-safe base64 chars>    messages are marked is_test and never billed
```

Each key carries:

- `scopes`: `messages:send`, `messages:read`, `otp`, `templates`, `contacts`, `devices`, `usage`
  (all by default). A missing scope answers `403 forbidden_scope`.
- `ip_allowlist`: IPs or CIDRs; the project's `allowed_ips` applies on top.
- `expires_at`, `revoked_at`.
- `require_signature`: when true every request must be signed (below).

## HMAC request signature

Optional per request, mandatory for keys with `require_signature`. The secret is the
API key itself, so nothing extra has to be exchanged.

```
X-Timestamp: <unix seconds>
X-Signature: hex( HMAC-SHA256( api_key, message ) )

message = timestamp + "\n" + METHOD + "\n" + path + "\n" + hex(sha256(body))
```

- `METHOD` upper case, `path` without query string (e.g. `/api/v1/messages`).
- `body` is the raw request bytes; for GET it is empty, so `hex(sha256(""))`.
- The timestamp must be within ±5 minutes (`HABARCHY_SIGNATURE_TOLERANCE`) of server time.

PHP example (Laravel):

```php
$ts   = (string) time();
$body = json_encode($payload);
$msg  = "$ts\nPOST\n/api/v1/messages\n" . hash('sha256', $body);
$sig  = hash_hmac('sha256', $msg, $apiKey);
Http::withHeaders(['X-Api-Key' => $apiKey, 'X-Timestamp' => $ts, 'X-Signature' => $sig])
    ->withBody($body, 'application/json')->post("$base/api/v1/messages");
```

The Go reference implementation is `projects.Sign` in
`backend/internal/app/projects/service.go`; the SDKs (step 7) mirror it.

## Webhook signatures (Habarchy → your server)

Every webhook POST carries:

```
X-Habarchy-Event: message.delivered
X-Habarchy-Delivery: <delivery uuid>
X-Habarchy-Timestamp: 1700000000
X-Habarchy-Signature: t=1700000000,v1=<hex HMAC-SHA256(webhook_secret, "1700000000." + raw_body)>
```

Verify in PHP:

```php
[$t, $v1] = [substr($parts[0], 2), substr($parts[1], 3)];   // from "t=...,v1=..."
$expected = hash_hmac('sha256', "$t.$rawBody", $webhookSecret);
abort_unless(hash_equals($expected, $v1) && abs(time() - (int) $t) < 300, 401);
```

The secret is set with `PATCH /api/admin/projects/{id}` `{ "webhook_secret": "..." }` and
stored encrypted. Without a secret the signature is still present (empty key).

## Error envelope

```json
{ "error": { "code": "unauthorized", "message": "invalid api key", "details": {} } }
```

Codes used by auth: `unauthorized` (401), `forbidden_scope` (403), `validation_failed` (400),
`not_found` (404), `conflict` (409).
