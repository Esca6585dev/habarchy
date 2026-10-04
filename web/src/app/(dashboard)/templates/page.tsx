"use client";
import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Template } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge } from "@/components/common/status-badge";
import { EmptyState } from "@/components/common/empty-state";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

export default function TemplatesPage() {
  const t = useTranslations("templates");
  const tc = useTranslations("common");
  const tch = useTranslations("channel");
  const { projectId, can } = useProject();
  const router = useRouter();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ key: "", channel: "sms", locale: "tk", subject: "", body: "" });
  const q = useQuery({
    queryKey: qk.templates(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Template[]>(await api.GET("/api/admin/projects/{project_id}/templates", { params: { path: { project_id: projectId! } } })),
  });
  const create = useMutation({
    mutationFn: async () =>
      unwrap<Template>(await api.POST("/api/admin/projects/{project_id}/templates", {
        params: { path: { project_id: projectId! } },
        body: { key: form.key, channel: form.channel as "sms", locale: form.locale as "tk", subject: form.subject, body: form.body },
      })),
    onSuccess: (tpl) => {
      toast.success(t("saved"));
      void qc.invalidateQueries({ queryKey: qk.templates(projectId ?? "") });
      setOpen(false);
      router.push(`/templates/${tpl.id}`);
    },
    onError: (e: ApiError) => toast.error(e.message + (e.details ? ` ${JSON.stringify(e.details)}` : "")),
  });

  return (
    <>
      <PageHeader
        title={t("title")}
        actions={can("developer") ? (
          <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild><Button size="sm" data-testid="new-template"><Plus /> {t("new")}</Button></DialogTrigger>
            <DialogContent>
              <DialogHeader><DialogTitle>{t("new")}</DialogTitle></DialogHeader>
              <div className="grid gap-3">
                <Field label={t("key")} htmlFor="tpl-key" hint="otp, welcome, password_reset…"><Input id="tpl-key" value={form.key} onChange={(e) => setForm({ ...form, key: e.target.value })} /></Field>
                <div className="grid grid-cols-2 gap-3">
                  <Field label={t("channel")}>
                    <Select value={form.channel} onValueChange={(v) => setForm({ ...form, channel: v })}>
                      <SelectTrigger><SelectValue /></SelectTrigger>
                      <SelectContent>{["sms", "email", "push", "telegram"].map((c) => <SelectItem key={c} value={c}>{tch(c as "sms")}</SelectItem>)}</SelectContent>
                    </Select>
                  </Field>
                  <Field label={t("locale")}>
                    <Select value={form.locale} onValueChange={(v) => setForm({ ...form, locale: v })}>
                      <SelectTrigger><SelectValue /></SelectTrigger>
                      <SelectContent>{["tk", "ru", "en"].map((l) => <SelectItem key={l} value={l}>{l.toUpperCase()}</SelectItem>)}</SelectContent>
                    </Select>
                  </Field>
                </div>
                {form.channel !== "sms" ? <Field label={t("subject")} htmlFor="tpl-subject"><Input id="tpl-subject" value={form.subject} onChange={(e) => setForm({ ...form, subject: e.target.value })} /></Field> : null}
                <Field label={t("body")} htmlFor="tpl-body" hint={t("variablesHint")}><Textarea id="tpl-body" rows={5} value={form.body} onChange={(e) => setForm({ ...form, body: e.target.value })} /></Field>
              </div>
              <DialogFooter>
                <Button variant="outline" onClick={() => setOpen(false)}>{tc("cancel")}</Button>
                <Button onClick={() => create.mutate()} disabled={create.isPending || !form.key || !form.body}>{tc("create")}</Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        ) : null}
      />
      {q.isLoading ? <Skeleton className="h-48" /> : !q.data?.length ? <EmptyState /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow><TableHead>{t("key")}</TableHead><TableHead>{t("channel")}</TableHead><TableHead>{t("locale")}</TableHead><TableHead>{t("variables")}</TableHead><TableHead>{t("version", { n: "" })}</TableHead><TableHead>{tc("status")}</TableHead><TableHead>{tc("updated")}</TableHead></TableRow></TableHeader>
            <TableBody>
              {q.data.map((tpl) => (
                <TableRow key={tpl.id}>
                  <TableCell className="font-medium"><Link href={`/templates/${tpl.id}`} className="hover:underline">{tpl.key}</Link></TableCell>
                  <TableCell><ChannelBadge channel={tpl.channel} /></TableCell>
                  <TableCell className="uppercase">{tpl.locale}</TableCell>
                  <TableCell className="space-x-1">{(tpl.required_vars ?? []).map((v) => <Badge key={v} variant="outline" className="font-mono">{v}</Badge>)}</TableCell>
                  <TableCell className="tabular">v{tpl.version}</TableCell>
                  <TableCell><Badge variant={tpl.is_active ? "success" : "outline"}>{tpl.is_active ? t("active") : t("inactive")}</Badge></TableCell>
                  <TableCell className="tabular text-xs text-muted-foreground">{formatDate(tpl.updated_at)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </>
  );
}
