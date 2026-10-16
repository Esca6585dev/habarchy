# scripts

Developer and test helpers for Habarchy.

## virtual_phone.py — sanal telefon (SMS gateway test guraly)

Android Gateway APK-ny gurmazdan serweriň SMS-gateway tarapyny synamak üçin.
Ol hakyky telefon ýaly jübütlenýär, outbox-y long-poll edýär, her SMS-i
"iberýär" (ekrana ýazýar), netijäni serwere hasabat berýär we heartbeat
iberýär. Diňe Python 3 standart kitaphanasyny ulanýar — gurnama gerek däl.

### Başlat

```bash
python3 scripts/virtual_phone.py \
  --url https://habarchy.example.tm \
  --key gw_xxxxxxxxxxxxxxxx
```

`gateway_key`-i admin panelden al: provaýder → **Pairing** (ýa-da
`habarchy://gateway?...&key=...` QR-dan). Ol `android_sms` tipli provaýderiň
açary.

### Doly synag (iki terminal)

1. **1-nji terminal** — sanal telefony işlet (diňläp dursun):

   ```bash
   python3 scripts/virtual_phone.py --url http://localhost:8080 --key gw_… --delivered
   ```

2. **2-nji terminal** — SMS-i nobata goý. Iň aňsaty — provaýderiň test-send
   endpointi (ýa-da adaty habar iberiş API-si):

   ```bash
   # admin JWT bilen:
   curl -X POST http://localhost:8080/api/v1/providers/<provider_id>/test \
     -H "Authorization: Bearer <JWT>" \
     -H "Content-Type: application/json" \
     -d '{"to":"+99365123456","text":"OTP kody: 4821"}'
   ```

   1-nji terminalda SMS peýda bolup, `sent` (we `--delivered` bilen
   `delivered`) hasabat berilýär. Admin panelde habar "ýetirildi" bolýar.

### Peýdaly baýdaklar

| Baýdak | Näme edýär |
|---|---|
| `--once` | Bir topar alyp işle, soň çyk (CI üçin) |
| `--fail` | Ähli SMS-i `failed` hasabat ber (fallback/retry synagy) |
| `--fail-rate 0.3` | ~30% SMS-i tötänleýin näsaz et |
| `--send-delay 2` | Her SMS-i "ibermek" üçin 2s garaş |
| `--delivered` | `sent`-den soň delivery receipt-de iber |
| `--inbound "+99365…:Salam"` | Bir galp gelen SMS ugrat we çyk |
| `--wait 20` | Her outbox soragyndaky long-poll sekundy |

Doly kömek: `python3 scripts/virtual_phone.py --help`
