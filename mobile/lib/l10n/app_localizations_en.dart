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
  String get serverUrl => 'CardDAV server';

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

  @override
  String get groups => 'Groups';

  @override
  String get contacts => 'Contacts';

  @override
  String get compose => 'Send';

  @override
  String get more => 'More';

  @override
  String get newGroup => 'New group';

  @override
  String get groupName => 'Group name';

  @override
  String get description => 'Description';

  @override
  String get members => 'Members';

  @override
  String get addMembers => 'Add members';

  @override
  String get pickContacts => 'Existing contacts';

  @override
  String get newContacts => 'New contacts';

  @override
  String get newContactsHint => 'One per line: Name, +99365…, mail@…';

  @override
  String get remove => 'Remove';

  @override
  String get sendToGroup => 'Send to group';

  @override
  String get name => 'Name';

  @override
  String get phone => 'Phone';

  @override
  String get whatsapp => 'WhatsApp';

  @override
  String get telegram => 'Telegram chat id';

  @override
  String get slack => 'Slack id';

  @override
  String get newContact => 'New contact';

  @override
  String get editContact => 'Edit contact';

  @override
  String get delete => 'Delete';

  @override
  String confirmDelete(String name) {
    return 'Delete “$name”?';
  }

  @override
  String get save => 'Save';

  @override
  String get template => 'Template';

  @override
  String get freeText => 'Free text';

  @override
  String get subject => 'Subject / title';

  @override
  String get body => 'Message';

  @override
  String get templateData => 'Template data (JSON)';

  @override
  String get sandboxMode => 'Test mode (sandbox, nothing is really sent)';

  @override
  String get recipients => 'Recipients';

  @override
  String get addresses => 'Addresses, one per line';

  @override
  String get addressesHint =>
      'Phones, e-mails or chat ids depending on the channel';

  @override
  String queuedRejected(int n, int m) {
    return '$n queued · $m rejected';
  }

  @override
  String get selectRecipients => 'Choose a group or type an address';

  @override
  String get noGroupsYet =>
      'No groups yet. Create “Colleagues”, “Customers”… and add members.';

  @override
  String added(int n, int created) {
    return '$n added · $created new contacts';
  }

  @override
  String get newProvider => 'New provider';

  @override
  String get editProvider => 'Edit provider';

  @override
  String get providerType => 'Type';

  @override
  String get priority => 'Priority (lower first)';

  @override
  String get rateLimit => 'Rate limit (msg/s, 0 = off)';

  @override
  String get active => 'Active';

  @override
  String get credentials => 'Credentials';

  @override
  String get keepCredentials => 'Leave empty to keep the current credentials';

  @override
  String get pairing => 'Pair phone';

  @override
  String get pairingHint =>
      'Install the Habarçy Gateway APK on a phone with a SIM card, enter this URL and key, press Start.';

  @override
  String get gatewayKey => 'Gateway key';

  @override
  String get apiUrl => 'API URL';

  @override
  String get online => 'Online';

  @override
  String get offline => 'Offline';

  @override
  String get whatsappHint =>
      'Free text works inside the 24 h window; otherwise an approved template is needed.';

  @override
  String get contactHint =>
      'At least one of phone, e-mail, WhatsApp, Telegram or Slack is required.';

  @override
  String get openMessages => 'Open message log';

  @override
  String get noMembers => 'No members yet';

  @override
  String get inboundSms => 'Inbound SMS';

  @override
  String get inboundHint =>
      'SMS received by gateway phones appear here when forwarding is on (provider setting + switch in the phone app).';

  @override
  String get importTitle => 'Import contacts';

  @override
  String get importFile => 'From a file';

  @override
  String get importFileHint =>
      'Excel, CSV, Word or vCard (.vcf). Existing contacts are merged, never duplicated.';

  @override
  String get importGoogle => 'Google Contacts';

  @override
  String get importGoogleHint =>
      'Opens Google’s consent screen; read-only, tokens are not stored.';

  @override
  String get importApple => 'Apple / iCloud Contacts';

  @override
  String get importAppleHint =>
      'Apple ID + app-specific password (appleid.apple.com). Used once, not stored.';

  @override
  String get appleId => 'Apple ID (e-mail)';

  @override
  String get appPassword => 'App-specific password';

  @override
  String get addToGroup => 'Add to group';

  @override
  String get noGroup => 'No group';

  @override
  String get importNow => 'Import';

  @override
  String importResult(int total, int created, int updated, int skipped) {
    return '$total rows: $created created, $updated updated, $skipped skipped';
  }

  @override
  String get googleNotConfigured =>
      'Google import is not configured on this server';
}
