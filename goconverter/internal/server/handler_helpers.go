package server

import (
	"path/filepath"
	"sort"
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

var canonicalFormatAliases = map[string][]string{
	"heif": {"heic"},
	"jpeg": {"jpg"},
	"tiff": {"tif"},
}

var aliasToCanonicalFormat = buildAliasToCanonicalFormat()

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
	return canonicalFormat(format)
}

func canonicalFormat(format string) string {
	normalized := strings.ToLower(strings.TrimSpace(format))
	if canonical, ok := aliasToCanonicalFormat[normalized]; ok {
		return canonical
	}

	return normalized
}

func aliasesForFormat(format string) []string {
	canonical := canonicalFormat(format)
	aliases, ok := canonicalFormatAliases[canonical]
	if !ok {
		return []string{canonical}
	}

	output := []string{canonical}
	output = append(output, aliases...)
	return output
}

func expandConversionFormatsWithAliases(canonicalFormats map[string][]string) map[string][]string {
	targetSetsBySource := make(map[string]map[string]struct{}, len(canonicalFormats))

	for sourceCanonical, canonicalTargets := range canonicalFormats {
		sourceFormats := aliasesForFormat(sourceCanonical)

		for _, source := range sourceFormats {
			targetSet, ok := targetSetsBySource[source]
			if !ok {
				targetSet = map[string]struct{}{}
				targetSetsBySource[source] = targetSet
			}

			for _, targetCanonical := range canonicalTargets {
				for _, target := range aliasesForFormat(targetCanonical) {
					targetSet[target] = struct{}{}
				}
			}
		}
	}

	output := make(map[string][]string, len(targetSetsBySource))
	for source, targetSet := range targetSetsBySource {
		targets := make([]string, 0, len(targetSet))
		for target := range targetSet {
			targets = append(targets, target)
		}
		sort.Strings(targets)
		output[source] = targets
	}

	return output
}

func buildAliasToCanonicalFormat() map[string]string {
	output := map[string]string{}
	for canonical, aliases := range canonicalFormatAliases {
		output[canonical] = canonical
		for _, alias := range aliases {
			output[alias] = canonical
		}
	}
	return output
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
	switch canonicalFormat(format) {
	case "jpeg":
		return "image/jpeg"
	case "avif":
		return "image/avif"
	case "gif":
		return "image/gif"
	case "heif":
		return "image/heif"
	case "png":
		return "image/png"
	case "tiff":
		return "image/tiff"
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
