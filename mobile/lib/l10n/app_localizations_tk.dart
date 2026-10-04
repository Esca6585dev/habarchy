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
  String get serverUrl => 'CardDAV serwer';

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

  @override
  String get groups => 'Toparlar';

  @override
  String get contacts => 'Kontaktlar';

  @override
  String get compose => 'Iber';

  @override
  String get more => 'Başga';

  @override
  String get newGroup => 'Täze topar';

  @override
  String get groupName => 'Toparyň ady';

  @override
  String get description => 'Düşündiriş';

  @override
  String get members => 'Agzalar';

  @override
  String get addMembers => 'Agza goş';

  @override
  String get pickContacts => 'Bar bolan kontaktlar';

  @override
  String get newContacts => 'Täze kontaktlar';

  @override
  String get newContactsHint => 'Her setirde biri: Ady, +99365…, mail@…';

  @override
  String get remove => 'Aýyr';

  @override
  String get sendToGroup => 'Topara iber';

  @override
  String get name => 'Ady';

  @override
  String get phone => 'Telefon';

  @override
  String get whatsapp => 'WhatsApp';

  @override
  String get telegram => 'Telegram chat id';

  @override
  String get slack => 'Slack id';

  @override
  String get newContact => 'Täze kontakt';

  @override
  String get editContact => 'Kontakty üýtget';

  @override
  String get delete => 'Poz';

  @override
  String confirmDelete(String name) {
    return '«$name» pozulsynmy?';
  }

  @override
  String get save => 'Ýatda sakla';

  @override
  String get template => 'Şablon';

  @override
  String get freeText => 'Erkin tekst';

  @override
  String get subject => 'Tema / sözbaşy';

  @override
  String get body => 'Habar';

  @override
  String get templateData => 'Şablon maglumatlary (JSON)';

  @override
  String get sandboxMode => 'Synag tertibi (sandbox, hakykatda iberilmeýär)';

  @override
  String get recipients => 'Alyjylar';

  @override
  String get addresses => 'Salgylar, her setirde biri';

  @override
  String get addressesHint => 'Kanala görä telefon, e-mail ýa-da chat id';

  @override
  String queuedRejected(int n, int m) {
    return '$n nobata goýuldy · $m ret edildi';
  }

  @override
  String get selectRecipients => 'Topar saýlaň ýa-da salgy ýazyň';

  @override
  String get noGroupsYet =>
      'Heniz topar ýok. «Işdeşler», «Müşderiler»… dörediň we agza goşuň.';

  @override
  String added(int n, int created) {
    return '$n goşuldy · $created täze kontakt';
  }

  @override
  String get newProvider => 'Täze provider';

  @override
  String get editProvider => 'Provideri üýtget';

  @override
  String get providerType => 'Görnüşi';

  @override
  String get priority => 'Ileri tutma (kiçisi öň)';

  @override
  String get rateLimit => 'Çäk (habar/s, 0 = ýok)';

  @override
  String get active => 'Işjeň';

  @override
  String get credentials => 'Maglumatlar';

  @override
  String get keepCredentials => 'Häzirki maglumatlary saklamak üçin boş goýuň';

  @override
  String get pairing => 'Telefony birikdir';

  @override
  String get pairingHint =>
      'SIM-kartaly telefona Habarçy Gateway APK-ny gurnaň, şu URL bilen açary giriziň, «Işe gir» basyň.';

  @override
  String get gatewayKey => 'Gateway açary';

  @override
  String get apiUrl => 'API URL';

  @override
  String get online => 'Onlaýn';

  @override
  String get offline => 'Oflaýn';

  @override
  String get whatsappHint =>
      'Erkin tekst 24 sagatlyk penjirede işleýär; başga ýagdaýda tassyklanan şablon gerek.';

  @override
  String get contactHint =>
      'Telefon, e-mail, WhatsApp, Telegram ýa-da Slack-dan iň azyndan biri gerek.';

  @override
  String get openMessages => 'Habar logyny aç';

  @override
  String get noMembers => 'Heniz agza ýok';

  @override
  String get inboundSms => 'Gelen SMS';

  @override
  String get inboundHint =>
      'Ugratmak açyk bolanda şlýuz-telefonlaryň alan SMS-leri şu ýerde görünýär (provider sazlamasy + telefon programmasyndaky düwme).';

  @override
  String get importTitle => 'Kontaktlary import et';

  @override
  String get importFile => 'Faýldan';

  @override
  String get importFileHint =>
      'Excel, CSV, Word ýa-da vCard (.vcf). Bar bolan kontaktlar birleşdirilýär, gaýtalanmaýar.';

  @override
  String get importGoogle => 'Google Contacts';

  @override
  String get importGoogleHint =>
      'Google-yň rugsat ekranyny açýar; diňe okalýar, token saklanmaýar.';

  @override
  String get importApple => 'Apple / iCloud Contacts';

  @override
  String get importAppleHint =>
      'Apple ID + app-specific parol (appleid.apple.com). Bir gezek ulanylýar, saklanmaýar.';

  @override
  String get appleId => 'Apple ID (e-mail)';

  @override
  String get appPassword => 'App-specific parol';

  @override
  String get addToGroup => 'Topara goş';

  @override
  String get noGroup => 'Toparsyz';

  @override
  String get importNow => 'Import et';

  @override
  String importResult(int total, int created, int updated, int skipped) {
    return '$total setir: $created döredildi, $updated täzelendi, $skipped geçildi';
  }

  @override
  String get googleNotConfigured => 'Bu serwerde Google import sazlanmadyk';
}
