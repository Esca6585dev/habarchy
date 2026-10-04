"use client";
import { useEffect, useRef, useState } from "react";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Apple, FileSpreadsheet, Upload } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Group, ImportResult } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

/** Import contacts from a file, Google Contacts or iCloud / CardDAV. */
export function ImportDialog({ open, onOpenChange, defaultGroupId }: { open: boolean; onOpenChange: (o: boolean) => void; defaultGroupId?: string }) {
  const t = useTranslations("contacts");
  const { projectId } = useProject();
  const qc = useQueryClient();
  const params = useSearchParams();
  const [groupId, setGroupId] = useState(defaultGroupId ?? "");
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<ImportResult | null>(null);
  const [busy, setBusy] = useState(false);
  const [apple, setApple] = useState({ server_url: "https://contacts.icloud.com/", username: "", password: "" });
  const fileRef = useRef<HTMLInputElement>(null);

  const groups = useQuery({
    queryKey: qk.groups(projectId ?? ""),
    enabled: !!projectId && open,
    queryFn: async () => unwrap<Group[]>(await api.GET("/api/admin/projects/{project_id}/groups", { params: { path: { project_id: projectId! } } })),
  });

  // Result of a Google round-trip arrives in the query string.
  useEffect(() => {
    const imported = params.get("imported");
    const err = params.get("import_error");
    if (imported) toast.success(t("importedToast", { n: Number(imported) }));
    if (err) toast.error(err);
    if (imported || err) {
      void qc.invalidateQueries({ queryKey: ["contacts", projectId] });
      window.history.replaceState(null, "", window.location.pathname);
    }
  }, [params, projectId, qc, t]);

  const done = (r: ImportResult) => {
    setResult(r);
    if (!r.dry_run) {
      toast.success(t("importedToast", { n: (r.created ?? 0) + (r.updated ?? 0) }));
      void qc.invalidateQueries({ queryKey: ["contacts", projectId] });
      void qc.invalidateQueries({ queryKey: qk.groups(projectId ?? "") });
    }
  };
  const fail = async (res: Response) => {
    const body = (await res.json().catch(() => null)) as { error?: { message?: string; details?: Record<string, unknown> } } | null;
    toast.error(`${body?.error?.message ?? res.statusText} ${body?.error?.details ? JSON.stringify(body.error.details) : ""}`);
  };

  const uploadFile = async (dryRun: boolean) => {
    if (!file || !projectId) return;
    setBusy(true);
    try {
      const fd = new FormData();
      fd.append("file", file);
      if (groupId) fd.append("group_id", groupId);
      if (dryRun) fd.append("dry_run", "true");
      const res = await fetch(`/api/backend/api/admin/projects/${projectId}/contacts/import`, { method: "POST", body: fd });
      if (!res.ok) return void (await fail(res));
      done(((await res.json()) as { data: ImportResult }).data);
    } finally {
      setBusy(false);
    }
  };
  const importApple = async (dryRun: boolean) => {
    if (!projectId) return;
    setBusy(true);
    try {
      const res = await fetch(`/api/backend/api/admin/projects/${projectId}/contacts/import/carddav`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...apple, group_id: groupId || undefined, dry_run: dryRun }),
      });
      if (!res.ok) return void (await fail(res));
      done(((await res.json()) as { data: ImportResult }).data);
    } finally {
      setBusy(false);
    }
  };
  const importGoogle = async () => {
    if (!projectId) return;
    const res = await api.GET("/api/admin/projects/{project_id}/contacts/import/google/url", { params: { path: { project_id: projectId }, query: { group_id: groupId || undefined, return_to: window.location.pathname } } });
    if (res.error || !res.response.ok) return void toast.error(t("googleNotConfigured"));
    window.location.href = (res.data as { data: { url: string } }).data.url;
  };

  return (
    <Dialog open={open} onOpenChange={(o) => { onOpenChange(o); if (!o) { setResult(null); setFile(null); } }}>
      <DialogContent wide>
        <DialogHeader><DialogTitle>{t("importTitle")}</DialogTitle><DialogDescription>{t("importFileHint")}</DialogDescription></DialogHeader>
        <Field label={t("addToGroup")}>
          <Select value={groupId || "__none"} onValueChange={(v) => setGroupId(v === "__none" ? "" : v)}>
            <SelectTrigger><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem value="__none">{t("noGroup")}</SelectItem>
              {(groups.data ?? []).map((g) => <SelectItem key={g.id} value={g.id!}>{g.name}</SelectItem>)}
            </SelectContent>
          </Select>
        </Field>
        <Tabs defaultValue="file" onValueChange={() => setResult(null)}>
          <TabsList>
            <TabsTrigger value="file"><FileSpreadsheet className="mr-1 h-4 w-4" />{t("importFile")}</TabsTrigger>
            <TabsTrigger value="google">{t("importGoogle")}</TabsTrigger>
            <TabsTrigger value="apple"><Apple className="mr-1 h-4 w-4" />{t("importApple")}</TabsTrigger>
          </TabsList>
          <TabsContent value="file" className="space-y-3">
            <input ref={fileRef} type="file" accept=".xlsx,.xlsm,.csv,.tsv,.txt,.vcf,.vcard,.docx" className="hidden" onChange={(e) => { setFile(e.target.files?.[0] ?? null); setResult(null); }} data-testid="import-file" />
            <div className="flex flex-wrap items-center gap-2">
              <Button variant="outline" onClick={() => fileRef.current?.click()}><Upload /> {t("chooseFile")}</Button>
              <span className="text-sm text-muted-foreground">{file ? `${file.name} · ${Math.round(file.size / 1024)} KB` : "—"}</span>
            </div>
            <div className="flex gap-2">
              <Button variant="secondary" disabled={!file || busy} onClick={() => uploadFile(true)}>{t("previewFirst")}</Button>
              <Button disabled={!file || busy} onClick={() => uploadFile(false)} data-testid="import-now">{t("importNow")}</Button>
            </div>
          </TabsContent>
          <TabsContent value="google" className="space-y-3">
            <p className="text-sm text-muted-foreground">{t("importGoogleHint")}</p>
            <Button onClick={importGoogle} disabled={busy}>{t("importGoogle")}</Button>
          </TabsContent>
          <TabsContent value="apple" className="space-y-3">
            <p className="text-sm text-muted-foreground">{t("importAppleHint")}</p>
            <div className="grid gap-3 sm:grid-cols-2">
              <Field label={t("appleId")} htmlFor="apple-id"><Input id="apple-id" value={apple.username} onChange={(e) => setApple({ ...apple, username: e.target.value })} placeholder="you@icloud.com" /></Field>
              <Field label={t("appPassword")} htmlFor="apple-pw"><Input id="apple-pw" type="password" value={apple.password} onChange={(e) => setApple({ ...apple, password: e.target.value })} placeholder="abcd-efgh-ijkl-mnop" /></Field>
              <div className="sm:col-span-2"><Field label={t("serverUrl")} htmlFor="apple-url"><Input id="apple-url" value={apple.server_url} onChange={(e) => setApple({ ...apple, server_url: e.target.value })} /></Field></div>
            </div>
            <div className="flex gap-2">
              <Button variant="secondary" disabled={!apple.username || !apple.password || busy} onClick={() => importApple(true)}>{t("previewFirst")}</Button>
              <Button disabled={!apple.username || !apple.password || busy} onClick={() => importApple(false)}>{t("importNow")}</Button>
            </div>
          </TabsContent>
        </Tabs>
        {result ? (
          <div className="space-y-2 rounded-lg border bg-muted/40 p-3 text-sm" data-testid="import-result">
            <p className="font-medium">{result.dry_run ? `${t("preview")} · ` : ""}{t("importResult", { total: result.total ?? 0, created: result.created ?? 0, updated: result.updated ?? 0, unchanged: result.unchanged ?? 0, skipped: result.skipped ?? 0 })}</p>
            {result.preview?.length ? (
              <ul className="max-h-40 overflow-auto font-mono text-xs text-muted-foreground">
                {result.preview.map((p, i) => { const r = p as Record<string, string>; return <li key={i}>{[r.Name, r.Phone, r.Email, r.WhatsApp, r.TelegramChatID, r.SlackID].filter(Boolean).join(" · ")}</li>; })}
              </ul>
            ) : null}
            {result.errors?.length ? <p className="text-xs text-destructive">{t("importErrors")}: {result.errors.slice(0, 5).map((e) => `#${e.row} ${e.reason}`).join("; ")}</p> : null}
          </div>
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
