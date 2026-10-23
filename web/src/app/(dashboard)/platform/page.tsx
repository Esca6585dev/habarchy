"use client";
import { useQuery } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate } from "@/lib/utils";

// The platform endpoints are owner-only and not in the typed OpenAPI client,
// so they are fetched through the backend proxy directly.
async function platformGet<T>(path: string): Promise<{ data: T; meta?: Record<string, unknown> }> {
  const res = await fetch(`/api/backend/api/admin/platform/${path}`);
  if (!res.ok) throw new Error(String(res.status));
  return (await res.json()) as { data: T; meta?: Record<string, unknown> };
}

type Tenant = { id: string; name: string; slug: string; status: string; created_at: string };
type PUser = { id: string; email: string; phone?: string; platform_admin: boolean; full_name: string; created_at: string };
type MsgRow = {
  message: { id: string; status: string; channel: string; to: string; created_at: string };
  tenant: { id: string; name: string; slug: string };
};
type Overview = { projects: number; users: number; messages_24h: number; pending: number };

function statusVariant(s: string): "default" | "secondary" | "destructive" | "outline" {
  if (s === "delivered" || s === "sent") return "default";
  if (s === "failed" || s === "cancelled") return "destructive";
  if (s === "queued" || s === "processing") return "secondary";
  return "outline";
}

export default function PlatformPage() {
  const t = useTranslations("platform");

  const overview = useQuery({ queryKey: ["platform", "overview"], queryFn: () => platformGet<Overview>("overview") });
  const messages = useQuery({ queryKey: ["platform", "messages"], queryFn: () => platformGet<MsgRow[]>("messages?limit=100") });
  const tenants = useQuery({ queryKey: ["platform", "tenants"], queryFn: () => platformGet<Tenant[]>("tenants") });
  const users = useQuery({ queryKey: ["platform", "users"], queryFn: () => platformGet<PUser[]>("users") });

  const kpis: { label: string; value: number | undefined }[] = [
    { label: t("tenants"), value: overview.data?.data.projects },
    { label: t("users"), value: overview.data?.data.users },
    { label: t("messages24h"), value: overview.data?.data.messages_24h },
    { label: t("pending"), value: overview.data?.data.pending },
  ];

  return (
    <>
      <PageHeader title={t("title")} description={t("subtitle")} />

      <div className="mb-5 grid grid-cols-2 gap-3 sm:grid-cols-4">
        {kpis.map((k) => (
          <div key={k.label} className="rounded-lg border bg-card p-4">
            <div className="text-xs font-medium uppercase tracking-wide text-muted-foreground">{k.label}</div>
            <div className="tabular mt-1 font-mono text-2xl font-semibold">{k.value ?? "—"}</div>
          </div>
        ))}
      </div>

      <Tabs defaultValue="messages">
        <TabsList>
          <TabsTrigger value="messages">{t("tabMessages")}</TabsTrigger>
          <TabsTrigger value="tenants">{t("tabTenants")}</TabsTrigger>
          <TabsTrigger value="users">{t("tabUsers")}</TabsTrigger>
        </TabsList>

        <TabsContent value="messages">
          {messages.isLoading ? <Skeleton className="h-64" /> : !messages.data?.data.length ? <EmptyState /> : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader><TableRow>
                  <TableHead>{t("when")}</TableHead><TableHead>{t("tenant")}</TableHead>
                  <TableHead>{t("to")}</TableHead><TableHead>{t("channel")}</TableHead><TableHead>{t("status")}</TableHead>
                </TableRow></TableHeader>
                <TableBody>
                  {messages.data.data.map((r) => (
                    <TableRow key={r.message.id}>
                      <TableCell className="tabular whitespace-nowrap text-xs text-muted-foreground">{formatDate(r.message.created_at, { dateStyle: "short", timeStyle: "short" })}</TableCell>
                      <TableCell className="text-sm font-medium">{r.tenant.name}</TableCell>
                      <TableCell className="font-mono text-xs">{r.message.to}</TableCell>
                      <TableCell><Badge variant="outline" className="font-mono">{r.message.channel}</Badge></TableCell>
                      <TableCell><Badge variant={statusVariant(r.message.status)} className="font-mono">{r.message.status}</Badge></TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </TabsContent>

        <TabsContent value="tenants">
          {tenants.isLoading ? <Skeleton className="h-64" /> : !tenants.data?.data.length ? <EmptyState /> : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader><TableRow><TableHead>{t("name")}</TableHead><TableHead>{t("slug")}</TableHead><TableHead>{t("status")}</TableHead><TableHead>{t("when")}</TableHead></TableRow></TableHeader>
                <TableBody>
                  {tenants.data.data.map((p) => (
                    <TableRow key={p.id}>
                      <TableCell className="text-sm font-medium">{p.name}</TableCell>
                      <TableCell className="font-mono text-xs">{p.slug}</TableCell>
                      <TableCell><Badge variant={p.status === "active" ? "default" : "secondary"}>{p.status}</Badge></TableCell>
                      <TableCell className="tabular whitespace-nowrap text-xs text-muted-foreground">{formatDate(p.created_at, { dateStyle: "short" })}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </TabsContent>

        <TabsContent value="users">
          {users.isLoading ? <Skeleton className="h-64" /> : !users.data?.data.length ? <EmptyState /> : (
            <div className="rounded-lg border">
              <Table>
                <TableHeader><TableRow><TableHead>{t("user")}</TableHead><TableHead>{t("phone")}</TableHead><TableHead>{t("role")}</TableHead><TableHead>{t("when")}</TableHead></TableRow></TableHeader>
                <TableBody>
                  {users.data.data.map((u) => (
                    <TableRow key={u.id}>
                      <TableCell className="text-sm">{u.full_name || u.email || "—"}</TableCell>
                      <TableCell className="font-mono text-xs">{u.phone || "—"}</TableCell>
                      <TableCell>{u.platform_admin ? <Badge>{t("owner")}</Badge> : <Badge variant="secondary">{t("client")}</Badge>}</TableCell>
                      <TableCell className="tabular whitespace-nowrap text-xs text-muted-foreground">{formatDate(u.created_at, { dateStyle: "short" })}</TableCell>
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
