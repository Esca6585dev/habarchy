"use client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { api, unwrap } from "@/lib/api/client";
import type { InboundSMS, Provider } from "@/lib/api/types";
import { qk } from "@/lib/hooks";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { formatDate } from "@/lib/utils";

export default function InboundPage() {
  const t = useTranslations("inbound");
  const { projectId } = useProject();
  const [offset, setOffset] = useState(0);
  const limit = 50;
  const q = useQuery({
    queryKey: ["inbound", projectId, offset],
    enabled: !!projectId,
    refetchInterval: 15_000,
    queryFn: async () => {
      const res = await api.GET("/api/admin/projects/{project_id}/inbound", { params: { path: { project_id: projectId! }, query: { limit, offset } } });
      return { rows: unwrap<InboundSMS[]>(res), total: Number((res.data as { meta?: { total?: number } })?.meta?.total ?? 0) };
    },
  });
  const providers = useQuery({
    queryKey: qk.providers(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Provider[]>(await api.GET("/api/admin/projects/{project_id}/providers", { params: { path: { project_id: projectId! } } })),
  });
  const nameOf = (id?: string) => providers.data?.find((p) => p.id === id)?.name ?? id?.slice(0, 8) ?? "";
  return (
    <>
      <PageHeader title={t("title")} description={t("hint")} />
      {q.isLoading ? <Skeleton className="h-48" /> : !q.data?.rows.length ? <EmptyState hint={t("empty")} /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow><TableHead>{t("receivedAt")}</TableHead><TableHead>{t("from")}</TableHead><TableHead>{t("text")}</TableHead><TableHead>{t("provider")}</TableHead></TableRow></TableHeader>
            <TableBody>
              {q.data.rows.map((r) => (
                <TableRow key={r.id}>
                  <TableCell className="tabular text-xs text-muted-foreground whitespace-nowrap">{formatDate(r.received_at, { dateStyle: "short", timeStyle: "short" })}</TableCell>
                  <TableCell className="font-mono text-xs">{r.from}</TableCell>
                  <TableCell className="max-w-xl whitespace-pre-wrap text-sm">{r.text}</TableCell>
                  <TableCell className="text-xs text-muted-foreground">{nameOf(r.provider_id)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <div className="flex items-center justify-between border-t p-3 text-xs text-muted-foreground">
            <span>{offset + 1}–{offset + q.data.rows.length} / {q.data.total}</span>
            <div className="space-x-2">
              <Button size="sm" variant="outline" disabled={offset === 0} onClick={() => setOffset(Math.max(0, offset - limit))}>‹</Button>
              <Button size="sm" variant="outline" disabled={offset + limit >= q.data.total} onClick={() => setOffset(offset + limit)}>›</Button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
