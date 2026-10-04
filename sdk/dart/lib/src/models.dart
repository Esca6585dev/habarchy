/// Error returned by the Habarchy API (non-2xx).
class HabarchyException implements Exception {
  HabarchyException(this.status, this.code, this.message, [this.details]);
  final int status;
  final String code;
  final String message;
  final Map<String, dynamic>? details;

  @override
  String toString() => 'HabarchyException($code, $status): $message';
}

/// "to": an address, or a contact reference.
class Recipient {
  const Recipient.address(String this.address)
      : contactId = null,
        externalId = null;
  const Recipient.contact(String this.contactId)
      : address = null,
        externalId = null;
  const Recipient.external(String this.externalId)
      : address = null,
        contactId = null;

  final String? address;
  final String? contactId;
  final String? externalId;

  Object toJson() {
    if (contactId == null && externalId == null) return address ?? '';
    return {
      if (contactId != null) 'contact_id': contactId,
      if (externalId != null) 'external_id': externalId,
      if (address != null) 'address': address,
    };
  }
}

class Accepted {
  Accepted.fromJson(Map<String, dynamic> j, {this.duplicate = false})
      : id = j['id'] as String,
        status = j['status'] as String,
        channel = j['channel'] as String? ?? '',
        to = j['to'] as String? ?? '';
  final String id;
  final String status;
  final String channel;
  final String to;
  final bool duplicate;
}

class Message {
  Message.fromJson(Map<String, dynamic> j)
      : id = j['id'] as String,
        status = j['status'] as String,
        channel = j['channel'] as String? ?? '',
        to = j['to'] as String? ?? '',
        body = j['body'] as String? ?? '',
        errorCode = j['error_code'] as String?,
        errorMessage = j['error_message'] as String?,
        attempts = (j['attempts'] as num?)?.toInt() ?? 0,
        sentAt = _date(j['sent_at']),
        deliveredAt = _date(j['delivered_at']),
        createdAt = _date(j['created_at']) ?? DateTime.now(),
        raw = j;
  final String id;
  final String status;
  final String channel;
  final String to;
  final String body;
  final String? errorCode;
  final String? errorMessage;
  final int attempts;
  final DateTime? sentAt;
  final DateTime? deliveredAt;
  final DateTime createdAt;

  /// Full JSON for fields not mapped here.
  final Map<String, dynamic> raw;
}

class MessageEvent {
  MessageEvent.fromJson(Map<String, dynamic> j)
      : type = j['type'] as String,
        payload = j['payload'],
        createdAt = _date(j['created_at']) ?? DateTime.now();
  final String type;
  final Object? payload;
  final DateTime createdAt;
}

class MessageDetail {
  MessageDetail.fromJson(Map<String, dynamic> j)
      : message = Message.fromJson(j['message'] as Map<String, dynamic>),
        events = ((j['events'] as List<dynamic>?) ?? []).map((e) => MessageEvent.fromJson(e as Map<String, dynamic>)).toList();
  final Message message;
  final List<MessageEvent> events;
}

class Batch {
  Batch.fromJson(Map<String, dynamic> j, {Map<String, dynamic>? meta})
      : id = j['id'] as String,
        status = j['status'] as String,
        total = (j['total'] as num?)?.toInt() ?? 0,
        queued = (j['queued'] as num?)?.toInt() ?? 0,
        sent = (j['sent'] as num?)?.toInt() ?? 0,
        delivered = (j['delivered'] as num?)?.toInt() ?? 0,
        failed = (j['failed'] as num?)?.toInt() ?? 0,
        accepted = (meta?['accepted'] as num?)?.toInt(),
        rejected = (meta?['rejected'] as num?)?.toInt();
  final String id;
  final String status;
  final int total;
  final int queued;
  final int sent;
  final int delivered;
  final int failed;
  final int? accepted;
  final int? rejected;
}

class OtpSent {
  OtpSent.fromJson(Map<String, dynamic> j)
      : messageId = j['message_id'] as String,
        to = j['to'] as String,
        channel = j['channel'] as String,
        expiresIn = (j['expires_in'] as num).toInt(),
        length = (j['length'] as num).toInt();
  final String messageId;
  final String to;
  final String channel;
  final int expiresIn;
  final int length;
}

class OtpVerified {
  OtpVerified.fromJson(Map<String, dynamic> j)
      : verified = j['verified'] == true,
        attemptsRemaining = (j['attempts_remaining'] as num?)?.toInt() ?? 0;
  final bool verified;
  final int attemptsRemaining;
}

class Device {
  Device.fromJson(Map<String, dynamic> j)
      : id = j['id'] as String,
        platform = j['platform'] as String,
        tokenHint = j['token_hint'] as String? ?? '',
        isActive = j['is_active'] == true;
  final String id;
  final String platform;
  final String tokenHint;
  final bool isActive;
}

DateTime? _date(Object? v) => v is String && v.isNotEmpty ? DateTime.tryParse(v) : null;
