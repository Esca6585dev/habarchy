"use client";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { api, unwrap } from "@/lib/api/client";
import type { Health } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge, StatusBadge } from "@/components/common/status-badge";
import { StatCard } from "@/components/dashboard/stat-card";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

export default function HealthPage() {
  const t = useTranslations("health");
  const { projectId } = useProject();
  const q = useQuery({
    queryKey: ["health", projectId],
    enabled: !!projectId,
    refetchInterval: 10_000,
    queryFn: async () => unwrap<Health>(await api.GET("/api/admin/projects/{project_id}/health", { params: { path: { project_id: projectId! } } })),
  });
  if (q.isLoading || !q.data) return <><PageHeader title={t("title")} /><Skeleton className="h-64" /></>;
  const h = q.data;
  const pending = (h.queues ?? []).reduce((a, x) => a + (x.pending ?? 0), 0);
  const failedToday = (h.queues ?? []).reduce((a, x) => a + (x.failed_today ?? 0), 0);
  return (
    <>
      <PageHeader title={t("title")} />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard label={t("pending")} value={pending} tone={pending > 100 ? "warning" : undefined} />
        <StatCard label={t("failedToday")} value={failedToday} tone={failedToday > 0 ? "destructive" : undefined} />
        <StatCard label={t("devices")} value={`${h.devices?.active ?? 0} / ${h.devices?.total ?? 0}`} />
        <StatCard label={t("contacts")} value={h.contacts ?? 0} />
      </div>
      <div className="mt-6 grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader><CardTitle>{t("queues")}</CardTitle></CardHeader>
          <CardContent className="p-0">
            <Table>
              <TableHeader><TableRow><TableHead>{t("queue")}</TableHead><TableHead className="text-right">{t("pending")}</TableHead><TableHead className="text-right">{t("active")}</TableHead><TableHead className="text-right">{t("scheduled")}</TableHead><TableHead className="text-right">{t("retry")}</TableHead><TableHead className="text-right">{t("processed")}</TableHead><TableHead className="text-right">{t("failedToday")}</TableHead></TableRow></TableHeader>
              <TableBody>
                {(h.queues ?? []).map((qq) => (
                  <TableRow key={qq.queue}>
                    <TableCell className="font-mono text-xs">{qq.queue}</TableCell>
                    <TableCell className="tabular text-right">{qq.pending}</TableCell>
                    <TableCell className="tabular text-right">{qq.active}</TableCell>
                    <TableCell className="tabular text-right">{qq.scheduled}</TableCell>
                    <TableCell className="tabular text-right">{qq.retry}</TableCell>
                    <TableCell className="tabular text-right">{qq.processed_today}</TableCell>
                    <TableCell className="tabular text-right text-destructive">{qq.failed_today}</TableCell>
                  </TableRow>
                ))}
                {!h.queues?.length ? <TableRow><TableCell colSpan={7} className="text-center text-muted-foreground">—</TableCell></TableRow> : null}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle>{t("providers")}</CardTitle></CardHeader>
          <CardContent className="p-0">
            <Table>
              <TableHeader><TableRow><TableHead>Provider</TableHead><TableHead>{t("ok")}</TableHead><TableHead>{t("failed")}</TableHead><TableHead>{t("lastSent")}</TableHead><TableHead /></TableRow></TableHeader>
              <TableBody>
                {(h.providers ?? []).map((p) => (
                  <TableRow key={p.id}>
                    <TableCell><div className="font-medium">{p.name}</div><div className="mt-0.5 flex gap-1 text-xs"><ChannelBadge channel={p.channel} /><span className="font-mono text-muted-foreground">{p.type}</span></div></TableCell>
                    <TableCell className="tabular text-success">{p.ok_count}</TableCell>
                    <TableCell className="tabular text-destructive">{p.failed_count}</TableCell>
                    <TableCell className="tabular text-xs text-muted-foreground">{formatDate(p.last_sent_at)}</TableCell>
                    <TableCell><StatusBadge status={p.status} /></TableCell>
                  </TableRow>
                ))}
                {!h.providers?.length ? <TableRow><TableCell colSpan={5} className="text-center text-muted-foreground">—</TableCell></TableRow> : null}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </>
  );
}
