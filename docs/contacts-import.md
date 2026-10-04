# Importing contacts

Contacts can be imported into a project (and straight into a group) from files or from an
online address book. Existing contacts are matched by `external_id`, phone or e-mail; only
their empty fields are filled, so re-importing the same file never creates duplicates.

| Source | How |
|--------|-----|
| **Excel** `.xlsx` | first sheet; header row (`Name / Ady / Имя`, `Phone / Telefon`, `Email`, `WhatsApp`, `Telegram`, `Slack`, `External ID`, `Tags / Topar`, `Locale`) or header-less rows classified by shape (phones, e-mails, Slack `U…`/`C…` ids, Telegram chat ids; the rest is the name) |
| **CSV / TSV / TXT** | separator auto-detected (`,` `;` tab `|`); free lines like `Aman Amanow, +99365123456` work too |
| **Word** `.docx` | table rows (header or free cells) and paragraphs, one contact per row / line |
| **vCard** `.vcf` | 2.1 / 3.0 / 4.0 exports from Google, iPhone, Outlook, Android: FN/N, TEL (first → phone, second → WhatsApp), EMAIL, CATEGORIES → tags, UID → `external_id` |
| **Google Contacts** | OAuth 2.0 consent (read-only People API); tokens are used once and never stored |
| **Apple / iCloud Contacts** | CardDAV with the Apple ID and an **app-specific password**; also works for Nextcloud, Fastmail, Zimbra… |

Endpoints (admin, role developer): `POST /api/admin/projects/{id}/contacts/import` (multipart
`file`, optional `group_id`, `dry_run=true`, `tag`), `POST …/contacts/import/carddav`
(`{server_url, username, password, group_id, dry_run}`), `GET …/contacts/import/google/url`
→ consent URL, callback `GET /api/admin/integrations/google/callback`. Public API:
`POST /api/v1/contacts/import` (API key, scope `contacts`). Max 10 MB and 5000 rows per import.
Every import is written to the audit log (`contacts.import`).

## Google Contacts setup (once per server)

1. Google Cloud console → APIs & Services → enable **People API**.
2. OAuth consent screen (external, add the admins as test users while unpublished).
3. Credentials → **OAuth client ID** (Web application); authorised redirect URI:
   `https://<HABARCHY_PUBLIC_URL>/api/admin/integrations/google/callback`.
4. Set `HABARCHY_GOOGLE_CLIENT_ID` and `HABARCHY_GOOGLE_CLIENT_SECRET` (deploy/.env), restart the api.

In the admin panel or the app: Contacts → Import → **Google**. The browser opens the consent
screen; after approval Habarchy imports the connections and returns to the contacts page (web) or
shows a result page (mobile).

## Apple Contacts (iCloud)

Apple has no OAuth API for contacts; iCloud exposes them over CardDAV. On appleid.apple.com →
Sign-In and Security → **App-Specific Passwords** → generate one (two-factor authentication must be
on). In Habarchy: Contacts → Import → **Apple / iCloud**, enter the Apple ID e-mail and that
password. Habarchy discovers the address books (`contacts.icloud.com` → `pNN-contacts.icloud.com`)
and imports every card. The password is sent once over HTTPS and not stored.


## Deleting contacts (soft delete)

Deleting a contact does **not** erase it: the row is kept with a `deleted_at`
timestamp and hidden from every normal query (lists, search, group members,
broadcasts, API lookups), so the project keeps it as a record and can bring it
back. Nothing removes it automatically.

- **Trash** — `GET /api/admin/projects/{id}/contacts/trash` lists soft-deleted
  contacts (admin panel: Contacts → Trash tab; app: Contacts → trash icon).
- **Restore** — `POST …/contacts/{cid}/restore` brings one back (also back into
  its groups). It fails with 409 if another active contact now uses the same
  `external_id`.
- **Purge** — `DELETE …/contacts/{cid}/purge` removes one permanently;
  `DELETE …/contacts/trash` empties the trash. Only an admin/developer can purge.

Re-importing the same `external_id` or phone while the original is in the trash
creates a fresh active contact; the trashed one stays untouched. Soft-deleted
contacts never receive messages and never count toward a group.
