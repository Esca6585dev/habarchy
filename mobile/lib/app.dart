import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/settings.dart';
import 'l10n/app_localizations.dart';
import 'l10n/tk_fallback.dart';
import 'router.dart';

class HabarchyApp extends ConsumerWidget {
  const HabarchyApp({super.key});

  static const seed = Color(0xFF0F7A5A);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);
    final mode = ref.watch(themeModeProvider);
    final locale = ref.watch(localeProvider);
    ThemeData theme(Brightness b) => ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: seed, brightness: b),
          useMaterial3: true,
          visualDensity: VisualDensity.adaptivePlatformDensity,
          cardTheme: const CardThemeData(margin: EdgeInsets.zero),
          inputDecorationTheme: const InputDecorationTheme(border: OutlineInputBorder(), isDense: true),
        );
    return MaterialApp.router(
      title: 'Habarchy',
      routerConfig: router,
      themeMode: mode,
      theme: theme(Brightness.light),
      darkTheme: theme(Brightness.dark),
      locale: locale,
      supportedLocales: AppLocalizations.supportedLocales,
      localizationsDelegates: const [AppLocalizations.delegate, ...habarchyLocalizationsDelegates],
      // Framework strings for Turkmen fall back to Russian (see l10n/tk_fallback.dart).
      localeResolutionCallback: (device, supported) {
        final want = locale ?? device;
        if (want != null && supported.any((l) => l.languageCode == want.languageCode)) return Locale(want.languageCode);
        return const Locale('tk');
      },
    );
  }
}
