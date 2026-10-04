import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import 'auth_image.dart';

/// Start a direct message or create a channel. Returns the channel id.
Future<String?> showNewChatSheet(BuildContext context, WidgetRef ref) async {
  final t = AppLocalizations.of(context);
  final repo = ref.read(adminRepositoryProvider);
  final me = ref.read(authControllerProvider).value;
  final myId = me is AuthSignedIn ? me.user.id : '';
  List<User> users = [];
  try {
    users = (await repo.users()).where((u) => u.id != myId).toList();
  } catch (_) {}
  if (!context.mounted) return null;

  return showModalBottomSheet<String>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (ctx) {
      final name = TextEditingController();
      var kind = 'public';
      final picked = <String>{};
      return StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.fromLTRB(16, 0, 16, MediaQuery.viewInsetsOf(ctx).bottom + 16),
          child: DefaultTabController(
            length: 2,
            child: Column(mainAxisSize: MainAxisSize.min, children: [
              TabBar(tabs: [Tab(text: t.chatDirect), Tab(text: t.chatChannel)]),
              SizedBox(
                height: 360,
                child: TabBarView(children: [
                  // Direct: pick a user
                  ListView(children: [
                    for (final u in users)
                      ListTile(
                        leading: ChatAvatar(name: u.fullName, avatarId: u.avatarId, size: 40),
                        title: Text(u.fullName.isEmpty ? u.email : u.fullName),
                        subtitle: Text(u.email),
                        onTap: () async {
                          final id = await repo.openDirect(u.id);
                          if (ctx.mounted) Navigator.pop(ctx, id);
                        },
                      ),
                  ]),
                  // Channel: kind + name + members
                  ListView(children: [
                    const SizedBox(height: 8),
                    SegmentedButton<String>(
                      segments: [ButtonSegment(value: 'public', label: Text(t.chatPublic)), ButtonSegment(value: 'private', label: Text(t.chatPrivate))],
                      selected: {kind},
                      onSelectionChanged: (s) => setState(() => kind = s.first),
                    ),
                    const SizedBox(height: 12),
                    TextField(controller: name, decoration: InputDecoration(labelText: t.chatName, hintText: 'general', border: const OutlineInputBorder())),
                    if (kind == 'private') ...[
                      const SizedBox(height: 12),
                      Text(t.members, style: Theme.of(ctx).textTheme.labelLarge),
                      for (final u in users)
                        CheckboxListTile(
                          value: picked.contains(u.id),
                          onChanged: (v) => setState(() => v == true ? picked.add(u.id) : picked.remove(u.id)),
                          title: Text(u.fullName.isEmpty ? u.email : u.fullName),
                          secondary: ChatAvatar(name: u.fullName, avatarId: u.avatarId, size: 32),
                          dense: true,
                        ),
                    ],
                    const SizedBox(height: 12),
                    FilledButton(
                      onPressed: () async {
                        if (name.text.trim().isEmpty) return;
                        try {
                          final id = await repo.createChatChannel(kind: kind, name: name.text.trim(), members: picked.toList());
                          if (ctx.mounted) Navigator.pop(ctx, id);
                        } catch (e) {
                          if (ctx.mounted) ScaffoldMessenger.of(ctx).showSnackBar(SnackBar(content: Text('$e')));
                        }
                      },
                      child: Text(t.save),
                    ),
                  ]),
                ]),
              ),
            ]),
          ),
        ),
      );
    },
  );
}
