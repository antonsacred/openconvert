<?php

namespace App\Tests\Service;

use App\Service\ConversionCatalogService;
use App\Service\ConverterApiClient;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\ArrayAdapter;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class ConversionCatalogServiceTest extends TestCase
{
    public function testGetFormatsParsesAndNormalizesResponse(): void
    {
        $responseBody = json_encode([
            'formats' => [
                ' PNG ' => [' JPG ', 'webp', 'webp'],
                'jpg' => ['png'],
            ],
        ], JSON_THROW_ON_ERROR);

        $httpClient = new MockHttpClient([
            new MockResponse($responseBody, ['http_code' => 200]),
        ]);

        $service = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081/'),
            new ArrayAdapter(),
        );

        self::assertSame([
            'jpg' => ['png'],
            'png' => ['jpg', 'webp'],
        ], $service->getFormats());
    }

    public function testGetFormatsUsesCache(): void
    {
        $calls = 0;
        $httpClient = new MockHttpClient(static function () use (&$calls): MockResponse {
            ++$calls;

            return new MockResponse('{"formats":{"png":["jpg"]}}', ['http_code' => 200]);
        });

        $service = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $first = $service->getFormats();
        $second = $service->getFormats();

        self::assertSame(['png' => ['jpg']], $first);
        self::assertSame($first, $second);
        self::assertSame(1, $calls);
    }

    public function testGetFormatsThrowsWhenConverterApiNotConfigured(): void
    {
        $service = new ConversionCatalogService(
            new ConverterApiClient(new MockHttpClient(), null),
            new ArrayAdapter(),
        );

        $this->expectException(\RuntimeException::class);
        $this->expectExceptionMessage('CONVERTER_API is not configured.');

        $service->getFormats();
    }

    public function testGetFormatsThrowsWhenPayloadIsInvalid(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{"formats":{"png":"jpg"}}', ['http_code' => 200]),
        ]);

        $service = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $this->expectException(\UnexpectedValueException::class);
        $this->expectExceptionMessage('Invalid converter response');

        $service->getFormats();
    }
}
