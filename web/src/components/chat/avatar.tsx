import { attachmentUrl } from "@/lib/chat";

/** Round avatar: the uploaded image, or initials on a tinted circle. */
export function Avatar({ name, avatarId, size = 36 }: { name?: string; avatarId?: string | null; size?: number }) {
  const url = attachmentUrl(avatarId);
  const initials = (name ?? "?").trim().split(/\s+/).slice(0, 2).map((w) => w[0]?.toUpperCase() ?? "").join("") || "?";
  if (url) {
    // eslint-disable-next-line @next/next/no-img-element
    return <img src={url} alt={name ?? ""} width={size} height={size} className="shrink-0 rounded-full object-cover" style={{ width: size, height: size }} />;
  }
  return (
    <span className="flex shrink-0 items-center justify-center rounded-full bg-primary/15 font-medium text-primary" style={{ width: size, height: size, fontSize: size * 0.4 }}>
      {initials}
    </span>
  );
}
