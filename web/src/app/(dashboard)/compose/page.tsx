"use client";
import { Suspense, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Send } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError, errorOf } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Group, Template } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge } from "@/components/common/status-badge";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";

const channels = ["sms", "whatsapp", "telegram", "email", "push", "slack"] as const;
type Channel = (typeof channels)[number];
type SendResult = { accepted: number; rejected: { to: string; reason: string }[]; batchId: string };

export default function ComposePage() {
  return <Suspense><Compose /></Suspense>;
}

function Compose() {
  const t = useTranslations("compose");
  const params = useSearchParams();
  const { projectId } = useProject();
  const [channel, setChannel] = useState<Channel>("sms");
  const [groupIds, setGroupIds] = useState<Set<string>>(() => new Set(params.get("group") ? [params.get("group")!] : []));
  const [addresses, setAddresses] = useState("");
  const [template, setTemplate] = useState<string>("");
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");
  const [data, setData] = useState("{}");
  const [sandbox, setSandbox] = useState(false);
  const [result, setResult] = useState<SendResult | null>(null);

  const groups = useQuery({
    queryKey: qk.groups(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Group[]>(await api.GET("/api/admin/projects/{project_id}/groups", { params: { path: { project_id: projectId! } } })),
  });
  const templates = useQuery({
    queryKey: qk.templates(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Template[]>(await api.GET("/api/admin/projects/{project_id}/templates", { params: { path: { project_id: projectId! } } })),
  });
  const templateKeys = useMemo(() => Array.from(new Set((templates.data ?? []).filter((x) => x.channel === channel).map((x) => x.key!))), [templates.data, channel]);
  const parsedData = useMemo(() => { try { return JSON.parse(data || "{}") as Record<string, unknown>; } catch { return null; } }, [data]);
  const toList = addresses.split("\n").map((s) => s.trim()).filter(Boolean);
  const canSend = (groupIds.size > 0 || toList.length > 0) && !!(template || body.trim()) && parsedData !== null;

  const send = useMutation({
    mutationFn: async (): Promise<SendResult> => {
      const res = await api.POST("/api/admin/projects/{project_id}/messages/send", {
        params: { path: { project_id: projectId! } },
        body: {
          channel, template: template || undefined, body: template ? undefined : body, subject: subject || undefined, title: channel === "push" ? subject || undefined : undefined,
          data: parsedData ?? {}, group_ids: Array.from(groupIds), to: toList, is_test: sandbox,
        },
      });
      if (res.error || !res.response.ok) {
        throw new ApiError(res.response.status, errorOf(res) ?? undefined);
      }
      const env = res.data as unknown as { data: { id: string }; meta: { accepted: number; rejected: { to: string; reason: string }[] } };
      return { accepted: env.meta.accepted, rejected: env.meta.rejected ?? [], batchId: env.data.id };
    },
    onSuccess: (r) => { setResult(r); toast.success(t("sent", { accepted: r.accepted, rejected: r.rejected.length })); },
    onError: (e: Error) => toast.error(e instanceof ApiError ? `${e.message} ${e.details ? JSON.stringify(e.details) : ""}` : e.message),
  });

  return (
    <>
      <PageHeader title={t("title")} />
      <div className="grid gap-6 lg:grid-cols-5">
        <Card className="lg:col-span-3">
          <CardHeader><CardTitle>{t("channel")}</CardTitle></CardHeader>
          <CardContent className="space-y-5">
            <div className="flex flex-wrap gap-2">
              {channels.map((c) => (
                <button key={c} type="button" onClick={() => { setChannel(c); setTemplate(""); }} className={`rounded-full border px-3 py-1 text-sm ${channel === c ? "border-primary bg-primary/10" : "hover:bg-accent"}`} data-testid={`channel-${c}`}>
                  <ChannelBadge channel={c} />
                </button>
              ))}
            </div>
            {channel === "whatsapp" ? <Alert><AlertDescription>{t("whatsappHint")}</AlertDescription></Alert> : null}
            <Field label={t("template")}>
              <Select value={template || "__free"} onValueChange={(v) => setTemplate(v === "__free" ? "" : v)}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="__free">{t("freeText")}</SelectItem>
                  {templateKeys.map((k) => <SelectItem key={k} value={k}>{k}</SelectItem>)}
                </SelectContent>
              </Select>
            </Field>
            {!template ? (
              <>
                {channel !== "sms" ? <Field label={t("subject")} htmlFor="c-subject"><Input id="c-subject" value={subject} onChange={(e) => setSubject(e.target.value)} /></Field> : null}
                <Field label={t("body")} htmlFor="c-body"><Textarea id="c-body" rows={6} value={body} onChange={(e) => setBody(e.target.value)} placeholder="Salam {{.name}}! …" data-testid="compose-body" /></Field>
              </>
            ) : null}
            <Field label={t("data")} htmlFor="c-data" error={parsedData ? undefined : "invalid JSON"}>
              <Textarea id="c-data" rows={3} className="font-mono text-xs" value={data} onChange={(e) => setData(e.target.value)} />
            </Field>
            <div className="flex items-center gap-2"><Switch id="c-sandbox" checked={sandbox} onCheckedChange={setSandbox} /><Label htmlFor="c-sandbox">{t("sandbox")}</Label></div>
          </CardContent>
        </Card>
        <div className="space-y-6 lg:col-span-2">
          <Card>
            <CardHeader><CardTitle>{t("targets")}</CardTitle></CardHeader>
            <CardContent className="space-y-4">
              <div>
                <p className="mb-2 text-sm font-medium">{t("groups")}</p>
                <div className="flex flex-wrap gap-2">
                  {(groups.data ?? []).map((g) => {
                    const on = groupIds.has(g.id!);
                    return (
                      <button key={g.id} type="button" onClick={() => { const n = new Set(groupIds); if (on) n.delete(g.id!); else n.add(g.id!); setGroupIds(n); }} className={`rounded-md border px-2.5 py-1 text-sm ${on ? "border-primary bg-primary/10" : "hover:bg-accent"}`} data-testid={`group-${g.id}`}>
                        {g.name} <Badge variant="secondary" className="ml-1">{g.member_count}</Badge>
                      </button>
                    );
                  })}
                  {!groups.data?.length ? <Link href="/groups" className="text-sm text-muted-foreground underline">{t("groups")} →</Link> : null}
                </div>
              </div>
              <Field label={t("addresses")} htmlFor="c-to" hint={t("addressesHint")}>
                <Textarea id="c-to" rows={4} className="font-mono text-xs" value={addresses} onChange={(e) => setAddresses(e.target.value)} placeholder={channel === "email" ? "user@example.tm" : channel === "slack" ? "C0123456789" : "+99365123456"} />
              </Field>
              <Button className="w-full" onClick={() => send.mutate()} disabled={!canSend || send.isPending} data-testid="compose-send"><Send /> {t("send")}</Button>
              {groupIds.size === 0 && toList.length === 0 ? <p className="text-xs text-muted-foreground">{t("pickTarget")}</p> : null}
            </CardContent>
          </Card>
          {result ? (
            <Card>
              <CardContent className="space-y-2 pt-6 text-sm">
                <p className="font-medium">{t("sent", { accepted: result.accepted, rejected: result.rejected.length })}</p>
                {result.rejected.slice(0, 10).map((r, i) => <p key={i} className="text-xs text-destructive">{r.to}: {r.reason}</p>)}
                <Button variant="outline" size="sm" asChild><Link href={`/messages?batch_id=${result.batchId}`}>{t("viewBatch")}</Link></Button>
              </CardContent>
            </Card>
          ) : null}
        </div>
      </div>
    </>
  );
}
