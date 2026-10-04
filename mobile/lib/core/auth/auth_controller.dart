import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../api/admin_repository.dart';
import '../api/api_client.dart';
import '../models/models.dart';
import '../storage/storage.dart';

sealed class AuthState {
  const AuthState();
}

class AuthUnknown extends AuthState {
  const AuthUnknown();
}

class AuthSignedOut extends AuthState {
  const AuthSignedOut();
}

class AuthSignedIn extends AuthState {
  const AuthSignedIn(this.user);
  final User user;
}

/// Holds the session. On start it tries the stored tokens; the API client
/// refreshes them transparently and signals expiry via sessionExpiredProvider.
class AuthController extends AsyncNotifier<AuthState> {
  @override
  Future<AuthState> build() async {
    ref.listen<int>(sessionExpiredProvider, (_, _) => state = const AsyncData(AuthSignedOut()));
    final store = ref.read(secureStoreProvider);
    if (await store.accessToken == null && await store.refreshToken == null) return const AuthSignedOut();
    try {
      final user = await ref.read(adminRepositoryProvider).me();
      return AuthSignedIn(user);
    } on ApiException catch (e) {
      if (e.isUnauthorized) {
        await store.clearTokens();
        return const AuthSignedOut();
      }
      rethrow;
    }
  }

  Future<void> login(String email, String password, {String? totpCode}) async {
    final repo = ref.read(adminRepositoryProvider);
    final r = await repo.login(email, password, totpCode: totpCode);
    await ref.read(secureStoreProvider).saveTokens(r.tokens.accessToken, r.tokens.refreshToken);
    state = AsyncData(AuthSignedIn(r.user));
  }

  Future<void> logout() async {
    final store = ref.read(secureStoreProvider);
    final refresh = await store.refreshToken;
    if (refresh != null) {
      try {
        await ref.read(adminRepositoryProvider).logout(refresh);
      } catch (_) {}
    }
    await store.clearTokens();
    await ref.read(prefsProvider).setProjectId(null);
    state = const AsyncData(AuthSignedOut());
  }
}

final authControllerProvider = AsyncNotifierProvider<AuthController, AuthState>(AuthController.new);

/// Projects of the signed-in user.
final projectsProvider = FutureProvider<List<Project>>((ref) async {
  final auth = ref.watch(authControllerProvider).value;
  if (auth is! AuthSignedIn) return const [];
  return ref.read(adminRepositoryProvider).projects();
});

/// Selected project id, persisted; falls back to the first project.
class SelectedProject extends Notifier<String?> {
  @override
  String? build() {
    final saved = ref.read(prefsProvider).projectId;
    final projects = ref.watch(projectsProvider).value;
    if (projects == null || projects.isEmpty) return saved;
    if (saved != null && projects.any((p) => p.id == saved)) return saved;
    return projects.first.id;
  }

  Future<void> select(String id) async {
    await ref.read(prefsProvider).setProjectId(id);
    state = id;
  }
}

final selectedProjectProvider = NotifierProvider<SelectedProject, String?>(SelectedProject.new);

final currentProjectProvider = Provider<Project?>((ref) {
  final id = ref.watch(selectedProjectProvider);
  final list = ref.watch(projectsProvider).value ?? const [];
  for (final p in list) {
    if (p.id == id) return p;
  }
  return null;
});
