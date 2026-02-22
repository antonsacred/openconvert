<?php

namespace App\Controller;

use App\Service\ConversionCatalogService;
use App\Service\FormatWikiInfoService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;

final class HomeController extends AbstractController
{
    private const string SOURCE_TOKEN = 'sourceplaceholder';
    private const string TARGET_TOKEN = 'targetplaceholder';

    private const array HERO_COPY_BY_FORMAT = [
        'jpg' => 'Convert JPG images for lighter files and broad compatibility across web and apps.',
        'png' => 'Convert PNG images while keeping transparency and sharp details for graphics.',
        'webp' => 'Convert WebP files for modern compression workflows and fast image delivery.',
    ];

    public function __construct(
        private readonly ConversionCatalogService $conversionCatalogService,
        private readonly FormatWikiInfoService $formatWikiInfoService,
    ) {
    }

    #[Route('/', name: 'app_home', methods: ['GET'])]
    public function home(): Response
    {
        return $this->renderConverterPage('home');
    }

    #[Route('/{source}-converter', name: 'app_source_converter', methods: ['GET'], requirements: ['source' => '[a-z0-9]+'])]
    public function sourceConverter(string $source): Response
    {
        return $this->renderConverterPage('source', $source);
    }

    #[Route('/{source}-to-{target}', name: 'app_pair_converter', methods: ['GET'], requirements: ['source' => '[a-z0-9]+', 'target' => '[a-z0-9]+'])]
    public function pairConverter(string $source, string $target): Response
    {
        return $this->renderConverterPage('pair', $source, $target);
    }

    private function renderConverterPage(string $pageMode, ?string $source = null, ?string $target = null): Response
    {
        $formatsBySource = $this->loadFormatsBySource();
        $selectedFrom = $source === null ? '' : strtolower(trim($source));
        $selectedTo = $target === null ? '' : strtolower(trim($target));

        if ($pageMode !== 'home' && $formatsBySource !== [] && !array_key_exists($selectedFrom, $formatsBySource)) {
            throw $this->createNotFoundException();
        }

        $availableTargets = [];
        if ($selectedFrom !== '' && array_key_exists($selectedFrom, $formatsBySource)) {
            $availableTargets = $formatsBySource[$selectedFrom];
        }

        if ($pageMode === 'pair' && $formatsBySource !== [] && !\in_array($selectedTo, $availableTargets, true)) {
            throw $this->createNotFoundException();
        }

        if ($pageMode !== 'pair') {
            $selectedTo = '';
        }

        $heroTitle = $this->buildHeroTitle($pageMode, $selectedFrom, $selectedTo);
        $heroDescription = $this->buildHeroDescription($pageMode, $selectedFrom, $selectedTo);

        $sourceWikiInfo = null;
        if ($selectedFrom !== '' && $pageMode !== 'home') {
            $sourceWikiInfo = $this->formatWikiInfoService->getFormatInfo($selectedFrom);
        }

        $targetWikiInfo = null;
        if ($selectedTo !== '' && $pageMode === 'pair') {
            $targetWikiInfo = $this->formatWikiInfoService->getFormatInfo($selectedTo);
        }

        return $this->render('home/index.html.twig', [
            'page_mode' => $pageMode,
            'page_title' => $heroTitle,
            'hero_title' => $heroTitle,
            'hero_description' => $heroDescription,
            'formats_by_source' => $formatsBySource,
            'available_targets' => $availableTargets,
            'selected_from' => $selectedFrom,
            'selected_to' => $selectedTo,
            'source_wiki_info' => $sourceWikiInfo,
            'target_wiki_info' => $targetWikiInfo,
            'source_page_template' => $this->generateUrl('app_source_converter', ['source' => self::SOURCE_TOKEN]),
            'pair_page_template' => $this->generateUrl('app_pair_converter', [
                'source' => self::SOURCE_TOKEN,
                'target' => self::TARGET_TOKEN,
            ]),
        ]);
    }

    /**
     * @return array<string, list<string>>
     */
    private function loadFormatsBySource(): array
    {
        try {
            return $this->conversionCatalogService->getFormats();
        } catch (\Throwable) {
            return [];
        }
    }

    private function buildHeroTitle(string $pageMode, string $source, string $target): string
    {
        if ($pageMode === 'pair') {
            return strtoupper($source).' to '.strtoupper($target).' Converter';
        }

        if ($pageMode === 'source') {
            return strtoupper($source).' Converter';
        }

        return 'File Converter';
    }

    private function buildHeroDescription(string $pageMode, string $source, string $target): string
    {
        if ($pageMode === 'pair') {
            return sprintf(
                'Convert %s files to %s with predictable output quality and fast processing.',
                strtoupper($source),
                strtoupper($target),
            );
        }

        if ($pageMode === 'source') {
            return self::HERO_COPY_BY_FORMAT[$source]
                ?? sprintf(
                    'Convert %s files into other supported formats with one streamlined workflow.',
                    strtoupper($source),
                );
        }

        return 'OpenConvert is an online file converter. Convert audio, video, documents, images, archives, ebooks, spreadsheets and more with one streamlined workflow.';
    }
}
