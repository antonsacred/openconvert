<?php

namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\Exception\TransportException;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class HealthControllerTest extends WebTestCase
{
    private ?string $originalServerConverterApi = null;
    private ?string $originalEnvConverterApi = null;
    private ?string $originalGetenvConverterApi = null;

    protected function setUp(): void
    {
        parent::setUp();

        $this->originalServerConverterApi = $_SERVER['CONVERTER_API'] ?? null;
        $this->originalEnvConverterApi = $_ENV['CONVERTER_API'] ?? null;
        $getenvValue = getenv('CONVERTER_API');
        $this->originalGetenvConverterApi = false === $getenvValue ? null : $getenvValue;
    }

    protected function tearDown(): void
    {
        $this->setConverterApi($this->originalServerConverterApi);

        if ($this->originalEnvConverterApi === null) {
            unset($_ENV['CONVERTER_API']);
        } else {
            $_ENV['CONVERTER_API'] = $this->originalEnvConverterApi;
        }

        if ($this->originalGetenvConverterApi === null) {
            putenv('CONVERTER_API');
        } else {
            putenv('CONVERTER_API='.$this->originalGetenvConverterApi);
        }

        parent::tearDown();
    }

    public function testHealthReportsConverterApiNotConfigured(): void
    {
        $this->setConverterApi(null);

        $client = static::createClient();
        $client->request('GET', '/health');

        self::assertResponseStatusCodeSame(503);

        $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);

        self::assertSame('degraded', $payload['status'] ?? null);
        self::assertSame('not_configured', $payload['checks']['converter_api']['status'] ?? null);
    }

    public function testHealthReportsConverterApiNotRunningWhenUnreachable(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $mockClient = new MockHttpClient(static function (): never {
            throw new TransportException('Connection refused');
        });

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $mockClient);
        $client->request('GET', '/health');

        self::assertResponseStatusCodeSame(503);

        $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);

        self::assertSame('degraded', $payload['status'] ?? null);
        self::assertSame('not_running', $payload['checks']['converter_api']['status'] ?? null);
        self::assertSame('http://converter-api:8081/health', $payload['checks']['converter_api']['url'] ?? null);
    }

    public function testHealthReportsConverterApiUpWhenHealthy(): void
    {
        $this->setConverterApi('http://converter-api:8081/');

        $mockClient = new MockHttpClient([
            new MockResponse('{"status":"ok"}', ['http_code' => 200]),
        ]);

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $mockClient);
        $client->request('GET', '/health');

        self::assertResponseIsSuccessful();

        $payload = json_decode($client->getResponse()->getContent(), true, 512, JSON_THROW_ON_ERROR);

        self::assertSame('ok', $payload['status'] ?? null);
        self::assertSame('up', $payload['checks']['converter_api']['status'] ?? null);
        self::assertSame('http://converter-api:8081/health', $payload['checks']['converter_api']['url'] ?? null);
    }

    private function setConverterApi(?string $value): void
    {
        if ($value === null) {
            unset($_SERVER['CONVERTER_API']);
            unset($_ENV['CONVERTER_API']);
            putenv('CONVERTER_API');

            return;
        }

        $_SERVER['CONVERTER_API'] = $value;
        $_ENV['CONVERTER_API'] = $value;
        putenv('CONVERTER_API='.$value);
    }
}
