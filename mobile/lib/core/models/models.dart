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
    @Default('') String bio,
    @JsonKey(name: 'avatar_id') String? avatarId,
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

@freezed
abstract class Contact with _$Contact {
  const factory Contact({
    required String id,
    @JsonKey(name: 'external_id') @Default('') String externalId,
    @Default('') String name,
    @Default('') String phone,
    @Default('') String email,
    @Default('') String whatsapp,
    @JsonKey(name: 'telegram_chat_id') @Default('') String telegramChatId,
    @JsonKey(name: 'slack_id') @Default('') String slackId,
    @Default('tk') String locale,
    @Default(<String>[]) List<String> tags,
    @JsonKey(name: 'created_at') DateTime? createdAt,
    @JsonKey(name: 'deleted_at') DateTime? deletedAt,
  }) = _Contact;
  factory Contact.fromJson(Map<String, dynamic> json) => _$ContactFromJson(json);
}

@freezed
abstract class Group with _$Group {
  const factory Group({
    required String id,
    required String name,
    @Default('') String description,
    @JsonKey(name: 'member_count') @Default(0) int memberCount,
  }) = _Group;
  factory Group.fromJson(Map<String, dynamic> json) => _$GroupFromJson(json);
}

@freezed
abstract class MembersResult with _$MembersResult {
  const factory MembersResult({
    @Default(0) int added,
    @JsonKey(name: 'created_contacts') @Default(0) int createdContacts,
    @JsonKey(name: 'not_found') @Default(<String>[]) List<String> notFound,
  }) = _MembersResult;
  factory MembersResult.fromJson(Map<String, dynamic> json) => _$MembersResultFromJson(json);
}

@freezed
abstract class ProviderDetail with _$ProviderDetail {
  const factory ProviderDetail({
    required String id,
    required String name,
    required String channel,
    required String type,
    @Default(100) int priority,
    @JsonKey(name: 'is_active') @Default(true) bool isActive,
    @JsonKey(name: 'rate_limit_per_sec') @Default(0) int rateLimitPerSec,
    @Default(<String, dynamic>{}) Map<String, dynamic> settings,
  }) = _ProviderDetail;
  factory ProviderDetail.fromJson(Map<String, dynamic> json) => _$ProviderDetailFromJson(json);
}

@freezed
abstract class Pairing with _$Pairing {
  const factory Pairing({
    @JsonKey(name: 'api_url') @Default('') String apiUrl,
    @JsonKey(name: 'gateway_key') @Default('') String gatewayKey,
    @Default('') String qr,
    @Default(false) bool online,
    @JsonKey(name: 'last_seen_at') DateTime? lastSeenAt,
  }) = _Pairing;
  factory Pairing.fromJson(Map<String, dynamic> json) => _$PairingFromJson(json);
}

/// Result of compose/send: the batch plus accepted / rejected counts.
class SendOutcome {
  const SendOutcome({required this.batchId, required this.accepted, required this.rejected});
  final String batchId;
  final int accepted;
  final List<Map<String, dynamic>> rejected;
}

@freezed
abstract class InboundSms with _$InboundSms {
  const factory InboundSms({
    required String id,
    @JsonKey(name: 'provider_id') @Default('') String providerId,
    @Default('') String from,
    @Default('') String text,
    @JsonKey(name: 'received_at') DateTime? receivedAt,
  }) = _InboundSms;
  factory InboundSms.fromJson(Map<String, dynamic> json) => _$InboundSmsFromJson(json);
}

@freezed
abstract class ChatChannel with _$ChatChannel {
  const factory ChatChannel({
    required String id,
    required String kind,
    @Default('') String name,
    @Default('') String topic,
    @Default(0) int unread,
    @JsonKey(name: 'last_body') @Default('') String lastBody,
    @JsonKey(name: 'last_at') DateTime? lastAt,
    @JsonKey(name: 'last_has_file') @Default(false) bool lastHasFile,
    @JsonKey(name: 'peer_id') String? peerId,
    @JsonKey(name: 'peer_name') @Default('') String peerName,
    @JsonKey(name: 'peer_avatar') String? peerAvatar,
  }) = _ChatChannel;
  factory ChatChannel.fromJson(Map<String, dynamic> json) => _$ChatChannelFromJson(json);
}

@freezed
abstract class ChatMessage with _$ChatMessage {
  const factory ChatMessage({
    required String id,
    @JsonKey(name: 'channel_id') required String channelId,
    @JsonKey(name: 'user_id') String? userId,
    @JsonKey(name: 'author_name') @Default('') String authorName,
    @JsonKey(name: 'author_avatar') String? authorAvatar,
    @Default('') String body,
    @JsonKey(name: 'attachment_id') String? attachmentId,
    @JsonKey(name: 'created_at') DateTime? createdAt,
  }) = _ChatMessage;
  factory ChatMessage.fromJson(Map<String, dynamic> json) => _$ChatMessageFromJson(json);
}

@freezed
abstract class ChatMember with _$ChatMember {
  const factory ChatMember({
    required String id,
    @JsonKey(name: 'full_name') @Default('') String fullName,
    @Default('') String email,
    @JsonKey(name: 'avatar_id') String? avatarId,
    @Default('member') String role,
  }) = _ChatMember;
  factory ChatMember.fromJson(Map<String, dynamic> json) => _$ChatMemberFromJson(json);
}
