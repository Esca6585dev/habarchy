import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/api/admin_repository.dart';

/// In-memory cache of attachment bytes, keyed by attachment id.
final _attachmentProvider = FutureProvider.autoDispose.family<Uint8List, String>((ref, id) async {
  ref.keepAlive();
  final bytes = await ref.read(adminRepositoryProvider).attachmentBytes(id);
  return Uint8List.fromList(bytes);
});

/// Round avatar: the uploaded image, or initials on a tinted circle.
class ChatAvatar extends ConsumerWidget {
  const ChatAvatar({super.key, this.name, this.avatarId, this.size = 36});
  final String? name;
  final String? avatarId;
  final double size;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final scheme = Theme.of(context).colorScheme;
    final placeholder = CircleAvatar(
      radius: size / 2,
      backgroundColor: scheme.primary.withValues(alpha: 0.15),
      child: Text(_initials(name), style: TextStyle(color: scheme.primary, fontSize: size * 0.38, fontWeight: FontWeight.w600)),
    );
    if (avatarId == null || avatarId!.isEmpty) return placeholder;
    final img = ref.watch(_attachmentProvider(avatarId!));
    return img.when(
      data: (bytes) => CircleAvatar(radius: size / 2, backgroundImage: MemoryImage(bytes)),
      loading: () => placeholder,
      error: (_, _) => placeholder,
    );
  }

  static String _initials(String? name) {
    final parts = (name ?? '').trim().split(RegExp(r'\s+')).where((p) => p.isNotEmpty).take(2);
    final s = parts.map((p) => p[0].toUpperCase()).join();
    return s.isEmpty ? '?' : s;
  }
}

/// Inline image (chat attachment) loaded with auth.
class ChatImage extends ConsumerWidget {
  const ChatImage({super.key, required this.attachmentId});
  final String attachmentId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final img = ref.watch(_attachmentProvider(attachmentId));
    return img.when(
      data: (bytes) => ClipRRect(borderRadius: BorderRadius.circular(12), child: Image.memory(bytes, fit: BoxFit.cover, errorBuilder: (_, _, _) => const SizedBox.shrink())),
      loading: () => const SizedBox(height: 100, child: Center(child: CircularProgressIndicator(strokeWidth: 2))),
      error: (_, _) => const Icon(Icons.broken_image_outlined),
    );
  }
}
