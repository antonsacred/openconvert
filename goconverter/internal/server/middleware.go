package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader("X-Request-Id"))
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set(requestIDContextKey, requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

func requestBodyLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodyBytes)
		c.Next()
	}
}

func requestLoggingMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		from, to := conversionFormatsFromContext(c)
		errorCode := c.GetString(errorCodeContextKey)

		payload := map[string]any{
			"level":       "info",
			"method":      c.Request.Method,
			"path":        c.FullPath(),
			"status_code": c.Writer.Status(),
			"duration_ms": time.Since(startedAt).Milliseconds(),
			"request_id":  requestIDFromContext(c),
			"client_ip":   c.ClientIP(),
			"from":        from,
			"to":          to,
		}
		if errorCode != "" {
			payload["error_code"] = errorCode
		}

		encodedPayload, err := json.Marshal(payload)
		if err != nil {
			logger.Printf(`{"level":"error","message":"failed to encode request log payload","status_code":%d}`, c.Writer.Status())
			return
		}

		logger.Print(string(encodedPayload))
	}
}

func generateRequestID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes[:])
}
