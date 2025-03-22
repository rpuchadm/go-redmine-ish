package main

import (
	"go-redmine-ish/config"
	"go-redmine-ish/handlers"
	"go-redmine-ish/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Cargar la configuración
	cfg := config.LoadConfig()

	// Crear un router Gin
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Origen permitido
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Tiempo de caché para las opciones preflight
	}))

	router.GET("/healthz", handlers.HealthzHandler)

	router.GET("/init", handlers.InitHandler(cfg))

	// Grupo de rutas con middleware de autenticación
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware(cfg))

	authGroup.GET("/auth", handlers.GetAuthHandler(cfg))

	authGroup.GET("/projects", handlers.GetProjectsHandler(cfg))
	authGroup.GET("/project/:id", handlers.GetProjectHandler(cfg))
	authGroup.POST("/project", handlers.CreateProjectHandler(cfg))
	authGroup.PUT("/project/:id", handlers.UpdateProjectHandler(cfg))
	authGroup.DELETE("/project/:id", handlers.DeleteProjectHandler(cfg))

	authGroup.GET("/users", handlers.GetUsersHandler(cfg))
	authGroup.GET("/user/:id", handlers.GetUserHandler(cfg))
	authGroup.POST("/user", handlers.CreateUserHandler(cfg))
	authGroup.PUT("/user/:id", handlers.UpdateUserHandler(cfg))
	authGroup.DELETE("/user/:id", handlers.DeleteUserHandler(cfg))

	authGroup.GET("/roles", handlers.GetRolesHandler(cfg))
	authGroup.GET("/role/:id", handlers.GetRoleHandler(cfg))
	authGroup.POST("/role", handlers.CreateRoleHandler(cfg))
	authGroup.PUT("/role/:id", handlers.UpdateRoleHandler(cfg))
	authGroup.DELETE("/role/:id", handlers.DeleteRoleHandler(cfg))

	authGroup.GET("/trackers", handlers.GetTrackersHandler(cfg))

	authGroup.GET("/issues", handlers.GetIssuesHandler(cfg))
	authGroup.GET("/issue/:id", handlers.GetIssueHandler(cfg))
	authGroup.POST("/issue", handlers.CreateIssueHandler(cfg))
	authGroup.PUT("/issue/:id", handlers.UpdateIssueHandler(cfg))
	authGroup.DELETE("/issue/:id", handlers.DeleteIssueHandler(cfg))

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
