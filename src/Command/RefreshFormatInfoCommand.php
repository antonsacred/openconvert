<?php

namespace App\Command;

use App\Service\ConversionCatalogService;
use App\Service\FormatInfoCatalog;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Component\DependencyInjection\Attribute\Autowire;
use Symfony\Contracts\HttpClient\HttpClientInterface;

#[AsCommand(
    name: 'app:format-info:refresh',
    description: 'Fetches format information and stores it in a local JSON file.',
)]
final class RefreshFormatInfoCommand extends Command
{
    public function __construct(
        private readonly ConversionCatalogService $conversionCatalogService,
        private readonly HttpClientInterface $httpClient,
        private readonly FormatInfoCatalog $formatInfoCatalog,
        #[Autowire('%kernel.project_dir%/config/format_info_data.json')]
        private readonly string $defaultOutputPath,
    ) {
        parent::__construct();
    }

    protected function configure(): void
    {
        $this->addOption(
            'output',
            null,
            InputOption::VALUE_REQUIRED,
            'Destination JSON file path.',
            $this->defaultOutputPath,
        );
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $io = new SymfonyStyle($input, $output);

        $outputPath = trim((string) $input->getOption('output'));
        if ($outputPath === '') {
            $io->error('Output path must not be empty.');

            return Command::FAILURE;
        }

        $formats = $this->resolveFormats($io);
        if ($formats === []) {
            $io->error('No formats found to refresh.');

            return Command::FAILURE;
        }

        $result = [];
        $successCount = 0;
        foreach ($formats as $format) {
            [$info, $isLiveData] = $this->resolveFormatInfo($format);
            $result[$format] = $info;
            if ($isLiveData) {
                ++$successCount;
            }
        }

        ksort($result);
        $payload = [
            'generatedAt' => gmdate('c'),
            'source' => 'manual-refresh-command',
            'formats' => $result,
        ];

        $encoded = json_encode($payload, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR);
        if (!\is_string($encoded)) {
            $io->error('Could not encode refreshed format info.');

            return Command::FAILURE;
        }

        $directory = \dirname($outputPath);
        if (!is_dir($directory)) {
            if (!@mkdir($directory, 0775, true) && !is_dir($directory)) {
                $io->error(sprintf('Could not create output directory: %s', $directory));

                return Command::FAILURE;
            }
        }

        if (@file_put_contents($outputPath, $encoded."\n") === false) {
            $io->error(sprintf('Could not write output file: %s', $outputPath));

            return Command::FAILURE;
        }

        $io->success(sprintf(
            'Stored %d format descriptions (%d fetched, %d fallback) in %s',
            \count($formats),
            $successCount,
            \count($formats) - $successCount,
            $outputPath,
        ));

        return Command::SUCCESS;
    }

    /**
     * @return list<string>
     */
    private function resolveFormats(SymfonyStyle $io): array
    {
        $formats = [];

        try {
            $formatsBySource = $this->conversionCatalogService->getFormats();
            foreach ($formatsBySource as $source => $targets) {
                $normalizedSource = strtolower(trim($source));
                if ($normalizedSource !== '') {
                    $formats[] = $normalizedSource;
                }

                foreach ($targets as $target) {
                    $normalizedTarget = strtolower(trim($target));
                    if ($normalizedTarget !== '') {
                        $formats[] = $normalizedTarget;
                    }
                }
            }
        } catch (\Throwable $exception) {
            $io->warning(sprintf(
                'Could not load formats from converter API (%s). Falling back to local known formats.',
                $exception->getMessage(),
            ));
        }

        $formats = array_values(array_unique([
            ...$formats,
            ...$this->formatInfoCatalog->knownFormats(),
        ]));
        sort($formats);

        return $formats;
    }

    /**
     * @return array{
     *     0: array{format: string, label: string, title: string, summary: string, url: string},
     *     1: bool
     * }
     */
    private function resolveFormatInfo(string $format): array
    {
        $fallbackInfo = $this->formatInfoCatalog->fallbackInfo($format);
        $pageTitle = $this->formatInfoCatalog->pageTitleFor($format);
        $summaryUrl = $this->formatInfoCatalog->buildWikipediaSummaryUrl($pageTitle);

        try {
            $response = $this->httpClient->request('GET', $summaryUrl, [
                'timeout' => 5,
                'headers' => [
                    'Accept' => 'application/json',
                ],
            ]);
            $statusCode = $response->getStatusCode();
            if ($statusCode < 200 || $statusCode >= 300) {
                return [$fallbackInfo, false];
            }

            $payload = $response->toArray(false);
        } catch (\Throwable) {
            return [$fallbackInfo, false];
        }

        if (!\is_array($payload)) {
            return [$fallbackInfo, false];
        }

        $title = $this->readString($payload, 'title') ?? $fallbackInfo['title'];
        $summary = $this->readString($payload, 'extract');
        if ($summary === null) {
            return [$fallbackInfo, false];
        }

        $url = $this->extractSummaryUrl($payload) ?? $fallbackInfo['url'];

        return [[
            'format' => $fallbackInfo['format'],
            'label' => $fallbackInfo['label'],
            'title' => $title,
            'summary' => $summary,
            'url' => $url,
        ], true];
    }

    /**
     * @param array<mixed> $payload
     */
    private function readString(array $payload, string $key): ?string
    {
        $value = $payload[$key] ?? null;
        if (!\is_string($value)) {
            return null;
        }

        $trimmed = trim($value);

        return $trimmed === '' ? null : $trimmed;
    }

    /**
     * @param array<mixed> $payload
     */
    private function extractSummaryUrl(array $payload): ?string
    {
        $desktop = $payload['content_urls']['desktop'] ?? null;
        if (!\is_array($desktop)) {
            return null;
        }

        $url = $desktop['page'] ?? null;
        if (!\is_string($url)) {
            return null;
        }

        $trimmed = trim($url);

        return $trimmed === '' ? null : $trimmed;
    }
}
