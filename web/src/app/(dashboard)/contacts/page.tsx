"use client";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useQuery } from "@tanstack/react-query";
import { api, unwrap } from "@/lib/api/client";
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
import { formatDate } from "@/lib/utils";

export default function ContactsPage() {
  const t = useTranslations("contacts");
  const tc = useTranslations("common");
  const { projectId } = useProject();
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
  const devices = useQuery({
    queryKey: ["devices", projectId],
    enabled: !!projectId,
    queryFn: async () => unwrap<AdminDevice[]>(await api.GET("/api/admin/projects/{project_id}/devices", { params: { path: { project_id: projectId! }, query: { limit: 200 } } })),
  });

  return (
    <>
      <PageHeader title={t("title")} />
      <Tabs defaultValue="contacts">
        <TabsList><TabsTrigger value="contacts">{t("contacts")} {contacts.data ? `(${contacts.data.total})` : ""}</TabsTrigger><TabsTrigger value="devices">{t("devices")} {devices.data ? `(${devices.data.length})` : ""}</TabsTrigger></TabsList>
        <TabsContent value="contacts">
          <form className="mb-3 max-w-sm" onSubmit={(e) => { e.preventDefault(); setOffset(0); setApplied(search.trim()); }}>
            <Input placeholder={`${tc("search")}…`} value={search} onChange={(e) => setSearch(e.target.value)} />
          </form>
          {contacts.isLoading ? <Skeleton className="h-48" /> : !contacts.data?.rows.length ? <EmptyState /> : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader><TableRow><TableHead>{t("externalId")}</TableHead><TableHead>{t("phone")}</TableHead><TableHead>{t("email")}</TableHead><TableHead>{t("telegram")}</TableHead><TableHead>{t("tags")}</TableHead><TableHead>{tc("created")}</TableHead></TableRow></TableHeader>
                <TableBody>
                  {contacts.data.rows.map((c) => (
                    <TableRow key={c.id}>
                      <TableCell className="font-mono text-xs">{c.external_id || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.phone || "—"}</TableCell>
                      <TableCell className="text-xs">{c.email || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{c.telegram_chat_id || "—"}</TableCell>
                      <TableCell className="space-x-1">{(c.tags ?? []).map((tag) => <Badge key={tag} variant="outline">{tag}</Badge>)}</TableCell>
                      <TableCell className="tabular text-xs text-muted-foreground">{formatDate(c.created_at)}</TableCell>
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
      </Tabs>
    </>
  );
}
