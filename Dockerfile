# syntax=docker/dockerfile:1.7

FROM composer:2 AS vendor
WORKDIR /app

COPY composer.json composer.lock symfony.lock ./
RUN composer install \
    --no-dev \
    --prefer-dist \
    --no-interaction \
    --no-progress \
    --no-scripts \
    --optimize-autoloader

FROM serversideup/php:8.5-fpm-nginx
WORKDIR /var/www/html

ENV APP_ENV=prod

COPY --chown=www-data:www-data . .
COPY --from=vendor --chown=www-data:www-data /app/vendor ./vendor

EXPOSE 8080
