<?php

namespace App\Dto;

final class ConvertResult
{
    /**
     * @param array<string, mixed> $payload
     */
    public function __construct(
        private readonly int $statusCode,
        private readonly array $payload,
    ) {
    }

    public function statusCode(): int
    {
        return $this->statusCode;
    }

    /**
     * @return array<string, mixed>
     */
    public function payload(): array
    {
        return $this->payload;
    }
}
