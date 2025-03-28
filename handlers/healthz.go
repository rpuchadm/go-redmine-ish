package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary: HealthzHandler
// @Description: Health check endpoint
// @Tags: health
// @Produce: json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func HealthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
