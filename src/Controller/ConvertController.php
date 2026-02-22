<?php

namespace App\Controller;

use App\Service\ConvertService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\File\UploadedFile;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;

final class ConvertController extends AbstractController
{
    public function __construct(
        private readonly ConvertService $convertService,
    ) {
    }

    #[Route('/api/convert', name: 'app_api_convert', methods: ['POST'])]
    public function convert(Request $request): JsonResponse
    {
        $from = strtolower(trim((string) $request->request->get('from', '')));
        $to = strtolower(trim((string) $request->request->get('to', '')));
        $file = $request->files->get('file');
        if ($from === '' || $to === '' || !$file instanceof UploadedFile) {
            return $this->errorResponse(
                Response::HTTP_BAD_REQUEST,
                'invalid_request',
                'from, to and file are required.',
            );
        }

        if (!$file->isValid()) {
            return $this->errorResponse(
                Response::HTTP_BAD_REQUEST,
                'invalid_request',
                'Uploaded file is invalid.',
            );
        }

        $inputBytes = file_get_contents($file->getPathname());
        if ($inputBytes === false) {
            return $this->errorResponse(
                Response::HTTP_BAD_REQUEST,
                'invalid_request',
                'Uploaded file could not be read.',
            );
        }

        $fileName = trim((string) $file->getClientOriginalName());
        if ($fileName === '') {
            $fileName = 'input.'.$from;
        }

        $result = $this->convertService->convert($from, $to, $fileName, base64_encode($inputBytes));

        return new JsonResponse($result->payload(), $result->statusCode());
    }

    private function errorResponse(int $statusCode, string $code, string $message): JsonResponse
    {
        return new JsonResponse([
            'error' => [
                'code' => $code,
                'message' => $message,
            ],
        ], $statusCode);
    }
}
