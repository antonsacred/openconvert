<?php

namespace App\Tests\Service;

use App\Service\ConversionFormatsService;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\ArrayAdapter;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class ConversionFormatsServiceTest extends TestCase
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

    public function testGetFormatsParsesAndNormalizesResponse(): void
    {
        $this->setConverterApi('http://converter-api:8081/');

        $responseBody = json_encode([
            'formats' => [
                ' PNG ' => [' JPG ', 'webp', 'webp'],
                'jpg' => ['png'],
            ],
        ], JSON_THROW_ON_ERROR);

        $httpClient = new MockHttpClient([
            new MockResponse($responseBody, ['http_code' => 200]),
        ]);

        $service = new ConversionFormatsService($httpClient, new ArrayAdapter());

        self::assertSame([
            'jpg' => ['png'],
            'png' => ['jpg', 'webp'],
        ], $service->getFormats());
    }

    public function testGetFormatsUsesCache(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $calls = 0;
        $httpClient = new MockHttpClient(static function () use (&$calls): MockResponse {
            ++$calls;

            return new MockResponse('{"formats":{"png":["jpg"]}}', ['http_code' => 200]);
        });

        $service = new ConversionFormatsService($httpClient, new ArrayAdapter());

        $first = $service->getFormats();
        $second = $service->getFormats();

        self::assertSame(['png' => ['jpg']], $first);
        self::assertSame($first, $second);
        self::assertSame(1, $calls);
    }

    public function testGetFormatsThrowsWhenConverterApiNotConfigured(): void
    {
        $this->setConverterApi(null);

        $service = new ConversionFormatsService(new MockHttpClient(), new ArrayAdapter());

        $this->expectException(\RuntimeException::class);
        $this->expectExceptionMessage('CONVERTER_API is not configured.');

        $service->getFormats();
    }

    public function testGetFormatsThrowsWhenPayloadIsInvalid(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $httpClient = new MockHttpClient([
            new MockResponse('{"formats":{"png":"jpg"}}', ['http_code' => 200]),
        ]);

        $service = new ConversionFormatsService($httpClient, new ArrayAdapter());

        $this->expectException(\UnexpectedValueException::class);
        $this->expectExceptionMessage('Invalid converter response');

        $service->getFormats();
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
