package server

import (
	"context"
	"net/http"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router *gin.Engine
	http   *http.Server
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
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
		http: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
	}

	srv.setupRoutes()

	return srv
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", handlers.HealthCheck)
	s.router.GET("/test-reload", handlers.TestReload)

	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheck)
		v1.GET("/test-reload", handlers.TestReload)
	}

	log.Info().Msg("Routes configured successfully")
}

func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
