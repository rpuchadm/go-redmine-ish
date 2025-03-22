package middleware

import (
	"fmt"
	"go-redmine-ish/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware de autenticación
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("AuthMiddleware")

		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		fmt.Printf("AuthMiddleware - authHeader %s\n", authHeader)

		// Verificar que el header tenga el formato correcto: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Extraer el token
		token := parts[1]

		fmt.Printf("AuthMiddleware - token %s\n", token)

		// Validar el token
		if token != cfg.AuthToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Si el token es válido, continuar con el siguiente handler
		c.Next()
	}
}
