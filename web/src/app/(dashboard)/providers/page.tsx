"use client";
import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { FlaskConical, Pencil, Plus, Share2, Smartphone, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { GatewayPairing, Provider } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge } from "@/components/common/status-badge";
import { EmptyState } from "@/components/common/empty-state";
import { ConfirmButton } from "@/components/common/confirm-button";
import { Field } from "@/components/common/field";
import { JsonView } from "@/components/common/json-view";
import { CopyButton } from "@/components/common/copy-button";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";

const types = ["http_sms", "smpp", "android_sms", "smtp", "fcm", "telegram_bot", "whatsapp_cloud", "slack"] as const;

type ShareRow = { grantee_project_id: string; name: string; slug: string; created_at: string };

/** Phones other tenants have shared with this one (read-only). */
function SharedWithMe({ projectId }: { projectId: string }) {
  const t = useTranslations("providers");
  const q = useQuery({
    queryKey: ["shared-with-me", projectId],
    queryFn: async () => {
      const res = await fetch(`/api/backend/api/admin/projects/${projectId}/shared-providers`);
      if (!res.ok) throw new Error(String(res.status));
      return ((await res.json()) as { data: { provider: Provider; owner: { name: string; slug: string } }[] }).data ?? [];
    },
  });
  if (!q.data?.length) return null;
  return (
    <div className="mb-5 rounded-lg border bg-muted/30 p-3">
      <div className="mb-2 flex items-center gap-2 text-sm font-medium"><Share2 className="h-4 w-4" /> {t("sharedWithMe")}</div>
      <ul className="divide-y rounded-md border bg-background text-sm">
        {q.data.map((s) => (
          <li key={s.provider.id} className="flex items-center justify-between gap-2 px-3 py-2">
            <span className="flex items-center gap-2"><span className="font-medium">{s.provider.name}</span> <ChannelBadge channel={s.provider.channel} /></span>
            <span className="text-xs text-muted-foreground">{t("from")}: {s.owner.name}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}

/** Manage which tenants may send through this phone. Uses the backend proxy
 *  directly (the share endpoints are not in the typed client). */
function ProviderShares({ projectId, providerId }: { projectId: string; providerId: string }) {
  const t = useTranslations("providers");
  const qc = useQueryClient();
  const [grantee, setGrantee] = useState("");
  const [busy, setBusy] = useState(false);
  const key = ["provider-shares", projectId, providerId];
  const shares = useQuery({
    queryKey: key,
    queryFn: async () => {
      const res = await fetch(`/api/backend/api/admin/projects/${projectId}/providers/${providerId}/shares`);
      if (!res.ok) throw new Error(String(res.status));
      return ((await res.json()) as { data: ShareRow[] }).data ?? [];
    },
  });
  const add = async () => {
    const id = grantee.trim();
    if (!id) return;
    setBusy(true);
    try {
      const res = await fetch(`/api/backend/api/admin/projects/${projectId}/providers/${providerId}/share`, {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ grantee_project_id: id }),
      });
      if (!res.ok) { const b = (await res.json().catch(() => null)) as { error?: { message?: string } } | null; toast.error(b?.error?.message ?? res.statusText); return; }
      setGrantee("");
      void qc.invalidateQueries({ queryKey: key });
    } finally { setBusy(false); }
  };
  const revoke = async (gid: string) => {
    const res = await fetch(`/api/backend/api/admin/projects/${projectId}/providers/${providerId}/share/${gid}`, { method: "DELETE" });
    if (res.ok) void qc.invalidateQueries({ queryKey: key }); else toast.error(res.statusText);
  };
  return (
    <div className="space-y-3 rounded-lg border bg-muted/40 p-4 text-sm">
      <div className="flex items-center gap-2 font-medium"><Share2 className="h-4 w-4" /> {t("shareTitle")}</div>
      <p className="text-muted-foreground">{t("shareHint")}</p>
      <div className="flex flex-wrap gap-2">
        <Input value={grantee} onChange={(e) => setGrantee(e.target.value)} placeholder={t("granteeId")} className="max-w-xs font-mono" />
        <Button size="sm" disabled={busy || !grantee.trim()} onClick={add}>{t("shareBtn")}</Button>
      </div>
      {shares.data?.length ? (
        <ul className="divide-y rounded-md border bg-background">
          {shares.data.map((s) => (
            <li key={s.grantee_project_id} className="flex items-center justify-between gap-2 px-3 py-2">
              <span><span className="font-medium">{s.name}</span> <code className="text-xs text-muted-foreground">{s.slug}</code></span>
              <Button size="icon" variant="ghost" onClick={() => revoke(s.grantee_project_id)} aria-label={t("revoke")}><Trash2 className="h-4 w-4" /></Button>
            </li>
          ))}
        </ul>
      ) : <p className="text-xs text-muted-foreground">{t("noShares")}</p>}
    </div>
  );
}
const samples: Record<string, string> = {
  http_sms: JSON.stringify({ url: "https://sms.example.tm/api/send", method: "POST", headers: { Authorization: "Bearer TOKEN" }, body_template: '{"to":"{{.To}}","text":{{.TextJSON}},"from":"{{.Sender}}"}', sender: "HABARCHY", success: { json_path: "status", json_equals: "OK" }, message_id: { json_path: "id" }, dlr: { message_id_param: "msgid", status_param: "status", delivered_values: ["DELIVRD"], failed_values: ["UNDELIV", "EXPIRED"] } }, null, 2),
  smpp: JSON.stringify({ host: "smsc.operator.tm", port: 2775, system_id: "habarchy", password: "secret", source_addr: "HABARCHY", source_ton: 5, source_npi: 0, dest_ton: 1, dest_npi: 1, enquire_link_sec: 60, request_dlr: true }, null, 2),
  android_sms: JSON.stringify({ gateway_key: "", sim_slot: -1, timeout_sec: 45, inbound_enabled: true }, null, 2),
  smtp: JSON.stringify({ host: "smtp.example.tm", port: 587, tls_mode: "starttls", username: "no-reply@example.tm", password: "secret", from_name: "Habarchy", from_email: "no-reply@example.tm" }, null, 2),
  fcm: JSON.stringify({ service_account: { type: "service_account", project_id: "my-firebase", private_key: "-----BEGIN PRIVATE KEY-----\n...", client_email: "firebase-adminsdk@my-firebase.iam.gserviceaccount.com" } }, null, 2),
  telegram_bot: JSON.stringify({ bot_token: "123456:ABC-DEF", parse_mode: "HTML" }, null, 2),
  whatsapp_cloud: JSON.stringify({ access_token: "EAAG…", phone_number_id: "123456789012345", api_version: "v20.0" }, null, 2),
  slack: JSON.stringify({ bot_token: "xoxb-…", default_channel: "C0123456789" }, null, 2),
};

type FormState = { id?: string; name: string; type: string; priority: number; is_active: boolean; rate_limit_per_sec: number; credentials: string };
const empty: FormState = { name: "", type: "http_sms", priority: 100, is_active: true, rate_limit_per_sec: 0, credentials: samples.http_sms };

export default function ProvidersPage() {
  const t = useTranslations("providers");
  const tc = useTranslations("common");
  const { projectId } = useProject();
  const qc = useQueryClient();
  const [form, setForm] = useState<FormState | null>(null);
  const [testing, setTesting] = useState<Provider | null>(null);
  const [testTo, setTestTo] = useState("");
  const [testResult, setTestResult] = useState<unknown>(null);
  const [detail, setDetail] = useState<Provider | null>(null);
  const [pairing, setPairing] = useState<GatewayPairing | null>(null);
  useEffect(() => {
    if (!detail || detail.type !== "android_sms" || !projectId) { setPairing(null); return; }
    let cancelled = false;
    void (async () => {
      const res = await api.GET("/api/admin/projects/{project_id}/providers/{provider_id}/pairing", { params: { path: { project_id: projectId, provider_id: detail.id! } } });
      if (!cancelled && res.data) setPairing(res.data.data ?? null);
    })();
    return () => { cancelled = true; };
  }, [detail, projectId]);

  const q = useQuery({
    queryKey: qk.providers(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Provider[]>(await api.GET("/api/admin/projects/{project_id}/providers", { params: { path: { project_id: projectId! } } })),
  });
  const invalidate = () => qc.invalidateQueries({ queryKey: qk.providers(projectId ?? "") });
  const save = useMutation({
    mutationFn: async (f: FormState) => {
      let creds: Record<string, unknown> | undefined;
      if (f.credentials.trim()) creds = JSON.parse(f.credentials) as Record<string, unknown>;
      const body = { name: f.name, type: f.type as Provider["type"], priority: f.priority, is_active: f.is_active, rate_limit_per_sec: f.rate_limit_per_sec, credentials: creds };
      if (f.id) return unwrap<Provider>(await api.PUT("/api/admin/projects/{project_id}/providers/{provider_id}", { params: { path: { project_id: projectId!, provider_id: f.id } }, body }));
      return unwrap<Provider>(await api.POST("/api/admin/projects/{project_id}/providers", { params: { path: { project_id: projectId! } }, body }));
    },
    onSuccess: () => { toast.success(tc("save")); void invalidate(); setForm(null); },
    onError: (e: Error) => toast.error(e instanceof ApiError ? `${e.message} ${e.details ? JSON.stringify(e.details) : ""}` : e.message),
  });
  const del = useMutation({
    mutationFn: async (id: string) => { await api.DELETE("/api/admin/projects/{project_id}/providers/{provider_id}", { params: { path: { project_id: projectId!, provider_id: id } } }); },
    onSuccess: () => void invalidate(),
  });
  const test = useMutation({
    mutationFn: async () => unwrap<unknown>(await api.POST("/api/admin/projects/{project_id}/providers/{provider_id}/test", { params: { path: { project_id: projectId!, provider_id: testing!.id! } }, body: { to: testTo } })),
    onSuccess: (r) => setTestResult(r),
    onError: (e: ApiError) => toast.error(e.message),
  });
  const openDetail = async (p: Provider) => {
    const full = unwrap<Provider>(await api.GET("/api/admin/projects/{project_id}/providers/{provider_id}", { params: { path: { project_id: projectId!, provider_id: p.id! } } }));
    setDetail(full);
  };

  return (
    <>
      <PageHeader title={t("title")} actions={<Button size="sm" onClick={() => setForm(empty)} data-testid="new-provider"><Plus /> {t("new")}</Button>} />
      {projectId ? <SharedWithMe projectId={projectId} /> : null}
      {q.isLoading ? <Skeleton className="h-48" /> : !q.data?.length ? <EmptyState /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow><TableHead>{tc("name")}</TableHead><TableHead>{t("channel")}</TableHead><TableHead>{t("type")}</TableHead><TableHead>{t("priority")}</TableHead><TableHead>{t("rateLimit")}</TableHead><TableHead>{tc("status")}</TableHead><TableHead className="text-right">{tc("actions")}</TableHead></TableRow></TableHeader>
            <TableBody>
              {q.data.map((p) => (
                <TableRow key={p.id}>
                  <TableCell className="font-medium"><button className="hover:underline" onClick={() => openDetail(p)}>{p.name}</button></TableCell>
                  <TableCell><ChannelBadge channel={p.channel} /></TableCell>
                  <TableCell className="font-mono text-xs">{p.type}</TableCell>
                  <TableCell className="tabular">{p.priority}</TableCell>
                  <TableCell className="tabular">{p.rate_limit_per_sec || "∞"}</TableCell>
                  <TableCell><Badge variant={p.is_active ? "success" : "outline"}>{p.is_active ? t("active") : tc("no")}</Badge></TableCell>
                  <TableCell className="space-x-1 text-right whitespace-nowrap">
                    <Button size="sm" variant="ghost" onClick={() => { setTesting(p); setTestResult(null); }}><FlaskConical /> {t("testSend")}</Button>
                    <Button size="sm" variant="ghost" onClick={() => setForm({ id: p.id, name: p.name ?? "", type: p.type ?? "http_sms", priority: p.priority ?? 100, is_active: p.is_active ?? true, rate_limit_per_sec: p.rate_limit_per_sec ?? 0, credentials: "" })}><Pencil /></Button>
                    <ConfirmButton size="sm" variant="ghost" title={tc("confirmDelete", { name: p.name ?? "" })} onConfirm={() => del.mutate(p.id!)}>{tc("delete")}</ConfirmButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      <Dialog open={!!form} onOpenChange={(o) => !o && setForm(null)}>
        <DialogContent wide>
          <DialogHeader><DialogTitle>{form?.id ? tc("edit") : t("new")}</DialogTitle><DialogDescription>{t("credentialsHint")}</DialogDescription></DialogHeader>
          {form ? (
            <div className="grid gap-3 sm:grid-cols-2">
              <Field label={tc("name")} htmlFor="p-name"><Input id="p-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></Field>
              <Field label={t("type")}>
                <Select value={form.type} disabled={!!form.id} onValueChange={(v) => setForm({ ...form, type: v, credentials: form.id ? form.credentials : samples[v] })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>{types.map((x) => <SelectItem key={x} value={x}>{x}</SelectItem>)}</SelectContent>
                </Select>
              </Field>
              <Field label={t("priority")} htmlFor="p-prio" hint={t("priorityHint")}><Input id="p-prio" type="number" value={form.priority} onChange={(e) => setForm({ ...form, priority: Number(e.target.value) })} /></Field>
              <Field label={t("rateLimit")} htmlFor="p-rate"><Input id="p-rate" type="number" value={form.rate_limit_per_sec} onChange={(e) => setForm({ ...form, rate_limit_per_sec: Number(e.target.value) })} /></Field>
              <div className="flex items-center gap-2 sm:col-span-2"><Switch id="p-active" checked={form.is_active} onCheckedChange={(v) => setForm({ ...form, is_active: v })} /><Label htmlFor="p-active">{t("active")}</Label></div>
              <div className="sm:col-span-2">
                <Field label={t("credentials")} htmlFor="p-creds"><Textarea id="p-creds" rows={12} className="font-mono text-xs" value={form.credentials} onChange={(e) => setForm({ ...form, credentials: e.target.value })} placeholder={form.id ? "{ … }" : undefined} /></Field>
              </div>
            </div>
          ) : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setForm(null)}>{tc("cancel")}</Button>
            <Button onClick={() => form && save.mutate(form)} disabled={save.isPending || !form?.name}>{tc("save")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!testing} onOpenChange={(o) => !o && setTesting(null)}>
        <DialogContent>
          <DialogHeader><DialogTitle>{t("testSend")} · {testing?.name}</DialogTitle></DialogHeader>
          <Field label={t("testTo")} htmlFor="test-to"><Input id="test-to" value={testTo} onChange={(e) => setTestTo(e.target.value)} placeholder="+99365123456 / user@example.tm / chat id / fcm token" /></Field>
          {testResult ? <JsonView value={testResult} /> : null}
          <DialogFooter>
            <Button onClick={() => test.mutate()} disabled={test.isPending || !testTo}>{t("testSend")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!detail} onOpenChange={(o) => !o && setDetail(null)}>
        <DialogContent wide>
          <DialogHeader><DialogTitle>{detail?.name}</DialogTitle><DialogDescription>{t("settings")}</DialogDescription></DialogHeader>
          {detail?.type === "android_sms" ? (
            <div className="space-y-3 rounded-lg border bg-muted/40 p-4 text-sm" data-testid="pairing">
              <div className="flex items-center gap-2 font-medium"><Smartphone className="h-4 w-4" /> {t("pairing")}
                {pairing ? <Badge variant={pairing.online ? "success" : "destructive"}>{pairing.online ? t("online") : t("offline")}</Badge> : null}
              </div>
              <p className="text-muted-foreground">{t("pairingHint")}</p>
              {pairing ? (
                <div className="grid gap-2 sm:grid-cols-[auto_1fr_auto] sm:items-center">
                  <span className="text-muted-foreground">API URL</span>
                  <code className="truncate rounded bg-background px-1.5 py-0.5 text-xs">{pairing.api_url}</code>
                  <CopyButton size="icon" value={pairing.api_url ?? ""} />
                  <span className="text-muted-foreground">{t("gatewayKey")}</span>
                  <code className="truncate rounded bg-background px-1.5 py-0.5 text-xs" data-testid="gateway-key">{pairing.gateway_key}</code>
                  <CopyButton size="icon" value={pairing.gateway_key ?? ""} />
                  <span className="text-muted-foreground">QR</span>
                  <code className="truncate rounded bg-background px-1.5 py-0.5 text-xs">{pairing.qr}</code>
                  <CopyButton size="icon" value={pairing.qr ?? ""} />
                  {pairing.last_seen_at ? <><span className="text-muted-foreground">{t("lastSeen")}</span><span className="sm:col-span-2">{new Date(pairing.last_seen_at).toLocaleString()}{pairing.device_info && typeof pairing.device_info === "object" ? ` · ${Object.entries(pairing.device_info as Record<string, unknown>).filter(([k]) => ["battery", "network", "operator", "model"].includes(k)).map(([k, v]) => `${k}: ${String(v)}`).join(" · ")}` : ""}</span></> : null}
                </div>
              ) : <Skeleton className="h-16" />}
            </div>
          ) : null}
          {detail?.type === "android_sms" && projectId && detail.id ? (
            <ProviderShares projectId={projectId} providerId={detail.id} />
          ) : null}
          {detail?.type === "http_sms" ? (
            <div className="flex items-center gap-2 text-sm">
              <span className="text-muted-foreground">{t("dlrUrl")}:</span>
              <code className="truncate rounded bg-muted px-1.5 py-0.5 text-xs">{`${typeof window !== "undefined" ? window.location.origin.replace(/:\d+$/, ":8080") : ""}/callbacks/sms/${detail.id}`}</code>
              <CopyButton size="icon" value={`/callbacks/sms/${detail.id}`} />
            </div>
          ) : null}
          <JsonView value={detail?.settings ?? {}} />
        </DialogContent>
      </Dialog>
    </>
  );
}
