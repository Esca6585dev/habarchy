import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import 'providers_screen.dart';

/// One credential field of a provider type.
class CredField {
  const CredField(this.key, this.label, {this.hint, this.secret = false, this.number = false, this.multiline = false, this.boolean = false, this.options, this.defaultValue});
  final String key;
  final String label;
  final String? hint;
  final bool secret;
  final bool number;
  final bool multiline;
  final bool boolean;
  final List<String>? options;
  final Object? defaultValue;
}

/// Field specs per provider type (mirrors docs/providers.md).
const providerTypes = <String, List<CredField>>{
  'android_sms': [
    CredField('sim_slot', 'SIM slot', options: ['-1', '0', '1'], defaultValue: '-1'),
    CredField('timeout_sec', 'Timeout (sec)', number: true, defaultValue: 45),
  ],
  'smtp': [
    CredField('host', 'SMTP host', hint: 'smtp.example.tm'),
    CredField('port', 'Port', number: true, defaultValue: 587),
    CredField('tls_mode', 'TLS', options: ['starttls', 'tls', 'none'], defaultValue: 'starttls'),
    CredField('username', 'Username'),
    CredField('password', 'Password', secret: true),
    CredField('from_name', 'From name', hint: 'Habarçy'),
    CredField('from_email', 'From e-mail', hint: 'no-reply@example.tm'),
  ],
  'smpp': [
    CredField('host', 'SMSC host', hint: 'smsc.operator.tm'),
    CredField('port', 'Port', number: true, defaultValue: 2775),
    CredField('system_id', 'System ID'),
    CredField('password', 'Password', secret: true),
    CredField('source_addr', 'Sender (source_addr)', hint: 'HABARCHY'),
    CredField('system_type', 'System type'),
    CredField('request_dlr', 'Request delivery reports', boolean: true, defaultValue: true),
    CredField('use_tls', 'TLS', boolean: true, defaultValue: false),
  ],
  'http_sms': [
    CredField('url', 'URL', hint: 'https://sms.example.tm/api/send'),
    CredField('method', 'Method', options: ['POST', 'GET'], defaultValue: 'POST'),
    CredField('body_template', 'Body template', hint: '{"to":"{{.To}}","text":{{.TextJSON}}}', multiline: true),
    CredField('sender', 'Sender'),
    CredField('headers_json', 'Headers (JSON)', hint: '{"Authorization":"Bearer TOKEN"}', multiline: true, secret: true),
    CredField('success_path', 'Success JSON path', hint: 'status'),
    CredField('success_equals', 'Success value', hint: 'OK'),
    CredField('message_id_path', 'Message id JSON path', hint: 'id'),
  ],
  'telegram_bot': [
    CredField('bot_token', 'Bot token', hint: '123456:ABC-DEF', secret: true),
    CredField('parse_mode', 'Parse mode', options: ['HTML', 'Markdown', 'MarkdownV2', ''], defaultValue: 'HTML'),
  ],
  'whatsapp_cloud': [
    CredField('access_token', 'Access token', secret: true),
    CredField('phone_number_id', 'Phone number ID'),
    CredField('api_version', 'API version', defaultValue: 'v20.0'),
  ],
  'slack': [
    CredField('bot_token', 'Bot token (xoxb-…)', secret: true),
    CredField('default_channel', 'Default channel', hint: 'C0123456789'),
    CredField('webhook_url', 'Incoming webhook URL (alternative)', secret: true),
  ],
  'fcm': [
    CredField('service_account_json', 'Service account JSON', multiline: true, secret: true),
  ],
};

/// Builds the credentials JSON from field values (type-specific shapes).
Map<String, dynamic> buildCredentials(String type, Map<String, Object?> v) {
  Object? parse(CredField f, Object? raw) {
    if (raw == null) return null;
    if (f.boolean) return raw == true || raw == 'true';
    if (f.number) return int.tryParse('$raw') ?? raw;
    if (f.key == 'sim_slot') return int.tryParse('$raw') ?? -1;
    return raw;
  }

  final out = <String, dynamic>{};
  for (final f in providerTypes[type] ?? const <CredField>[]) {
    final raw = v[f.key];
    if (raw == null || raw == '') continue;
    out[f.key] = parse(f, raw);
  }
  switch (type) {
    case 'http_sms':
      if (out.remove('headers_json') case final String h when h.trim().isNotEmpty) out['headers'] = jsonDecode(h);
      final sp = out.remove('success_path'), se = out.remove('success_equals');
      if (sp != null) out['success'] = {'json_path': sp, 'json_equals': ?se};
      if (out.remove('message_id_path') case final String mp) out['message_id'] = {'json_path': mp};
    case 'fcm':
      if (out.remove('service_account_json') case final String j when j.trim().isNotEmpty) out['service_account'] = jsonDecode(j);
  }
  return out;
}

