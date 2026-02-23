<?php

namespace App\Tests\Controller;

use App\Tests\Support\ConverterApiClientOverrideTrait;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\Exception\TransportException;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class HealthControllerTest extends WebTestCase
{
    use ConverterApiClientOverrideTrait;

    public function testHealthReportsConverterApiNotConfigured(): void
    {
        $client = static::createClient();
        $this->overrideConverterApiClient(null);
        $client->request('GET', '/health');

        self::assertResponseStatusCodeSame(503);

        $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);

        self::assertSame('degraded', $payload['status'] ?? null);
        self::assertSame('not_configured', $payload['checks']['converter_api']['status'] ?? null);
        self::assertNotSame('', $client->getResponse()->headers->get('X-Request-Id') ?? '');
    }

    public function testHealthReportsConverterApiNotRunningWhenUnreachable(): void
    {
        $mockClient = new MockHttpClient(static function (): never {
            throw new TransportException('Connection refused');
        });

        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $mockClient);
        $client->request('GET', '/health');

        self::assertResponseStatusCodeSame(503);

        $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);

        self::assertSame('degraded', $payload['status'] ?? null);
        self::assertSame('not_running', $payload['checks']['converter_api']['status'] ?? null);
        self::assertSame('http://converter-api:8081/health', $payload['checks']['converter_api']['url'] ?? null);
        self::assertNotSame('', $client->getResponse()->headers->get('X-Request-Id') ?? '');
    }

    public function testHealthReportsConverterApiUpWhenHealthy(): void
    {
        $mockClient = new MockHttpClient([
            new MockResponse('{"status":"ok"}', ['http_code' => 200]),
        ]);

        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081/', $mockClient);
        $client->request('GET', '/health');

        self::assertResponseIsSuccessful();

        $payload = json_decode($client->getResponse()->getContent(), true, 512, JSON_THROW_ON_ERROR);

        self::assertSame('ok', $payload['status'] ?? null);
        self::assertSame('up', $payload['checks']['converter_api']['status'] ?? null);
        self::assertSame('http://converter-api:8081/health', $payload['checks']['converter_api']['url'] ?? null);
        self::assertNotSame('', $client->getResponse()->headers->get('X-Request-Id') ?? '');
    }

    public function testHealthPropagatesIncomingRequestIdHeader(): void
    {
        $mockClient = new MockHttpClient([
            new MockResponse('{"status":"ok"}', ['http_code' => 200]),
        ]);

        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $mockClient);
        $client->request('GET', '/health', server: [
            'HTTP_X_REQUEST_ID' => 'req-health-001',
        ]);

        self::assertResponseIsSuccessful();
        self::assertSame('req-health-001', $client->getResponse()->headers->get('X-Request-Id'));
    }
}
