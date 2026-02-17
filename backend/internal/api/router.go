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
	engine             *gin.Engine
	authHandler        *handlers.AuthHandler
	userHandler        *handlers.UserHandler
	quizHandler        *handlers.QuizHandler
	forumHandler       *handlers.ForumHandler
	coursesHandler     *handlers.CoursesHandler
	simulatorHandler   *handlers.SimulatorHandler
	pvpHandler         *handlers.PvPHandler
	rankingsHandler    *handlers.RankingsHandler
	tokensHandler      *handlers.TokensHandler
	tournamentsHandler *handlers.TournamentsHandler
	jwtManager         *jwt.JWTManager
	redis              *database.RedisClient
	config             *config.Config
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	quizHandler *handlers.QuizHandler,
	forumHandler *handlers.ForumHandler,
	coursesHandler *handlers.CoursesHandler,
	simulatorHandler *handlers.SimulatorHandler,
	pvpHandler *handlers.PvPHandler,
	rankingsHandler *handlers.RankingsHandler,
	tokensHandler *handlers.TokensHandler,
	tournamentsHandler *handlers.TournamentsHandler,
	jwtManager *jwt.JWTManager,
	redis *database.RedisClient,
	cfg *config.Config,
) *Router {
	return &Router{
		engine:             gin.Default(),
		authHandler:        authHandler,
		userHandler:        userHandler,
		quizHandler:        quizHandler,
		forumHandler:       forumHandler,
		coursesHandler:     coursesHandler,
		simulatorHandler:   simulatorHandler,
		pvpHandler:         pvpHandler,
		rankingsHandler:    rankingsHandler,
		tokensHandler:      tokensHandler,
		tournamentsHandler: tournamentsHandler,
		jwtManager:         jwtManager,
		redis:              redis,
		config:             cfg,
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
			"version": "1.0.0",
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

		// Quiz routes (protegidas)
		quiz := v1.Group("/quiz")
		quiz.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			quiz.GET("/:difficulty", r.quizHandler.GetDailyQuiz)
			quiz.POST("/submit", r.quizHandler.SubmitQuiz)
			quiz.GET("/history", r.quizHandler.GetQuizHistory)
		}

		// Forum routes (parcialmente públicas)
		forum := v1.Group("/forum")
		{
			// Rutas públicas (lectura)
			forum.GET("/posts", r.forumHandler.GetPosts)
			forum.GET("/posts/:id", r.forumHandler.GetPostByID)

			// Rutas protegidas (escritura)
			forumAuth := forum.Group("")
			forumAuth.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				forumAuth.POST("/posts", r.forumHandler.CreatePost)
				forumAuth.PUT("/posts/:id", r.forumHandler.UpdatePost)
				forumAuth.DELETE("/posts/:id", r.forumHandler.DeletePost)
				forumAuth.POST("/replies", r.forumHandler.CreateReply)
				forumAuth.DELETE("/replies/:id", r.forumHandler.DeleteReply)
				forumAuth.POST("/react", r.forumHandler.ReactToPost)
				forumAuth.DELETE("/react", r.forumHandler.RemoveReaction)
			}
		}

		// Courses routes (protegidas)
		courses := v1.Group("/courses")
		courses.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			courses.GET("", r.coursesHandler.GetAllCourses)
			courses.GET("/:id", r.coursesHandler.GetCourseByID)
			courses.GET("/lessons/:id", r.coursesHandler.GetLessonByID)
			courses.POST("/lessons/:id/complete", r.coursesHandler.CompleteLesson)
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

		// PvP routes
		pvp := v1.Group("/pvp")
		{
			// WebSocket endpoint
			pvp.GET("/ws", middleware.WebSocketAuthMiddleware(r.jwtManager), r.pvpHandler.WebSocket)

			// REST endpoints
			pvpRest := pvp.Group("")
			pvpRest.Use(middleware.AuthMiddleware(r.jwtManager))
			{
				pvpRest.POST("/queue/join", r.pvpHandler.JoinQueue)
				pvpRest.POST("/queue/leave", r.pvpHandler.LeaveQueue)
				pvpRest.POST("/submit", r.pvpHandler.SubmitDecision)
				pvpRest.GET("/history", r.pvpHandler.GetHistory)
			}
		}

		// Rankings routes (protegidas)
		rankings := v1.Group("/rankings")
		rankings.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			rankings.GET("/global", r.rankingsHandler.GetGlobalLeaderboard)
			rankings.GET("/school/:school_id", r.rankingsHandler.GetSchoolLeaderboard)
			rankings.GET("/my-school", r.rankingsHandler.GetMySchoolLeaderboard)
			rankings.GET("/my-position", r.rankingsHandler.GetMyPosition)
			rankings.GET("/profile/:user_id", r.rankingsHandler.GetPublicProfile)
			rankings.GET("/achievements", r.rankingsHandler.GetMyAchievements)

			// Admin endpoints
			rankings.POST("/admin/update-cache", r.rankingsHandler.UpdateCache)
		}

		// Tokens routes (protegidas)
		tokens := v1.Group("/tokens")
		tokens.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			tokens.GET("/balance", r.tokensHandler.GetMyTokens)
			tokens.GET("/transactions", r.tokensHandler.GetTransactionHistory)
		}

		// Tournaments routes (protegidas)
		tournaments := v1.Group("/tournaments")
		tournaments.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			tournaments.GET("", r.tournamentsHandler.GetActiveTournaments)
			tournaments.GET("/my-tournaments", r.tournamentsHandler.GetMyTournaments)
			tournaments.GET("/:tournament_id", r.tournamentsHandler.GetTournamentDetails)
			tournaments.POST("/join", r.tournamentsHandler.JoinTournament)
			tournaments.GET("/:tournament_id/standings", r.tournamentsHandler.GetTournamentStandings)
			tournaments.GET("/:tournament_id/bracket", r.tournamentsHandler.GetTournamentBracket)
		}
	}

	return r.engine
}

func (r *Router) Run(port string) error {
	return r.engine.Run(":" + port)
}
