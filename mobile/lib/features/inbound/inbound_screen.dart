import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';

final inboundProvider = FutureProvider.autoDispose<List<InboundSms>>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const [];
  return ref.read(adminRepositoryProvider).inbound(pid);
});

/// SMS received by the project's gateway phones.
class InboundScreen extends ConsumerWidget {
  const InboundScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final list = ref.watch(inboundProvider);
    return Scaffold(
      appBar: AppBar(title: Text(t.inboundSms), actions: [IconButton(icon: const Icon(Icons.refresh), onPressed: () => ref.invalidate(inboundProvider))]),
      body: list.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(inboundProvider)),
        data: (rows) => rows.isEmpty
            ? EmptyState(text: t.inboundHint, icon: Icons.move_to_inbox_outlined)
            : RefreshIndicator(
                onRefresh: () => ref.refresh(inboundProvider.future),
                child: ListView.separated(
                  itemCount: rows.length,
                  separatorBuilder: (_, _) => const Divider(height: 1),
                  itemBuilder: (context, i) {
                    final m = rows[i];
                    return ListTile(
                      leading: const Icon(Icons.sms_outlined),
                      title: Text(m.from, style: const TextStyle(fontFamily: 'monospace')),
                      subtitle: Text(m.text),
                      trailing: Text(shortDate(m.receivedAt), style: Theme.of(context).textTheme.labelSmall),
                      isThreeLine: m.text.length > 60,
                    );
                  },
                ),
              ),
      ),
    );
  }
}
