import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_ru.dart';
import 'app_localizations_tk.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('en'),
    Locale('ru'),
    Locale('tk'),
  ];

  /// No description provided for @appTitle.
  ///
  /// In en, this message translates to:
  /// **'Habarchy'**
  String get appTitle;

  /// No description provided for @signIn.
  ///
  /// In en, this message translates to:
  /// **'Sign in'**
  String get signIn;

  /// No description provided for @email.
  ///
  /// In en, this message translates to:
  /// **'Email'**
  String get email;

  /// No description provided for @password.
  ///
  /// In en, this message translates to:
  /// **'Password'**
  String get password;

  /// No description provided for @totpCode.
  ///
  /// In en, this message translates to:
  /// **'Authenticator code'**
  String get totpCode;

  /// No description provided for @totpRequired.
  ///
  /// In en, this message translates to:
  /// **'Two-factor authentication is enabled. Enter the 6-digit code.'**
  String get totpRequired;

  /// No description provided for @invalidCredentials.
  ///
  /// In en, this message translates to:
  /// **'Invalid email or password'**
  String get invalidCredentials;

  /// No description provided for @serverUrl.
  ///
  /// In en, this message translates to:
  /// **'CardDAV server'**
  String get serverUrl;

  /// No description provided for @dashboard.
  ///
  /// In en, this message translates to:
  /// **'Dashboard'**
  String get dashboard;

  /// No description provided for @messages.
  ///
  /// In en, this message translates to:
  /// **'Messages'**
  String get messages;

  /// No description provided for @templates.
  ///
  /// In en, this message translates to:
  /// **'Templates'**
  String get templates;

  /// No description provided for @providers.
  ///
  /// In en, this message translates to:
  /// **'Providers'**
  String get providers;

  /// No description provided for @apiKeys.
  ///
  /// In en, this message translates to:
  /// **'API keys'**
  String get apiKeys;

  /// No description provided for @settings.
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get settings;

  /// No description provided for @logout.
  ///
  /// In en, this message translates to:
  /// **'Sign out'**
  String get logout;

  /// No description provided for @projects.
  ///
  /// In en, this message translates to:
  /// **'Projects'**
  String get projects;

  /// No description provided for @noProjects.
  ///
  /// In en, this message translates to:
  /// **'No projects yet. Create one in the web admin.'**
  String get noProjects;

  /// No description provided for @today.
  ///
  /// In en, this message translates to:
  /// **'Today'**
  String get today;

  /// No description provided for @total.
  ///
  /// In en, this message translates to:
  /// **'Total'**
  String get total;

  /// No description provided for @sent.
  ///
  /// In en, this message translates to:
  /// **'Sent'**
  String get sent;

  /// No description provided for @delivered.
  ///
  /// In en, this message translates to:
  /// **'Delivered'**
  String get delivered;

  /// No description provided for @failed.
  ///
  /// In en, this message translates to:
  /// **'Failed'**
  String get failed;

  /// No description provided for @pending.
  ///
  /// In en, this message translates to:
  /// **'Pending'**
  String get pending;

  /// No description provided for @cost.
  ///
  /// In en, this message translates to:
  /// **'Cost'**
  String get cost;

  /// No description provided for @perDay.
  ///
  /// In en, this message translates to:
  /// **'Messages per day'**
  String get perDay;

  /// No description provided for @recentFailures.
  ///
  /// In en, this message translates to:
  /// **'Recent failures'**
  String get recentFailures;

  /// No description provided for @nothingHere.
  ///
  /// In en, this message translates to:
  /// **'Nothing here yet'**
  String get nothingHere;

  /// No description provided for @offlineCached.
  ///
  /// In en, this message translates to:
  /// **'Offline — showing the last loaded data'**
  String get offlineCached;

  /// No description provided for @search.
  ///
  /// In en, this message translates to:
  /// **'Search'**
  String get search;

  /// No description provided for @all.
  ///
  /// In en, this message translates to:
  /// **'All'**
  String get all;

  /// No description provided for @status.
  ///
  /// In en, this message translates to:
  /// **'Status'**
  String get status;

  /// No description provided for @channel.
  ///
  /// In en, this message translates to:
  /// **'Channel'**
  String get channel;

  /// No description provided for @resend.
  ///
  /// In en, this message translates to:
  /// **'Resend'**
  String get resend;

  /// No description provided for @cancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get cancel;

  /// No description provided for @timeline.
  ///
  /// In en, this message translates to:
  /// **'Timeline'**
  String get timeline;

  /// No description provided for @rawResponse.
  ///
  /// In en, this message translates to:
  /// **'Raw provider response'**
  String get rawResponse;

  /// No description provided for @attempts.
  ///
  /// In en, this message translates to:
  /// **'Attempts'**
  String get attempts;

  /// No description provided for @provider.
  ///
  /// In en, this message translates to:
  /// **'Provider'**
  String get provider;

  /// No description provided for @preview.
  ///
  /// In en, this message translates to:
  /// **'Preview'**
  String get preview;

  /// No description provided for @sampleData.
  ///
  /// In en, this message translates to:
  /// **'Sample data (JSON)'**
  String get sampleData;

  /// No description provided for @testSend.
  ///
  /// In en, this message translates to:
  /// **'Test send'**
  String get testSend;

  /// No description provided for @recipient.
  ///
  /// In en, this message translates to:
  /// **'Recipient'**
  String get recipient;

  /// No description provided for @send.
  ///
  /// In en, this message translates to:
  /// **'Send'**
  String get send;

  /// No description provided for @revoke.
  ///
  /// In en, this message translates to:
  /// **'Revoke'**
  String get revoke;

  /// No description provided for @revoked.
  ///
  /// In en, this message translates to:
  /// **'Revoked'**
  String get revoked;

  /// No description provided for @live.
  ///
  /// In en, this message translates to:
  /// **'live'**
  String get live;

  /// No description provided for @test.
  ///
  /// In en, this message translates to:
  /// **'test'**
  String get test;

  /// No description provided for @language.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get language;

  /// No description provided for @theme.
  ///
  /// In en, this message translates to:
  /// **'Theme'**
  String get theme;

  /// No description provided for @themeSystem.
  ///
  /// In en, this message translates to:
  /// **'System'**
  String get themeSystem;

  /// No description provided for @themeLight.
  ///
  /// In en, this message translates to:
  /// **'Light'**
  String get themeLight;

  /// No description provided for @themeDark.
  ///
  /// In en, this message translates to:
  /// **'Dark'**
  String get themeDark;

  /// No description provided for @pushDemo.
  ///
  /// In en, this message translates to:
  /// **'Push demo mode'**
  String get pushDemo;

  /// No description provided for @pushDemoHint.
  ///
  /// In en, this message translates to:
  /// **'Registers this device\'s FCM token with a project (POST /api/v1/devices) so you can receive test pushes here.'**
  String get pushDemoHint;

  /// No description provided for @pushNotConfigured.
  ///
  /// In en, this message translates to:
  /// **'Firebase is not configured for this build. Run `flutterfire configure` and rebuild.'**
  String get pushNotConfigured;

  /// No description provided for @apiKeyForPush.
  ///
  /// In en, this message translates to:
  /// **'Project API key (devices scope)'**
  String get apiKeyForPush;

  /// No description provided for @externalId.
  ///
  /// In en, this message translates to:
  /// **'External ID (contact)'**
  String get externalId;

  /// No description provided for @registerDevice.
  ///
  /// In en, this message translates to:
  /// **'Register device'**
  String get registerDevice;

  /// No description provided for @deviceRegistered.
  ///
  /// In en, this message translates to:
  /// **'Device registered'**
  String get deviceRegistered;

  /// No description provided for @fcmToken.
  ///
  /// In en, this message translates to:
  /// **'FCM token'**
  String get fcmToken;

  /// No description provided for @incomingPushes.
  ///
  /// In en, this message translates to:
  /// **'Incoming pushes'**
  String get incomingPushes;

  /// No description provided for @noPushesYet.
  ///
  /// In en, this message translates to:
  /// **'No pushes received yet'**
  String get noPushesYet;

  /// No description provided for @copy.
  ///
  /// In en, this message translates to:
  /// **'Copy'**
  String get copy;

  /// No description provided for @copied.
  ///
  /// In en, this message translates to:
  /// **'Copied'**
  String get copied;

  /// No description provided for @error.
  ///
  /// In en, this message translates to:
  /// **'Something went wrong'**
  String get error;

  /// No description provided for @retry.
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get retry;

  /// No description provided for @loading.
  ///
  /// In en, this message translates to:
  /// **'Loading…'**
  String get loading;

  /// No description provided for @version.
  ///
  /// In en, this message translates to:
  /// **'Version {version}'**
  String version(String version);

  /// No description provided for @statusQueued.
  ///
  /// In en, this message translates to:
  /// **'Queued'**
  String get statusQueued;

  /// No description provided for @statusProcessing.
  ///
  /// In en, this message translates to:
  /// **'Processing'**
  String get statusProcessing;

  /// No description provided for @statusSent.
  ///
  /// In en, this message translates to:
  /// **'Sent'**
  String get statusSent;

  /// No description provided for @statusDelivered.
  ///
  /// In en, this message translates to:
  /// **'Delivered'**
  String get statusDelivered;

  /// No description provided for @statusFailed.
  ///
  /// In en, this message translates to:
  /// **'Failed'**
  String get statusFailed;

  /// No description provided for @statusCancelled.
  ///
  /// In en, this message translates to:
  /// **'Cancelled'**
  String get statusCancelled;

  /// No description provided for @healthy.
  ///
  /// In en, this message translates to:
  /// **'Healthy'**
  String get healthy;

  /// No description provided for @degraded.
  ///
  /// In en, this message translates to:
  /// **'Degraded'**
  String get degraded;

  /// No description provided for @failing.
  ///
  /// In en, this message translates to:
  /// **'Failing'**
  String get failing;

  /// No description provided for @idle.
  ///
  /// In en, this message translates to:
  /// **'Idle'**
  String get idle;

  /// No description provided for @disabled.
  ///
  /// In en, this message translates to:
  /// **'Disabled'**
  String get disabled;

  /// No description provided for @queues.
  ///
  /// In en, this message translates to:
  /// **'Queues'**
  String get queues;

  /// No description provided for @systemHealth.
  ///
  /// In en, this message translates to:
  /// **'System health'**
  String get systemHealth;

  /// No description provided for @switchProject.
  ///
  /// In en, this message translates to:
  /// **'Switch project'**
  String get switchProject;

  /// No description provided for @role.
  ///
  /// In en, this message translates to:
  /// **'Role'**
  String get role;

  /// No description provided for @confirmResend.
  ///
  /// In en, this message translates to:
  /// **'Queue this message again?'**
  String get confirmResend;

  /// No description provided for @yes.
  ///
  /// In en, this message translates to:
  /// **'Yes'**
  String get yes;

  /// No description provided for @no.
  ///
  /// In en, this message translates to:
  /// **'No'**
  String get no;

  /// No description provided for @requiredVars.
  ///
  /// In en, this message translates to:
  /// **'Variables'**
  String get requiredVars;

  /// No description provided for @groups.
  ///
  /// In en, this message translates to:
  /// **'Groups'**
  String get groups;

  /// No description provided for @contacts.
  ///
  /// In en, this message translates to:
  /// **'Contacts'**
  String get contacts;

  /// No description provided for @compose.
  ///
  /// In en, this message translates to:
  /// **'Send'**
  String get compose;

  /// No description provided for @more.
  ///
  /// In en, this message translates to:
  /// **'More'**
  String get more;

  /// No description provided for @newGroup.
  ///
  /// In en, this message translates to:
  /// **'New group'**
  String get newGroup;

  /// No description provided for @groupName.
  ///
  /// In en, this message translates to:
  /// **'Group name'**
  String get groupName;

  /// No description provided for @description.
  ///
  /// In en, this message translates to:
  /// **'Description'**
  String get description;

  /// No description provided for @members.
  ///
  /// In en, this message translates to:
  /// **'Members'**
  String get members;

  /// No description provided for @addMembers.
  ///
  /// In en, this message translates to:
  /// **'Add members'**
  String get addMembers;

  /// No description provided for @pickContacts.
  ///
  /// In en, this message translates to:
  /// **'Existing contacts'**
  String get pickContacts;

  /// No description provided for @newContacts.
  ///
  /// In en, this message translates to:
  /// **'New contacts'**
  String get newContacts;

  /// No description provided for @newContactsHint.
  ///
  /// In en, this message translates to:
  /// **'One per line: Name, +99365…, mail@…'**
  String get newContactsHint;

  /// No description provided for @remove.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get remove;

  /// No description provided for @sendToGroup.
  ///
  /// In en, this message translates to:
  /// **'Send to group'**
  String get sendToGroup;

  /// No description provided for @name.
  ///
  /// In en, this message translates to:
  /// **'Name'**
  String get name;

  /// No description provided for @phone.
  ///
  /// In en, this message translates to:
  /// **'Phone'**
  String get phone;

  /// No description provided for @whatsapp.
  ///
  /// In en, this message translates to:
  /// **'WhatsApp'**
  String get whatsapp;

  /// No description provided for @telegram.
  ///
  /// In en, this message translates to:
  /// **'Telegram chat id'**
  String get telegram;

  /// No description provided for @slack.
  ///
  /// In en, this message translates to:
  /// **'Slack id'**
  String get slack;

  /// No description provided for @newContact.
  ///
  /// In en, this message translates to:
  /// **'New contact'**
  String get newContact;

  /// No description provided for @editContact.
  ///
  /// In en, this message translates to:
  /// **'Edit contact'**
  String get editContact;

  /// No description provided for @delete.
  ///
  /// In en, this message translates to:
  /// **'Delete'**
  String get delete;

  /// No description provided for @confirmDelete.
  ///
  /// In en, this message translates to:
  /// **'Delete “{name}”?'**
  String confirmDelete(String name);

  /// No description provided for @save.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get save;

  /// No description provided for @template.
  ///
  /// In en, this message translates to:
  /// **'Template'**
  String get template;

  /// No description provided for @freeText.
  ///
  /// In en, this message translates to:
  /// **'Free text'**
  String get freeText;

  /// No description provided for @subject.
  ///
  /// In en, this message translates to:
  /// **'Subject / title'**
  String get subject;

  /// No description provided for @body.
  ///
  /// In en, this message translates to:
  /// **'Message'**
  String get body;

  /// No description provided for @templateData.
  ///
  /// In en, this message translates to:
  /// **'Template data (JSON)'**
  String get templateData;

  /// No description provided for @sandboxMode.
  ///
  /// In en, this message translates to:
  /// **'Test mode (sandbox, nothing is really sent)'**
  String get sandboxMode;

  /// No description provided for @recipients.
  ///
  /// In en, this message translates to:
  /// **'Recipients'**
  String get recipients;

  /// No description provided for @addresses.
  ///
  /// In en, this message translates to:
  /// **'Addresses, one per line'**
  String get addresses;

  /// No description provided for @addressesHint.
  ///
  /// In en, this message translates to:
  /// **'Phones, e-mails or chat ids depending on the channel'**
  String get addressesHint;

  /// No description provided for @queuedRejected.
  ///
  /// In en, this message translates to:
  /// **'{n} queued · {m} rejected'**
  String queuedRejected(int n, int m);

  /// No description provided for @selectRecipients.
  ///
  /// In en, this message translates to:
  /// **'Choose a group or type an address'**
  String get selectRecipients;

  /// No description provided for @noGroupsYet.
  ///
  /// In en, this message translates to:
  /// **'No groups yet. Create “Colleagues”, “Customers”… and add members.'**
  String get noGroupsYet;

  /// No description provided for @added.
  ///
  /// In en, this message translates to:
  /// **'{n} added · {created} new contacts'**
  String added(int n, int created);

  /// No description provided for @newProvider.
  ///
  /// In en, this message translates to:
  /// **'New provider'**
  String get newProvider;

  /// No description provided for @editProvider.
  ///
  /// In en, this message translates to:
  /// **'Edit provider'**
  String get editProvider;

  /// No description provided for @providerType.
  ///
  /// In en, this message translates to:
  /// **'Type'**
  String get providerType;

  /// No description provided for @priority.
  ///
  /// In en, this message translates to:
  /// **'Priority (lower first)'**
  String get priority;

  /// No description provided for @rateLimit.
  ///
  /// In en, this message translates to:
  /// **'Rate limit (msg/s, 0 = off)'**
  String get rateLimit;

  /// No description provided for @active.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get active;

  /// No description provided for @credentials.
  ///
  /// In en, this message translates to:
  /// **'Credentials'**
  String get credentials;

  /// No description provided for @keepCredentials.
  ///
  /// In en, this message translates to:
  /// **'Leave empty to keep the current credentials'**
  String get keepCredentials;

  /// No description provided for @pairing.
  ///
  /// In en, this message translates to:
  /// **'Pair phone'**
  String get pairing;

  /// No description provided for @pairingHint.
  ///
  /// In en, this message translates to:
  /// **'Install the Habarçy Gateway APK on a phone with a SIM card, enter this URL and key, press Start.'**
  String get pairingHint;

  /// No description provided for @gatewayKey.
  ///
  /// In en, this message translates to:
  /// **'Gateway key'**
  String get gatewayKey;

  /// No description provided for @apiUrl.
  ///
  /// In en, this message translates to:
  /// **'API URL'**
  String get apiUrl;

  /// No description provided for @online.
  ///
  /// In en, this message translates to:
  /// **'Online'**
  String get online;

  /// No description provided for @offline.
  ///
  /// In en, this message translates to:
  /// **'Offline'**
  String get offline;

  /// No description provided for @whatsappHint.
  ///
  /// In en, this message translates to:
  /// **'Free text works inside the 24 h window; otherwise an approved template is needed.'**
  String get whatsappHint;

  /// No description provided for @contactHint.
  ///
  /// In en, this message translates to:
  /// **'At least one of phone, e-mail, WhatsApp, Telegram or Slack is required.'**
  String get contactHint;

  /// No description provided for @openMessages.
  ///
  /// In en, this message translates to:
  /// **'Open message log'**
  String get openMessages;

  /// No description provided for @noMembers.
  ///
  /// In en, this message translates to:
  /// **'No members yet'**
  String get noMembers;

  /// No description provided for @inboundSms.
  ///
  /// In en, this message translates to:
  /// **'Inbound SMS'**
  String get inboundSms;

  /// No description provided for @inboundHint.
  ///
  /// In en, this message translates to:
  /// **'SMS received by gateway phones appear here when forwarding is on (provider setting + switch in the phone app).'**
  String get inboundHint;

  /// No description provided for @importTitle.
  ///
  /// In en, this message translates to:
  /// **'Import contacts'**
  String get importTitle;

  /// No description provided for @importFile.
  ///
  /// In en, this message translates to:
  /// **'From a file'**
  String get importFile;

  /// No description provided for @importFileHint.
  ///
  /// In en, this message translates to:
  /// **'Excel, CSV, Word or vCard (.vcf). Existing contacts are merged, never duplicated.'**
  String get importFileHint;

  /// No description provided for @importGoogle.
  ///
  /// In en, this message translates to:
  /// **'Google Contacts'**
  String get importGoogle;

  /// No description provided for @importGoogleHint.
  ///
  /// In en, this message translates to:
  /// **'Opens Google’s consent screen; read-only, tokens are not stored.'**
  String get importGoogleHint;

  /// No description provided for @importApple.
  ///
  /// In en, this message translates to:
  /// **'Apple / iCloud Contacts'**
  String get importApple;

  /// No description provided for @importAppleHint.
  ///
  /// In en, this message translates to:
  /// **'Apple ID + app-specific password (appleid.apple.com). Used once, not stored.'**
  String get importAppleHint;

  /// No description provided for @appleId.
  ///
  /// In en, this message translates to:
  /// **'Apple ID (e-mail)'**
  String get appleId;

  /// No description provided for @appPassword.
  ///
  /// In en, this message translates to:
  /// **'App-specific password'**
  String get appPassword;

  /// No description provided for @addToGroup.
  ///
  /// In en, this message translates to:
  /// **'Add to group'**
  String get addToGroup;

  /// No description provided for @noGroup.
  ///
  /// In en, this message translates to:
  /// **'No group'**
  String get noGroup;

  /// No description provided for @importNow.
  ///
  /// In en, this message translates to:
  /// **'Import'**
  String get importNow;

  /// No description provided for @importResult.
  ///
  /// In en, this message translates to:
  /// **'{total} rows: {created} created, {updated} updated, {skipped} skipped'**
  String importResult(int total, int created, int updated, int skipped);

  /// No description provided for @googleNotConfigured.
  ///
  /// In en, this message translates to:
  /// **'Google import is not configured on this server'**
  String get googleNotConfigured;

  /// No description provided for @trash.
  ///
  /// In en, this message translates to:
  /// **'Trash'**
  String get trash;

  /// No description provided for @restore.
  ///
  /// In en, this message translates to:
  /// **'Restore'**
  String get restore;

  /// No description provided for @purge.
  ///
  /// In en, this message translates to:
  /// **'Delete permanently'**
  String get purge;

  /// No description provided for @purgeAll.
  ///
  /// In en, this message translates to:
  /// **'Empty trash'**
  String get purgeAll;

  /// No description provided for @emptyTrash.
  ///
  /// In en, this message translates to:
  /// **'Trash is empty'**
  String get emptyTrash;

  /// No description provided for @softDeleteHint.
  ///
  /// In en, this message translates to:
  /// **'Deleted contacts are kept here and can be restored. They are only removed for good when you purge them.'**
  String get softDeleteHint;

  /// No description provided for @restored.
  ///
  /// In en, this message translates to:
  /// **'Restored'**
  String get restored;

  /// No description provided for @purgedAll.
  ///
  /// In en, this message translates to:
  /// **'{n} permanently deleted'**
  String purgedAll(int n);

  /// No description provided for @chat.
  ///
  /// In en, this message translates to:
  /// **'Team chat'**
  String get chat;

  /// No description provided for @chatNew.
  ///
  /// In en, this message translates to:
  /// **'New'**
  String get chatNew;

  /// No description provided for @chatDirect.
  ///
  /// In en, this message translates to:
  /// **'Direct'**
  String get chatDirect;

  /// No description provided for @chatChannel.
  ///
  /// In en, this message translates to:
  /// **'Channel'**
  String get chatChannel;

  /// No description provided for @chatPublic.
  ///
  /// In en, this message translates to:
  /// **'Public'**
  String get chatPublic;

  /// No description provided for @chatPrivate.
  ///
  /// In en, this message translates to:
  /// **'Private'**
  String get chatPrivate;

  /// No description provided for @chatName.
  ///
  /// In en, this message translates to:
  /// **'Name'**
  String get chatName;

  /// No description provided for @chatMessage.
  ///
  /// In en, this message translates to:
  /// **'Write a message…'**
  String get chatMessage;

  /// No description provided for @chatNoMessages.
  ///
  /// In en, this message translates to:
  /// **'No messages yet'**
  String get chatNoMessages;

  /// No description provided for @chatEmpty.
  ///
  /// In en, this message translates to:
  /// **'No chats yet. Start a direct message or create a channel.'**
  String get chatEmpty;

  /// No description provided for @profileTitle.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get profileTitle;

  /// No description provided for @uploadPhoto.
  ///
  /// In en, this message translates to:
  /// **'Upload photo'**
  String get uploadPhoto;

  /// No description provided for @profileSaved.
  ///
  /// In en, this message translates to:
  /// **'Profile saved'**
  String get profileSaved;

  /// No description provided for @bio.
  ///
  /// In en, this message translates to:
  /// **'About'**
  String get bio;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['en', 'ru', 'tk'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'ru':
      return AppLocalizationsRu();
    case 'tk':
      return AppLocalizationsTk();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
