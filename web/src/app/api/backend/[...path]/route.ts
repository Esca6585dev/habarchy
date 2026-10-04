import { NextRequest, NextResponse } from "next/server";
import { ACCESS_COOKIE, REFRESH_COOKIE } from "@/i18n/config";
import { serverEnv } from "@/lib/env";
import { refreshTokens, setAuthCookies, clearAuthCookies } from "@/lib/auth-cookies";

export const dynamic = "force-dynamic";

/**
 * Transparent proxy to the Habarchy API for browser code. Adds the
 * Authorization header from the httpOnly cookie, refreshes an expired
 * access token once, and streams responses (SSE included) back.
 */
async function handle(req: NextRequest, ctx: { params: Promise<{ path: string[] }> }) {
  const { path } = await ctx.params;
  const target = `${serverEnv.apiUrl}/${path.join("/")}${req.nextUrl.search}`;
  let access = req.cookies.get(ACCESS_COOKIE)?.value;
  const refresh = req.cookies.get(REFRESH_COOKIE)?.value;
  let refreshed: Awaited<ReturnType<typeof refreshTokens>> = null;

  if (!access && refresh) {
    refreshed = await refreshTokens(refresh);
    access = refreshed?.access_token;
  }
  if (!access) {
    return NextResponse.json({ error: { code: "unauthorized", message: "not signed in" } }, { status: 401 });
  }

  const body = req.method === "GET" || req.method === "HEAD" ? undefined : await req.arrayBuffer();
  const forward = async (token: string) =>
    fetch(target, {
      method: req.method,
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": req.headers.get("content-type") ?? "application/json",
        Accept: req.headers.get("accept") ?? "application/json",
        "X-Forwarded-For": req.headers.get("x-forwarded-for") ?? "",
        "User-Agent": req.headers.get("user-agent") ?? "habarchy-web",
      },
      body,
      cache: "no-store",
      // @ts-expect-error Node fetch option for streaming request bodies
      duplex: "half",
    });

  let upstream = await forward(access);
  if (upstream.status === 401 && refresh && !refreshed) {
    refreshed = await refreshTokens(refresh);
    if (refreshed) upstream = await forward(refreshed.access_token);
  }

  const headers = new Headers();
  for (const h of ["content-type", "content-disposition", "cache-control", "x-request-id"]) {
    const v = upstream.headers.get(h);
    if (v) headers.set(h, v);
  }
  const res = new NextResponse(upstream.body, { status: upstream.status, headers });
  if (refreshed) setAuthCookies(res, refreshed);
  if (upstream.status === 401 && !refreshed) clearAuthCookies(res);
  return res;
}

export { handle as GET, handle as POST, handle as PUT, handle as PATCH, handle as DELETE };
