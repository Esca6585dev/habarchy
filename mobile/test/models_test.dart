import 'package:flutter_test/flutter_test.dart';
import 'package:habarchy_admin/core/models/models.dart';

void main() {
  test('Message parses the admin JSON shape', () {
    final m = Message.fromJson({
      'id': '01a1', 'status': 'delivered', 'channel': 'sms', 'to': '+99365123456', 'template': 'otp', 'body': 'Kod: 1',
      'attempts': 1, 'is_test': true, 'cost_micros': 0, 'currency': 'TMT', 'metadata': {'k': 'v'},
      'created_at': '2026-10-04T06:48:07Z', 'sent_at': '2026-10-04T06:48:08Z', 'provider_message_id': 'gw-1',
    });
    expect(m.status, 'delivered');
    expect(m.isTest, isTrue);
    expect(m.sentAt, isNotNull);
    expect(m.metadata['k'], 'v');
  });

  test('Dashboard tolerates missing sections', () {
    final d = Dashboard.fromJson({'totals': {'total': 3, 'failed': 1}, 'daily': [{'day': '2026-10-04', 'channel': 'sms', 'total': 3}]});
    expect(d.totals.total, 3);
    expect(d.daily.single.channel, 'sms');
    expect(d.byChannel, isEmpty);
    expect(d.latency.p95SentSec, 0);
  });

  test('ApiError details carry totp_required', () {
    final e = ApiError.fromJson({'code': 'unauthorized', 'message': 'totp code required', 'details': {'totp_required': true}});
    expect(e.details?['totp_required'], true);
  });
}
