"use client";
import { Bar, BarChart, CartesianGrid, Legend, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { useTranslations } from "next-intl";
import type { Dashboard } from "@/lib/api/types";

const colors: Record<string, string> = {
  sms: "var(--chart-sms)",
  email: "var(--chart-email)",
  push: "var(--chart-push)",
  telegram: "var(--chart-telegram)",
};

/** Stacked bars per day, one series per channel (sent + delivered). */
export function DailyChart({ data }: { data: Dashboard }) {
  const t = useTranslations("channel");
  type Row = Record<string, string | number>;
  const byDay = new Map<string, Row>();
  const channels = new Set<string>();
  for (const p of data.daily ?? []) {
    if (!p.day || !p.channel) continue;
    channels.add(p.channel);
    const row: Row = byDay.get(p.day) ?? { day: p.day };
    row[p.channel] = Number(row[p.channel] ?? 0) + (p.total ?? 0);
    byDay.set(p.day, row);
  }
  const rows = Array.from(byDay.values()).sort((a, b) => String(a.day).localeCompare(String(b.day)));
  if (!rows.length) return <p className="py-10 text-center text-sm text-muted-foreground">—</p>;
  return (
    <ResponsiveContainer width="100%" height={260}>
      <BarChart data={rows} margin={{ top: 8, right: 8, left: -16, bottom: 0 }}>
        <CartesianGrid vertical={false} stroke="var(--border)" />
        <XAxis dataKey="day" tickFormatter={(d: string) => d.slice(5)} tick={{ fontSize: 12, fill: "var(--muted-foreground)" }} axisLine={false} tickLine={false} />
        <YAxis allowDecimals={false} tick={{ fontSize: 12, fill: "var(--muted-foreground)" }} axisLine={false} tickLine={false} />
        <Tooltip
          cursor={{ fill: "var(--accent)", opacity: 0.4 }}
          contentStyle={{ background: "var(--popover)", border: "1px solid var(--border)", borderRadius: 8, fontSize: 12, color: "var(--popover-foreground)" }}
        />
        <Legend wrapperStyle={{ fontSize: 12 }} formatter={(v: string) => t(v as "sms")} />
        {Array.from(channels).map((c) => (
          <Bar key={c} dataKey={c} stackId="a" fill={colors[c] ?? "var(--primary)"} radius={[2, 2, 0, 0]} />
        ))}
      </BarChart>
    </ResponsiveContainer>
  );
}
