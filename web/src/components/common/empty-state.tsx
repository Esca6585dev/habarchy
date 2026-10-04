"use client";
import { Inbox } from "lucide-react";
import { useTranslations } from "next-intl";

export function EmptyState({ title, hint, action }: { title?: string; hint?: string; action?: React.ReactNode }) {
  const t = useTranslations("common");
  return (
    <div className="flex flex-col items-center justify-center rounded-lg border border-dashed px-6 py-12 text-center">
      <Inbox className="mb-3 h-8 w-8 text-muted-foreground" />
      <p className="font-medium">{title ?? t("empty")}</p>
      {hint ? <p className="mt-1 max-w-sm text-sm text-muted-foreground">{hint}</p> : null}
      {action ? <div className="mt-4">{action}</div> : null}
    </div>
  );
}
