package server

import (
	"net/http"

	apidocs "goconverter/internal/apidocs"
	"goconverter/internal/converter"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

type ConversionsResponse struct {
	Output [][2]string `json:"output"`
}

// healthHandler godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

// listConversionsHandler godoc
// @Summary List supported conversions
// @Tags conversions
// @Produce json
// @Success 200 {object} ConversionsResponse
// @Router /conversions [get]
func listConversionsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, ConversionsResponse{Output: converter.PossibleConversions()})
}

// openAPISpecHandler godoc
// @Summary OpenAPI spec
// @Tags docs
// @Produce json
// @Success 200 {string} string "OpenAPI JSON"
// @Router /openapi.json [get]
func openAPISpecHandler(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(apidocs.SwaggerInfo.ReadDoc()))
}
