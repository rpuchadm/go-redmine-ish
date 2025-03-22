package main

import (
	"go-redmine-ish/config"
	"go-redmine-ish/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	// Cargar la configuración
	cfg := config.LoadConfig()

	// Crear un router Gin
	router := gin.Default()

	router.GET("/healthz", handlers.HealthzHandler)
	router.GET("/init", handlers.InitHandler(cfg))
	router.GET("/projects", handlers.GetProjectsHandler(cfg))

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
