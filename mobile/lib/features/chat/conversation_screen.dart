import 'dart:async';

import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import '../contacts/contacts_screen.dart' show confirm;
import 'auth_image.dart';
import 'chat_screen.dart';

final chatMessagesProvider = FutureProvider.autoDispose.family<List<ChatMessage>, String>((ref, channelId) async {
  final msgs = await ref.read(adminRepositoryProvider).chatMessages(channelId);
  return msgs.reversed.toList(); // oldest first for display
});

final meUserProvider = FutureProvider.autoDispose<User>((ref) async => ref.read(adminRepositoryProvider).me());

/// One conversation; polls for new messages and marks the channel read.
class ChatConversationScreen extends ConsumerStatefulWidget {
  const ChatConversationScreen({super.key, required this.channelId});
  final String channelId;
  @override
  ConsumerState<ChatConversationScreen> createState() => _ChatConversationScreenState();
}

class _ChatConversationScreenState extends ConsumerState<ChatConversationScreen> {
  final _text = TextEditingController();
  final _scroll = ScrollController();
  Timer? _poll;
  bool _sending = false;

  @override
  void initState() {
    super.initState();
    _markRead();
    _poll = Timer.periodic(const Duration(seconds: 4), (_) {
      ref.invalidate(chatMessagesProvider(widget.channelId));
      _markRead();
    });
  }

  void _markRead() {
    ref.read(adminRepositoryProvider).markChatRead(widget.channelId).ignore();
    ref.invalidate(chatChannelsProvider);
  }

  @override
  void dispose() {
    _poll?.cancel();
    _text.dispose();
    _scroll.dispose();
    super.dispose();
  }

  void _jumpToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scroll.hasClients) _scroll.jumpTo(_scroll.position.maxScrollExtent);
    });
  }

  Future<void> _send({String? attachmentId}) async {
    final body = _text.text.trim();
    if (body.isEmpty && attachmentId == null) return;
    setState(() => _sending = true);
    try {
      await ref.read(adminRepositoryProvider).postChatMessage(widget.channelId, body: body, attachmentId: attachmentId);
      _text.clear();
      ref.invalidate(chatMessagesProvider(widget.channelId));
      ref.invalidate(chatChannelsProvider);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  Future<void> _attach() async {
    final picked = await FilePicker.pickFiles(type: FileType.image);
    if (picked.isEmpty) return;
    final f = picked.first;
    final bytes = await f.readAsBytes();
    try {
      final id = await ref.read(adminRepositoryProvider).uploadAttachment(filename: f.name, bytes: bytes);
      await _send(attachmentId: id);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final messages = ref.watch(chatMessagesProvider(widget.channelId));
    final me = ref.watch(meUserProvider).value;
    final title = ref.watch(chatChannelsProvider).value?.where((c) => c.id == widget.channelId).firstOrNull;
    return Scaffold(
      appBar: AppBar(title: Text(title == null ? t.chat : (title.kind == 'direct' ? title.peerName : title.name))),
      body: Column(children: [
        Expanded(
          child: messages.when(
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(chatMessagesProvider(widget.channelId))),
            data: (rows) {
              _jumpToBottom();
              if (rows.isEmpty) return EmptyState(text: t.chatNoMessages, icon: Icons.chat_bubble_outline);
              return ListView.builder(
                controller: _scroll,
                padding: const EdgeInsets.all(12),
                itemCount: rows.length,
                itemBuilder: (context, i) => _Bubble(message: rows[i], mine: rows[i].userId != null && rows[i].userId == me?.id, onDelete: () async {
                  if (!await confirm(context, t.delete)) return;
                  await ref.read(adminRepositoryProvider).deleteChatMessage(widget.channelId, rows[i].id);
                  ref.invalidate(chatMessagesProvider(widget.channelId));
                }),
              );
            },
          ),
        ),
        SafeArea(
          top: false,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(8, 4, 8, 8),
            child: Row(children: [
              IconButton(icon: const Icon(Icons.image_outlined), onPressed: _sending ? null : _attach),
              Expanded(
                child: TextField(
                  controller: _text,
                  minLines: 1,
                  maxLines: 4,
                  textCapitalization: TextCapitalization.sentences,
                  decoration: InputDecoration(hintText: t.chatMessage, isDense: true, border: const OutlineInputBorder()),
                ),
              ),
              IconButton(icon: const Icon(Icons.send), onPressed: _sending ? null : () => _send()),
            ]),
          ),
        ),
      ]),
    );
  }
}

class _Bubble extends StatelessWidget {
  const _Bubble({required this.message, required this.mine, required this.onDelete});
  final ChatMessage message;
  final bool mine;
  final VoidCallback onDelete;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final bubble = Container(
      constraints: BoxConstraints(maxWidth: MediaQuery.sizeOf(context).width * 0.72),
      margin: const EdgeInsets.symmetric(vertical: 3),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: mine ? scheme.primary : scheme.surfaceContainerHighest,
        borderRadius: BorderRadius.only(
          topLeft: const Radius.circular(14),
          topRight: const Radius.circular(14),
          bottomLeft: Radius.circular(mine ? 14 : 4),
          bottomRight: Radius.circular(mine ? 4 : 14),
        ),
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        if (!mine) Text(message.authorName, style: TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: scheme.primary)),
        if (message.attachmentId != null) Padding(padding: const EdgeInsets.only(top: 4, bottom: 4), child: ChatImage(attachmentId: message.attachmentId!)),
        if (message.body.isNotEmpty) Text(message.body, style: TextStyle(color: mine ? scheme.onPrimary : scheme.onSurface)),
        Text(shortDate(message.createdAt), style: TextStyle(fontSize: 10, color: (mine ? scheme.onPrimary : scheme.onSurfaceVariant).withValues(alpha: 0.7))),
      ]),
    );
    return Row(
      mainAxisAlignment: mine ? MainAxisAlignment.end : MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        if (!mine) Padding(padding: const EdgeInsets.only(right: 6), child: ChatAvatar(name: message.authorName, avatarId: message.authorAvatar, size: 28)),
        mine ? GestureDetector(onLongPress: onDelete, child: bubble) : bubble,
      ],
    );
  }
}
