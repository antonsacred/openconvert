<?php

namespace App\Tests\Service;

use App\Dto\HealthStatusResult;
use App\Service\ConverterApiClient;
use App\Service\ConverterHealthService;
use PHPUnit\Framework\TestCase;
use Symfony\Component\HttpClient\Exception\TransportException;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class ConverterHealthServiceTest extends TestCase
{
    public function testGetHealthStatusReturnsNotConfiguredWhenConverterApiIsMissing(): void
    {
        $service = new ConverterHealthService(new ConverterApiClient(new MockHttpClient(), null));
        $result = $service->getHealthStatus();

        self::assertInstanceOf(HealthStatusResult::class, $result);
        self::assertSame(503, $result->statusCode());
        self::assertSame('not_configured', $result->payload()['checks']['converter_api']['status'] ?? null);
    }

    public function testGetHealthStatusReturnsUpForHealthyConverterApi(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{"status":"ok"}', ['http_code' => 200]),
        ]);

        $service = new ConverterHealthService(new ConverterApiClient($httpClient, 'http://converter-api:8081/'));
        $result = $service->getHealthStatus();

        self::assertSame(200, $result->statusCode());
        self::assertSame('up', $result->payload()['checks']['converter_api']['status'] ?? null);
        self::assertSame('http://converter-api:8081/health', $result->payload()['checks']['converter_api']['url'] ?? null);
    }

    public function testGetHealthStatusReturnsNotRunningWhenUnreachable(): void
    {
        $httpClient = new MockHttpClient(static function (): never {
            throw new TransportException('Connection refused');
        });

        $service = new ConverterHealthService(new ConverterApiClient($httpClient, 'http://converter-api:8081'));
        $result = $service->getHealthStatus();

        self::assertSame(503, $result->statusCode());
        self::assertSame('not_running', $result->payload()['checks']['converter_api']['status'] ?? null);
        self::assertSame('http://converter-api:8081/health', $result->payload()['checks']['converter_api']['url'] ?? null);
    }
}
