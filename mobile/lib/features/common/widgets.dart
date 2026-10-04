import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../l10n/app_localizations.dart';

/// Colour for a message / provider / event status.
Color statusColor(BuildContext context, String status) {
  final cs = Theme.of(context).colorScheme;
  switch (status) {
    case 'delivered':
    case 'sent':
    case 'healthy':
    case 'provider.succeeded':
      return Colors.green.shade600;
    case 'processing':
    case 'queued':
    case 'degraded':
    case 'provider.attempt':
    case 'queued.retry':
      return Colors.amber.shade700;
    case 'failed':
    case 'down':
    case 'provider.failed':
      return cs.error;
    case 'cancelled':
    case 'disabled':
      return cs.outline;
    default:
      return cs.primary;
  }
}

/// Localised label for a message status; falls back to the raw value.
String statusLabel(AppLocalizations t, String status) {
  switch (status) {
    case 'queued':
      return t.statusQueued;
    case 'processing':
      return t.statusProcessing;
    case 'sent':
      return t.statusSent;
    case 'delivered':
      return t.statusDelivered;
    case 'failed':
      return t.statusFailed;
    case 'cancelled':
      return t.statusCancelled;
    default:
      return status;
  }
}

class StatusChip extends StatelessWidget {
  const StatusChip(this.status, {super.key});
  final String status;

  @override
  Widget build(BuildContext context) {
    final color = statusColor(context, status);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: color.withValues(alpha: 0.5)),
      ),
      child: Text(
        statusLabel(AppLocalizations.of(context), status),
        style: TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: color),
      ),
    );
  }
}

Color channelColor(String channel) {
  switch (channel) {
    case 'sms':
      return const Color(0xFF0F7A5A);
    case 'email':
      return const Color(0xFF2563EB);
    case 'push':
      return const Color(0xFFD97706);
    case 'telegram':
      return const Color(0xFF0EA5E9);
    default:
      return Colors.grey;
  }
}

IconData channelIcon(String channel) {
  switch (channel) {
    case 'sms':
      return Icons.sms_outlined;
    case 'email':
      return Icons.mail_outline;
    case 'push':
      return Icons.notifications_outlined;
    case 'telegram':
      return Icons.send_outlined;
    default:
      return Icons.hub_outlined;
  }
}

class ChannelChip extends StatelessWidget {
  const ChannelChip(this.channel, {super.key});
  final String channel;

  @override
  Widget build(BuildContext context) {
    final color = channelColor(channel);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(999),
      ),
      child: Row(mainAxisSize: MainAxisSize.min, children: [
        Icon(channelIcon(channel), size: 12, color: color),
        const SizedBox(width: 4),
        Text(channel.toUpperCase(), style: TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: color)),
      ]),
    );
  }
}

class EmptyState extends StatelessWidget {
  const EmptyState({super.key, this.text, this.icon = Icons.inbox_outlined});
  final String? text;
  final IconData icon;

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Icon(icon, size: 40, color: Theme.of(context).colorScheme.outline),
          const SizedBox(height: 12),
          Text(text ?? t.nothingHere, style: Theme.of(context).textTheme.bodyMedium, textAlign: TextAlign.center),
        ]),
      ),
    );
  }
}

class ErrorRetry extends StatelessWidget {
  const ErrorRetry({super.key, required this.error, required this.onRetry});
  final Object error;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Icon(Icons.error_outline, size: 40, color: Theme.of(context).colorScheme.error),
          const SizedBox(height: 12),
          Text(t.error, style: Theme.of(context).textTheme.titleSmall),
          const SizedBox(height: 4),
          Text('$error', style: Theme.of(context).textTheme.bodySmall, textAlign: TextAlign.center, maxLines: 3, overflow: TextOverflow.ellipsis),
          const SizedBox(height: 12),
          FilledButton.tonal(onPressed: onRetry, child: Text(t.retry)),
        ]),
      ),
    );
  }
}

/// `dd.MM HH:mm` in local time, or `—` when null.
String shortDate(DateTime? d) {
  if (d == null) return '—';
  return DateFormat('dd.MM HH:mm').format(d.toLocal());
}

/// Money from micro-units (1 TMT = 1_000_000) formatted with the currency.
String money(int micros, [String currency = 'TMT']) {
  if (micros == 0) return '—';
  final v = micros / 1000000;
  return '${NumberFormat('#,##0.00##').format(v)} $currency';
}
