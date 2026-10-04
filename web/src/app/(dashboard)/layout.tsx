import { cookies } from "next/headers";
import { PROJECT_COOKIE } from "@/i18n/config";
import { ProjectProvider } from "@/components/project-context";
import { Sidebar } from "@/components/shell/sidebar";
import { Topbar } from "@/components/shell/topbar";

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const store = await cookies();
  const projectId = store.get(PROJECT_COOKIE)?.value ?? null;
  return (
    <ProjectProvider initialProjectId={projectId}>
      <div className="flex min-h-dvh">
        <Sidebar className="sticky top-0 hidden h-dvh lg:flex" />
        <div className="flex min-w-0 flex-1 flex-col">
          <Topbar />
          <main className="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6">{children}</main>
        </div>
      </div>
    </ProjectProvider>
  );
}
