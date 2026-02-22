<?php

namespace App\Controller;

use App\Service\ConversionExecutionService;
use Psr\Log\LoggerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Component\HttpFoundation\File\UploadedFile;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\RateLimiter\RateLimiterFactory;
use Symfony\Component\Routing\Attribute\Route;

final class ConvertController extends AbstractController
{
    public function __construct(
        private readonly ConversionExecutionService $conversionExecutionService,
        #[Autowire(service: 'limiter.convert_api')]
        private readonly RateLimiterFactory $convertRateLimiter,
        #[Autowire('%env(int:APP_MAX_UPLOAD_BYTES)%')]
        private readonly int $maxUploadBytes,
        private readonly LoggerInterface $logger,
    ) {
    }

    #[Route('/api/convert', name: 'app_api_convert', methods: ['POST'])]
    public function convert(Request $request): JsonResponse
    {
        $requestId = $this->resolveRequestId($request);
        $requestStart = microtime(true);

        $limiterKey = trim((string) ($request->getClientIp() ?? 'unknown'));
        $rateLimit = $this->convertRateLimiter->create($limiterKey)->consume(1);
        if (!$rateLimit->isAccepted()) {
            $headers = [];
            $retryAfter = $rateLimit->getRetryAfter();
            if ($retryAfter !== null) {
                $retryAfterSeconds = max(1, $retryAfter->getTimestamp() - time());
                $headers['Retry-After'] = (string) $retryAfterSeconds;
            }

            return $this->errorResponse(
                Response::HTTP_TOO_MANY_REQUESTS,
                'rate_limited',
                'Too many conversion requests. Please retry shortly.',
                $requestId,
                $headers,
            );
        }

        $from = strtolower(trim((string) $request->request->get('from', '')));
        $to = strtolower(trim((string) $request->request->get('to', '')));
        $file = $request->files->get('file');
        if ($from === '' || $to === '' || !$file instanceof UploadedFile) {
            return $this->errorResponse(
                Response::HTTP_BAD_REQUEST,
                'invalid_request',
                'from, to and file are required.',
                $requestId,
            );
        }

        if (!$file->isValid()) {
            return $this->errorResponse(
                Response::HTTP_BAD_REQUEST,
                'invalid_request',
                'Uploaded file is invalid.',
                $requestId,
            );
        }

        $fileSize = $file->getSize();
        if (!\is_int($fileSize) || $fileSize < 0) {
            $detectedFileSize = @filesize($file->getPathname());
            $fileSize = \is_int($detectedFileSize) ? $detectedFileSize : 0;
        }

        if ($fileSize > $this->maxUploadBytes) {
            return $this->errorResponse(
                Response::HTTP_REQUEST_ENTITY_TOO_LARGE,
                'payload_too_large',
                sprintf('Uploaded file exceeds %d bytes limit.', $this->maxUploadBytes),
                $requestId,
            );
        }

        $inputBytes = file_get_contents($file->getPathname());
        if ($inputBytes === false) {
            return $this->errorResponse(
                Response::HTTP_BAD_REQUEST,
                'invalid_request',
                'Uploaded file could not be read.',
                $requestId,
            );
        }

        $fileName = trim((string) $file->getClientOriginalName());
        if ($fileName === '') {
            $fileName = 'input.'.$from;
        }

        $result = $this->conversionExecutionService->convert($from, $to, $fileName, base64_encode($inputBytes), $requestId);

        $durationMs = (int) round((microtime(true) - $requestStart) * 1000);
        $statusCode = $result->statusCode();
        $logLevel = $statusCode >= 500 ? 'error' : ($statusCode >= 400 ? 'warning' : 'info');
        $this->logger->log($logLevel, 'convert_request', [
            'request_id' => $requestId,
            'from' => $from,
            'to' => $to,
            'status_code' => $statusCode,
            'duration_ms' => $durationMs,
            'client_ip' => $request->getClientIp(),
        ]);

        return new JsonResponse($result->payload(), $statusCode, [
            'X-Request-Id' => $requestId,
        ]);
    }

    /**
     * @param array<string, string> $headers
     */
    private function errorResponse(int $statusCode, string $code, string $message, string $requestId, array $headers = []): JsonResponse
    {
        $this->logger->warning('convert_request_failed', [
            'request_id' => $requestId,
            'status_code' => $statusCode,
            'error_code' => $code,
            'error_message' => $message,
        ]);

        return new JsonResponse([
            'error' => [
                'code' => $code,
                'message' => $message,
                'requestId' => $requestId,
            ],
        ], $statusCode, [
            ...$headers,
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
