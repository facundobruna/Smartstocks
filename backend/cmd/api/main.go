package main

import (
	"github.com/smartstocks/backend/internal/api"
	"github.com/smartstocks/backend/internal/api/handlers"
	"github.com/smartstocks/backend/internal/config"
	"github.com/smartstocks/backend/internal/repository"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/pkg/database"
	"github.com/smartstocks/backend/pkg/jwt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Conectar a MySQL
	mysqlDB, err := database.NewMySQL(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	// Conectar a Redis
	redisClient, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Inicializar JWT Manager
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Inicializar repositorios
	userRepo := repository.NewUserRepository(mysqlDB.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(mysqlDB.DB)
	schoolRepo := repository.NewSchoolRepository(mysqlDB.DB)
	simulatorRepo := repository.NewSimulatorRepository(mysqlDB.DB)

	// Inicializar servicios de IA
	simulatorAIService := services.NewSimulatorAIService(
		cfg.OpenAI.APIKey,
		cfg.OpenAI.APIURL,
		cfg.OpenAI.Model,
	)

	// Inicializar servicios
	authService := services.NewAuthService(
		userRepo,
		refreshTokenRepo,
		schoolRepo,
		jwtManager,
		cfg.JWT.RefreshTokenExpirationDays,
	)

	simulatorService := services.NewSimulatorService(
		simulatorRepo,
		userRepo,
		simulatorAIService,
	)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo, schoolRepo)
	simulatorHandler := handlers.NewSimulatorHandler(simulatorService)

	// Configurar router
	router := api.NewRouter(
		authHandler,
		userHandler,
		simulatorHandler,
		jwtManager,
		redisClient,
		cfg,
	)

	engine := router.Setup()

	// Manejo de shutdown graceful
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
		log.Printf("📊 Environment: %s", cfg.Server.GinMode)
		log.Printf("✅ MySQL connected to %s", cfg.Database.Host)
		log.Printf("✅ Redis connected to %s", cfg.Redis.Host)
		if cfg.OpenAI.APIKey != "" {
			log.Printf("✅ OpenAI configured (model: %s)", cfg.OpenAI.Model)
		} else {
			log.Printf("⚠️  OpenAI not configured - using fallback scenarios")
		}

		if err := engine.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	log.Println("✅ Server stopped gracefully")
}
