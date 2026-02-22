<?php

namespace App\Dto;

final class HealthStatusResult
{
    /**
     * @param array{status: string, checks: array{converter_api: array<string, mixed>}} $payload
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
     * @return array{status: string, checks: array{converter_api: array<string, mixed>}}
     */
    public function payload(): array
    {
        return $this->payload;
    }
}
