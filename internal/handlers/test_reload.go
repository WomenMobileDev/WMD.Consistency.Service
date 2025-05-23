package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestReload is a handler to test live reloading
func TestReload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Live reloading is working!",
		"version": "1.0.0",
	})
}
