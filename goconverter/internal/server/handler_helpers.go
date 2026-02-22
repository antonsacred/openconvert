package server

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	requestIDContextKey  = "request_id"
	errorCodeContextKey  = "error_code"
	fromFormatContextKey = "from"
	toFormatContextKey   = "to"
)

var conversionConcurrencyMu sync.Mutex
var currentConcurrentConversions int

func writeError(c *gin.Context, statusCode int, code string, message string) {
	c.Set(errorCodeContextKey, code)
	c.JSON(statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestIDFromContext(c),
		},
	})
}

func normalizeFormat(format string) string {
	return strings.ToLower(strings.TrimSpace(format))
}

func outputFileName(inputFileName string, targetFormat string) string {
	trimmed := strings.TrimSpace(inputFileName)
	if trimmed == "" {
		return "converted." + targetFormat
	}

	ext := filepath.Ext(trimmed)
	base := strings.TrimSuffix(trimmed, ext)
	if base == "" {
		base = "converted"
	}

	return base + "." + targetFormat
}

func mimeTypeByFormat(format string) string {
	switch format {
	case "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func requestIDFromContext(c *gin.Context) string {
	requestID, ok := c.Get(requestIDContextKey)
	if !ok {
		return ""
	}

	requestIDString, ok := requestID.(string)
	if !ok {
		return ""
	}

	return requestIDString
}

func markConversionFormats(c *gin.Context, from string, to string) {
	c.Set(fromFormatContextKey, from)
	c.Set(toFormatContextKey, to)
}

func conversionFormatsFromContext(c *gin.Context) (string, string) {
	from := ""
	if fromAny, ok := c.Get(fromFormatContextKey); ok {
		if fromString, valid := fromAny.(string); valid {
			from = fromString
		}
	}

	to := ""
	if toAny, ok := c.Get(toFormatContextKey); ok {
		if toString, valid := toAny.(string); valid {
			to = toString
		}
	}

	return from, to
}

func tryAcquireConversionSlot() bool {
	conversionConcurrencyMu.Lock()
	defer conversionConcurrencyMu.Unlock()

	if maxConcurrentConversions <= 0 || currentConcurrentConversions >= maxConcurrentConversions {
		return false
	}

	currentConcurrentConversions++

	return true
}

func releaseConversionSlot() {
	conversionConcurrencyMu.Lock()
	defer conversionConcurrencyMu.Unlock()

	if currentConcurrentConversions > 0 {
		currentConcurrentConversions--
	}
}
