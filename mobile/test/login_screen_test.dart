import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:habarchy_admin/core/storage/storage.dart';
import 'package:habarchy_admin/features/login/login_screen.dart';
import 'package:habarchy_admin/l10n/app_localizations.dart';
import 'package:habarchy_admin/l10n/tk_fallback.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  testWidgets('login form validates and renders in Turkmen by default', (tester) async {
    SharedPreferences.setMockInitialValues({});
    final prefs = await Prefs.load();
    await tester.pumpWidget(ProviderScope(
      overrides: [prefsProvider.overrideWithValue(prefs)],
      child: const MaterialApp(
        locale: Locale('tk'),
        localizationsDelegates: [AppLocalizations.delegate, ...habarchyLocalizationsDelegates],
        supportedLocales: AppLocalizations.supportedLocales,
        home: LoginScreen(),
      ),
    ));
    expect(find.text('Habarçy'), findsOneWidget);
    expect(find.byKey(const Key('sign-in')), findsOneWidget);

    // Empty submit shows validation messages, no network call.
    await tester.tap(find.byKey(const Key('sign-in')));
    await tester.pump();
    expect(find.text('E-mail'), findsWidgets);
    expect(find.text('Parol'), findsWidgets);
  });
}
