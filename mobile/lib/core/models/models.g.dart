// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_ApiError _$ApiErrorFromJson(Map<String, dynamic> json) => _ApiError(
  code: json['code'] as String,
  message: json['message'] as String,
  details: json['details'] as Map<String, dynamic>?,
);

Map<String, dynamic> _$ApiErrorToJson(_ApiError instance) => <String, dynamic>{
  'code': instance.code,
  'message': instance.message,
  'details': instance.details,
};

_Tokens _$TokensFromJson(Map<String, dynamic> json) => _Tokens(
  accessToken: json['access_token'] as String,
  refreshToken: json['refresh_token'] as String,
  expiresIn: (json['expires_in'] as num?)?.toInt() ?? 900,
);

Map<String, dynamic> _$TokensToJson(_Tokens instance) => <String, dynamic>{
  'access_token': instance.accessToken,
  'refresh_token': instance.refreshToken,
  'expires_in': instance.expiresIn,
};

_User _$UserFromJson(Map<String, dynamic> json) => _User(
  id: json['id'] as String,
  email: json['email'] as String,
  fullName: json['full_name'] as String? ?? '',
  totpEnabled: json['totp_enabled'] as bool? ?? false,
);

Map<String, dynamic> _$UserToJson(_User instance) => <String, dynamic>{
  'id': instance.id,
  'email': instance.email,
  'full_name': instance.fullName,
  'totp_enabled': instance.totpEnabled,
};

_Project _$ProjectFromJson(Map<String, dynamic> json) => _Project(
  id: json['id'] as String,
  name: json['name'] as String,
  slug: json['slug'] as String,
  status: json['status'] as String? ?? 'active',
  role: json['role'] as String? ?? 'viewer',
  dailyQuota: (json['daily_quota'] as num?)?.toInt() ?? 0,
  monthlyQuota: (json['monthly_quota'] as num?)?.toInt() ?? 0,
  defaultLocale: json['default_locale'] as String? ?? 'tk',
);

Map<String, dynamic> _$ProjectToJson(_Project instance) => <String, dynamic>{
  'id': instance.id,
  'name': instance.name,
  'slug': instance.slug,
  'status': instance.status,
  'role': instance.role,
  'daily_quota': instance.dailyQuota,
  'monthly_quota': instance.monthlyQuota,
  'default_locale': instance.defaultLocale,
};

_Message _$MessageFromJson(Map<String, dynamic> json) => _Message(
  id: json['id'] as String,
  status: json['status'] as String,
  channel: json['channel'] as String,
  to: json['to'] as String,
  template: json['template'] as String? ?? '',
  subject: json['subject'] as String? ?? '',
  body: json['body'] as String? ?? '',
  priority: json['priority'] as String? ?? 'normal',
  providerMessageId: json['provider_message_id'] as String? ?? '',
  errorCode: json['error_code'] as String? ?? '',
  errorMessage: json['error_message'] as String? ?? '',
  attempts: (json['attempts'] as num?)?.toInt() ?? 0,
  isTest: json['is_test'] as bool? ?? false,
  costMicros: (json['cost_micros'] as num?)?.toInt() ?? 0,
  currency: json['currency'] as String? ?? 'TMT',
  createdAt: DateTime.parse(json['created_at'] as String),
  sentAt: json['sent_at'] == null
      ? null
      : DateTime.parse(json['sent_at'] as String),
  deliveredAt: json['delivered_at'] == null
      ? null
      : DateTime.parse(json['delivered_at'] as String),
  scheduledAt: json['scheduled_at'] == null
      ? null
      : DateTime.parse(json['scheduled_at'] as String),
  metadata:
      json['metadata'] as Map<String, dynamic>? ?? const <String, dynamic>{},
);

