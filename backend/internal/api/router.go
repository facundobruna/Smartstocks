package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartstocks/backend/internal/api/handlers"
	"github.com/smartstocks/backend/internal/api/middleware"
	"github.com/smartstocks/backend/internal/config"
	"github.com/smartstocks/backend/pkg/database"
	"github.com/smartstocks/backend/pkg/jwt"
)

type Router struct {
	engine           *gin.Engine
	authHandler      *handlers.AuthHandler
	userHandler      *handlers.UserHandler
	simulatorHandler *handlers.SimulatorHandler
	jwtManager       *jwt.JWTManager
	redis            *database.RedisClient
	config           *config.Config
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	simulatorHandler *handlers.SimulatorHandler,
	jwtManager *jwt.JWTManager,
	redis *database.RedisClient,
	cfg *config.Config,
) *Router {
	return &Router{
		engine:           gin.Default(),
		authHandler:      authHandler,
		userHandler:      userHandler,
		simulatorHandler: simulatorHandler,
		jwtManager:       jwtManager,
		redis:            redis,
		config:           cfg,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Middlewares globales
	r.engine.Use(middleware.CORSMiddleware(r.config.CORS.AllowedOrigins))
	r.engine.Use(middleware.RateLimitMiddleware(r.redis, 100))

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "smart-stocks-api",
		})
	})

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes (públicas)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/verify-email", r.authHandler.VerifyEmail)
			auth.POST("/logout", r.authHandler.Logout)
		}

		// Schools route (pública)
		v1.GET("/schools", r.userHandler.GetSchools)

		// User routes (protegidas)
		user := v1.Group("/user")
		user.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			user.GET("/profile", r.userHandler.GetProfile)
			user.PUT("/profile", r.userHandler.UpdateProfile)
			user.GET("/stats", r.userHandler.GetUserStats)
		}

		// Simulator routes (protegidas)
		simulator := v1.Group("/simulator")
		simulator.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			simulator.GET("/:difficulty", r.simulatorHandler.GetScenario)
			simulator.POST("/submit", r.simulatorHandler.SubmitDecision)
			simulator.GET("/history", r.simulatorHandler.GetHistory)
			simulator.GET("/cooldown/:difficulty", r.simulatorHandler.GetCooldownStatus)
			simulator.GET("/stats", r.simulatorHandler.GetStats)
		}
	}

	return r.engine
}

func (r *Router) Run(port string) error {
	return r.engine.Run(":" + port)
}
