import 'package:freezed_annotation/freezed_annotation.dart';

part 'models.freezed.dart';
part 'models.g.dart';

/// Every API response is wrapped in `{data, meta, error}`.
@freezed
abstract class ApiError with _$ApiError {
  const factory ApiError({
    required String code,
    required String message,
    Map<String, dynamic>? details,
  }) = _ApiError;
  factory ApiError.fromJson(Map<String, dynamic> json) => _$ApiErrorFromJson(json);
}

@freezed
abstract class Tokens with _$Tokens {
  const factory Tokens({
    @JsonKey(name: 'access_token') required String accessToken,
    @JsonKey(name: 'refresh_token') required String refreshToken,
    @JsonKey(name: 'expires_in') @Default(900) int expiresIn,
  }) = _Tokens;
  factory Tokens.fromJson(Map<String, dynamic> json) => _$TokensFromJson(json);
}

@freezed
abstract class User with _$User {
  const factory User({
    required String id,
    required String email,
    @JsonKey(name: 'full_name') @Default('') String fullName,
    @JsonKey(name: 'totp_enabled') @Default(false) bool totpEnabled,
  }) = _User;
  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

@freezed
abstract class Project with _$Project {
  const factory Project({
    required String id,
    required String name,
    required String slug,
    @Default('active') String status,
    @Default('viewer') String role,
    @JsonKey(name: 'daily_quota') @Default(0) int dailyQuota,
    @JsonKey(name: 'monthly_quota') @Default(0) int monthlyQuota,
    @JsonKey(name: 'default_locale') @Default('tk') String defaultLocale,
  }) = _Project;
  factory Project.fromJson(Map<String, dynamic> json) => _$ProjectFromJson(json);
}

@freezed
abstract class Message with _$Message {
  const factory Message({
    required String id,
    required String status,
    required String channel,
    required String to,
    @Default('') String template,
    @Default('') String subject,
    @Default('') String body,
    @Default('normal') String priority,
    @JsonKey(name: 'provider_message_id') @Default('') String providerMessageId,
    @JsonKey(name: 'error_code') @Default('') String errorCode,
    @JsonKey(name: 'error_message') @Default('') String errorMessage,
    @Default(0) int attempts,
    @JsonKey(name: 'is_test') @Default(false) bool isTest,
    @JsonKey(name: 'cost_micros') @Default(0) int costMicros,
    @Default('TMT') String currency,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'sent_at') DateTime? sentAt,
    @JsonKey(name: 'delivered_at') DateTime? deliveredAt,
    @JsonKey(name: 'scheduled_at') DateTime? scheduledAt,
    @Default(<String, dynamic>{}) Map<String, dynamic> metadata,
  }) = _Message;
  factory Message.fromJson(Map<String, dynamic> json) => _$MessageFromJson(json);
}

@freezed
abstract class MessageEvent with _$MessageEvent {
  const factory MessageEvent({
    required String id,
    required String type,
    @Default(<String, dynamic>{}) Map<String, dynamic> payload,
    @JsonKey(name: 'created_at') required DateTime createdAt,
  }) = _MessageEvent;
  factory MessageEvent.fromJson(Map<String, dynamic> json) => _$MessageEventFromJson(json);
}

@freezed
abstract class MessageDetail with _$MessageDetail {
  const factory MessageDetail({
    required Message message,
    @Default(<MessageEvent>[]) List<MessageEvent> events,
  }) = _MessageDetail;
  factory MessageDetail.fromJson(Map<String, dynamic> json) => _$MessageDetailFromJson(json);
}

@freezed
abstract class Template with _$Template {
  const factory Template({
    required String id,
    required String key,
    required String channel,
    @Default('tk') String locale,
    @Default('') String subject,
    @Default('') String body,
    @JsonKey(name: 'required_vars') @Default(<String>[]) List<String> requiredVars,
    @Default(1) int version,
    @JsonKey(name: 'is_active') @Default(true) bool isActive,
  }) = _Template;
  factory Template.fromJson(Map<String, dynamic> json) => _$TemplateFromJson(json);
}

@freezed
abstract class Preview with _$Preview {
  const factory Preview({
    @Default('') String subject,
    @Default('') String body,
    @JsonKey(name: 'required_vars') @Default(<String>[]) List<String> requiredVars,
    @JsonKey(name: 'missing_vars') @Default(<String>[]) List<String> missingVars,
  }) = _Preview;
  factory Preview.fromJson(Map<String, dynamic> json) => _$PreviewFromJson(json);
}

