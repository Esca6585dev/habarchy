import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import 'contact_form.dart';

final contactSearchProvider = NotifierProvider<ContactSearch, String>(ContactSearch.new);

class ContactSearch extends Notifier<String> {
  @override
  String build() => '';
  void set(String v) => state = v;
}

final contactsProvider = FutureProvider.autoDispose<List<Contact>>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const [];
  return ref.read(adminRepositoryProvider).contacts(pid, search: ref.watch(contactSearchProvider));
});

class ContactsScreen extends ConsumerWidget {
  const ContactsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final list = ref.watch(contactsProvider);
    final project = ref.watch(currentProjectProvider);
    final editable = project != null && project.role != 'viewer';
    return Scaffold(
      appBar: AppBar(title: Text(t.contacts)),
      floatingActionButton: editable
          ? FloatingActionButton(
              onPressed: () async {
                final body = await showContactForm(context);
                if (body == null || project.id.isEmpty) return;
                try {
                  await ref.read(adminRepositoryProvider).createContact(project.id, body);
                  ref.invalidate(contactsProvider);
                } catch (e) {
                  if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
                }
              },
              child: const Icon(Icons.person_add_alt_1),
            )
          : null,
      body: Column(children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
          child: TextField(
            decoration: InputDecoration(prefixIcon: const Icon(Icons.search), hintText: t.search),
            onSubmitted: (v) => ref.read(contactSearchProvider.notifier).set(v.trim()),
          ),
        ),
        Expanded(
          child: list.when(
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(contactsProvider)),
            data: (rows) => rows.isEmpty
                ? const EmptyState(icon: Icons.person_off_outlined)
                : RefreshIndicator(
                    onRefresh: () => ref.refresh(contactsProvider.future),
                    child: ListView.separated(
                      itemCount: rows.length,
                      separatorBuilder: (_, _) => const Divider(height: 1),
                      itemBuilder: (context, i) {
                        final c = rows[i];
                        return ListTile(
                          leading: CircleAvatar(child: Text(contactLabel(c).substring(0, 1).toUpperCase())),
                          title: Text(contactLabel(c)),
                          subtitle: Text(contactSubtitle(c), maxLines: 1, overflow: TextOverflow.ellipsis),
                          onTap: editable
                              ? () async {
                                  final body = await showContactForm(context, initial: c);
                                  if (body == null) return;
                                  try {
                                    await ref.read(adminRepositoryProvider).updateContact(project.id, c.id, body);
                                    ref.invalidate(contactsProvider);
                                  } catch (e) {
                                    if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
                                  }
                                }
                              : null,
                          trailing: editable
                              ? IconButton(
                                  icon: const Icon(Icons.delete_outline),
                                  onPressed: () async {
                                    final ok = await confirm(context, t.confirmDelete(contactLabel(c)));
                                    if (!ok) return;
                                    await ref.read(adminRepositoryProvider).deleteContact(project.id, c.id);
                                    ref.invalidate(contactsProvider);
                                  },
                                )
                              : null,
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

/// Yes / no dialog.
Future<bool> confirm(BuildContext context, String text) async {
  final t = AppLocalizations.of(context);
  return await showDialog<bool>(
        context: context,
        builder: (ctx) => AlertDialog(
          content: Text(text),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx, false), child: Text(t.no)),
            FilledButton(onPressed: () => Navigator.pop(ctx, true), child: Text(t.yes)),
          ],
        ),
      ) ??
      false;
}
