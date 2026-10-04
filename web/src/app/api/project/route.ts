import { NextResponse } from "next/server";
import { PROJECT_COOKIE } from "@/i18n/config";

/** Remembers the selected project in a cookie so server components know it. */
export async function POST(req: Request) {
  const { project_id } = (await req.json()) as { project_id?: string };
  const res = NextResponse.json({ data: { project_id } });
  if (project_id) res.cookies.set(PROJECT_COOKIE, project_id, { path: "/", maxAge: 365 * 24 * 3600, sameSite: "lax" });
  else res.cookies.set(PROJECT_COOKIE, "", { path: "/", maxAge: 0 });
  return res;
}
