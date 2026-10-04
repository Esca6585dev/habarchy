import 'package:flutter/material.dart';

import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';

/// Create / edit contact dialog. Returns the request body or null.
Future<Map<String, dynamic>?> showContactForm(BuildContext context, {Contact? initial}) {
  final t = AppLocalizations.of(context);
  final c = {
    'name': TextEditingController(text: initial?.name ?? ''),
    'phone': TextEditingController(text: initial?.phone ?? ''),
    'email': TextEditingController(text: initial?.email ?? ''),
    'whatsapp': TextEditingController(text: initial?.whatsapp ?? ''),
    'telegram_chat_id': TextEditingController(text: initial?.telegramChatId ?? ''),
    'slack_id': TextEditingController(text: initial?.slackId ?? ''),
    'external_id': TextEditingController(text: initial?.externalId ?? ''),
  };
  return showDialog<Map<String, dynamic>>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: Text(initial == null ? t.newContact : t.editContact),
      content: SingleChildScrollView(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          TextField(controller: c['name'], decoration: InputDecoration(labelText: t.name), textCapitalization: TextCapitalization.words),
          const SizedBox(height: 8),
          TextField(controller: c['phone'], decoration: InputDecoration(labelText: t.phone, hintText: '+99365123456'), keyboardType: TextInputType.phone),
          const SizedBox(height: 8),
          TextField(controller: c['email'], decoration: InputDecoration(labelText: t.email), keyboardType: TextInputType.emailAddress),
          const SizedBox(height: 8),
          TextField(controller: c['whatsapp'], decoration: InputDecoration(labelText: t.whatsapp, hintText: '= ${t.phone}'), keyboardType: TextInputType.phone),
          const SizedBox(height: 8),
          TextField(controller: c['telegram_chat_id'], decoration: InputDecoration(labelText: t.telegram)),
          const SizedBox(height: 8),
          TextField(controller: c['slack_id'], decoration: InputDecoration(labelText: t.slack, hintText: 'U0123ABCD')),
          const SizedBox(height: 8),
          TextField(controller: c['external_id'], decoration: const InputDecoration(labelText: 'External ID')),
          const SizedBox(height: 8),
          Text(t.contactHint, style: Theme.of(ctx).textTheme.bodySmall),
        ]),
      ),
      actions: [
        TextButton(onPressed: () => Navigator.pop(ctx), child: Text(t.cancel)),
        FilledButton(onPressed: () => Navigator.pop(ctx, {for (final e in c.entries) e.key: e.value.text.trim()}), child: Text(t.save)),
      ],
    ),
  );
}

/// "Name, +99365…, mail@…" → contact body, or null when no address.
Map<String, dynamic>? parseInlineContact(String line) {
  final parts = line.split(RegExp(r'[,;\t]')).map((p) => p.trim()).where((p) => p.isNotEmpty).toList();
  if (parts.isEmpty) return null;
  final out = <String, dynamic>{};
  for (final p in parts) {
    if (p.contains('@')) {
      out['email'] = p;
    } else if (RegExp(r'^\+?[\d\s()-]{6,}$').hasMatch(p)) {
      out.putIfAbsent('phone', () => p);
    } else if (RegExp(r'^[UC][A-Z0-9]{6,}$').hasMatch(p)) {
      out['slack_id'] = p;
    } else if (RegExp(r'^-?\d{5,}$').hasMatch(p)) {
      out['telegram_chat_id'] = p;
    } else {
      out.putIfAbsent('name', () => p);
    }
  }
  return out.containsKey('phone') || out.containsKey('email') || out.containsKey('slack_id') || out.containsKey('telegram_chat_id') ? out : null;
}

String contactLabel(Contact c) {
  for (final v in [c.name, c.externalId, c.phone, c.email, c.telegramChatId, c.slackId]) {
    if (v.isNotEmpty) return v;
  }
  return c.id.substring(0, 8);
}

String contactSubtitle(Contact c) => [c.phone, c.email, if (c.whatsapp.isNotEmpty && c.whatsapp != c.phone) 'WA ${c.whatsapp}', c.telegramChatId, c.slackId].where((v) => v.isNotEmpty).join(' · ');
