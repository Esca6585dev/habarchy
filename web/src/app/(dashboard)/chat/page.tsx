"use client";
import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslations } from "next-intl";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Hash, ImagePlus, Lock, MessageSquarePlus, Plus, Send, Trash2, Users } from "lucide-react";
import { toast } from "sonner";
import { api, unwrap } from "@/lib/api/client";
import { useMe } from "@/lib/hooks";
import type { User } from "@/lib/api/types";
import type { ChatChannel, ChatMember, ChatMessage } from "@/lib/chat";
import { attachmentUrl, channelTitle } from "@/lib/chat";
import { Avatar } from "@/components/chat/avatar";
import { PageHeader } from "@/components/common/page-header";
import { EmptyState } from "@/components/common/empty-state";
import { Field } from "@/components/common/field";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";

export default function ChatPage() {
  const t = useTranslations("chat");
  const tc = useTranslations("common");
  const qc = useQueryClient();
  const { data: me } = useMe();
  const [active, setActive] = useState<string | null>(null);
  const [text, setText] = useState("");
  const [creating, setCreating] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  const channels = useQuery({
    queryKey: ["chat", "channels"],
    queryFn: async () => unwrap<ChatChannel[]>(await api.GET("/api/admin/chat/channels")),
    refetchInterval: 20_000,
  });
  const current = channels.data?.find((c) => c.id === active) ?? null;
  const messages = useQuery({
    queryKey: ["chat", "messages", active],
    enabled: !!active,
    queryFn: async () => unwrap<ChatMessage[]>(await api.GET("/api/admin/chat/channels/{channel_id}/messages", { params: { path: { channel_id: active! } } })),
  });
  const members = useQuery({
    queryKey: ["chat", "members", active],
    enabled: !!active && current?.kind !== "direct",
    queryFn: async () => unwrap<ChatMember[]>(await api.GET("/api/admin/chat/channels/{channel_id}/members", { params: { path: { channel_id: active! } } })),
  });

  // Live updates via SSE: refetch the affected channel's messages and the list.
  useEffect(() => {
    let es: EventSource | null = null;
    let closed = false;
    let retry: ReturnType<typeof setTimeout>;
    const connect = () => {
      es = new EventSource("/api/backend/api/admin/chat/stream");
      const onEv = (e: MessageEvent) => {
        try {
          const ev = JSON.parse(e.data) as { channel_id: string };
          void qc.invalidateQueries({ queryKey: ["chat", "channels"] });
          void qc.invalidateQueries({ queryKey: ["chat", "messages", ev.channel_id] });
        } catch { /* ignore */ }
      };
      for (const typ of ["message", "message.deleted", "channel"]) es.addEventListener(typ, onEv as EventListener);
      es.onerror = () => { es?.close(); if (!closed) retry = setTimeout(connect, 3000); };
    };
    connect();
    return () => { closed = true; es?.close(); clearTimeout(retry); };
  }, [qc]);

  // Mark read + scroll to bottom when a channel opens or new messages arrive.
  useEffect(() => {
    if (!active) return;
    void api.POST("/api/admin/chat/channels/{channel_id}/read", { params: { path: { channel_id: active } } }).then(() => qc.invalidateQueries({ queryKey: ["chat", "channels"] }));
  }, [active, messages.data, qc]);
  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: "auto" }); }, [messages.data, active]);

  const ordered = useMemo(() => [...(messages.data ?? [])].reverse(), [messages.data]);

  const send = useMutation({
    mutationFn: async (attachmentId?: string) => {
      if (!active) return;
      unwrap<ChatMessage>(await api.POST("/api/admin/chat/channels/{channel_id}/messages", { params: { path: { channel_id: active } }, body: { body: text.trim(), attachment_id: attachmentId } }));
    },
    onSuccess: () => { setText(""); void qc.invalidateQueries({ queryKey: ["chat", "messages", active] }); void qc.invalidateQueries({ queryKey: ["chat", "channels"] }); },
    onError: (e: Error) => toast.error(e.message),
  });
  const del = useMutation({
    mutationFn: async (id: string) => { await api.DELETE("/api/admin/chat/channels/{channel_id}/messages/{message_id}", { params: { path: { channel_id: active!, message_id: id } } }); },
    onSuccess: () => void qc.invalidateQueries({ queryKey: ["chat", "messages", active] }),
  });

  const uploadAndSend = async (file: File) => {
    const fd = new FormData();
    fd.append("file", file);
    const res = await fetch("/api/backend/api/admin/attachments", { method: "POST", body: fd });
    if (!res.ok) return void toast.error(tc("error"));
    const id = ((await res.json()) as { data: { id: string } }).data.id;
    send.mutate(id);
  };

  return (
    <>
      <PageHeader title={t("title")} actions={<Button size="sm" onClick={() => setCreating(true)} data-testid="new-channel"><MessageSquarePlus /> {t("new")}</Button>} />
      <div className="grid h-[calc(100vh-11rem)] grid-cols-1 gap-4 md:grid-cols-[18rem_1fr]">
        {/* Channel list */}
        <div className="overflow-y-auto rounded-lg border">
          {channels.isLoading ? <Skeleton className="h-64 m-2" /> : !channels.data?.length ? <EmptyState hint={t("empty")} /> : (
            <ul className="divide-y">
              {channels.data.map((c) => (
                <li key={c.id}>
                  <button className={`flex w-full items-center gap-2 px-3 py-2.5 text-left hover:bg-accent ${active === c.id ? "bg-accent" : ""}`} onClick={() => setActive(c.id)} data-testid={`channel-${c.id}`}>
                    {c.kind === "direct" ? <Avatar name={c.peer_name} avatarId={c.peer_avatar} size={32} /> : <span className="flex h-8 w-8 items-center justify-center rounded-full bg-muted">{c.kind === "private" ? <Lock className="h-4 w-4" /> : <Hash className="h-4 w-4" />}</span>}
                    <span className="min-w-0 flex-1">
                      <span className="flex items-center gap-1"><span className="truncate font-medium">{channelTitle(c)}</span>{c.unread > 0 ? <Badge className="ml-auto">{c.unread}</Badge> : null}</span>
                      <span className="block truncate text-xs text-muted-foreground">{c.last_has_file ? "📷 " : ""}{c.last_body || t("noMessages")}</span>
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Message pane */}
        <div className="flex min-h-0 flex-col rounded-lg border">
          {!current ? <div className="flex flex-1 items-center justify-center text-sm text-muted-foreground">{t("pick")}</div> : (
            <>
              <div className="flex items-center gap-2 border-b px-4 py-2.5">
                {current.kind === "direct" ? <Avatar name={current.peer_name} avatarId={current.peer_avatar} size={28} /> : null}
                <span className="font-medium">{channelTitle(current)}</span>
                {current.kind !== "direct" && members.data ? <span className="flex items-center gap-1 text-xs text-muted-foreground"><Users className="h-3.5 w-3.5" />{members.data.length}</span> : null}
                {current.topic ? <span className="truncate text-xs text-muted-foreground">· {current.topic}</span> : null}
              </div>
              <div className="flex-1 space-y-3 overflow-y-auto p-4">
                {messages.isLoading ? <Skeleton className="h-40" /> : ordered.length === 0 ? <EmptyState hint={t("noMessages")} /> : ordered.map((m) => {
                  const mine = m.user_id === (me as User | undefined)?.id;
                  return (
                    <div key={m.id} className={`flex gap-2 ${mine ? "flex-row-reverse" : ""}`}>
                      <Avatar name={m.author_name} avatarId={m.author_avatar} size={32} />
                      <div className={`group max-w-[75%] ${mine ? "items-end text-right" : ""}`}>
                        <div className="flex items-center gap-2 text-xs text-muted-foreground"><span className="font-medium text-foreground">{m.author_name}</span><time>{new Date(m.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time>{mine ? <button className="opacity-0 group-hover:opacity-100" onClick={() => del.mutate(m.id)} title={tc("delete")}><Trash2 className="h-3 w-3" /></button> : null}</div>
                        <div className={`mt-0.5 inline-block rounded-2xl px-3 py-1.5 text-sm ${mine ? "rounded-br-sm bg-primary text-primary-foreground" : "rounded-bl-sm bg-muted"}`}>
                          {m.attachment_id ? (
                            // eslint-disable-next-line @next/next/no-img-element
                            <a href={attachmentUrl(m.attachment_id)} target="_blank" rel="noreferrer"><img src={attachmentUrl(m.attachment_id)} alt="" className="mb-1 max-h-60 rounded-lg" /></a>
                          ) : null}
                          {m.body ? <span className="whitespace-pre-wrap break-words">{m.body}</span> : null}
                        </div>
                      </div>
                    </div>
                  );
                })}
                <div ref={bottomRef} />
              </div>
              <form className="flex items-end gap-2 border-t p-3" onSubmit={(e) => { e.preventDefault(); if (text.trim()) send.mutate(undefined); }}>
                <input ref={fileRef} type="file" accept="image/*" className="hidden" onChange={(e) => { const f = e.target.files?.[0]; if (f) void uploadAndSend(f); e.target.value = ""; }} />
                <Button type="button" size="icon" variant="ghost" onClick={() => fileRef.current?.click()} title={t("attachImage")}><ImagePlus /></Button>
                <Textarea rows={1} value={text} onChange={(e) => setText(e.target.value)} onKeyDown={(e) => { if (e.key === "Enter" && !e.shiftKey) { e.preventDefault(); if (text.trim()) send.mutate(undefined); } }} placeholder={t("message")} className="max-h-32 min-h-10 resize-none" data-testid="chat-input" />
                <Button type="submit" size="icon" disabled={!text.trim() || send.isPending} data-testid="chat-send"><Send /></Button>
              </form>
            </>
          )}
        </div>
      </div>
      <NewChannelDialog open={creating} onOpenChange={setCreating} onOpen={setActive} />
    </>
  );
}

function NewChannelDialog({ open, onOpenChange, onOpen }: { open: boolean; onOpenChange: (o: boolean) => void; onOpen: (id: string) => void }) {
  const t = useTranslations("chat");
  const tc = useTranslations("common");
  const qc = useQueryClient();
  const { data: me } = useMe();
  const [tab, setTab] = useState<"channel" | "direct">("channel");
  const [kind, setKind] = useState("public");
  const [name, setName] = useState("");
  const [picked, setPicked] = useState<Set<string>>(new Set());

  const users = useQuery({
    queryKey: ["users"],
    enabled: open,
    queryFn: async () => unwrap<User[]>(await api.GET("/api/admin/users")),
  });
  const others = (users.data ?? []).filter((u) => u.id && u.id !== (me as User | undefined)?.id) as (User & { id: string })[];

  const create = useMutation({
    mutationFn: async () => unwrap<{ id: string }>(await api.POST("/api/admin/chat/channels", { body: { kind: kind as "public" | "private", name: name.trim(), members: Array.from(picked) } })),
    onSuccess: (c) => { void qc.invalidateQueries({ queryKey: ["chat", "channels"] }); onOpen(c.id); onOpenChange(false); setName(""); setPicked(new Set()); },
    onError: (e: Error) => toast.error(e.message),
  });
  const direct = useMutation({
    mutationFn: async (userId: string) => unwrap<{ id: string }>(await api.POST("/api/admin/chat/direct", { body: { user_id: userId } })),
    onSuccess: (c) => { void qc.invalidateQueries({ queryKey: ["chat", "channels"] }); onOpen(c.id); onOpenChange(false); },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader><DialogTitle>{t("new")}</DialogTitle></DialogHeader>
        <div className="mb-2 flex gap-2">
          <Button size="sm" variant={tab === "channel" ? "default" : "outline"} onClick={() => setTab("channel")}><Plus /> {t("channel")}</Button>
          <Button size="sm" variant={tab === "direct" ? "default" : "outline"} onClick={() => setTab("direct")}><Avatar name="" size={16} /> {t("direct")}</Button>
        </div>
        {tab === "channel" ? (
          <div className="space-y-3">
            <Field label={t("kind")}>
              <Select value={kind} onValueChange={setKind}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent><SelectItem value="public">{t("public")}</SelectItem><SelectItem value="private">{t("private")}</SelectItem></SelectContent>
              </Select>
            </Field>
            <Field label={t("name")} htmlFor="ch-name"><Input id="ch-name" value={name} onChange={(e) => setName(e.target.value)} placeholder="general" /></Field>
            {kind === "private" ? (
              <Field label={t("members")}>
                <div className="max-h-48 space-y-1 overflow-auto rounded-md border p-2 text-sm">
                  {others.map((u) => { const on = picked.has(u.id); return (
                    <label key={u.id} className="flex cursor-pointer items-center gap-2 rounded px-1 py-0.5 hover:bg-accent">
                      <input type="checkbox" checked={on} onChange={() => { const n = new Set(picked); if (on) n.delete(u.id); else n.add(u.id); setPicked(n); }} />
                      <Avatar name={u.full_name} avatarId={u.avatar_id} size={24} /> {u.full_name || u.email}
                    </label>
                  ); })}
                </div>
              </Field>
            ) : null}
            <DialogFooter><Button variant="outline" onClick={() => onOpenChange(false)}>{tc("cancel")}</Button><Button onClick={() => create.mutate()} disabled={!name.trim() || create.isPending}>{tc("save")}</Button></DialogFooter>
          </div>
        ) : (
          <div className="max-h-72 space-y-1 overflow-auto">
            {others.map((u) => (
              <button key={u.id} className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-left hover:bg-accent" onClick={() => direct.mutate(u.id)}>
                <Avatar name={u.full_name} avatarId={u.avatar_id} size={32} /> <span>{u.full_name || u.email}</span>
              </button>
            ))}
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
