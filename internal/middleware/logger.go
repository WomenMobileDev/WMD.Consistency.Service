package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger is a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)

		// Log request details
		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Dur("latency", latency).
			Int("size", c.Writer.Size()).
			Msg("Request processed")
	}
}
