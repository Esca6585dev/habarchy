import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/api/api_client.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/config.dart';
import '../../core/push/push_service.dart';
import '../../core/settings.dart';
import '../../core/storage/storage.dart';
import '../../l10n/app_localizations.dart';

class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final auth = ref.watch(authControllerProvider).value;
    final user = auth is AuthSignedIn ? auth.user : null;
    final locale = ref.watch(localeProvider);
    final mode = ref.watch(themeModeProvider);
    return Scaffold(
      appBar: AppBar(title: Text(t.settings)),
      body: ListView(children: [
        if (user != null) ListTile(leading: const Icon(Icons.person_outline), title: Text(user.fullName.isEmpty ? user.email : user.fullName), subtitle: Text(user.email)),
        ListTile(leading: const Icon(Icons.dns_outlined), title: Text(t.serverUrl), subtitle: Text(ref.watch(apiUrlProvider))),
        const Divider(),
        ListTile(
          leading: const Icon(Icons.language),
          title: Text(t.language),
          trailing: DropdownButton<String>(
            value: locale?.languageCode ?? 'system',
            underline: const SizedBox(),
            items: const [
              DropdownMenuItem(value: 'system', child: Text('Auto')),
              DropdownMenuItem(value: 'tk', child: Text('Türkmençe')),
              DropdownMenuItem(value: 'ru', child: Text('Русский')),
              DropdownMenuItem(value: 'en', child: Text('English')),
            ],
            onChanged: (v) => ref.read(localeProvider.notifier).set(v == null || v == 'system' ? null : Locale(v)),
          ),
        ),
        ListTile(
          leading: const Icon(Icons.brightness_6_outlined),
          title: Text(t.theme),
          trailing: SegmentedButton<ThemeMode>(
            segments: [
              ButtonSegment(value: ThemeMode.system, label: Text(t.themeSystem)),
              ButtonSegment(value: ThemeMode.light, icon: const Icon(Icons.light_mode_outlined)),
              ButtonSegment(value: ThemeMode.dark, icon: const Icon(Icons.dark_mode_outlined)),
            ],
            selected: {mode},
            onSelectionChanged: (s) => ref.read(themeModeProvider.notifier).set(s.first),
            showSelectedIcon: false,
          ),
        ),
        const Divider(),
        const _PushDemoSection(),
        const Divider(),
        ListTile(
          leading: const Icon(Icons.logout),
          title: Text(t.logout),
          onTap: () => ref.read(authControllerProvider.notifier).logout(),
        ),
        Padding(padding: const EdgeInsets.all(16), child: Text(t.version('0.6.0 · ${AppConfig.flavor}'), style: Theme.of(context).textTheme.labelSmall, textAlign: TextAlign.center)),
      ]),
    );
  }
}

/// Registers this device's FCM token with a project through the public API
/// and lists incoming pushes. Doubles as the firebase_messaging reference.
class _PushDemoSection extends ConsumerStatefulWidget {
  const _PushDemoSection();
  @override
  ConsumerState<_PushDemoSection> createState() => _PushDemoSectionState();
}

class _PushDemoSectionState extends ConsumerState<_PushDemoSection> {
  final _apiKey = TextEditingController();
  late final _externalId = TextEditingController(text: ref.read(prefsProvider).pushExternalId ?? '');
  String? _token;
  String? _status;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    ref.read(secureStoreProvider).pushApiKey.then((v) {
      if (v != null && mounted) setState(() => _apiKey.text = v);
    });
    ref.read(pushServiceProvider).token().then((v) {
      if (mounted) setState(() => _token = v);
    });
  }

  @override
  void dispose() {
    _apiKey.dispose();
    _externalId.dispose();
    super.dispose();
  }

  Future<void> _register() async {
    final t = AppLocalizations.of(context);
    final push = ref.read(pushServiceProvider);
    setState(() {
      _busy = true;
      _status = null;
    });
    try {
      if (!await push.requestPermission()) throw Exception('permission denied');
      final token = await push.token();
      if (token == null) throw Exception('no token');
      await ref.read(secureStoreProvider).savePushApiKey(_apiKey.text.trim());
      await ref.read(prefsProvider).setPushExternalId(_externalId.text.trim());
      await ref.read(adminRepositoryProvider).registerDevice(
            apiKey: _apiKey.text.trim(),
            token: token,
            platform: push.platformName,
            externalId: _externalId.text.trim(),
            appVersion: '0.6.0',
          );
      setState(() {
        _token = token;
        _status = t.deviceRegistered;
      });
    } on ApiException catch (e) {
      setState(() => _status = e.toString());
    } catch (e) {
      setState(() => _status = '${t.error}: $e');
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final push = ref.watch(pushServiceProvider);
    final pushes = ref.watch(receivedPushesProvider);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [const Icon(Icons.notifications_active_outlined), const SizedBox(width: 12), Text(t.pushDemo, style: Theme.of(context).textTheme.titleMedium)]),
        const SizedBox(height: 4),
        Text(t.pushDemoHint, style: Theme.of(context).textTheme.bodySmall),
        const SizedBox(height: 12),
        if (!push.isAvailable)
          Card(color: Theme.of(context).colorScheme.tertiaryContainer, child: ListTile(leading: const Icon(Icons.info_outline), title: Text(t.pushNotConfigured, style: const TextStyle(fontSize: 13))))
        else ...[
          TextField(controller: _apiKey, decoration: InputDecoration(labelText: t.apiKeyForPush, hintText: 'hb_test_…'), obscureText: true),
          const SizedBox(height: 8),
          TextField(controller: _externalId, decoration: InputDecoration(labelText: t.externalId, hintText: 'user-42')),
          const SizedBox(height: 8),
          FilledButton.icon(onPressed: _busy || _apiKey.text.isEmpty ? null : _register, icon: const Icon(Icons.app_registration), label: Text(t.registerDevice)),
          if (_status != null) Padding(padding: const EdgeInsets.only(top: 6), child: Text(_status!, style: Theme.of(context).textTheme.bodySmall)),
          if (_token != null)
            ListTile(
              contentPadding: EdgeInsets.zero,
              title: Text(t.fcmToken, style: Theme.of(context).textTheme.labelSmall),
              subtitle: Text(_token!, maxLines: 2, overflow: TextOverflow.ellipsis, style: const TextStyle(fontFamily: 'monospace', fontSize: 11)),
              trailing: IconButton(
                icon: const Icon(Icons.copy, size: 18),
                onPressed: () {
                  Clipboard.setData(ClipboardData(text: _token!));
                  ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(t.copied)));
                },
              ),
            ),
          const SizedBox(height: 8),
          Text(t.incomingPushes, style: Theme.of(context).textTheme.labelLarge),
          if (pushes.isEmpty) Padding(padding: const EdgeInsets.symmetric(vertical: 8), child: Text(t.noPushesYet, style: Theme.of(context).textTheme.bodySmall)),
          for (final p in pushes)
            Card(
              child: ListTile(
                dense: true,
                leading: Icon(p.openedFromTray ? Icons.open_in_new : Icons.notifications),
                title: Text(p.title.isEmpty ? '(data)' : p.title),
                subtitle: Text('${p.body}\n${p.data}'),
                isThreeLine: true,
              ),
            ),
        ],
      ]),
    );
  }
}
