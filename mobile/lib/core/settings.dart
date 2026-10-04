import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'storage/storage.dart';

class LocaleNotifier extends Notifier<Locale?> {
  @override
  Locale? build() {
    final s = ref.read(prefsProvider).locale;
    return s == null ? null : Locale(s);
  }

  Future<void> set(Locale? l) async {
    await ref.read(prefsProvider).setLocale(l?.languageCode);
    state = l;
  }
}

final localeProvider = NotifierProvider<LocaleNotifier, Locale?>(LocaleNotifier.new);

class ThemeModeNotifier extends Notifier<ThemeMode> {
  @override
  ThemeMode build() => switch (ref.read(prefsProvider).themeMode) { 'light' => ThemeMode.light, 'dark' => ThemeMode.dark, _ => ThemeMode.system };

  Future<void> set(ThemeMode m) async {
    await ref.read(prefsProvider).setThemeMode(m.name);
    state = m;
  }
}

final themeModeProvider = NotifierProvider<ThemeModeNotifier, ThemeMode>(ThemeModeNotifier.new);
