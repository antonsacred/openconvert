<?php

namespace App\Controller;

use Psr\Log\LoggerInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;

final class SentryTestController extends AbstractController
{
    public function __construct(
        private readonly LoggerInterface $logger,
    ) {
    }

    #[Route('/_sentry-test', name: 'app_sentry_test', methods: ['GET'])]
    public function testLog(): Response
    {
        // Tests Monolog integration logs to Sentry.
        $this->logger->error('My custom logged error.', ['some' => 'Context Data']);

        // Tests uncaught exception logging to Sentry.
        throw new \RuntimeException('Example exception.');
    }

    #[Route('/_sentry-test-warning', name: 'app_sentry_test_warning', methods: ['GET'])]
    public function testWarningLog(): Response
    {
        $this->logger->warning('My custom logged warning.', ['some' => 'Context Data']);

        return new Response('Sentry warning test log sent.', Response::HTTP_OK, [
            'Content-Type' => 'text/plain; charset=UTF-8',
        ]);
    }
}
