import { NextResponse } from "next/server";
import { serverEnv } from "@/lib/env";
import { setAuthCookies, type Tokens } from "@/lib/auth-cookies";

/** Verifies the phone OTP and, on success, lands tokens in httpOnly cookies. */
export async function POST(req: Request) {
  const body = await req.text();
  const upstream = await fetch(`${serverEnv.apiUrl}/api/auth/phone/verify`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Forwarded-For": req.headers.get("x-forwarded-for") ?? "" },
    body,
    cache: "no-store",
  });
  const json = (await upstream.json()) as { data?: { tokens: Tokens; user: unknown; created?: boolean }; error?: unknown };
  if (!upstream.ok || !json.data) {
    return NextResponse.json(json, { status: upstream.status });
  }
  const res = NextResponse.json({ data: { user: json.data.user, created: json.data.created } });
  setAuthCookies(res, json.data.tokens);
  return res;
}
