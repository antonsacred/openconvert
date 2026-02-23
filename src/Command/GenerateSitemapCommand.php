<?php

namespace App\Command;

use App\Service\SitemapService;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Style\SymfonyStyle;
use Symfony\Component\DependencyInjection\Attribute\Autowire;

#[AsCommand(
    name: 'app:sitemap:generate',
    description: 'Generates sitemap.xml for home/source/pair conversion pages.',
)]
final class GenerateSitemapCommand extends Command
{
    public function __construct(
        private readonly SitemapService $sitemapService,
        #[Autowire('%kernel.project_dir%/public/sitemap.xml')]
        private readonly string $defaultOutputPath,
    ) {
        parent::__construct();
    }

    protected function configure(): void
    {
        $this->addOption(
            'hostname',
            null,
            InputOption::VALUE_REQUIRED,
            'Hostname used to generate sitemap URLs (for example: openconvert.example.com).',
        );

        $this->addOption(
            'output',
            null,
            InputOption::VALUE_REQUIRED,
            'Destination sitemap file path.',
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

        try {
            $result = $this->sitemapService->generate((string) $input->getOption('hostname'));
        } catch (\InvalidArgumentException $exception) {
            $io->error($exception->getMessage());

            return Command::FAILURE;
        } catch (\Throwable $exception) {
            $io->error(sprintf(
                'Could not generate sitemap: %s',
                $exception->getMessage(),
            ));

            return Command::FAILURE;
        }

        $directory = \dirname($outputPath);
        if (!is_dir($directory)) {
            if (!@mkdir($directory, 0775, true) && !is_dir($directory)) {
                $io->error(sprintf('Could not create output directory: %s', $directory));

                return Command::FAILURE;
            }
        }

        if (@file_put_contents($outputPath, $result->xml()) === false) {
            $io->error(sprintf('Could not write output file: %s', $outputPath));

            return Command::FAILURE;
        }

        $io->success(sprintf(
            'Generated sitemap with %d URLs at %s',
            $result->urlCount(),
            $outputPath,
        ));

        return Command::SUCCESS;
    }
}
