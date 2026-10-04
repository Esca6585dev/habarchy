import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'core/auth/auth_controller.dart';
import 'features/api_keys/api_keys_screen.dart';
import 'features/dashboard/dashboard_screen.dart';
import 'features/login/login_screen.dart';
import 'features/messages/message_detail_screen.dart';
import 'features/messages/messages_screen.dart';
import 'features/providers/providers_screen.dart';
import 'features/settings/settings_screen.dart';
import 'features/shell/shell_scaffold.dart';
import 'features/templates/templates_screen.dart';

/// Bridges Riverpod state changes to go_router's refreshListenable.
class _AuthListenable extends ChangeNotifier {
  _AuthListenable(Ref ref) {
    ref.listen(authControllerProvider, (_, _) => notifyListeners());
  }
}

final routerProvider = Provider<GoRouter>((ref) {
  final listenable = _AuthListenable(ref);
  ref.onDispose(listenable.dispose);
  return GoRouter(
    initialLocation: '/',
    refreshListenable: listenable,
    redirect: (context, state) {
      final auth = ref.read(authControllerProvider);
      if (auth.isLoading) return null;
      final signedIn = auth.value is AuthSignedIn;
      final onLogin = state.matchedLocation == '/login';
      if (!signedIn && !onLogin) return '/login';
      if (signedIn && onLogin) return '/';
      return null;
    },
    routes: [
      GoRoute(path: '/login', builder: (_, _) => const LoginScreen()),
      StatefulShellRoute.indexedStack(
        builder: (context, state, shell) => ShellScaffold(shell: shell),
        branches: [
          StatefulShellBranch(routes: [GoRoute(path: '/', builder: (_, _) => const DashboardScreen())]),
          StatefulShellBranch(routes: [
            GoRoute(
              path: '/messages',
              builder: (_, _) => const MessagesScreen(),
              routes: [GoRoute(path: ':id', builder: (_, s) => MessageDetailScreen(id: s.pathParameters['id']!))],
            ),
          ]),
          StatefulShellBranch(routes: [GoRoute(path: '/templates', builder: (_, _) => const TemplatesScreen())]),
          StatefulShellBranch(routes: [GoRoute(path: '/providers', builder: (_, _) => const ProvidersScreen())]),
          StatefulShellBranch(routes: [GoRoute(path: '/api-keys', builder: (_, _) => const ApiKeysScreen())]),
          StatefulShellBranch(routes: [GoRoute(path: '/settings', builder: (_, _) => const SettingsScreen())]),
        ],
      ),
    ],
  );
});
