"use client";
import { use } from "react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Ban, RotateCcw } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import type { Message, MessageEvent, WebhookDelivery } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge, StatusBadge } from "@/components/common/status-badge";
import { JsonView } from "@/components/common/json-view";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

interface Detail { message: Message; events: MessageEvent[]; webhooks: WebhookDelivery[] }

export default function MessageDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const t = useTranslations("messages");
  const tc = useTranslations("common");
  const { projectId, can } = useProject();
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["message", projectId, id],
    enabled: !!projectId,
    refetchInterval: (query) => (["queued", "processing", "sent"].includes(query.state.data?.message.status ?? "") ? 5000 : false),
    queryFn: async () => unwrap<Detail>(await api.GET("/api/admin/projects/{project_id}/messages/{message_id}", { params: { path: { project_id: projectId!, message_id: id } } })),
  });
  const act = useMutation({
    mutationFn: async (kind: "resend" | "cancel") => {
      const path = kind === "resend" ? "/api/admin/projects/{project_id}/messages/{message_id}/resend" : "/api/admin/projects/{project_id}/messages/{message_id}/cancel";
      return unwrap<Message>(await api.POST(path, { params: { path: { project_id: projectId!, message_id: id } } }));
    },
    onSuccess: (_, kind) => {
      toast.success(kind === "resend" ? t("resent") : tc("save"));
      void qc.invalidateQueries({ queryKey: ["message", projectId, id] });
    },
    onError: (e: ApiError) => toast.error(e.message),
  });

  if (q.isLoading || !q.data) return <Skeleton className="h-64" />;
  const { message: m, events, webhooks } = q.data;
  const fields: Array<[string, React.ReactNode]> = [
    [t("to"), <span key="to" className="font-mono">{m.to}</span>],
    [t("template"), m.template ? `${m.template} v${m.template_version ?? ""}` : "—"],
    [t("provider"), m.provider_message_id ? <span key="p" className="font-mono text-xs">{m.provider_message_id}</span> : "—"],
    [t("attempts"), m.attempts],
    ["Priority", m.priority],
    [tc("created"), formatDate(m.created_at)],
    ["Sent", formatDate(m.sent_at)],
    ["Delivered", formatDate(m.delivered_at)],
    ["Scheduled", formatDate(m.scheduled_at)],
    ["Cost", `${((m.cost_micros ?? 0) / 1e6).toFixed(4)} ${m.currency ?? ""}`],
  ];
  return (
    <>
      <PageHeader
        title={t("detail")}
        description={m.id}
        actions={
          <>
            <Button variant="outline" size="sm" asChild><Link href="/messages"><ArrowLeft /> {t("title")}</Link></Button>
            {can("developer") && ["failed", "cancelled"].includes(m.status ?? "") ? (
              <Button size="sm" onClick={() => act.mutate("resend")} disabled={act.isPending}><RotateCcw /> {t("resend")}</Button>
            ) : null}
            {can("developer") && m.status === "queued" ? (
              <Button size="sm" variant="destructive" onClick={() => act.mutate("cancel")} disabled={act.isPending}><Ban /> {t("cancel")}</Button>
            ) : null}
          </>
        }
      />
      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <CardHeader><CardTitle className="flex items-center gap-2"><StatusBadge status={m.status} /><ChannelBadge channel={m.channel} />{m.is_test ? <span className="text-xs uppercase text-muted-foreground">test</span> : null}</CardTitle></CardHeader>
          <CardContent>
            <dl className="space-y-2 text-sm">
              {fields.map(([k, v]) => (
                <div key={k} className="flex justify-between gap-4"><dt className="text-muted-foreground">{k}</dt><dd className="truncate text-right">{v}</dd></div>
              ))}
            </dl>
            {m.error_code ? <p className="mt-3 rounded-md bg-destructive/10 p-2 text-xs text-destructive"><b>{m.error_code}</b> {m.error_message}</p> : null}
            <div className="mt-4">
              <p className="mb-1 text-xs font-medium uppercase text-muted-foreground">{t("body")}</p>
              {m.subject ? <p className="text-sm font-medium">{m.subject}</p> : null}
              <pre className="whitespace-pre-wrap rounded-md border bg-muted/40 p-2 text-sm">{m.body}</pre>
            </div>
            <div className="mt-4">
              <p className="mb-1 text-xs font-medium uppercase text-muted-foreground">{t("metadata")}</p>
              <JsonView value={m.metadata} />
            </div>
          </CardContent>
        </Card>
        <div className="space-y-6 lg:col-span-2">
          <Card>
            <CardHeader><CardTitle>{t("timeline")}</CardTitle></CardHeader>
            <CardContent>
              <ol className="relative space-y-4 border-l pl-4">
                {events.map((e) => (
                  <li key={e.id} className="relative">
                    <span className="absolute -left-[21px] top-1.5 h-2.5 w-2.5 rounded-full bg-primary" />
                    <div className="flex flex-wrap items-center gap-2 text-sm">
                      <span className="font-medium">{e.type}</span>
                      <span className="tabular text-xs text-muted-foreground">{formatDate(e.created_at, { dateStyle: "short", timeStyle: "medium" })}</span>
                    </div>
                    {e.payload && Object.keys(e.payload).length ? <JsonView value={e.payload} className="mt-1 max-h-48" /> : null}
                  </li>
                ))}
              </ol>
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle>{t("webhooks")}</CardTitle></CardHeader>
            <CardContent>
              {webhooks.length === 0 ? <p className="text-sm text-muted-foreground">—</p> : (
                <ul className="divide-y text-sm">
                  {webhooks.map((w) => (
                    <li key={w.id} className="flex flex-wrap items-center gap-2 py-2">
                      <StatusBadge status={w.status} />
                      <span className="font-medium">{w.event}</span>
                      <span className="tabular text-xs text-muted-foreground">HTTP {w.response_code ?? "—"} · {w.attempts}×</span>
                      <Link href={`/webhooks?message_id=${m.id}`} className="ml-auto text-xs hover:underline">{tc("edit")}</Link>
                    </li>
                  ))}
                </ul>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
