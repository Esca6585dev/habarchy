export function JsonView({ value, className = "" }: { value: unknown; className?: string }) {
  let text = "";
  try {
    text = typeof value === "string" ? value : JSON.stringify(value, null, 2);
  } catch {
    text = String(value);
  }
  return (
    <pre className={`max-h-80 overflow-auto rounded-md border bg-muted/50 p-3 font-mono text-xs leading-relaxed ${className}`}>
      {text || "{}"}
    </pre>
  );
}
