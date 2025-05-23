package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HealthCheck(c *gin.Context) {
	middleware.RespondWithOK(c, gin.H{
		"status": "up and running",
	})
}

// DBHealthCheck checks if the database connection is healthy using a database.Database instance
func DBHealthCheck(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			middleware.RespondWithError(c, http.StatusServiceUnavailable, "Database connection not initialized", nil)
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Ping the database
		err := db.Ping(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Database health check failed")
			middleware.RespondWithError(c, http.StatusServiceUnavailable, "Database connection failed", err.Error())
			return
		}

		// Additional GORM-specific health check
		var result int64
		db.DB.Raw("SELECT 1").Scan(&result)

		middleware.RespondWithOK(c, gin.H{
			"status": "success",
			"orm":    "GORM connection verified",
		})
	}
}
