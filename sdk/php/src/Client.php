<?php

declare(strict_types=1);

namespace Habarchy;

/**
 * Minimal client for the Habarchy public API.
 *
 *   $hb = new Client('https://habarchy.example.tm', $apiKey, ['sign' => true]);
 *   $acc = $hb->sendMessage(['channel' => 'sms', 'to' => '+99365123456', 'template' => 'otp', 'data' => ['code' => '4821']]);
 */
class Client
{
    private string $baseUrl;
    private bool $sign;
    private int $timeout;
    /** @var callable(string,string,array<string,string>,string):array{0:int,1:string} */
    private $transport;
    /** @var callable():int */
    private $now;

    /**
     * @param array{sign?:bool,timeout?:int,transport?:callable,now?:callable} $options
     */
    public function __construct(string $baseUrl, private readonly string $apiKey, array $options = [])
    {
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->sign = (bool) ($options['sign'] ?? false);
        $this->timeout = (int) ($options['timeout'] ?? 20);
        $this->transport = $options['transport'] ?? [$this, 'curlTransport'];
        $this->now = $options['now'] ?? static fn (): int => time();
    }

    /** hex(HMAC-SHA256(apiKey, ts "\n" METHOD "\n" path "\n" hex(sha256(body)))) */
    public static function signature(string $apiKey, string $timestamp, string $method, string $path, string $body): string
    {
        $msg = $timestamp . "\n" . strtoupper($method) . "\n" . $path . "\n" . hash('sha256', $body);

        return hash_hmac('sha256', $msg, $apiKey);
    }

    /** @return array<string,mixed> `id`, `status`, `channel`, `to`, `duplicate` */
    public function sendMessage(array $request): array
    {
        [$data, $meta] = $this->request('POST', '/api/v1/messages', $request);
        $data['duplicate'] = (bool) ($meta['duplicate'] ?? false);

        return $data;
    }

    /** @return array<string,mixed> batch with `accepted` / `rejected` merged from meta */
    public function sendBatch(array $request): array
    {
        [$data, $meta] = $this->request('POST', '/api/v1/messages/batch', $request);
        $data['accepted'] = $meta['accepted'] ?? null;
        $data['rejected'] = $meta['rejected'] ?? null;

        return $data;
    }

    /** @return array{message:array<string,mixed>,events:array<int,array<string,mixed>>} */
    public function getMessage(string $id): array
    {
        return $this->request('GET', '/api/v1/messages/' . rawurlencode($id))[0];
    }

    public function cancelMessage(string $id): array
    {
        return $this->request('POST', '/api/v1/messages/' . rawurlencode($id) . '/cancel')[0];
    }

    public function getBatch(string $id): array
    {
        return $this->request('GET', '/api/v1/batches/' . rawurlencode($id))[0];
    }

    /** @return array{message_id:string,to:string,channel:string,expires_in:int,length:int} */
    public function sendOTP(string $to, array $options = []): array
    {
        return $this->request('POST', '/api/v1/otp/send', ['to' => $to] + $options)[0];
    }

    /** @return array{verified:bool,attempts_remaining:int} */
    public function verifyOTP(string $to, string $code): array
    {
        return $this->request('POST', '/api/v1/otp/verify', ['to' => $to, 'code' => $code])[0];
    }

    public function registerDevice(string $token, string $platform, array $options = []): array
    {
        return $this->request('POST', '/api/v1/devices', ['token' => $token, 'platform' => $platform] + $options)[0];
    }

    /** @return array{0:array<string,mixed>,1:array<string,mixed>} [data, meta] */
    private function request(string $method, string $path, ?array $body = null): array
    {
        $payload = $body === null ? '' : json_encode($body, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR);
        $headers = ['X-Api-Key' => $this->apiKey, 'Accept' => 'application/json', 'User-Agent' => 'habarchy-php/1.0'];
        if ($payload !== '') {
            $headers['Content-Type'] = 'application/json';
        }
        if ($this->sign) {
            $ts = (string) ($this->now)();
            $headers['X-Timestamp'] = $ts;
            $headers['X-Signature'] = self::signature($this->apiKey, $ts, $method, $path, $payload);
        }
        [$status, $text] = ($this->transport)($method, $this->baseUrl . $path, $headers, $payload);
        $env = [];
        if ($text !== '') {
            $env = json_decode($text, true);
            if (!is_array($env)) {
                if ($status >= 300) {
                    throw new HabarchyException($status, 'http_error', substr($text, 0, 200));
                }
                throw new HabarchyException($status, 'bad_response', 'response is not JSON');
            }
        }
        if ($status >= 300) {
            $e = $env['error'] ?? ['code' => 'http_error', 'message' => 'HTTP ' . $status];
            throw new HabarchyException($status, (string) ($e['code'] ?? 'http_error'), (string) ($e['message'] ?? ''), $e['details'] ?? null);
        }

        return [$env['data'] ?? [], is_array($env['meta'] ?? null) ? $env['meta'] : []];
    }

    /** @param array<string,string> $headers @return array{0:int,1:string} */
    private function curlTransport(string $method, string $url, array $headers, string $body): array
    {
        $ch = curl_init($url);
        $hdr = [];
        foreach ($headers as $k => $v) {
            $hdr[] = $k . ': ' . $v;
        }
        curl_setopt_array($ch, [
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_HTTPHEADER => $hdr,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT => $this->timeout,
            CURLOPT_CONNECTTIMEOUT => 10,
        ]);
        if ($body !== '') {
            curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
        }
        $text = curl_exec($ch);
        if ($text === false) {
            $err = curl_error($ch);
            curl_close($ch);
            throw new HabarchyException(0, 'network_error', $err);
        }
        $status = (int) curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
        curl_close($ch);

        return [$status, (string) $text];
    }
}
