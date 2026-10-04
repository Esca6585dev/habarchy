import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../l10n/app_localizations.dart';

/// Secondary destinations (templates, providers, API keys, contacts, settings).
class MoreScreen extends StatelessWidget {
  const MoreScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final items = [
      (Icons.forum_outlined, t.chat, '/chat'),
      (Icons.account_circle_outlined, t.profileTitle, '/profile'),
      (Icons.description_outlined, t.templates, '/templates'),
      (Icons.power_outlined, t.providers, '/providers'),
      (Icons.move_to_inbox_outlined, t.inboundSms, '/inbound'),
      (Icons.people_outline, t.contacts, '/contacts'),
      (Icons.key_outlined, t.apiKeys, '/api-keys'),
      (Icons.settings_outlined, t.settings, '/settings'),
    ];
    return Scaffold(
      appBar: AppBar(title: Text(t.more)),
      body: ListView(children: [
        for (final i in items)
          ListTile(leading: Icon(i.$1), title: Text(i.$2), trailing: const Icon(Icons.chevron_right), onTap: () => context.go(i.$3)),
      ]),
    );
  }
}
