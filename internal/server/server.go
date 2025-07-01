package server

import (
	"context"
	"net/http"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/router"
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
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to database, continuing without database connection")
	} else {
		// Run database migrations
		if err := db.RunMigrations(); err != nil {
			log.Error().Err(err).Msg("Failed to run database migrations")
		} else {
			log.Info().Msg("Database migrations completed successfully")
		}
	}
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Env == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Use the router from the router package
	r := router.SetupRouter(db)

	srv := &Server{
		router: r,
		config: cfg,
		db:     db,
		http: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: r,
		},
	}

	return srv
}

func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.db != nil {
		s.db.Close()
	}
	return s.http.Shutdown(ctx)
}
