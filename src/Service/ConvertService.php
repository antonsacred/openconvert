<?php

namespace App\Service;

use App\Dto\ConvertResult;
use App\Dto\HealthStatusResult;
use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Contracts\Cache\CacheInterface;
use Symfony\Contracts\Cache\ItemInterface;
use Symfony\Contracts\HttpClient\Exception\DecodingExceptionInterface;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class ConvertService
{
    private const string CACHE_KEY_PREFIX = 'converter_api_v1_conversions';
    private const int CACHE_TTL_SECONDS = 300;

    public function __construct(
        private readonly HttpClientInterface $httpClient,
        #[Autowire(service: 'cache.app')]
        private readonly CacheInterface $cache,
        #[Autowire('%env(default::CONVERTER_API)%')]
        private readonly ?string $converterApi = null,
    ) {
    }

    /**
     * @return array<string, list<string>>
     */
    public function getFormats(): array
    {
        $conversionsUrl = $this->buildUrl('/v1/conversions');
        if ($conversionsUrl === null) {
            throw new \RuntimeException('CONVERTER_API is not configured.');
        }

        $cacheKey = self::CACHE_KEY_PREFIX.'_'.sha1($conversionsUrl);

        return $this->cache->get($cacheKey, function (ItemInterface $item) use ($conversionsUrl): array {
            $item->expiresAfter(self::CACHE_TTL_SECONDS);

            try {
                $response = $this->httpClient->request('GET', $conversionsUrl, [
                    'timeout' => 4,
                ]);
                $statusCode = $response->getStatusCode();
                if ($statusCode < 200 || $statusCode >= 300) {
                    throw new \RuntimeException(sprintf(
                        'Converter API returned unexpected status code %d for %s.',
                        $statusCode,
                        $conversionsUrl,
                    ));
                }

                $payload = $response->toArray(false);
            } catch (TransportExceptionInterface | DecodingExceptionInterface $exception) {
                throw new \RuntimeException(
                    sprintf('Failed to fetch conversions from %s.', $conversionsUrl),
                    previous: $exception,
                );
            }

            if (!\is_array($payload)) {
                throw new \UnexpectedValueException('Invalid converter response: expected JSON object.');
            }

            return $this->parseFormats($payload);
        });
    }

    public function getHealthStatus(): HealthStatusResult
    {
        $healthUrl = $this->buildUrl('/health');
        if ($healthUrl === null) {
            return new HealthStatusResult(503, [
                'status' => 'degraded',
                'checks' => [
                    'converter_api' => [
                        'status' => 'not_configured',
                        'message' => 'CONVERTER_API is not configured.',
                    ],
                ],
            ]);
        }

        try {
            $response = $this->httpClient->request('GET', $healthUrl, [
                'timeout' => 2,
            ]);
            $statusCode = $response->getStatusCode();
            if ($statusCode < 200 || $statusCode >= 300) {
                return new HealthStatusResult(503, [
                    'status' => 'degraded',
                    'checks' => [
                        'converter_api' => [
                            'status' => 'not_running',
                            'message' => 'Converter API health endpoint did not return a successful status.',
                            'url' => $healthUrl,
                            'http_status' => $statusCode,
                        ],
                    ],
                ]);
            }
        } catch (TransportExceptionInterface) {
            return new HealthStatusResult(503, [
                'status' => 'degraded',
                'checks' => [
                    'converter_api' => [
                        'status' => 'not_running',
                        'message' => 'Converter API health endpoint is not reachable.',
                        'url' => $healthUrl,
                    ],
                ],
            ]);
        }

        return new HealthStatusResult(200, [
            'status' => 'ok',
            'checks' => [
                'converter_api' => [
                    'status' => 'up',
                    'url' => $healthUrl,
                ],
            ],
        ]);
    }

    public function convert(string $from, string $to, string $fileName, string $contentBase64, string $requestId): ConvertResult
    {
        $convertUrl = $this->buildUrl('/v1/convert');
        if ($convertUrl === null) {
            return $this->errorResult(
                503,
                'converter_api_not_configured',
                'CONVERTER_API is not configured.',
                $requestId,
            );
        }

        $normalizedFrom = strtolower(trim($from));
        $normalizedTo = strtolower(trim($to));
        $normalizedFileName = trim($fileName);

        if ($normalizedFrom === '' || $normalizedTo === '' || $normalizedFileName === '' || $contentBase64 === '') {
            return $this->errorResult(
                400,
                'invalid_request',
                'from, to, fileName and contentBase64 are required.',
                $requestId,
            );
        }

        try {
            $upstreamResponse = $this->httpClient->request('POST', $convertUrl, [
                'timeout' => 30,
                'headers' => [
                    'X-Request-Id' => $requestId,
                ],
                'json' => [
                    'from' => $normalizedFrom,
                    'to' => $normalizedTo,
                    'fileName' => $normalizedFileName,
                    'contentBase64' => $contentBase64,
                ],
            ]);
        } catch (TransportExceptionInterface) {
            return $this->errorResult(
                503,
                'converter_api_unreachable',
                'Converter API is not reachable.',
                $requestId,
            );
        }

        $upstreamStatusCode = $upstreamResponse->getStatusCode();
        $upstreamBody = $upstreamResponse->getContent(false);

        $upstreamPayload = null;
        if ($upstreamBody !== '') {
            try {
                $upstreamPayload = json_decode($upstreamBody, true, 512, JSON_THROW_ON_ERROR);
            } catch (\JsonException) {
                if ($upstreamStatusCode >= 200 && $upstreamStatusCode < 300) {
                    return $this->errorResult(
                        502,
                        'invalid_upstream_response',
                        'Converter API returned invalid JSON payload.',
                        $requestId,
                    );
                }
            }
        }

        if ($upstreamStatusCode < 200 || $upstreamStatusCode >= 300) {
            $errorCode = 'conversion_failed';
            $errorMessage = sprintf('Converter API returned HTTP %d.', $upstreamStatusCode);
            if (\is_array($upstreamPayload) && isset($upstreamPayload['error']) && \is_array($upstreamPayload['error'])) {
                if (isset($upstreamPayload['error']['code']) && \is_string($upstreamPayload['error']['code']) && trim($upstreamPayload['error']['code']) !== '') {
                    $errorCode = trim($upstreamPayload['error']['code']);
                }
                if (isset($upstreamPayload['error']['message']) && \is_string($upstreamPayload['error']['message']) && trim($upstreamPayload['error']['message']) !== '') {
                    $errorMessage = trim($upstreamPayload['error']['message']);
                }
            }

            return $this->errorResult($upstreamStatusCode, $errorCode, $errorMessage, $requestId);
        }

        if (!\is_array($upstreamPayload)
            || !isset($upstreamPayload['fileName'])
            || !\is_string($upstreamPayload['fileName'])
            || !isset($upstreamPayload['mimeType'])
            || !\is_string($upstreamPayload['mimeType'])
            || !isset($upstreamPayload['contentBase64'])
            || !\is_string($upstreamPayload['contentBase64'])
        ) {
            return $this->errorResult(
                502,
                'invalid_upstream_response',
                'Converter API response is missing required fields.',
                $requestId,
            );
        }

        return new ConvertResult(200, [
            'fileName' => $upstreamPayload['fileName'],
            'mimeType' => $upstreamPayload['mimeType'],
            'contentBase64' => $upstreamPayload['contentBase64'],
        ]);
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

    private function errorResult(int $statusCode, string $code, string $message, string $requestId): ConvertResult
    {
        return new ConvertResult($statusCode, [
            'error' => [
                'code' => $code,
                'message' => $message,
                'requestId' => $requestId,
            ],
        ]);
    }

    private function buildUrl(string $path): ?string
    {
        $baseUrl = $this->normalizeConverterApi($this->converterApi);
        if ($baseUrl === null) {
            return null;
        }

        return rtrim($baseUrl, '/').$path;
    }

    private function normalizeConverterApi(?string $converterApi): ?string
    {
        if ($converterApi === null) {
            return null;
        }

        $trimmed = trim($converterApi);

        return $trimmed === '' ? null : $trimmed;
    }
}
