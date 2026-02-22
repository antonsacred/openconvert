<?php

namespace App\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;
use Symfony\Contracts\HttpClient\Exception\TransportExceptionInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class HealthController extends AbstractController
{
    public function __construct(
        private readonly HttpClientInterface $httpClient,
    ) {
    }

    #[Route('/health', name: 'app_health', methods: ['GET'])]
    public function __invoke(): JsonResponse
    {
        $converterApi = $this->resolveConverterApi();
        if ($converterApi === null) {
            return new JsonResponse([
                'status' => 'degraded',
                'checks' => [
                    'converter_api' => [
                        'status' => 'not_configured',
                        'message' => 'CONVERTER_API is not configured.',
                    ],
                ],
            ], Response::HTTP_SERVICE_UNAVAILABLE);
        }

        $healthUrl = rtrim($converterApi, '/').'/health';

        try {
            $response = $this->httpClient->request('GET', $healthUrl, [
                'timeout' => 2,
            ]);
            $statusCode = $response->getStatusCode();
            if ($statusCode < 200 || $statusCode >= 300) {
                return new JsonResponse([
                    'status' => 'degraded',
                    'checks' => [
                        'converter_api' => [
                            'status' => 'not_running',
                            'message' => 'Converter API health endpoint did not return a successful status.',
                            'url' => $healthUrl,
                            'http_status' => $statusCode,
                        ],
                    ],
                ], Response::HTTP_SERVICE_UNAVAILABLE);
            }
        } catch (TransportExceptionInterface $exception) {
            return new JsonResponse([
                'status' => 'degraded',
                'checks' => [
                    'converter_api' => [
                        'status' => 'not_running',
                        'message' => 'Converter API health endpoint is not reachable.',
                        'url' => $healthUrl,
                    ],
                ],
            ], Response::HTTP_SERVICE_UNAVAILABLE);
        }

        return new JsonResponse([
            'status' => 'ok',
            'checks' => [
                'converter_api' => [
                    'status' => 'up',
                    'url' => $healthUrl,
                ],
            ],
        ]);
    }

    private function resolveConverterApi(): ?string
    {
        $candidates = [
            $_SERVER['CONVERTER_API'] ?? null,
            $_ENV['CONVERTER_API'] ?? null,
            getenv('CONVERTER_API') ?: null,
        ];

        foreach ($candidates as $candidate) {
            if (!\is_string($candidate)) {
                continue;
            }

            $trimmed = trim($candidate);
            if ($trimmed !== '') {
                return $trimmed;
            }
        }

        return null;
    }
}
