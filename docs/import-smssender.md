# sms-sender nomerlerini Habarça import etmek

**TK:** Beýleki repozitoriýanyň (`sms-sender-app`) backend-indäki kontaktlary
(nomerleri) Habarçynyň bir müşderisine (tenant) geçirýän gural. Çeşme:
`https://atom.com.tm/elitepanel/sms`, `Integration-Key` header bilen,
`GET /contacts/all/` → `[{id, title, phone_number, is_sent}]`.

Her nomer E.164-e normalizlenýär (türkmen default), bir tenant-a goýulýar we
**başga müşderiniňki bilen garyşmaýar**. Gaýtadan işletseň üýtgedýär,
goşmaýar (idempotent): `external_id = "smssender:<telefon>"`.

## Işletmek

`backend/.env`-de DB sazlamalary bar bolmaly (API bilen deň). Soň:

```sh
cd backend

# 1) Ilki gury-barlag (dry run): näçe nomer, näçesi nädogry — hiç zat ýazylmaýar
make import-smssender ARGS='-project demo -key <Integration-Key> -dry-run'

# 2) Hakyky import (demo tenant-a)
make import-smssender ARGS='-project demo -key <Integration-Key>'

# ýa-da göni:
SMSSENDER_KEY=<key> go run ./cmd/import-smssender -project demo
```

Flaglar:

| flag | düşündiriş |
|--|--|
| `-project` | **hökmany** — Habarçy tenant slug (nirä import) |
| `-key` | **hökmany** — çeşme `Integration-Key` (ýa `SMSSENDER_KEY` env) |
| `-source` | çeşme base URL (default `https://atom.com.tm/elitepanel/sms`) |
| `-tag` | import edilen kontaktlara bellik (default `smssender`) |
| `-dry-run` | diňe barla, ýazma |
| `-limit` | iň köp näçe import (0 = hemmesi) |

## Netije

```
fetched 1240 contacts from https://atom.com.tm/elitepanel/sms → tenant "demo"
done: created=1198 updated=0 invalid=30 blank=12 dry_run=false
note: 30 numbers could not be normalized to a valid phone and were skipped
```

- **created** — täze goşuldy · **updated** — öňden bardy, täzelendi
- **invalid** — E.164-e öwrülmedik (nädogry) nomerler, geçirilmedi
- **blank** — boş telefon meýdany

Import edilen kontaktlaryň `attributes`-inde `source=smssender`, `is_sent`,
`source_id` saklanýar; bellik (`tag`) bilen admin panelde süzüp bolýar.

> ⚠️ Çeşme host-y (`atom.com.tm`) işleýän maşynyňyzdan elýeterli bolmaly.
> `Integration-Key` gizlin — repozitoriýa salmaň, diňe `.env`/secret-de saklaň.
