// Chat API shapes (the endpoints return ad-hoc JSON, not schema components).
export type ChatChannel = {
  id: string;
  kind: "public" | "private" | "direct";
  name: string;
  topic: string;
  created_by?: string;
  created_at: string;
  unread: number;
  last_body?: string;
  last_at?: string | null;
  last_has_file?: boolean;
  peer_id?: string;
  peer_name?: string;
  peer_avatar?: string | null;
};

export type ChatMessage = {
  id: string;
  channel_id: string;
  user_id?: string;
  author_name: string;
  author_avatar?: string | null;
  body: string;
  attachment_id?: string | null;
  attachment_type?: string;
  attachment_name?: string;
  created_at: string;
  edited_at?: string | null;
};

export type ChatMember = {
  id: string;
  full_name: string;
  email: string;
  avatar_id?: string | null;
  bio?: string;
  role: string;
};

/** URL of an attachment image, served through the proxy with the auth cookie. */
export function attachmentUrl(id?: string | null): string | undefined {
  return id ? `/api/backend/api/admin/attachments/${id}` : undefined;
}

/** A channel's display name: the peer for DMs, otherwise the channel name. */
export function channelTitle(c: ChatChannel): string {
  return c.kind === "direct" ? c.peer_name || "—" : c.name;
}
