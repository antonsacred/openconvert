package server

import (
	"net/http"

	apidocs "goconverter/internal/apidocs"

	"github.com/gin-gonic/gin"
)

// openAPISpecHandler godoc
// @Summary OpenAPI spec
// @Tags docs
// @Produce json
// @Success 200 {string} string "OpenAPI JSON"
// @Router /openapi.json [get]
func openAPISpecHandler(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(apidocs.SwaggerInfo.ReadDoc()))
}
