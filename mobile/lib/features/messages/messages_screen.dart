import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';

class MessageFilter {
  const MessageFilter({this.status, this.channel, this.search = ''});
  final String? status;
  final String? channel;
  final String search;
  MessageFilter copyWith({String? Function()? status, String? Function()? channel, String? search}) =>
      MessageFilter(status: status != null ? status() : this.status, channel: channel != null ? channel() : this.channel, search: search ?? this.search);
}

class MessageFilterNotifier extends Notifier<MessageFilter> {
  @override
  MessageFilter build() => const MessageFilter();
  void set(MessageFilter f) => state = f;
}

final messageFilterProvider = NotifierProvider<MessageFilterNotifier, MessageFilter>(MessageFilterNotifier.new);

/// Infinite list state: pages are appended as the user scrolls.
class MessagesPager extends AsyncNotifier<({List<Message> rows, String? next})> {
  @override
  Future<({List<Message> rows, String? next})> build() async {
    final pid = ref.watch(selectedProjectProvider);
    final f = ref.watch(messageFilterProvider);
    if (pid == null) return (rows: const <Message>[], next: null);
    return ref.read(adminRepositoryProvider).messages(pid, status: f.status, channel: f.channel, search: f.search);
  }

  Future<void> loadMore() async {
    final cur = state.value;
    final pid = ref.read(selectedProjectProvider);
    if (cur == null || cur.next == null || pid == null) return;
    final f = ref.read(messageFilterProvider);
    final more = await ref.read(adminRepositoryProvider).messages(pid, status: f.status, channel: f.channel, search: f.search, cursor: cur.next);
    state = AsyncData((rows: [...cur.rows, ...more.rows], next: more.next));
  }
}

final messagesPagerProvider = AsyncNotifierProvider<MessagesPager, ({List<Message> rows, String? next})>(MessagesPager.new);

class MessagesScreen extends ConsumerStatefulWidget {
  const MessagesScreen({super.key});
  @override
  ConsumerState<MessagesScreen> createState() => _MessagesScreenState();
}

class _MessagesScreenState extends ConsumerState<MessagesScreen> {
  final _search = TextEditingController();
  final _scroll = ScrollController();

  @override
  void initState() {
    super.initState();
    _scroll.addListener(() {
      if (_scroll.position.pixels > _scroll.position.maxScrollExtent - 300) ref.read(messagesPagerProvider.notifier).loadMore();
    });
  }

  @override
  void dispose() {
    _search.dispose();
    _scroll.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final filter = ref.watch(messageFilterProvider);
    final pager = ref.watch(messagesPagerProvider);
    const statuses = ['queued', 'processing', 'sent', 'delivered', 'failed', 'cancelled'];
    const channels = ['sms', 'email', 'push', 'telegram'];
    return Scaffold(
      appBar: AppBar(
        title: Text(t.messages),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(96),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(12, 0, 12, 8),
            child: Column(children: [
              TextField(
                controller: _search,
                decoration: InputDecoration(hintText: t.search, prefixIcon: const Icon(Icons.search), isDense: true, suffixIcon: _search.text.isEmpty ? null : IconButton(icon: const Icon(Icons.clear), onPressed: () {
                  _search.clear();
                  ref.read(messageFilterProvider.notifier).set(filter.copyWith(search: ''));
                })),
                textInputAction: TextInputAction.search,
                onSubmitted: (v) => ref.read(messageFilterProvider.notifier).set(filter.copyWith(search: v.trim())),
              ),
              const SizedBox(height: 6),
              SizedBox(
                height: 32,
                child: ListView(scrollDirection: Axis.horizontal, children: [
                  _chip(t.all, filter.status == null && filter.channel == null, () => ref.read(messageFilterProvider.notifier).set(filter.copyWith(status: () => null, channel: () => null))),
                  for (final s in statuses) _chip(statusLabel(t, s), filter.status == s, () => ref.read(messageFilterProvider.notifier).set(filter.copyWith(status: () => filter.status == s ? null : s))),
                  for (final c in channels) _chip(c.toUpperCase(), filter.channel == c, () => ref.read(messageFilterProvider.notifier).set(filter.copyWith(channel: () => filter.channel == c ? null : c))),
                ]),
              ),
            ]),
          ),
        ),
      ),
      body: pager.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(messagesPagerProvider)),
        data: (page) => page.rows.isEmpty
            ? const EmptyState()
            : RefreshIndicator(
                onRefresh: () => ref.refresh(messagesPagerProvider.future),
                child: ListView.separated(
                  controller: _scroll,
                  itemCount: page.rows.length + (page.next != null ? 1 : 0),
                  separatorBuilder: (_, _) => const Divider(height: 1),
                  itemBuilder: (context, i) {
                    if (i >= page.rows.length) return const Padding(padding: EdgeInsets.all(16), child: Center(child: CircularProgressIndicator()));
                    final m = page.rows[i];
                    return ListTile(
                      leading: Icon(channelIcon(m.channel), color: channelColor(m.channel)),
                      title: Row(children: [
                        Expanded(child: Text(m.to, style: const TextStyle(fontFamily: 'monospace'), overflow: TextOverflow.ellipsis)),
                        StatusChip(m.status),
                      ]),
                      subtitle: Text('${m.template.isNotEmpty ? '${m.template} · ' : ''}${m.errorCode.isNotEmpty ? m.errorCode : m.body}', maxLines: 1, overflow: TextOverflow.ellipsis),
                      trailing: Text(shortDate(m.createdAt), style: Theme.of(context).textTheme.labelSmall),
                      onTap: () => context.go('/messages/${m.id}'),
                    );
                  },
                ),
              ),
      ),
    );
  }

  Widget _chip(String label, bool selected, VoidCallback onTap) => Padding(
        padding: const EdgeInsets.only(right: 6),
        child: FilterChip(label: Text(label), selected: selected, onSelected: (_) => onTap(), visualDensity: VisualDensity.compact, showCheckmark: false),
      );
}
