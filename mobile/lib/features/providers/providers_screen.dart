import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';

final healthProvider = FutureProvider.autoDispose<Health>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const Health();
  return ref.read(adminRepositoryProvider).health(pid);
});

/// Provider status (from /health) with test-send.
class ProvidersScreen extends ConsumerWidget {
  const ProvidersScreen({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final health = ref.watch(healthProvider);
    final project = ref.watch(currentProjectProvider);
    final isAdmin = project != null && (project.role == 'admin' || project.role == 'owner');
    return Scaffold(
      appBar: AppBar(title: Text(t.providers), actions: [IconButton(icon: const Icon(Icons.refresh), onPressed: () => ref.invalidate(healthProvider))]),
      body: health.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(healthProvider)),
        data: (h) => RefreshIndicator(
          onRefresh: () => ref.refresh(healthProvider.future),
          child: ListView(padding: const EdgeInsets.all(16), children: [
            if (h.providers.isEmpty) const EmptyState(icon: Icons.power_off_outlined),
            for (final p in h.providers)
              Card(
                child: ListTile(
                  leading: Icon(channelIcon(p.channel), color: channelColor(p.channel)),
                  title: Text(p.name),
                  subtitle: Text('${p.type} · ✓ ${p.okCount} · ✕ ${p.failedCount}${p.lastSentAt != null ? ' · ${shortDate(p.lastSentAt)}' : ''}'),
                  trailing: StatusChip(p.status),
                  onTap: isAdmin && project.id.isNotEmpty ? () => _testSend(context, ref, project.id, p) : null,
                ),
              ),
            const SizedBox(height: 16),
            Text(t.queues, style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            Card(
              child: Column(children: [
                for (final q in h.queues)
                  ListTile(
                    dense: true,
                    title: Text(q.queue, style: const TextStyle(fontFamily: 'monospace')),
                    subtitle: Text('${t.pending}: ${q.pending} · active: ${q.active} · retry: ${q.retry}'),
                    trailing: Text('${q.processedToday} / ${q.failedToday}', style: TextStyle(color: q.failedToday > 0 ? Theme.of(context).colorScheme.error : null)),
                  ),
                if (h.queues.isEmpty) ListTile(dense: true, title: Text(t.nothingHere)),
              ]),
            ),
          ]),
        ),
      ),
    );
  }

  Future<void> _testSend(BuildContext context, WidgetRef ref, String projectId, ProviderHealth p) async {
    final t = AppLocalizations.of(context);
    final to = TextEditingController();
    Map<String, dynamic>? result;
    await showDialog<void>(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => AlertDialog(
          title: Text('${t.testSend} · ${p.name}'),
          content: Column(mainAxisSize: MainAxisSize.min, children: [
            TextField(controller: to, decoration: InputDecoration(labelText: t.recipient, hintText: '+99365123456')),
            if (result != null) ...[
              const SizedBox(height: 12),
              Align(alignment: Alignment.centerLeft, child: SelectableText(const JsonEncoder.withIndent('  ').convert(result), style: const TextStyle(fontFamily: 'monospace', fontSize: 11))),
            ],
          ]),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx), child: Text(t.cancel)),
            FilledButton(
              onPressed: () async {
                try {
                  final r = await ref.read(adminRepositoryProvider).testSend(projectId, p.id, to.text.trim());
                  setState(() => result = r);
                } catch (e) {
                  setState(() => result = {'error': '$e'});
                }
              },
              child: Text(t.send),
            ),
          ],
        ),
      ),
    );
  }
}
