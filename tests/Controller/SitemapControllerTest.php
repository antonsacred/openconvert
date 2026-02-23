<?php

namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class SitemapControllerTest extends WebTestCase
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

    public function testDynamicSitemapReturnsHttpsUrlsForCurrentHost(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $calls = 0;
        $mockClient = new MockHttpClient(static function (string $method, string $url) use (&$calls): MockResponse {
            if (!str_contains($url, '/v1/conversions')) {
                return new MockResponse('{}', ['http_code' => 404]);
            }

            ++$calls;

            return new MockResponse('{"formats":{"png":["jpg"],"jpg":["png"]}}', ['http_code' => 200]);
        });

        $client = static::createClient();
        static::getContainer()->get('cache.app')->clear();
        $client->disableReboot();
        static::getContainer()->set(HttpClientInterface::class, $mockClient);

        $client->request('GET', '/sitemap.xml', server: ['HTTP_HOST' => 'convert.example.com']);

        self::assertResponseIsSuccessful();
        self::assertResponseHeaderSame('Content-Type', 'application/xml; charset=UTF-8');
        self::assertStringContainsString('public', (string) $client->getResponse()->headers->get('Cache-Control'));
        self::assertStringContainsString('<loc>https://convert.example.com/</loc>', $client->getResponse()->getContent());
        self::assertStringContainsString('<loc>https://convert.example.com/png-converter</loc>', $client->getResponse()->getContent());
        self::assertStringContainsString('<loc>https://convert.example.com/png-to-jpg</loc>', $client->getResponse()->getContent());

        $client->request('GET', '/sitemap.xml', server: ['HTTP_HOST' => 'convert.example.com']);
        self::assertResponseIsSuccessful();
        self::assertSame(1, $calls);
    }

    public function testDynamicSitemapReturns503WhenGeneratorFails(): void
    {
        $this->setConverterApi('http://converter-api-unavailable:8081');

        $mockClient = new MockHttpClient([
            new MockResponse('{}', ['http_code' => 503]),
        ]);

        $client = static::createClient();
        static::getContainer()->get('cache.app')->clear();
        static::getContainer()->set(HttpClientInterface::class, $mockClient);
        $client->request('GET', '/sitemap.xml', server: ['HTTP_HOST' => 'down.example.com']);

        self::assertResponseStatusCodeSame(503);
        self::assertSame('Sitemap is temporarily unavailable.', $client->getResponse()->getContent(false));
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
