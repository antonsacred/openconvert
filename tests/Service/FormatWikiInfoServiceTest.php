<?php

namespace App\Tests\Service;

use App\Service\FormatInfoCatalog;
use App\Service\FormatWikiInfoService;
use PHPUnit\Framework\TestCase;

final class FormatWikiInfoServiceTest extends TestCase
{
    public function testGetFormatInfoReadsFromStoredFile(): void
    {
        $path = tempnam(sys_get_temp_dir(), 'format-info-');
        self::assertNotFalse($path);
        file_put_contents($path, json_encode([
            'formats' => [
                'png' => [
                    'title' => 'Portable Network Graphics',
                    'summary' => 'PNG is a raster graphics file format.',
                    'url' => 'https://example.org/png',
                ],
            ],
        ], JSON_THROW_ON_ERROR));

        $service = new FormatWikiInfoService(new FormatInfoCatalog(), $path);

        try {
            self::assertSame([
                'format' => 'png',
                'label' => 'PNG',
                'title' => 'Portable Network Graphics',
                'summary' => 'PNG is a raster graphics file format.',
                'url' => 'https://example.org/png',
            ], $service->getFormatInfo('png'));
        } finally {
            @unlink($path);
        }
    }

    public function testGetFormatInfoFallsBackWhenStoredDataIsMissing(): void
    {
        $service = new FormatWikiInfoService(
            new FormatInfoCatalog(),
            sys_get_temp_dir().'/missing-format-info-'.bin2hex(random_bytes(4)).'.json',
        );

        $info = $service->getFormatInfo('jpg');

        self::assertSame('jpg', $info['format']);
        self::assertSame('JPG', $info['label']);
        self::assertSame('JPEG', $info['title']);
        self::assertStringContainsString('JPG', $info['summary']);
        self::assertSame('https://en.wikipedia.org/wiki/JPEG', $info['url']);
    }
}
