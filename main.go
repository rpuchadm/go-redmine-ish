package main

import (
	"go-redmine-ish/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Crear un router Gin
	router := gin.Default()

	// Definir el endpoint /healthz
	router.GET("/healthz", handlers.HealthzHandler)

	// Definir el endpoint /init
	router.GET("/init", handlers.InitHandler)

	// Iniciar el servidor
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
