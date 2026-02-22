<?php

namespace App\Service;

use App\Dto\HealthStatusResult;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;

final class ConverterHealthService
{
    public function __construct(
        private readonly ConverterApiClient $converterApiClient,
    ) {
    }

    public function getHealthStatus(): HealthStatusResult
    {
        $healthUrl = $this->converterApiClient->endpoint('/health');
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
            $response = $this->converterApiClient->request('GET', '/health', [
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
}
