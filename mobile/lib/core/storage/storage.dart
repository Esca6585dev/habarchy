import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// Secrets (tokens, push API key) live in the platform keychain / keystore.
class SecureStore {
  SecureStore([FlutterSecureStorage? storage]) : _s = storage ?? const FlutterSecureStorage(aOptions: AndroidOptions());
  final FlutterSecureStorage _s;

  static const _access = 'access_token';
  static const _refresh = 'refresh_token';
  static const _pushKey = 'push_api_key';

  Future<String?> get accessToken => _s.read(key: _access);
  Future<String?> get refreshToken => _s.read(key: _refresh);
  Future<void> saveTokens(String access, String refresh) async {
    await _s.write(key: _access, value: access);
    await _s.write(key: _refresh, value: refresh);
  }

  Future<void> clearTokens() async {
    await _s.delete(key: _access);
    await _s.delete(key: _refresh);
  }

  Future<String?> get pushApiKey => _s.read(key: _pushKey);
  Future<void> savePushApiKey(String? v) => v == null ? _s.delete(key: _pushKey) : _s.write(key: _pushKey, value: v);
}

/// Non-secret preferences and the offline cache.
class Prefs {
  Prefs(this._p);
  final SharedPreferences _p;

  static Future<Prefs> load() async => Prefs(await SharedPreferences.getInstance());

  String? get apiUrl => _p.getString('api_url');
  Future<void> setApiUrl(String v) => _p.setString('api_url', v);

  String? get projectId => _p.getString('project_id');
  Future<void> setProjectId(String? v) => v == null ? _p.remove('project_id') : _p.setString('project_id', v);

  String? get locale => _p.getString('locale');
  Future<void> setLocale(String? v) => v == null ? _p.remove('locale') : _p.setString('locale', v);

  String get themeMode => _p.getString('theme_mode') ?? 'system';
  Future<void> setThemeMode(String v) => _p.setString('theme_mode', v);

  /// Last dashboard payload per project, for offline display.
  String? cachedDashboard(String projectId) => _p.getString('dash_$projectId');
  Future<void> cacheDashboard(String projectId, String json) => _p.setString('dash_$projectId', json);

  String? get pushExternalId => _p.getString('push_external_id');
  Future<void> setPushExternalId(String v) => _p.setString('push_external_id', v);
}

final secureStoreProvider = Provider<SecureStore>((_) => SecureStore());

/// Overridden in main() after SharedPreferences has loaded.
final prefsProvider = Provider<Prefs>((_) => throw UnimplementedError('prefsProvider must be overridden'));
