import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

/// Flutter ships no Material/Cupertino strings for Turkmen (`tk`), so the
/// framework widgets (date pickers, tooltips, dialogs) fall back to Russian,
/// which every Turkmen user reads. App strings themselves come from
/// `AppLocalizations` and are fully translated.
const Locale _fallback = Locale('ru');

class TkMaterialLocalizationsDelegate extends LocalizationsDelegate<MaterialLocalizations> {
  const TkMaterialLocalizationsDelegate();
  @override
  bool isSupported(Locale locale) => locale.languageCode == 'tk';
  @override
  Future<MaterialLocalizations> load(Locale locale) => GlobalMaterialLocalizations.delegate.load(_fallback);
  @override
  bool shouldReload(TkMaterialLocalizationsDelegate old) => false;
}

class TkCupertinoLocalizationsDelegate extends LocalizationsDelegate<CupertinoLocalizations> {
  const TkCupertinoLocalizationsDelegate();
  @override
  bool isSupported(Locale locale) => locale.languageCode == 'tk';
  @override
  Future<CupertinoLocalizations> load(Locale locale) => GlobalCupertinoLocalizations.delegate.load(_fallback);
  @override
  bool shouldReload(TkCupertinoLocalizationsDelegate old) => false;
}

class TkWidgetsLocalizationsDelegate extends LocalizationsDelegate<WidgetsLocalizations> {
  const TkWidgetsLocalizationsDelegate();
  @override
  bool isSupported(Locale locale) => locale.languageCode == 'tk';
  @override
  Future<WidgetsLocalizations> load(Locale locale) => GlobalWidgetsLocalizations.delegate.load(_fallback);
  @override
  bool shouldReload(TkWidgetsLocalizationsDelegate old) => false;
}

/// Everything a `MaterialApp` needs for tk/ru/en.
const List<LocalizationsDelegate<dynamic>> habarchyLocalizationsDelegates = [
  TkMaterialLocalizationsDelegate(),
  TkCupertinoLocalizationsDelegate(),
  TkWidgetsLocalizationsDelegate(),
  GlobalMaterialLocalizations.delegate,
  GlobalWidgetsLocalizations.delegate,
  GlobalCupertinoLocalizations.delegate,
];