/// Create / edit a provider with a typed form (no JSON typing on a phone).
class ProviderEditScreen extends ConsumerStatefulWidget {
  const ProviderEditScreen({super.key, required this.projectId, this.initial});
  final String projectId;
  final ProviderDetail? initial;

  @override
  ConsumerState<ProviderEditScreen> createState() => _ProviderEditScreenState();
}

class _ProviderEditScreenState extends ConsumerState<ProviderEditScreen> {
  late String _type = widget.initial?.type ?? 'android_sms';
  late final _name = TextEditingController(text: widget.initial?.name ?? '');
  late final _priority = TextEditingController(text: '${widget.initial?.priority ?? 100}');
  late final _rate = TextEditingController(text: '${widget.initial?.rateLimitPerSec ?? 0}');
  late bool _active = widget.initial?.isActive ?? true;
  final Map<String, Object?> _values = {};
  final Map<String, TextEditingController> _ctl = {};
  bool _saving = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _resetFields();
  }

  void _resetFields() {
    _values.clear();
    for (final f in providerTypes[_type] ?? const <CredField>[]) {
      Object? initial = f.defaultValue;
      final s = widget.initial?.settings;
      if (s != null && s.containsKey(f.key) && !f.secret) initial = s[f.key];
      _values[f.key] = initial is bool ? initial : (initial?.toString());
      _ctl[f.key]?.dispose();
      _ctl[f.key] = TextEditingController(text: initial is bool ? '' : (initial?.toString() ?? ''));
    }
  }

  @override
  void dispose() {
    _name.dispose();
    _priority.dispose();
    _rate.dispose();
    for (final c in _ctl.values) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _save() async {
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      final body = <String, dynamic>{
        'name': _name.text.trim(),
        'type': _type,
        'priority': int.tryParse(_priority.text) ?? 100,
        'rate_limit_per_sec': int.tryParse(_rate.text) ?? 0,
        'is_active': _active,
      };
      final creds = buildCredentials(_type, _values);
      final touched = creds.entries.any((e) => e.value != null && e.value != '' && e.value != providerTypes[_type]?.firstWhere((f) => f.key == e.key, orElse: () => const CredField('', '')).defaultValue);
      if (widget.initial == null || touched) body['credentials'] = creds;
      final repo = ref.read(adminRepositoryProvider);
      final saved = widget.initial == null ? await repo.createProvider(widget.projectId, body) : await repo.updateProvider(widget.projectId, widget.initial!.id, body);
      ref.invalidate(healthProvider);
      if (!mounted) return;
      if (saved.type == 'android_sms' && widget.initial == null) {
        await showPairing(context, ref, widget.projectId, saved.id);
      }
      if (mounted) Navigator.pop(context, saved);
    } catch (e) {
      setState(() => _error = '$e');
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final fields = providerTypes[_type] ?? const <CredField>[];
    return Scaffold(
      appBar: AppBar(title: Text(widget.initial == null ? t.newProvider : t.editProvider)),
      body: ListView(padding: const EdgeInsets.all(16), children: [
        TextField(controller: _name, key: const Key('provider-name'), decoration: InputDecoration(labelText: t.name), textCapitalization: TextCapitalization.sentences),
        const SizedBox(height: 12),
        DropdownButtonFormField<String>(
          initialValue: _type,
          decoration: InputDecoration(labelText: t.providerType),
          items: [for (final k in providerTypes.keys) DropdownMenuItem(value: k, child: Row(children: [Icon(channelIcon(_channelOf(k)), size: 16, color: channelColor(_channelOf(k))), const SizedBox(width: 8), Text(k)]))],
          onChanged: widget.initial != null
              ? null
              : (v) => setState(() {
                    _type = v ?? _type;
                    _resetFields();
                  }),
        ),
        const SizedBox(height: 12),
        Row(children: [
          Expanded(child: TextField(controller: _priority, decoration: InputDecoration(labelText: t.priority), keyboardType: TextInputType.number, inputFormatters: [FilteringTextInputFormatter.digitsOnly])),
          const SizedBox(width: 12),
          Expanded(child: TextField(controller: _rate, decoration: InputDecoration(labelText: t.rateLimit), keyboardType: TextInputType.number, inputFormatters: [FilteringTextInputFormatter.digitsOnly])),
        ]),
        SwitchListTile(contentPadding: EdgeInsets.zero, value: _active, onChanged: (v) => setState(() => _active = v), title: Text(t.active)),
        const Divider(),
        Text(t.credentials, style: Theme.of(context).textTheme.titleSmall),
        if (widget.initial != null) Text(t.keepCredentials, style: Theme.of(context).textTheme.bodySmall),
        if (_type == 'android_sms') Padding(padding: const EdgeInsets.only(top: 4), child: Text(t.pairingHint, style: Theme.of(context).textTheme.bodySmall)),
        const SizedBox(height: 8),
        for (final f in fields) ...[
          if (f.boolean)
            SwitchListTile(contentPadding: EdgeInsets.zero, value: _values[f.key] == true, onChanged: (v) => setState(() => _values[f.key] = v), title: Text(f.label))
          else if (f.options != null)
            DropdownButtonFormField<String>(
              initialValue: '${_values[f.key] ?? f.defaultValue ?? ''}',
              decoration: InputDecoration(labelText: f.label),
              items: [for (final o in f.options!) DropdownMenuItem(value: o, child: Text(o.isEmpty ? '—' : o))],
              onChanged: (v) => setState(() => _values[f.key] = v),
            )
          else
            TextField(
              controller: _ctl[f.key],
              key: Key('cred-${f.key}'),
              obscureText: f.secret && !f.multiline,
              maxLines: f.multiline ? 4 : 1,
              keyboardType: f.number ? TextInputType.number : (f.multiline ? TextInputType.multiline : TextInputType.text),
              style: f.multiline ? const TextStyle(fontFamily: 'monospace', fontSize: 12) : null,
              decoration: InputDecoration(labelText: f.label, hintText: f.hint),
              onChanged: (v) => _values[f.key] = v,
            ),
          const SizedBox(height: 10),
        ],
        if (_error != null) Text(_error!, style: TextStyle(color: Theme.of(context).colorScheme.error)),
        const SizedBox(height: 8),
        FilledButton.icon(
          key: const Key('save-provider'),
          onPressed: _saving || _name.text.trim().isEmpty ? null : _save,
          icon: const Icon(Icons.save_outlined),
          label: Text(t.save),
        ),
      ]),
    );
  }

  static String _channelOf(String type) => switch (type) {
        'smtp' => 'email',
        'fcm' => 'push',
        'telegram_bot' => 'telegram',
        'whatsapp_cloud' => 'whatsapp',
        'slack' => 'slack',
        _ => 'sms',
      };
}

