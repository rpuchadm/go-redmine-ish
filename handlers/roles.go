package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRolesHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		roles, err := models.GetAllRoles(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountRoles(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := gin.H{
			"roles": roles,
			"count": count,
		}

		c.JSON(http.StatusOK, data)
	}
}
