# Habarchy SDK — TypeScript / Node

```sh
npm install @habarchy/sdk        # or: npm install ../path/to/habarchy/sdk/ts
```

```ts
import { Habarchy } from "@habarchy/sdk";

const hb = new Habarchy("https://habarchy.example.tm", process.env.HABARCHY_API_KEY!, { sign: true });
const acc = await hb.sendMessage({ channel: "sms", to: "+99365123456", template: "otp", data: { code: "4821", minutes: 5 }, idempotency_key: "order-77" });
const detail = await hb.getMessage(acc.id);           // detail.message.status, detail.events
```

Methods: `sendMessage`, `sendBatch`, `getMessage`, `cancelMessage`, `getBatch`, `sendOTP`,
`verifyOTP`, `registerDevice`. `to` is an address string or `{ external_id }` / `{ contact_id }`;
`channel: "auto"` picks push → sms → email for a contact.

Errors throw `HabarchyError { status, code, message, details }`. Node ≥ 20 (global `fetch`,
Web Crypto); zero runtime dependencies. `{ sign: true }` adds `X-Timestamp` / `X-Signature`
(HMAC-SHA256, see docs/auth.md). Server-side only — never ship the API key to a browser.

```sh
npm install && npm test
```
