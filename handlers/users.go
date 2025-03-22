package handlers

import (
	"go-redmine-ish/config"
	"go-redmine-ish/database"
	"go-redmine-ish/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsersHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Inicializar la base de datos
		db, err := database.InitDB(cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer db.Close()

		users, err := models.GetAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := models.CountUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := gin.H{
			"users": users,
			"count": count,
		}

		c.JSON(http.StatusOK, data)
	}
}
