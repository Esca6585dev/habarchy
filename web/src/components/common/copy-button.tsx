"use client";
import { useState } from "react";
import { Check, Copy } from "lucide-react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";

export function CopyButton({ value, size = "sm" }: { value: string; size?: "sm" | "icon" }) {
  const t = useTranslations("common");
  const [done, setDone] = useState(false);
  return (
    <Button
      type="button"
      variant="outline"
      size={size}
      onClick={async () => {
        try {
          await navigator.clipboard.writeText(value);
          setDone(true);
          setTimeout(() => setDone(false), 1500);
        } catch {
          /* clipboard unavailable */
        }
      }}
      aria-label={t("copy")}
    >
      {done ? <Check /> : <Copy />}
      {size === "sm" ? (done ? t("copied") : t("copy")) : null}
    </Button>
  );
}
