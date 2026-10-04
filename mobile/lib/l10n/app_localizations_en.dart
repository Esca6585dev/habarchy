// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Habarchy';

  @override
  String get signIn => 'Sign in';

  @override
  String get email => 'Email';

  @override
  String get password => 'Password';

  @override
  String get totpCode => 'Authenticator code';

  @override
  String get totpRequired =>
      'Two-factor authentication is enabled. Enter the 6-digit code.';

  @override
  String get invalidCredentials => 'Invalid email or password';

  @override
  String get serverUrl => 'Server';

  @override
  String get dashboard => 'Dashboard';

  @override
  String get messages => 'Messages';

  @override
  String get templates => 'Templates';

  @override
  String get providers => 'Providers';

  @override
  String get apiKeys => 'API keys';

  @override
  String get settings => 'Settings';

  @override
  String get logout => 'Sign out';

  @override
  String get projects => 'Projects';

  @override
  String get noProjects => 'No projects yet. Create one in the web admin.';

  @override
  String get today => 'Today';

  @override
  String get total => 'Total';

  @override
  String get sent => 'Sent';

  @override
  String get delivered => 'Delivered';

  @override
  String get failed => 'Failed';

  @override
  String get pending => 'Pending';

  @override
  String get cost => 'Cost';

  @override
  String get perDay => 'Messages per day';

  @override
  String get recentFailures => 'Recent failures';

  @override
  String get nothingHere => 'Nothing here yet';

  @override
  String get offlineCached => 'Offline — showing the last loaded data';

  @override
  String get search => 'Search';

  @override
  String get all => 'All';

  @override
  String get status => 'Status';

  @override
  String get channel => 'Channel';

  @override
  String get resend => 'Resend';

  @override
  String get cancel => 'Cancel';

  @override
  String get timeline => 'Timeline';

  @override
  String get rawResponse => 'Raw provider response';

  @override
  String get attempts => 'Attempts';

  @override
  String get provider => 'Provider';

  @override
  String get preview => 'Preview';

  @override
  String get sampleData => 'Sample data (JSON)';

  @override
  String get testSend => 'Test send';

  @override
  String get recipient => 'Recipient';

  @override
  String get send => 'Send';

  @override
  String get revoke => 'Revoke';

  @override
  String get revoked => 'Revoked';

  @override
  String get live => 'live';

  @override
  String get test => 'test';

  @override
  String get language => 'Language';

  @override
  String get theme => 'Theme';

  @override
  String get themeSystem => 'System';

  @override
  String get themeLight => 'Light';

  @override
  String get themeDark => 'Dark';

  @override
  String get pushDemo => 'Push demo mode';

  @override
  String get pushDemoHint =>
      'Registers this device\'s FCM token with a project (POST /api/v1/devices) so you can receive test pushes here.';

  @override
  String get pushNotConfigured =>
      'Firebase is not configured for this build. Run `flutterfire configure` and rebuild.';

  @override
  String get apiKeyForPush => 'Project API key (devices scope)';

  @override
  String get externalId => 'External ID (contact)';

  @override
  String get registerDevice => 'Register device';

  @override
  String get deviceRegistered => 'Device registered';

  @override
  String get fcmToken => 'FCM token';

  @override
  String get incomingPushes => 'Incoming pushes';

  @override
  String get noPushesYet => 'No pushes received yet';

  @override
  String get copy => 'Copy';

  @override
  String get copied => 'Copied';

  @override
  String get error => 'Something went wrong';

  @override
  String get retry => 'Retry';

  @override
  String get loading => 'Loading…';

  @override
  String version(String version) {
    return 'Version $version';
  }

  @override
  String get statusQueued => 'Queued';

  @override
  String get statusProcessing => 'Processing';

  @override
  String get statusSent => 'Sent';

  @override
  String get statusDelivered => 'Delivered';

  @override
  String get statusFailed => 'Failed';

  @override
  String get statusCancelled => 'Cancelled';

  @override
  String get healthy => 'Healthy';

  @override
  String get degraded => 'Degraded';

  @override
  String get failing => 'Failing';

  @override
  String get idle => 'Idle';

  @override
  String get disabled => 'Disabled';

  @override
  String get queues => 'Queues';

  @override
  String get systemHealth => 'System health';

  @override
  String get switchProject => 'Switch project';

  @override
  String get role => 'Role';

  @override
  String get confirmResend => 'Queue this message again?';

  @override
  String get yes => 'Yes';

  @override
  String get no => 'No';

  @override
  String get requiredVars => 'Variables';
}
