import 'dart:convert';

import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/api/admin_repository.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/models/models.dart';
import '../../core/storage/storage.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';

/// Loads the dashboard; on network failure serves the cached copy.
final dashboardProvider = FutureProvider.autoDispose<({Dashboard data, bool fromCache})>((ref) async {
  final pid = ref.watch(selectedProjectProvider);
  if (pid == null) throw StateError('no project');
  final prefs = ref.read(prefsProvider);
  try {
    final raw = await ref.read(adminRepositoryProvider).dashboardRaw(pid, days: 7);
    await prefs.cacheDashboard(pid, jsonEncode(raw));
    return (data: Dashboard.fromJson(raw), fromCache: false);
  } catch (e) {
    final cached = prefs.cachedDashboard(pid);
    if (cached == null) rethrow;
    return (data: Dashboard.fromJson(jsonDecode(cached) as Map<String, dynamic>), fromCache: true);
  }
});

class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = AppLocalizations.of(context);
    final projects = ref.watch(projectsProvider);
    final pid = ref.watch(selectedProjectProvider);
    if (projects.isLoading) return const Center(child: CircularProgressIndicator());
    if (pid == null) return EmptyState(text: t.noProjects, icon: Icons.folder_off_outlined);
    final dash = ref.watch(dashboardProvider);
    return Scaffold(
      appBar: AppBar(title: Text(t.dashboard), actions: [
        IconButton(icon: const Icon(Icons.refresh), onPressed: () => ref.invalidate(dashboardProvider)),
      ]),
      body: dash.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(dashboardProvider)),
        data: (r) => RefreshIndicator(
          onRefresh: () => ref.refresh(dashboardProvider.future),
          child: ListView(
            padding: const EdgeInsets.all(16),
            children: [
              if (r.fromCache)
                Card(
                  color: Theme.of(context).colorScheme.tertiaryContainer,
                  child: ListTile(leading: const Icon(Icons.cloud_off), title: Text(t.offlineCached), dense: true),
                ),
              _Counters(totals: r.data.totals, latency: r.data.latency),
              const SizedBox(height: 16),
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(t.perDay, style: Theme.of(context).textTheme.titleMedium),
                    const SizedBox(height: 12),
                    SizedBox(height: 180, child: _DailyChart(points: r.data.daily)),
                  ]),
                ),
              ),
              const SizedBox(height: 16),
              Text(t.recentFailures, style: Theme.of(context).textTheme.titleMedium),
              const SizedBox(height: 8),
              if (r.data.recentFailures.isEmpty)
                Card(child: ListTile(leading: const Icon(Icons.check_circle_outline, color: Colors.green), title: Text(t.nothingHere)))
              else
                for (final f in r.data.recentFailures)
                  Card(
                    child: ListTile(
                      leading: Icon(channelIcon(f['Channel'] as String? ?? ''), color: Theme.of(context).colorScheme.error),
                      title: Text(f['ToAddress'] as String? ?? '', style: const TextStyle(fontFamily: 'monospace')),
                      subtitle: Text('${f['ErrorCode'] ?? ''} ${f['ErrorMessage'] ?? ''}', maxLines: 2, overflow: TextOverflow.ellipsis),
                      trailing: const Icon(Icons.chevron_right),
                      onTap: () => context.go('/messages/${f['ID']}'),
                    ),
                  ),
            ],
          ),
        ),
      ),
    );
  }
}

class _Counters extends StatelessWidget {
  const _Counters({required this.totals, required this.latency});
  final Totals totals;
  final Latency latency;
  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    Widget tile(String label, String value, Color? color) => Card(
          child: Padding(
            padding: const EdgeInsets.all(12),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(label.toUpperCase(), style: Theme.of(context).textTheme.labelSmall?.copyWith(color: Theme.of(context).colorScheme.onSurfaceVariant)),
              const SizedBox(height: 4),
              Text(value, style: Theme.of(context).textTheme.headlineSmall?.copyWith(color: color, fontWeight: FontWeight.w600, fontFeatures: const [FontFeature.tabularFigures()])),
            ]),
          ),
        );
    final p95 = latency.p95SentSec < 1 ? '${(latency.p95SentSec * 1000).round()} ms' : '${latency.p95SentSec.toStringAsFixed(1)} s';
    return GridView.count(
      crossAxisCount: MediaQuery.sizeOf(context).width >= 800 ? 4 : 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      childAspectRatio: 1.9,
      mainAxisSpacing: 8,
      crossAxisSpacing: 8,
      children: [
        tile(t.total, '${totals.total}', null),
        tile(t.delivered, '${totals.delivered}', Colors.green.shade600),
        tile(t.failed, '${totals.failed}', totals.failed > 0 ? Theme.of(context).colorScheme.error : null),
        tile('p95', p95, null),
      ],
    );
  }
}

class _DailyChart extends StatelessWidget {
  const _DailyChart({required this.points});
  final List<DailyPoint> points;
  @override
  Widget build(BuildContext context) {
    if (points.isEmpty) return const EmptyState(icon: Icons.bar_chart);
    final days = points.map((p) => p.day).toSet().toList()..sort();
    final channels = points.map((p) => p.channel).toSet().toList()..sort();
    double maxY = 0;
    final groups = <BarChartGroupData>[];
    for (var i = 0; i < days.length; i++) {
      double from = 0;
      final stacks = <BarChartRodStackItem>[];
      for (final c in channels) {
        final v = points.where((p) => p.day == days[i] && p.channel == c).fold<int>(0, (a, p) => a + p.total).toDouble();
        if (v > 0) stacks.add(BarChartRodStackItem(from, from + v, channelColor(c)));
        from += v;
      }
      maxY = from > maxY ? from : maxY;
      groups.add(BarChartGroupData(x: i, barRods: [BarChartRodData(toY: from, rodStackItems: stacks, width: 14, borderRadius: BorderRadius.circular(3))]));
    }
    return BarChart(BarChartData(
      maxY: maxY == 0 ? 1 : maxY * 1.15,
      barGroups: groups,
      gridData: FlGridData(show: true, drawVerticalLine: false, horizontalInterval: maxY <= 5 ? 1 : null),
      borderData: FlBorderData(show: false),
      titlesData: FlTitlesData(
        topTitles: const AxisTitles(),
        rightTitles: const AxisTitles(),
        leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, reservedSize: 28, interval: maxY <= 5 ? 1 : null, getTitlesWidget: (v, _) => Text(v.toInt().toString(), style: const TextStyle(fontSize: 10)))),
        bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, getTitlesWidget: (v, _) {
          final i = v.toInt();
          return Padding(padding: const EdgeInsets.only(top: 4), child: Text(i >= 0 && i < days.length ? days[i].substring(5) : '', style: const TextStyle(fontSize: 10)));
        })),
      ),
      barTouchData: BarTouchData(touchTooltipData: BarTouchTooltipData(getTooltipItem: (g, _, rod, _) => BarTooltipItem('${rod.toY.toInt()}', const TextStyle(fontWeight: FontWeight.bold)))),
    ));
  }
}
