package handlers

import (
	"strconv"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AchievementHandler handles achievement-related requests
type AchievementHandler struct {
	achievementService *service.AchievementService
}

// NewAchievementHandler creates a new achievement handler
func NewAchievementHandler(achievementService *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{
		achievementService: achievementService,
	}
}

// ListAchievements handles listing all achievements for the current user
func (h *AchievementHandler) ListAchievements() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Call the service to list achievements
		achievements, err := h.achievementService.ListAchievements(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list achievements")
			middleware.RespondWithInternalError(c, "Failed to list achievements")
			return
		}

		middleware.RespondWithOK(c, achievements)
	}
}

// GetAchievement handles getting a specific achievement
func (h *AchievementHandler) GetAchievement() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Get the achievement ID from the URL
		achievementID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			middleware.RespondWithBadRequest(c, "Invalid achievement ID")
			return
		}

		// Call the service to get the achievement
		achievement, err := h.achievementService.GetAchievement(c.Request.Context(), userID, uint(achievementID))
		if err != nil {
			if err.Error() == "achievement not found" {
				middleware.RespondWithNotFound(c, "Achievement")
				return
			}
			log.Error().Err(err).Msg("Failed to get achievement")
			middleware.RespondWithInternalError(c, "Failed to get achievement")
			return
		}

		middleware.RespondWithOK(c, achievement)
	}
}

// ListHabitAchievements handles listing all achievements for a habit
func (h *AchievementHandler) ListHabitAchievements() gin.HandlerFunc {
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

		// Call the service to list habit achievements
		achievements, err := h.achievementService.ListHabitAchievements(c.Request.Context(), userID, uint(habitID))
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to list habit achievements")
			middleware.RespondWithInternalError(c, "Failed to list habit achievements")
			return
		}

		middleware.RespondWithOK(c, achievements)
	}
}
