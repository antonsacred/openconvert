<?php

namespace App\Service;

use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Contracts\HttpClient\HttpClientInterface;
use Symfony\Contracts\HttpClient\ResponseInterface;

final class ConverterApiClient
{
    public function __construct(
        private readonly HttpClientInterface $httpClient,
        #[Autowire('%env(default::CONVERTER_API)%')]
        private readonly ?string $converterApi = null,
    ) {
    }

    public function endpoint(string $path): ?string
    {
        $baseUrl = $this->normalizeConverterApi($this->converterApi);
        if ($baseUrl === null) {
            return null;
        }

        return rtrim($baseUrl, '/').$path;
    }

    /**
     * @param array<string, mixed> $options
     */
    public function request(string $method, string $path, array $options = []): ResponseInterface
    {
        $url = $this->endpoint($path);
        if ($url === null) {
            throw new \RuntimeException('CONVERTER_API is not configured.');
        }

        return $this->httpClient->request($method, $url, $options);
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
