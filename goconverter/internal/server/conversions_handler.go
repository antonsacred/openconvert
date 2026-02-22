package server

import (
	"net/http"

	"goconverter/internal/converter"

	"github.com/gin-gonic/gin"
)

// listConversionsHandler godoc
// @Summary List supported conversions
// @Tags conversions
// @Produce json
// @Success 200 {object} ConversionsResponse
// @Router /v1/conversions [get]
func listConversionsHandler(c *gin.Context) {
	canonicalFormats := converter.ConversionTargetsBySource()

	c.JSON(http.StatusOK, ConversionsResponse{
		Formats: expandConversionFormatsWithAliases(canonicalFormats),
	})
}
