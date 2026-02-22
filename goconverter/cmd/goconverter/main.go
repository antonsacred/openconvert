//go:generate go run github.com/swaggo/swag/cmd/swag init -g main.go -d .,../../internal/server -o ../../internal/apidocs --parseInternal

// @title GoConverter API
// @version 1.0
// @description API for discovering available file conversions.
// @BasePath /

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"goconverter/internal/server"
)

func main() {
	port := readEnvString("PORT", "8081")

	maxDecodedFileSizeBytes := readEnvInt("GO_CONVERTER_MAX_DECODED_FILE_SIZE_BYTES", 50*1024*1024)
	defaultMaxRequestBodyBytes := base64.StdEncoding.EncodedLen(maxDecodedFileSizeBytes) + (2 * 1024 * 1024)
	maxRequestBodyBytes := readEnvInt64("GO_CONVERTER_MAX_REQUEST_BODY_BYTES", int64(defaultMaxRequestBodyBytes))
	maxConcurrentConversions := readEnvInt("GO_CONVERTER_MAX_CONCURRENT_CONVERSIONS", 4)
	server.ConfigureRuntimeLimits(maxDecodedFileSizeBytes, maxRequestBodyBytes, maxConcurrentConversions)

	router := server.NewRouter()

	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: time.Duration(readEnvInt("GO_CONVERTER_READ_HEADER_TIMEOUT_SECONDS", 5)) * time.Second,
		ReadTimeout:       time.Duration(readEnvInt("GO_CONVERTER_READ_TIMEOUT_SECONDS", 30)) * time.Second,
		WriteTimeout:      time.Duration(readEnvInt("GO_CONVERTER_WRITE_TIMEOUT_SECONDS", 60)) * time.Second,
		IdleTimeout:       time.Duration(readEnvInt("GO_CONVERTER_IDLE_TIMEOUT_SECONDS", 120)) * time.Second,
	}

	log.Printf("starting goconverter on %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func readEnvString(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}

func readEnvInt(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue <= 0 {
		log.Printf("invalid %s=%q, using fallback %d", name, value, fallback)
		return fallback
	}

	return parsedValue
}

func readEnvInt64(name string, fallback int64) int64 {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsedValue <= 0 {
		log.Printf("invalid %s=%q, using fallback %d", name, value, fallback)
		return fallback
	}

	return parsedValue
}
