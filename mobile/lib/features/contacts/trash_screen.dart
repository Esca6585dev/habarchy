import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import 'contact_form.dart';
import 'contacts_screen.dart';

final deletedContactsProvider = FutureProvider.autoDispose<({List<Contact> rows, int total})>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return (rows: <Contact>[], total: 0);
  return ref.read(adminRepositoryProvider).deletedContacts(pid);
});

/// Trash: soft-deleted contacts, restorable or permanently removable.
class ContactsTrashScreen extends ConsumerWidget {
  const ContactsTrashScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final list = ref.watch(deletedContactsProvider);
    final project = ref.watch(currentProjectProvider);
    final pid = project?.id ?? '';
    return Scaffold(
      appBar: AppBar(
        title: Text(t.trash),
        actions: [
          if ((list.value?.rows.isNotEmpty ?? false) && pid.isNotEmpty)
            IconButton(
              icon: const Icon(Icons.delete_forever_outlined),
              tooltip: t.purgeAll,
              onPressed: () async {
                if (!await confirm(context, t.purgeAll)) return;
                final n = await ref.read(adminRepositoryProvider).purgeAllContacts(pid);
                ref.invalidate(deletedContactsProvider);
                if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(t.purgedAll(n))));
              },
            ),
        ],
      ),
      body: Column(children: [
        Padding(padding: const EdgeInsets.all(16), child: Text(t.softDeleteHint, style: Theme.of(context).textTheme.bodySmall)),
        Expanded(
          child: list.when(
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(deletedContactsProvider)),
            data: (d) => d.rows.isEmpty
                ? EmptyState(text: t.emptyTrash, icon: Icons.delete_outline)
                : RefreshIndicator(
                    onRefresh: () => ref.refresh(deletedContactsProvider.future),
                    child: ListView.separated(
                      itemCount: d.rows.length,
                      separatorBuilder: (_, _) => const Divider(height: 1),
                      itemBuilder: (context, i) {
                        final c = d.rows[i];
                        return ListTile(
                          leading: const Icon(Icons.person_off_outlined),
                          title: Text(contactLabel(c)),
                          subtitle: Text([contactSubtitle(c), if (c.deletedAt != null) '· ${shortDate(c.deletedAt)}'].where((s) => s.isNotEmpty).join(' ')),
                          trailing: Row(mainAxisSize: MainAxisSize.min, children: [
                            IconButton(
                              icon: const Icon(Icons.restore),
                              tooltip: t.restore,
                              onPressed: () async {
                                try {
                                  await ref.read(adminRepositoryProvider).restoreContact(pid, c.id);
                                  ref.invalidate(deletedContactsProvider);
                                  ref.invalidate(contactsProvider);
                                  if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(t.restored)));
                                } catch (e) {
                                  if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
                                }
                              },
                            ),
                            IconButton(
                              icon: const Icon(Icons.delete_forever_outlined),
                              tooltip: t.purge,
                              onPressed: () async {
                                if (!await confirm(context, '${t.purge}: ${contactLabel(c)}')) return;
                                await ref.read(adminRepositoryProvider).purgeContact(pid, c.id);
                                ref.invalidate(deletedContactsProvider);
                              },
                            ),
                          ]),
                        );
                      },
                    ),
                  ),
          ),
        ),
      ]),
    );
  }
}
