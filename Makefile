start:
	docker compose up -d

stop:
	docker compose down

test:
	./bin/phpunit
