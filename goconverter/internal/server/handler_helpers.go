package server

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
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
