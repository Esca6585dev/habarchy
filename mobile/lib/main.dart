import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'app.dart';
import 'core/push/push_service.dart';
import 'core/storage/storage.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final prefs = await Prefs.load();
  final push = PushService();
  await push.init(); // no-op until Firebase is configured
  runApp(
    ProviderScope(
      overrides: [
        prefsProvider.overrideWithValue(prefs),
        pushServiceProvider.overrideWithValue(push),
      ],
      child: const HabarchyApp(),
    ),
  );
}
