# OpenConvert

OpenConvert is a Symfony frontend/API gateway paired with a Go conversion service (`goconverter`).

## Local Runtime

- Start: `make start`
- Stop: `make stop`
- Test: `make test`
- Go tests: `cd goconverter && go test ./...`

## Production Environment Variables

### Symfony app

- `CONVERTER_API`: Base URL of converter service (required).
- `APP_MAX_UPLOAD_BYTES`: Maximum uploaded file size accepted by `/api/convert` (default `10485760`).
- `APP_CONVERT_RATE_LIMIT`: Per-IP allowed `/api/convert` requests per minute (default `30`).

### Go converter

- `PORT`: HTTP port (default `8081`).
- `GO_CONVERTER_MAX_DECODED_FILE_SIZE_BYTES`: Max decoded file size (default `52428800`).
- `GO_CONVERTER_MAX_REQUEST_BODY_BYTES`: Max HTTP request body bytes (default derived from decoded limit + overhead).
- `GO_CONVERTER_MAX_CONCURRENT_CONVERSIONS`: Max in-flight conversions (default `4`).
- `GO_CONVERTER_READ_HEADER_TIMEOUT_SECONDS`: Read header timeout (default `5`).
- `GO_CONVERTER_READ_TIMEOUT_SECONDS`: Read timeout (default `30`).
- `GO_CONVERTER_WRITE_TIMEOUT_SECONDS`: Write timeout (default `60`).
- `GO_CONVERTER_IDLE_TIMEOUT_SECONDS`: Idle timeout (default `120`).

## Operational Notes

- `/api/convert` is rate limited per client IP.
- Requests and error payloads include `X-Request-Id` / `error.requestId` for cross-service tracing.
- `/health` verifies converter reachability through `CONVERTER_API/health`.
- Format info shown on converter pages is loaded from `config/format_info_data.json` (no runtime Wikipedia API calls).
- Refresh format info manually when needed:
  - `php bin/console app:format-info:refresh`
  - optional output path override: `php bin/console app:format-info:refresh --output=config/format_info_data.json`
- Generate sitemap manually when needed:
  - `php bin/console app:sitemap:generate --hostname=openconvert.example.com`
  - optional output path override: `php bin/console app:sitemap:generate --hostname=openconvert.example.com --output=public/sitemap.xml`
- Dynamic sitemap is available at `/sitemap.xml` and always emits `https://` URLs for the current request host.
