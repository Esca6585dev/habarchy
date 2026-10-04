import { NextResponse } from "next/server";
import { serverEnv } from "@/lib/env";
import { setAuthCookies, type Tokens } from "@/lib/auth-cookies";

/** Proxies the login so tokens land in httpOnly cookies, never in JS. */
export async function POST(req: Request) {
  const body = await req.text();
  const upstream = await fetch(`${serverEnv.apiUrl}/api/admin/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Forwarded-For": req.headers.get("x-forwarded-for") ?? "" },
    body,
    cache: "no-store",
  });
  const json = (await upstream.json()) as { data?: { tokens: Tokens; user: unknown }; error?: unknown };
  if (!upstream.ok || !json.data) {
    return NextResponse.json(json, { status: upstream.status });
  }
  const res = NextResponse.json({ data: { user: json.data.user } });
  setAuthCookies(res, json.data.tokens);
  return res;
}
