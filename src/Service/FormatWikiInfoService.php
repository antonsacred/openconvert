<?php

namespace App\Service;

use Symfony\Component\DependencyInjection\Attribute\Autowire;

final class FormatWikiInfoService
{
    /**
     * @var array<string, array{format: string, label: string, title: string, summary: string, url: string}>|null
     */
    private ?array $loadedFormatInfo = null;

    public function __construct(
        private readonly FormatInfoCatalog $formatInfoCatalog,
        #[Autowire('%kernel.project_dir%/config/format_info_data.json')]
        private readonly string $formatInfoDataPath,
    ) {
    }

    /**
     * @return array{
     *     format: string,
     *     label: string,
     *     title: string,
     *     summary: string,
     *     url: string
     * }
     */
    public function getFormatInfo(string $format): array
    {
        $normalizedFormat = strtolower(trim($format));
        if ($normalizedFormat === '') {
            throw new \InvalidArgumentException('Format must not be empty.');
        }

        $storedInfo = $this->loadStoredInfo();
        if (isset($storedInfo[$normalizedFormat])) {
            return $storedInfo[$normalizedFormat];
        }

        return $this->formatInfoCatalog->fallbackInfo($normalizedFormat);
    }

    /**
     * @return array<string, array{format: string, label: string, title: string, summary: string, url: string}>
     */
    private function loadStoredInfo(): array
    {
        if ($this->loadedFormatInfo !== null) {
            return $this->loadedFormatInfo;
        }

        if (!is_file($this->formatInfoDataPath)) {
            $this->loadedFormatInfo = [];

            return $this->loadedFormatInfo;
        }

        $raw = @file_get_contents($this->formatInfoDataPath);
        if (!\is_string($raw) || $raw === '') {
            $this->loadedFormatInfo = [];

            return $this->loadedFormatInfo;
        }

        try {
            $decoded = json_decode($raw, true, 512, JSON_THROW_ON_ERROR);
        } catch (\JsonException) {
            $this->loadedFormatInfo = [];

            return $this->loadedFormatInfo;
        }

        if (!\is_array($decoded) || !isset($decoded['formats']) || !\is_array($decoded['formats'])) {
            $this->loadedFormatInfo = [];

            return $this->loadedFormatInfo;
        }

        $loaded = [];
        foreach ($decoded['formats'] as $format => $info) {
            if (!\is_string($format) || !\is_array($info)) {
                continue;
            }

            $normalizedFormat = strtolower(trim($format));
            $title = $this->readString($info, 'title');
            $summary = $this->readString($info, 'summary');
            $url = $this->readString($info, 'url');
            if ($normalizedFormat === '' || $title === null || $summary === null || $url === null) {
                continue;
            }

            $loaded[$normalizedFormat] = [
                'format' => $normalizedFormat,
                'label' => strtoupper($normalizedFormat),
                'title' => $title,
                'summary' => $summary,
                'url' => $url,
            ];
        }

        $this->loadedFormatInfo = $loaded;

        return $this->loadedFormatInfo;
    }

    /**
     * @param array<string, mixed> $payload
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
}
