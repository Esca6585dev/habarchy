"use client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { api, unwrap } from "@/lib/api/client";
import type { AuditLog } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { JsonView } from "@/components/common/json-view";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

export default function AuditPage() {
  const t = useTranslations("audit");
  const [offset, setOffset] = useState(0);
  const [detail, setDetail] = useState<AuditLog | null>(null);
  const { projectId } = useProject();
  const limit = 50;
  const q = useQuery({
    queryKey: ["audit", projectId, offset],
    enabled: !!projectId,
    queryFn: async () => {
      const res = await api.GET("/api/admin/projects/{project_id}/audit-logs", { params: { path: { project_id: projectId! }, query: { limit, offset } } });
      return { rows: unwrap<AuditLog[]>(res), total: Number((res.data as { meta?: { total?: number } })?.meta?.total ?? 0) };
    },
  });
  return (
    <>
      <PageHeader title={t("title")} />
      {q.isLoading ? <Skeleton className="h-48" /> : !q.data?.rows.length ? <EmptyState /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow><TableHead>{t("when")}</TableHead><TableHead>{t("action")}</TableHead><TableHead>{t("entity")}</TableHead><TableHead>{t("user")}</TableHead><TableHead>{t("ip")}</TableHead></TableRow></TableHeader>
            <TableBody>
              {q.data.rows.map((a) => (
                <TableRow key={a.ID} className="cursor-pointer" onClick={() => setDetail(a)}>
                  <TableCell className="tabular whitespace-nowrap text-xs text-muted-foreground">{formatDate(a.CreatedAt, { dateStyle: "short", timeStyle: "medium" })}</TableCell>
                  <TableCell><Badge variant="secondary" className="font-mono">{a.Action}</Badge></TableCell>
                  <TableCell className="font-mono text-xs">{a.EntityType} {a.EntityID?.slice(0, 8)}</TableCell>
                  <TableCell className="font-mono text-xs">{a.UserID?.slice(0, 8) ?? "—"}</TableCell>
                  <TableCell className="font-mono text-xs">{a.Ip}</TableCell>
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
      <Dialog open={!!detail} onOpenChange={(o) => !o && setDetail(null)}>
        <DialogContent>
          <DialogHeader><DialogTitle>{detail?.Action}</DialogTitle></DialogHeader>
          <p className="text-xs text-muted-foreground">{detail?.UserAgent}</p>
          <JsonView value={detail?.Changes} />
        </DialogContent>
      </Dialog>
    </>
  );
}