@freezed
abstract class ProviderInfo with _$ProviderInfo {
  const factory ProviderInfo({
    required String id,
    required String name,
    required String channel,
    required String type,
    @Default(100) int priority,
    @JsonKey(name: 'is_active') @Default(true) bool isActive,
    @JsonKey(name: 'rate_limit_per_sec') @Default(0) int rateLimitPerSec,
  }) = _ProviderInfo;
  factory ProviderInfo.fromJson(Map<String, dynamic> json) => _$ProviderInfoFromJson(json);
}

@freezed
abstract class ProviderHealth with _$ProviderHealth {
  const factory ProviderHealth({
    required String id,
    required String name,
    required String channel,
    required String type,
    @Default('idle') String status,
    @JsonKey(name: 'ok_count') @Default(0) int okCount,
    @JsonKey(name: 'failed_count') @Default(0) int failedCount,
    @JsonKey(name: 'last_sent_at') DateTime? lastSentAt,
  }) = _ProviderHealth;
  factory ProviderHealth.fromJson(Map<String, dynamic> json) => _$ProviderHealthFromJson(json);
}

@freezed
abstract class QueueStats with _$QueueStats {
  const factory QueueStats({
    required String queue,
    @Default(0) int pending,
    @Default(0) int active,
    @Default(0) int scheduled,
    @Default(0) int retry,
    @JsonKey(name: 'processed_today') @Default(0) int processedToday,
    @JsonKey(name: 'failed_today') @Default(0) int failedToday,
  }) = _QueueStats;
  factory QueueStats.fromJson(Map<String, dynamic> json) => _$QueueStatsFromJson(json);
}

@freezed
abstract class Health with _$Health {
  const factory Health({
    @Default(<QueueStats>[]) List<QueueStats> queues,
    @Default(<ProviderHealth>[]) List<ProviderHealth> providers,
    @Default(0) int contacts,
  }) = _Health;
  factory Health.fromJson(Map<String, dynamic> json) => _$HealthFromJson(json);
}

@freezed
abstract class ApiKey with _$ApiKey {
  const factory ApiKey({
    required String id,
    required String name,
    required String prefix,
    required String hint,
    String? key,
    @Default(<String>[]) List<String> scopes,
    @JsonKey(name: 'last_used_at') DateTime? lastUsedAt,
    @JsonKey(name: 'revoked_at') DateTime? revokedAt,
    @JsonKey(name: 'created_at') required DateTime createdAt,
  }) = _ApiKey;
  factory ApiKey.fromJson(Map<String, dynamic> json) => _$ApiKeyFromJson(json);
}

@freezed
abstract class Totals with _$Totals {
  const factory Totals({
    @Default(0) int total,
    @Default(0) int sent,
    @Default(0) int delivered,
    @Default(0) int failed,
    @Default(0) int pending,
    @JsonKey(name: 'cost_micros') @Default(0) int costMicros,
  }) = _Totals;
  factory Totals.fromJson(Map<String, dynamic> json) => _$TotalsFromJson(json);
}

@freezed
abstract class DailyPoint with _$DailyPoint {
  const factory DailyPoint({
    required String day,
    required String channel,
    @Default(0) int total,
    @Default(0) int sent,
    @Default(0) int delivered,
    @Default(0) int failed,
  }) = _DailyPoint;
  factory DailyPoint.fromJson(Map<String, dynamic> json) => _$DailyPointFromJson(json);
}

@freezed
abstract class Latency with _$Latency {
  const factory Latency({
    @JsonKey(name: 'p50_sent_sec') @Default(0) double p50SentSec,
    @JsonKey(name: 'p95_sent_sec') @Default(0) double p95SentSec,
    @JsonKey(name: 'p95_delivered_sec') @Default(0) double p95DeliveredSec,
    @Default(0) int samples,
  }) = _Latency;
  factory Latency.fromJson(Map<String, dynamic> json) => _$LatencyFromJson(json);
}

@freezed
abstract class Dashboard with _$Dashboard {
  const factory Dashboard({
    @Default(Totals()) Totals totals,
    @Default(Latency()) Latency latency,
    @Default(<DailyPoint>[]) List<DailyPoint> daily,
    @JsonKey(name: 'recent_failures') @Default(<Map<String, dynamic>>[]) List<Map<String, dynamic>> recentFailures,
    @JsonKey(name: 'by_channel') @Default(<String, Totals>{}) Map<String, Totals> byChannel,
  }) = _Dashboard;
  factory Dashboard.fromJson(Map<String, dynamic> json) => _$DashboardFromJson(json);
}

@freezed
abstract class Device with _$Device {
  const factory Device({
    required String id,
    required String platform,
    @JsonKey(name: 'token_hint') @Default('') String tokenHint,
    @JsonKey(name: 'is_active') @Default(true) bool isActive,
  }) = _Device;
  factory Device.fromJson(Map<String, dynamic> json) => _$DeviceFromJson(json);
}
