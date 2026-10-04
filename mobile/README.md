# Habarçy — Flutter admin app (`mobile/`)

**TK:** Habarçy admin paneliniň Android/iOS görnüşi: dashboard, habar logy (filter, gözleg, timeline, resend/cancel), şablonlar we preview, provider saglygy + test iberme, API açarlary, sazlamalar (dil tk/ru/en, tema, API adresi). Şeýle-de Firebase push-y alýan **demo müşderi**: enjam tokenini `/api/v1/devices` arkaly bellige alyp, gelen push-lary görkezýär — öz Flutter programmaňyza goşmak üçin nusga.

**EN:** Flutter 3 client for the Habarchy admin API (Riverpod 3, go_router, dio, freezed, fl_chart) plus a push demo built on `firebase_messaging`.

Tabs: **Dashboard · Messages · Send · Groups · More** (templates, providers, contacts, API keys,
settings). *Send* composes to groups and/or typed addresses over SMS, WhatsApp, Telegram, e-mail,
push or Slack, with templates or free text and a sandbox switch. *Groups* manages contact groups
(colleagues, classmates, customers…) with a contact picker and "Name, +99365…" quick entry.
*Providers* has typed forms for every provider type (SMTP, SMPP, HTTP SMS, Telegram, WhatsApp Cloud,
Slack, FCM, Android phone gateway with pairing key).

```
lib/
  main.dart, app.dart, router.dart          # bootstrap, Material 3 theme, auth-aware routes
  core/config.dart                          # --dart-define API_URL / FLAVOR
  core/api/{api_client,admin_repository}    # dio + bearer + single-flight refresh, typed calls
  core/auth/auth_controller.dart            # login / logout / projects / selected project
  core/models/models.dart                   # freezed DTOs (Message, Template, Health, Dashboard…)
  core/push/push_service.dart               # FCM + local notifications, background handler
  core/storage/storage.dart                 # secure storage (tokens) + prefs (locale, theme, url)
  l10n/app_{tk,ru,en}.arb                   # strings, default tk; tk_fallback.dart for framework strings
  features/{login,dashboard,messages,templates,providers,api_keys,settings,shell}
test/                                        # models, api client (refresh flow), login widget test
```

## Run

```sh
cd mobile
flutter pub get
dart run build_runner build --delete-conflicting-outputs   # freezed / json_serializable
flutter gen-l10n                                            # ARB → lib/l10n/*.dart (also runs on build)

# Android emulator → API on the host machine
flutter run --flavor dev --dart-define=API_URL=http://10.0.2.2:8080
# physical phone on the same Wi-Fi
flutter run --flavor dev --dart-define=API_URL=http://192.168.1.20:8080
# iOS simulator
flutter run --dart-define=API_URL=http://localhost:8080

flutter analyze && flutter test
```

The API address can also be changed inside the app (Settings → API URL); it is stored in shared preferences and overrides `API_URL`.

## Build

| Target | Command |
|--------|---------|
| Android dev APK (cleartext http allowed, id `tm.habarchy.admin.dev`) | `flutter build apk --release --flavor dev --dart-define=API_URL=http://10.0.2.2:8080` |
| Android prod APK (https only) | `flutter build apk --release --flavor prod --dart-define=API_URL=https://habarchy.example.tm --dart-define=FLAVOR=prod` |
| Android App Bundle | `flutter build appbundle --release --flavor prod --dart-define=API_URL=https://…` |
| iOS | `flutter build ipa --release --dart-define=API_URL=https://…` (flavors are Android-only; use Xcode schemes if you need them on iOS) |

Output: `build/app/outputs/flutter-apk/app-<flavor>-release.apk`, `build/app/outputs/bundle/<flavor>Release/`.

### Release signing (Android)

```sh
keytool -genkey -v -keystore android/keystore/habarchy-release.jks -keyalg RSA -keysize 2048 -validity 10000 -alias habarchy
cat > android/key.properties <<'EOT'
storeFile=../keystore/habarchy-release.jks
storePassword=CHANGE_ME
keyAlias=habarchy
keyPassword=CHANGE_ME
EOT
```

Both files are git-ignored. Without `key.properties` release builds fall back to the debug key so `flutter run --release` still works.

## Push (Firebase Cloud Messaging)

The app runs without Firebase; push demo is disabled until it is configured.

1. Create a Firebase project and add the Android app `tm.habarchy.admin` (+ `tm.habarchy.admin.dev`) and the iOS bundle.
2. `dart pub global activate flutterfire_cli && flutterfire configure` — writes `lib/firebase_options.dart`, `android/app/google-services.json`, `ios/Runner/GoogleService-Info.plist` (the last two are git-ignored).
3. Uncomment `id("com.google.gms.google-services")` in `android/app/build.gradle.kts` and add the classpath to `android/settings.gradle.kts` as flutterfire prints.
4. iOS: enable *Push Notifications* and *Background Modes → Remote notifications* in Xcode, upload the APNs key to Firebase.
5. In the Habarchy admin add an `fcm` provider (service-account JSON) to the project, create an API key with `devices:write messages:write`.
6. In the app: Settings → Push demo → paste the API key and an `external_id` → **Register**. The token is posted to `POST /api/v1/devices`; sending `{"channel":"push","to":{"external_id":"…"}}` now reaches the phone. Received pushes are listed on the same screen (foreground messages are shown through `flutter_local_notifications` on Android).

### Using the same code in your own Flutter app

Copy `lib/core/push/push_service.dart` and `registerDevice` from `lib/core/api/admin_repository.dart`, or use `sdk/dart` (`HabarchyClient.registerDevice`). Call it after login with your user id as `external_id`, and again whenever `FirebaseMessaging.instance.onTokenRefresh` fires.

## Localisation

Strings live in `lib/l10n/app_{tk,ru,en}.arb`; `flutter gen-l10n` regenerates `lib/l10n/app_localizations*.dart`. Flutter has no Material strings for Turkmen, so `lib/l10n/tk_fallback.dart` serves the framework's Russian strings for `tk` (dialog buttons, date pickers); every app string is translated.
