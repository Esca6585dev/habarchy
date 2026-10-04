// ignore: unused_import
import 'package:intl/intl.dart' as intl;

import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Russian (`ru`).
class AppLocalizationsRu extends AppLocalizations {
  AppLocalizationsRu([String locale = 'ru']) : super(locale);

  @override
  String get appTitle => 'Habarchy';

  @override
  String get signIn => 'Войти';

  @override
  String get email => 'E-mail';

  @override
  String get password => 'Пароль';

  @override
  String get totpCode => 'Код аутентификатора';

  @override
  String get totpRequired =>
      'Включена двухфакторная аутентификация. Введите 6-значный код.';

  @override
  String get invalidCredentials => 'Неверный e-mail или пароль';

  @override
  String get serverUrl => 'Сервер';

  @override
  String get dashboard => 'Панель';

  @override
  String get messages => 'Сообщения';

  @override
  String get templates => 'Шаблоны';

  @override
  String get providers => 'Провайдеры';

  @override
  String get apiKeys => 'API-ключи';

  @override
  String get settings => 'Настройки';

  @override
  String get logout => 'Выйти';

  @override
  String get projects => 'Проекты';

  @override
  String get noProjects => 'Проектов пока нет. Создайте его в веб-панели.';

  @override
  String get today => 'Сегодня';

  @override
  String get total => 'Всего';

  @override
  String get sent => 'Отправлено';

  @override
  String get delivered => 'Доставлено';

  @override
  String get failed => 'Ошибки';

  @override
  String get pending => 'В ожидании';

  @override
  String get cost => 'Расходы';

  @override
  String get perDay => 'Сообщений в день';

  @override
  String get recentFailures => 'Последние ошибки';

  @override
  String get nothingHere => 'Пока пусто';

  @override
  String get offlineCached => 'Офлайн — показаны последние загруженные данные';

  @override
  String get search => 'Поиск';

  @override
  String get all => 'Все';

  @override
  String get status => 'Статус';

  @override
  String get channel => 'Канал';

  @override
  String get resend => 'Отправить повторно';

  @override
  String get cancel => 'Отмена';

  @override
  String get timeline => 'Хронология';

  @override
  String get rawResponse => 'Ответ провайдера';

  @override
  String get attempts => 'Попытки';

  @override
  String get provider => 'Провайдер';

  @override
  String get preview => 'Предпросмотр';

  @override
  String get sampleData => 'Пример данных (JSON)';

  @override
  String get testSend => 'Тестовая отправка';

  @override
  String get recipient => 'Получатель';

  @override
  String get send => 'Отправить';

  @override
  String get revoke => 'Отозвать';

  @override
  String get revoked => 'Отозван';

  @override
  String get live => 'live';

  @override
  String get test => 'test';

  @override
  String get language => 'Язык';

  @override
  String get theme => 'Тема';

  @override
  String get themeSystem => 'Системная';

  @override
  String get themeLight => 'Светлая';

  @override
  String get themeDark => 'Тёмная';

  @override
  String get pushDemo => 'Демо-режим push';

  @override
  String get pushDemoHint =>
      'Регистрирует FCM-токен этого устройства в проекте (POST /api/v1/devices), чтобы получать тестовые push здесь.';

  @override
  String get pushNotConfigured =>
      'Firebase не настроен для этой сборки. Выполните `flutterfire configure` и пересоберите.';

  @override
  String get apiKeyForPush => 'API-ключ проекта (права devices)';

  @override
  String get externalId => 'Внешний ID (контакт)';

  @override
  String get registerDevice => 'Зарегистрировать устройство';

  @override
  String get deviceRegistered => 'Устройство зарегистрировано';

  @override
  String get fcmToken => 'FCM-токен';

  @override
  String get incomingPushes => 'Входящие push';

  @override
  String get noPushesYet => 'Push ещё не получены';

  @override
  String get copy => 'Копировать';

  @override
  String get copied => 'Скопировано';

  @override
  String get error => 'Что-то пошло не так';

  @override
  String get retry => 'Повторить';

  @override
  String get loading => 'Загрузка…';

  @override
  String version(String version) {
    return 'Версия $version';
  }

  @override
  String get statusQueued => 'В очереди';

  @override
  String get statusProcessing => 'Обработка';

  @override
  String get statusSent => 'Отправлено';

  @override
  String get statusDelivered => 'Доставлено';

  @override
  String get statusFailed => 'Ошибка';

  @override
  String get statusCancelled => 'Отменено';

