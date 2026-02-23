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

FROM node:22-alpine AS node_deps
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

FROM serversideup/php:8.5-fpm-nginx AS build
WORKDIR /var/www/html

ENV APP_ENV=prod \
    APP_DEBUG=0

COPY --chown=www-data:www-data . .
COPY --from=vendor --chown=www-data:www-data /app/vendor ./vendor
COPY --from=node_deps --chown=www-data:www-data /app/node_modules ./node_modules

RUN php bin/console importmap:install --no-interaction
RUN php bin/console typescript:build
RUN php bin/console tailwind:build
RUN php bin/console asset-map:compile
RUN rm -rf node_modules

FROM serversideup/php:8.5-fpm-nginx
WORKDIR /var/www/html

ENV APP_ENV=prod \
    APP_DEBUG=0

COPY --from=build --chown=www-data:www-data /var/www/html /var/www/html

EXPOSE 8080
