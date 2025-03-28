package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetSettingsHandlerData struct {
	Trackers []models.Tracker `json:"trackers"`
}

// @Summary: GetSettingsHandler
// @Description: Get settings
// @Tags: settings
// @Produce: json
// @Success 200 {object} GetSettingsHandlerData
// @Failure 500 {object} map[string]string
// @Router /settings [get]
// @Security BearerAuth
func GetSettingsHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		data := GetSettingsHandlerData{}

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(trackers) > 0 {
			data.Trackers = trackers
		}

		c.JSON(http.StatusOK, data)

	}
}
