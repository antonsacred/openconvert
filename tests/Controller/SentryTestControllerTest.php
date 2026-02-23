<?php

namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;

final class SentryTestControllerTest extends WebTestCase
{
    public function testErrorEndpointReturnsServerError(): void
    {
        $client = static::createClient();
        $client->request('GET', '/_sentry-test');

        self::assertResponseStatusCodeSame(500);
    }

    public function testWarningEndpointReturnsSuccessResponse(): void
    {
        $client = static::createClient();
        $client->request('GET', '/_sentry-test-warning');

        self::assertResponseIsSuccessful();
        self::assertResponseHeaderSame('Content-Type', 'text/plain; charset=UTF-8');
        self::assertSame('Sentry warning test log sent.', $client->getResponse()->getContent());
    }
}
