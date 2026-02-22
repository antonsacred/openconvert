package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"goconverter/internal/converter"

	"github.com/gin-gonic/gin"
)

// convertHandler godoc
// @Summary Convert file
// @Tags conversions
// @Accept json
// @Produce json
// @Param request body ConvertRequest true "Conversion request"
// @Success 200 {object} ConvertResponse
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 415 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/convert [post]
func convertHandler(c *gin.Context) {
	var request ConvertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return
	}

	from := normalizeFormat(request.From)
	to := normalizeFormat(request.To)
	fileName := strings.TrimSpace(request.FileName)
	contentBase64 := strings.TrimSpace(request.ContentBase64)
	if from == "" || to == "" || fileName == "" || contentBase64 == "" {
		writeError(c, http.StatusBadRequest, "invalid_request", "from, to, fileName, and contentBase64 are required")
		return
	}

	converterImplementation, ok := converter.FindConverter(from, to)
	if !ok {
		writeError(c, http.StatusUnsupportedMediaType, "unsupported_conversion_pair", fmt.Sprintf("conversion from %s to %s is not supported", from, to))
		return
	}

	if len(contentBase64) > base64.StdEncoding.EncodedLen(maxDecodedFileSizeBytes) {
		writeError(c, http.StatusRequestEntityTooLarge, "payload_too_large", "decoded input file exceeds 50MB limit")
		return
	}

	inputBytes, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid_base64", "contentBase64 must be valid base64")
		return
	}
	if len(inputBytes) > maxDecodedFileSizeBytes {
		writeError(c, http.StatusRequestEntityTooLarge, "payload_too_large", "decoded input file exceeds 50MB limit")
		return
	}

	outputBytes, err := converterImplementation.Convert(inputBytes)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "conversion_failed", "failed to convert file")
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{
		From:          from,
		To:            to,
		FileName:      outputFileName(fileName, to),
		MimeType:      mimeTypeByFormat(to),
		ContentBase64: base64.StdEncoding.EncodeToString(outputBytes),
	})
}
