<?php

namespace App\Tests\Command;

use App\Command\GenerateSitemapCommand;
use App\Service\ConversionCatalogService;
use App\Service\ConverterApiClient;
use App\Service\SitemapService;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\ArrayAdapter;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Tester\CommandTester;
use Symfony\Component\HttpClient\MockHttpClient;
use Symfony\Component\HttpClient\Response\MockResponse;
use Symfony\Component\Routing\Generator\UrlGenerator;
use Symfony\Component\Routing\Route;
use Symfony\Component\Routing\RouteCollection;
use Symfony\Component\Routing\RequestContext;

final class GenerateSitemapCommandTest extends TestCase
{
    public function testGenerateSitemapCreatesAllPageUrls(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{"formats":{"png":["jpg","webp"],"jpg":["png"]}}', ['http_code' => 200]),
        ]);

        $conversionCatalogService = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $command = new GenerateSitemapCommand(
            new SitemapService($conversionCatalogService, $this->createUrlGenerator()),
            '/tmp/sitemap.xml',
        );

        $tester = new CommandTester($command);
        $outputPath = sys_get_temp_dir().'/sitemap-'.bin2hex(random_bytes(8)).'.xml';

        try {
            self::assertSame(Command::SUCCESS, $tester->execute([
                '--hostname' => 'convert.example.com',
                '--output' => $outputPath,
            ]));

            self::assertFileExists($outputPath);
            $xml = (string) file_get_contents($outputPath);
            self::assertStringContainsString('<urlset', $xml);
            self::assertStringContainsString('<loc>https://convert.example.com/</loc>', $xml);
            self::assertStringContainsString('<loc>https://convert.example.com/png-converter</loc>', $xml);
            self::assertStringContainsString('<loc>https://convert.example.com/jpg-converter</loc>', $xml);
            self::assertStringContainsString('<loc>https://convert.example.com/png-to-jpg</loc>', $xml);
            self::assertStringContainsString('<loc>https://convert.example.com/png-to-webp</loc>', $xml);
            self::assertStringContainsString('<loc>https://convert.example.com/jpg-to-png</loc>', $xml);

            preg_match_all('/<loc>/', $xml, $matches);
            self::assertCount(6, $matches[0]);
        } finally {
            @unlink($outputPath);
        }
    }

    public function testGenerateSitemapFailsWhenFormatsCannotBeLoaded(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{}', ['http_code' => 503]),
        ]);

        $conversionCatalogService = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $command = new GenerateSitemapCommand(
            new SitemapService($conversionCatalogService, $this->createUrlGenerator()),
            '/tmp/sitemap.xml',
        );

        $tester = new CommandTester($command);
        $outputPath = sys_get_temp_dir().'/sitemap-'.bin2hex(random_bytes(8)).'.xml';

        try {
            self::assertSame(Command::FAILURE, $tester->execute([
                '--hostname' => 'convert.example.com',
                '--output' => $outputPath,
            ]));
            self::assertFileDoesNotExist($outputPath);
        } finally {
            @unlink($outputPath);
        }
    }

    public function testGenerateSitemapFailsForInvalidHostname(): void
    {
        $httpClient = new MockHttpClient([
            new MockResponse('{"formats":{"png":["jpg"]}}', ['http_code' => 200]),
        ]);

        $conversionCatalogService = new ConversionCatalogService(
            new ConverterApiClient($httpClient, 'http://converter-api:8081'),
            new ArrayAdapter(),
        );

        $command = new GenerateSitemapCommand(
            new SitemapService($conversionCatalogService, $this->createUrlGenerator()),
            '/tmp/sitemap.xml',
        );

        $tester = new CommandTester($command);
        $outputPath = sys_get_temp_dir().'/sitemap-'.bin2hex(random_bytes(8)).'.xml';

        try {
            self::assertSame(Command::FAILURE, $tester->execute([
                '--hostname' => 'convert.example.com/path',
                '--output' => $outputPath,
            ]));
            self::assertFileDoesNotExist($outputPath);
        } finally {
            @unlink($outputPath);
        }
    }

    private function createUrlGenerator(): UrlGenerator
    {
        $routes = new RouteCollection();
        $routes->add('app_home', new Route('/'));
        $routes->add('app_source_converter', new Route('/{source}-converter'));
        $routes->add('app_pair_converter', new Route('/{source}-to-{target}'));

        $context = new RequestContext();
        $context->setScheme('http');
        $context->setHost('router-context.test');

        return new UrlGenerator($routes, $context);
    }
}
