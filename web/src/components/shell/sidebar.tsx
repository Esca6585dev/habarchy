"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useTranslations } from "next-intl";
import {
  Activity, BookOpen, Contact, FileText, Gauge, KeyRound, LayoutDashboard, ListOrdered, MessageSquare, Plug, Settings, ShieldCheck, Webhook, Send, Users } from "lucide-react";
import { cn } from "@/lib/utils";
import { ProjectSwitcher } from "./project-switcher";
import { useProject } from "@/components/project-context";

export function Sidebar({ className, onNavigate }: { className?: string; onNavigate?: () => void }) {
  const t = useTranslations("nav");
  const pathname = usePathname();
  const { can } = useProject();
  const items = [
    { href: "/", label: t("dashboard"), icon: LayoutDashboard },
    { href: "/messages", label: t("messages"), icon: MessageSquare },
    { href: "/templates", label: t("templates"), icon: FileText },
    { href: "/providers", label: t("providers"), icon: Plug, min: "admin" as const },
    { href: "/api-keys", label: t("apiKeys"), icon: KeyRound, min: "developer" as const },
    { href: "/compose", label: t("compose"), icon: Send, min: "developer" as const },
    { href: "/groups", label: t("groups"), icon: Users },
    { href: "/contacts", label: t("contacts"), icon: Contact },
    { href: "/webhooks", label: t("webhooks"), icon: Webhook },
    { href: "/usage", label: t("usage"), icon: Gauge },
    { href: "/health", label: t("health"), icon: Activity },
    { href: "/audit", label: t("audit"), icon: ShieldCheck, min: "admin" as const },
    { href: "/settings", label: t("settings"), icon: Settings, min: "admin" as const },
  ];
  return (
    <aside className={cn("flex h-full w-64 flex-col gap-4 border-r bg-sidebar p-4", className)}>
      <Link href="/" className="flex items-center gap-2 px-1 text-lg font-semibold tracking-tight" onClick={onNavigate}>
        <span className="flex h-7 w-7 items-center justify-center rounded-md bg-primary text-primary-foreground text-sm">H</span>
        Habarchy
      </Link>
      <ProjectSwitcher />
      <nav className="flex flex-1 flex-col gap-0.5 overflow-y-auto">
        {items
          .filter((i) => !i.min || can(i.min))
          .map(({ href, label, icon: Icon }) => {
            const active = href === "/" ? pathname === "/" : pathname.startsWith(href);
            return (
              <Link
                key={href}
                href={href}
                onClick={onNavigate}
                className={cn(
                  "flex items-center gap-2.5 rounded-md px-2.5 py-2 text-sm transition-colors hover:bg-accent hover:text-accent-foreground",
                  active ? "bg-accent font-medium text-accent-foreground" : "text-muted-foreground",
                )}
              >
                <Icon className="h-4 w-4" />
                {label}
              </Link>
            );
          })}
      </nav>
      <div className="flex flex-col gap-0.5 border-t pt-3 text-sm">
        <a href="/api-docs" target="_blank" rel="noreferrer" className="flex items-center gap-2.5 rounded-md px-2.5 py-2 text-muted-foreground hover:bg-accent">
          <BookOpen className="h-4 w-4" /> {t("docs")}
        </a>
        {can("admin") ? (
          <a href="/admin/queues/" target="_blank" rel="noreferrer" className="flex items-center gap-2.5 rounded-md px-2.5 py-2 text-muted-foreground hover:bg-accent">
            <ListOrdered className="h-4 w-4" /> {t("queues")}
          </a>
        ) : null}
      </div>
    </aside>
  );
}
