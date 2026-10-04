import type { NextResponse } from "next/server";
import { ACCESS_COOKIE, REFRESH_COOKIE } from "@/i18n/config";
import { serverEnv } from "@/lib/env";

export interface Tokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export function setAuthCookies(res: NextResponse, tokens: Tokens) {
  const base = { httpOnly: true, sameSite: "lax" as const, secure: serverEnv.secureCookies, path: "/" };
  res.cookies.set(ACCESS_COOKIE, tokens.access_token, { ...base, maxAge: Math.max(60, tokens.expires_in) });
  res.cookies.set(REFRESH_COOKIE, tokens.refresh_token, { ...base, maxAge: 30 * 24 * 3600 });
}

export function clearAuthCookies(res: NextResponse) {
  res.cookies.set(ACCESS_COOKIE, "", { path: "/", maxAge: 0 });
  res.cookies.set(REFRESH_COOKIE, "", { path: "/", maxAge: 0 });
}

/** Calls the backend refresh endpoint; returns new tokens or null. */
export async function refreshTokens(refreshToken: string): Promise<Tokens | null> {
  const res = await fetch(`${serverEnv.apiUrl}/api/admin/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
    cache: "no-store",
  });
  if (!res.ok) return null;
  const body = (await res.json()) as { data?: { tokens?: Tokens } };
  return body.data?.tokens ?? null;
}
