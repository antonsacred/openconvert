<?php

namespace App\Dto;

final class SitemapResult
{
    public function __construct(
        private readonly string $xml,
        private readonly int $urlCount,
    ) {
    }

    public function xml(): string
    {
        return $this->xml;
    }

    public function urlCount(): int
    {
        return $this->urlCount;
    }
}
