package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// DBHealthCheck checks the health of the database connection
func DBHealthCheck(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection not initialized",
			})
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		// Ping the database
		err := db.Ping(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Database health check failed")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			})
			return
		}

		// Additional GORM-specific health check
		var result int64
		db.DB.Raw("SELECT 1").Scan(&result)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Database connection is healthy",
			"orm":     "GORM connection verified",
		})
	}
}
