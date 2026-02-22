package server

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(requestLoggingMiddleware(log.Default()))

	router.GET("/health", healthHandler)
	router.GET("/openapi.json", openAPISpecHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	v1.GET("/conversions", listConversionsHandler)
	v1.POST("/convert", requestBodyLimitMiddleware(), convertHandler)

	return router
}
