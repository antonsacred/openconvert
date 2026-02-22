package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthHandler godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
