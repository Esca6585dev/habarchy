import 'package:dio/dio.dart';
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

  Future<ProviderDetail> provider(String projectId, String id) async =>
      ProviderDetail.fromJson(await _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/providers/$id'));

  Future<ProviderDetail> createProvider(String projectId, Map<String, dynamic> body) async =>
      ProviderDetail.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/providers', data: body));

  Future<ProviderDetail> updateProvider(String projectId, String id, Map<String, dynamic> body) async =>
      ProviderDetail.fromJson(await _api.request<Map<String, dynamic>>('PUT', '/api/admin/projects/$projectId/providers/$id', data: body));

  Future<void> deleteProvider(String projectId, String id) => _api.request<void>('DELETE', '/api/admin/projects/$projectId/providers/$id');

  Future<Pairing> pairing(String projectId, String id) async =>
      Pairing.fromJson(await _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/providers/$id/pairing'));

  // ---- contacts ----
  Future<List<Contact>> contacts(String projectId, {String? search, int limit = 100, int offset = 0}) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/contacts',
              query: {if (search != null && search.isNotEmpty) 'search': search, 'limit': limit, 'offset': offset}))
          .map((e) => Contact.fromJson(e as Map<String, dynamic>))
          .toList();

  Future<Contact> createContact(String projectId, Map<String, dynamic> body) async =>
      Contact.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/contacts', data: body));

  Future<Contact> updateContact(String projectId, String id, Map<String, dynamic> body) async =>
      Contact.fromJson(await _api.request<Map<String, dynamic>>('PUT', '/api/admin/projects/$projectId/contacts/$id', data: body));

  Future<void> deleteContact(String projectId, String id) => _api.request<void>('DELETE', '/api/admin/projects/$projectId/contacts/$id');

  Future<({List<Contact> rows, int total})> deletedContacts(String projectId, {int limit = 100, int offset = 0}) async {
    final (data, meta) = await _api.requestWithMeta<List<dynamic>>('/api/admin/projects/$projectId/contacts/trash', query: {'limit': limit, 'offset': offset});
    return (rows: data.map((e) => Contact.fromJson(e as Map<String, dynamic>)).toList(), total: (meta['total'] as num?)?.toInt() ?? 0);
  }

  Future<Contact> restoreContact(String projectId, String id) async =>
      Contact.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/contacts/$id/restore'));

  Future<void> purgeContact(String projectId, String id) => _api.request<void>('DELETE', '/api/admin/projects/$projectId/contacts/$id/purge');

  Future<int> purgeAllContacts(String projectId) async {
    final data = await _api.request<Map<String, dynamic>>('DELETE', '/api/admin/projects/$projectId/contacts/trash');
    return (data['purged'] as num?)?.toInt() ?? 0;
  }

  // ---- groups ----
  Future<List<Group>> groups(String projectId) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/groups')).map((e) => Group.fromJson(e as Map<String, dynamic>)).toList();

  Future<Group> createGroup(String projectId, String name, {String description = ''}) async =>
      Group.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/groups', data: {'name': name, 'description': description}));

  Future<Group> updateGroup(String projectId, String id, String name, {String description = ''}) async =>
      Group.fromJson(await _api.request<Map<String, dynamic>>('PUT', '/api/admin/projects/$projectId/groups/$id', data: {'name': name, 'description': description}));

  Future<void> deleteGroup(String projectId, String id) => _api.request<void>('DELETE', '/api/admin/projects/$projectId/groups/$id');

  Future<List<Contact>> groupMembers(String projectId, String id) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/groups/$id/members', query: {'limit': 500}))
          .map((e) => Contact.fromJson(e as Map<String, dynamic>))
          .toList();

  Future<MembersResult> addGroupMembers(String projectId, String id, {List<String> contactIds = const [], List<Map<String, dynamic>> contacts = const []}) async =>
      MembersResult.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/groups/$id/members',
          data: {'contact_ids': contactIds, 'contacts': contacts}));

  Future<void> removeGroupMember(String projectId, String id, String contactId) =>
      _api.request<void>('DELETE', '/api/admin/projects/$projectId/groups/$id/members/$contactId');

  Future<List<InboundSms>> inbound(String projectId, {int limit = 100}) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/inbound', query: {'limit': limit})).map((e) => InboundSms.fromJson(e as Map<String, dynamic>)).toList();

  // ---- import ----
  Future<Map<String, dynamic>> importFile(String projectId, {required String filename, required List<int> bytes, String? groupId, bool dryRun = false}) async {
    final form = FormData.fromMap({
      'file': MultipartFile.fromBytes(bytes, filename: filename),
      if (groupId != null && groupId.isNotEmpty) 'group_id': groupId,
      if (dryRun) 'dry_run': 'true',
    });
    return _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/contacts/import', data: form);
  }

  Future<Map<String, dynamic>> importCardDAV(String projectId, {required String username, required String password, String serverUrl = '', String? groupId, bool dryRun = false}) =>
      _api.request<Map<String, dynamic>>('POST', '/api/admin/projects/$projectId/contacts/import/carddav', data: {
        'server_url': serverUrl,
        'username': username,
        'password': password,
        if (groupId != null && groupId.isNotEmpty) 'group_id': groupId,
        'dry_run': dryRun,
      });

  Future<String> googleImportUrl(String projectId, {String? groupId}) async {
    final data = await _api.request<Map<String, dynamic>>('GET', '/api/admin/projects/$projectId/contacts/import/google/url', query: {if (groupId != null && groupId.isNotEmpty) 'group_id': groupId});
    return data['url'] as String;
  }

  // ---- compose ----
  Future<SendOutcome> send(String projectId, Map<String, dynamic> body) async {
    final (data, meta) = await _api.requestWithMeta<Map<String, dynamic>>('/api/admin/projects/$projectId/messages/send', method: 'POST', data: body);
    return SendOutcome(
      batchId: data['id'] as String? ?? '',
      accepted: (meta['accepted'] as num?)?.toInt() ?? 0,
      rejected: ((meta['rejected'] as List<dynamic>?) ?? const []).map((e) => e as Map<String, dynamic>).toList(),
    );
  }

  Future<List<ApiKey>> apiKeys(String projectId) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/projects/$projectId/api-keys')).map((e) => ApiKey.fromJson(e as Map<String, dynamic>)).toList();

  Future<void> revokeApiKey(String projectId, String keyId) => _api.request<void>('DELETE', '/api/admin/projects/$projectId/api-keys/$keyId');

  /// Public API: register this device's FCM token with a project API key.
  Future<Device> registerDevice({required String apiKey, required String token, required String platform, String? externalId, String? appVersion}) async =>
      Device.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/v1/devices',
          data: {'token': token, 'platform': platform, if (externalId != null && externalId.isNotEmpty) 'external_id': externalId, 'app_version': ?appVersion},
          headers: {'X-Api-Key': apiKey},
          noAuth: true));

  // ---- profile + chat ----
  Future<User> updateProfile({required String fullName, String bio = '', String? avatarId}) async =>
      User.fromJson(await _api.request<Map<String, dynamic>>('PUT', '/api/admin/me/profile', data: {'full_name': fullName, 'bio': bio, 'avatar_id': ?avatarId}));

  Future<String> uploadAttachment({required String filename, required List<int> bytes}) async {
    final form = FormData.fromMap({'file': MultipartFile.fromBytes(bytes, filename: filename)});
    final data = await _api.request<Map<String, dynamic>>('POST', '/api/admin/attachments', data: form);
    return data['id'] as String;
  }

  Future<List<int>> attachmentBytes(String id) => _api.bytes('/api/admin/attachments/$id');

  Future<List<ChatChannel>> chatChannels() async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/chat/channels')).map((e) => ChatChannel.fromJson(e as Map<String, dynamic>)).toList();

  Future<List<ChatMessage>> chatMessages(String channelId, {String? before}) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/chat/channels/$channelId/messages', query: {'before': ?before}))
          .map((e) => ChatMessage.fromJson(e as Map<String, dynamic>))
          .toList();

  Future<List<ChatMember>> chatMembers(String channelId) async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/chat/channels/$channelId/members')).map((e) => ChatMember.fromJson(e as Map<String, dynamic>)).toList();

  Future<ChatMessage> postChatMessage(String channelId, {String body = '', String? attachmentId}) async =>
      ChatMessage.fromJson(await _api.request<Map<String, dynamic>>('POST', '/api/admin/chat/channels/$channelId/messages', data: {'body': body, 'attachment_id': ?attachmentId}));

  Future<void> deleteChatMessage(String channelId, String messageId) => _api.request<void>('DELETE', '/api/admin/chat/channels/$channelId/messages/$messageId');

  Future<void> markChatRead(String channelId) => _api.request<void>('POST', '/api/admin/chat/channels/$channelId/read');

  Future<String> createChatChannel({required String kind, required String name, List<String> members = const []}) async {
    final data = await _api.request<Map<String, dynamic>>('POST', '/api/admin/chat/channels', data: {'kind': kind, 'name': name, 'members': members});
    return data['id'] as String;
  }

  Future<String> openDirect(String userId) async {
    final data = await _api.request<Map<String, dynamic>>('POST', '/api/admin/chat/direct', data: {'user_id': userId});
    return data['id'] as String;
  }

  Future<List<User>> users() async =>
      (await _api.request<List<dynamic>>('GET', '/api/admin/users')).map((e) => User.fromJson(e as Map<String, dynamic>)).toList();
}

final adminRepositoryProvider = Provider<AdminRepository>((ref) => AdminRepository(ref.watch(apiClientProvider)));
