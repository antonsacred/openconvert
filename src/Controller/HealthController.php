<?php

namespace App\Controller;

use App\Service\ConvertService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Attribute\Route;

final class HealthController extends AbstractController
{
    public function __construct(
        private readonly ConvertService $convertService,
    ) {
    }

    #[Route('/health', name: 'app_health', methods: ['GET'])]
    public function __invoke(): JsonResponse
    {
        $result = $this->convertService->getHealthStatus();

        return new JsonResponse($result['payload'], $result['statusCode']);
    }
}
