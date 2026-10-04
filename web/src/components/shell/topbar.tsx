"use client";
import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { LogOut, Menu, UserRound } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { LocaleSwitcher } from "./locale-switcher";
import { ThemeToggle } from "./theme-toggle";
import { Sidebar } from "./sidebar";
import { useMe } from "@/lib/hooks";

export function Topbar() {
  const t = useTranslations("nav");
  const router = useRouter();
  const { data: me } = useMe();
  const [open, setOpen] = useState(false);
  return (
    <header className="sticky top-0 z-30 flex h-14 items-center gap-2 border-b bg-background/95 px-4 backdrop-blur">
      <Button variant="ghost" size="icon" className="lg:hidden" onClick={() => setOpen(true)} aria-label="Menu">
        <Menu />
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="left-0 top-0 h-full max-h-none w-72 translate-x-0 translate-y-0 rounded-none p-0">
          <DialogTitle className="sr-only">Menu</DialogTitle>
          <Sidebar className="border-r-0" onNavigate={() => setOpen(false)} />
        </DialogContent>
      </Dialog>
      <div className="ml-auto flex items-center gap-1">
        <LocaleSwitcher />
        <ThemeToggle />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="sm" data-testid="user-menu">
              <UserRound />
              <span className="hidden max-w-[12rem] truncate sm:inline">{me?.email}</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuLabel className="truncate">{me?.full_name || me?.email}</DropdownMenuLabel>
            <DropdownMenuItem asChild>
              <Link href="/profile">
                <UserRound /> {t("profile")}
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link href="/projects">{t("projects")}</Link>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onSelect={async () => {
                await fetch("/api/auth/logout", { method: "POST" });
                router.replace("/login");
                router.refresh();
              }}
            >
              <LogOut /> {t("logout")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
