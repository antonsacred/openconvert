<?php

namespace App\Tests\Controller;

use App\Tests\Support\ConverterApiClientOverrideTrait;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class SitemapControllerTest extends WebTestCase
{
    use ConverterApiClientOverrideTrait;

    public function testDynamicSitemapReturnsHttpsUrlsForCurrentHost(): void
    {
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
        $this->overrideConverterApiClient('http://converter-api:8081', $mockClient);
        $client->disableReboot();

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
        $mockClient = new MockHttpClient([
            new MockResponse('{}', ['http_code' => 503]),
        ]);

        $client = static::createClient();
        static::getContainer()->get('cache.app')->clear();
        $this->overrideConverterApiClient('http://converter-api-unavailable:8081', $mockClient);
        $client->request('GET', '/sitemap.xml', server: ['HTTP_HOST' => 'down.example.com']);

        self::assertResponseStatusCodeSame(503);
        self::assertSame('Sitemap is temporarily unavailable.', $client->getResponse()->getContent(false));
    }
}
