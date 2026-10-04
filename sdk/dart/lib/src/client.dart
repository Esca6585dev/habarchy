import 'dart:convert';
import 'dart:io';

import 'package:crypto/crypto.dart';

import 'models.dart';

/// Pluggable transport (for tests): returns (status, body).
typedef Transport = Future<(int, String)> Function(String method, Uri url, Map<String, String> headers, String body);

/// hex(HMAC-SHA256(apiKey, ts "\n" METHOD "\n" path "\n" hex(sha256(body))))
String habarchySignature(String apiKey, String timestamp, String method, String path, String body) {
  final bodyHash = sha256.convert(utf8.encode(body)).toString();
  final msg = '$timestamp\n${method.toUpperCase()}\n$path\n$bodyHash';
  return Hmac(sha256, utf8.encode(apiKey)).convert(utf8.encode(msg)).toString();
}

/// Habarchy public API client for one project (one API key).
class HabarchyClient {
  HabarchyClient(String baseUrl, this.apiKey, {this.sign = false, Transport? transport, DateTime Function()? now, Duration timeout = const Duration(seconds: 20)})
      : baseUrl = baseUrl.replaceAll(RegExp(r'/+$'), ''),
        _transport = transport ?? _ioTransport(timeout),
        _now = now ?? DateTime.now;

  final String baseUrl;
  final String apiKey;

  /// Adds X-Timestamp / X-Signature to every request.
  final bool sign;
  final Transport _transport;
  final DateTime Function() _now;

  Future<Accepted> sendMessage({
    required String channel,
    required Recipient to,
    String? template,
    Map<String, dynamic>? data,
    String? subject,
    String? title,
    String? body,
    String? locale,
    DateTime? scheduledAt,
    String? priority,
    String? idempotencyKey,
    Map<String, dynamic>? metadata,
  }) async {
    final (d, meta) = await _request('POST', '/api/v1/messages', {
      'channel': channel,
      'to': to.toJson(),
      if (template != null) 'template': template,
      if (data != null) 'data': data,
      if (subject != null) 'subject': subject,
      if (title != null) 'title': title,
      if (body != null) 'body': body,
      if (locale != null) 'locale': locale,
      if (scheduledAt != null) 'scheduled_at': scheduledAt.toUtc().toIso8601String(),
      if (priority != null) 'priority': priority,
      if (idempotencyKey != null) 'idempotency_key': idempotencyKey,
      if (metadata != null) 'metadata': metadata,
    });
    return Accepted.fromJson(d as Map<String, dynamic>, duplicate: meta?['duplicate'] == true);
  }

  Future<Batch> sendBatch({
    required String channel,
    required List<({Recipient to, Map<String, dynamic>? data})> recipients,
    String? template,
    String? subject,
    String? title,
    String? body,
    String? locale,
    String? priority,
    String? idempotencyKey,
  }) async {
    final (d, meta) = await _request('POST', '/api/v1/messages/batch', {
      'channel': channel,
      if (template != null) 'template': template,
      if (subject != null) 'subject': subject,
      if (title != null) 'title': title,
      if (body != null) 'body': body,
      if (locale != null) 'locale': locale,
      if (priority != null) 'priority': priority,
      if (idempotencyKey != null) 'idempotency_key': idempotencyKey,
      'recipients': [
        for (final r in recipients) {'to': r.to.toJson(), if (r.data != null) 'data': r.data},
      ],
    });
    return Batch.fromJson(d as Map<String, dynamic>, meta: meta);
  }

  Future<MessageDetail> getMessage(String id) async {
    final (d, _) = await _request('GET', '/api/v1/messages/${Uri.encodeComponent(id)}');
    return MessageDetail.fromJson(d as Map<String, dynamic>);
  }

  Future<Message> cancelMessage(String id) async {
    final (d, _) = await _request('POST', '/api/v1/messages/${Uri.encodeComponent(id)}/cancel');
    return Message.fromJson(d as Map<String, dynamic>);
  }

  Future<Batch> getBatch(String id) async {
    final (d, _) = await _request('GET', '/api/v1/batches/${Uri.encodeComponent(id)}');
    return Batch.fromJson(d as Map<String, dynamic>);
  }

  Future<OtpSent> sendOtp({required String to, String channel = 'sms', int? length, int? ttl, String? template, String? locale, Map<String, dynamic>? data}) async {
    final (d, _) = await _request('POST', '/api/v1/otp/send', {
      'channel': channel,
      'to': to,
      if (length != null) 'length': length,
      if (ttl != null) 'ttl': ttl,
      if (template != null) 'template': template,
      if (locale != null) 'locale': locale,
      if (data != null) 'data': data,
    });
    return OtpSent.fromJson(d as Map<String, dynamic>);
  }

  Future<OtpVerified> verifyOtp({required String to, required String code}) async {
    final (d, _) = await _request('POST', '/api/v1/otp/verify', {'to': to, 'code': code});
    return OtpVerified.fromJson(d as Map<String, dynamic>);
  }

  /// Registers this device's FCM token (call after login and on token refresh).
  Future<Device> registerDevice({required String token, required String platform, String? appVersion, String? contactId, String? externalId}) async {
    final (d, _) = await _request('POST', '/api/v1/devices', {
      'token': token,
      'platform': platform,
      if (appVersion != null) 'app_version': appVersion,
      if (contactId != null) 'contact_id': contactId,
      if (externalId != null) 'external_id': externalId,
    });
    return Device.fromJson(d as Map<String, dynamic>);
  }

  Future<(Object?, Map<String, dynamic>?)> _request(String method, String path, [Map<String, dynamic>? body]) async {
    final payload = body == null ? '' : jsonEncode(body);
    final headers = <String, String>{'X-Api-Key': apiKey, 'Accept': 'application/json', 'User-Agent': 'habarchy-dart/1.0'};
    if (payload.isNotEmpty) headers['Content-Type'] = 'application/json';
    if (sign) {
      final ts = (_now().millisecondsSinceEpoch ~/ 1000).toString();
      headers['X-Timestamp'] = ts;
      headers['X-Signature'] = habarchySignature(apiKey, ts, method, path, payload);
    }
    final (status, text) = await _transport(method, Uri.parse(baseUrl + path), headers, payload);
    Map<String, dynamic> env = {};
    if (text.isNotEmpty) {
      try {
        env = jsonDecode(text) as Map<String, dynamic>;
      } catch (_) {
        if (status >= 300) throw HabarchyException(status, 'http_error', text.length > 200 ? text.substring(0, 200) : text);
        throw HabarchyException(status, 'bad_response', 'response is not JSON');
      }
    }
    if (status >= 300) {
      final e = env['error'] as Map<String, dynamic>? ?? {'code': 'http_error', 'message': 'HTTP $status'};
      throw HabarchyException(status, e['code'] as String? ?? 'http_error', e['message'] as String? ?? '', e['details'] as Map<String, dynamic>?);
    }
    return (env['data'], env['meta'] as Map<String, dynamic>?);
  }

  static Transport _ioTransport(Duration timeout) {
    final client = HttpClient()..connectionTimeout = timeout;
    return (method, url, headers, body) async {
      final req = await client.openUrl(method, url);
      headers.forEach(req.headers.set);
      if (body.isNotEmpty) req.write(body);
      final res = await req.close().timeout(timeout);
      return (res.statusCode, await res.transform(utf8.decoder).join());
    };
  }
}
