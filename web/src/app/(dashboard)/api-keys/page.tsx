"use client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { KeyRound, Plus } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { APIKey } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { ConfirmButton } from "@/components/common/confirm-button";
import { CopyButton } from "@/components/common/copy-button";
import { Field } from "@/components/common/field";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

const allScopes = ["messages:send", "messages:read", "otp", "templates", "contacts", "devices", "usage"] as const;

export default function ApiKeysPage() {
  const t = useTranslations("apiKeys");
  const tc = useTranslations("common");
  const { projectId, can } = useProject();
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [created, setCreated] = useState<APIKey | null>(null);
  const [form, setForm] = useState({ name: "", live: false, scopes: [...allScopes] as string[], ips: "", requireSignature: false });

  const q = useQuery({
    queryKey: qk.apiKeys(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<APIKey[]>(await api.GET("/api/admin/projects/{project_id}/api-keys", { params: { path: { project_id: projectId! } } })),
  });
  const create = useMutation({
    mutationFn: async () =>
      unwrap<APIKey>(await api.POST("/api/admin/projects/{project_id}/api-keys", {
        params: { path: { project_id: projectId! } },
        body: { name: form.name, live: form.live, scopes: form.scopes as APIKey["scopes"] as never, ip_allowlist: form.ips.split("\n").map((s) => s.trim()).filter(Boolean), require_signature: form.requireSignature },
      })),
    onSuccess: (k) => { setCreated(k); setOpen(false); void qc.invalidateQueries({ queryKey: qk.apiKeys(projectId ?? "") }); },
    onError: (e: ApiError) => toast.error(e.message + (e.details ? ` ${JSON.stringify(e.details)}` : "")),
  });
  const revoke = useMutation({
    mutationFn: async (id: string) => { await api.DELETE("/api/admin/projects/{project_id}/api-keys/{key_id}", { params: { path: { project_id: projectId!, key_id: id } } }); },
    onSuccess: () => { toast.success(t("revoked")); void qc.invalidateQueries({ queryKey: qk.apiKeys(projectId ?? "") }); },
  });

  return (
    <>
      <PageHeader title={t("title")} actions={can("admin") ? <Button size="sm" onClick={() => setOpen(true)} data-testid="new-api-key"><Plus /> {t("new")}</Button> : null} />
      {created ? (
        <Alert variant="success" className="mb-4" data-testid="created-key">
          <KeyRound className="h-4 w-4" />
          <AlertTitle>{t("createdTitle")}</AlertTitle>
          <AlertDescription>
            <p className="mb-2">{t("createdHint")}</p>
            <div className="flex items-center gap-2"><code className="break-all rounded bg-muted px-2 py-1 font-mono text-xs">{created.key}</code><CopyButton value={created.key ?? ""} /></div>
          </AlertDescription>
        </Alert>
      ) : null}
      {q.isLoading ? <Skeleton className="h-48" /> : !q.data?.length ? <EmptyState /> : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader><TableRow><TableHead>{tc("name")}</TableHead><TableHead>{t("hint")}</TableHead><TableHead>{t("scopes")}</TableHead><TableHead>{t("lastUsed")}</TableHead><TableHead>{tc("status")}</TableHead><TableHead className="text-right">{tc("actions")}</TableHead></TableRow></TableHeader>
            <TableBody>
              {q.data.map((k) => (
                <TableRow key={k.id}>
                  <TableCell className="font-medium">{k.name}</TableCell>
                  <TableCell className="font-mono text-xs">{k.prefix}…{k.hint} {k.require_signature ? <Badge variant="outline" className="ml-1">HMAC</Badge> : null}</TableCell>
                  <TableCell className="max-w-[16rem] truncate text-xs text-muted-foreground">{(k.scopes ?? []).join(", ")}</TableCell>
                  <TableCell className="tabular text-xs text-muted-foreground">{k.last_used_at ? formatDate(k.last_used_at) : t("never")}</TableCell>
                  <TableCell>{k.revoked_at ? <Badge variant="outline">{t("revoked")}</Badge> : <Badge variant={k.prefix === "hb_live_" ? "success" : "info"}>{k.prefix === "hb_live_" ? "live" : "test"}</Badge>}</TableCell>
                  <TableCell className="text-right">{!k.revoked_at && can("admin") ? <ConfirmButton size="sm" variant="ghost" title={tc("confirmDelete", { name: k.name ?? "" })} confirmLabel={t("revoke")} onConfirm={() => revoke.mutate(k.id!)}>{t("revoke")}</ConfirmButton> : null}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader><DialogTitle>{t("new")}</DialogTitle></DialogHeader>
          <div className="grid gap-3">
            <Field label={tc("name")} htmlFor="k-name"><Input id="k-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></Field>
            <div className="flex items-center gap-2"><Switch id="k-live" checked={form.live} onCheckedChange={(v) => setForm({ ...form, live: v })} /><Label htmlFor="k-live">{form.live ? t("live") : t("test")}</Label></div>
            <Field label={t("scopes")}>
              <div className="grid grid-cols-2 gap-2">
                {allScopes.map((s) => (
                  <label key={s} className="flex items-center gap-2 text-sm">
                    <input type="checkbox" className="accent-primary" checked={form.scopes.includes(s)} onChange={(e) => setForm({ ...form, scopes: e.target.checked ? [...form.scopes, s] : form.scopes.filter((x) => x !== s) })} />
                    <span className="font-mono text-xs">{s}</span>
                  </label>
                ))}
              </div>
            </Field>
            <Field label={t("ipAllowlist")} htmlFor="k-ips"><Textarea id="k-ips" rows={2} value={form.ips} onChange={(e) => setForm({ ...form, ips: e.target.value })} placeholder="10.0.0.0/8" /></Field>
            <div className="flex items-center gap-2"><Switch id="k-sig" checked={form.requireSignature} onCheckedChange={(v) => setForm({ ...form, requireSignature: v })} /><Label htmlFor="k-sig">{t("requireSignature")}</Label></div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>{tc("cancel")}</Button>
            <Button onClick={() => create.mutate()} disabled={create.isPending || !form.name}>{tc("create")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
