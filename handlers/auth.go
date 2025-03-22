package handlers

import (
	"go-redmine-ish/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAuthHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
