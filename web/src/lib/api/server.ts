import "server-only";
import { cookies } from "next/headers";
import { ACCESS_COOKIE } from "@/i18n/config";
import { serverEnv } from "@/lib/env";
import type { Envelope } from "./types";

/** Server-component fetch against the backend with the cookie token. */
export async function serverApi<T>(path: string, init?: RequestInit): Promise<Envelope<T> | null> {
  const store = await cookies();
  const token = store.get(ACCESS_COOKIE)?.value;
  if (!token) return null;
  const res = await fetch(`${serverEnv.apiUrl}${path}`, {
    ...init,
    headers: { ...(init?.headers ?? {}), Authorization: `Bearer ${token}`, Accept: "application/json" },
    cache: "no-store",
  });
  if (res.status === 401) return null;
  return (await res.json()) as Envelope<T>;
}
