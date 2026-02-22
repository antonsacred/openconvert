<?php

namespace App\Tests\Service;

use App\Dto\ConvertResult;
use App\Service\ConvertService;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\ArrayAdapter;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class ConvertServiceTest extends TestCase
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

        $service = new ConvertService($httpClient, new ArrayAdapter(), 'http://converter-api:8081/');

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

        $service = new ConvertService($httpClient, new ArrayAdapter(), 'http://converter-api:8081');

        $first = $service->getFormats();
        $second = $service->getFormats();

        self::assertSame(['png' => ['jpg']], $first);
        self::assertSame($first, $second);
        self::assertSame(1, $calls);
    }

    public function testGetFormatsThrowsWhenConverterApiNotConfigured(): void
    {
        $service = new ConvertService(new MockHttpClient(), new ArrayAdapter(), null);

        $this->expectException(\RuntimeException::class);
        $this->expectExceptionMessage('CONVERTER_API is not configured.');

        $service->getFormats();
    }

    public function testGetFormatsThrowsWhenPayloadIsInvalid(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{"formats":{"png":"jpg"}}', ['http_code' => 200]),
        ]);

        $service = new ConvertService($httpClient, new ArrayAdapter(), 'http://converter-api:8081');

        $this->expectException(\UnexpectedValueException::class);
        $this->expectExceptionMessage('Invalid converter response');

        $service->getFormats();
    }

    public function testConvertReturnsDto(): void
    {
        $httpClient = new MockHttpClient(static function (string $method, string $url, array $options): MockResponse {
            self::assertSame('POST', $method);
            self::assertSame('http://converter-api:8081/v1/convert', $url);
            self::assertArrayHasKey('body', $options);
            self::assertIsString($options['body']);

            return new MockResponse('{"fileName":"sample.jpg","mimeType":"image/jpeg","contentBase64":"Y29udmVydGVk"}', ['http_code' => 200]);
        });

        $service = new ConvertService($httpClient, new ArrayAdapter(), 'http://converter-api:8081');
        $result = $service->convert('png', 'jpg', 'sample.png', base64_encode('test-bytes'));

        self::assertInstanceOf(ConvertResult::class, $result);
        self::assertSame(200, $result->statusCode());
        self::assertSame('sample.jpg', $result->payload()['fileName'] ?? null);
    }

    public function testConvertReturnsErrorDtoWhenApiNotConfigured(): void
    {
        $service = new ConvertService(new MockHttpClient(), new ArrayAdapter(), null);
        $result = $service->convert('png', 'jpg', 'sample.png', base64_encode('test-bytes'));

        self::assertInstanceOf(ConvertResult::class, $result);
        self::assertSame(503, $result->statusCode());
        self::assertSame('converter_api_not_configured', $result->payload()['error']['code'] ?? null);
    }
}
