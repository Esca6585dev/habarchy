"use client";
import { useState } from "react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil, Plus, Send, UserPlus, Users } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap, ApiError } from "@/lib/api/client";
import { qk } from "@/lib/hooks";
import type { Contact, Group, GroupMembersResult } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { ConfirmButton } from "@/components/common/confirm-button";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { parseInlineContact } from "@/lib/contacts";

type GroupForm = { id?: string; name: string; description: string };

export default function GroupsPage() {
  const t = useTranslations("groups");
  const tc = useTranslations("common");
  const tct = useTranslations("contacts");
  const { projectId, can } = useProject();
  const qc = useQueryClient();
  const editable = can("developer");
  const [form, setForm] = useState<GroupForm | null>(null);
  const [open, setOpen] = useState<Group | null>(null);
  const [adding, setAdding] = useState(false);
  const [search, setSearch] = useState("");
  const [picked, setPicked] = useState<Set<string>>(new Set());
  const [inline, setInline] = useState("");

  const groups = useQuery({
    queryKey: qk.groups(projectId ?? ""),
    enabled: !!projectId,
    queryFn: async () => unwrap<Group[]>(await api.GET("/api/admin/projects/{project_id}/groups", { params: { path: { project_id: projectId! } } })),
  });
  const members = useQuery({
    queryKey: qk.groupMembers(projectId ?? "", open?.id ?? ""),
    enabled: !!projectId && !!open,
    queryFn: async () => unwrap<Contact[]>(await api.GET("/api/admin/projects/{project_id}/groups/{group_id}/members", { params: { path: { project_id: projectId!, group_id: open!.id! }, query: { limit: 500 } } })),
  });
  const candidates = useQuery({
    queryKey: qk.contacts(projectId ?? "", { search, forGroup: open?.id }),
    enabled: !!projectId && adding,
    queryFn: async () => unwrap<Contact[]>(await api.GET("/api/admin/projects/{project_id}/contacts", { params: { path: { project_id: projectId! }, query: { search: search || undefined, limit: 50 } } })),
  });
  const invalidate = () => {
    void qc.invalidateQueries({ queryKey: qk.groups(projectId ?? "") });
    if (open) void qc.invalidateQueries({ queryKey: qk.groupMembers(projectId ?? "", open.id!) });
  };
  const onError = (e: Error) => toast.error(e instanceof ApiError ? `${e.message} ${e.details ? JSON.stringify(e.details) : ""}` : e.message);

  const save = useMutation({
    mutationFn: async (f: GroupForm) => {
      const body = { name: f.name, description: f.description };
      if (f.id) return unwrap<Group>(await api.PUT("/api/admin/projects/{project_id}/groups/{group_id}", { params: { path: { project_id: projectId!, group_id: f.id } }, body }));
      return unwrap<Group>(await api.POST("/api/admin/projects/{project_id}/groups", { params: { path: { project_id: projectId! } }, body }));
    },
    onSuccess: (g) => { toast.success(tc("save")); invalidate(); setForm(null); if (open && g.id === open.id) setOpen(g); },
    onError,
  });
  const del = useMutation({
    mutationFn: async (id: string) => { await api.DELETE("/api/admin/projects/{project_id}/groups/{group_id}", { params: { path: { project_id: projectId!, group_id: id } } }); },
    onSuccess: () => { invalidate(); setOpen(null); },
    onError,
  });
  const add = useMutation({
    mutationFn: async () => {
      const contacts = inline.split("\n").map(parseInlineContact).filter((c): c is Record<string, string> => !!c);
      return unwrap<GroupMembersResult>(await api.POST("/api/admin/projects/{project_id}/groups/{group_id}/members", {
        params: { path: { project_id: projectId!, group_id: open!.id! } },
        body: { contact_ids: Array.from(picked), contacts },
      }));
    },
    onSuccess: (r) => { toast.success(t("added", { n: r.added ?? 0, created: r.created_contacts ?? 0 })); invalidate(); setAdding(false); setPicked(new Set()); setInline(""); },
    onError,
  });
  const remove = useMutation({
    mutationFn: async (contactId: string) => { await api.DELETE("/api/admin/projects/{project_id}/groups/{group_id}/members/{contact_id}", { params: { path: { project_id: projectId!, group_id: open!.id!, contact_id: contactId } } }); },
    onSuccess: invalidate,
    onError,
  });

  const label = (c: Contact) => c.name || c.external_id || c.phone || c.email || c.telegram_chat_id || c.slack_id || c.id?.slice(0, 8);

  return (
    <>
      <PageHeader title={t("title")} actions={editable ? <Button size="sm" onClick={() => setForm({ name: "", description: "" })} data-testid="new-group"><Plus /> {t("new")}</Button> : null} />
      {groups.isLoading ? <Skeleton className="h-48" /> : !groups.data?.length ? <EmptyState hint={t("empty")} /> : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {groups.data.map((g) => (
            <Card key={g.id} className="flex flex-col">
              <CardHeader className="pb-2">
                <CardTitle className="flex items-center gap-2 text-base"><Users className="h-4 w-4 text-muted-foreground" /> {g.name}</CardTitle>
                {g.description ? <p className="text-sm text-muted-foreground">{g.description}</p> : null}
              </CardHeader>
              <CardContent className="mt-auto flex items-center justify-between gap-2">
                <Badge variant="secondary">{t("memberCount", { n: g.member_count ?? 0 })}</Badge>
                <div className="flex gap-1">
                  <Button size="sm" variant="outline" onClick={() => setOpen(g)}>{t("members")}</Button>
                  {editable ? <Button size="sm" asChild><Link href={`/compose?group=${g.id}`}><Send /> {t("sendTo")}</Link></Button> : null}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={!!form} onOpenChange={(o) => !o && setForm(null)}>
        <DialogContent>
          <DialogHeader><DialogTitle>{form?.id ? tc("edit") : t("new")}</DialogTitle></DialogHeader>
          {form ? (
            <div className="space-y-3">
              <Field label={t("name")} htmlFor="g-name"><Input id="g-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="Işdeşler" /></Field>
              <Field label={t("description")} htmlFor="g-desc"><Input id="g-desc" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></Field>
            </div>
          ) : null}
          <DialogFooter>
            <Button variant="outline" onClick={() => setForm(null)}>{tc("cancel")}</Button>
            <Button onClick={() => form && save.mutate(form)} disabled={save.isPending || !form?.name.trim()}>{tc("save")}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!open} onOpenChange={(o) => { if (!o) { setOpen(null); setAdding(false); } }}>
        <DialogContent wide>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">{open?.name}
              {editable && open ? <Button size="sm" variant="ghost" onClick={() => setForm({ id: open.id, name: open.name ?? "", description: open.description ?? "" })}><Pencil /></Button> : null}
            </DialogTitle>
            <DialogDescription>{open?.description || t("members")}</DialogDescription>
          </DialogHeader>
          {editable ? (
            <div className="flex flex-wrap gap-2">
              <Button size="sm" variant={adding ? "secondary" : "default"} onClick={() => setAdding((v) => !v)}><UserPlus /> {t("addMembers")}</Button>
              {open ? <Button size="sm" variant="outline" asChild><Link href={`/compose?group=${open.id}`}><Send /> {t("sendTo")}</Link></Button> : null}
              {open ? <ConfirmButton size="sm" variant="destructive" title={tc("confirmDelete", { name: open.name ?? "" })} description={t("deleteHint")} onConfirm={() => del.mutate(open.id!)}>{tc("delete")}</ConfirmButton> : null}
            </div>
          ) : null}
          {adding ? (
            <div className="grid gap-3 rounded-lg border p-3 sm:grid-cols-2">
              <div className="space-y-2">
                <p className="text-sm font-medium">{t("pickContacts")}</p>
                <Input placeholder={`${tc("search")}…`} value={search} onChange={(e) => setSearch(e.target.value)} />
                <div className="max-h-56 space-y-1 overflow-auto text-sm">
                  {(candidates.data ?? []).map((c) => {
                    const checked = picked.has(c.id!);
                    return (
                      <label key={c.id} className="flex cursor-pointer items-center gap-2 rounded px-1 py-0.5 hover:bg-accent">
                        <input type="checkbox" checked={checked} onChange={() => { const n = new Set(picked); if (checked) n.delete(c.id!); else n.add(c.id!); setPicked(n); }} />
                        <span className="truncate">{label(c)}</span>
                        <span className="ml-auto truncate font-mono text-xs text-muted-foreground">{c.phone || c.email}</span>
                      </label>
                    );
                  })}
                </div>
              </div>
              <div className="space-y-2">
                <p className="text-sm font-medium">{t("newContacts")}</p>
                <Textarea rows={8} className="font-mono text-xs" value={inline} onChange={(e) => setInline(e.target.value)} placeholder={"Aman Amanow, +99365123456\nMaral, maral@example.tm"} />
                <p className="text-xs text-muted-foreground">{t("addHint")}</p>
                <Button size="sm" onClick={() => add.mutate()} disabled={add.isPending || (picked.size === 0 && !inline.trim())} data-testid="add-members">{t("addMembers")}</Button>
              </div>
            </div>
          ) : null}
          {members.isLoading ? <Skeleton className="h-32" /> : !members.data?.length ? <EmptyState /> : (
            <div className="max-h-[50vh] overflow-auto rounded-lg border">
              <Table>
                <TableHeader><TableRow><TableHead>{tct("name")}</TableHead><TableHead>{tct("phone")}</TableHead><TableHead>{tct("email")}</TableHead><TableHead>{tct("telegram")} / {tct("slack")}</TableHead><TableHead /></TableRow></TableHeader>
                <TableBody>
                  {members.data.map((c) => (
                    <TableRow key={c.id}>
                      <TableCell className="font-medium">{c.name || c.external_id || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.phone || "—"}{c.whatsapp && c.whatsapp !== c.phone ? <span className="ml-1 text-chart-whatsapp">WA {c.whatsapp}</span> : null}</TableCell>
                      <TableCell className="text-xs">{c.email || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.telegram_chat_id || c.slack_id || "—"}</TableCell>
                      <TableCell className="text-right">{editable ? <Button size="sm" variant="ghost" onClick={() => remove.mutate(c.id!)}>{t("remove")}</Button> : null}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  );
}
