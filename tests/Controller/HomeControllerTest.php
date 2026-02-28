<?php

namespace App\Tests\Controller;

use App\Tests\Support\ConverterApiClientOverrideTrait;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class HomeControllerTest extends WebTestCase
{
    use ConverterApiClientOverrideTrait;

    public function testHomePageShowsLandingState(): void
    {
        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $this->createMockHttpClient());
        $crawler = $client->request('GET', '/');

        self::assertResponseIsSuccessful();
        self::assertCount(1, $crawler->filter('html[lang="en"]'));
        self::assertCount(1, $crawler->filter('meta[name="description"]'));
        self::assertSame(
            'OpenConvert is an online file converter. Convert audio, video, documents, images, archives, ebooks, spreadsheets and more with one streamlined workflow.',
            $crawler->filter('meta[name="description"]')->attr('content'),
        );
        self::assertSelectorTextContains('h1', 'File Converter');
        self::assertGreaterThan(0, $crawler->filter('select[name="from"] option[value="png"]')->count());
        self::assertNotNull($crawler->filter('select[name="to"]')->attr('disabled'));
        self::assertCount(1, $crawler->filter('[data-controller="upload-queue"]'));
        self::assertCount(1, $crawler->filter('[data-upload-queue-convert-url-value="/api/convert"]'));
        self::assertCount(1, $crawler->filter('input[type="file"][data-upload-queue-target="fileInput"]'));
        self::assertCount(1, $crawler->filter('[data-upload-queue-target="fileList"]'));
        self::assertCount(1, $crawler->filter('[data-upload-queue-target="error"]'));
        self::assertGreaterThan(0, $crawler->filter('button[data-action="click->upload-queue#openFilePicker"]')->count());
        self::assertCount(1, $crawler->filter('button[data-upload-queue-target="downloadAllButton"][data-action="click->upload-queue#downloadAll"]'));
        self::assertCount(1, $crawler->filter('button[data-upload-queue-target="convertButton"]'));
        self::assertCount(1, $crawler->filter('button[aria-label="Open navigation menu"]'));
        self::assertGreaterThan(0, $crawler->filter('.navbar-center a[href="/heic-converter"]')->count());
        self::assertGreaterThan(0, $crawler->filter('.dropdown-content a[href="/heic-converter"]')->count());
        self::assertCount(1, $crawler->filter('head link[rel="preload"][as="font"][type="font/woff2"][href="/fonts/space-grotesk-latin-wght-normal.woff2"]'));
    }

    public function testUploadControlsArePresentOnSourceAndPairPages(): void
    {
        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $this->createMockHttpClient());

        $sourceCrawler = $client->request('GET', '/png-converter');
        self::assertResponseIsSuccessful();
        self::assertCount(1, $sourceCrawler->filter('[data-controller="upload-queue"]'));
        self::assertCount(1, $sourceCrawler->filter('input[type="file"][data-upload-queue-target="fileInput"]'));
        self::assertCount(1, $sourceCrawler->filter('[data-upload-queue-target="fileList"]'));
        self::assertCount(1, $sourceCrawler->filter('button[data-upload-queue-target="downloadAllButton"]'));

        $pairCrawler = $client->request('GET', '/png-to-jpg');
        self::assertResponseIsSuccessful();
        self::assertCount(1, $pairCrawler->filter('[data-controller="upload-queue"]'));
        self::assertCount(1, $pairCrawler->filter('input[type="file"][data-upload-queue-target="fileInput"]'));
        self::assertCount(1, $pairCrawler->filter('[data-upload-queue-target="fileList"]'));
        self::assertCount(1, $pairCrawler->filter('button[data-upload-queue-target="downloadAllButton"]'));
    }

    public function testSourceConverterPageShowsWikiInfoAndTargets(): void
    {
        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $this->createMockHttpClient());
        $crawler = $client->request('GET', '/png-converter');

        self::assertResponseIsSuccessful();
        self::assertCount(1, $crawler->filter('meta[name="description"]'));
        self::assertSame(
            'Convert PNG images while keeping transparency and sharp details for graphics.',
            $crawler->filter('meta[name="description"]')->attr('content'),
        );
        self::assertSelectorTextContains('h1', 'PNG Converter');
        self::assertStringContainsString('Portable Network Graphics', $client->getResponse()->getContent());
        self::assertGreaterThan(0, $crawler->filter('a[href="/png-to-jpg"]')->count());
        self::assertGreaterThan(0, $crawler->filter('a[href="/png-to-webp"]')->count());
    }

    public function testPairConverterPageShowsBothFormatInfos(): void
    {
        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $this->createMockHttpClient());
        $client->request('GET', '/png-to-jpg');

        self::assertResponseIsSuccessful();
        $crawler = $client->getCrawler();
        self::assertCount(1, $crawler->filter('meta[name="description"]'));
        self::assertSame(
            'Convert PNG files to JPG with predictable output quality and fast processing.',
            $crawler->filter('meta[name="description"]')->attr('content'),
        );
        self::assertSelectorTextContains('h1', 'PNG to JPG Converter');
        self::assertStringContainsString('Portable Network Graphics', $client->getResponse()->getContent());
        self::assertStringContainsString('JPEG', $client->getResponse()->getContent());
    }

    public function testInvalidSourceReturnsNotFound(): void
    {
        $client = static::createClient();
        $this->overrideConverterApiClient('http://converter-api:8081', $this->createMockHttpClient());
        $client->request('GET', '/docx-converter');

        self::assertResponseStatusCodeSame(404);
    }

    private function createMockHttpClient(): MockHttpClient
    {
        return new MockHttpClient(static function (string $method, string $url): MockResponse {
            if (str_contains($url, '/v1/conversions')) {
                return new MockResponse('{"formats":{"jpg":["png"],"png":["jpg","webp"]}}', ['http_code' => 200]);
            }

            return new MockResponse('{}', ['http_code' => 404]);
        });
    }
}
