package handlers

import (
	"net/http"

	"go-redmine-ish/config"
	"go-redmine-ish/database"

	"github.com/gin-gonic/gin"
)

func InitHandler(c *gin.Context) {
	// Cargar la configuración
	cfg := config.LoadConfig()

	// Inicializar la base de datos
	db, err := database.InitDB(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Base de datos inicializada correctamente"})
}
