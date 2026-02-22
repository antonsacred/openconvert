<?php

namespace App\Controller;

use App\Service\ConverterHealthService;
use Psr\Log\LoggerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\Routing\Attribute\Route;

final class HealthController extends AbstractController
{
    public function __construct(
        private readonly ConverterHealthService $converterHealthService,
        private readonly LoggerInterface $logger,
    ) {
    }

    #[Route('/health', name: 'app_health', methods: ['GET'])]
    public function __invoke(Request $request): JsonResponse
    {
        $requestId = $this->resolveRequestId($request);
        $result = $this->converterHealthService->getHealthStatus();

        $this->logger->info('health_check', [
            'request_id' => $requestId,
            'status_code' => $result->statusCode(),
            'health_status' => $result->payload()['status'] ?? 'unknown',
        ]);

        return new JsonResponse($result->payload(), $result->statusCode(), [
            'X-Request-Id' => $requestId,
        ]);
    }

    private function resolveRequestId(Request $request): string
    {
        $incomingRequestId = trim((string) $request->headers->get('X-Request-Id', ''));
        if ($incomingRequestId !== '' && preg_match('/^[A-Za-z0-9._:-]{6,128}$/', $incomingRequestId) === 1) {
            return $incomingRequestId;
        }

        try {
            return bin2hex(random_bytes(16));
        } catch (\Throwable) {
            return str_replace('.', '', uniqid('req', true));
        }
    }
}
