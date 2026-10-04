import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/api_client.dart';
import '../../core/auth/auth_controller.dart';
import '../../core/config.dart';
import '../../l10n/app_localizations.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});
  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _form = GlobalKey<FormState>();
  final _email = TextEditingController();
  final _password = TextEditingController();
  final _totp = TextEditingController();
  late final _server = TextEditingController(text: ref.read(apiUrlProvider));
  bool _needTotp = false;
  bool _busy = false;
  bool _showServer = false;
  String? _error;

  @override
  void dispose() {
    _email.dispose();
    _password.dispose();
    _totp.dispose();
    _server.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final t = AppLocalizations.of(context);
    if (!_form.currentState!.validate()) return;
    setState(() {
      _busy = true;
      _error = null;
    });
    try {
      if (_server.text.trim() != ref.read(apiUrlProvider)) await ref.read(apiUrlProvider.notifier).set(_server.text);
      await ref.read(authControllerProvider.notifier).login(_email.text.trim(), _password.text, totpCode: _needTotp ? _totp.text.trim() : null);
    } on ApiException catch (e) {
      setState(() {
        if (e.totpRequired) {
          _needTotp = true;
        } else {
          _error = e.isUnauthorized ? t.invalidCredentials : e.error.message;
        }
      });
    } catch (e) {
      setState(() => _error = '${t.error}: $e');
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 400),
            child: Form(
              key: _form,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Row(children: [
                    CircleAvatar(backgroundColor: Theme.of(context).colorScheme.primary, child: const Text('H', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold))),
                    const SizedBox(width: 12),
                    Text(t.appTitle, style: Theme.of(context).textTheme.headlineSmall),
                    const Spacer(),
                    if (AppConfig.isDev) Chip(label: Text(AppConfig.flavor), visualDensity: VisualDensity.compact),
                  ]),
                  const SizedBox(height: 24),
                  TextFormField(
                    key: const Key('email'),
                    controller: _email,
                    decoration: InputDecoration(labelText: t.email, prefixIcon: const Icon(Icons.mail_outline)),
                    keyboardType: TextInputType.emailAddress,
                    autofillHints: const [AutofillHints.username],
                    textInputAction: TextInputAction.next,
                    validator: (v) => (v == null || !v.contains('@')) ? t.email : null,
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    key: const Key('password'),
                    controller: _password,
                    decoration: InputDecoration(labelText: t.password, prefixIcon: const Icon(Icons.lock_outline)),
                    obscureText: true,
                    autofillHints: const [AutofillHints.password],
                    textInputAction: _needTotp ? TextInputAction.next : TextInputAction.done,
                    onFieldSubmitted: (_) => _needTotp ? null : _submit(),
                    validator: (v) => (v == null || v.isEmpty) ? t.password : null,
                  ),
                  if (_needTotp) ...[
                    const SizedBox(height: 12),
                    TextFormField(
                      key: const Key('totp'),
                      controller: _totp,
                      autofocus: true,
                      decoration: InputDecoration(labelText: t.totpCode, prefixIcon: const Icon(Icons.verified_user_outlined), helperText: t.totpRequired, helperMaxLines: 3),
                      keyboardType: TextInputType.number,
                      maxLength: 6,
                      onFieldSubmitted: (_) => _submit(),
                    ),
                  ],
                  if (_error != null) ...[
                    const SizedBox(height: 12),
                    Text(_error!, style: TextStyle(color: Theme.of(context).colorScheme.error)),
                  ],
                  const SizedBox(height: 16),
                  FilledButton(
                    key: const Key('sign-in'),
                    onPressed: _busy ? null : _submit,
                    child: _busy ? const SizedBox(height: 18, width: 18, child: CircularProgressIndicator(strokeWidth: 2)) : Text(t.signIn),
                  ),
                  const SizedBox(height: 8),
                  TextButton.icon(
                    onPressed: () => setState(() => _showServer = !_showServer),
                    icon: const Icon(Icons.dns_outlined, size: 18),
                    label: Text(_showServer ? t.serverUrl : _server.text, overflow: TextOverflow.ellipsis),
                  ),
                  if (_showServer)
                    TextFormField(
                      controller: _server,
                      decoration: InputDecoration(labelText: t.serverUrl, hintText: 'https://habarchy.example.tm'),
                      keyboardType: TextInputType.url,
                      validator: (v) => (v == null || !v.startsWith('http')) ? 'http(s)://' : null,
                    ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
