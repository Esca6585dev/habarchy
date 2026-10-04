import { NextResponse } from "next/server";
import { cookies } from "next/headers";
import { REFRESH_COOKIE } from "@/i18n/config";
import { serverEnv } from "@/lib/env";
import { clearAuthCookies } from "@/lib/auth-cookies";

export async function POST() {
  const store = await cookies();
  const refresh = store.get(REFRESH_COOKIE)?.value;
  if (refresh) {
    await fetch(`${serverEnv.apiUrl}/api/admin/auth/logout`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refresh }),
      cache: "no-store",
    }).catch(() => undefined);
  }
  const res = NextResponse.json({ data: { ok: true } });
  clearAuthCookies(res);
  return res;
}
