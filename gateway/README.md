# Habarçy Gateway (`gateway/`) — telefon SMS şlýuzy

**TK:** SIM-kartaly Android telefony Habarçy üçin SMS iberiji (provider) edýän programma.
Serwerde `android_sms` görnüşli provider döredýärsiňiz, telefona şu APK-ny gurnaýarsyňyz,
API URL bilen gateway açaryny girizip **Işe gir** basýarsyňyz. Şondan soň
`POST /api/v1/messages` arkaly iberilen SMS-ler telefonyň SIM-kartasyndan gidýär,
iberildi / gowşuryldy ýagdaýy serwere gaýdýar, telefona gelen SMS-ler islege görä serwere
ugradylýar.

**EN:** Turns an Android phone with a SIM card into an SMS provider for Habarchy. The app
runs as a foreground service: it long-polls the server's outbox, sends each SMS with
`SmsManager`, reports sent / failed / delivered, sends a heartbeat every minute and can
forward received SMS.

```
gateway/
  lib/main.dart              single-screen Flutter UI (tk/ru/en): URL, key, start/stop, counters, log
  lib/gateway_channel.dart   MethodChannel wrapper (habarchy/gateway)
  android/.../gateway/
    GatewayService.kt        foreground service: heartbeat + long poll + send loop, wake lock
    SmsSender.kt             SmsManager, SIM slot selection, multipart, SENT/DELIVERED intents
    SmsResultReceiver.kt     collects part results → POST /outbox/{id}/result
    Reporter.kt              background HTTP reporting with retry
    SmsInboundReceiver.kt    SMS_RECEIVED → POST /inbound (optional)
    BootReceiver.kt          restarts the service after reboot / app update
    Api.kt, GatewayStore.kt  HTTP client (no deps), SharedPreferences settings + log
```

## 1. Server side

1. Admin panel → **Providers → Add provider**, type `android_sms`, name e.g. *Office phone*.
   Leave `gateway_key` empty: the server generates one. Optional: `sim_slot` (0/1 for dual SIM),
   `timeout_sec` (default 45), `rate_limit_per_sec` (e.g. 1), `priority`.
2. Click the provider → **Pairing** shows the API URL, the gateway key and a
   `habarchy://gateway?url=…&key=…` payload (make a QR from it if you like).
3. `HABARCHY_PUBLIC_URL` must be the address the phone can reach (https in production).

## 2. Phone side

```sh
cd gateway
flutter pub get
flutter build apk --release            # build/app/outputs/flutter-apk/app-release.apk
# or run on a connected phone
flutter run --release
```

Install the APK (sideload; SMS-sending apps are not accepted on Google Play), open it,
paste the API URL and gateway key, press **Start**, grant the SMS permission and accept the
battery-optimisation exemption when asked. The notification *Işleýär · <provider>* means the
phone is connected; the admin **Health** page shows the provider as *healthy/idle* with battery,
network and operator, or *offline* when no heartbeat arrived for 90 s.

Tips:

- Keep the phone on a charger with Wi-Fi or mobile data; the service holds a partial wake lock.
- Dual SIM: set `sim_slot` in the provider credentials; the app needs *Phone* permission to
  pick the SIM (otherwise the default SMS SIM is used).
- Some vendors (Xiaomi, Huawei, Samsung) kill background apps aggressively: enable *Autostart*
  and set battery to *No restrictions* for Habarçy Gateway.
- Operators cap SMS per hour per SIM; set `rate_limit_per_sec` accordingly and keep an
  `http_sms`/`smpp` provider as fallback (higher `priority` number).
- **Test SMS** in the app sends directly from the phone without the server (checks the SIM).
- **Received SMS** (*Forward received SMS to the server*): off by default. Turning it on asks for
  the SMS read permission and enables the SMS receiver at the OS level; the app also reads the inbox
  for SMS that arrived while it was killed, so none are lost (nothing older than the moment you
  switched it on is ever sent). Turning it off, or setting `inbound_enabled: false` on the provider
  in the admin panel, disables the receiver component completely, so no background work happens for
  incoming SMS. Each SMS is forwarded once (deduplicated) and is visible in the admin panel under
  *Inbound SMS*, via `GET /api/v1/inbound`, and as the `sms.inbound` webhook.

Release signing is the same as the admin app (`android/key.properties`, git-ignored). The app
has no Firebase or third-party dependencies.

## 3. Protocol (for other implementations)

All requests carry `X-Gateway-Key: <key>`.

| Call | Purpose |
|------|---------|
| `GET /api/gateway/v1/me` | provider name, project id, pending count, `inbound_enabled`, `sim_slot` |
| `GET /api/gateway/v1/outbox?wait=20&limit=5` | long poll; returns `[{id, message_id, to, text, sim_slot, expires_at}]` and marks them leased |
| `POST /api/gateway/v1/outbox/{id}/result` | `{status: sent\|failed\|delivered, error_code?, error_message?, parts?}` |
| `POST /api/gateway/v1/heartbeat` | any JSON (`battery`, `network`, `operator`, `model`, `app_version`) |
| `POST /api/gateway/v1/inbound` | `{from, text, received_at}` → stored, `GET /api/v1/inbound`, `sms.inbound` webhook; 403 when the provider has `inbound_enabled: false` |

`error_code` values map to Android `SmsManager` results: `generic_failure`, `radio_off`,
`null_pdu`, `no_service`, `limit_exceeded`, `short_code_not_allowed`, `network_reject`,
`invalid_number`. Items whose `expires_at` passed must not be sent.
