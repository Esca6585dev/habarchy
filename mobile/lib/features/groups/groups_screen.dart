import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import '../contacts/contact_form.dart';
import '../contacts/contacts_screen.dart';

final groupsProvider = FutureProvider.autoDispose<List<Group>>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const [];
  return ref.read(adminRepositoryProvider).groups(pid);
});

final groupMembersProvider = FutureProvider.autoDispose.family<List<Contact>, String>((ref, groupId) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const [];
  return ref.read(adminRepositoryProvider).groupMembers(pid, groupId);
});

/// Contact groups: colleagues, classmates, customers…
class GroupsScreen extends ConsumerWidget {
  const GroupsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final list = ref.watch(groupsProvider);
    final project = ref.watch(currentProjectProvider);
    final editable = project != null && project.role != 'viewer';
    return Scaffold(
      appBar: AppBar(title: Text(t.groups)),
      floatingActionButton: editable
          ? FloatingActionButton.extended(
              key: const Key('new-group'),
              onPressed: () => _editGroup(context, ref, project.id),
              icon: const Icon(Icons.group_add_outlined),
              label: Text(t.newGroup),
            )
          : null,
      body: list.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(groupsProvider)),
        data: (rows) => rows.isEmpty
            ? EmptyState(text: t.noGroupsYet, icon: Icons.groups_outlined)
            : RefreshIndicator(
                onRefresh: () => ref.refresh(groupsProvider.future),
                child: ListView.builder(
                  padding: const EdgeInsets.fromLTRB(12, 8, 12, 88),
                  itemCount: rows.length,
                  itemBuilder: (context, i) {
                    final g = rows[i];
                    return Card(
                      child: ListTile(
                        leading: const CircleAvatar(child: Icon(Icons.groups_outlined)),
                        title: Text(g.name),
                        subtitle: Text([if (g.description.isNotEmpty) g.description, '${g.memberCount} ${t.members.toLowerCase()}'].join(' · ')),
                        trailing: editable
                            ? IconButton(icon: const Icon(Icons.send_outlined), tooltip: t.sendToGroup, onPressed: () => context.go('/compose?group=${g.id}'))
                            : null,
                        onTap: () => Navigator.of(context).push(MaterialPageRoute<void>(builder: (_) => GroupDetailScreen(group: g))),
                      ),
                    );
                  },
                ),
              ),
      ),
    );
  }
}

Future<void> _editGroup(BuildContext context, WidgetRef ref, String projectId, {Group? initial}) async {
  final t = AppLocalizations.of(context);
  final name = TextEditingController(text: initial?.name ?? '');
  final desc = TextEditingController(text: initial?.description ?? '');
  final ok = await showDialog<bool>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: Text(initial == null ? t.newGroup : t.groupName),
      content: Column(mainAxisSize: MainAxisSize.min, children: [
        TextField(controller: name, key: const Key('group-name'), decoration: InputDecoration(labelText: t.groupName, hintText: 'Işdeşler'), autofocus: true, textCapitalization: TextCapitalization.sentences),
        const SizedBox(height: 8),
        TextField(controller: desc, decoration: InputDecoration(labelText: t.description)),
      ]),
      actions: [
        TextButton(onPressed: () => Navigator.pop(ctx, false), child: Text(t.cancel)),
        FilledButton(onPressed: () => Navigator.pop(ctx, true), child: Text(t.save)),
      ],
    ),
  );
  if (ok != true || name.text.trim().isEmpty) return;
  try {
    final repo = ref.read(adminRepositoryProvider);
    if (initial == null) {
      await repo.createGroup(projectId, name.text.trim(), description: desc.text.trim());
    } else {
      await repo.updateGroup(projectId, initial.id, name.text.trim(), description: desc.text.trim());
    }
    ref.invalidate(groupsProvider);
  } catch (e) {
    if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
  }
}

