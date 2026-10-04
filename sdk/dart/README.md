# Habarchy SDK — Dart / Flutter

```yaml
dependencies:
  habarchy:
    git: { url: https://github.com/Esca6585dev/habarchy.git, path: sdk/dart }
```

```dart
final hb = HabarchyClient('https://habarchy.example.tm', apiKey, sign: true);
final acc = await hb.sendMessage(channel: 'sms', to: const Recipient.address('+99365123456'), template: 'otp', data: {'code': '4821', 'minutes': 5});
final detail = await hb.getMessage(acc.id);                     // detail.message.status
await hb.registerDevice(token: fcmToken, platform: 'android', externalId: userId);   // push
```

Methods: `sendMessage`, `sendBatch`, `getMessage`, `cancelMessage`, `getBatch`, `sendOtp`,
`verifyOtp`, `registerDevice`. Errors: `HabarchyException(status, code, message, details)`.
`sign: true` adds the HMAC headers (see docs/auth.md). Only `package:crypto` is required.

In a Flutter app the API key belongs to **your backend**, not the app — except for
`registerDevice`, for which you can issue a key with only the `devices` scope.

```sh
dart pub get && dart test
```
