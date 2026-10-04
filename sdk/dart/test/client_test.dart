import 'dart:convert';

import 'package:crypto/crypto.dart';
import 'package:habarchy/habarchy.dart';
import 'package:test/test.dart';

const key = 'hb_live_testkey';

String serverSig(String ts, String method, String path, String body) {
  final bodyHash = sha256.convert(utf8.encode(body)).toString();
  return Hmac(sha256, utf8.encode(key)).convert(utf8.encode('$ts\n$method\n$path\n$bodyHash')).toString();
}

void main() {
  test('signature matches server algorithm', () {
    expect(habarchySignature(key, '1700000000', 'post', '/api/v1/messages', '{"a":1}'), serverSig('1700000000', 'POST', '/api/v1/messages', '{"a":1}'));
  });

  test('sendMessage signs and unwraps', () async {
    late Map<String, String> seenHeaders;
    late String seenBody;
    final c = HabarchyClient('https://h.example.tm/', key, sign: true, now: () => DateTime.fromMillisecondsSinceEpoch(1700000000000), transport: (method, url, headers, body) async {
      expect(url.toString(), 'https://h.example.tm/api/v1/messages');
      seenHeaders = headers;
      seenBody = body;
      return (202, '{"data":{"id":"m1","status":"queued","channel":"sms","to":"+99365123456"},"meta":null,"error":null}');
    });
    final acc = await c.sendMessage(channel: 'sms', to: const Recipient.address('+99365123456'), template: 'otp', data: {'code': '1234'});
    expect(acc.id, 'm1');
    expect(acc.duplicate, isFalse);
    expect(seenHeaders['X-Signature'], serverSig('1700000000', 'POST', '/api/v1/messages', seenBody));
    expect(jsonDecode(seenBody)['to'], '+99365123456');
  });

  test('recipient shapes, errors and batch meta', () async {
    expect(const Recipient.external('u1').toJson(), {'external_id': 'u1'});
    final c = HabarchyClient('https://h.example.tm', key, transport: (method, url, headers, body) async {
      if (url.path.endsWith('/otp/verify')) return (429, '{"error":{"code":"rate_limited","message":"slow","details":{"retry_after":30}}}');
      if (url.path.endsWith('/messages/batch')) return (202, '{"data":{"id":"b1","status":"queued","total":2},"meta":{"accepted":2,"rejected":0}}');
      return (504, 'gateway timeout');
    });
    await expectLater(c.verifyOtp(to: '+993', code: '0000'), throwsA(isA<HabarchyException>().having((e) => e.code, 'code', 'rate_limited').having((e) => e.status, 'status', 429)));
    final b = await c.sendBatch(channel: 'sms', body: 'x', recipients: [(to: const Recipient.address('+1'), data: null), (to: const Recipient.external('u2'), data: {'n': 1})]);
    expect(b.accepted, 2);
    await expectLater(c.getMessage('m1'), throwsA(isA<HabarchyException>().having((e) => e.code, 'code', 'http_error')));
  });
}
