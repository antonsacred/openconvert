<?php

namespace App\Tests\Controller;

use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class HomeControllerTest extends WebTestCase
{
    private ?string $originalServerConverterApi = null;
    private ?string $originalEnvConverterApi = null;
    private ?string $originalGetenvConverterApi = null;

    protected function setUp(): void
    {
        parent::setUp();

        $this->originalServerConverterApi = $_SERVER['CONVERTER_API'] ?? null;
        $this->originalEnvConverterApi = $_ENV['CONVERTER_API'] ?? null;
        $getenvValue = getenv('CONVERTER_API');
        $this->originalGetenvConverterApi = false === $getenvValue ? null : $getenvValue;
    }

    protected function tearDown(): void
    {
        $this->setConverterApi($this->originalServerConverterApi);

        if ($this->originalEnvConverterApi === null) {
            unset($_ENV['CONVERTER_API']);
        } else {
            $_ENV['CONVERTER_API'] = $this->originalEnvConverterApi;
        }

        if ($this->originalGetenvConverterApi === null) {
            putenv('CONVERTER_API');
        } else {
            putenv('CONVERTER_API='.$this->originalGetenvConverterApi);
        }

        parent::tearDown();
    }

    public function testHomeUsesConverterFormatsAndKeepsTargetDisabledWithoutSource(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $mockHttpClient = new MockHttpClient([
            new MockResponse('{"formats":{"jpg":["png"],"png":["jpg","webp"]}}', ['http_code' => 200]),
        ]);

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $mockHttpClient);
        $crawler = $client->request('GET', '/');

        self::assertResponseIsSuccessful();
        self::assertGreaterThan(0, $crawler->filter('select[name="from"] option[value="png"]')->count());
        self::assertGreaterThan(0, $crawler->filter('select[name="from"] option[value="jpg"]')->count());

        $toSelect = $crawler->filter('select[name="to"]');
        self::assertCount(1, $toSelect);
        self::assertNotNull($toSelect->attr('disabled'));
    }

    public function testHomeShowsTargetsForSelectedSource(): void
    {
        $this->setConverterApi('http://converter-api:8081');

        $mockHttpClient = new MockHttpClient([
            new MockResponse('{"formats":{"jpg":["png"],"png":["jpg","webp"]}}', ['http_code' => 200]),
        ]);

        $client = static::createClient();
        static::getContainer()->set(HttpClientInterface::class, $mockHttpClient);
        $crawler = $client->request('GET', '/?from=png&to=webp');

        self::assertResponseIsSuccessful();

        $toSelect = $crawler->filter('select[name="to"]');
        self::assertCount(1, $toSelect);
        self::assertNull($toSelect->attr('disabled'));
        self::assertGreaterThan(0, $crawler->filter('select[name="to"] option[value="jpg"]')->count());
        self::assertGreaterThan(0, $crawler->filter('select[name="to"] option[value="webp"]')->count());
        self::assertGreaterThan(0, $crawler->filter('select[name="to"] option[value="webp"][selected]')->count());
    }

    private function setConverterApi(?string $value): void
    {
        if ($value === null) {
            unset($_SERVER['CONVERTER_API']);
            unset($_ENV['CONVERTER_API']);
            putenv('CONVERTER_API');

            return;
        }

        $_SERVER['CONVERTER_API'] = $value;
        $_ENV['CONVERTER_API'] = $value;
        putenv('CONVERTER_API='.$value);
    }
}
