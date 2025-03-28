package handlers

import (
	"go-redmine-ish/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary: GetAuthHandler
// @Description: Get authentication status
// @Tags: auth
// @Produce: json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth [get]
// @Security BearerAuth
func GetAuthHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
