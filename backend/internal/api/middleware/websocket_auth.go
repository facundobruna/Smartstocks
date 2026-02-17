package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/pkg/jwt"
)

// WebSocketAuthMiddleware verifica el JWT token para conexiones WebSocket
// Lee el token del query parameter en lugar del header
func WebSocketAuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero intentar obtener del header (estándar)
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// Si no hay en header, buscar en query parameter
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		// Si aún no hay token, rechazar
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required. Use: ws://host/path?token=YOUR_TOKEN",
			})
			c.Abort()
			return
		}

		// Validar token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Guardar claims en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}
