"use client";
import Link from "next/link";
import { useTranslations } from "next-intl";
import { ChevronsUpDown, FolderKanban, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { useProject } from "@/components/project-context";

export function ProjectSwitcher() {
  const { project, projects, setProjectId } = useProject();
  const t = useTranslations("projects");
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="w-full justify-between" data-testid="project-switcher">
          <span className="flex min-w-0 items-center gap-2">
            <FolderKanban className="shrink-0" />
            <span className="truncate">{project?.name ?? t("title")}</span>
          </span>
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-64" align="start">
        <DropdownMenuLabel>{t("title")}</DropdownMenuLabel>
        {projects.map((p) => (
          <DropdownMenuItem key={p.id} onSelect={() => setProjectId(p.id!)}>
            <span className="truncate">{p.name}</span>
            <span className="ml-auto text-xs text-muted-foreground">{p.role}</span>
          </DropdownMenuItem>
        ))}
        <DropdownMenuSeparator />
        <DropdownMenuItem asChild>
          <Link href="/projects">
            <Plus /> {t("new")}
          </Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
