"use client";
import { Suspense, useState } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil, Plus, RotateCcw, Trash2, Upload } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import type { AdminDevice, Contact } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { ConfirmButton } from "@/components/common/confirm-button";
import { Field } from "@/components/common/field";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { formatDate } from "@/lib/utils";
import { ImportDialog } from "@/components/contacts/import-dialog";

type ContactForm = { id?: string; name: string; external_id: string; phone: string; email: string; whatsapp: string; telegram_chat_id: string; slack_id: string };
const emptyContact: ContactForm = { name: "", external_id: "", phone: "", email: "", whatsapp: "", telegram_chat_id: "", slack_id: "" };

export default function ContactsPage() {
  return <Suspense><Contacts /></Suspense>;
}

function Contacts() {
  const t = useTranslations("contacts");
  const tc = useTranslations("common");
  const { projectId, can } = useProject();
  const qc = useQueryClient();
  const editable = can("developer");
  const [form, setForm] = useState<ContactForm | null>(null);
  const [importing, setImporting] = useState(false);
  const [search, setSearch] = useState("");
  const [applied, setApplied] = useState("");
  const [offset, setOffset] = useState(0);
  const limit = 50;
  const contacts = useQuery({
    queryKey: ["contacts", projectId, { applied, offset }],
    enabled: !!projectId,
    queryFn: async () => {
      const res = await api.GET("/api/admin/projects/{project_id}/contacts", { params: { path: { project_id: projectId! }, query: { search: applied || undefined, limit, offset } } });
      return { rows: unwrap<Contact[]>(res), total: Number((res.data as { meta?: { total?: number } })?.meta?.total ?? 0) };
    },
  });
  const save = useMutation({
    mutationFn: async (f: ContactForm) => {
      const body = { name: f.name, external_id: f.external_id, phone: f.phone, email: f.email, whatsapp: f.whatsapp, telegram_chat_id: f.telegram_chat_id, slack_id: f.slack_id };
      if (f.id) return unwrap<Contact>(await api.PUT("/api/admin/projects/{project_id}/contacts/{contact_id}", { params: { path: { project_id: projectId!, contact_id: f.id } }, body }));
      return unwrap<Contact>(await api.POST("/api/admin/projects/{project_id}/contacts", { params: { path: { project_id: projectId! } }, body }));
    },
    onSuccess: () => { toast.success(tc("save")); void qc.invalidateQueries({ queryKey: ["contacts", projectId] }); setForm(null); },
    onError: (e: Error) => toast.error(e instanceof ApiError ? `${e.message} ${e.details ? JSON.stringify(e.details) : ""}` : e.message),
  });
  const del = useMutation({
    mutationFn: async (id: string) => { await api.DELETE("/api/admin/projects/{project_id}/contacts/{contact_id}", { params: { path: { project_id: projectId!, contact_id: id } } }); },
    onSuccess: () => { void qc.invalidateQueries({ queryKey: ["contacts", projectId] }); void qc.invalidateQueries({ queryKey: ["contacts-trash", projectId] }); },
  });
  const trash = useQuery({
    queryKey: ["contacts-trash", projectId, offset],
    enabled: !!projectId && editable,
    queryFn: async () => {
      const res = await api.GET("/api/admin/projects/{project_id}/contacts/trash", { params: { path: { project_id: projectId! }, query: { limit, offset } } });
      return { rows: unwrap<Contact[]>(res), total: Number((res.data as { meta?: { total?: number } })?.meta?.total ?? 0) };
    },
  });
  const invalidateBoth = () => { void qc.invalidateQueries({ queryKey: ["contacts", projectId] }); void qc.invalidateQueries({ queryKey: ["contacts-trash", projectId] }); };
  const restore = useMutation({
    mutationFn: async (id: string) => unwrap<Contact>(await api.POST("/api/admin/projects/{project_id}/contacts/{contact_id}/restore", { params: { path: { project_id: projectId!, contact_id: id } } })),
    onSuccess: () => { toast.success(t("restoredToast")); invalidateBoth(); },
    onError: (e: Error) => toast.error(e instanceof ApiError ? `${e.message} ${e.details ? JSON.stringify(e.details) : ""}` : e.message),
  });
  const purge = useMutation({
    mutationFn: async (id: string) => { await api.DELETE("/api/admin/projects/{project_id}/contacts/{contact_id}/purge", { params: { path: { project_id: projectId!, contact_id: id } } }); },
    onSuccess: () => void qc.invalidateQueries({ queryKey: ["contacts-trash", projectId] }),
  });
  const purgeAll = useMutation({
    mutationFn: async () => { const r = await api.DELETE("/api/admin/projects/{project_id}/contacts/trash", { params: { path: { project_id: projectId! } } }); return Number((r.data as { data?: { purged?: number } })?.data?.purged ?? 0); },
    onSuccess: (n) => { toast.success(t("purgedToast", { n })); void qc.invalidateQueries({ queryKey: ["contacts-trash", projectId] }); },
  });
  const devices = useQuery({
    queryKey: ["devices", projectId],
    enabled: !!projectId,
    queryFn: async () => unwrap<AdminDevice[]>(await api.GET("/api/admin/projects/{project_id}/devices", { params: { path: { project_id: projectId! }, query: { limit: 200 } } })),
  });

  return (
    <>
      <PageHeader title={t("title")} actions={editable ? <><Button size="sm" variant="outline" onClick={() => setImporting(true)} data-testid="import-contacts"><Upload /> {t("import")}</Button><Button size="sm" onClick={() => setForm(emptyContact)} data-testid="new-contact"><Plus /> {t("new")}</Button></> : null} />
      {editable ? <ImportDialog open={importing} onOpenChange={setImporting} /> : null}
      <Tabs defaultValue="contacts">
        <TabsList><TabsTrigger value="contacts">{t("contacts")} {contacts.data ? `(${contacts.data.total})` : ""}</TabsTrigger><TabsTrigger value="devices">{t("devices")} {devices.data ? `(${devices.data.length})` : ""}</TabsTrigger>{editable ? <TabsTrigger value="trash">{t("trash")} {trash.data?.total ? `(${trash.data.total})` : ""}</TabsTrigger> : null}</TabsList>
        <TabsContent value="contacts">
          <form className="mb-3 max-w-sm" onSubmit={(e) => { e.preventDefault(); setOffset(0); setApplied(search.trim()); }}>
            <Input placeholder={`${tc("search")}…`} value={search} onChange={(e) => setSearch(e.target.value)} />
          </form>
          {contacts.isLoading ? <Skeleton className="h-48" /> : !contacts.data?.rows.length ? <EmptyState /> : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader><TableRow><TableHead>{t("name")}</TableHead><TableHead>{t("externalId")}</TableHead><TableHead>{t("phone")}</TableHead><TableHead>{t("email")}</TableHead><TableHead>{t("telegram")} / {t("slack")}</TableHead><TableHead>{t("tags")}</TableHead><TableHead>{tc("created")}</TableHead>{editable ? <TableHead /> : null}</TableRow></TableHeader>
                <TableBody>
                  {contacts.data.rows.map((c) => (
                    <TableRow key={c.id}>
                      <TableCell className="font-medium">{c.name || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.external_id || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.phone || "—"}{c.whatsapp && c.whatsapp !== c.phone ? <span className="ml-1 text-chart-whatsapp">WA {c.whatsapp}</span> : null}</TableCell>
                      <TableCell className="text-xs">{c.email || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.telegram_chat_id || c.slack_id || "—"}</TableCell>
                      <TableCell className="space-x-1">{(c.tags ?? []).map((tag) => <Badge key={tag} variant="outline">{tag}</Badge>)}</TableCell>
                      <TableCell className="tabular text-xs text-muted-foreground">{formatDate(c.created_at)}</TableCell>
                      {editable ? (
                        <TableCell className="space-x-1 text-right whitespace-nowrap">
                          <Button size="sm" variant="ghost" onClick={() => setForm({ id: c.id, name: c.name ?? "", external_id: c.external_id ?? "", phone: c.phone ?? "", email: c.email ?? "", whatsapp: c.whatsapp ?? "", telegram_chat_id: c.telegram_chat_id ?? "", slack_id: c.slack_id ?? "" })}><Pencil /></Button>
                          <ConfirmButton size="sm" variant="ghost" title={tc("confirmDelete", { name: c.name || c.phone || c.email || "" })} onConfirm={() => del.mutate(c.id!)}>{tc("delete")}</ConfirmButton>
                        </TableCell>
                      ) : null}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              <div className="flex items-center justify-between border-t p-3 text-xs text-muted-foreground">
                <span>{offset + 1}–{offset + contacts.data.rows.length} / {contacts.data.total}</span>
                <div className="space-x-2">
                  <Button size="sm" variant="outline" disabled={offset === 0} onClick={() => setOffset(Math.max(0, offset - limit))}>‹</Button>
                  <Button size="sm" variant="outline" disabled={offset + limit >= contacts.data.total} onClick={() => setOffset(offset + limit)}>›</Button>
                </div>
              </div>
            </div>
          )}
        </TabsContent>
        <TabsContent value="devices">
          {devices.isLoading ? <Skeleton className="h-48" /> : !devices.data?.length ? <EmptyState /> : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader><TableRow><TableHead>{t("platform")}</TableHead><TableHead>{t("token")}</TableHead><TableHead>Contact</TableHead><TableHead>App</TableHead><TableHead>{t("active")}</TableHead><TableHead>{t("lastSeen")}</TableHead></TableRow></TableHeader>
                <TableBody>
                  {devices.data.map((d) => (
                    <TableRow key={d.ID}>
                      <TableCell><Badge variant="secondary">{d.Platform}</Badge></TableCell>
                      <TableCell className="font-mono text-xs">{d.fcm_token}</TableCell>
                      <TableCell className="font-mono text-xs">{d.ContactID ? d.ContactID.slice(0, 8) : "—"}</TableCell>
                      <TableCell className="text-xs">{d.AppVersion || "—"}</TableCell>
                      <TableCell>{d.IsActive ? <Badge variant="success">{tc("yes")}</Badge> : <Badge variant="outline">{tc("no")}</Badge>}</TableCell>
                      <TableCell className="tabular text-xs text-muted-foreground">{formatDate(d.LastSeenAt)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </TabsContent>
        {editable ? (
          <TabsContent value="trash">
            <div className="mb-3 flex items-center justify-between gap-2">
              <p className="text-xs text-muted-foreground">{t("softDeleteHint")}</p>
              {trash.data?.rows.length ? <ConfirmButton size="sm" variant="destructive" title={t("purgeAll")} onConfirm={() => purgeAll.mutate()}><Trash2 /> {t("purgeAll")}</ConfirmButton> : null}
            </div>
            {trash.isLoading ? <Skeleton className="h-48" /> : !trash.data?.rows.length ? <EmptyState hint={t("emptyTrash")} /> : (
              <div className="rounded-lg border">
                <Table>
                  <TableHeader><TableRow><TableHead>{t("name")}</TableHead><TableHead>{t("phone")}</TableHead><TableHead>{t("email")}</TableHead><TableHead>{t("deletedAt")}</TableHead><TableHead /></TableRow></TableHeader>
                  <TableBody>
                    {trash.data.rows.map((c) => (
                      <TableRow key={c.id} className="text-muted-foreground">
                        <TableCell className="font-medium">{c.name || "—"}</TableCell>
                        <TableCell className="font-mono text-xs">{c.phone || "—"}</TableCell>
                        <TableCell className="text-xs">{c.email || "—"}</TableCell>
                        <TableCell className="tabular text-xs">{formatDate(c.deleted_at)}</TableCell>
                        <TableCell className="space-x-1 text-right whitespace-nowrap">
                          <Button size="sm" variant="ghost" onClick={() => restore.mutate(c.id!)}><RotateCcw /> {t("restore")}</Button>
                          <ConfirmButton size="sm" variant="ghost" title={t("purge")} description={tc("confirmDelete", { name: c.name || c.phone || c.email || "" })} onConfirm={() => purge.mutate(c.id!)}>{t("purge")}</ConfirmButton>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </TabsContent>
        ) : null}
      </Tabs>

      <Dialog open={!!form} onOpenChange={(o) => !o && setForm(null)}>
        <DialogContent>
          <DialogHeader><DialogTitle>{form?.id ? t("edit") : t("new")}</DialogTitle></DialogHeader>
          {form ? (
            <div className="grid gap-3 sm:grid-cols-2">
              <Field label={t("name")} htmlFor="ct-name"><Input id="ct-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} /></Field>
              <Field label={t("externalId")} htmlFor="ct-ext"><Input id="ct-ext" value={form.external_id} onChange={(e) => setForm({ ...form, external_id: e.target.value })} /></Field>
              <Field label={t("phone")} htmlFor="ct-phone"><Input id="ct-phone" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} placeholder="+99365123456" /></Field>
              <Field label={t("whatsapp")} htmlFor="ct-wa"><Input id="ct-wa" value={form.whatsapp} onChange={(e) => setForm({ ...form, whatsapp: e.target.value })} placeholder="= phone" /></Field>
              <Field label={t("email")} htmlFor="ct-mail"><Input id="ct-mail" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} /></Field>
              <Field label={t("telegram")} htmlFor="ct-tg"><Input id="ct-tg" value={form.telegram_chat_id} onChange={(e) => setForm({ ...form, telegram_chat_id: e.target.value })} /></Field>
              <Field label={t("slack")} htmlFor="ct-slack"><Input id="ct-slack" value={form.slack_id} onChange={(e) => setForm({ ...form, slack_id: e.target.value })} placeholder="U0123ABCD" /></Field>
            </div>
          ) : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setForm(null)}>{tc("cancel")}</Button>
            <Button onClick={() => form && save.mutate(form)} disabled={save.isPending}>{tc("save")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
