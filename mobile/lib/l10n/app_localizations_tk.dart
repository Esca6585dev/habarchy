// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Turkmen (`tk`).
class AppLocalizationsTk extends AppLocalizations {
  AppLocalizationsTk([String locale = 'tk']) : super(locale);

  @override
  String get appTitle => 'Habarçy';

  @override
  String get signIn => 'Gir';

  @override
  String get email => 'E-mail';

  @override
  String get password => 'Parol';

  @override
  String get totpCode => 'Autentifikator kody';

  @override
  String get totpRequired =>
      'Iki basgançakly tassyklama açyk. 6 sanly kody giriziň.';

  @override
  String get invalidCredentials => 'E-mail ýa-da parol nädogry';

  @override
  String get serverUrl => 'Serwer';

  @override
  String get dashboard => 'Panel';

  @override
  String get messages => 'Habarlar';

  @override
  String get templates => 'Şablonlar';

  @override
  String get providers => 'Providerler';

  @override
  String get apiKeys => 'API açarlary';

  @override
  String get settings => 'Sazlamalar';

  @override
  String get logout => 'Çykmak';

  @override
  String get projects => 'Proýektler';

  @override
  String get noProjects => 'Entek proýekt ýok. Web adminde dörediň.';

  @override
  String get today => 'Şu gün';

  @override
  String get total => 'Jemi';

  @override
  String get sent => 'Iberildi';

  @override
  String get delivered => 'Gowşuryldy';

  @override
  String get failed => 'Şowsuz';

  @override
  String get pending => 'Garaşýar';

  @override
  String get cost => 'Çykdajy';

  @override
  String get perDay => 'Günde habarlar';

  @override
  String get recentFailures => 'Soňky şowsuzlyklar';

  @override
  String get nothingHere => 'Häzirlikçe boş';

  @override
  String get offlineCached => 'Oflaýn — soňky ýüklenen maglumat görkezilýär';

  @override
  String get search => 'Gözle';

  @override
  String get all => 'Hemmesi';

  @override
  String get status => 'Ýagdaý';

  @override
  String get channel => 'Kanal';

  @override
  String get resend => 'Gaýtadan iber';

  @override
  String get cancel => 'Ýatyr';

  @override
  String get timeline => 'Wakalar';

  @override
  String get rawResponse => 'Provideriň çig jogaby';

  @override
  String get attempts => 'Synanyşyklar';

  @override
  String get provider => 'Provider';

  @override
  String get preview => 'Deslapky görnüş';

  @override
  String get sampleData => 'Nusga maglumat (JSON)';

  @override
  String get testSend => 'Synag ibermek';

  @override
  String get recipient => 'Alyjy';

  @override
  String get send => 'Iber';

  @override
  String get revoke => 'Ýatyr';

  @override
  String get revoked => 'Ýatyryldy';

  @override
  String get live => 'live';

  @override
  String get test => 'test';

  @override
  String get language => 'Dil';

  @override
  String get theme => 'Tema';

  @override
  String get themeSystem => 'Ulgam';

  @override
  String get themeLight => 'Açyk';

  @override
  String get themeDark => 'Garaňky';

  @override
  String get pushDemo => 'Push demo rejimi';

  @override
  String get pushDemoHint =>
      'Bu enjamyň FCM tokenini proýekte hasaba alýar (POST /api/v1/devices), şonda synag pushlaryny şu ýerde alyp bilersiňiz.';

  @override
  String get pushNotConfigured =>
      'Bu gurnama üçin Firebase sazlanmadyk. `flutterfire configure` işlediň we täzeden guruň.';

  @override
  String get apiKeyForPush => 'Proýektiň API açary (devices scope)';

  @override
  String get externalId => 'Daşky ID (kontakt)';

  @override
  String get registerDevice => 'Enjamy hasaba al';

  @override
  String get deviceRegistered => 'Enjam hasaba alyndy';

  @override
  String get fcmToken => 'FCM token';

  @override
  String get incomingPushes => 'Gelen pushlar';

  @override
  String get noPushesYet => 'Entek push gelmedi';

  @override
  String get copy => 'Göçür';

  @override
  String get copied => 'Göçürildi';

  @override
  String get error => 'Bir ýalňyşlyk ýüze çykdy';

  @override
  String get retry => 'Gaýtala';

  @override
  String get loading => 'Ýüklenýär…';

  @override
  String version(String version) {
    return 'Wersiýa $version';
  }

  @override
  String get statusQueued => 'Nobatda';

  @override
  String get statusProcessing => 'Işlenýär';

  @override
  String get statusSent => 'Iberildi';

  @override
  String get statusDelivered => 'Gowşuryldy';

  @override
  String get statusFailed => 'Şowsuz';

  @override
  String get statusCancelled => 'Ýatyryldy';

  @override
  String get healthy => 'Sagdyn';

  @override
  String get degraded => 'Pese gaçan';

  @override
  String get failing => 'Şowsuz';

  @override
  String get idle => 'Boş';

  @override
  String get disabled => 'Öçürilen';

  @override
  String get queues => 'Nobatlar';

  @override
  String get systemHealth => 'Ulgam ýagdaýy';

  @override
  String get switchProject => 'Proýekti çalyş';

  @override
  String get role => 'Rol';

  @override
  String get confirmResend => 'Bu habar täzeden nobata goýulsynmy?';

  @override
  String get yes => 'Hawa';

  @override
  String get no => 'Ýok';

  @override
  String get requiredVars => 'Üýtgeýjiler';
}
