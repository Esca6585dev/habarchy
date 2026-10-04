"use client";
import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { useQueryClient } from "@tanstack/react-query";
import { Radio } from "lucide-react";
import { ChannelBadge, StatusBadge } from "@/components/common/status-badge";
import { formatDate } from "@/lib/utils";

interface LiveEvent {
  type: string;
  message_id?: string;
  channel?: string;
  status?: string;
  to?: string;
  provider?: string;
  error_code?: string;
  at: string;
}

/** Subscribes to the SSE stream (through the proxy) and shows the last events. */
export function LiveFeed({ projectId }: { projectId: string }) {
  const t = useTranslations("dashboard");
  const qc = useQueryClient();
  const [events, setEvents] = useState<LiveEvent[]>([]);
  const [connected, setConnected] = useState(false);
  const retry = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => {
    let es: EventSource | null = null;
    let closed = false;
    const connect = () => {
      es = new EventSource(`/api/backend/api/admin/stream?project_id=${projectId}`);
      es.addEventListener("ready", () => setConnected(true));
      const onMsg = (e: MessageEvent) => {
        try {
          const ev = JSON.parse(e.data) as LiveEvent;
          setEvents((prev) => [ev, ...prev].slice(0, 30));
          if (ev.type !== "message.queued") void qc.invalidateQueries({ queryKey: ["dashboard", projectId] });
        } catch {
          /* ignore */
        }
      };
      for (const typ of ["message.queued", "message.sent", "message.delivered", "message.failed", "message.cancelled"]) {
        es.addEventListener(typ, onMsg as EventListener);
      }
      es.onerror = () => {
        setConnected(false);
        es?.close();
        if (!closed) retry.current = setTimeout(connect, 3000);
      };
    };
    connect();
    return () => {
      closed = true;
      es?.close();
      if (retry.current) clearTimeout(retry.current);
    };
  }, [projectId, qc]);

  return (
    <div>
      <div className="mb-2 flex items-center gap-2 text-xs text-muted-foreground">
        <Radio className={`h-3.5 w-3.5 ${connected ? "text-success" : "text-muted-foreground"}`} />
        {t("liveHint")}
      </div>
      <ul className="divide-y rounded-md border">
        {events.length === 0 ? (
          <li className="px-3 py-6 text-center text-sm text-muted-foreground">{t("waiting")}</li>
        ) : (
          events.map((ev, i) => (
            <li key={`${ev.message_id}-${ev.type}-${i}`} className="flex items-center gap-3 px-3 py-2 text-sm">
              <StatusBadge status={ev.status} />
              <ChannelBadge channel={ev.channel} />
              <Link href={`/messages/${ev.message_id}`} className="min-w-0 flex-1 truncate font-mono text-xs hover:underline">
                {ev.to}
              </Link>
              {ev.error_code ? <span className="truncate text-xs text-destructive">{ev.error_code}</span> : null}
              <span className="tabular shrink-0 text-xs text-muted-foreground">{formatDate(ev.at, { timeStyle: "medium" })}</span>
            </li>
          ))
        )}
      </ul>
    </div>
  );
}
