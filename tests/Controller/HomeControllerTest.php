<?php

namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class HomeControllerTest extends WebTestCase
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

    public function testHomePageShowsLandingState(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $this->createMockHttpClient());
        $crawler = $client->request('GET', '/');

        self::assertResponseIsSuccessful();
        self::assertSelectorTextContains('h1', 'File Converter');
        self::assertGreaterThan(0, $crawler->filter('select[name="from"] option[value="png"]')->count());
        self::assertNotNull($crawler->filter('select[name="to"]')->attr('disabled'));
        self::assertCount(1, $crawler->filter('[data-controller="upload-queue"]'));
        self::assertCount(1, $crawler->filter('[data-upload-queue-convert-url-value="/api/convert"]'));
        self::assertCount(1, $crawler->filter('input[type="file"][data-upload-queue-target="fileInput"]'));
        self::assertCount(1, $crawler->filter('[data-upload-queue-target="fileList"]'));
        self::assertCount(1, $crawler->filter('[data-upload-queue-target="error"]'));
        self::assertGreaterThan(0, $crawler->filter('button[data-action="click->upload-queue#openFilePicker"]')->count());
        self::assertCount(1, $crawler->filter('button[data-upload-queue-target="convertButton"]'));
    }

    public function testUploadControlsArePresentOnSourceAndPairPages(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $this->createMockHttpClient());

        $sourceCrawler = $client->request('GET', '/png-converter');
        self::assertResponseIsSuccessful();
        self::assertCount(1, $sourceCrawler->filter('[data-controller="upload-queue"]'));
        self::assertCount(1, $sourceCrawler->filter('input[type="file"][data-upload-queue-target="fileInput"]'));
        self::assertCount(1, $sourceCrawler->filter('[data-upload-queue-target="fileList"]'));

        $pairCrawler = $client->request('GET', '/png-to-jpg');
        self::assertResponseIsSuccessful();
        self::assertCount(1, $pairCrawler->filter('[data-controller="upload-queue"]'));
        self::assertCount(1, $pairCrawler->filter('input[type="file"][data-upload-queue-target="fileInput"]'));
        self::assertCount(1, $pairCrawler->filter('[data-upload-queue-target="fileList"]'));
    }

    public function testSourceConverterPageShowsWikiInfoAndTargets(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $this->createMockHttpClient());
        $crawler = $client->request('GET', '/png-converter');

        self::assertResponseIsSuccessful();
        self::assertSelectorTextContains('h1', 'PNG Converter');
        self::assertStringContainsString('Portable Network Graphics', $client->getResponse()->getContent());
        self::assertGreaterThan(0, $crawler->filter('a[href="/png-to-jpg"]')->count());
        self::assertGreaterThan(0, $crawler->filter('a[href="/png-to-webp"]')->count());
    }

    public function testPairConverterPageShowsBothFormatInfos(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $this->createMockHttpClient());
        $client->request('GET', '/png-to-jpg');

        self::assertResponseIsSuccessful();
        self::assertSelectorTextContains('h1', 'PNG to JPG Converter');
        self::assertStringContainsString('Portable Network Graphics', $client->getResponse()->getContent());
        self::assertStringContainsString('JPEG', $client->getResponse()->getContent());
    }

    public function testInvalidSourceReturnsNotFound(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $this->createMockHttpClient());
        $client->request('GET', '/docx-converter');

        self::assertResponseStatusCodeSame(404);
    }

    private function createMockHttpClient(): MockHttpClient
    {
        return new MockHttpClient(static function (string $method, string $url): MockResponse {
            if (str_contains($url, '/v1/conversions')) {
                return new MockResponse('{"formats":{"jpg":["png"],"png":["jpg","webp"]}}', ['http_code' => 200]);
            }

            if (str_contains($url, '/page/summary/Portable_Network_Graphics')) {
                return new MockResponse('{"title":"Portable Network Graphics","extract":"PNG is a raster graphics file format.","content_urls":{"desktop":{"page":"https://en.wikipedia.org/wiki/Portable_Network_Graphics"}}}', ['http_code' => 200]);
            }

            if (str_contains($url, '/page/summary/JPEG')) {
                return new MockResponse('{"title":"JPEG","extract":"JPEG is a commonly used method of lossy compression for digital images.","content_urls":{"desktop":{"page":"https://en.wikipedia.org/wiki/JPEG"}}}', ['http_code' => 200]);
            }

            if (str_contains($url, '/page/summary/WebP')) {
                return new MockResponse('{"title":"WebP","extract":"WebP is an image format employing both lossy and lossless compression.","content_urls":{"desktop":{"page":"https://en.wikipedia.org/wiki/WebP"}}}', ['http_code' => 200]);
            }

            return new MockResponse('{}', ['http_code' => 404]);
        });
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
