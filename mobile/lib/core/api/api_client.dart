import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_riverpod/legacy.dart';

import '../config.dart';
import '../models/models.dart';
import '../storage/storage.dart';

/// Thrown for any non-2xx envelope; carries the API error code.
class ApiException implements Exception {
  ApiException(this.statusCode, this.error);
  final int statusCode;
  final ApiError error;
  bool get isUnauthorized => statusCode == 401;
  bool get totpRequired => error.details?['totp_required'] == true;
  @override
  String toString() => '${error.code}: ${error.message}';
}

/// Dio wrapper: base URL, bearer token, single-flight refresh on 401.
class ApiClient {
  ApiClient({required this.baseUrl, required SecureStore store, required this.onSessionExpired, Dio? dio})
      : _store = store, // ignore: prefer_initializing_formals
        dio = dio ?? Dio() {
    this.dio.options
      ..baseUrl = baseUrl
      ..connectTimeout = const Duration(seconds: 15)
      ..receiveTimeout = const Duration(seconds: 30)
      ..headers = {'Accept': 'application/json'}
      ..validateStatus = (_) => true; // we map errors ourselves
    this.dio.interceptors.add(InterceptorsWrapper(onRequest: _onRequest, onResponse: _onResponse));
  }

  final String baseUrl;
  final Dio dio;
  final SecureStore _store;
  final void Function() onSessionExpired;
  Future<bool>? _refreshing;

  Future<void> _onRequest(RequestOptions o, RequestInterceptorHandler h) async {
    if (o.extra['noAuth'] != true) {
      final t = await _store.accessToken;
      if (t != null) o.headers['Authorization'] = 'Bearer $t';
    }
    h.next(o);
  }

  Future<void> _onResponse(Response r, ResponseInterceptorHandler h) async {
    if (r.statusCode == 401 && r.requestOptions.extra['noAuth'] != true && r.requestOptions.extra['retried'] != true) {
      final ok = await _refreshOnce();
      if (ok) {
        final opts = r.requestOptions..extra['retried'] = true;
        final t = await _store.accessToken;
        opts.headers['Authorization'] = 'Bearer $t';
        try {
          final retry = await dio.fetch<dynamic>(opts);
          return h.resolve(retry);
        } on DioException catch (e) {
          return h.reject(e);
        }
      }
      onSessionExpired();
    }
    h.next(r);
  }

  /// Refreshes the token pair; concurrent 401s share one refresh call.
  Future<bool> _refreshOnce() {
    return _refreshing ??= () async {
      try {
        final refresh = await _store.refreshToken;
        if (refresh == null) return false;
        final res = await dio.post<Map<String, dynamic>>(
          '/api/admin/auth/refresh',
          data: {'refresh_token': refresh},
          options: Options(extra: {'noAuth': true}),
        );
        if (res.statusCode != 200 || res.data?['data'] == null) {
          await _store.clearTokens();
          return false;
        }
        final tokens = Tokens.fromJson((res.data!['data'] as Map<String, dynamic>)['tokens'] as Map<String, dynamic>);
        await _store.saveTokens(tokens.accessToken, tokens.refreshToken);
        return true;
      } catch (_) {
        return false;
      } finally {
        _refreshing = null;
      }
    }();
  }

  /// GET/POST/... returning the `data` member of the envelope, or throwing ApiException.
  Future<T> request<T>(String method, String path, {Object? data, Map<String, dynamic>? query, Map<String, String>? headers, bool noAuth = false}) async {
    final res = await dio.request<dynamic>(
      path,
      data: data,
      queryParameters: query,
      options: Options(method: method, headers: headers, extra: {'noAuth': noAuth}),
    );
    final code = res.statusCode ?? 0;
    if (code >= 200 && code < 300) {
      if (res.data is Map<String, dynamic>) return (res.data as Map<String, dynamic>)['data'] as T;
      return null as T;
    }
    final body = res.data;
    final err = body is Map<String, dynamic> && body['error'] is Map<String, dynamic>
        ? ApiError.fromJson(body['error'] as Map<String, dynamic>)
        : ApiError(code: 'http_$code', message: 'HTTP $code');
    throw ApiException(code, err);
  }

  /// Returns data and meta (for paginated lists).
  Future<(T, Map<String, dynamic>)> requestWithMeta<T>(String path, {Map<String, dynamic>? query}) async {
    final res = await dio.get<dynamic>(path, queryParameters: query);
    final code = res.statusCode ?? 0;
    if (code >= 200 && code < 300 && res.data is Map<String, dynamic>) {
      final m = res.data as Map<String, dynamic>;
      return (m['data'] as T, (m['meta'] as Map<String, dynamic>?) ?? const {});
    }
    final body = res.data;
    final err = body is Map<String, dynamic> && body['error'] is Map<String, dynamic>
        ? ApiError.fromJson(body['error'] as Map<String, dynamic>)
        : ApiError(code: 'http_$code', message: 'HTTP $code');
    throw ApiException(code, err);
  }
}

/// Current API base URL (persisted; defaults to the build-time value).
class ApiUrlNotifier extends Notifier<String> {
  @override
  String build() => ref.read(prefsProvider).apiUrl ?? AppConfig.defaultApiUrl;
  Future<void> set(String url) async {
    final clean = url.trim().replaceAll(RegExp(r'/+$'), '');
    await ref.read(prefsProvider).setApiUrl(clean);
    state = clean;
  }
}

final apiUrlProvider = NotifierProvider<ApiUrlNotifier, String>(ApiUrlNotifier.new);

/// Set by the auth controller so the client can sign the user out.
final sessionExpiredProvider = StateProvider<int>((_) => 0);

final apiClientProvider = Provider<ApiClient>((ref) {
  final url = ref.watch(apiUrlProvider);
  return ApiClient(
    baseUrl: url,
    store: ref.read(secureStoreProvider),
    onSessionExpired: () => ref.read(sessionExpiredProvider.notifier).state++,
  );
});
