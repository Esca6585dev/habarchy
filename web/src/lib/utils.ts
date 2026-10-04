import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(value?: string | null, opts?: Intl.DateTimeFormatOptions) {
  if (!value) return "—";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleString(undefined, opts ?? { dateStyle: "medium", timeStyle: "short" });
}

export function formatRelative(value?: string | null) {
  if (!value) return "—";
  const diff = Date.now() - new Date(value).getTime();
  const s = Math.round(diff / 1000);
  if (s < 60) return `${s}s`;
  const m = Math.round(s / 60);
  if (m < 60) return `${m}m`;
  const h = Math.round(m / 60);
  if (h < 48) return `${h}h`;
  return `${Math.round(h / 24)}d`;
}

export function formatMoney(micros: number, currency = "TMT") {
  return `${(micros / 1_000_000).toFixed(2)} ${currency}`;
}

export function truncate(s: string, n = 60) {
  return s.length > n ? s.slice(0, n) + "…" : s;
}
