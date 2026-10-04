"use client";
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useProjects } from "@/lib/hooks";
import type { Project } from "@/lib/api/types";

interface Ctx {
  projectId: string | null;
  project: Project | null;
  projects: Project[];
  isLoading: boolean;
  setProjectId: (id: string) => void;
  role: Project["role"] | undefined;
  can: (min: "viewer" | "developer" | "admin" | "owner") => boolean;
}

const ProjectContext = createContext<Ctx | null>(null);
const rank = { viewer: 1, developer: 2, admin: 3, owner: 4 } as const;

export function ProjectProvider({ initialProjectId, children }: { initialProjectId: string | null; children: React.ReactNode }) {
  const { data: projects = [], isLoading } = useProjects();
  const [projectId, setId] = useState<string | null>(initialProjectId);
  const qc = useQueryClient();

  // Fall back to the first project when the cookie is missing or stale.
  useEffect(() => {
    if (isLoading) return;
    if (projects.length && (!projectId || !projects.some((p) => p.id === projectId))) {
      const first = projects[0].id!;
      setId(first);
      void fetch("/api/project", { method: "POST", body: JSON.stringify({ project_id: first }) });
    }
  }, [projects, projectId, isLoading]);

  const setProjectId = useCallback(
    (id: string) => {
      setId(id);
      void fetch("/api/project", { method: "POST", body: JSON.stringify({ project_id: id }) }).then(() => qc.invalidateQueries());
    },
    [qc],
  );

  const project = useMemo(() => projects.find((p) => p.id === projectId) ?? null, [projects, projectId]);
  const role = project?.role;
  const can = useCallback((min: keyof typeof rank) => (role ? rank[role as keyof typeof rank] >= rank[min] : false), [role]);

  return (
    <ProjectContext.Provider value={{ projectId, project, projects, isLoading, setProjectId, role, can }}>{children}</ProjectContext.Provider>
  );
}

export function useProject() {
  const ctx = useContext(ProjectContext);
  if (!ctx) throw new Error("useProject outside ProjectProvider");
  return ctx;
}

/** Throws a stable project id for pages; renders nothing until known. */
export function useRequiredProjectId(): string {
  const { projectId } = useProject();
  return projectId ?? "";
}
