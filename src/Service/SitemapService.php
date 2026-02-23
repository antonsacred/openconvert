<?php

namespace App\Service;

use App\Dto\SitemapResult;
use Symfony\Component\Routing\Generator\UrlGeneratorInterface;

final class SitemapService
{
    public function __construct(
        private readonly ConversionCatalogService $conversionCatalogService,
        private readonly UrlGeneratorInterface $urlGenerator,
    ) {
    }

    public function generate(string $hostname): SitemapResult
    {
        $normalizedHostname = $this->normalizeHostname($hostname);
        if ($normalizedHostname === null) {
            throw new \InvalidArgumentException('Hostname is required and must be a valid host name (without path).');
        }

        $formatsBySource = $this->conversionCatalogService->getFormats();

        $urls = [];
        $urls[$this->absoluteUrl($normalizedHostname, 'app_home')] = true;

        foreach ($formatsBySource as $source => $targets) {
            $normalizedSource = strtolower(trim((string) $source));
            if ($normalizedSource === '') {
                continue;
            }

            $sourceUrl = $this->absoluteUrl($normalizedHostname, 'app_source_converter', [
                'source' => $normalizedSource,
            ]);
            $urls[$sourceUrl] = true;

            foreach ($targets as $target) {
                $normalizedTarget = strtolower(trim((string) $target));
                if ($normalizedTarget === '') {
                    continue;
                }

                $pairUrl = $this->absoluteUrl($normalizedHostname, 'app_pair_converter', [
                    'source' => $normalizedSource,
                    'target' => $normalizedTarget,
                ]);
                $urls[$pairUrl] = true;
            }
        }

        $urlList = array_keys($urls);
        sort($urlList);

        return new SitemapResult($this->buildXml($urlList), \count($urlList));
    }

    /**
     * @param list<string> $urls
     */
    private function buildXml(array $urls): string
    {
        $xml = new \DOMDocument('1.0', 'UTF-8');
        $xml->formatOutput = true;

        $urlset = $xml->createElement('urlset');
        $urlset->setAttribute('xmlns', 'http://www.sitemaps.org/schemas/sitemap/0.9');
        $xml->appendChild($urlset);

        foreach ($urls as $url) {
            $urlElement = $xml->createElement('url');
            $locElement = $xml->createElement('loc');
            $locElement->appendChild($xml->createTextNode($url));
            $urlElement->appendChild($locElement);
            $urlset->appendChild($urlElement);
        }

        $result = $xml->saveXML();
        if (!\is_string($result)) {
            throw new \RuntimeException('Could not render sitemap XML.');
        }

        return $result;
    }

    /**
     * @param array<string, mixed> $parameters
     */
    private function absoluteUrl(string $hostname, string $route, array $parameters = []): string
    {
        $path = $this->urlGenerator->generate($route, $parameters, UrlGeneratorInterface::ABSOLUTE_PATH);

        return 'https://'.$hostname.$path;
    }

    private function normalizeHostname(string $hostname): ?string
    {
        $trimmedHostname = trim($hostname);
        if ($trimmedHostname === '') {
            return null;
        }

        $candidate = str_contains($trimmedHostname, '://')
            ? $trimmedHostname
            : 'https://'.$trimmedHostname;
        $parsed = parse_url($candidate);

        if (!\is_array($parsed)) {
            return null;
        }

        $host = trim((string) ($parsed['host'] ?? ''));
        if ($host === '') {
            return null;
        }

        $path = trim((string) ($parsed['path'] ?? ''));
        if ($path !== '' && $path !== '/') {
            return null;
        }

        if (isset($parsed['query']) || isset($parsed['fragment']) || isset($parsed['user']) || isset($parsed['pass'])) {
            return null;
        }

        $port = $parsed['port'] ?? null;
        if ($port !== null) {
            $port = (int) $port;
            if ($port < 1 || $port > 65535) {
                return null;
            }
        }

        $normalizedHost = strtolower($host);

        return $port === null ? $normalizedHost : sprintf('%s:%d', $normalizedHost, $port);
    }
}