Map<String, dynamic> _$MessageToJson(_Message instance) => <String, dynamic>{
  'id': instance.id,
  'status': instance.status,
  'channel': instance.channel,
  'to': instance.to,
  'template': instance.template,
  'subject': instance.subject,
  'body': instance.body,
  'priority': instance.priority,
  'provider_message_id': instance.providerMessageId,
  'error_code': instance.errorCode,
  'error_message': instance.errorMessage,
  'attempts': instance.attempts,
  'is_test': instance.isTest,
  'cost_micros': instance.costMicros,
  'currency': instance.currency,
  'created_at': instance.createdAt.toIso8601String(),
  'sent_at': instance.sentAt?.toIso8601String(),
  'delivered_at': instance.deliveredAt?.toIso8601String(),
  'scheduled_at': instance.scheduledAt?.toIso8601String(),
  'metadata': instance.metadata,
};

_MessageEvent _$MessageEventFromJson(Map<String, dynamic> json) =>
    _MessageEvent(
      id: json['id'] as String,
      type: json['type'] as String,
      payload:
          json['payload'] as Map<String, dynamic>? ?? const <String, dynamic>{},
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$MessageEventToJson(_MessageEvent instance) =>
    <String, dynamic>{
      'id': instance.id,
      'type': instance.type,
      'payload': instance.payload,
      'created_at': instance.createdAt.toIso8601String(),
    };

_MessageDetail _$MessageDetailFromJson(Map<String, dynamic> json) =>
    _MessageDetail(
      message: Message.fromJson(json['message'] as Map<String, dynamic>),
      events:
          (json['events'] as List<dynamic>?)
              ?.map((e) => MessageEvent.fromJson(e as Map<String, dynamic>))
              .toList() ??
          const <MessageEvent>[],
    );

Map<String, dynamic> _$MessageDetailToJson(_MessageDetail instance) =>
    <String, dynamic>{'message': instance.message, 'events': instance.events};

_Template _$TemplateFromJson(Map<String, dynamic> json) => _Template(
  id: json['id'] as String,
  key: json['key'] as String,
  channel: json['channel'] as String,
  locale: json['locale'] as String? ?? 'tk',
  subject: json['subject'] as String? ?? '',
  body: json['body'] as String? ?? '',
  requiredVars:
      (json['required_vars'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList() ??
      const <String>[],
  version: (json['version'] as num?)?.toInt() ?? 1,
  isActive: json['is_active'] as bool? ?? true,
);

Map<String, dynamic> _$TemplateToJson(_Template instance) => <String, dynamic>{
  'id': instance.id,
  'key': instance.key,
  'channel': instance.channel,
  'locale': instance.locale,
  'subject': instance.subject,
  'body': instance.body,
  'required_vars': instance.requiredVars,
  'version': instance.version,
  'is_active': instance.isActive,
};

_Preview _$PreviewFromJson(Map<String, dynamic> json) => _Preview(
  subject: json['subject'] as String? ?? '',
  body: json['body'] as String? ?? '',
  requiredVars:
      (json['required_vars'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList() ??
      const <String>[],
  missingVars:
      (json['missing_vars'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList() ??
      const <String>[],
);

Map<String, dynamic> _$PreviewToJson(_Preview instance) => <String, dynamic>{
  'subject': instance.subject,
  'body': instance.body,
  'required_vars': instance.requiredVars,
  'missing_vars': instance.missingVars,
};

_ProviderInfo _$ProviderInfoFromJson(Map<String, dynamic> json) =>
    _ProviderInfo(
      id: json['id'] as String,
      name: json['name'] as String,
      channel: json['channel'] as String,
      type: json['type'] as String,
      priority: (json['priority'] as num?)?.toInt() ?? 100,
      isActive: json['is_active'] as bool? ?? true,
      rateLimitPerSec: (json['rate_limit_per_sec'] as num?)?.toInt() ?? 0,
    );

Map<String, dynamic> _$ProviderInfoToJson(_ProviderInfo instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'channel': instance.channel,
      'type': instance.type,
      'priority': instance.priority,
      'is_active': instance.isActive,
      'rate_limit_per_sec': instance.rateLimitPerSec,
    };

_ProviderHealth _$ProviderHealthFromJson(Map<String, dynamic> json) =>
    _ProviderHealth(
      id: json['id'] as String,
      name: json['name'] as String,
      channel: json['channel'] as String,
      type: json['type'] as String,
      status: json['status'] as String? ?? 'idle',
      okCount: (json['ok_count'] as num?)?.toInt() ?? 0,
      failedCount: (json['failed_count'] as num?)?.toInt() ?? 0,
      lastSentAt: json['last_sent_at'] == null
          ? null
          : DateTime.parse(json['last_sent_at'] as String),
    );

Map<String, dynamic> _$ProviderHealthToJson(_ProviderHealth instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'channel': instance.channel,
      'type': instance.type,
      'status': instance.status,
      'ok_count': instance.okCount,
      'failed_count': instance.failedCount,
      'last_sent_at': instance.lastSentAt?.toIso8601String(),
    };

_QueueStats _$QueueStatsFromJson(Map<String, dynamic> json) => _QueueStats(
  queue: json['queue'] as String,
  pending: (json['pending'] as num?)?.toInt() ?? 0,
  active: (json['active'] as num?)?.toInt() ?? 0,
  scheduled: (json['scheduled'] as num?)?.toInt() ?? 0,
  retry: (json['retry'] as num?)?.toInt() ?? 0,
  processedToday: (json['processed_today'] as num?)?.toInt() ?? 0,
  failedToday: (json['failed_today'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$QueueStatsToJson(_QueueStats instance) =>
    <String, dynamic>{
      'queue': instance.queue,
      'pending': instance.pending,
      'active': instance.active,
      'scheduled': instance.scheduled,
      'retry': instance.retry,
      'processed_today': instance.processedToday,
      'failed_today': instance.failedToday,
    };

_Health _$HealthFromJson(Map<String, dynamic> json) => _Health(
  queues:
      (json['queues'] as List<dynamic>?)
          ?.map((e) => QueueStats.fromJson(e as Map<String, dynamic>))
          .toList() ??
      const <QueueStats>[],
  providers:
      (json['providers'] as List<dynamic>?)
          ?.map((e) => ProviderHealth.fromJson(e as Map<String, dynamic>))
          .toList() ??
      const <ProviderHealth>[],
  contacts: (json['contacts'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$HealthToJson(_Health instance) => <String, dynamic>{
  'queues': instance.queues,
  'providers': instance.providers,
  'contacts': instance.contacts,
};

_ApiKey _$ApiKeyFromJson(Map<String, dynamic> json) => _ApiKey(
  id: json['id'] as String,
  name: json['name'] as String,
  prefix: json['prefix'] as String,
  hint: json['hint'] as String,
  key: json['key'] as String?,
  scopes:
      (json['scopes'] as List<dynamic>?)?.map((e) => e as String).toList() ??
      const <String>[],
  lastUsedAt: json['last_used_at'] == null
      ? null
      : DateTime.parse(json['last_used_at'] as String),
  revokedAt: json['revoked_at'] == null
      ? null
      : DateTime.parse(json['revoked_at'] as String),
  createdAt: DateTime.parse(json['created_at'] as String),
);

Map<String, dynamic> _$ApiKeyToJson(_ApiKey instance) => <String, dynamic>{
  'id': instance.id,
  'name': instance.name,
  'prefix': instance.prefix,
  'hint': instance.hint,
  'key': instance.key,
  'scopes': instance.scopes,
  'last_used_at': instance.lastUsedAt?.toIso8601String(),
  'revoked_at': instance.revokedAt?.toIso8601String(),
  'created_at': instance.createdAt.toIso8601String(),
};

_Totals _$TotalsFromJson(Map<String, dynamic> json) => _Totals(
  total: (json['total'] as num?)?.toInt() ?? 0,
  sent: (json['sent'] as num?)?.toInt() ?? 0,
  delivered: (json['delivered'] as num?)?.toInt() ?? 0,
  failed: (json['failed'] as num?)?.toInt() ?? 0,
  pending: (json['pending'] as num?)?.toInt() ?? 0,
  costMicros: (json['cost_micros'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$TotalsToJson(_Totals instance) => <String, dynamic>{
  'total': instance.total,
  'sent': instance.sent,
  'delivered': instance.delivered,
  'failed': instance.failed,
  'pending': instance.pending,
  'cost_micros': instance.costMicros,
};

_DailyPoint _$DailyPointFromJson(Map<String, dynamic> json) => _DailyPoint(
  day: json['day'] as String,
  channel: json['channel'] as String,
  total: (json['total'] as num?)?.toInt() ?? 0,
  sent: (json['sent'] as num?)?.toInt() ?? 0,
  delivered: (json['delivered'] as num?)?.toInt() ?? 0,
  failed: (json['failed'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$DailyPointToJson(_DailyPoint instance) =>
    <String, dynamic>{
      'day': instance.day,
      'channel': instance.channel,
      'total': instance.total,
      'sent': instance.sent,
      'delivered': instance.delivered,
      'failed': instance.failed,
    };

_Latency _$LatencyFromJson(Map<String, dynamic> json) => _Latency(
  p50SentSec: (json['p50_sent_sec'] as num?)?.toDouble() ?? 0,
  p95SentSec: (json['p95_sent_sec'] as num?)?.toDouble() ?? 0,
  p95DeliveredSec: (json['p95_delivered_sec'] as num?)?.toDouble() ?? 0,
  samples: (json['samples'] as num?)?.toInt() ?? 0,
);

Map<String, dynamic> _$LatencyToJson(_Latency instance) => <String, dynamic>{
  'p50_sent_sec': instance.p50SentSec,
  'p95_sent_sec': instance.p95SentSec,
  'p95_delivered_sec': instance.p95DeliveredSec,
  'samples': instance.samples,
};

_Dashboard _$DashboardFromJson(Map<String, dynamic> json) => _Dashboard(
  totals: json['totals'] == null
      ? const Totals()
      : Totals.fromJson(json['totals'] as Map<String, dynamic>),
  latency: json['latency'] == null
      ? const Latency()
      : Latency.fromJson(json['latency'] as Map<String, dynamic>),
  daily:
      (json['daily'] as List<dynamic>?)
          ?.map((e) => DailyPoint.fromJson(e as Map<String, dynamic>))
          .toList() ??
      const <DailyPoint>[],
  recentFailures:
      (json['recent_failures'] as List<dynamic>?)
          ?.map((e) => e as Map<String, dynamic>)
          .toList() ??
      const <Map<String, dynamic>>[],
  byChannel:
      (json['by_channel'] as Map<String, dynamic>?)?.map(
        (k, e) => MapEntry(k, Totals.fromJson(e as Map<String, dynamic>)),
      ) ??
      const <String, Totals>{},
);

Map<String, dynamic> _$DashboardToJson(_Dashboard instance) =>
    <String, dynamic>{
      'totals': instance.totals,
      'latency': instance.latency,
      'daily': instance.daily,
      'recent_failures': instance.recentFailures,
      'by_channel': instance.byChannel,
    };

_Device _$DeviceFromJson(Map<String, dynamic> json) => _Device(
  id: json['id'] as String,
  platform: json['platform'] as String,
  tokenHint: json['token_hint'] as String? ?? '',
  isActive: json['is_active'] as bool? ?? true,
);

Map<String, dynamic> _$DeviceToJson(_Device instance) => <String, dynamic>{
  'id': instance.id,
  'platform': instance.platform,
  'token_hint': instance.tokenHint,
  'is_active': instance.isActive,
};
