import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import '../groups/groups_screen.dart';
import '../templates/templates_screen.dart';

const composeChannels = ['sms', 'whatsapp', 'telegram', 'email', 'push', 'slack'];

/// Compose a message to groups and/or addresses on any channel.
class ComposeScreen extends ConsumerStatefulWidget {
  const ComposeScreen({super.key, this.initialGroupId});
  final String? initialGroupId;

  @override
  ConsumerState<ComposeScreen> createState() => _ComposeScreenState();
}

class _ComposeScreenState extends ConsumerState<ComposeScreen> {
  String _channel = 'sms';
  final Set<String> _groups = {};
  String _template = '';
  final _subject = TextEditingController();
  final _body = TextEditingController();
  final _addresses = TextEditingController();
  final _data = TextEditingController(text: '{}');
  bool _sandbox = false;
  bool _sending = false;

  @override
  void initState() {
    super.initState();
    if (widget.initialGroupId != null) _groups.add(widget.initialGroupId!);
  }

  @override
  void didUpdateWidget(covariant ComposeScreen old) {
    super.didUpdateWidget(old);
    if (widget.initialGroupId != null && widget.initialGroupId != old.initialGroupId) setState(() => _groups.add(widget.initialGroupId!));
  }

  @override
  void dispose() {
    _subject.dispose();
    _body.dispose();
    _addresses.dispose();
    _data.dispose();
    super.dispose();
  }

  List<String> get _toList => _addresses.text.split('\n').map((s) => s.trim()).where((s) => s.isNotEmpty).toList();

  Map<String, dynamic>? get _parsedData {
    try {
      final v = jsonDecode(_data.text.trim().isEmpty ? '{}' : _data.text);
      return v is Map<String, dynamic> ? v : null;
    } catch (_) {
      return null;
    }
  }

