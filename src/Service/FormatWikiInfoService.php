<?php

namespace App\Service;

use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Contracts\Cache\CacheInterface;
use Symfony\Contracts\Cache\ItemInterface;
use Symfony\Contracts\HttpClient\Exception\DecodingExceptionInterface;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class FormatWikiInfoService
{
    private const string CACHE_KEY_PREFIX = 'wiki_info';
    private const int CACHE_TTL_SECONDS = 86400;

    private const array WIKIPEDIA_PAGE_BY_FORMAT = [
        'bmp' => 'BMP_file_format',
        'gif' => 'GIF',
        'jpeg' => 'JPEG',
        'jpg' => 'JPEG',
        'mp3' => 'MP3',
        'mp4' => 'MPEG-4_Part_14',
        'pdf' => 'PDF',
        'png' => 'Portable_Network_Graphics',
        'svg' => 'Scalable_Vector_Graphics',
        'tif' => 'TIFF',
        'tiff' => 'TIFF',
        'webp' => 'WebP',
    ];

    private const array FALLBACK_SUMMARY_BY_FORMAT = [
        'jpg' => 'JPG (JPEG) is a lossy raster image format widely used for photos and web images.',
        'png' => 'PNG is a lossless raster image format that supports transparency and sharp graphics.',
        'webp' => 'WebP is an image format that supports efficient lossy and lossless compression.',
        'pdf' => 'PDF is a document format designed for consistent viewing and printing across platforms.',
        'mp3' => 'MP3 is a compressed digital audio format commonly used for music and spoken content.',
        'mp4' => 'MP4 is a multimedia container format commonly used for video and audio distribution.',
    ];

    public function __construct(
        private readonly HttpClientInterface $httpClient,
        #[Autowire(service: 'cache.app')]
        private readonly CacheInterface $cache,
    ) {
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
    public function getFormatInfo(string $format): array
    {
        $normalizedFormat = strtolower(trim($format));
        if ($normalizedFormat === '') {
            throw new \InvalidArgumentException('Format must not be empty.');
        }

        $pageTitle = self::WIKIPEDIA_PAGE_BY_FORMAT[$normalizedFormat] ?? strtoupper($normalizedFormat).'_file_format';
        $fallbackUrl = 'https://en.wikipedia.org/wiki/'.rawurlencode($pageTitle);

        $cacheKey = self::CACHE_KEY_PREFIX.'_'.sha1($normalizedFormat.'|'.$pageTitle);

        return $this->cache->get($cacheKey, function (ItemInterface $item) use ($fallbackUrl, $normalizedFormat, $pageTitle): array {
            $item->expiresAfter(self::CACHE_TTL_SECONDS);

            $summaryUrl = 'https://en.wikipedia.org/api/rest_v1/page/summary/'.rawurlencode($pageTitle);

            try {
                $response = $this->httpClient->request('GET', $summaryUrl, [
                    'timeout' => 4,
                    'headers' => [
                        'Accept' => 'application/json',
                    ],
                ]);

                $statusCode = $response->getStatusCode();
                if ($statusCode >= 200 && $statusCode < 300) {
                    $payload = $response->toArray(false);

                    if (\is_array($payload)) {
                        $title = $this->extractString($payload, 'title');
                        $summary = $this->extractString($payload, 'extract');
                        $url = $this->extractSummaryUrl($payload) ?? $fallbackUrl;

                        if ($summary !== null && $summary !== '') {
                            return [
                                'format' => $normalizedFormat,
                                'label' => strtoupper($normalizedFormat),
                                'title' => $title !== null && $title !== '' ? $title : strtoupper($normalizedFormat),
                                'summary' => $summary,
                                'url' => $url,
                            ];
                        }
                    }
                }
            } catch (TransportExceptionInterface | DecodingExceptionInterface) {
            }

            return [
                'format' => $normalizedFormat,
                'label' => strtoupper($normalizedFormat),
                'title' => strtoupper($normalizedFormat),
                'summary' => $this->fallbackSummary($normalizedFormat),
                'url' => $fallbackUrl,
            ];
        });
    }

    private function fallbackSummary(string $format): string
    {
        if (isset(self::FALLBACK_SUMMARY_BY_FORMAT[$format])) {
            return self::FALLBACK_SUMMARY_BY_FORMAT[$format];
        }

        return sprintf(
            '%s is a file format used for storing and exchanging digital media or document content.',
            strtoupper($format),
        );
    }

    /**
     * @param array<mixed> $payload
     */
    private function extractString(array $payload, string $key): ?string
    {
        $value = $payload[$key] ?? null;
        if (!\is_string($value)) {
            return null;
        }

        return trim($value);
    }

    /**
     * @param array<mixed> $payload
     */
    private function extractSummaryUrl(array $payload): ?string
    {
        $desktop = $payload['content_urls']['desktop'] ?? null;
        if (!\is_array($desktop)) {
            return null;
        }

        $url = $desktop['page'] ?? null;
        if (!\is_string($url)) {
            return null;
        }

        return trim($url);
    }
}
