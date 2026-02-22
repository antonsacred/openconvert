APP_NAME := goconverter
CMD_PATH := ./cmd/$(APP_NAME)

.PHONY: help test generate build run fmt tidy

help:
	@echo "Available targets:"
	@echo "  make test      - Run all tests"
	@echo "  make generate  - Generate OpenAPI docs"
	@echo "  make build     - Build binary to ./bin/$(APP_NAME)"
	@echo "  make run       - Run the application locally"
	@echo "  make fmt       - Format Go code"
	@echo "  make tidy      - Tidy Go module dependencies"

test:
	go test ./...

generate:
	go generate ./cmd/$(APP_NAME)

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

fmt:
	go fmt ./...

tidy:
	go mod tidy
