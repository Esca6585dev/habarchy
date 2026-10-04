import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'gateway_channel.dart';
import 'l10n_fallback.dart';
import 'strings.dart';

void main() {
  runApp(const GatewayApp());
}

class GatewayApp extends StatefulWidget {
  const GatewayApp({super.key, this.channel});
  final GatewayChannel? channel;

  @override
  State<GatewayApp> createState() => _GatewayAppState();
}

class _GatewayAppState extends State<GatewayApp> {
  Locale? _locale;

  @override
  Widget build(BuildContext context) {
    final scheme = ColorScheme.fromSeed(seedColor: const Color(0xFF0F7A5A));
    return MaterialApp(
      title: 'Habarçy Gateway',
      debugShowCheckedModeBanner: false,
      locale: _locale,
      supportedLocales: const [Locale('tk'), Locale('ru'), Locale('en')],
      localizationsDelegates: gatewayLocalizationsDelegates,
      localeResolutionCallback: (device, supported) {
        if (_locale != null) return _locale;
        if (device != null && supported.any((l) => l.languageCode == device.languageCode)) return Locale(device.languageCode);
        return const Locale('tk');
      },
      theme: ThemeData(colorScheme: scheme, useMaterial3: true, inputDecorationTheme: const InputDecorationTheme(border: OutlineInputBorder(), isDense: true)),
      darkTheme: ThemeData(colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF0F7A5A), brightness: Brightness.dark), useMaterial3: true),
      home: HomeScreen(channel: widget.channel ?? GatewayChannel(), onLocale: (l) => setState(() => _locale = l)),
    );
  }
}

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key, required this.channel, required this.onLocale});
  final GatewayChannel channel;
  final ValueChanged<Locale> onLocale;

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final _url = TextEditingController();
  final _key = TextEditingController();
  final _testTo = TextEditingController();
  bool _forwardInbound = false;
  bool _hasPermissions = true;
  bool _ignoringBattery = true;
  GatewayStatus _status = GatewayStatus.empty;
  List<String> _log = const [];
  Timer? _timer;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    _load();
    _timer = Timer.periodic(const Duration(seconds: 2), (_) => _refresh());
  }

  @override
  void dispose() {
    _timer?.cancel();
    _url.dispose();
    _key.dispose();
    _testTo.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final cfg = await widget.channel.config();
      _url.text = cfg.url;
      _key.text = cfg.key;
      _forwardInbound = cfg.forwardInbound;
      _hasPermissions = await widget.channel.hasPermissions();
      _ignoringBattery = await widget.channel.isIgnoringBatteryOptimizations();
    } on MissingPluginException {
      // Running outside Android (tests / desktop): keep defaults.
    } on PlatformException {
      // ignore
    }
    await _refresh();
  }

  Future<void> _refresh() async {
    try {
      final s = await widget.channel.status();
      final l = await widget.channel.log();
      if (mounted) setState(() { _status = s; _log = l; });
    } on MissingPluginException {
      // not on Android
    } on PlatformException {
      // ignore
    }
  }

  Future<void> _save() async {
    final t = S.of(context);
    if (_url.text.trim().isEmpty || _key.text.trim().isEmpty) {
      _snack(t('fillAll'));
      return;
    }
    await widget.channel.saveConfig(url: _url.text.trim(), key: _key.text.trim(), forwardInbound: _forwardInbound);
    _snack(t('saved'));
  }

  Future<void> _toggle() async {
    setState(() => _busy = true);
    try {
      if (_status.running) {
        await widget.channel.stop();
      } else {
        await _save();
        if (!_hasPermissions) {
          _hasPermissions = await widget.channel.requestPermissions();
          if (!_hasPermissions) return;
        }
        await widget.channel.start();
      }
    } on PlatformException catch (e) {
      _snack(e.message ?? e.code);
    } finally {
      if (mounted) setState(() => _busy = false);
      await _refresh();
    }
  }

  void _snack(String msg) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(msg)));
  }

  String _ago(S t, DateTime? d) {
    if (d == null) return t('never');
    return t('secondsAgo', {'n': DateTime.now().difference(d).inSeconds});
  }

  @override
  Widget build(BuildContext context) {
    final t = S.of(context);
    final cs = Theme.of(context).colorScheme;
    final running = _status.running;
    final online = running && _status.connected;
    return Scaffold(
      appBar: AppBar(
        title: Text(t('title')),
        actions: [
          PopupMenuButton<Locale>(
            icon: const Icon(Icons.language),
            tooltip: t('language'),
            onSelected: widget.onLocale,
            itemBuilder: (_) => const [
              PopupMenuItem(value: Locale('tk'), child: Text('Türkmençe')),
              PopupMenuItem(value: Locale('ru'), child: Text('Русский')),
              PopupMenuItem(value: Locale('en'), child: Text('English')),
            ],
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Status card
          Card(
            color: online ? Colors.green.withValues(alpha: 0.12) : (running ? Colors.amber.withValues(alpha: 0.15) : cs.surfaceContainerHighest),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Row(children: [
                  Icon(online ? Icons.cloud_done : (running ? Icons.cloud_off : Icons.pause_circle_outline), color: online ? Colors.green.shade700 : (running ? Colors.amber.shade800 : cs.outline), size: 32),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                      Text(running ? t('running') : t('stopped'), style: Theme.of(context).textTheme.titleMedium, key: const Key('status-title')),
                      Text(running ? (online ? t('connected') : t('disconnected')) : t('intro'), style: Theme.of(context).textTheme.bodySmall),
                    ]),
                  ),
                ]),
                const SizedBox(height: 12),
                Wrap(spacing: 16, runSpacing: 8, children: [
                  if (_status.providerName.isNotEmpty) _stat(t('provider'), _status.providerName),
                  _stat(t('pending'), '${_status.pending}'),
                  _stat(t('sent'), '${_status.sent}'),
                  _stat(t('delivered'), '${_status.delivered}'),
                  _stat(t('failed'), '${_status.failed}'),
                  if (_forwardInbound) _stat(t('inbound'), '${_status.inbound}'),
                  _stat(t('lastPoll'), _ago(t, _status.lastPollAt)),
                ]),
                if (_status.lastError.isNotEmpty) ...[
                  const SizedBox(height: 8),
                  Text('${t('lastError')}: ${_status.lastError}', style: TextStyle(color: cs.error, fontSize: 12)),
                ],
                const SizedBox(height: 12),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton.icon(
                    key: const Key('toggle'),
                    onPressed: _busy ? null : _toggle,
                    icon: Icon(running ? Icons.stop : Icons.play_arrow),
                    label: Text(running ? t('stop') : t('start')),
                    style: running ? FilledButton.styleFrom(backgroundColor: cs.error) : null,
                  ),
                ),
              ]),
            ),
          ),
          if (!_hasPermissions) _warning(Icons.sms_failed_outlined, t('permissions'), t('grant'), () async {
            _hasPermissions = await widget.channel.requestPermissions();
            setState(() {});
          }),
          if (!_ignoringBattery) _warning(Icons.battery_alert, t('battery'), t('disable'), () async {
            await widget.channel.requestIgnoreBatteryOptimizations();
            _ignoringBattery = await widget.channel.isIgnoringBatteryOptimizations();
            setState(() {});
          }),
          const SizedBox(height: 16),
          // Config
          Text(t('apiUrl'), style: Theme.of(context).textTheme.labelLarge),
          const SizedBox(height: 6),
          TextField(controller: _url, key: const Key('url'), keyboardType: TextInputType.url, decoration: const InputDecoration(hintText: 'https://habarchy.example.tm')),
          const SizedBox(height: 12),
          Text(t('gatewayKey'), style: Theme.of(context).textTheme.labelLarge),
          const SizedBox(height: 6),
          TextField(controller: _key, key: const Key('key'), decoration: InputDecoration(hintText: 'gw_…', helperText: t('keyHint'), helperMaxLines: 2)),
          SwitchListTile(
            contentPadding: EdgeInsets.zero,
            value: _forwardInbound,
            onChanged: (v) async {
              setState(() => _forwardInbound = v);
              if (v) {
                try {
                  final ok = await widget.channel.requestInboundPermissions();
                  if (!ok) _snack(t('inboundNoPermission'));
                } on MissingPluginException {
                  // not on Android
                } on PlatformException {
                  // ignore
                }
              }
              await _save();
            },
            title: Text(t('forwardInbound')),
            subtitle: Text(
              !_forwardInbound
                  ? t('inboundOff')
                  : !_status.serverInboundEnabled
                      ? t('inboundServerOff')
                      : _status.inboundActive
                          ? t('inboundOn')
                          : t('inboundNoPermission'),
              style: TextStyle(fontSize: 12, color: _forwardInbound && !_status.inboundActive ? cs.error : null),
            ),
          ),
          Align(alignment: Alignment.centerRight, child: OutlinedButton.icon(onPressed: _save, icon: const Icon(Icons.save_outlined), label: Text(t('save')))),
          const SizedBox(height: 16),
          // Test SMS
          ExpansionTile(
            title: Text(t('testSms')),
            tilePadding: EdgeInsets.zero,
            children: [
              Row(children: [
                Expanded(child: TextField(controller: _testTo, keyboardType: TextInputType.phone, decoration: InputDecoration(labelText: t('phone'), hintText: '+99365123456'))),
                const SizedBox(width: 8),
                FilledButton.tonal(
                  onPressed: () async {
                    try {
                      await widget.channel.testSms(_testTo.text.trim(), 'Habarçy gateway test ✓');
                      _snack(t('sent'));
                    } on PlatformException catch (e) {
                      _snack(e.message ?? e.code);
                    }
                  },
                  child: Text(t('send')),
                ),
              ]),
              const SizedBox(height: 8),
            ],
          ),
          const SizedBox(height: 8),
          // Log
          Row(children: [
            Text(t('log'), style: Theme.of(context).textTheme.labelLarge),
            const Spacer(),
            TextButton(onPressed: () async { await widget.channel.clearLog(); await _refresh(); }, child: Text(t('clear'))),
          ]),
          Container(
            constraints: const BoxConstraints(minHeight: 120),
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(color: cs.surfaceContainerHighest, borderRadius: BorderRadius.circular(8)),
            child: _log.isEmpty
                ? Text('—', style: TextStyle(color: cs.outline))
                : Column(crossAxisAlignment: CrossAxisAlignment.start, children: [for (final l in _log) Text(l, style: const TextStyle(fontFamily: 'monospace', fontSize: 12))]),
          ),
        ],
      ),
    );
  }

  Widget _stat(String label, String value) => Column(crossAxisAlignment: CrossAxisAlignment.start, mainAxisSize: MainAxisSize.min, children: [
        Text(label, style: Theme.of(context).textTheme.labelSmall),
        Text(value, style: Theme.of(context).textTheme.titleSmall),
      ]);

  Widget _warning(IconData icon, String text, String action, VoidCallback onTap) => Card(
        margin: const EdgeInsets.only(top: 8),
        child: ListTile(leading: Icon(icon, color: Theme.of(context).colorScheme.error), title: Text(text, style: const TextStyle(fontSize: 13)), trailing: TextButton(onPressed: onTap, child: Text(action))),
      );
}
