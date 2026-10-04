"use client";
import { Suspense, useState } from "react";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { RotateCcw } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import type { WebhookDelivery } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { StatusBadge } from "@/components/common/status-badge";
import { EmptyState } from "@/components/common/empty-state";
import { JsonView } from "@/components/common/json-view";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate, truncate } from "@/lib/utils";

function WebhooksInner() {
  const t = useTranslations("webhooks");
  const tc = useTranslations("common");
  const { projectId, can } = useProject();
  const params = useSearchParams();
  const qc = useQueryClient();
  const [onlyFailed, setOnlyFailed] = useState(false);
  const [detail, setDetail] = useState<WebhookDelivery | null>(null);
  const q = useQuery({
    queryKey: ["webhooks", projectId, { onlyFailed }],
    enabled: !!projectId,
    refetchInterval: 15_000,
    queryFn: async () => unwrap<WebhookDelivery[]>(await api.GET("/api/admin/projects/{project_id}/webhooks", { params: { path: { project_id: projectId! }, query: { failed: onlyFailed || undefined, limit: 100 } } })),
  });
  const resend = useMutation({
    mutationFn: async (id: string) => unwrap<WebhookDelivery>(await api.POST("/api/admin/projects/{project_id}/webhooks/{delivery_id}/resend", { params: { path: { project_id: projectId!, delivery_id: id } } })),
    onSuccess: () => { toast.success(t("resend")); void qc.invalidateQueries({ queryKey: ["webhooks", projectId] }); },
    onError: (e: ApiError) => toast.error(e.message),
  });
  const filterMsg = params.get("message_id");
  const rows = (q.data ?? []).filter((w) => !filterMsg || w.message_id === filterMsg);
  return (
    <>
      <PageHeader title={t("title")} actions={<div className="flex items-center gap-2"><Switch id="only-failed" checked={onlyFailed} onCheckedChange={setOnlyFailed} /><Label htmlFor="only-failed" className="text-xs">{t("onlyFailed")}</Label></div>} />
      {q.isLoading ? <Skeleton className="h-48" /> : rows.length === 0 ? <EmptyState /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow><TableHead>{tc("status")}</TableHead><TableHead>{t("event")}</TableHead><TableHead>{t("url")}</TableHead><TableHead>{t("response")}</TableHead><TableHead>{t("attempts")}</TableHead><TableHead>{t("nextRetry")}</TableHead><TableHead>{tc("created")}</TableHead><TableHead /></TableRow></TableHeader>
            <TableBody>
              {rows.map((w) => (
                <TableRow key={w.id} className="cursor-pointer" onClick={() => setDetail(w)}>
                  <TableCell><StatusBadge status={w.status} /></TableCell>
                  <TableCell className="font-medium">{w.event}</TableCell>
                  <TableCell className="max-w-[16rem] truncate font-mono text-xs">{truncate(w.url ?? "", 50)}</TableCell>
                  <TableCell className="tabular text-xs">{w.response_code ?? "—"}</TableCell>
                  <TableCell className="tabular">{w.attempts}</TableCell>
                  <TableCell className="tabular text-xs text-muted-foreground">{formatDate(w.next_retry_at)}</TableCell>
                  <TableCell className="tabular text-xs text-muted-foreground">{formatDate(w.created_at)}</TableCell>
                  <TableCell className="text-right">{can("developer") && w.status !== "delivered" ? <Button size="sm" variant="ghost" onClick={(e) => { e.stopPropagation(); resend.mutate(w.id!); }}><RotateCcw /> {t("resend")}</Button> : null}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
      <Dialog open={!!detail} onOpenChange={(o) => !o && setDetail(null)}>
        <DialogContent wide>
          <DialogHeader><DialogTitle>{detail?.event} · <StatusBadge status={detail?.status} /></DialogTitle></DialogHeader>
          <p className="break-all font-mono text-xs">{detail?.url}</p>
          <p className="text-xs font-medium uppercase text-muted-foreground">{t("payload")}</p>
          <JsonView value={detail?.payload} />
          <p className="text-xs font-medium uppercase text-muted-foreground">{t("response")} {detail?.response_code ?? ""}</p>
          <JsonView value={detail?.response_body || "—"} className="max-h-40" />
        </DialogContent>
      </Dialog>
    </>
  );
}

export default function WebhooksPage() {
  return <Suspense><WebhooksInner /></Suspense>;
}
