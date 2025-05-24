package service

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/middleware"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
	"github.com/rs/zerolog/log"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

// RegisterRequest represents the request for user registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ForgotPasswordRequest represents the request for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request for reset password
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// AuthResponse represents the response for authentication
type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Register handles user registration
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if the email is already registered
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Create a new user
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	// Set the password
	if err := user.SetPassword(req.Password); err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New("failed to hash password")
	}

	// Save the user to the database
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	// Generate a JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Return the token and user
	return &AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// Login handles user login
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Find the user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Check the password
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate a JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Return the token and user
	return &AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// ForgotPassword handles forgot password requests
func (s *AuthService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Find the user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		// Don't reveal that the email doesn't exist
		return nil
	}

	// In a real application, generate a reset token and send an email
	// For this example, we'll just log it
	log.Info().Str("email", user.Email).Msg("Password reset requested")

	return nil
}

// ResetPassword handles password reset
func (s *AuthService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// In a real application, verify the token and find the user
	// For this example, we'll just log it
	log.Info().Str("token", req.Token).Msg("Password reset token received")

	// Mock response for demonstration
	return nil
}
