<?php

namespace App\Tests\Service;

use App\Dto\ConvertResult;
use App\Service\ConversionExecutionService;
use App\Service\ConverterApiClient;
use PHPUnit\Framework\TestCase;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class ConversionExecutionServiceTest extends TestCase
{
    public function testConvertReturnsDto(): void
    {
        $httpClient = new MockHttpClient(static function (string $method, string $url, array $options): MockResponse {
            self::assertSame('POST', $method);
            self::assertSame('http://converter-api:8081/v1/convert', $url);
            self::assertArrayHasKey('body', $options);
            self::assertIsString($options['body']);
            self::assertArrayHasKey('normalized_headers', $options);
            self::assertSame(['X-Request-Id: req-123'], $options['normalized_headers']['x-request-id'] ?? null);

            return new MockResponse('{"fileName":"sample.jpg","mimeType":"image/jpeg","contentBase64":"Y29udmVydGVk"}', ['http_code' => 200]);
        });

        $service = new ConversionExecutionService(new ConverterApiClient($httpClient, 'http://converter-api:8081'));
        $result = $service->convert('png', 'jpg', 'sample.png', base64_encode('test-bytes'), 'req-123');

        self::assertInstanceOf(ConvertResult::class, $result);
        self::assertSame(200, $result->statusCode());
        self::assertSame('sample.jpg', $result->payload()['fileName'] ?? null);
    }

    public function testConvertReturnsErrorDtoWhenApiNotConfigured(): void
    {
        $service = new ConversionExecutionService(new ConverterApiClient(new MockHttpClient(), null));
        $result = $service->convert('png', 'jpg', 'sample.png', base64_encode('test-bytes'), 'req-123');

        self::assertInstanceOf(ConvertResult::class, $result);
        self::assertSame(503, $result->statusCode());
        self::assertSame('converter_api_not_configured', $result->payload()['error']['code'] ?? null);
        self::assertSame('req-123', $result->payload()['error']['requestId'] ?? null);
    }
}
