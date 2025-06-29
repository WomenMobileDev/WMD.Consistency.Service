package handlers

import (
	"net/http"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to register the user
		response, err := h.authService.Register(c.Request.Context(), req)
		if err != nil {
			if err.Error() == "email already registered" {
				middleware.RespondWithConflict(c, err.Error())
				return
			}
			log.Error().Err(err).Msg("Failed to register user")
			middleware.RespondWithInternalError(c, "Failed to register user")
			return
		}

		middleware.RespondWithCreated(c, response)
	}
}

// Login handles user login
func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to login the user
		response, err := h.authService.Login(c.Request.Context(), req)
		if err != nil {
			if err.Error() == "invalid email or password" {
				middleware.RespondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
				return
			}
			log.Error().Err(err).Msg("Failed to login user")
			middleware.RespondWithInternalError(c, "Failed to login user")
			return
		}

		middleware.RespondWithOK(c, response)
	}
}

// ForgotPassword handles forgot password requests
func (h *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.ForgotPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to handle forgot password
		err := h.authService.ForgotPassword(c.Request.Context(), req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to process forgot password request")
			middleware.RespondWithInternalError(c, "Failed to process forgot password request")
			return
		}

		middleware.RespondWithSuccess(c, http.StatusOK, "If your email is registered, you will receive a password reset link", nil)
	}
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.ResetPasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			middleware.RespondWithBadRequest(c, err.Error())
			return
		}

		// Call the service to handle reset password
		err := h.authService.ResetPassword(c.Request.Context(), req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to reset password")
			middleware.RespondWithInternalError(c, "Failed to reset password")
			return
		}

		middleware.RespondWithSuccess(c, http.StatusOK, "Password has been reset successfully", nil)
	}
}
