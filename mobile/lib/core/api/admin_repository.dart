import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../models/models.dart';
import 'api_client.dart';

/// Typed access to the admin API.
class AdminRepository {
  AdminRepository(this._api);
  final ApiClient _api;

  Future<({Tokens tokens, User user})> login(String email, String password, {String? totpCode}) async {
    final data = await _api.request<Map<String, dynamic>>('POST', '/api/admin/auth/login',
        data: {'email': email, 'password': password, if (totpCode != null && totpCode.isNotEmpty) 'totp_code': totpCode}, noAuth: true);
    return (tokens: Tokens.fromJson(data['tokens'] as Map<String, dynamic>), user: User.fromJson(data['user'] as Map<String, dynamic>));
  }

  Future<void> logout(String refreshToken) =>
      _api.request<void>('POST', '/api/admin/auth/logout', data: {'refresh_token': refreshToken}, noAuth: true);

  Future<User> me() async => User.fromJson(await _api.request<Map<String, dynamic>>('GET', '/api/admin/me'));

  Future<List<Project>> projects() async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects')).map((e) => Project.fromJson(e as Map<String, dynamic>)).toList();

  Future<Dashboard> dashboard(String projectId, {int days = 7}) async =>
      Dashboard.fromJson(await _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/dashboard', query: {'days': days}));

  Future<Map<String, dynamic>> dashboardRaw(String projectId, {int days = 7}) =>
      _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/dashboard', query: {'days': days});

  Future<Health> health(String projectId) async => Health.fromJson(await _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/health'));

  Future<({List<Message> rows, String? next})> messages(String projectId, {String? status, String? channel, String? search, String? cursor, int limit = 30}) async {
    final (data, meta) = await _api.requestWithMeta<List<dynamic>>('/api/admin/projects/$projectId/messages', query: {
      'status': ?status,
      'channel': ?channel,
      if (search != null && search.isNotEmpty) 'search': search,
      'cursor': ?cursor,
      'limit': limit,
    });
    return (rows: data.map((e) => Message.fromJson(e as Map<String, dynamic>)).toList(), next: meta['next_cursor'] as String?);
  }

  Future<MessageDetail> message(String projectId, String id) async =>
      MessageDetail.fromJson(await _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/messages/$id'));

  Future<Message> resend(String projectId, String id) async =>
      Message.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/messages/$id/resend'));

  Future<Message> cancel(String projectId, String id) async =>
      Message.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/messages/$id/cancel'));

  Future<List<Template>> templates(String projectId) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/templates')).map((e) => Template.fromJson(e as Map<String, dynamic>)).toList();

  Future<Preview> previewTemplate(String projectId, String templateId, Map<String, dynamic> data) async => Preview.fromJson(
      await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/templates/$templateId/preview', data: {'data': data}));

  Future<List<ProviderInfo>> providers(String projectId) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/providers')).map((e) => ProviderInfo.fromJson(e as Map<String, dynamic>)).toList();

  Future<Map<String, dynamic>> testSend(String projectId, String providerId, String to, {String? text}) =>
      _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/providers/$providerId/test', data: {'to': to, 'text': ?text});

  Future<List<ApiKey>> apiKeys(String projectId) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/api-keys')).map((e) => ApiKey.fromJson(e as Map<String, dynamic>)).toList();

  Future<void> revokeApiKey(String projectId, String keyId) => _api.request<void>('DELETE', '/api/admin/projects/$projectId/api-keys/$keyId');

  /// Public API: register this device's FCM token with a project API key.
  Future<Device> registerDevice({required String apiKey, required String token, required String platform, String? externalId, String? appVersion}) async =>
      Device.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/v1/devices',
          data: {'token': token, 'platform': platform, if (externalId != null && externalId.isNotEmpty) 'external_id': externalId, 'app_version': ?appVersion},
          headers: {'X-Api-Key': apiKey},
          noAuth: true));
}

final adminRepositoryProvider = Provider<AdminRepository>((ref) => AdminRepository(ref.watch(apiClientProvider)));