  Future<void> _send(String projectId) async {
    final t = AppLocalizations.of(context);
    setState(() => _sending = true);
    try {
      final r = await ref.read(adminRepositoryProvider).send(projectId, {
        'channel': _channel,
        if (_template.isNotEmpty) 'template': _template,
        if (_template.isEmpty) 'body': _body.text,
        if (_subject.text.trim().isNotEmpty) 'subject': _subject.text.trim(),
        if (_channel == 'push' && _subject.text.trim().isNotEmpty) 'title': _subject.text.trim(),
        'data': _parsedData ?? {},
        'group_ids': _groups.toList(),
        'to': _toList,
        'is_test': _sandbox,
      });
      if (!mounted) return;
      await showDialog<void>(
        context: context,
        builder: (ctx) => AlertDialog(
          icon: Icon(r.rejected.isEmpty ? Icons.check_circle_outline : Icons.warning_amber_outlined, color: r.rejected.isEmpty ? Colors.green : Colors.amber),
          title: Text(t.queuedRejected(r.accepted, r.rejected.length)),
          content: r.rejected.isEmpty
              ? null
              : SizedBox(
                  width: double.maxFinite,
                  child: ListView(shrinkWrap: true, children: [
                    for (final x in r.rejected.take(20)) ListTile(dense: true, title: Text('${x['to']}'), subtitle: Text('${x['reason']}')),
                  ]),
                ),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(ctx);
                context.go('/messages');
              },
              child: Text(t.openMessages),
            ),
            FilledButton(onPressed: () => Navigator.pop(ctx), child: const Text('OK')),
          ],
        ),
      );
      if (r.rejected.isEmpty && mounted) setState(() => _body.clear());
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final project = ref.watch(currentProjectProvider);
    final groups = ref.watch(groupsProvider).value ?? const [];
    final templates = (ref.watch(templatesProvider).value ?? const []).where((x) => x.channel == _channel).map((x) => x.key).toSet().toList()..sort();
    if (!templates.contains(_template)) _template = '';
    final hasTarget = _groups.isNotEmpty || _toList.isNotEmpty;
    final canSend = project != null && hasTarget && (_template.isNotEmpty || _body.text.trim().isNotEmpty) && _parsedData != null && !_sending;
    return Scaffold(
      appBar: AppBar(title: Text(t.compose)),
      body: ListView(padding: const EdgeInsets.fromLTRB(16, 8, 16, 96), children: [
        Text(t.channel, style: Theme.of(context).textTheme.labelLarge),
        const SizedBox(height: 6),
        Wrap(spacing: 6, children: [
          for (final c in composeChannels)
            ChoiceChip(
              key: Key('channel-$c'),
              avatar: Icon(channelIcon(c), size: 16, color: channelColor(c)),
              label: Text(c.toUpperCase()),
              selected: _channel == c,
              onSelected: (_) => setState(() => _channel = c),
            ),
        ]),
        if (_channel == 'whatsapp') Padding(padding: const EdgeInsets.only(top: 6), child: Text(t.whatsappHint, style: Theme.of(context).textTheme.bodySmall)),
        const SizedBox(height: 16),
        Text(t.recipients, style: Theme.of(context).textTheme.labelLarge),
        const SizedBox(height: 6),
        if (groups.isEmpty)
          TextButton.icon(onPressed: () => context.go('/groups'), icon: const Icon(Icons.group_add_outlined), label: Text(t.noGroupsYet, maxLines: 2))
        else
          Wrap(spacing: 6, runSpacing: -6, children: [
            for (final g in groups)
              FilterChip(
                key: Key('group-${g.id}'),
                label: Text('${g.name} (${g.memberCount})'),
                selected: _groups.contains(g.id),
                onSelected: (v) => setState(() => v ? _groups.add(g.id) : _groups.remove(g.id)),
              ),
          ]),
        const SizedBox(height: 8),
        TextField(
          controller: _addresses,
          key: const Key('addresses'),
          maxLines: 3,
          style: const TextStyle(fontFamily: 'monospace', fontSize: 13),
          decoration: InputDecoration(labelText: t.addresses, helperText: t.addressesHint, hintText: _channel == 'email' ? 'user@example.tm' : (_channel == 'slack' ? 'C0123456789' : '+99365123456')),
          onChanged: (_) => setState(() {}),
        ),
        const SizedBox(height: 16),
        DropdownButtonFormField<String>(
          initialValue: _template.isEmpty ? '' : _template,
          decoration: InputDecoration(labelText: t.template),
          items: [DropdownMenuItem(value: '', child: Text(t.freeText)), for (final k in templates) DropdownMenuItem(value: k, child: Text(k))],
          onChanged: (v) => setState(() => _template = v ?? ''),
        ),
        const SizedBox(height: 8),
        if (_template.isEmpty) ...[
          if (_channel != 'sms') ...[
            TextField(controller: _subject, decoration: InputDecoration(labelText: t.subject)),
            const SizedBox(height: 8),
          ],
          TextField(
            controller: _body,
            key: const Key('body'),
            maxLines: 5,
            decoration: InputDecoration(labelText: t.body, hintText: 'Salam {{.name}}! …', alignLabelWithHint: true),
            onChanged: (_) => setState(() {}),
          ),
          const SizedBox(height: 8),
        ],
        ExpansionTile(
          title: Text(t.templateData, style: Theme.of(context).textTheme.labelLarge),
          tilePadding: EdgeInsets.zero,
          children: [
            TextField(
              controller: _data,
              maxLines: 3,
              style: const TextStyle(fontFamily: 'monospace', fontSize: 12),
              decoration: InputDecoration(errorText: _parsedData == null ? 'invalid JSON' : null),
              onChanged: (_) => setState(() {}),
            ),
          ],
        ),
        SwitchListTile(contentPadding: EdgeInsets.zero, value: _sandbox, onChanged: (v) => setState(() => _sandbox = v), title: Text(t.sandboxMode)),
        const SizedBox(height: 8),
        FilledButton.icon(
          key: const Key('send'),
          onPressed: canSend ? () => _send(project.id) : null,
          icon: _sending ? const SizedBox(width: 16, height: 16, child: CircularProgressIndicator(strokeWidth: 2)) : const Icon(Icons.send),
          label: Text(t.send),
        ),
        if (!hasTarget) Padding(padding: const EdgeInsets.only(top: 6), child: Text(t.selectRecipients, style: Theme.of(context).textTheme.bodySmall, textAlign: TextAlign.center)),
      ]),
    );
  }
}