  @override
  String get healthy => 'В норме';

  @override
  String get degraded => 'Деградация';

  @override
  String get failing => 'Сбои';

  @override
  String get idle => 'Простой';

  @override
  String get disabled => 'Отключён';

  @override
  String get queues => 'Очереди';

  @override
  String get systemHealth => 'Состояние системы';

  @override
  String get switchProject => 'Сменить проект';

  @override
  String get role => 'Роль';

  @override
  String get confirmResend => 'Поставить сообщение в очередь снова?';

  @override
  String get yes => 'Да';

  @override
  String get no => 'Нет';

  @override
  String get requiredVars => 'Переменные';

  @override
  String get groups => 'Группы';

  @override
  String get contacts => 'Контакты';

  @override
  String get compose => 'Отправить';

  @override
  String get more => 'Ещё';

  @override
  String get newGroup => 'Новая группа';

  @override
  String get groupName => 'Название группы';

  @override
  String get description => 'Описание';

  @override
  String get members => 'Участники';

  @override
  String get addMembers => 'Добавить участников';

  @override
  String get pickContacts => 'Существующие контакты';

  @override
  String get newContacts => 'Новые контакты';

  @override
  String get newContactsHint => 'По одному в строке: Имя, +99365…, mail@…';

  @override
  String get remove => 'Убрать';

  @override
  String get sendToGroup => 'Отправить группе';

  @override
  String get name => 'Имя';

  @override
  String get phone => 'Телефон';

  @override
  String get whatsapp => 'WhatsApp';

  @override
  String get telegram => 'Telegram chat id';

  @override
  String get slack => 'Slack id';

  @override
  String get newContact => 'Новый контакт';

  @override
  String get editContact => 'Изменить контакт';

  @override
  String get delete => 'Удалить';

  @override
  String confirmDelete(String name) {
    return 'Удалить «$name»?';
  }

  @override
  String get save => 'Сохранить';

  @override
  String get template => 'Шаблон';

  @override
  String get freeText => 'Свободный текст';

  @override
  String get subject => 'Тема / заголовок';

  @override
  String get body => 'Сообщение';

  @override
  String get templateData => 'Данные шаблона (JSON)';

  @override
  String get sandboxMode => 'Тестовый режим (sandbox, реально не отправляется)';

  @override
  String get recipients => 'Получатели';

  @override
  String get addresses => 'Адреса, по одному в строке';

  @override
  String get addressesHint =>
      'Телефоны, e-mail или chat id в зависимости от канала';

  @override
  String queuedRejected(int n, int m) {
    return '$n в очереди · $m отклонено';
  }

  @override
  String get selectRecipients => 'Выберите группу или введите адрес';

  @override
  String get noGroupsYet =>
      'Групп пока нет. Создайте «Коллеги», «Клиенты»… и добавьте участников.';

  @override
  String added(int n, int created) {
    return '$n добавлено · $created новых контактов';
  }

  @override
  String get newProvider => 'Новый провайдер';

  @override
  String get editProvider => 'Изменить провайдер';

  @override
  String get providerType => 'Тип';

  @override
  String get priority => 'Приоритет (меньше — раньше)';

  @override
  String get rateLimit => 'Лимит (сообщ./с, 0 = нет)';

  @override
  String get active => 'Активен';

  @override
  String get credentials => 'Данные доступа';

  @override
  String get keepCredentials =>
      'Оставьте пустым, чтобы сохранить текущие данные';

  @override
  String get pairing => 'Подключить телефон';

  @override
  String get pairingHint =>
      'Установите APK Habarçy Gateway на телефон с SIM-картой, введите этот URL и ключ, нажмите «Запустить».';

  @override
  String get gatewayKey => 'Ключ шлюза';

  @override
  String get apiUrl => 'API URL';

  @override
  String get online => 'Онлайн';

  @override
  String get offline => 'Офлайн';

  @override
  String get whatsappHint =>
      'Свободный текст работает в 24-часовом окне; иначе нужен одобренный шаблон.';

  @override
  String get contactHint =>
      'Нужен хотя бы телефон, e-mail, WhatsApp, Telegram или Slack.';

  @override
  String get openMessages => 'Открыть журнал';

  @override
  String get noMembers => 'Участников пока нет';

  @override
  String get inboundSms => 'Входящие SMS';

  @override
  String get inboundHint =>
      'Здесь появляются SMS, полученные телефонами-шлюзами при включённой пересылке (настройка провайдера + переключатель в приложении).';
}
