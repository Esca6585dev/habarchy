import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';

final templatesProvider = FutureProvider.autoDispose<List<Template>>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) return const [];
  return ref.read(adminRepositoryProvider).templates(pid);
});

class TemplatesScreen extends ConsumerWidget {
  const TemplatesScreen({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final list = ref.watch(templatesProvider);
    return Scaffold(
      appBar: AppBar(title: Text(t.templates)),
      body: list.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(templatesProvider)),
        data: (rows) => rows.isEmpty
            ? const EmptyState()
            : RefreshIndicator(
                onRefresh: () => ref.refresh(templatesProvider.future),
                child: ListView.separated(
                  itemCount: rows.length,
                  separatorBuilder: (_, _) => const Divider(height: 1),
                  itemBuilder: (context, i) {
                    final tpl = rows[i];
                    return ListTile(
                      leading: Icon(channelIcon(tpl.channel), color: channelColor(tpl.channel)),
                      title: Text('${tpl.key}  ·  ${tpl.locale.toUpperCase()}  ·  v${tpl.version}'),
                      subtitle: Text(tpl.body, maxLines: 2, overflow: TextOverflow.ellipsis),
                      trailing: tpl.isActive ? null : const Icon(Icons.visibility_off_outlined, size: 18),
                      onTap: () => showModalBottomSheet<void>(context: context, isScrollControlled: true, showDragHandle: true, builder: (_) => _PreviewSheet(tpl: tpl)),
                    );
                  },
                ),
              ),
      ),
    );
  }
}

class _PreviewSheet extends ConsumerStatefulWidget {
  const _PreviewSheet({required this.tpl});
  final Template tpl;
  @override
  ConsumerState<_PreviewSheet> createState() => _PreviewSheetState();
}

class _PreviewSheetState extends ConsumerState<_PreviewSheet> {
  late final TextEditingController _data;
  Preview? _preview;
  String? _error;

  @override
  void initState() {
    super.initState();
    final sample = {for (final v in widget.tpl.requiredVars) v: v == 'code' ? '4821' : v};
    _data = TextEditingController(text: const JsonEncoder.withIndent('  ').convert(sample));
    _render();
  }

  Future<void> _render() async {
    final pid = ref.read(selectedProjectProvider);
    if (pid == null) return;
    try {
      final data = jsonDecode(_data.text) as Map<String, dynamic>;
      final p = await ref.read(adminRepositoryProvider).previewTemplate(pid, widget.tpl.id, data);
      if (mounted) {
        setState(() {
          _preview = p;
          _error = null;
        });
      }
    } catch (e) {
      if (mounted) setState(() => _error = '$e');
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    return Padding(
      padding: EdgeInsets.fromLTRB(16, 0, 16, MediaQuery.viewInsetsOf(context).bottom + 16),
      child: ListView(shrinkWrap: true, children: [
        Row(children: [ChannelChip(widget.tpl.channel), const SizedBox(width: 8), Text(widget.tpl.key, style: Theme.of(context).textTheme.titleMedium)]),
        const SizedBox(height: 8),
        Wrap(spacing: 6, children: [for (final v in widget.tpl.requiredVars) Chip(label: Text(v, style: const TextStyle(fontFamily: 'monospace', fontSize: 12)), visualDensity: VisualDensity.compact)]),
        const SizedBox(height: 8),
        TextField(controller: _data, maxLines: 4, style: const TextStyle(fontFamily: 'monospace', fontSize: 12), decoration: InputDecoration(labelText: t.sampleData), onChanged: (_) => _render()),
        const SizedBox(height: 12),
        Text(t.preview, style: Theme.of(context).textTheme.labelLarge),
        const SizedBox(height: 6),
        if (_error != null) Text(_error!, style: TextStyle(color: Theme.of(context).colorScheme.error)),
        if (_preview != null)
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: Theme.of(context).colorScheme.primaryContainer, borderRadius: BorderRadius.circular(14)),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              if (_preview!.subject.isNotEmpty) Text(_preview!.subject, style: const TextStyle(fontWeight: FontWeight.bold)),
              SelectableText(_preview!.body),
            ]),
          ),
        if (_preview != null && _preview!.missingVars.isNotEmpty)
          Padding(padding: const EdgeInsets.only(top: 6), child: Text('${t.requiredVars}: ${_preview!.missingVars.join(', ')}', style: TextStyle(color: Theme.of(context).colorScheme.error, fontSize: 12))),
      ]),
    );
  }
}
