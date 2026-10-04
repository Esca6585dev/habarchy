import 'package:flutter_test/flutter_test.dart';
import 'package:habarchy_admin/core/models/models.dart';
import 'package:habarchy_admin/features/contacts/contact_form.dart';
import 'package:habarchy_admin/features/providers/provider_edit_screen.dart';

void main() {
  test('contact and group models parse API JSON', () {
    final c = Contact.fromJson({'id': 'c1', 'name': 'Aman', 'phone': '+99365123456', 'whatsapp': '', 'tags': ['vip'], 'created_at': '2026-10-04T10:00:00Z'});
    expect(c.name, 'Aman');
    expect(contactLabel(c), 'Aman');
    expect(contactSubtitle(c), '+99365123456');
    final g = Group.fromJson({'id': 'g1', 'name': 'Işdeşler', 'member_count': 3});
    expect(g.memberCount, 3);
    expect(MembersResult.fromJson({'added': 2, 'created_contacts': 1, 'not_found': []}).createdContacts, 1);
  });

  test('inline contact lines are parsed', () {
    expect(parseInlineContact('Aman Amanow, +99365123456'), {'name': 'Aman Amanow', 'phone': '+99365123456'});
    expect(parseInlineContact('maral@example.tm'), {'email': 'maral@example.tm'});
    expect(parseInlineContact('U0123ABCD, Ops'), {'slack_id': 'U0123ABCD', 'name': 'Ops'});
    expect(parseInlineContact('just a name'), isNull);
  });

  test('provider credentials are built per type', () {
    expect(buildCredentials('smtp', {'host': 'smtp.x', 'port': '587', 'tls_mode': 'starttls', 'password': 's'}), {'host': 'smtp.x', 'port': 587, 'tls_mode': 'starttls', 'password': 's'});
    expect(buildCredentials('android_sms', {'sim_slot': '1', 'timeout_sec': '30'}), {'sim_slot': 1, 'timeout_sec': 30});
    final http = buildCredentials('http_sms', {'url': 'https://x', 'method': 'POST', 'headers_json': '{"Authorization":"Bearer t"}', 'success_path': 'status', 'success_equals': 'OK', 'message_id_path': 'id'});
    expect(http['headers'], {'Authorization': 'Bearer t'});
    expect(http['success'], {'json_path': 'status', 'json_equals': 'OK'});
    expect(http['message_id'], {'json_path': 'id'});
    expect(buildCredentials('smpp', {'host': 'h', 'request_dlr': true, 'use_tls': false}), {'host': 'h', 'request_dlr': true, 'use_tls': false});
  });
}
