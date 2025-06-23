package handlers

import (
	"net/http"
	"strconv"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CheckInHandler handles check-in-related requests
type CheckInHandler struct {
	checkInService *service.CheckInService
}

// NewCheckInHandler creates a new check-in handler
func NewCheckInHandler(checkInService *service.CheckInService) *CheckInHandler {
	return &CheckInHandler{
		checkInService: checkInService,
	}
}

// CheckIn handles checking in for a habit
func (h *CheckInHandler) CheckIn() gin.HandlerFunc {
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
		var req service.CheckInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to check in
		checkIn, err := h.checkInService.CheckIn(c.Request.Context(), userID, uint(habitID), req)
		if err != nil {
			if appErr, ok := err.(*models.AppError); ok {
				middleware.RespondWithError(c, http.StatusBadRequest, appErr.Code, appErr.Message, appErr.Details)
			} else {
				log.Error().Err(err).Msg("Failed to check in")
				middleware.RespondWithInternalError(c, err.Error())
			}
			return
		}

		middleware.RespondWithCreated(c, checkIn)
	}
}

// ListCheckIns handles listing all check-ins for a habit
func (h *CheckInHandler) ListCheckIns() gin.HandlerFunc {
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

		// Call the service to list check-ins
		checkIns, err := h.checkInService.ListCheckIns(c.Request.Context(), userID, uint(habitID))
		if err != nil {
			if err.Error() == "habit not found" {
				middleware.RespondWithNotFound(c, "Habit")
				return
			} else if err.Error() == "forbidden" {
				middleware.RespondWithForbidden(c)
				return
			}
			log.Error().Err(err).Msg("Failed to list check-ins")
			middleware.RespondWithInternalError(c, "Failed to list check-ins")
			return
		}

		middleware.RespondWithOK(c, checkIns)
	}
}
