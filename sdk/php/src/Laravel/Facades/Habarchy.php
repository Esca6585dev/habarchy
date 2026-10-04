<?php

declare(strict_types=1);

namespace Habarchy\Laravel\Facades;

use Illuminate\Support\Facades\Facade;

/**
 * @method static array sendMessage(array $request)
 * @method static array sendBatch(array $request)
 * @method static array getMessage(string $id)
 * @method static array cancelMessage(string $id)
 * @method static array getBatch(string $id)
 * @method static array sendOTP(string $to, array $options = [])
 * @method static array verifyOTP(string $to, string $code)
 * @method static array registerDevice(string $token, string $platform, array $options = [])
 *
 * @see \Habarchy\Client
 */
final class Habarchy extends Facade
{
    protected static function getFacadeAccessor(): string
    {
        return \Habarchy\Client::class;
    }
}
