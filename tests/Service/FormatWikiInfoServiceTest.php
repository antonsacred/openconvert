<?php

namespace App\Tests\Service;

use App\Service\FormatWikiInfoService;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\ArrayAdapter;
use Symfony\Component\HttpClient\Exception\TransportException;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;

final class FormatWikiInfoServiceTest extends TestCase
{
    public function testGetFormatInfoUsesWikipediaSummaryWhenAvailable(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{"title":"Portable Network Graphics","extract":"PNG is a raster graphics file format.","content_urls":{"desktop":{"page":"https://en.wikipedia.org/wiki/Portable_Network_Graphics"}}}', ['http_code' => 200]),
        ]);

        $service = new FormatWikiInfoService($httpClient, new ArrayAdapter());

        self::assertSame([
            'format' => 'png',
            'label' => 'PNG',
            'title' => 'Portable Network Graphics',
            'summary' => 'PNG is a raster graphics file format.',
            'url' => 'https://en.wikipedia.org/wiki/Portable_Network_Graphics',
        ], $service->getFormatInfo('png'));
    }

    public function testGetFormatInfoFallsBackWhenWikipediaIsUnavailable(): void
    {
        $httpClient = new MockHttpClient(static function (): never {
            throw new TransportException('Connection refused');
        });

        $service = new FormatWikiInfoService($httpClient, new ArrayAdapter());

        $info = $service->getFormatInfo('jpg');

        self::assertSame('jpg', $info['format']);
        self::assertSame('JPG', $info['label']);
        self::assertSame('JPG', $info['title']);
        self::assertStringContainsString('JPG', $info['summary']);
        self::assertSame('https://en.wikipedia.org/wiki/JPEG', $info['url']);
    }
}
