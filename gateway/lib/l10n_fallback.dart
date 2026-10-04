import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

/// Flutter ships framework strings (dialog buttons, tooltips) for en/ru but
/// not for Turkmen. These delegates accept every locale and serve the
/// framework defaults, so the app never crashes on a `tk` phone. All visible
/// app strings come from [S] and are translated.
class AnyMaterialLocalizationsDelegate extends LocalizationsDelegate<MaterialLocalizations> {
  const AnyMaterialLocalizationsDelegate();
  @override
  bool isSupported(Locale locale) => true;
  @override
  Future<MaterialLocalizations> load(Locale locale) => DefaultMaterialLocalizations.load(locale);
  @override
  bool shouldReload(AnyMaterialLocalizationsDelegate old) => false;
}

class AnyCupertinoLocalizationsDelegate extends LocalizationsDelegate<CupertinoLocalizations> {
  const AnyCupertinoLocalizationsDelegate();
  @override
  bool isSupported(Locale locale) => true;
  @override
  Future<CupertinoLocalizations> load(Locale locale) => DefaultCupertinoLocalizations.load(locale);
  @override
  bool shouldReload(AnyCupertinoLocalizationsDelegate old) => false;
}

class AnyWidgetsLocalizationsDelegate extends LocalizationsDelegate<WidgetsLocalizations> {
  const AnyWidgetsLocalizationsDelegate();
  @override
  bool isSupported(Locale locale) => true;
  @override
  Future<WidgetsLocalizations> load(Locale locale) => DefaultWidgetsLocalizations.load(locale);
  @override
  bool shouldReload(AnyWidgetsLocalizationsDelegate old) => false;
}

const List<LocalizationsDelegate<dynamic>> gatewayLocalizationsDelegates = [
  AnyMaterialLocalizationsDelegate(),
  AnyCupertinoLocalizationsDelegate(),
  AnyWidgetsLocalizationsDelegate(),
];
