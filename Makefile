start:
	docker compose up -d

stop:
	docker compose down

test:
	./bin/phpunit

test-go:
	cd goconverter && go test ./...

test-all:
	./bin/phpunit
	cd goconverter && go test ./...
