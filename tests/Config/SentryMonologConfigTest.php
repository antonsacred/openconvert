<?php

namespace App\Tests\Config;

use Monolog\Logger;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Yaml\Yaml;

final class SentryMonologConfigTest extends TestCase
{
    public function testProdConfigEnablesMonologSentryHandler(): void
    {
        $config = Yaml::parseFile(
            __DIR__.'/../../config/packages/sentry.yaml',
            Yaml::PARSE_CONSTANT,
        );

        self::assertIsArray($config);

        $prodConfig = $config['when@prod'] ?? null;
        self::assertIsArray($prodConfig);

        self::assertSame(false, $prodConfig['sentry']['register_error_listener'] ?? null);
        self::assertSame(false, $prodConfig['sentry']['register_error_handler'] ?? null);
        self::assertSame(0.1, $prodConfig['sentry']['options']['traces_sample_rate'] ?? null);

        self::assertSame(
            'service',
            $prodConfig['monolog']['handlers']['sentry']['type'] ?? null,
        );
        self::assertSame(
            'Sentry\\Monolog\\Handler',
            $prodConfig['monolog']['handlers']['sentry']['id'] ?? null,
        );

        self::assertSame(
            '@Sentry\\State\\HubInterface',
            $prodConfig['services']['Sentry\\Monolog\\Handler']['arguments']['$hub'] ?? null,
        );
        self::assertSame(
            Logger::WARNING,
            $prodConfig['services']['Sentry\\Monolog\\Handler']['arguments']['$level'] ?? null,
        );
        self::assertSame(
            true,
            $prodConfig['services']['Sentry\\Monolog\\Handler']['arguments']['$fillExtraContext'] ?? null,
        );
    }
}
