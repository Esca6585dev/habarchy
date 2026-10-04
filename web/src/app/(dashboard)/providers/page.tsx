"use client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { FlaskConical, Pencil, Plus } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Provider } from "@/lib/api/types";
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

const types = ["http_sms", "smpp", "smtp", "fcm", "telegram_bot"] as const;
const samples: Record<string, string> = {
  http_sms: JSON.stringify({ url: "https://sms.example.tm/api/send", method: "POST", headers: { Authorization: "Bearer TOKEN" }, body_template: '{"to":"{{.To}}","text":{{.TextJSON}},"from":"{{.Sender}}"}', sender: "HABARCHY", success: { json_path: "status", json_equals: "OK" }, message_id: { json_path: "id" }, dlr: { message_id_param: "msgid", status_param: "status", delivered_values: ["DELIVRD"], failed_values: ["UNDELIV", "EXPIRED"] } }, null, 2),
  smpp: JSON.stringify({ host: "smsc.operator.tm", port: 2775, system_id: "habarchy", password: "secret", source_addr: "HABARCHY", source_ton: 5, source_npi: 0, dest_ton: 1, dest_npi: 1, enquire_link_sec: 60, request_dlr: true }, null, 2),
  smtp: JSON.stringify({ host: "smtp.example.tm", port: 587, tls_mode: "starttls", username: "no-reply@example.tm", password: "secret", from_name: "Habarchy", from_email: "no-reply@example.tm" }, null, 2),
  fcm: JSON.stringify({ service_account: { type: "service_account", project_id: "my-firebase", private_key: "-----BEGIN PRIVATE KEY-----\n...", client_email: "firebase-adminsdk@my-firebase.iam.gserviceaccount.com" } }, null, 2),
  telegram_bot: JSON.stringify({ bot_token: "123456:ABC-DEF", parse_mode: "HTML" }, null, 2),
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
