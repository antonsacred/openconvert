<?php

namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\Exception\TransportException;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;
use Symfony\Component\HttpFoundation\File\UploadedFile;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class ConvertControllerTest extends WebTestCase
{
    /**
     * @var array<string, array{server: ?string, env: ?string, getenv: ?string}>
     */
    private array $originalEnvVars = [];

    protected function setUp(): void
    {
        parent::setUp();

        $this->originalEnvVars['CONVERTER_API'] = $this->snapshotEnvVar('CONVERTER_API');
        $this->originalEnvVars['APP_MAX_UPLOAD_BYTES'] = $this->snapshotEnvVar('APP_MAX_UPLOAD_BYTES');
        $this->originalEnvVars['APP_CONVERT_RATE_LIMIT'] = $this->snapshotEnvVar('APP_CONVERT_RATE_LIMIT');
    }

    protected function tearDown(): void
    {
        foreach ($this->originalEnvVars as $name => $snapshot) {
            $this->restoreEnvVar($name, $snapshot);
        }

        parent::tearDown();
    }

    public function testConvertReturnsConvertedPayloadAndRequestIdHeader(): void
    {
        $this->setEnvVar('CONVERTER_API', 'http://converter-api:8081');

        [$uploadedFile, $inputFilePath] = $this->createUploadedFile('sample.png');
        try {
            $mockClient = new MockHttpClient(static function (string $method, string $url, array $options): MockResponse {
                self::assertSame('POST', $method);
                self::assertSame('http://converter-api:8081/v1/convert', $url);
                self::assertArrayHasKey('body', $options);
                self::assertIsString($options['body']);
                $payload = json_decode($options['body'], true, 512, JSON_THROW_ON_ERROR);
                self::assertIsArray($payload);
                self::assertSame('png', $payload['from'] ?? null);
                self::assertSame('jpg', $payload['to'] ?? null);
                self::assertSame('sample.png', $payload['fileName'] ?? null);
                self::assertSame(base64_encode('test-image-bytes'), $payload['contentBase64'] ?? null);

                $normalizedHeaders = $options['normalized_headers'] ?? [];
                self::assertIsArray($normalizedHeaders);
                $requestIdHeaderValues = $normalizedHeaders['x-request-id'] ?? [];
                self::assertIsArray($requestIdHeaderValues);
                self::assertStringContainsString('req-abc-123', implode('|', $requestIdHeaderValues));

                return new MockResponse('{"fileName":"sample.jpg","mimeType":"image/jpeg","contentBase64":"Y29udmVydGVk"}', ['http_code' => 200]);
            });

            $client = static::createClient();
            static::getContainer()->set(HttpClientInterface::class, $mockClient);
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $uploadedFile,
            ], [
                'HTTP_X_REQUEST_ID' => 'req-abc-123',
            ]);

            self::assertResponseIsSuccessful();
            self::assertSame('req-abc-123', $client->getResponse()->headers->get('X-Request-Id'));

            $payload = json_decode($client->getResponse()->getContent(), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('sample.jpg', $payload['fileName'] ?? null);
            self::assertSame('image/jpeg', $payload['mimeType'] ?? null);
            self::assertSame('Y29udmVydGVk', $payload['contentBase64'] ?? null);
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertReportsConverterApiNotConfigured(): void
    {
        $this->setEnvVar('CONVERTER_API', null);

        [$uploadedFile, $inputFilePath] = $this->createUploadedFile('sample.png');
        try {
            $client = static::createClient();
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $uploadedFile,
            ]);

            self::assertResponseStatusCodeSame(503);
            $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('converter_api_not_configured', $payload['error']['code'] ?? null);
            self::assertNotSame('', $payload['error']['requestId'] ?? '');
            self::assertSame($payload['error']['requestId'] ?? null, $client->getResponse()->headers->get('X-Request-Id'));
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertForwardsUpstreamError(): void
    {
        $this->setEnvVar('CONVERTER_API', 'http://converter-api:8081');

        [$uploadedFile, $inputFilePath] = $this->createUploadedFile('sample.png');
        try {
            $mockClient = new MockHttpClient([
                new MockResponse('{"error":{"code":"unsupported_conversion_pair","message":"conversion from png to pdf is not supported"}}', ['http_code' => 415]),
            ]);

            $client = static::createClient();
            static::getContainer()->set(HttpClientInterface::class, $mockClient);
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'pdf',
            ], [
                'file' => $uploadedFile,
            ]);

            self::assertResponseStatusCodeSame(415);
            $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('unsupported_conversion_pair', $payload['error']['code'] ?? null);
            self::assertNotSame('', $payload['error']['requestId'] ?? '');
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertReturnsBadGatewayWhenUpstreamPayloadIsInvalid(): void
    {
        $this->setEnvVar('CONVERTER_API', 'http://converter-api:8081');

        [$uploadedFile, $inputFilePath] = $this->createUploadedFile('sample.png');
        try {
            $mockClient = new MockHttpClient([
                new MockResponse('{"unexpected":"payload"}', ['http_code' => 200]),
            ]);

            $client = static::createClient();
            static::getContainer()->set(HttpClientInterface::class, $mockClient);
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $uploadedFile,
            ]);

            self::assertResponseStatusCodeSame(502);
            $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('invalid_upstream_response', $payload['error']['code'] ?? null);
            self::assertNotSame('', $payload['error']['requestId'] ?? '');
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertReportsUpstreamUnreachable(): void
    {
        $this->setEnvVar('CONVERTER_API', 'http://converter-api:8081');

        [$uploadedFile, $inputFilePath] = $this->createUploadedFile('sample.png');
        try {
            $mockClient = new MockHttpClient(static function (): never {
                throw new TransportException('Connection refused');
            });

            $client = static::createClient();
            static::getContainer()->set(HttpClientInterface::class, $mockClient);
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $uploadedFile,
            ]);

            self::assertResponseStatusCodeSame(503);
            $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('converter_api_unreachable', $payload['error']['code'] ?? null);
            self::assertNotSame('', $payload['error']['requestId'] ?? '');
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertRejectsPayloadWhenUploadedFileIsTooLarge(): void
    {
        $this->setEnvVar('CONVERTER_API', 'http://converter-api:8081');
        $this->setEnvVar('APP_MAX_UPLOAD_BYTES', '4');

        [$uploadedFile, $inputFilePath] = $this->createUploadedFile('sample.png');
        try {
            $client = static::createClient();
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $uploadedFile,
            ]);

            self::assertResponseStatusCodeSame(413);
            $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('payload_too_large', $payload['error']['code'] ?? null);
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertRateLimitsRequests(): void
    {
        $this->setEnvVar('CONVERTER_API', 'http://converter-api:8081');
        $this->setEnvVar('APP_CONVERT_RATE_LIMIT', '1');
        $clientIp = sprintf('203.0.113.%d', random_int(11, 250));

        $calls = 0;
        $mockClient = new MockHttpClient(static function () use (&$calls): MockResponse {
            ++$calls;

            return new MockResponse('{"fileName":"sample.jpg","mimeType":"image/jpeg","contentBase64":"Y29udmVydGVk"}', ['http_code' => 200]);
        });

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $mockClient);

        [$firstUploadedFile, $firstInputFilePath] = $this->createUploadedFile('first.png');
        try {
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $firstUploadedFile,
            ], [
                'REMOTE_ADDR' => $clientIp,
            ]);
            self::assertResponseIsSuccessful();
        } finally {
            @unlink($firstInputFilePath);
        }

        [$secondUploadedFile, $secondInputFilePath] = $this->createUploadedFile('second.png');
        try {
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $secondUploadedFile,
            ], [
                'REMOTE_ADDR' => $clientIp,
            ]);

            self::assertResponseStatusCodeSame(429);
            $payload = json_decode($client->getResponse()->getContent(false), true, 512, JSON_THROW_ON_ERROR);
            self::assertSame('rate_limited', $payload['error']['code'] ?? null);
            self::assertNotNull($client->getResponse()->headers->get('Retry-After'));
            self::assertSame(1, $calls);
        } finally {
            @unlink($secondInputFilePath);
        }
    }

    /**
     * @return array{server: ?string, env: ?string, getenv: ?string}
     */
    private function snapshotEnvVar(string $name): array
    {
        $getenvValue = getenv($name);

        return [
            'server' => $_SERVER[$name] ?? null,
            'env' => $_ENV[$name] ?? null,
            'getenv' => false === $getenvValue ? null : $getenvValue,
        ];
    }

    /**
     * @param array{server: ?string, env: ?string, getenv: ?string} $snapshot
     */
    private function restoreEnvVar(string $name, array $snapshot): void
    {
        if ($snapshot['server'] === null) {
            unset($_SERVER[$name]);
        } else {
            $_SERVER[$name] = $snapshot['server'];
        }

        if ($snapshot['env'] === null) {
            unset($_ENV[$name]);
        } else {
            $_ENV[$name] = $snapshot['env'];
        }

        if ($snapshot['getenv'] === null) {
            putenv($name);
        } else {
            putenv(sprintf('%s=%s', $name, $snapshot['getenv']));
        }
    }

    private function setEnvVar(string $name, ?string $value): void
    {
        if ($value === null) {
            unset($_SERVER[$name]);
            unset($_ENV[$name]);
            putenv($name);

            return;
        }

        $_SERVER[$name] = $value;
        $_ENV[$name] = $value;
        putenv(sprintf('%s=%s', $name, $value));
    }

    /**
     * @return array{0: UploadedFile, 1: string}
     */
    private function createUploadedFile(string $fileName): array
    {
        $inputFilePath = tempnam(sys_get_temp_dir(), 'convert-input-');
        self::assertNotFalse($inputFilePath);
        file_put_contents($inputFilePath, 'test-image-bytes');

        return [
            new UploadedFile(
                $inputFilePath,
                $fileName,
                'application/octet-stream',
                null,
                true,
            ),
            $inputFilePath,
        ];
    }
}
