import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import 'messages_screen.dart';

final messageDetailProvider = FutureProvider.autoDispose.family<MessageDetail, String>((ref, id) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) throw StateError('no project');
  return ref.read(adminRepositoryProvider).message(pid, id);
});

class MessageDetailScreen extends ConsumerWidget {
  const MessageDetailScreen({super.key, required this.id});
  final String id;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final detail = ref.watch(messageDetailProvider(id));
    final project = ref.watch(currentProjectProvider);
    final canAct = project != null && project.role != 'viewer';
    return Scaffold(
      appBar: AppBar(title: Text(t.messages)),
      body: detail.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(messageDetailProvider(id))),
        data: (d) {
          final m = d.message;
          return ListView(padding: const EdgeInsets.all(16), children: [
            Row(children: [StatusChip(m.status), const SizedBox(width: 8), ChannelChip(m.channel), if (m.isTest) ...[const SizedBox(width: 8), const Text('TEST', style: TextStyle(fontSize: 11, color: Colors.grey))]]),
            const SizedBox(height: 12),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(children: [
                  _row(t.recipient, m.to, mono: true),
                  _row('ID', m.id, mono: true),
                  if (m.template.isNotEmpty) _row(t.templates, m.template),
                  _row(t.attempts, '${m.attempts}'),
                  if (m.providerMessageId.isNotEmpty) _row(t.provider, m.providerMessageId, mono: true),
                  _row(t.sent, shortDate(m.sentAt)),
                  _row(t.delivered, shortDate(m.deliveredAt)),
                  _row(t.cost, money(m.costMicros, m.currency)),
                ]),
              ),
            ),
            if (m.errorCode.isNotEmpty) ...[
              const SizedBox(height: 12),
              Card(color: Theme.of(context).colorScheme.errorContainer, child: ListTile(leading: const Icon(Icons.error_outline), title: Text(m.errorCode), subtitle: Text(m.errorMessage))),
            ],
            const SizedBox(height: 12),
            if (m.subject.isNotEmpty) Text(m.subject, style: Theme.of(context).textTheme.titleMedium),
            Card(child: Padding(padding: const EdgeInsets.all(12), child: SelectableText(m.body))),
            const SizedBox(height: 16),
            Text(t.timeline, style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            for (final e in d.events) _EventTile(e: e),
            const SizedBox(height: 24),
            if (canAct && (m.status == 'failed' || m.status == 'cancelled'))
              FilledButton.icon(
                icon: const Icon(Icons.replay),
                label: Text(t.resend),
                onPressed: () async {
                  final ok = await showDialog<bool>(
                    context: context,
                    builder: (ctx) => AlertDialog(
                      title: Text(t.resend),
                      content: Text(t.confirmResend),
                      actions: [TextButton(onPressed: () => Navigator.pop(ctx, false), child: Text(t.no)), FilledButton(onPressed: () => Navigator.pop(ctx, true), child: Text(t.yes))],
                    ),
                  );
                  if (ok != true) return;
                  try {
                    await ref.read(adminRepositoryProvider).resend(project.id, m.id);
                    ref.invalidate(messageDetailProvider(id));
                    ref.invalidate(messagesPagerProvider);
                  } catch (e) {
                    if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
                  }
                },
              ),
            if (canAct && m.status == 'queued')
              OutlinedButton.icon(
                icon: const Icon(Icons.block),
                label: Text(t.cancel),
                onPressed: () async {
                  try {
                    await ref.read(adminRepositoryProvider).cancel(project.id, m.id);
                    ref.invalidate(messageDetailProvider(id));
                    ref.invalidate(messagesPagerProvider);
                  } catch (e) {
                    if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
                  }
                },
              ),
          ]);
        },
      ),
    );
  }

  Widget _row(String k, String v, {bool mono = false}) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 3),
        child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
          SizedBox(width: 110, child: Text(k, style: const TextStyle(color: Colors.grey))),
          Expanded(child: SelectableText(v, style: mono ? const TextStyle(fontFamily: 'monospace', fontSize: 13) : null)),
        ]),
      );
}

class _EventTile extends StatelessWidget {
  const _EventTile({required this.e});
  final MessageEvent e;
  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final hasPayload = e.payload.isNotEmpty;
    return Card(
      child: hasPayload
          ? ExpansionTile(
              leading: Icon(Icons.circle, size: 10, color: statusColor(context, e.type)),
              title: Text(e.type),
              subtitle: Text(shortDate(e.createdAt)),
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                  child: Align(alignment: Alignment.centerLeft, child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(t.rawResponse, style: Theme.of(context).textTheme.labelSmall),
                    const SizedBox(height: 4),
                    SelectableText(const JsonEncoder.withIndent('  ').convert(e.payload), style: const TextStyle(fontFamily: 'monospace', fontSize: 11)),
                  ])),
                ),
              ],
            )
          : ListTile(leading: Icon(Icons.circle, size: 10, color: statusColor(context, e.type)), title: Text(e.type), subtitle: Text(shortDate(e.createdAt)), dense: true),
    );
  }
}
