<?php

namespace App\Controller;

use App\Service\ConversionFormatsService;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Attribute\Route;

final class HomeController extends AbstractController
{
    public function __construct(
        private readonly ConversionFormatsService $conversionFormatsService,
    ) {
    }

    #[Route('/', name: 'app_home', methods: ['GET'])]
    public function __invoke(Request $request): Response
    {
        $formatsBySource = $this->conversionFormatsService->getFormats();

        $selectedFrom = (string) $request->query->get('from', '');
        if (!array_key_exists($selectedFrom, $formatsBySource)) {
            $selectedFrom = '';
        }

        $availableTargets = $selectedFrom === '' ? [] : $formatsBySource[$selectedFrom];

        $selectedTo = (string) $request->query->get('to', '');
        if ( !\in_array($selectedTo, $availableTargets, true)) {
            $selectedTo = '';
        }

        return $this->render('home/index.html.twig', [
            'formats_by_source' => $formatsBySource,
            'available_targets' => $availableTargets,
            'selected_from' => $selectedFrom,
            'selected_to' => $selectedTo,
        ]);
    }
}
