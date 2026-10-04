import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/api/admin_repository.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import 'auth_image.dart';
import 'new_chat_sheet.dart';

final chatChannelsProvider = FutureProvider.autoDispose<List<ChatChannel>>((ref) async {
  return ref.read(adminRepositoryProvider).chatChannels();
});

/// Channel list; polls every few seconds for new messages / unread counts.
class ChatScreen extends ConsumerStatefulWidget {
  const ChatScreen({super.key});
  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  Timer? _poll;
  @override
  void initState() {
    super.initState();
    _poll = Timer.periodic(const Duration(seconds: 6), (_) => ref.invalidate(chatChannelsProvider));
  }

  @override
  void dispose() {
    _poll?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final list = ref.watch(chatChannelsProvider);
    return Scaffold(
      appBar: AppBar(title: Text(t.chat)),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () async {
          final id = await showNewChatSheet(context, ref);
          if (id != null && context.mounted) context.push('/chat/$id');
        },
        icon: const Icon(Icons.add_comment_outlined),
        label: Text(t.chatNew),
      ),
      body: list.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(chatChannelsProvider)),
        data: (rows) => rows.isEmpty
            ? EmptyState(text: t.chatEmpty, icon: Icons.forum_outlined)
            : RefreshIndicator(
                onRefresh: () => ref.refresh(chatChannelsProvider.future),
                child: ListView.separated(
                  itemCount: rows.length,
                  separatorBuilder: (_, _) => const Divider(height: 1),
                  itemBuilder: (context, i) {
                    final c = rows[i];
                    final title = c.kind == 'direct' ? (c.peerName.isEmpty ? '—' : c.peerName) : c.name;
                    return ListTile(
                      leading: c.kind == 'direct'
                          ? ChatAvatar(name: c.peerName, avatarId: c.peerAvatar, size: 44)
                          : CircleAvatar(radius: 22, child: Icon(c.kind == 'private' ? Icons.lock_outline : Icons.tag)),
                      title: Text(title, maxLines: 1, overflow: TextOverflow.ellipsis),
                      subtitle: Text('${c.lastHasFile ? '📷 ' : ''}${c.lastBody.isEmpty ? t.chatNoMessages : c.lastBody}', maxLines: 1, overflow: TextOverflow.ellipsis),
                      trailing: c.unread > 0 ? Badge(label: Text('${c.unread}')) : null,
                      onTap: () => context.push('/chat/${c.id}'),
                    );
                  },
                ),
              ),
      ),
    );
  }
}
