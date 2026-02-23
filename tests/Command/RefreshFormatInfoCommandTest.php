<?php

namespace App\Tests\Command;

use App\Command\RefreshFormatInfoCommand;
use App\Service\ConversionCatalogService;
use App\Service\ConverterApiClient;
use App\Service\FormatInfoCatalog;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\ArrayAdapter;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Tester\CommandTester;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class RefreshFormatInfoCommandTest extends TestCase
{
    public function testExecuteStoresFetchedAndFallbackData(): void
    {
        $httpClient = new MockHttpClient(static function (string $method, string $url): MockResponse {
            if ('GET' === $method && 'http://converter-api:8081/v1/conversions' === $url) {
                return new MockResponse('{"formats":{"png":["jpg"],"avif":["png"]}}', ['http_code' => 200]);
            }

            if (str_contains($url, '/page/summary/Portable_Network_Graphics')) {
                return new MockResponse('{"title":"Portable Network Graphics","extract":"PNG is a raster graphics file format.","content_urls":{"desktop":{"page":"https://example.org/png"}}}', ['http_code' => 200]);
            }

            return new MockResponse('{}', ['http_code' => 404]);
        });

        $conversionCatalogService = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $command = new RefreshFormatInfoCommand(
            $conversionCatalogService,
            $httpClient,
            new FormatInfoCatalog(),
            '/tmp/format-info-default.json',
        );
        $tester = new CommandTester($command);
        $outputPath = sys_get_temp_dir().'/format-info-refresh-'.bin2hex(random_bytes(8)).'.json';

        try {
            self::assertSame(Command::SUCCESS, $tester->execute([
                '--output' => $outputPath,
            ]));

            self::assertFileExists($outputPath);
            $payload = json_decode((string) file_get_contents($outputPath), true, 512, JSON_THROW_ON_ERROR);

            self::assertIsArray($payload);
            self::assertArrayHasKey('formats', $payload);
            self::assertIsArray($payload['formats']);
            self::assertSame('Portable Network Graphics', $payload['formats']['png']['title'] ?? null);
            self::assertSame('https://example.org/png', $payload['formats']['png']['url'] ?? null);
            self::assertArrayHasKey('jpg', $payload['formats']);
            self::assertArrayHasKey('avif', $payload['formats']);
            self::assertStringContainsString('JPG', (string) ($payload['formats']['jpg']['summary'] ?? ''));
            self::assertStringContainsString('AVIF', (string) ($payload['formats']['avif']['summary'] ?? ''));
        } finally {
            @unlink($outputPath);
        }
    }

    public function testExecuteFailsWhenFormatsCannotBeLoaded(): void
    {
        $httpClient = new MockHttpClient(static function (string $method, string $url): MockResponse {
            if ('GET' === $method && 'http://converter-api:8081/v1/conversions' === $url) {
                return new MockResponse('{}', ['http_code' => 503]);
            }

            return new MockResponse('{}', ['http_code' => 404]);
        });

        $conversionCatalogService = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $command = new RefreshFormatInfoCommand(
            $conversionCatalogService,
            $httpClient,
            new FormatInfoCatalog(),
            '/tmp/format-info-default.json',
        );
        $tester = new CommandTester($command);
        $outputPath = sys_get_temp_dir().'/format-info-refresh-'.bin2hex(random_bytes(8)).'.json';

        try {
            self::assertSame(Command::FAILURE, $tester->execute([
                '--output' => $outputPath,
            ]));
            self::assertFileDoesNotExist($outputPath);
        } finally {
            @unlink($outputPath);
        }
    }
}
