import { NextRequest, NextResponse } from "next/server";
import { ACCESS_COOKIE, REFRESH_COOKIE } from "@/i18n/config";

const PUBLIC = ["/login", "/api/auth", "/api/locale", "/_next", "/favicon.ico", "/icon", "/manifest"];

/**
 * Guards the dashboard: without an access or refresh cookie the user goes
 * to /login. Token refresh itself happens in the API proxy so edge
 * middleware stays cheap.
 */
export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;
  if (PUBLIC.some((p) => pathname.startsWith(p))) return NextResponse.next();
  const hasSession = req.cookies.has(ACCESS_COOKIE) || req.cookies.has(REFRESH_COOKIE);
  if (!hasSession) {
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    url.searchParams.set("next", pathname);
    return NextResponse.redirect(url);
  }
  return NextResponse.next();
}

export const config = { matcher: ["/((?!_next/static|_next/image|.*\\.(?:png|svg|ico|webp)).*)"] };
