package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetTrackersHandlerData struct {
	Trackers []models.Tracker `json:"trackers"`
	Count    int              `json:"count"`
}

// @Summary: GetTrackersHandler
// @Description: Get all trackers
// @Tags: trackers
// @Produce: json
// @Success 200 {object} GetTrackersHandlerData
// @Failure 500 {object} map[string]string
// @Router /trackers [get]
// @Security BearerAuth
func GetTrackersHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		trackers, err := models.GetAllTrackers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountTrackers(db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := GetTrackersHandlerData{
			Trackers: trackers,
			Count:    count,
		}

		c.JSON(http.StatusOK, data)
	}
}
