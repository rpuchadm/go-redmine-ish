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
	router.GET("/project/:id", handlers.GetProjectHandler(cfg))
	router.POST("/project", handlers.CreateProjectHandler(cfg))
	router.PUT("/project/:id", handlers.UpdateProjectHandler(cfg))
	router.DELETE("/project/:id", handlers.DeleteProjectHandler(cfg))

	router.GET("/users", handlers.GetUsersHandler(cfg))
	router.GET("/user/:id", handlers.GetUserHandler(cfg))
	router.POST("/user", handlers.CreateUserHandler(cfg))
	router.PUT("/user/:id", handlers.UpdateUserHandler(cfg))
	router.DELETE("/user/:id", handlers.DeleteUserHandler(cfg))

	router.GET("/roles", handlers.GetRolesHandler(cfg))
	router.GET("/role/:id", handlers.GetRoleHandler(cfg))
	router.POST("/role", handlers.CreateRoleHandler(cfg))
	router.PUT("/role/:id", handlers.UpdateRoleHandler(cfg))
	router.DELETE("/role/:id", handlers.DeleteRoleHandler(cfg))

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
