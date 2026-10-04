/**
 * Habarchy SDK — minimal client for the public API.
 * Works in Node ≥ 20 (global fetch + Web Crypto). Keep the API key on the server.
 */

export type Channel = "sms" | "email" | "push" | "telegram" | "auto";
export type Priority = "high" | "normal" | "low";

/** "to": an address string, or a contact reference. */
export type Recipient = string | { contact_id?: string; external_id?: string; address?: string };

export interface SendMessageRequest {
  channel: Channel;
  to: Recipient;
  template?: string;
  data?: Record<string, unknown>;
  subject?: string;
  title?: string;
  body?: string;
  locale?: "tk" | "ru" | "en";
  scheduled_at?: string | Date;
  priority?: Priority;
  idempotency_key?: string;
  metadata?: Record<string, unknown>;
}

export interface Accepted {
  id: string;
  status: string;
  channel: string;
  to: string;
  scheduled_at?: string | null;
  /** true when the idempotency key was already used (HTTP 200). */
  duplicate: boolean;
}

export interface Message {
  id: string;
  status: "queued" | "processing" | "sent" | "delivered" | "failed" | "cancelled";
  channel: string;
  to: string;
  contact_id?: string;
  batch_id?: string;
  template?: string;
  subject?: string;
  body: string;
  priority: string;
  provider_id?: string;
  provider_message_id?: string;
  error_code?: string;
  error_message?: string;
  attempts: number;
  scheduled_at?: string | null;
  sent_at?: string | null;
  delivered_at?: string | null;
  cost_micros: number;
  currency: string;
  metadata: unknown;
  is_test: boolean;
  created_at: string;
  updated_at: string;
}

export interface MessageEvent {
  id: string;
  type: string;
  provider_id?: string;
  payload: unknown;
  created_at: string;
}

export interface MessageDetail {
  message: Message;
  events: MessageEvent[];
}

export interface SendBatchRequest {
  channel: Channel;
  template?: string;
  subject?: string;
  title?: string;
  body?: string;
  locale?: "tk" | "ru" | "en";
  scheduled_at?: string | Date;
  priority?: Priority;
  idempotency_key?: string;
  metadata?: Record<string, unknown>;
  recipients: { to: Recipient; data?: Record<string, unknown> }[];
}

export interface Batch {
  id: string;
  status: string;
  channel: string;
  template?: string;
  total: number;
  queued: number;
  sent: number;
  delivered: number;
  failed: number;
  created_at: string;
  completed_at?: string | null;
  accepted?: number;
  rejected?: number;
}

export interface SendOTPRequest {
  channel?: "sms" | "email" | "telegram";
  to: string;
  length?: number;
  ttl?: number;
  template?: string;
  locale?: "tk" | "ru" | "en";
  data?: Record<string, unknown>;
}

export interface OTPSent {
  message_id: string;
  to: string;
  channel: string;
  expires_in: number;
  length: number;
}

export interface OTPVerified {
  verified: boolean;
  attempts_remaining: number;
}

export interface RegisterDeviceRequest {
  token: string;
  platform: "android" | "ios" | "web";
  app_version?: string;
  contact_id?: string;
  external_id?: string;
}

export interface Device {
  id: string;
  contact_id?: string;
  platform: string;
  token_hint: string;
  app_version: string;
  is_active: boolean;
  last_seen_at: string;
  created_at: string;
}

export class HabarchyError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string,
    public readonly details?: Record<string, unknown>,
  ) {
    super(`habarchy: ${code} (${status}): ${message}`);
    this.name = "HabarchyError";
  }
}

export interface ClientOptions {
  /** Add X-Timestamp / X-Signature to every request. */
  sign?: boolean;
  /** Custom fetch (tests, polyfills). */
  fetch?: typeof fetch;
  /** Request timeout in ms (default 20000). */
  timeoutMs?: number;
  /** Clock override for signing (tests). */
  now?: () => Date;
}

const enc = new TextEncoder();

function hex(buf: ArrayBuffer): string {
  return Array.from(new Uint8Array(buf), (b) => b.toString(16).padStart(2, "0")).join("");
}

/** hex(HMAC-SHA256(apiKey, ts "\n" METHOD "\n" path "\n" hex(sha256(body)))) */
export async function signature(apiKey: string, timestamp: string, method: string, path: string, body: string): Promise<string> {
  const subtle = globalThis.crypto.subtle;
  const bodyHash = hex(await subtle.digest("SHA-256", enc.encode(body)));
  const key = await subtle.importKey("raw", enc.encode(apiKey), { name: "HMAC", hash: "SHA-256" }, false, ["sign"]);
  const msg = `${timestamp}\n${method.toUpperCase()}\n${path}\n${bodyHash}`;
  return hex(await subtle.sign("HMAC", key, enc.encode(msg)));
}

