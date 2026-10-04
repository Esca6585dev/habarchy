/// Build-time configuration.
///
/// Pass values with `--dart-define`, e.g.
/// `flutter build apk --release --flavor prod --dart-define=API_URL=https://habarchy.example.tm`.
class AppConfig {
  AppConfig._();

  /// Backend base URL. Can be overridden at runtime from the login screen.
  static const defaultApiUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'http://10.0.2.2:8080', // Android emulator -> host machine
  );

  /// dev | prod. Only affects labels and logging verbosity.
  static const flavor = String.fromEnvironment('FLAVOR', defaultValue: 'dev');

  static bool get isDev => flavor != 'prod';
}
