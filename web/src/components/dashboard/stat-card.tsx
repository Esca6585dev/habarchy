import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";

export function StatCard({ label, value, sub, tone }: { label: string; value: string | number; sub?: string; tone?: "success" | "destructive" | "info" | "warning" }) {
  return (
    <Card>
      <CardContent className="p-4">
        <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">{label}</p>
        <p className={cn("tabular mt-1 text-2xl font-semibold", tone === "success" && "text-success", tone === "destructive" && "text-destructive", tone === "info" && "text-info", tone === "warning" && "text-warning")}>
          {value}
        </p>
        {sub ? <p className="mt-0.5 text-xs text-muted-foreground">{sub}</p> : null}
      </CardContent>
    </Card>
  );
}