/// Members of one group with add / remove.
class GroupDetailScreen extends ConsumerWidget {
  const GroupDetailScreen({super.key, required this.group});
  final Group group;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final members = ref.watch(groupMembersProvider(group.id));
    final project = ref.watch(currentProjectProvider);
    final editable = project != null && project.role != 'viewer';
    return Scaffold(
      appBar: AppBar(
        title: Text(group.name),
        actions: [
          if (editable) IconButton(icon: const Icon(Icons.edit_outlined), onPressed: () => _editGroup(context, ref, project.id, initial: group)),
          if (editable)
            IconButton(
              icon: const Icon(Icons.delete_outline),
              onPressed: () async {
                if (!await confirm(context, t.confirmDelete(group.name))) return;
                await ref.read(adminRepositoryProvider).deleteGroup(project.id, group.id);
                ref.invalidate(groupsProvider);
                if (context.mounted) Navigator.pop(context);
              },
            ),
        ],
      ),
      floatingActionButton: editable
          ? Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.end, children: [
              FloatingActionButton.small(heroTag: 'add', onPressed: () => _addMembers(context, ref, project.id), child: const Icon(Icons.person_add_alt_1)),
              const SizedBox(height: 8),
              FloatingActionButton.extended(
                heroTag: 'send',
                onPressed: () {
                  Navigator.pop(context);
                  context.go('/compose?group=${group.id}');
                },
                icon: const Icon(Icons.send_outlined),
                label: Text(t.sendToGroup),
              ),
            ])
          : null,
      body: members.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(groupMembersProvider(group.id))),
        data: (rows) => rows.isEmpty
            ? EmptyState(text: t.noMembers, icon: Icons.person_outline)
            : ListView.separated(
                padding: const EdgeInsets.only(bottom: 120),
                itemCount: rows.length,
                separatorBuilder: (_, _) => const Divider(height: 1),
                itemBuilder: (context, i) {
                  final c = rows[i];
                  return ListTile(
                    leading: CircleAvatar(child: Text(contactLabel(c).substring(0, 1).toUpperCase())),
                    title: Text(contactLabel(c)),
                    subtitle: Text(contactSubtitle(c), maxLines: 1, overflow: TextOverflow.ellipsis),
                    trailing: editable
                        ? IconButton(
                            icon: const Icon(Icons.remove_circle_outline),
                            tooltip: t.remove,
                            onPressed: () async {
                              await ref.read(adminRepositoryProvider).removeGroupMember(project.id, group.id, c.id);
                              ref.invalidate(groupMembersProvider(group.id));
                              ref.invalidate(groupsProvider);
                            },
                          )
                        : null,
                  );
                },
              ),
      ),
    );
  }

  Future<void> _addMembers(BuildContext context, WidgetRef ref, String projectId) async {
    final t = AppLocalizations.of(context);
    final repo = ref.read(adminRepositoryProvider);
    final picked = <String>{};
    final inline = TextEditingController();
    final search = TextEditingController();
    List<Contact> candidates = [];
    try {
      candidates = await repo.contacts(projectId, limit: 100);
    } catch (_) {}
    if (!context.mounted) return;
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.fromLTRB(16, 0, 16, MediaQuery.viewInsetsOf(ctx).bottom + 16),
          child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(t.addMembers, style: Theme.of(ctx).textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(t.pickContacts, style: Theme.of(ctx).textTheme.labelLarge),
            TextField(
              controller: search,
              decoration: InputDecoration(prefixIcon: const Icon(Icons.search), hintText: t.search, isDense: true),
              onSubmitted: (v) async {
                candidates = await repo.contacts(projectId, search: v.trim(), limit: 100);
                setState(() {});
              },
            ),
            const SizedBox(height: 4),
            ConstrainedBox(
              constraints: const BoxConstraints(maxHeight: 180),
              child: candidates.isEmpty
                  ? Padding(padding: const EdgeInsets.all(8), child: Text(t.nothingHere))
                  : ListView(
                      shrinkWrap: true,
                      children: [
                        for (final c in candidates)
                          CheckboxListTile(
                            dense: true,
                            value: picked.contains(c.id),
                            onChanged: (v) => setState(() => v == true ? picked.add(c.id) : picked.remove(c.id)),
                            title: Text(contactLabel(c)),
                            subtitle: Text(contactSubtitle(c), maxLines: 1, overflow: TextOverflow.ellipsis),
                          ),
                      ],
                    ),
            ),
            const SizedBox(height: 8),
            Text(t.newContacts, style: Theme.of(ctx).textTheme.labelLarge),
            TextField(
              controller: inline,
              key: const Key('inline-contacts'),
              maxLines: 4,
              style: const TextStyle(fontFamily: 'monospace', fontSize: 13),
              decoration: InputDecoration(hintText: 'Aman Amanow, +99365123456\nMaral, maral@example.tm', helperText: t.newContactsHint, helperMaxLines: 2),
            ),
            const SizedBox(height: 12),
            Row(mainAxisAlignment: MainAxisAlignment.end, children: [
              TextButton(onPressed: () => Navigator.pop(ctx, false), child: Text(t.cancel)),
              const SizedBox(width: 8),
              FilledButton(key: const Key('add-members'), onPressed: () => Navigator.pop(ctx, true), child: Text(t.addMembers)),
            ]),
          ]),
        ),
      ),
    );
    if (ok != true) return;
    final contacts = inline.text.split('\n').map(parseInlineContact).whereType<Map<String, dynamic>>().toList();
    if (picked.isEmpty && contacts.isEmpty) return;
    try {
      final r = await repo.addGroupMembers(projectId, group.id, contactIds: picked.toList(), contacts: contacts);
      ref.invalidate(groupMembersProvider(group.id));
      ref.invalidate(groupsProvider);
      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(t.added(r.added, r.createdContacts))));
    } catch (e) {
      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    }
  }
}
