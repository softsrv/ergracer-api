package api

import (
	"database/sql"

	"ergracer-api/internal/api/handlers"
	"ergracer-api/internal/config"
	"ergracer-api/internal/middleware"
	"ergracer-api/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	db     *sql.DB
	config *config.Config
}

func NewServer(db *sql.DB, config *config.Config) *Server {
	router := gin.New()

	// Set trusted proxies for security
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	server := &Server{
		router: router,
		db:     db,
		config: config,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	userService := services.NewUserService(s.db)
	sessionService := services.NewSessionService(s.db)
	friendshipService := services.NewFriendshipService(s.db)
	raceService := services.NewRaceService(s.db)

	authHandler := handlers.NewAuthHandler(userService, sessionService, s.config)
	friendsHandler := handlers.NewFriendsHandler(friendshipService, userService)
	racesHandler := handlers.NewRacesHandler(raceService)
	historyHandler := handlers.NewHistoryHandler(s.db)

	api := s.router.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/verify-email", authHandler.VerifyEmail)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthRequired(s.config.JWTSecret()))
	{
		protected.GET("/profile", authHandler.GetProfile)

		friends := protected.Group("/friends")
		{
			friends.POST("/invite", friendsHandler.InviteFriend)
			friends.POST("/accept/:friendId", friendsHandler.AcceptFriendship)
			friends.GET("/", friendsHandler.GetFriends)
			friends.GET("/invitations", friendsHandler.GetPendingInvitations)
		}

		races := protected.Group("/races")
		{
			races.POST("/", racesHandler.CreateRace)
			races.POST("/join", racesHandler.JoinRace)
			races.GET("/:uuid", racesHandler.GetRace)
			races.POST("/:raceId/ready", racesHandler.SetReady)
			races.POST("/:raceId/progress", racesHandler.UpdateProgress)
			races.POST("/:raceId/start", racesHandler.StartRace)
		}

		protected.GET("/history", historyHandler.GetUserRaceHistory)
	}

	s.router.HEAD("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
