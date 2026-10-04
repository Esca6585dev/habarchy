import { test } from "node:test";
import assert from "node:assert/strict";
import { createHmac, createHash } from "node:crypto";
import { Habarchy, HabarchyError, signature } from "./index.js";

const KEY = "hb_live_testkey";

function serverSignature(ts: string, method: string, path: string, body: string): string {
  const bodyHash = createHash("sha256").update(body).digest("hex");
  return createHmac("sha256", KEY).update(`${ts}\n${method}\n${path}\n${bodyHash}`).digest("hex");
}

test("signature matches the server algorithm", async () => {
  const body = JSON.stringify({ channel: "sms", to: "+99365123456" });
  assert.equal(await signature(KEY, "1700000000", "post", "/api/v1/messages", body), serverSignature("1700000000", "POST", "/api/v1/messages", body));
  assert.equal(await signature(KEY, "1", "GET", "/api/v1/messages/x", ""), serverSignature("1", "GET", "/api/v1/messages/x", ""));
});

test("sendMessage sends signed request and unwraps the envelope", async () => {
  const seen: { url: string; init: RequestInit }[] = [];
  const fakeFetch: typeof fetch = async (url, init) => {
    seen.push({ url: String(url), init: init ?? {} });
    const headers = init?.headers as Record<string, string>;
    const body = String(init?.body ?? "");
    assert.equal(headers["X-Api-Key"], KEY);
    assert.equal(headers["X-Signature"], serverSignature(headers["X-Timestamp"], "POST", "/api/v1/messages", body));
    return new Response(JSON.stringify({ data: { id: "m1", status: "queued", channel: "sms", to: "+99365123456" }, meta: null, error: null }), { status: 202 });
  };
  const c = new Habarchy("https://h.example.tm/", KEY, { sign: true, fetch: fakeFetch, now: () => new Date(1_700_000_000_000) });
  const acc = await c.sendMessage({ channel: "sms", to: "+99365123456", template: "otp", data: { code: "1234" } });
  assert.equal(acc.id, "m1");
  assert.equal(acc.duplicate, false);
  assert.equal(seen[0].url, "https://h.example.tm/api/v1/messages");
});

test("errors become HabarchyError; batch meta is merged", async () => {
  const fakeFetch: typeof fetch = async (url) => {
    const u = String(url);
    if (u.endsWith("/otp/verify")) return new Response(JSON.stringify({ error: { code: "rate_limited", message: "slow down", details: { retry_after: 30 } } }), { status: 429 });
    if (u.endsWith("/messages/batch")) return new Response(JSON.stringify({ data: { id: "b1", status: "queued", total: 2 }, meta: { accepted: 2, rejected: 0 } }), { status: 202 });
    return new Response("gateway timeout", { status: 504 });
  };
  const c = new Habarchy("https://h.example.tm", KEY, { fetch: fakeFetch });
  await assert.rejects(c.verifyOTP("+993", "0000"), (e: unknown) => e instanceof HabarchyError && e.code === "rate_limited" && e.status === 429 && e.details?.retry_after === 30);
  const b = await c.sendBatch({ channel: "sms", body: "x", recipients: [{ to: "+1" }, { to: { external_id: "u2" } }] });
  assert.equal(b.accepted, 2);
  await assert.rejects(c.getMessage("m1"), (e: unknown) => e instanceof HabarchyError && e.code === "http_error" && e.status === 504);
});
