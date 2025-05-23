package server

import (
	"context"
	"net/http"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router *gin.Engine
	http   *http.Server
	config *config.Config
	db     *database.Database
}

func NewServer(cfg *config.Config) *Server {
	// Initialize database connection
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to database, continuing without database connection")
	}
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Env == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	srv := &Server{
		router: router,
		config: cfg,
		db:     db,
		http: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
	}

	srv.setupRoutes()

	return srv
}

func (s *Server) setupRoutes() {
	// Root level routes
	s.router.GET("/health", handlers.HealthCheck)
	s.router.GET("/test-reload", handlers.TestReload)
	s.router.GET("/db-health", handlers.DBHealthCheck(s.db))

	// API v1 group
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheck)
		v1.GET("/test-reload", handlers.TestReload)
		v1.GET("/db-health", handlers.DBHealthCheck(s.db))
	}

	log.Info().Msg("Routes configured successfully")
}

func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Close database connection if it exists
	if s.db != nil {
		s.db.Close()
	}
	return s.http.Shutdown(ctx)
}
