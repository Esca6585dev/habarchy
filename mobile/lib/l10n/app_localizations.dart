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
  /// **'Server'**
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
