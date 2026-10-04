// Placeholder Firebase configuration.
//
// Run `dart pub global activate flutterfire_cli && flutterfire configure`
// in mobile/ to generate the real file (it overwrites this one), and add
// android/app/google-services.json + ios/Runner/GoogleService-Info.plist.
// With the placeholder values the app runs normally; only push demo mode
// is disabled.
import 'package:firebase_core/firebase_core.dart' show FirebaseOptions;

class DefaultFirebaseOptions {
  static const _placeholder = 'REPLACE_ME';

  static bool get isConfigured => currentPlatform.apiKey != _placeholder;

  static FirebaseOptions get currentPlatform => const FirebaseOptions(
        apiKey: _placeholder,
        appId: '1:000000000000:android:0000000000000000',
        messagingSenderId: '000000000000',
        projectId: 'habarchy-demo',
      );
}
