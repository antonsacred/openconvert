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

    public function testConvertReturnsConvertedPayload(): void
    {
        $this->setConverterApi('http://converter-api:8081');

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

                return new MockResponse('{"fileName":"sample.jpg","mimeType":"image/jpeg","contentBase64":"Y29udmVydGVk"}', ['http_code' => 200]);
            });

            $client = static::createClient();
            static::getContainer()->set(HttpClientInterface::class, $mockClient);
            $client->request('POST', '/api/convert', [
                'from' => 'png',
                'to' => 'jpg',
            ], [
                'file' => $uploadedFile,
            ]);

            self::assertResponseIsSuccessful();

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
        $this->setConverterApi(null);

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
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertForwardsUpstreamError(): void
    {
        $this->setConverterApi('http://converter-api:8081');

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
        } finally {
            @unlink($inputFilePath);
        }
    }

    public function testConvertReportsUpstreamUnreachable(): void
    {
        $this->setConverterApi('http://converter-api:8081');

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
        } finally {
            @unlink($inputFilePath);
        }
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
