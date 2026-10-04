"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Check, Plus } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Project } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export default function ProjectsPage() {
  const t = useTranslations("projects");
  const tc = useTranslations("common");
  const { projects, projectId, setProjectId } = useProject();
  const router = useRouter();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState({ name: "", slug: "", daily_quota: 0, monthly_quota: 0, default_locale: "tk" });
  const create = useMutation({
    mutationFn: async () => unwrap<Project>(await api.POST("/api/admin/projects", { body: { ...form, default_locale: form.default_locale as Project["default_locale"] as never } })),
    onSuccess: async (p) => {
      toast.success(t("created"));
      await qc.invalidateQueries({ queryKey: qk.projects });
      setProjectId(p.id!);
      setOpen(false);
      router.push("/");
    },
    onError: (e: ApiError) => toast.error(e.message + (e.details ? ` ${JSON.stringify(e.details)}` : "")),
  });
  return (
    <>
      <PageHeader title={t("title")} actions={<Button size="sm" onClick={() => setOpen(true)} data-testid="new-project"><Plus /> {t("new")}</Button>} />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {projects.map((p) => (
          <Card key={p.id} className={p.id === projectId ? "border-primary" : undefined}>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">{p.name}{p.id === projectId ? <Badge>{t("current")}</Badge> : null}</CardTitle>
              <CardDescription className="font-mono text-xs">{p.slug}</CardDescription>
            </CardHeader>
            <CardContent className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">{t("role")}: <b>{p.role}</b></span>
              {p.id !== projectId ? <Button size="sm" variant="outline" onClick={() => { setProjectId(p.id!); router.push("/"); }}>{t("switch")}</Button> : <Check className="h-4 w-4 text-primary" />}
            </CardContent>
          </Card>
        ))}
      </div>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader><DialogTitle>{t("new")}</DialogTitle></DialogHeader>
          <div className="grid gap-3">
            <Field label={t("name")} htmlFor="pr-name"><Input id="pr-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} data-testid="project-name" /></Field>
            <Field label={t("slug")} htmlFor="pr-slug" hint="auto"><Input id="pr-slug" value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} /></Field>
            <div className="grid grid-cols-2 gap-3">
              <Field label={t("dailyQuota")} htmlFor="pr-dq" hint={t("unlimited")}><Input id="pr-dq" type="number" value={form.daily_quota} onChange={(e) => setForm({ ...form, daily_quota: Number(e.target.value) })} /></Field>
              <Field label={t("monthlyQuota")} htmlFor="pr-mq" hint={t("unlimited")}><Input id="pr-mq" type="number" value={form.monthly_quota} onChange={(e) => setForm({ ...form, monthly_quota: Number(e.target.value) })} /></Field>
            </div>
            <Field label={t("defaultLocale")}>
              <Select value={form.default_locale} onValueChange={(v) => setForm({ ...form, default_locale: v })}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>{["tk", "ru", "en"].map((l) => <SelectItem key={l} value={l}>{l.toUpperCase()}</SelectItem>)}</SelectContent>
              </Select>
            </Field>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>{tc("cancel")}</Button>
            <Button onClick={() => create.mutate()} disabled={create.isPending || form.name.length < 2} data-testid="create-project">{tc("create")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
