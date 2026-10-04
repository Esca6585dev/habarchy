/** "Name, +99365…, mail@…" → contact request. Phones start with + or digits, e-mails contain @. */
export function parseInlineContact(line: string): Record<string, string> | null {
  const parts = line.split(/[,;\t]/).map((p) => p.trim()).filter(Boolean);
  if (!parts.length) return null;
  const out: Record<string, string> = {};
  for (const p of parts) {
    if (p.includes("@")) out.email = p;
    else if (/^\+?[\d\s()-]{6,}$/.test(p)) out.phone ??= p;
    else if (/^[UC][A-Z0-9]{6,}$/.test(p)) out.slack_id = p;
    else if (/^-?\d{5,}$/.test(p)) out.telegram_chat_id = p;
    else out.name ??= p;
  }
  return out.phone || out.email || out.slack_id || out.telegram_chat_id ? out : null;
}

/** Display label for a contact row. */
export function contactLabel(c: { name?: string; external_id?: string; phone?: string; email?: string; telegram_chat_id?: string; slack_id?: string; id?: string }): string {
  return c.name || c.external_id || c.phone || c.email || c.telegram_chat_id || c.slack_id || c.id?.slice(0, 8) || "";
}