export class Habarchy {
  private readonly baseUrl: string;
  private readonly fetchImpl: typeof fetch;
  private readonly sign: boolean;
  private readonly timeoutMs: number;
  private readonly now: () => Date;

  constructor(baseUrl: string, private readonly apiKey: string, opts: ClientOptions = {}) {
    this.baseUrl = baseUrl.replace(/\/+$/, "");
    this.fetchImpl = opts.fetch ?? globalThis.fetch.bind(globalThis);
    this.sign = opts.sign ?? false;
    this.timeoutMs = opts.timeoutMs ?? 20_000;
    this.now = opts.now ?? (() => new Date());
  }

  /** Queue one message (HTTP 202). */
  async sendMessage(req: SendMessageRequest): Promise<Accepted> {
    const { data, meta } = await this.request<Omit<Accepted, "duplicate">>("POST", "/api/v1/messages", normalizeDates(req));
    return { ...data, duplicate: meta?.duplicate === true };
  }

  async sendBatch(req: SendBatchRequest): Promise<Batch> {
    const { data, meta } = await this.request<Batch>("POST", "/api/v1/messages/batch", normalizeDates(req));
    return { ...data, accepted: num(meta?.accepted), rejected: num(meta?.rejected) };
  }

  async getMessage(id: string): Promise<MessageDetail> {
    return (await this.request<MessageDetail>("GET", `/api/v1/messages/${encodeURIComponent(id)}`)).data;
  }

  async cancelMessage(id: string): Promise<Message> {
    return (await this.request<Message>("POST", `/api/v1/messages/${encodeURIComponent(id)}/cancel`)).data;
  }

  async getBatch(id: string): Promise<Batch> {
    return (await this.request<Batch>("GET", `/api/v1/batches/${encodeURIComponent(id)}`)).data;
  }

  async sendOTP(req: SendOTPRequest): Promise<OTPSent> {
    return (await this.request<OTPSent>("POST", "/api/v1/otp/send", req)).data;
  }

  async verifyOTP(to: string, code: string): Promise<OTPVerified> {
    return (await this.request<OTPVerified>("POST", "/api/v1/otp/verify", { to, code })).data;
  }

  async registerDevice(req: RegisterDeviceRequest): Promise<Device> {
    return (await this.request<Device>("POST", "/api/v1/devices", req)).data;
  }

  private async request<T>(method: string, path: string, body?: unknown): Promise<{ data: T; meta?: Record<string, unknown> }> {
    const payload = body === undefined ? "" : JSON.stringify(body);
    const headers: Record<string, string> = { "X-Api-Key": this.apiKey, Accept: "application/json", "User-Agent": "habarchy-ts/1.0" };
    if (payload) headers["Content-Type"] = "application/json";
    if (this.sign) {
      const ts = Math.floor(this.now().getTime() / 1000).toString();
      headers["X-Timestamp"] = ts;
      headers["X-Signature"] = await signature(this.apiKey, ts, method, path, payload);
    }
    const ctrl = new AbortController();
    const timer = setTimeout(() => ctrl.abort(), this.timeoutMs);
    let res: Response;
    try {
      res = await this.fetchImpl(this.baseUrl + path, { method, headers, body: payload || undefined, signal: ctrl.signal });
    } finally {
      clearTimeout(timer);
    }
    const text = await res.text();
    let env: { data?: T; meta?: Record<string, unknown>; error?: { code: string; message: string; details?: Record<string, unknown> } } = {};
    if (text) {
      try {
        env = JSON.parse(text);
      } catch {
        if (!res.ok) throw new HabarchyError(res.status, "http_error", text.slice(0, 200));
        throw new HabarchyError(res.status, "bad_response", "response is not JSON");
      }
    }
    if (!res.ok) {
      const e = env.error ?? { code: "http_error", message: res.statusText };
      throw new HabarchyError(res.status, e.code, e.message, e.details);
    }
    return { data: env.data as T, meta: env.meta };
  }
}

function num(v: unknown): number | undefined {
  return typeof v === "number" ? v : undefined;
}

function normalizeDates<T extends { scheduled_at?: string | Date }>(req: T): T {
  if (req.scheduled_at instanceof Date) return { ...req, scheduled_at: req.scheduled_at.toISOString() };
  return req;
}

export default Habarchy;
