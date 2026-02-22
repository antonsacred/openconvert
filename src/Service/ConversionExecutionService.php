<?php

namespace App\Service;

use App\Dto\ConvertResult;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;

final class ConversionExecutionService
{
    public function __construct(
        private readonly ConverterApiClient $converterApiClient,
    ) {
    }

    public function convert(string $from, string $to, string $fileName, string $contentBase64, string $requestId): ConvertResult
    {
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
            $upstreamResponse = $this->converterApiClient->request('POST', '/v1/convert', [
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
        } catch (\RuntimeException $exception) {
            if ($exception->getMessage() === 'CONVERTER_API is not configured.') {
                return $this->errorResult(
                    503,
                    'converter_api_not_configured',
                    'CONVERTER_API is not configured.',
                    $requestId,
                );
            }

            throw $exception;
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
}
