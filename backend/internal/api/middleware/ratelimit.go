package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/pkg/database"
	"github.com/smartstocks/backend/pkg/utils"
)

// RateLimitMiddleware limita las peticiones por IP
func RateLimitMiddleware(redis *database.RedisClient, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		// Incrementar contador
		count, err := redis.Increment(ctx, key)
		if err != nil {
			// Si hay error en Redis, permitir la petición
			c.Next()
			return
		}

		// Si es el primer request, establecer expiración de 1 minuto
		if count == 1 {
			_ = redis.SetExpire(ctx, key, time.Minute)
		}

		// Verificar límite
		if count > int64(requestsPerMinute) {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
