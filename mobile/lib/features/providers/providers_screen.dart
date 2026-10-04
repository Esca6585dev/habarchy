import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import '../contacts/contacts_screen.dart' show confirm;
import 'provider_edit_screen.dart';

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
      floatingActionButton: isAdmin && project.id.isNotEmpty
          ? FloatingActionButton.extended(
              key: const Key('new-provider'),
              onPressed: () => Navigator.of(context).push(MaterialPageRoute<void>(builder: (_) => ProviderEditScreen(projectId: project.id))),
              icon: const Icon(Icons.add),
              label: Text(t.newProvider),
            )
          : null,
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
                  onTap: isAdmin && project.id.isNotEmpty ? () => _actions(context, ref, project.id, p) : null,
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

  Future<void> _actions(BuildContext context, WidgetRef ref, String projectId, ProviderHealth p) async {
    final t = AppLocalizations.of(context);
    await showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (ctx) => ListView(shrinkWrap: true, children: [
        ListTile(leading: Icon(channelIcon(p.channel), color: channelColor(p.channel)), title: Text(p.name), subtitle: Text(p.type)),
        if (p.type == 'android_sms')
          ListTile(
            leading: const Icon(Icons.phone_android),
            title: Text(t.pairing),
            onTap: () {
              Navigator.pop(ctx);
              showPairing(context, ref, projectId, p.id);
            },
          ),
        ListTile(
          leading: const Icon(Icons.science_outlined),
          title: Text(t.testSend),
          onTap: () {
            Navigator.pop(ctx);
            _testSend(context, ref, projectId, p);
          },
        ),
        ListTile(
          leading: const Icon(Icons.edit_outlined),
          title: Text(t.editProvider),
          onTap: () async {
            Navigator.pop(ctx);
            try {
              final detail = await ref.read(adminRepositoryProvider).provider(projectId, p.id);
              if (context.mounted) await Navigator.of(context).push(MaterialPageRoute<void>(builder: (_) => ProviderEditScreen(projectId: projectId, initial: detail)));
            } catch (e) {
              if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
            }
          },
        ),
        ListTile(
          leading: Icon(Icons.delete_outline, color: Theme.of(ctx).colorScheme.error),
          title: Text(t.delete),
          onTap: () async {
            Navigator.pop(ctx);
            if (!await confirm(context, t.confirmDelete(p.name))) return;
            await ref.read(adminRepositoryProvider).deleteProvider(projectId, p.id);
            ref.invalidate(healthProvider);
          },
        ),
      ]),
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
