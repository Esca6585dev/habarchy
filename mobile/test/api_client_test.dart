import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:habarchy_admin/core/api/api_client.dart';
import 'package:habarchy_admin/core/storage/storage.dart';
import 'package:mocktail/mocktail.dart';

class _MockStorage extends Mock implements FlutterSecureStorage {}

/// Fake transport: scripted responses per request index.
class _ScriptedAdapter implements HttpClientAdapter {
  _ScriptedAdapter(this.script);
  final List<ResponseBody Function(RequestOptions)> script;
  final List<RequestOptions> seen = [];
  int i = 0;
  @override
  Future<ResponseBody> fetch(RequestOptions options, Stream<List<int>>? requestStream, Future<void>? cancelFuture) async {
    seen.add(options);
    return script[i++ < script.length ? i - 1 : script.length - 1](options);
  }

  @override
  void close({bool force = false}) {}
}

ResponseBody _json(int status, String body) =>
    ResponseBody.fromString(body, status, headers: {Headers.contentTypeHeader: ['application/json']});

void main() {
  late _MockStorage storage;
  late SecureStore store;
  final tokens = <String, String?>{};

  setUp(() {
    storage = _MockStorage();
    tokens
      ..clear()
      ..['access_token'] = 'old-access'
      ..['refresh_token'] = 'refresh-1';
    when(() => storage.read(key: any(named: 'key'))).thenAnswer((inv) async => tokens[inv.namedArguments[#key]]);
    when(() => storage.write(key: any(named: 'key'), value: any(named: 'value'))).thenAnswer((inv) async {
      tokens[inv.namedArguments[#key] as String] = inv.namedArguments[#value] as String?;
    });
    when(() => storage.delete(key: any(named: 'key'))).thenAnswer((inv) async => tokens.remove(inv.namedArguments[#key]));
    store = SecureStore(storage);
  });

  test('adds the bearer token and unwraps the envelope', () async {
    final adapter = _ScriptedAdapter([(_) => _json(200, '{"data":{"id":"u1","email":"a@b.tm"}}')]);
    final client = ApiClient(baseUrl: 'http://x', store: store, onSessionExpired: () {}, dio: Dio()..httpClientAdapter = adapter);
    final data = await client.request<Map<String, dynamic>>('GET', '/api/admin/me');
    expect(data['email'], 'a@b.tm');
    expect(adapter.seen.single.headers['Authorization'], 'Bearer old-access');
  });

  test('refreshes once on 401 and retries with the new token', () async {
    final adapter = _ScriptedAdapter([
      (_) => _json(401, '{"error":{"code":"unauthorized","message":"expired"}}'),
      (o) {
        expect(o.path, contains('/auth/refresh'));
        return _json(200, '{"data":{"tokens":{"access_token":"new-access","refresh_token":"refresh-2","expires_in":900}}}');
      },
      (o) {
        expect(o.headers['Authorization'], 'Bearer new-access');
        return _json(200, '{"data":[{"id":"p1","name":"P","slug":"p"}]}');
      },
    ]);
    var expired = 0;
    final client = ApiClient(baseUrl: 'http://x', store: store, onSessionExpired: () => expired++, dio: Dio()..httpClientAdapter = adapter);
    final data = await client.request<List<dynamic>>('GET', '/api/admin/projects');
    expect(data, hasLength(1));
    expect(tokens['refresh_token'], 'refresh-2');
    expect(expired, 0);
  });

  test('signals session expiry when the refresh fails', () async {
    final adapter = _ScriptedAdapter([
      (_) => _json(401, '{"error":{"code":"unauthorized","message":"expired"}}'),
      (_) => _json(401, '{"error":{"code":"unauthorized","message":"refresh token reused"}}'),
    ]);
    var expired = 0;
    final client = ApiClient(baseUrl: 'http://x', store: store, onSessionExpired: () => expired++, dio: Dio()..httpClientAdapter = adapter);
    await expectLater(client.request<dynamic>('GET', '/api/admin/me'), throwsA(isA<ApiException>().having((e) => e.isUnauthorized, 'unauthorized', true)));
    expect(expired, 1);
    expect(tokens['access_token'], isNull);
  });

  test('maps error envelopes to ApiException with totp flag', () async {
    final adapter = _ScriptedAdapter([(_) => _json(401, '{"error":{"code":"unauthorized","message":"totp code required","details":{"totp_required":true}}}')]);
    final client = ApiClient(baseUrl: 'http://x', store: store, onSessionExpired: () {}, dio: Dio()..httpClientAdapter = adapter);
    try {
      await client.request<dynamic>('POST', '/api/admin/auth/login', data: {}, noAuth: true);
      fail('expected exception');
    } on ApiException catch (e) {
      expect(e.totpRequired, isTrue);
    }
  });
}
