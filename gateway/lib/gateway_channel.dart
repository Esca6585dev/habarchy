import 'package:flutter/services.dart';

/// Status snapshot reported by the native service.
class GatewayStatus {
  const GatewayStatus({
    required this.running,
    required this.enabled,
    required this.connected,
    required this.providerName,
    required this.pending,
    required this.sent,
    required this.failed,
    required this.delivered,
    required this.inbound,
    required this.lastPollAt,
    required this.lastHeartbeatAt,
    required this.lastError,
    this.forwardInbound = false,
    this.serverInboundEnabled = true,
    this.inboundActive = false,
  });

  factory GatewayStatus.fromMap(Map<Object?, Object?> m) => GatewayStatus(
        running: m['running'] == true,
        enabled: m['enabled'] == true,
        connected: m['connected'] == true,
        providerName: (m['providerName'] as String?) ?? '',
        pending: (m['pending'] as num?)?.toInt() ?? 0,
        sent: (m['sent'] as num?)?.toInt() ?? 0,
        failed: (m['failed'] as num?)?.toInt() ?? 0,
        delivered: (m['delivered'] as num?)?.toInt() ?? 0,
        inbound: (m['inbound'] as num?)?.toInt() ?? 0,
        lastPollAt: _time(m['lastPollAt']),
        lastHeartbeatAt: _time(m['lastHeartbeatAt']),
        lastError: (m['lastError'] as String?) ?? '',
        forwardInbound: m['forwardInbound'] == true,
        serverInboundEnabled: m['serverInboundEnabled'] != false,
        inboundActive: m['inboundActive'] == true,
      );

  static const empty = GatewayStatus(
    running: false, enabled: false, connected: false, providerName: '', pending: 0, sent: 0, failed: 0, delivered: 0, inbound: 0,
    lastPollAt: null, lastHeartbeatAt: null, lastError: '',
  );

  static DateTime? _time(Object? v) {
    final ms = (v as num?)?.toInt() ?? 0;
    return ms == 0 ? null : DateTime.fromMillisecondsSinceEpoch(ms);
  }

  final bool running;
  final bool enabled;
  final bool connected;
  final String providerName;
  final int pending;
  final int sent;
  final int failed;
  final int delivered;
  final int inbound;
  final DateTime? lastPollAt;
  final DateTime? lastHeartbeatAt;
  final String lastError;
  final bool forwardInbound;
  final bool serverInboundEnabled;
  final bool inboundActive;
}

class GatewayConfig {
  const GatewayConfig({this.url = '', this.key = '', this.enabled = false, this.forwardInbound = false});
  factory GatewayConfig.fromMap(Map<Object?, Object?> m) => GatewayConfig(
        url: (m['url'] as String?) ?? '',
        key: (m['key'] as String?) ?? '',
        enabled: m['enabled'] == true,
        forwardInbound: m['forwardInbound'] == true,
      );
  final String url;
  final String key;
  final bool enabled;
  final bool forwardInbound;
}

/// Thin wrapper around the `habarchy/gateway` MethodChannel.
class GatewayChannel {
  GatewayChannel([MethodChannel? channel]) : _ch = channel ?? const MethodChannel('habarchy/gateway');
  final MethodChannel _ch;

  Future<GatewayConfig> config() async => GatewayConfig.fromMap((await _ch.invokeMethod<Map<Object?, Object?>>('getConfig')) ?? {});
  Future<void> saveConfig({required String url, required String key, required bool forwardInbound}) =>
      _ch.invokeMethod('saveConfig', {'url': url, 'key': key, 'forwardInbound': forwardInbound});
  Future<void> start() => _ch.invokeMethod('start');
  Future<void> stop() => _ch.invokeMethod('stop');
  Future<GatewayStatus> status() async => GatewayStatus.fromMap((await _ch.invokeMethod<Map<Object?, Object?>>('getStatus')) ?? {});
  Future<List<String>> log() async => ((await _ch.invokeMethod<List<Object?>>('getLog')) ?? []).cast<String>();
  Future<void> clearLog() => _ch.invokeMethod('clearLog');
  Future<bool> hasPermissions() async => (await _ch.invokeMethod<bool>('hasPermissions')) ?? false;
  Future<bool> requestPermissions() async => (await _ch.invokeMethod<bool>('requestPermissions')) ?? false;
  Future<bool> isIgnoringBatteryOptimizations() async => (await _ch.invokeMethod<bool>('isIgnoringBatteryOptimizations')) ?? true;
  Future<void> requestIgnoreBatteryOptimizations() => _ch.invokeMethod('requestIgnoreBatteryOptimizations');
  Future<void> testSms(String to, String text) => _ch.invokeMethod('testSms', {'to': to, 'text': text});
  Future<bool> requestInboundPermissions() async => (await _ch.invokeMethod<bool>('requestInboundPermissions')) ?? false;
  Future<bool> hasInboundPermissions() async => (await _ch.invokeMethod<bool>('hasInboundPermissions')) ?? false;
}
