"use client";
import { useTranslations } from "next-intl";
import { Badge } from "@/components/ui/badge";

const variants: Record<string, "success" | "destructive" | "warning" | "info" | "secondary" | "outline"> = {
  delivered: "success",
  sent: "info",
  queued: "secondary",
  processing: "warning",
  failed: "destructive",
  cancelled: "outline",
  pending: "secondary",
  completed: "success",
  healthy: "success",
  degraded: "warning",
  failing: "destructive",
  idle: "secondary",
  disabled: "outline",
};

export function StatusBadge({ status }: { status?: string | null }) {
  const t = useTranslations("status");
  const th = useTranslations("health");
  const tw = useTranslations("webhooks");
  if (!status) return null;
  let label: string = status;
  if (["queued", "processing", "sent", "delivered", "failed", "cancelled"].includes(status)) label = t(status as "queued");
  else if (["healthy", "degraded", "failing", "idle", "disabled"].includes(status)) label = th(status as "healthy");
  else if (status === "pending") label = tw("pending");
  return <Badge variant={variants[status] ?? "secondary"}>{label}</Badge>;
}

const channelColor: Record<string, string> = {
  sms: "bg-chart-sms/15 text-chart-sms",
  email: "bg-chart-email/15 text-chart-email",
  push: "bg-chart-push/20 text-chart-push",
  telegram: "bg-chart-telegram/15 text-chart-telegram",
};

export function ChannelBadge({ channel }: { channel?: string | null }) {
  const t = useTranslations("channel");
  if (!channel) return null;
  return (
    <span className={`inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium ${channelColor[channel] ?? "bg-muted"}`}>
      {t(channel as "sms")}
    </span>
  );
}
