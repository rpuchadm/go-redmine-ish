package handlers

import (
	"net/http"

	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"

	"github.com/gin-gonic/gin"
)

func InitHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		err = models.DropProjectsTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = models.CreateProjectsTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = models.TestProjectsTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Base de datos inicializada correctamente"})
	}
}