/// Shows the API URL + gateway key for the phone gateway app.
Future<void> showPairing(BuildContext context, WidgetRef ref, String projectId, String providerId) async {
  final t = AppLocalizations.of(context);
  Pairing p;
  try {
    p = await ref.read(adminRepositoryProvider).pairing(projectId, providerId);
  } catch (e) {
    if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    return;
  }
  if (!context.mounted) return;
  await showDialog<void>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: Row(children: [
        const Icon(Icons.phone_android),
        const SizedBox(width: 8),
        Expanded(child: Text(t.pairing)),
        Chip(label: Text(p.online ? t.online : t.offline), backgroundColor: (p.online ? Colors.green : Colors.red).withValues(alpha: 0.15), visualDensity: VisualDensity.compact),
      ]),
      content: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(t.pairingHint, style: Theme.of(ctx).textTheme.bodySmall),
        const SizedBox(height: 12),
        _copyRow(ctx, t.apiUrl, p.apiUrl),
        const SizedBox(height: 8),
        _copyRow(ctx, t.gatewayKey, p.gatewayKey),
        if (p.lastSeenAt != null) Padding(padding: const EdgeInsets.only(top: 8), child: Text(shortDate(p.lastSeenAt), style: Theme.of(ctx).textTheme.bodySmall)),
      ]),
      actions: [FilledButton(onPressed: () => Navigator.pop(ctx), child: const Text('OK'))],
    ),
  );
}

Widget _copyRow(BuildContext ctx, String label, String value) => Row(children: [
      Expanded(
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(label, style: Theme.of(ctx).textTheme.labelSmall),
          SelectableText(value, style: const TextStyle(fontFamily: 'monospace', fontSize: 12)),
        ]),
      ),
      IconButton(
        icon: const Icon(Icons.copy, size: 18),
        onPressed: () {
          Clipboard.setData(ClipboardData(text: value));
          ScaffoldMessenger.of(ctx).showSnackBar(SnackBar(content: Text(AppLocalizations.of(ctx).copied)));
        },
      ),
    ]);
