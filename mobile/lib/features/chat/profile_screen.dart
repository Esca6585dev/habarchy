import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';
import '../../l10n/app_localizations.dart';
import '../common/widgets.dart';
import 'auth_image.dart';
import 'conversation_screen.dart' show meUserProvider;

/// Edit the signed-in user's name, bio and photo.
class ProfileScreen extends ConsumerStatefulWidget {
  const ProfileScreen({super.key});
  @override
  ConsumerState<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends ConsumerState<ProfileScreen> {
  final _name = TextEditingController();
  final _bio = TextEditingController();
  String? _avatarId;
  bool _loaded = false;
  bool _saving = false;

  @override
  void dispose() {
    _name.dispose();
    _bio.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    setState(() => _saving = true);
    try {
      await ref.read(adminRepositoryProvider).updateProfile(fullName: _name.text.trim(), bio: _bio.text.trim(), avatarId: _avatarId);
      ref.invalidate(meUserProvider);
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(AppLocalizations.of(context).profileSaved)));
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _pickAvatar() async {
    final picked = await FilePicker.pickFiles(type: FileType.image);
    if (picked.isEmpty) return;
    final f = picked.first;
    final bytes = await f.readAsBytes();
    try {
      final id = await ref.read(adminRepositoryProvider).uploadAttachment(filename: f.name, bytes: bytes);
      setState(() => _avatarId = id);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = AppLocalizations.of(context);
    final me = ref.watch(meUserProvider);
    return Scaffold(
      appBar: AppBar(title: Text(t.profileTitle)),
      body: me.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => ErrorRetry(error: e, onRetry: () => ref.invalidate(meUserProvider)),
        data: (u) {
          if (!_loaded) {
            _name.text = u.fullName;
            _bio.text = u.bio;
            _avatarId = u.avatarId;
            _loaded = true;
          }
          return ListView(padding: const EdgeInsets.all(16), children: [
            Center(
              child: Column(children: [
                ChatAvatar(name: _name.text.isEmpty ? u.email : _name.text, avatarId: _avatarId, size: 96),
                TextButton.icon(onPressed: _pickAvatar, icon: const Icon(Icons.photo_camera_outlined), label: Text(t.uploadPhoto)),
              ]),
            ),
            const SizedBox(height: 8),
            TextField(controller: _name, decoration: InputDecoration(labelText: t.name, border: const OutlineInputBorder()), textCapitalization: TextCapitalization.words),
            const SizedBox(height: 12),
            TextField(controller: _bio, maxLines: 3, decoration: InputDecoration(labelText: t.bio, border: const OutlineInputBorder())),
            const SizedBox(height: 8),
            Text(u.email, style: Theme.of(context).textTheme.bodySmall),
            const SizedBox(height: 16),
            FilledButton(onPressed: _saving ? null : _save, child: Text(t.save)),
          ]);
        },
      ),
    );
  }
}
