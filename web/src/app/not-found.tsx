import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <div className="flex min-h-dvh flex-col items-center justify-center gap-4 p-6 text-center">
      <p className="text-6xl font-semibold tabular">404</p>
      <p className="text-muted-foreground">Page not found</p>
      <Button asChild><Link href="/">Habarchy</Link></Button>
    </div>
  );
}
