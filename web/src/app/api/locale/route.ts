import { NextResponse } from "next/server";
import { isLocale, LOCALE_COOKIE } from "@/i18n/config";

export async function POST(req: Request) {
  const { locale } = (await req.json()) as { locale?: string };
  if (!isLocale(locale)) return NextResponse.json({ error: { code: "validation_failed", message: "bad locale" } }, { status: 400 });
  const res = NextResponse.json({ data: { locale } });
  res.cookies.set(LOCALE_COOKIE, locale, { path: "/", maxAge: 365 * 24 * 3600, sameSite: "lax" });
  return res;
}
