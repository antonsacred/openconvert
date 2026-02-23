<?php

namespace App\Tests\Support;

use App\Service\ConverterApiClient;
use Symfony\Contracts\HttpClient\HttpClientInterface;

trait ConverterApiClientOverrideTrait
{
    protected function overrideConverterApiClient(?string $converterApi, ?HttpClientInterface $httpClient = null): ConverterApiClient
    {
        $container = static::getContainer();

        if ($httpClient !== null) {
            $container->set(HttpClientInterface::class, $httpClient);
        } else {
            $httpClient = $container->get(HttpClientInterface::class);
        }

        $converterApiClient = new ConverterApiClient($httpClient, $converterApi);
        $container->set(ConverterApiClient::class, $converterApiClient);

        return $converterApiClient;
    }
}
