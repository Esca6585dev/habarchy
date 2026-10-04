"use client";
import createClient, { type Middleware } from "openapi-fetch";
import type { paths } from "./schema";
import type { APIError } from "./types";

/**
 * Browser-side API client. Requests go to the Next.js proxy
 * (/api/backend/...) which attaches the admin JWT from the httpOnly
 * cookie and refreshes it when expired. The browser never sees tokens.
 */
export const api = createClient<paths>({ baseUrl: "/api/backend" });

export class ApiError extends Error {
  code: string;
  status: number;
  details?: Record<string, unknown>;
  constructor(status: number, err?: APIError) {
    super(err?.message ?? `HTTP ${status}`);
    this.code = err?.code ?? "internal_error";
    this.status = status;
    this.details = err?.details;
  }
}

const unauthorized: Middleware = {
  async onResponse({ response }) {
    if (response.status === 401 && typeof window !== "undefined" && !window.location.pathname.startsWith("/login")) {
      window.location.assign(`/login?next=${encodeURIComponent(window.location.pathname)}`);
    }
    return response;
  },
};
api.use(unauthorized);

/** Unwrap an openapi-fetch result into data or throw ApiError. */
export function unwrap<T>(res: { data?: unknown; error?: unknown; response: Response }): T {
  if (res.error || !res.response.ok) {
    const err = (res.error as { error?: APIError } | undefined)?.error;
    throw new ApiError(res.response.status, err);
  }
  return (res.data as { data: T }).data;
}

/** Extract the API error from an openapi-fetch result, if any. */
export function errorOf(res: { error?: unknown; response: Response }): APIError | undefined {
  const body = res.error as unknown as { error?: APIError } | undefined;
  return body?.error;
}

/** Throw ApiError when the response is not 2xx (for 204 endpoints). */
export function ensureOk(res: { error?: unknown; response: Response }): void {
  if (res.error || !res.response.ok) throw new ApiError(res.response.status, errorOf(res));
}

/** Generic fetch through the proxy for endpoints without typed helpers (CSV, raw). */
export async function rawFetch(path: string, init?: RequestInit) {
  const res = await fetch(`/api/backend${path}`, init);
  if (res.status === 401 && typeof window !== "undefined") {
    window.location.assign("/login");
  }
  return res;
}
