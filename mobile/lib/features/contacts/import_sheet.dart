import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../core/api/admin_repository.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';

/// Import contacts: file (xlsx / csv / docx / vcf), Google Contacts, iCloud / CardDAV.
Future<bool> showImportSheet(BuildContext context, WidgetRef ref, String projectId, {List<Group> groups = const [], String? groupId}) async {
  final t = AppLocalizations.of(context);
  final repo = ref.read(adminRepositoryProvider);
  var selectedGroup = groupId ?? '';
  var changed = false;

  String summary(Map<String, dynamic> r) =>
      t.importResult((r['total'] as num?)?.toInt() ?? 0, (r['created'] as num?)?.toInt() ?? 0, (r['updated'] as num?)?.toInt() ?? 0, (r['skipped'] as num?)?.toInt() ?? 0);

  Future<void> showResult(BuildContext ctx, Map<String, dynamic> r) async {
    final preview = (r['preview'] as List<dynamic>? ?? const []).cast<Map<String, dynamic>>();
    final errors = (r['errors'] as List<dynamic>? ?? const []).cast<Map<String, dynamic>>();
    await showDialog<void>(
      context: ctx,
      builder: (d) => AlertDialog(
        title: Text(r['dry_run'] == true ? t.preview : t.importTitle),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView(shrinkWrap: true, children: [
            Text(summary(r)),
            const SizedBox(height: 8),
            for (final p in preview.take(10))
              Text([p['Name'], p['Phone'], p['Email']].where((v) => v != null && '$v'.isNotEmpty).join(' · '), style: const TextStyle(fontFamily: 'monospace', fontSize: 12)),
            for (final e in errors.take(5)) Text('#${e['row']} ${e['reason']}', style: TextStyle(color: Theme.of(d).colorScheme.error, fontSize: 12)),
          ]),
        ),
        actions: [FilledButton(onPressed: () => Navigator.pop(d), child: const Text('OK'))],
      ),
    );
    if (r['dry_run'] != true) changed = true;
  }

  Future<void> run(BuildContext ctx, Future<Map<String, dynamic>> Function() action) async {
    try {
      final r = await action();
      if (ctx.mounted) await showResult(ctx, r);
    } catch (e) {
      if (ctx.mounted) ScaffoldMessenger.of(ctx).showSnackBar(SnackBar(content: Text('$e')));
    }
  }

  await showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (ctx) => StatefulBuilder(
      builder: (ctx, setState) => Padding(
        padding: EdgeInsets.fromLTRB(16, 0, 16, MediaQuery.viewInsetsOf(ctx).bottom + 16),
        child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(t.importTitle, style: Theme.of(ctx).textTheme.titleMedium),
          const SizedBox(height: 8),
          if (groups.isNotEmpty)
            DropdownButtonFormField<String>(
              initialValue: selectedGroup,
              decoration: InputDecoration(labelText: t.addToGroup),
              items: [DropdownMenuItem(value: '', child: Text(t.noGroup)), for (final g in groups) DropdownMenuItem(value: g.id, child: Text(g.name))],
              onChanged: (v) => setState(() => selectedGroup = v ?? ''),
            ),
          const SizedBox(height: 8),
          ListTile(
            leading: const Icon(Icons.upload_file_outlined),
            title: Text(t.importFile),
            subtitle: Text(t.importFileHint, style: Theme.of(ctx).textTheme.bodySmall),
            onTap: () async {
              final picked = await FilePicker.pickFiles(type: FileType.custom, allowedExtensions: ['xlsx', 'xlsm', 'csv', 'tsv', 'txt', 'vcf', 'docx'], withData: true);
              final f = picked?.files.single;
              if (f == null || f.bytes == null || !ctx.mounted) return;
              await run(ctx, () => repo.importFile(projectId, filename: f.name, bytes: f.bytes!, groupId: selectedGroup));
            },
          ),
          ListTile(
            leading: const Icon(Icons.g_mobiledata, size: 32),
            title: Text(t.importGoogle),
            subtitle: Text(t.importGoogleHint, style: Theme.of(ctx).textTheme.bodySmall),
            onTap: () async {
              try {
                final url = await repo.googleImportUrl(projectId, groupId: selectedGroup);
                await launchUrl(Uri.parse(url), mode: LaunchMode.externalApplication);
                changed = true;
              } catch (e) {
                if (ctx.mounted) ScaffoldMessenger.of(ctx).showSnackBar(SnackBar(content: Text(t.googleNotConfigured)));
              }
            },
          ),
          ListTile(
            leading: const Icon(Icons.apple),
            title: Text(t.importApple),
            subtitle: Text(t.importAppleHint, style: Theme.of(ctx).textTheme.bodySmall),
            onTap: () async {
              final user = TextEditingController();
              final pass = TextEditingController();
              final server = TextEditingController(text: 'https://contacts.icloud.com/');
              final ok = await showDialog<bool>(
                context: ctx,
                builder: (d) => AlertDialog(
                  title: Text(t.importApple),
                  content: Column(mainAxisSize: MainAxisSize.min, children: [
                    TextField(controller: user, decoration: InputDecoration(labelText: t.appleId, hintText: 'you@icloud.com'), keyboardType: TextInputType.emailAddress),
                    const SizedBox(height: 8),
                    TextField(controller: pass, decoration: InputDecoration(labelText: t.appPassword, hintText: 'abcd-efgh-ijkl-mnop'), obscureText: true),
                    const SizedBox(height: 8),
                    TextField(controller: server, decoration: InputDecoration(labelText: t.serverUrl)),
                  ]),
                  actions: [
                    TextButton(onPressed: () => Navigator.pop(d, false), child: Text(t.cancel)),
                    FilledButton(onPressed: () => Navigator.pop(d, true), child: Text(t.importNow)),
                  ],
                ),
              );
              if (ok != true || !ctx.mounted) return;
              await run(ctx, () => repo.importCardDAV(projectId, username: user.text.trim(), password: pass.text, serverUrl: server.text.trim(), groupId: selectedGroup));
            },
          ),
        ]),
      ),
    ),
  );
  return changed;
}
