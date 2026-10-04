"use client";
import { useState } from "react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { api, unwrap } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Dashboard } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { StatCard } from "@/components/dashboard/stat-card";
import { DailyChart } from "@/components/dashboard/daily-chart";
import { LiveFeed } from "@/components/dashboard/live-feed";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { EmptyState } from "@/components/common/empty-state";
import { ChannelBadge, StatusBadge } from "@/components/common/status-badge";
import { formatMoney, formatRelative } from "@/lib/utils";

export default function DashboardPage() {
  const t = useTranslations("dashboard");
  const tc = useTranslations("common");
  const { projectId, project, projects, isLoading } = useProject();
  const [days, setDays] = useState(14);
  const q = useQuery({
    queryKey: qk.dashboard(projectId ?? "", days),
    enabled: !!projectId,
    refetchInterval: 30_000,
    queryFn: async () =>
      unwrap<Dashboard>(await api.GET("/api/admin/projects/{project_id}/dashboard", { params: { path: { project_id: projectId! }, query: { days } } })),
  });

  if (!isLoading && projects.length === 0) {
    return (
      <>
        <PageHeader title={t("title")} />
        <EmptyState title={t("noProject")} action={<Button asChild><Link href="/projects">{tc("create")}</Link></Button>} />
      </>
    );
  }
  const d = q.data;
  const fmtSec = (s?: number) => (s === undefined ? "—" : s < 1 ? `${Math.round(s * 1000)} ms` : `${s.toFixed(1)} s`);

  return (
    <>
      <PageHeader
        title={project?.name ?? t("title")}
        description={project?.slug}
        actions={[7, 14, 30].map((n) => (
          <Button key={n} size="sm" variant={days === n ? "default" : "outline"} onClick={() => setDays(n)}>
            {tc("days", { n })}
          </Button>
        ))}
      />
      {!d ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">{Array.from({ length: 4 }).map((_, i) => <Skeleton key={i} className="h-24" />)}</div>
      ) : (
        <>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <StatCard label={t("total")} value={d.totals?.total ?? 0} sub={`${t("pending")}: ${d.totals?.pending ?? 0}`} />
            <StatCard label={t("delivered")} value={d.totals?.delivered ?? 0} sub={`${t("sent")}: ${d.totals?.sent ?? 0}`} tone="success" />
            <StatCard label={t("failed")} value={d.totals?.failed ?? 0} tone={(d.totals?.failed ?? 0) > 0 ? "destructive" : undefined} />
            <StatCard label={t("cost")} value={formatMoney(d.totals?.cost_micros ?? 0)} sub={`${t("latencyP95")}: ${fmtSec(d.latency?.p95_sent_sec)} · ${t("samples", { n: d.latency?.samples ?? 0 })}`} />
          </div>
          <div className="mt-6 grid gap-6 lg:grid-cols-3">
            <Card className="lg:col-span-2">
              <CardHeader><CardTitle>{t("perDay")}</CardTitle></CardHeader>
              <CardContent><DailyChart data={d} /></CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle>{t("byChannel")}</CardTitle></CardHeader>
              <CardContent className="space-y-3">
                {Object.entries(d.by_channel ?? {}).length === 0 ? <p className="text-sm text-muted-foreground">—</p> : null}
                {Object.entries(d.by_channel ?? {}).map(([ch, tot]) => (
                  <div key={ch} className="flex items-center justify-between text-sm">
                    <ChannelBadge channel={ch} />
                    <span className="tabular text-muted-foreground">
                      <span className="text-success">{tot.delivered ?? 0}</span> / <span className="text-info">{tot.sent ?? 0}</span> / <span className="text-destructive">{tot.failed ?? 0}</span>
                    </span>
                  </div>
                ))}
                <div className="border-t pt-3 text-xs text-muted-foreground">
                  {t("latencyP50")}: {fmtSec(d.latency?.p50_sent_sec)} · {t("latencyDelivered")}: {fmtSec(d.latency?.p95_delivered_sec)}
                </div>
              </CardContent>
            </Card>
          </div>
          <div className="mt-6 grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader><CardTitle>{t("live")}</CardTitle></CardHeader>
              <CardContent>{projectId ? <LiveFeed projectId={projectId} /> : null}</CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle>{t("recentFailures")}</CardTitle></CardHeader>
              <CardContent>
                {(d.recent_failures ?? []).length === 0 ? (
                  <p className="text-sm text-muted-foreground">—</p>
                ) : (
                  <ul className="divide-y">
                    {(d.recent_failures as Array<Record<string, string>>).map((m) => (
                      <li key={m.ID} className="flex items-center gap-3 py-2 text-sm">
                        <StatusBadge status="failed" />
                        <ChannelBadge channel={m.Channel} />
                        <Link href={`/messages/${m.ID}`} className="min-w-0 flex-1 truncate font-mono text-xs hover:underline">{m.ToAddress}</Link>
                        <span className="truncate text-xs text-destructive">{m.ErrorCode}</span>
                        <span className="tabular shrink-0 text-xs text-muted-foreground">{formatRelative(m.CreatedAt)}</span>
                      </li>
                    ))}
                  </ul>
                )}
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </>
  );
}
