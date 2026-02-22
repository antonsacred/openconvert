<?php

namespace App\Service;

use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Contracts\Cache\CacheInterface;
use Symfony\Contracts\Cache\ItemInterface;
use Symfony\Contracts\HttpClient\Exception\DecodingExceptionInterface;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;

final class ConversionCatalogService
{
    private const string CACHE_KEY_PREFIX = 'converter_api_v1_conversions';
    private const int CACHE_TTL_SECONDS = 300;

    public function __construct(
        private readonly ConverterApiClient $converterApiClient,
        #[Autowire(service: 'cache.app')]
        private readonly CacheInterface $cache,
    ) {
    }

    /**
     * @return array<string, list<string>>
     */
    public function getFormats(): array
    {
        $conversionsUrl = $this->converterApiClient->endpoint('/v1/conversions');
        if ($conversionsUrl === null) {
            throw new \RuntimeException('CONVERTER_API is not configured.');
        }

        $cacheKey = self::CACHE_KEY_PREFIX.'_'.sha1($conversionsUrl);

        return $this->cache->get($cacheKey, function (ItemInterface $item): array {
            $item->expiresAfter(self::CACHE_TTL_SECONDS);

            try {
                $response = $this->converterApiClient->request('GET', '/v1/conversions', [
                    'timeout' => 4,
                ]);
                $statusCode = $response->getStatusCode();
                if ($statusCode < 200 || $statusCode >= 300) {
                    throw new \RuntimeException(sprintf(
                        'Converter API returned unexpected status code %d for /v1/conversions.',
                        $statusCode,
                    ));
                }

                $payload = $response->toArray(false);
            } catch (TransportExceptionInterface | DecodingExceptionInterface $exception) {
                throw new \RuntimeException('Failed to fetch conversions from converter API.', previous: $exception);
            }

            if (!\is_array($payload)) {
                throw new \UnexpectedValueException('Invalid converter response: expected JSON object.');
            }

            return $this->parseFormats($payload);
        });
    }

    /**
     * @param array<mixed> $payload
     *
     * @return array<string, list<string>>
     */
    private function parseFormats(array $payload): array
    {
        if (!isset($payload['formats']) || !\is_array($payload['formats'])) {
            throw new \UnexpectedValueException('Invalid converter response: "formats" field is required.');
        }

        $formats = [];
        foreach ($payload['formats'] as $source => $targets) {
            if (!\is_string($source)) {
                throw new \UnexpectedValueException('Invalid converter response: format keys must be strings.');
            }

            $normalizedSource = strtolower(trim($source));
            if ($normalizedSource === '') {
                throw new \UnexpectedValueException('Invalid converter response: format keys cannot be empty.');
            }

            if (!\is_array($targets)) {
                throw new \UnexpectedValueException(sprintf(
                    'Invalid converter response: targets for "%s" must be an array.',
                    $normalizedSource,
                ));
            }

            $normalizedTargets = [];
            foreach ($targets as $target) {
                if (!\is_string($target)) {
                    throw new \UnexpectedValueException(sprintf(
                        'Invalid converter response: target format for "%s" must be a string.',
                        $normalizedSource,
                    ));
                }

                $normalizedTarget = strtolower(trim($target));
                if ($normalizedTarget === '') {
                    throw new \UnexpectedValueException(sprintf(
                        'Invalid converter response: target format for "%s" cannot be empty.',
                        $normalizedSource,
                    ));
                }

                $normalizedTargets[] = $normalizedTarget;
            }

            $formats[$normalizedSource] = array_values(array_unique($normalizedTargets));
        }

        ksort($formats);

        return $formats;
    }
}
