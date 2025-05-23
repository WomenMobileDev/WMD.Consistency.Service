package handlers

import (
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile handles getting the current user's profile
func (h *UserHandler) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Call the service to get the user profile
		user, err := h.userService.GetProfile(c.Request.Context(), userID)
		if err != nil {
			if err.Error() == "user not found" {
				middleware.RespondWithNotFound(c, "User")
				return
			}
			log.Error().Err(err).Msg("Failed to get user profile")
			middleware.RespondWithInternalError(c, "Failed to get user profile")
			return
		}

		middleware.RespondWithOK(c, user)
	}
}

// UpdateProfile handles updating the current user's profile
func (h *UserHandler) UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user ID from the context
		userID, err := middleware.GetUserID(c)
		if err != nil {
			middleware.RespondWithUnauthorized(c)
			return
		}

		// Parse the request body
		var req service.UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to update the user profile
		user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
		if err != nil {
			if err.Error() == "user not found" {
				middleware.RespondWithNotFound(c, "User")
				return
			}
			log.Error().Err(err).Msg("Failed to update user profile")
			middleware.RespondWithInternalError(c, "Failed to update user profile")
			return
		}

		middleware.RespondWithOK(c, user)
	}
}
