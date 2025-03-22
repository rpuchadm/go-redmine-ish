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
	router.GET("/users", handlers.GetUsersHandler(cfg))
	router.GET("/roles", handlers.GetRolesHandler(cfg))
	router.GET("/trackers", handlers.GetTrackersHandler(cfg))
	router.GET("/issues", handlers.GetIssuesHandler(cfg))
	router.GET("/issue/:id", handlers.GetIssueHandler(cfg))
	router.POST("/issue", handlers.CreateIssueHandler(cfg))
	router.PUT("/issue/:id", handlers.UpdateIssueHandler(cfg))
	router.DELETE("/issue/:id", handlers.DeleteIssueHandler(cfg))

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
