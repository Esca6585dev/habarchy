<?php

declare(strict_types=1);

// Dependency-free test runner: php tests/run.php
spl_autoload_register(static function (string $class): void {
    if (str_starts_with($class, 'Habarchy\\')) {
        $file = __DIR__ . '/../src/' . str_replace('\\', '/', substr($class, 9)) . '.php';
        if (is_file($file)) {
            require $file;
        }
    }
});

use Habarchy\Client;
use Habarchy\HabarchyException;

$failures = 0;
function check(bool $ok, string $what): void
{
    global $failures;
    echo ($ok ? 'ok   ' : 'FAIL ') . $what . PHP_EOL;
    if (!$ok) {
        $failures++;
    }
}

const KEY = 'hb_live_testkey';
$serverSig = static fn (string $ts, string $m, string $p, string $b): string => hash_hmac('sha256', "$ts\n$m\n$p\n" . hash('sha256', $b), KEY);

check(Client::signature(KEY, '1700000000', 'post', '/api/v1/messages', '{"a":1}') === $serverSig('1700000000', 'POST', '/api/v1/messages', '{"a":1}'), 'signature matches server algorithm');

$seen = [];
$client = new Client('https://h.example.tm/', KEY, [
    'sign' => true,
    'now' => static fn (): int => 1700000000,
    'transport' => static function (string $method, string $url, array $headers, string $body) use (&$seen): array {
        $seen = compact('method', 'url', 'headers', 'body');
        if (str_ends_with($url, '/otp/verify')) {
            return [429, '{"error":{"code":"rate_limited","message":"slow","details":{"retry_after":30}}}'];
        }
        if (str_ends_with($url, '/messages/batch')) {
            return [202, '{"data":{"id":"b1","status":"queued","total":2},"meta":{"accepted":2,"rejected":0}}'];
        }
        if (str_contains($url, '/messages/m1')) {
            return [200, '{"data":{"message":{"id":"m1","status":"delivered"},"events":[{"type":"sent"}]}}'];
        }
        if (str_ends_with($url, '/messages')) {
            return [202, '{"data":{"id":"m1","status":"queued","channel":"sms","to":"+99365123456"},"meta":null,"error":null}'];
        }
        return [504, 'gateway timeout'];
    },
]);

$acc = $client->sendMessage(['channel' => 'sms', 'to' => '+99365123456', 'template' => 'otp', 'data' => ['code' => '4821']]);
check($acc['id'] === 'm1' && $acc['duplicate'] === false, 'sendMessage unwraps envelope');
check($seen['url'] === 'https://h.example.tm/api/v1/messages', 'base url trailing slash handled');
check($seen['headers']['X-Api-Key'] === KEY && $seen['headers']['X-Timestamp'] === '1700000000', 'auth headers');
check($seen['headers']['X-Signature'] === $serverSig('1700000000', 'POST', '/api/v1/messages', $seen['body']), 'request is signed');
check(json_decode($seen['body'], true)['to'] === '+99365123456', 'body encodes recipient');

$b = $client->sendBatch(['channel' => 'sms', 'body' => 'x', 'recipients' => [['to' => '+1'], ['to' => ['external_id' => 'u2']]]]);
check($b['id'] === 'b1' && $b['accepted'] === 2, 'sendBatch merges meta');

$d = $client->getMessage('m1');
check($d['message']['status'] === 'delivered' && count($d['events']) === 1, 'getMessage returns message + events');

try {
    $client->verifyOTP('+993', '0000');
    check(false, 'verifyOTP throws');
} catch (HabarchyException $e) {
    check($e->errorCode === 'rate_limited' && $e->status === 429 && ($e->details['retry_after'] ?? null) === 30, 'API error mapped to HabarchyException');
}
try {
    $client->registerDevice('tok', 'android');
    check(false, 'non-JSON error throws');
} catch (HabarchyException $e) {
    check($e->errorCode === 'http_error' && $e->status === 504, 'non-JSON error mapped');
}

echo PHP_EOL . ($failures === 0 ? 'ALL OK' : "$failures FAILED") . PHP_EOL;
exit($failures === 0 ? 0 : 1);
