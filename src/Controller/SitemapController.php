<?php

namespace App\Controller;

use App\Service\SitemapService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;
use Symfony\Contracts\Cache\CacheInterface;
use Symfony\Contracts\Cache\ItemInterface;

final class SitemapController extends AbstractController
{
    private const int CACHE_TTL_SECONDS = 600;

    public function __construct(
        private readonly SitemapService $sitemapService,
        #[Autowire(service: 'cache.app')]
        private readonly CacheInterface $cache,
    ) {
    }

    #[Route('/sitemap.xml', name: 'app_sitemap', methods: ['GET'])]
    public function __invoke(Request $request): Response
    {
        $hostname = $request->getHttpHost();
        $cacheKey = 'sitemap_xml_'.sha1(strtolower(trim($hostname)));

        try {
            $xml = $this->cache->get($cacheKey, function (ItemInterface $item) use ($hostname): string {
                $item->expiresAfter(self::CACHE_TTL_SECONDS);

                return $this->sitemapService->generate($hostname)->xml();
            });
        } catch (\InvalidArgumentException) {
            return new Response('Invalid host.', Response::HTTP_BAD_REQUEST, [
                'Content-Type' => 'text/plain; charset=UTF-8',
            ]);
        } catch (\Throwable) {
            return new Response('Sitemap is temporarily unavailable.', Response::HTTP_SERVICE_UNAVAILABLE, [
                'Content-Type' => 'text/plain; charset=UTF-8',
            ]);
        }

        $response = new Response($xml, Response::HTTP_OK, [
            'Content-Type' => 'application/xml; charset=UTF-8',
        ]);
        $response->setPublic();
        $response->setMaxAge(self::CACHE_TTL_SECONDS);
        $response->setSharedMaxAge(self::CACHE_TTL_SECONDS);

        return $response;
    }
}
