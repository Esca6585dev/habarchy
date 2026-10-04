"use client";
import { useQuery } from "@tanstack/react-query";
import { api, unwrap } from "@/lib/api/client";
import type { Project, User } from "@/lib/api/types";

export const qk = {
  me: ["me"] as const,
  projects: ["projects"] as const,
  project: (id: string) => ["project", id] as const,
  members: (id: string) => ["members", id] as const,
  apiKeys: (id: string) => ["api-keys", id] as const,
  providers: (id: string) => ["providers", id] as const,
  templates: (id: string) => ["templates", id] as const,
  template: (id: string, tid: string) => ["template", id, tid] as const,
  versions: (id: string, tid: string) => ["template-versions", id, tid] as const,
  messages: (id: string, f: Record<string, unknown>) => ["messages", id, f] as const,
  message: (id: string, mid: string) => ["message", id, mid] as const,
  webhooks: (id: string, f: Record<string, unknown>) => ["webhooks", id, f] as const,
  contacts: (id: string, f: Record<string, unknown>) => ["contacts", id, f] as const,
  devices: (id: string) => ["devices", id] as const,
  dashboard: (id: string, days: number) => ["dashboard", id, days] as const,
  usage: (id: string, f: Record<string, unknown>) => ["usage", id, f] as const,
  health: (id: string) => ["health", id] as const,
  audit: (id: string, page: number) => ["audit", id, page] as const,
  users: ["users"] as const,
};

export function useMe() {
  return useQuery({
    queryKey: qk.me,
    queryFn: async () => unwrap<User>(await api.GET("/api/admin/me")),
  });
}

export function useProjects() {
  return useQuery({
    queryKey: qk.projects,
    queryFn: async () => unwrap<Project[]>(await api.GET("/api/admin/projects")),
  });
}
