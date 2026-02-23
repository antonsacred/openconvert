<?php

namespace App\Service;

final class FormatInfoCatalog
{
    private const array WIKIPEDIA_PAGE_BY_FORMAT = [
        'avif' => 'AVIF',
        'bmp' => 'BMP_file_format',
        'gif' => 'GIF',
        'heic' => 'High_Efficiency_Image_File_Format',
        'heif' => 'High_Efficiency_Image_File_Format',
        'jpeg' => 'JPEG',
        'jpg' => 'JPEG',
        'magick' => 'ImageMagick',
        'mp3' => 'MP3',
        'mp4' => 'MPEG-4_Part_14',
        'pdf' => 'PDF',
        'png' => 'Portable_Network_Graphics',
        'svg' => 'Scalable_Vector_Graphics',
        'tif' => 'TIFF',
        'tiff' => 'TIFF',
        'webp' => 'WebP',
    ];

    private const array TITLE_BY_FORMAT = [
        'avif' => 'AVIF',
        'bmp' => 'BMP file format',
        'gif' => 'GIF',
        'heic' => 'High Efficiency Image File Format',
        'heif' => 'High Efficiency Image File Format',
        'jpeg' => 'JPEG',
        'jpg' => 'JPEG',
        'magick' => 'ImageMagick',
        'mp3' => 'MP3',
        'mp4' => 'MPEG-4 Part 14',
        'pdf' => 'PDF',
        'png' => 'Portable Network Graphics',
        'svg' => 'Scalable Vector Graphics',
        'tif' => 'TIFF',
        'tiff' => 'TIFF',
        'webp' => 'WebP',
    ];

    private const array FALLBACK_SUMMARY_BY_FORMAT = [
        'avif' => 'AVIF is a modern image format based on AV1 compression with high quality at small file sizes.',
        'bmp' => 'BMP is an uncompressed raster image format designed for compatibility in Windows environments.',
        'gif' => 'GIF is a raster image format known for short animations and indexed-color graphics.',
        'heic' => 'HEIC is a container format for high-efficiency images, commonly used by Apple devices.',
        'heif' => 'HEIF is a container format for high-efficiency images and image sequences.',
        'jpeg' => 'JPEG is a lossy raster image format widely used for photos and web images.',
        'jpg' => 'JPG (JPEG) is a lossy raster image format widely used for photos and web images.',
        'magick' => 'ImageMagick is a software suite and format family used for broad image processing workflows.',
        'mp3' => 'MP3 is a compressed digital audio format commonly used for music and spoken content.',
        'mp4' => 'MP4 is a multimedia container format commonly used for video and audio distribution.',
        'pdf' => 'PDF is a document format designed for consistent viewing and printing across platforms.',
        'png' => 'PNG is a lossless raster image format that supports transparency and sharp graphics.',
        'svg' => 'SVG is a vector graphics format that scales cleanly without losing visual quality.',
        'tif' => 'TIFF is a high-quality image format commonly used for scanning and professional workflows.',
        'tiff' => 'TIFF is a high-quality image format commonly used for scanning and professional workflows.',
        'webp' => 'WebP is an image format that supports efficient lossy and lossless compression.',
    ];

    public function pageTitleFor(string $format): string
    {
        $normalizedFormat = strtolower(trim($format));
        if ($normalizedFormat === '') {
            throw new \InvalidArgumentException('Format must not be empty.');
        }

        return self::WIKIPEDIA_PAGE_BY_FORMAT[$normalizedFormat] ?? strtoupper($normalizedFormat).'_file_format';
    }

    /**
     * @return array{
     *     format: string,
     *     label: string,
     *     title: string,
     *     summary: string,
     *     url: string
     * }
     */
    public function fallbackInfo(string $format): array
    {
        $normalizedFormat = strtolower(trim($format));
        if ($normalizedFormat === '') {
            throw new \InvalidArgumentException('Format must not be empty.');
        }

        $pageTitle = $this->pageTitleFor($normalizedFormat);
        $title = self::TITLE_BY_FORMAT[$normalizedFormat] ?? strtoupper($normalizedFormat);
        $summary = self::FALLBACK_SUMMARY_BY_FORMAT[$normalizedFormat]
            ?? sprintf(
                '%s is a file format used for storing and exchanging digital media or document content.',
                strtoupper($normalizedFormat),
            );

        return [
            'format' => $normalizedFormat,
            'label' => strtoupper($normalizedFormat),
            'title' => $title,
            'summary' => $summary,
            'url' => $this->buildWikipediaUrl($pageTitle),
        ];
    }

    public function buildWikipediaSummaryUrl(string $pageTitle): string
    {
        return 'https://en.wikipedia.org/api/rest_v1/page/summary/'.rawurlencode($pageTitle);
    }

    public function buildWikipediaUrl(string $pageTitle): string
    {
        return 'https://en.wikipedia.org/wiki/'.rawurlencode($pageTitle);
    }
}
