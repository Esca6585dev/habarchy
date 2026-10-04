"use client";
import { useMemo, useState } from "react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useInfiniteQuery } from "@tanstack/react-query";
import { RefreshCw } from "lucide-react";
import { api, unwrap } from "@/lib/api/client";
import type { Message } from "@/lib/api/types";
import { useProject } from "@/components/project-context";
import { PageHeader } from "@/components/common/page-header";
import { ChannelBadge, StatusBadge } from "@/components/common/status-badge";
import { EmptyState } from "@/components/common/empty-state";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDate, truncate } from "@/lib/utils";

const statuses = ["queued", "processing", "sent", "delivered", "failed", "cancelled"] as const;
const channels = ["sms", "email", "push", "telegram"] as const;

export default function MessagesPage() {
  const t = useTranslations("messages");
  const tc = useTranslations("common");
  const ts = useTranslations("status");
  const tch = useTranslations("channel");
  const { projectId } = useProject();
  const [status, setStatus] = useState<string>("all");
  const [channel, setChannel] = useState<string>("all");
  const [search, setSearch] = useState("");
  const [includeTest, setIncludeTest] = useState(true);
  const [debounced, setDebounced] = useState("");

  const filters = useMemo(() => ({ status, channel, search: debounced, includeTest }), [status, channel, debounced, includeTest]);
  const q = useInfiniteQuery({
    queryKey: ["messages", projectId, filters],
    enabled: !!projectId,
    initialPageParam: undefined as string | undefined,
    refetchInterval: 15_000,
    queryFn: async ({ pageParam }) => {
      const res = await api.GET("/api/admin/projects/{project_id}/messages", {
        params: {
          path: { project_id: projectId! },
          query: {
            status: status === "all" ? undefined : (status as Message["status"]),
            channel: channel === "all" ? undefined : (channel as Message["channel"]),
            search: debounced || undefined,
            include_test: includeTest,
            cursor: pageParam,
            limit: 50,
          },
        },
      });
      const data = unwrap<Message[]>(res);
      const meta = (res.data as { meta?: { next_cursor?: string | null } })?.meta;
      return { rows: data, next: meta?.next_cursor ?? undefined };
    },
    getNextPageParam: (last) => last.next,
  });
  const rows = q.data?.pages.flatMap((p) => p.rows) ?? [];

  return (
    <>
      <PageHeader
        title={t("title")}
        actions={<Button variant="outline" size="sm" onClick={() => q.refetch()}><RefreshCw className={q.isFetching ? "animate-spin" : ""} /> {tc("refresh")}</Button>}
      />
      <div className="mb-4 flex flex-wrap items-end gap-3">
        <form className="flex-1 min-w-[12rem]" onSubmit={(e) => { e.preventDefault(); setDebounced(search.trim()); }}>
          <Input placeholder={`${tc("search")}…`} value={search} onChange={(e) => setSearch(e.target.value)} onBlur={() => setDebounced(search.trim())} data-testid="messages-search" />
        </form>
        <Select value={status} onValueChange={setStatus}>
          <SelectTrigger className="w-40"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{tc("all")}</SelectItem>
            {statuses.map((s) => <SelectItem key={s} value={s}>{ts(s)}</SelectItem>)}
          </SelectContent>
        </Select>
        <Select value={channel} onValueChange={setChannel}>
          <SelectTrigger className="w-36"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{tc("all")}</SelectItem>
            {channels.map((c) => <SelectItem key={c} value={c}>{tch(c)}</SelectItem>)}
          </SelectContent>
        </Select>
        <div className="flex items-center gap-2">
          <Switch id="include-test" checked={includeTest} onCheckedChange={setIncludeTest} />
          <Label htmlFor="include-test" className="text-xs">{t("includeTest")}</Label>
        </div>
      </div>
      {q.isLoading ? (
        <Skeleton className="h-64" />
      ) : rows.length === 0 ? (
        <EmptyState />
      ) : (
        <div className="rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("status")}</TableHead>
                <TableHead>{t("channel")}</TableHead>
                <TableHead>{t("to")}</TableHead>
                <TableHead>{t("template")} / {t("body")}</TableHead>
                <TableHead>{t("attempts")}</TableHead>
                <TableHead>{t("error")}</TableHead>
                <TableHead>{tc("created")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((m) => (
                <TableRow key={m.id} className="cursor-pointer" data-testid="message-row">
                  <TableCell><StatusBadge status={m.status} /></TableCell>
                  <TableCell><ChannelBadge channel={m.channel} /> {m.is_test ? <span className="ml-1 text-[10px] uppercase text-muted-foreground">test</span> : null}</TableCell>
                  <TableCell className="font-mono text-xs"><Link href={`/messages/${m.id}`} className="hover:underline">{m.to}</Link></TableCell>
                  <TableCell className="max-w-[20rem] truncate text-xs text-muted-foreground">{m.template ? <span className="font-medium text-foreground">{m.template}</span> : null} {truncate(m.body ?? "", 70)}</TableCell>
                  <TableCell className="tabular">{m.attempts}</TableCell>
                  <TableCell className="max-w-[12rem] truncate text-xs text-destructive">{m.error_code}</TableCell>
                  <TableCell className="tabular whitespace-nowrap text-xs text-muted-foreground">{formatDate(m.created_at)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          {q.hasNextPage ? (
            <div className="border-t p-3 text-center">
              <Button variant="outline" size="sm" onClick={() => q.fetchNextPage()} disabled={q.isFetchingNextPage}>{tc("loadMore")}</Button>
            </div>
          ) : null}
        </div>
      )}
    </>
  );
}
