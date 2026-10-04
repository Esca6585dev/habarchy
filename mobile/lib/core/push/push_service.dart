import 'dart:async';
import 'dart:io' show Platform;

import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../firebase_options.dart';

/// A push received while the app was open or tapped from the tray.
class ReceivedPush {
  ReceivedPush({required this.title, required this.body, required this.data, required this.at, this.openedFromTray = false});
  final String title;
  final String body;
  final Map<String, dynamic> data;
  final DateTime at;
  final bool openedFromTray;
}

/// Background handler must be a top-level function.
@pragma('vm:entry-point')
Future<void> habarchyBackgroundHandler(RemoteMessage message) async {
  // Nothing to do: the system tray shows notification messages; data-only
  // messages could be persisted here. Keep it cheap.
  debugPrint('[push:bg] ${message.messageId}');
}

/// Wraps Firebase Messaging: initialisation, Android channel, foreground
/// display with flutter_local_notifications, token access.
class PushService {
  PushService();

  static const channel = AndroidNotificationChannel(
    'habarchy_default',
    'Habarchy',
    description: 'Habarchy notifications',
    importance: Importance.high,
  );

  final _local = FlutterLocalNotificationsPlugin();
  final _received = StreamController<ReceivedPush>.broadcast();
  bool _ready = false;

  bool get isAvailable => !kIsWeb && DefaultFirebaseOptions.isConfigured && _ready;
  Stream<ReceivedPush> get received => _received.stream;

  Future<void> init() async {
    if (kIsWeb || !DefaultFirebaseOptions.isConfigured) return;
    try {
      await Firebase.initializeApp(options: DefaultFirebaseOptions.currentPlatform);
      FirebaseMessaging.onBackgroundMessage(habarchyBackgroundHandler);

      const android = AndroidInitializationSettings('@mipmap/ic_launcher');
      const ios = DarwinInitializationSettings(requestAlertPermission: false, requestBadgePermission: false, requestSoundPermission: false);
      await _local.initialize(settings: const InitializationSettings(android: android, iOS: ios));
      await _local.resolvePlatformSpecificImplementation<AndroidFlutterLocalNotificationsPlugin>()?.createNotificationChannel(channel);

      // iOS: show alerts while foregrounded too.
      await FirebaseMessaging.instance.setForegroundNotificationPresentationOptions(alert: true, badge: true, sound: true);

      FirebaseMessaging.onMessage.listen((m) {
        final push = _toPush(m);
        _received.add(push);
        // Android does not show notification messages in the foreground; do it ourselves.
        if (!kIsWeb && Platform.isAndroid && m.notification != null) {
          _local.show(
            id: m.hashCode,
            title: push.title,
            body: push.body,
            notificationDetails: const NotificationDetails(android: AndroidNotificationDetails('habarchy_default', 'Habarchy', importance: Importance.high, priority: Priority.high)),
          );
        }
      });
      FirebaseMessaging.onMessageOpenedApp.listen((m) => _received.add(_toPush(m, opened: true)));
      final initial = await FirebaseMessaging.instance.getInitialMessage();
      if (initial != null) _received.add(_toPush(initial, opened: true));
      _ready = true;
    } catch (e) {
      debugPrint('[push] init failed: $e');
      _ready = false;
    }
  }

  ReceivedPush _toPush(RemoteMessage m, {bool opened = false}) => ReceivedPush(
        title: m.notification?.title ?? (m.data['title'] as String? ?? ''),
        body: m.notification?.body ?? (m.data['body'] as String? ?? ''),
        data: m.data,
        at: DateTime.now(),
        openedFromTray: opened,
      );

  Future<bool> requestPermission() async {
    if (!isAvailable) return false;
    final s = await FirebaseMessaging.instance.requestPermission(alert: true, badge: true, sound: true);
    return s.authorizationStatus == AuthorizationStatus.authorized || s.authorizationStatus == AuthorizationStatus.provisional;
  }

  Future<String?> token() async {
    if (!isAvailable) return null;
    return FirebaseMessaging.instance.getToken();
  }

  Stream<String> get tokenRefresh => isAvailable ? FirebaseMessaging.instance.onTokenRefresh : const Stream.empty();

  String get platformName => kIsWeb ? 'web' : Platform.isIOS ? 'ios' : 'android';
}

final pushServiceProvider = Provider<PushService>((_) => PushService());

/// Rolling list of received pushes for the demo screen.
class ReceivedPushes extends Notifier<List<ReceivedPush>> {
  @override
  List<ReceivedPush> build() {
    final sub = ref.read(pushServiceProvider).received.listen((p) => state = [p, ...state].take(50).toList());
    ref.onDispose(sub.cancel);
    return const [];
  }
}

final receivedPushesProvider = NotifierProvider<ReceivedPushes, List<ReceivedPush>>(ReceivedPushes.new);
