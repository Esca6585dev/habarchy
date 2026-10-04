import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:habarchy_gateway/gateway_channel.dart';
import 'package:habarchy_gateway/main.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  testWidgets('shows status from the native channel and starts the service', (tester) async {
    const channel = MethodChannel('habarchy/gateway');
    final calls = <String>[];
    var running = false;
    tester.binding.defaultBinaryMessenger.setMockMethodCallHandler(channel, (call) async {
      calls.add(call.method);
      switch (call.method) {
        case 'getConfig':
          return {'url': 'http://10.0.2.2:8080', 'key': 'gw_test_key_0123456789', 'enabled': false, 'forwardInbound': false};
        case 'getStatus':
          return {'running': running, 'enabled': running, 'connected': running, 'providerName': 'Phone', 'pending': 2, 'sent': 5, 'failed': 1, 'delivered': 4, 'inbound': 0, 'lastPollAt': DateTime.now().millisecondsSinceEpoch, 'lastHeartbeatAt': 0, 'lastError': ''};
        case 'getLog':
          return ['12:00:00  connected'];
        case 'hasPermissions':
        case 'isIgnoringBatteryOptimizations':
          return true;
        case 'start':
          running = true;
          return true;
        case 'saveConfig':
          return true;
      }
      return null;
    });

    await tester.pumpWidget(GatewayApp(channel: GatewayChannel(channel)));
    await tester.pumpAndSettle();

    expect(find.text('Durdy'), findsOneWidget); // default locale tk
    expect(find.text('http://10.0.2.2:8080'), findsOneWidget);
    expect(find.textContaining('connected'), findsOneWidget);

    await tester.tap(find.byKey(const Key('toggle')));
    await tester.pumpAndSettle();
    expect(calls, containsAll(['saveConfig', 'start']));
    expect(find.text('Işleýär'), findsOneWidget);
    expect(find.text('Phone'), findsOneWidget);
  });

  test('GatewayStatus parses numbers and timestamps', () {
    final s = GatewayStatus.fromMap({'running': true, 'sent': 3, 'lastPollAt': 0, 'lastHeartbeatAt': 1700000000000});
    expect(s.running, isTrue);
    expect(s.sent, 3);
    expect(s.lastPollAt, isNull);
    expect(s.lastHeartbeatAt, isNotNull);
  });
}
