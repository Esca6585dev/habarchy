"use client";
import { use, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, History, Save } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError, errorOf } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Preview, Template, TemplateVersion } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge } from "@/components/common/status-badge";
import { ConfirmButton } from "@/components/common/confirm-button";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

export default function TemplateEditorPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const t = useTranslations("templates");
  const tc = useTranslations("common");
  const { projectId, can } = useProject();
  const router = useRouter();
  const qc = useQueryClient();
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");
  const [active, setActive] = useState(true);
  const [sample, setSample] = useState('{\n  "code": "4821"\n}');
  const [preview, setPreview] = useState<Preview | null>(null);
  const [previewErr, setPreviewErr] = useState<string | null>(null);

  const q = useQuery({
    queryKey: qk.template(projectId ?? "", id),
    enabled: !!projectId,
    queryFn: async () => unwrap<Template>(await api.GET("/api/admin/projects/{project_id}/templates/{template_id}", { params: { path: { project_id: projectId!, template_id: id } } })),
  });
  const versions = useQuery({
    queryKey: qk.versions(projectId ?? "", id),
    enabled: !!projectId,
    queryFn: async () => unwrap<TemplateVersion[]>(await api.GET("/api/admin/projects/{project_id}/templates/{template_id}/versions", { params: { path: { project_id: projectId!, template_id: id } } })),
  });
  useEffect(() => {
    if (q.data) {
      setSubject(q.data.subject ?? "");
      setBody(q.data.body ?? "");
      setActive(q.data.is_active ?? true);
    }
  }, [q.data]);

  const sampleData = useMemo(() => {
    try {
      return JSON.parse(sample) as Record<string, unknown>;
    } catch {
      return null;
    }
  }, [sample]);

  // Live preview, debounced.
  useEffect(() => {
    if (!projectId || !q.data || !body) return;
    const h = setTimeout(async () => {
      const res = await api.POST("/api/admin/projects/{project_id}/templates/preview", {
        params: { path: { project_id: projectId } },
        body: { channel: q.data!.channel ?? "sms", subject, body, data: sampleData ?? {} },
      });
      if (res.error || !res.response.ok) {
        const err = errorOf(res);
        setPreviewErr((err?.details?.body as string | undefined) ?? err?.message ?? "error");
        setPreview(null);
      } else {
        setPreviewErr(null);
        setPreview((res.data as { data: Preview }).data);
      }
    }, 350);
    return () => clearTimeout(h);
  }, [body, subject, sampleData, projectId, q.data]);

  const save = useMutation({
    mutationFn: async () => unwrap<Template>(await api.PUT("/api/admin/projects/{project_id}/templates/{template_id}", { params: { path: { project_id: projectId!, template_id: id } }, body: { subject, body, is_active: active } })),
    onSuccess: () => {
      toast.success(t("saved"));
      void qc.invalidateQueries({ queryKey: qk.template(projectId ?? "", id) });
      void qc.invalidateQueries({ queryKey: qk.versions(projectId ?? "", id) });
      void qc.invalidateQueries({ queryKey: qk.templates(projectId ?? "") });
    },
    onError: (e: ApiError) => toast.error(e.message + (e.details ? ` ${JSON.stringify(e.details)}` : "")),
  });
  const restore = useMutation({
    mutationFn: async (version: number) => unwrap<Template>(await api.POST("/api/admin/projects/{project_id}/templates/{template_id}/versions/{version}/restore", { params: { path: { project_id: projectId!, template_id: id, version } } })),
    onSuccess: () => {
      toast.success(t("saved"));
      void qc.invalidateQueries({ queryKey: qk.template(projectId ?? "", id) });
      void qc.invalidateQueries({ queryKey: qk.versions(projectId ?? "", id) });
    },
  });
  const del = useMutation({
    mutationFn: async () => { await api.DELETE("/api/admin/projects/{project_id}/templates/{template_id}", { params: { path: { project_id: projectId!, template_id: id } } }); },
    onSuccess: () => { void qc.invalidateQueries({ queryKey: qk.templates(projectId ?? "") }); router.push("/templates"); },
  });

  if (q.isLoading || !q.data) return <Skeleton className="h-64" />;
  const tpl = q.data;
  const editable = can("developer");
  return (
    <>
      <PageHeader
        title={tpl.key ?? ""}
        description={`${t("version", { n: tpl.version ?? 1 })} · ${tpl.locale?.toUpperCase()}`}
        actions={
          <>
            <ChannelBadge channel={tpl.channel} />
            <Button variant="outline" size="sm" asChild><Link href="/templates"><ArrowLeft /> {t("title")}</Link></Button>
            {editable ? <ConfirmButton variant="destructive" size="sm" title={tc("confirmDelete", { name: tpl.key })} onConfirm={() => del.mutate()}>{tc("delete")}</ConfirmButton> : null}
            {editable ? <Button size="sm" onClick={() => save.mutate()} disabled={save.isPending} data-testid="save-template"><Save /> {tc("save")}</Button> : null}
          </>
        }
      />
      <div className="grid gap-6 lg:grid-cols-5">
        <Card className="lg:col-span-3">
          <CardHeader><CardTitle>{t("editor")}</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            {tpl.channel !== "sms" ? <Field label={t("subject")} htmlFor="subject"><Input id="subject" value={subject} onChange={(e) => setSubject(e.target.value)} disabled={!editable} /></Field> : null}
            <Field label={t("body")} htmlFor="body" hint={t("variablesHint")}>
              <Textarea id="body" rows={12} className="font-mono text-sm" value={body} onChange={(e) => setBody(e.target.value)} disabled={!editable} />
            </Field>
            <div className="flex items-center gap-2"><Switch id="active" checked={active} onCheckedChange={setActive} disabled={!editable} /><Label htmlFor="active">{t("active")}</Label></div>
          </CardContent>
        </Card>
        <div className="space-y-6 lg:col-span-2">
          <Card>
            <CardHeader><CardTitle>{t("variables")}</CardTitle></CardHeader>
            <CardContent className="space-y-3">
              <div className="flex flex-wrap gap-1">
                {(preview?.required_vars ?? tpl.required_vars ?? []).map((v) => (
                  <Badge key={v} variant={preview?.missing_vars?.includes(v) ? "destructive" : "outline"} className="font-mono">{v}</Badge>
                ))}
              </div>
              <Field label={t("sampleData")} htmlFor="sample" error={sampleData ? undefined : "invalid JSON"}>
                <Textarea id="sample" rows={5} className="font-mono text-xs" value={sample} onChange={(e) => setSample(e.target.value)} />
              </Field>
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle>{t("preview")}</CardTitle></CardHeader>
            <CardContent>
              {previewErr ? <p className="text-sm text-destructive">{previewErr}</p> : null}
              {preview ? (
                <div className="space-y-2">
                  {preview.subject ? <p className="text-sm font-medium">{preview.subject}</p> : null}
                  {tpl.channel === "email" ? (
                    <iframe title="preview" sandbox="" srcDoc={preview.body} className="h-56 w-full rounded-md border bg-white" />
                  ) : (
                    <div className="rounded-2xl rounded-bl-sm bg-primary/10 px-3 py-2 text-sm whitespace-pre-wrap">{preview.body}</div>
                  )}
                  {preview.missing_vars?.length ? <p className="text-xs text-destructive">{t("missing", { vars: preview.missing_vars.join(", ") })}</p> : null}
                </div>
              ) : null}
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle className="flex items-center gap-2"><History className="h-4 w-4" /> {t("history")}</CardTitle></CardHeader>
            <CardContent>
              <ul className="divide-y text-sm">
                {(versions.data ?? []).map((v) => (
                  <li key={v.version} className="flex items-center gap-3 py-2">
                    <span className="tabular font-medium">v{v.version}</span>
                    <span className="min-w-0 flex-1 truncate text-xs text-muted-foreground">{v.body}</span>
                    <span className="tabular text-xs text-muted-foreground">{formatDate(v.created_at, { dateStyle: "short", timeStyle: "short" })}</span>
                    {editable && v.version !== tpl.version ? <Button size="sm" variant="ghost" onClick={() => restore.mutate(v.version!)}>{t("restore")}</Button> : null}
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
