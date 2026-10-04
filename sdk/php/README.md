# Habarchy SDK — PHP / Laravel

Zero-dependency client (`ext-curl`, `ext-json`), PHP ≥ 8.1, with a Laravel 10–12 service
provider and facade. Made for projects like **tds.gov.tm** that live in a separate repository.

## Install into a Laravel 12 project (Composer path repository)

```sh
# in your Laravel project
composer config repositories.habarchy path ../habarchy/sdk/php     # relative or absolute path to this folder
composer require habarchy/sdk:@dev
php artisan vendor:publish --tag=habarchy-config                  # optional: config/habarchy.php
```

`.env`:

```
HABARCHY_URL=https://habarchy.example.tm
HABARCHY_API_KEY=hb_live_xxx
HABARCHY_SIGN=true
```

The provider and the `Habarchy` facade are auto-discovered.

```php
use Habarchy\Laravel\Facades\Habarchy;

$acc = Habarchy::sendMessage(['channel' => 'sms', 'to' => $user->phone, 'template' => 'otp', 'data' => ['code' => $code, 'minutes' => 5], 'idempotency_key' => "otp-{$user->id}-{$nonce}"]);
$otp = Habarchy::sendOTP($user->phone);                     // server-generated code
$ok  = Habarchy::verifyOTP($user->phone, $request->code)['verified'];
```

Or inject `Habarchy\Client` in a constructor. Without Laravel:

```php
$hb = new \Habarchy\Client('https://habarchy.example.tm', getenv('HABARCHY_API_KEY'), ['sign' => true]);
$hb->sendMessage(['channel' => 'email', 'to' => 'user@example.tm', 'subject' => 'Salam', 'body' => '<p>Hoş geldiňiz!</p>']);
```

Methods: `sendMessage`, `sendBatch`, `getMessage`, `cancelMessage`, `getBatch`, `sendOTP`,
`verifyOTP`, `registerDevice`. Errors throw `Habarchy\HabarchyException` (`status`,
`errorCode`, `details`). `sign => true` adds `X-Timestamp` / `X-Signature` (HMAC-SHA256,
docs/auth.md).

```sh
composer validate && composer test      # php tests/run.php
```
