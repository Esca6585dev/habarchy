"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { api, unwrap, ApiError, ensureOk } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Member, Project } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ConfirmButton } from "@/components/common/confirm-button";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

const channels = ["push", "sms", "email", "telegram"] as const;
const roles = ["owner", "admin", "developer", "viewer"] as const;

export default function SettingsPage() {
  const t = useTranslations("projects");
  const tc = useTranslations("common");
  const tch = useTranslations("channel");
  const { project, projectId, can, setProjectId, projects } = useProject();
  const qc = useQueryClient();
  const router = useRouter();
  const [form, setForm] = useState({ name: "", status: "active", daily_quota: 0, monthly_quota: 0, webhook_url: "", webhook_secret: "", default_locale: "tk", allowed_ips: "", auto_channel_order: ["push", "sms", "email"] as string[] });
  const [member, setMember] = useState({ email: "", role: "developer" });
  useEffect(() => {
    if (project) {
      setForm({
        name: project.name ?? "", status: project.status ?? "active", daily_quota: project.daily_quota ?? 0, monthly_quota: project.monthly_quota ?? 0,
        webhook_url: project.webhook_url ?? "", webhook_secret: "", default_locale: project.default_locale ?? "tk",
        allowed_ips: (project.allowed_ips ?? []).join("\n"), auto_channel_order: (project.auto_channel_order ?? []) as string[],
      });
    }
  }, [project]);
  const members = useQuery({
    queryKey: qk.members(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Member[]>(await api.GET("/api/admin/projects/{project_id}/members", { params: { path: { project_id: projectId! } } })),
  });
  const save = useMutation({
    mutationFn: async () =>
      unwrap<Project>(await api.PATCH("/api/admin/projects/{project_id}", {
        params: { path: { project_id: projectId! } },
        body: {
          name: form.name, status: form.status as "active", daily_quota: form.daily_quota, monthly_quota: form.monthly_quota, webhook_url: form.webhook_url,
          webhook_secret: form.webhook_secret || undefined, default_locale: form.default_locale as "tk",
          allowed_ips: form.allowed_ips.split("\n").map((s) => s.trim()).filter(Boolean), auto_channel_order: form.auto_channel_order as Array<"sms">,
        },
      })),
    onSuccess: () => { toast.success(t("saved")); void qc.invalidateQueries({ queryKey: qk.projects }); setForm((f) => ({ ...f, webhook_secret: "" })); },
    onError: (e: ApiError) => toast.error(e.message + (e.details ? ` ${JSON.stringify(e.details)}` : "")),
  });
  const setMemberM = useMutation({
    mutationFn: async () => unwrap<Member>(await api.PUT("/api/admin/projects/{project_id}/members", { params: { path: { project_id: projectId! } }, body: { email: member.email, role: member.role as "viewer" } })),
    onSuccess: () => { void qc.invalidateQueries({ queryKey: qk.members(projectId ?? "") }); setMember({ email: "", role: "developer" }); },
    onError: (e: ApiError) => toast.error(e.message),
  });
  const removeM = useMutation({
    mutationFn: async (userId: string) => { ensureOk(await api.DELETE("/api/admin/projects/{project_id}/members/{user_id}", { params: { path: { project_id: projectId!, user_id: userId } } })); },
    onSuccess: () => void qc.invalidateQueries({ queryKey: qk.members(projectId ?? "") }),
    onError: (e: ApiError) => toast.error(e.message),
  });
  const del = useMutation({
    mutationFn: async () => { await api.DELETE("/api/admin/projects/{project_id}", { params: { path: { project_id: projectId! } } }); },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: qk.projects });
      const next = projects.find((p) => p.id !== projectId);
      if (next) setProjectId(next.id!);
      router.push("/projects");
    },
  });
  if (!project) return null;
  const move = (i: number, dir: -1 | 1) => {
    const arr = [...form.auto_channel_order];
    const j = i + dir;
    if (j < 0 || j >= arr.length) return;
    [arr[i], arr[j]] = [arr[j], arr[i]];
    setForm({ ...form, auto_channel_order: arr });
  };
  return (
    <>
      <PageHeader title={t("settings")} description={project.slug} actions={can("admin") ? <Button size="sm" onClick={() => save.mutate()} disabled={save.isPending}>{tc("save")}</Button> : null} />
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader><CardTitle>{t("settings")}</CardTitle></CardHeader>
          <CardContent className="grid gap-3">
            <Field label={t("name")} htmlFor="s-name"><Input id="s-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></Field>
            <div className="grid grid-cols-2 gap-3">
              <Field label={t("statusLabel")}>
                <Select value={form.status} onValueChange={(v) => setForm({ ...form, status: v })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>{["active", "suspended", "archived"].map((s) => <SelectItem key={s} value={s}>{s}</SelectItem>)}</SelectContent>
                </Select>
              </Field>
              <Field label={t("defaultLocale")}>
                <Select value={form.default_locale} onValueChange={(v) => setForm({ ...form, default_locale: v })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>{["tk", "ru", "en"].map((l) => <SelectItem key={l} value={l}>{l.toUpperCase()}</SelectItem>)}</SelectContent>
                </Select>
              </Field>
              <Field label={t("dailyQuota")} htmlFor="s-dq" hint={t("unlimited")}><Input id="s-dq" type="number" value={form.daily_quota} onChange={(e) => setForm({ ...form, daily_quota: Number(e.target.value) })} /></Field>
              <Field label={t("monthlyQuota")} htmlFor="s-mq" hint={t("unlimited")}><Input id="s-mq" type="number" value={form.monthly_quota} onChange={(e) => setForm({ ...form, monthly_quota: Number(e.target.value) })} /></Field>
            </div>
            <Field label={t("webhookUrl")} htmlFor="s-wh"><Input id="s-wh" placeholder="https://" value={form.webhook_url} onChange={(e) => setForm({ ...form, webhook_url: e.target.value })} /></Field>
            <Field label={t("webhookSecret")} htmlFor="s-ws" hint={project.has_webhook_secret ? `✓ ${t("webhookSecretHint")}` : t("webhookSecretHint")}><Input id="s-ws" type="password" autoComplete="off" value={form.webhook_secret} onChange={(e) => setForm({ ...form, webhook_secret: e.target.value })} /></Field>
            <Field label={t("allowedIps")} htmlFor="s-ips"><Textarea id="s-ips" rows={2} value={form.allowed_ips} onChange={(e) => setForm({ ...form, allowed_ips: e.target.value })} /></Field>
            <Field label={t("autoOrder")}>
              <ul className="space-y-1">
                {form.auto_channel_order.map((c, i) => (
                  <li key={c} className="flex items-center gap-2 rounded-md border px-2 py-1 text-sm">
                    <span className="tabular w-4 text-muted-foreground">{i + 1}</span><span className="flex-1">{tch(c as "sms")}</span>
                    <Button size="sm" variant="ghost" onClick={() => move(i, -1)} disabled={i === 0}>↑</Button>
                    <Button size="sm" variant="ghost" onClick={() => move(i, 1)} disabled={i === form.auto_channel_order.length - 1}>↓</Button>
                    <Button size="sm" variant="ghost" onClick={() => setForm({ ...form, auto_channel_order: form.auto_channel_order.filter((x) => x !== c) })}>×</Button>
                  </li>
                ))}
                {channels.filter((c) => !form.auto_channel_order.includes(c)).map((c) => (
                  <li key={c}><Button size="sm" variant="outline" onClick={() => setForm({ ...form, auto_channel_order: [...form.auto_channel_order, c] })}>+ {tch(c)}</Button></li>
                ))}
              </ul>
            </Field>
          </CardContent>
        </Card>
        <div className="space-y-6">
          <Card>
            <CardHeader><CardTitle>{t("members")}</CardTitle></CardHeader>
            <CardContent className="space-y-3">
              <Table>
                <TableHeader><TableRow><TableHead>Email</TableHead><TableHead>{tc("role")}</TableHead><TableHead /></TableRow></TableHeader>
                <TableBody>
                  {(members.data ?? []).map((m) => (
                    <TableRow key={m.user_id}>
                      <TableCell><div>{m.email}</div><div className="text-xs text-muted-foreground">{m.full_name}</div></TableCell>
                      <TableCell><Badge variant="secondary">{m.role}</Badge></TableCell>
                      <TableCell className="text-right">{can("admin") ? <ConfirmButton size="sm" variant="ghost" title={tc("confirmDelete", { name: m.email ?? "" })} confirmLabel={t("remove")} onConfirm={() => removeM.mutate(m.user_id!)}>{t("remove")}</ConfirmButton> : null}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              {can("admin") ? (
                <form className="flex flex-wrap items-end gap-2" onSubmit={(e) => { e.preventDefault(); setMemberM.mutate(); }}>
                  <div className="min-w-[12rem] flex-1"><Field label={t("memberEmail")} htmlFor="m-email"><Input id="m-email" type="email" value={member.email} onChange={(e) => setMember({ ...member, email: e.target.value })} /></Field></div>
                  <Select value={member.role} onValueChange={(v) => setMember({ ...member, role: v })}>
                    <SelectTrigger className="w-36"><SelectValue /></SelectTrigger>
                    <SelectContent>{roles.filter((r) => r !== "owner" || can("owner")).map((r) => <SelectItem key={r} value={r}>{r}</SelectItem>)}</SelectContent>
                  </Select>
                  <Button type="submit" disabled={!member.email || setMemberM.isPending}>{t("addMember")}</Button>
                </form>
              ) : null}
            </CardContent>
          </Card>
          {can("owner") ? (
            <Card className="border-destructive/40">
              <CardHeader><CardTitle className="text-destructive">{t("dangerZone")}</CardTitle><CardDescription>{t("deleteHint")}</CardDescription></CardHeader>
              <CardContent><ConfirmButton variant="destructive" title={tc("confirmDelete", { name: project.name ?? "" })} description={t("deleteHint")} onConfirm={() => del.mutate()}>{t("deleteProject")}</ConfirmButton></CardContent>
            </Card>
          ) : null}
        </div>
      </div>
    </>
  );
}
