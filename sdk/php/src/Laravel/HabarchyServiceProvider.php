<?php

declare(strict_types=1);

namespace Habarchy\Laravel;

use Habarchy\Client;
use Illuminate\Support\ServiceProvider;

/** Laravel 10–12: binds Habarchy\Client as a singleton from config/habarchy.php. */
final class HabarchyServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->mergeConfigFrom(__DIR__ . '/../../config/habarchy.php', 'habarchy');
        $this->app->singleton(Client::class, function ($app): Client {
            $cfg = $app['config']->get('habarchy');

            return new Client((string) $cfg['url'], (string) $cfg['api_key'], [
                'sign' => (bool) ($cfg['sign'] ?? false),
                'timeout' => (int) ($cfg['timeout'] ?? 20),
            ]);
        });
        $this->app->alias(Client::class, 'habarchy');
    }

    public function boot(): void
    {
        if ($this->app->runningInConsole()) {
            $this->publishes([__DIR__ . '/../../config/habarchy.php' => $this->app->configPath('habarchy.php')], 'habarchy-config');
        }
    }
}
