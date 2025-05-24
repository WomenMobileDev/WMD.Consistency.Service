package handlers

import (
	"strconv"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// StreakHandler handles streak-related requests
type StreakHandler struct {
	streakService *service.StreakService
}

// NewStreakHandler creates a new streak handler
func NewStreakHandler(streakService *service.StreakService) *StreakHandler {
	return &StreakHandler{
		streakService: streakService,
	}
}

// CreateStreak handles creating a new streak for a habit
func (h *StreakHandler) CreateStreak() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Get the habit ID from the URL
		habitID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			middleware.RespondWithBadRequest(c, "Invalid habit ID")
			return
		}

		// Parse the request body
		var req service.CreateStreakRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to create the streak
		streak, err := h.streakService.CreateStreak(c.Request.Context(), userID, uint(habitID), req)
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "an active streak already exists for this habit" {
				middleware.RespondWithConflict(c, err.Error())
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to create streak")
			middleware.RespondWithInternalError(c, "Failed to create streak")
			return
		}

		middleware.RespondWithCreated(c, streak)
	}
}

// GetCurrentStreak handles getting the current active streak for a habit
func (h *StreakHandler) GetCurrentStreak() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Get the habit ID from the URL
		habitID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			middleware.RespondWithBadRequest(c, "Invalid habit ID")
			return
		}

		// Call the service to get the current streak
		streak, err := h.streakService.GetCurrentStreak(c.Request.Context(), userID, uint(habitID))
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "no active streak found" {
				middleware.RespondWithNotFound(c, "Active streak")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to get current streak")
			middleware.RespondWithInternalError(c, "Failed to get current streak")
			return
		}

		middleware.RespondWithOK(c, streak)
	}
}

// ListStreaks handles listing all streaks for a habit
func (h *StreakHandler) ListStreaks() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Get the habit ID from the URL
		habitID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			middleware.RespondWithBadRequest(c, "Invalid habit ID")
			return
		}

		// Call the service to list the streaks
		streaks, err := h.streakService.ListStreaks(c.Request.Context(), userID, uint(habitID))
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to list streaks")
			middleware.RespondWithInternalError(c, "Failed to list streaks")
			return
		}

		middleware.RespondWithOK(c, streaks)
	}
}
