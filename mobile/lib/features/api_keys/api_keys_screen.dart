import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';

final apiKeysProvider = FutureProvider.autoDispose<List<ApiKey>>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const [];
  return ref.read(adminRepositoryProvider).apiKeys(pid);
});

class ApiKeysScreen extends ConsumerWidget {
  const ApiKeysScreen({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final keys = ref.watch(apiKeysProvider);
    final project = ref.watch(currentProjectProvider);
    final isAdmin = project != null && (project.role == 'admin' || project.role == 'owner');
    return Scaffold(
      appBar: AppBar(title: Text(t.apiKeys)),
      body: keys.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(apiKeysProvider)),
        data: (rows) => rows.isEmpty
            ? const EmptyState(icon: Icons.key_off_outlined)
            : ListView.separated(
                itemCount: rows.length,
                separatorBuilder: (_, _) => const Divider(height: 1),
                itemBuilder: (context, i) {
                  final k = rows[i];
                  final live = k.prefix == 'hb_live_';
                  return ListTile(
                    leading: Icon(k.revokedAt != null ? Icons.key_off : Icons.key, color: k.revokedAt != null ? Colors.grey : live ? Colors.green : Colors.blue),
                    title: Text(k.name),
                    subtitle: Text('${k.prefix}…${k.hint} · ${k.revokedAt != null ? t.revoked : live ? t.live : t.test}\n${k.scopes.join(', ')}'),
                    isThreeLine: true,
                    trailing: isAdmin && k.revokedAt == null
                        ? IconButton(
                            icon: const Icon(Icons.delete_outline),
                            tooltip: t.revoke,
                            onPressed: () async {
                              final ok = await showDialog<bool>(
                                context: context,
                                builder: (ctx) => AlertDialog(
                                  title: Text(t.revoke),
                                  content: Text(k.name),
                                  actions: [TextButton(onPressed: () => Navigator.pop(ctx, false), child: Text(t.no)), FilledButton(onPressed: () => Navigator.pop(ctx, true), child: Text(t.yes))],
                                ),
                              );
                              if (ok != true) return;
                              await ref.read(adminRepositoryProvider).revokeApiKey(project.id, k.id);
                              ref.invalidate(apiKeysProvider);
                            },
                          )
                        : null,
                  );
                },
              ),
      ),
    );
  }
}
