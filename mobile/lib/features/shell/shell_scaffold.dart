import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/auth/auth_controller.dart';
import '../../l10n/app_localizations.dart';

/// Bottom navigation on phones, a navigation rail on tablets.
class ShellScaffold extends ConsumerWidget {
  const ShellScaffold({super.key, required this.shell});
  final StatefulNavigationShell shell;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final wide = MediaQuery.sizeOf(context).width >= 800;
    final items = [
      (Icons.dashboard_outlined, Icons.dashboard, t.dashboard),
      (Icons.forum_outlined, Icons.forum, t.messages),
      (Icons.description_outlined, Icons.description, t.templates),
      (Icons.power_outlined, Icons.power, t.providers),
      (Icons.key_outlined, Icons.key, t.apiKeys),
      (Icons.settings_outlined, Icons.settings, t.settings),
    ];
    void go(int i) => shell.goBranch(i, initialLocation: i == shell.currentIndex);

    final project = ref.watch(currentProjectProvider);
    final body = Column(children: [
      if (project != null) _ProjectBar(projectName: project.name, role: project.role),
      Expanded(child: shell),
    ]);

    if (wide) {
      return Scaffold(
        body: Row(children: [
          NavigationRail(
            selectedIndex: shell.currentIndex,
            onDestinationSelected: go,
            labelType: NavigationRailLabelType.all,
            leading: const Padding(padding: EdgeInsets.symmetric(vertical: 12), child: _Logo()),
            destinations: [for (final i in items) NavigationRailDestination(icon: Icon(i.$1), selectedIcon: Icon(i.$2), label: Text(i.$3))],
          ),
          const VerticalDivider(width: 1),
          Expanded(child: body),
        ]),
      );
    }
    return Scaffold(
      body: body,
      bottomNavigationBar: NavigationBar(
        selectedIndex: shell.currentIndex,
        onDestinationSelected: go,
        labelBehavior: NavigationDestinationLabelBehavior.alwaysHide,
        destinations: [for (final i in items) NavigationDestination(icon: Icon(i.$1), selectedIcon: Icon(i.$2), label: i.$3)],
      ),
    );
  }
}

class _Logo extends StatelessWidget {
  const _Logo();
  @override
  Widget build(BuildContext context) => CircleAvatar(
        backgroundColor: Theme.of(context).colorScheme.primary,
        child: Text('H', style: TextStyle(color: Theme.of(context).colorScheme.onPrimary, fontWeight: FontWeight.bold)),
      );
}

/// Project switcher strip shown under the status bar.
class _ProjectBar extends ConsumerWidget {
  const _ProjectBar({required this.projectName, required this.role});
  final String projectName;
  final String role;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final projects = ref.watch(projectsProvider).value ?? const [];
    final scheme = Theme.of(context).colorScheme;
    return Material(
      color: scheme.surfaceContainerHighest,
      child: SafeArea(
        bottom: false,
        child: InkWell(
          onTap: projects.length > 1
              ? () => showModalBottomSheet<void>(
                    context: context,
                    showDragHandle: true,
                    builder: (ctx) => ListView(
                      shrinkWrap: true,
                      children: [
                        ListTile(title: Text(t.switchProject, style: Theme.of(ctx).textTheme.titleMedium)),
                        for (final p in projects)
                          ListTile(
                            leading: const Icon(Icons.folder_outlined),
                            title: Text(p.name),
                            subtitle: Text('${p.slug} · ${t.role}: ${p.role}'),
                            selected: p.id == ref.read(selectedProjectProvider),
                            onTap: () {
                              ref.read(selectedProjectProvider.notifier).select(p.id);
                              Navigator.of(ctx).pop();
                            },
                          ),
                      ],
                    ),
                  )
              : null,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(children: [
              Icon(Icons.folder_outlined, size: 18, color: scheme.onSurfaceVariant),
              const SizedBox(width: 8),
              Expanded(child: Text(projectName, style: Theme.of(context).textTheme.labelLarge, overflow: TextOverflow.ellipsis)),
              Text(role, style: Theme.of(context).textTheme.labelSmall?.copyWith(color: scheme.onSurfaceVariant)),
              if (projects.length > 1) Icon(Icons.unfold_more, size: 18, color: scheme.onSurfaceVariant),
            ]),
          ),
        ),
      ),
    );
  }
}
