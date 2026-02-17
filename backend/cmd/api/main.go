package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/smartstocks/backend/internal/api"
	"github.com/smartstocks/backend/internal/api/handlers"
	"github.com/smartstocks/backend/internal/config"
	"github.com/smartstocks/backend/internal/repository"
	"github.com/smartstocks/backend/internal/services"
	"github.com/smartstocks/backend/internal/websocket"
	"github.com/smartstocks/backend/pkg/database"
	"github.com/smartstocks/backend/pkg/jwt"
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

	// Inicializar WebSocket Manager
	wsManager := websocket.NewManager()
	go wsManager.Run()

	// Inicializar repositorios
	userRepo := repository.NewUserRepository(mysqlDB.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(mysqlDB.DB)
	schoolRepo := repository.NewSchoolRepository(mysqlDB.DB)
	simulatorRepo := repository.NewSimulatorRepository(mysqlDB.DB)
	pvpRepo := repository.NewPvPRepository(mysqlDB.DB)
	rankingsRepo := repository.NewRankingsRepository(mysqlDB.DB)
	tokensRepo := repository.NewTokensRepository(mysqlDB.DB)
	tournamentsRepo := repository.NewTournamentsRepository(mysqlDB.DB)

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

	pvpService := services.NewPvPService(
		pvpRepo,
		simulatorRepo,
		userRepo,
		simulatorAIService,
	)

	rankingsService := services.NewRankingsService(
		rankingsRepo,
		userRepo,
	)

	tokensService := services.NewTokensService(
		tokensRepo,
	)

	tournamentsService := services.NewTournamentsService(
		tournamentsRepo,
		userRepo,
	)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo, schoolRepo)
	simulatorHandler := handlers.NewSimulatorHandler(simulatorService)
	pvpHandler := handlers.NewPvPHandler(pvpService, wsManager)
	rankingsHandler := handlers.NewRankingsHandler(rankingsService)
	tokensHandler := handlers.NewTokensHandler(tokensService)
	tournamentsHandler := handlers.NewTournamentsHandler(tournamentsService)

	// Configurar router
	router := api.NewRouter(
		authHandler,
		userHandler,
		simulatorHandler,
		pvpHandler,
		rankingsHandler,
		tokensHandler,
		tournamentsHandler,
		jwtManager,
		redisClient,
		cfg,
	)

	engine := router.Setup()

	// Manejo de shutdown graceful
	go func() {
		log.Println("========================================")
		log.Printf("🚀 Smart Stocks API v1.0.0")
		log.Println("========================================")
		log.Printf("📡 Server starting on port %s", cfg.Server.Port)
		log.Printf("📊 Environment: %s", cfg.Server.GinMode)
		log.Printf("✅ MySQL connected to %s", cfg.Database.Host)
		log.Printf("✅ Redis connected to %s", cfg.Redis.Host)
		log.Printf("🔌 WebSocket manager running")
		log.Printf("🏆 Rankings system enabled")
		log.Printf("💰 Tokens system enabled")
		log.Printf("🎮 Tournaments system enabled")
		if cfg.OpenAI.APIKey != "" {
			log.Printf("✅ OpenAI configured (model: %s)", cfg.OpenAI.Model)
		} else {
			log.Printf("⚠️  OpenAI not configured - using fallback scenarios")
		}
		log.Println("========================================")
		log.Printf("📖 API Documentation: http://localhost:%s/health", cfg.Server.Port)
		log.Println("========================================")

		if err := engine.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("========================================")
	log.Println("🛑 Shutting down server...")
	log.Println("✅ Server stopped gracefully")
	log.Println("========================================")
}
