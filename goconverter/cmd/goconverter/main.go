//go:generate go run github.com/swaggo/swag/cmd/swag init -g main.go -d .,../../internal/server -o ../../internal/apidocs --parseInternal

// @title GoConverter API
// @version 1.0
// @description API for discovering available file conversions.
// @BasePath /

package main

import (
	"log"
	"os"

	"goconverter/internal/server"
)

func main() {
	router := server.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
