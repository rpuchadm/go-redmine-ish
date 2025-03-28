// @title Go Redmine-ish API
// @version 1.0
// @description API para gestionar proyectos, usuarios, roles, categorías y problemas al estilo de Redmine.

package main

import (
	"go-redmine-ish/config"
	"go-redmine-ish/docs" // docs is generated by Swag CLI, you have to import it.
	"go-redmine-ish/handlers"
	"go-redmine-ish/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	// Swagger
	docs.SwaggerInfo.Title = "Go Redmine-ish API"
	docs.SwaggerInfo.Description = "API para gestionar proyectos, usuarios, roles, categorías y problemas al estilo de Redmine."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "https://issues.mydomain.com/"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https"}

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

	authGroup.GET("/category/:id", handlers.GetCategoryHandler(cfg))
	authGroup.POST("/category", handlers.CreateCategoryHandler(cfg))
	authGroup.PUT("/category/:id", handlers.UpdateCategoryHandler(cfg))
	authGroup.DELETE("/category/:id", handlers.DeleteCategoryHandler(cfg))

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

	authGroup.GET("/settings", handlers.GetSettingsHandler(cfg))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
