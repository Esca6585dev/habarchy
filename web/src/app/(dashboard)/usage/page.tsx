"use client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { Download } from "lucide-react";
import { api, unwrap, rawFetch } from "@/lib/api/client";
import type { UsageRow } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge } from "@/components/common/status-badge";
import { EmptyState } from "@/components/common/empty-state";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatMoney } from "@/lib/utils";

export default function UsagePage() {
  const t = useTranslations("usage");
  const tc = useTranslations("common");
  const { projectId } = useProject();
  const [days, setDays] = useState(30);
  const [groupBy, setGroupBy] = useState<"" | "day" | "channel">("");
  const q = useQuery({
    queryKey: ["usage", projectId, { days, groupBy }],
    enabled: !!projectId,
    queryFn: async () => unwrap<UsageRow[]>(await api.GET("/api/admin/projects/{project_id}/usage", { params: { path: { project_id: projectId! }, query: { days, group_by: groupBy || undefined } } })),
  });
  const totals = (q.data ?? []).reduce<{ sent: number; delivered: number; failed: number; cost: number }>(
    (a, r) => ({ sent: a.sent + (r.sent ?? 0), delivered: a.delivered + (r.delivered ?? 0), failed: a.failed + (r.failed ?? 0), cost: a.cost + (r.cost_micros ?? 0) }),
    { sent: 0, delivered: 0, failed: 0, cost: 0 },
  );
  const exportCsv = async () => {
    const res = await rawFetch(`/api/admin/projects/${projectId}/usage?days=${days}&format=csv${groupBy ? `&group_by=${groupBy}` : ""}`);
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "habarchy-usage.csv";
    a.click();
    URL.revokeObjectURL(url);
  };
  return (
    <>
      <PageHeader
        title={t("title")}
        description={t("hint")}
        actions={
          <>
            {[7, 30, 90].map((n) => <Button key={n} size="sm" variant={days === n ? "default" : "outline"} onClick={() => setDays(n)}>{tc("days", { n })}</Button>)}
            <Select value={groupBy || "both"} onValueChange={(v) => setGroupBy(v === "both" ? "" : (v as "day" | "channel"))}>
              <SelectTrigger className="w-40"><SelectValue /></SelectTrigger>
              <SelectContent><SelectItem value="both">{t("both")}</SelectItem><SelectItem value="day">{t("day")}</SelectItem><SelectItem value="channel">{t("channel")}</SelectItem></SelectContent>
            </Select>
            <Button size="sm" variant="outline" onClick={exportCsv}><Download /> {tc("export")}</Button>
          </>
        }
      />
      {q.isLoading ? <Skeleton className="h-48" /> : !q.data?.length ? <EmptyState /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow>{groupBy !== "channel" ? <TableHead>{t("day")}</TableHead> : null}{groupBy !== "day" ? <TableHead>{t("channel")}</TableHead> : null}<TableHead className="text-right">{t("sent")}</TableHead><TableHead className="text-right">{t("delivered")}</TableHead><TableHead className="text-right">{t("failed")}</TableHead><TableHead className="text-right">{t("cost")}</TableHead></TableRow></TableHeader>
            <TableBody>
              {q.data.map((r, i) => (
                <TableRow key={i}>
                  {groupBy !== "channel" ? <TableCell className="tabular">{r.day}</TableCell> : null}
                  {groupBy !== "day" ? <TableCell><ChannelBadge channel={r.channel} /></TableCell> : null}
                  <TableCell className="tabular text-right">{r.sent}</TableCell>
                  <TableCell className="tabular text-right text-success">{r.delivered}</TableCell>
                  <TableCell className="tabular text-right text-destructive">{r.failed}</TableCell>
                  <TableCell className="tabular text-right">{formatMoney(r.cost_micros ?? 0, r.currency)}</TableCell>
                </TableRow>
              ))}
              <TableRow className="font-medium">
                <TableCell colSpan={groupBy ? 1 : 2}>Σ</TableCell>
                <TableCell className="tabular text-right">{totals.sent}</TableCell>
                <TableCell className="tabular text-right text-success">{totals.delivered}</TableCell>
                <TableCell className="tabular text-right text-destructive">{totals.failed}</TableCell>
                <TableCell className="tabular text-right">{formatMoney(totals.cost)}</TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      )}
    </>
  );
}
