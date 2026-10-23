import { NextResponse } from "next/server";
import { serverEnv } from "@/lib/env";

/** Proxies the phone-login OTP request (no cookie; code is sent by SMS). */
export async function POST(req: Request) {
  const body = await req.text();
  const upstream = await fetch(`${serverEnv.apiUrl}/api/auth/phone/start`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Forwarded-For": req.headers.get("x-forwarded-for") ?? "" },
    body,
    cache: "no-store",
  });
  const json = (await upstream.json()) as unknown;
  return NextResponse.json(json, { status: upstream.status });
}
