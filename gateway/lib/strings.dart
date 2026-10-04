import 'package:flutter/widgets.dart';

/// Tiny tk / ru / en string table (the app has a single screen).
class S {
  S(this.code);
  final String code;

  static S of(BuildContext context) {
    final lang = Localizations.maybeLocaleOf(context)?.languageCode ?? 'tk';
    return S(const ['tk', 'ru', 'en'].contains(lang) ? lang : 'tk');
  }

  static const _t = <String, Map<String, String>>{
    'title': {'tk': 'Habarçy SMS Gateway', 'ru': 'Habarçy SMS-шлюз', 'en': 'Habarchy SMS Gateway'},
    'intro': {
      'tk': 'Bu telefon Habarçy üçin SMS iberiji bolup işleýär: serwerden gelen habarlary SIM-kartaňyzdan iberýär.',
      'ru': 'Этот телефон отправляет SMS для Habarçy: сообщения с сервера уходят через вашу SIM-карту.',
      'en': 'This phone sends SMS for Habarchy: messages from the server go out through your SIM card.',
    },
    'apiUrl': {'tk': 'Serwer (API) URL', 'ru': 'Адрес сервера (API)', 'en': 'Server (API) URL'},
    'gatewayKey': {'tk': 'Gateway açary', 'ru': 'Ключ шлюза', 'en': 'Gateway key'},
    'keyHint': {
      'tk': 'Admin panel → Providerler → android_sms → Jübütleşdirme',
      'ru': 'Админ-панель → Провайдеры → android_sms → Сопряжение',
      'en': 'Admin panel → Providers → android_sms → Pairing',
    },
    'forwardInbound': {'tk': 'Gelen SMS-leri serwere ugrat', 'ru': 'Пересылать входящие SMS на сервер', 'en': 'Forward received SMS to the server'},
    'save': {'tk': 'Ýatda sakla', 'ru': 'Сохранить', 'en': 'Save'},
    'start': {'tk': 'Işe gir', 'ru': 'Запустить', 'en': 'Start'},
    'stop': {'tk': 'Duruz', 'ru': 'Остановить', 'en': 'Stop'},
    'status': {'tk': 'Ýagdaý', 'ru': 'Состояние', 'en': 'Status'},
    'running': {'tk': 'Işleýär', 'ru': 'Работает', 'en': 'Running'},
    'stopped': {'tk': 'Durdy', 'ru': 'Остановлен', 'en': 'Stopped'},
    'connected': {'tk': 'Serwere birikdi', 'ru': 'Подключено к серверу', 'en': 'Connected to server'},
    'disconnected': {'tk': 'Serwer bilen baglanyşyk ýok', 'ru': 'Нет связи с сервером', 'en': 'No connection to server'},
    'provider': {'tk': 'Provider', 'ru': 'Провайдер', 'en': 'Provider'},
    'pending': {'tk': 'Nobatda', 'ru': 'В очереди', 'en': 'Pending'},
    'sent': {'tk': 'Iberildi', 'ru': 'Отправлено', 'en': 'Sent'},
    'delivered': {'tk': 'Gowşuryldy', 'ru': 'Доставлено', 'en': 'Delivered'},
    'failed': {'tk': 'Şowsuz', 'ru': 'Ошибки', 'en': 'Failed'},
    'inbound': {'tk': 'Gelen', 'ru': 'Входящие', 'en': 'Inbound'},
    'lastPoll': {'tk': 'Soňky sorag', 'ru': 'Последний опрос', 'en': 'Last poll'},
    'lastError': {'tk': 'Soňky ýalňyşlyk', 'ru': 'Последняя ошибка', 'en': 'Last error'},
    'permissions': {'tk': 'SMS rugsady gerek', 'ru': 'Нужно разрешение на SMS', 'en': 'SMS permission needed'},
    'grant': {'tk': 'Rugsat ber', 'ru': 'Разрешить', 'en': 'Grant'},
    'battery': {
      'tk': 'Programmanyň fonda durmagy üçin batareýa optimizasiýasyny öçüriň',
      'ru': 'Отключите оптимизацию батареи, чтобы приложение работало в фоне',
      'en': 'Disable battery optimisation so the app keeps running in the background',
    },
    'disable': {'tk': 'Öçür', 'ru': 'Отключить', 'en': 'Disable'},
    'log': {'tk': 'Žurnal', 'ru': 'Журнал', 'en': 'Log'},
    'clear': {'tk': 'Arassala', 'ru': 'Очистить', 'en': 'Clear'},
    'testSms': {'tk': 'Synag SMS', 'ru': 'Тестовое SMS', 'en': 'Test SMS'},
    'phone': {'tk': 'Telefon belgisi', 'ru': 'Номер телефона', 'en': 'Phone number'},
    'send': {'tk': 'Iber', 'ru': 'Отправить', 'en': 'Send'},
    'saved': {'tk': 'Ýatda saklandy', 'ru': 'Сохранено', 'en': 'Saved'},
    'fillAll': {'tk': 'URL we açary giriziň', 'ru': 'Введите URL и ключ', 'en': 'Enter the URL and key'},
    'never': {'tk': 'heniz ýok', 'ru': 'ещё нет', 'en': 'never'},
    'secondsAgo': {'tk': '{n} s öň', 'ru': '{n} с назад', 'en': '{n}s ago'},
    'language': {'tk': 'Dil', 'ru': 'Язык', 'en': 'Language'},
    'inboundOn': {'tk': 'Gelen SMS-ler serwere ugradylýar', 'ru': 'Входящие SMS пересылаются на сервер', 'en': 'Received SMS are forwarded to the server'},
    'inboundServerOff': {'tk': 'Serwerde (provider sazlamasynda) ýapyk — fonda hiç zat işlemeýär', 'ru': 'Отключено на сервере (в настройках провайдера) — в фоне ничего не работает', 'en': 'Disabled on the server (provider settings) — nothing runs in the background'},
    'inboundNoPermission': {'tk': 'SMS okamak rugsady berilmedi', 'ru': 'Нет разрешения на чтение SMS', 'en': 'SMS read permission not granted'},
    'inboundOff': {'tk': 'Ýapyk — gelen SMS-ler okalmaýar, fonda işlemeýär', 'ru': 'Выключено — входящие SMS не читаются, в фоне ничего не работает', 'en': 'Off — received SMS are not read, nothing runs in the background'},
  };

  String call(String key, [Map<String, Object>? args]) {
    var s = _t[key]?[code] ?? _t[key]?['en'] ?? key;
    args?.forEach((k, v) => s = s.replaceAll('{$k}', '$v'));
    return s;
  }
}
