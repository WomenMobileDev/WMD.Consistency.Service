package handlers

import (
	"net/http"
	"strconv"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HabitHandler handles habit-related requests
type HabitHandler struct {
	habitService *service.HabitService
}

// NewHabitHandler creates a new habit handler
func NewHabitHandler(habitService *service.HabitService) *HabitHandler {
	return &HabitHandler{
		habitService: habitService,
	}
}

// CreateHabit handles creating a new habit
func (h *HabitHandler) CreateHabit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Parse the request body
		var req service.CreateHabitRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to create the habit
		habit, err := h.habitService.CreateHabit(c.Request.Context(), userID, req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create habit")
			middleware.RespondWithInternalError(c, "Failed to create habit")
			return
		}

		middleware.RespondWithCreated(c, habit)
	}
}

// GetHabit handles getting a habit by ID
func (h *HabitHandler) GetHabit() gin.HandlerFunc {
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

		// Call the service to get the habit
		habit, err := h.habitService.GetHabit(c.Request.Context(), userID, uint(habitID))
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to get habit")
			middleware.RespondWithInternalError(c, "Failed to get habit")
			return
		}

		middleware.RespondWithOK(c, habit)
	}
}

func (h *HabitHandler) ListHabits() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		habits, err := h.habitService.ListHabits(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list habits")
			middleware.RespondWithInternalError(c, "Failed to list habits")
			return
		}

		middleware.RespondWithOK(c, habits)
	}
}

// UpdateHabit handles updating a habit
func (h *HabitHandler) UpdateHabit() gin.HandlerFunc {
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
		var req service.UpdateHabitRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to update the habit
		habit, err := h.habitService.UpdateHabit(c.Request.Context(), userID, uint(habitID), req)
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to update habit")
			middleware.RespondWithInternalError(c, "Failed to update habit")
			return
		}

		middleware.RespondWithOK(c, habit)
	}
}

// DeleteHabit handles deleting a habit
func (h *HabitHandler) DeleteHabit() gin.HandlerFunc {
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

		// Call the service to delete the habit
		err = h.habitService.DeleteHabit(c.Request.Context(), userID, uint(habitID))
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to delete habit")
			middleware.RespondWithInternalError(c, "Failed to delete habit")
			return
		}

		middleware.RespondWithSuccess(c, http.StatusOK, "Habit deleted successfully", nil)
	}
}
