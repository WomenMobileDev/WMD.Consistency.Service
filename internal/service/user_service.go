package service

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// UpdateProfileRequest represents the request for updating a user profile
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetProfile gets the user profile
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*models.UserResponse, error) {
	// Find the user by ID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Return the user profile
	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile updates the user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*models.UserResponse, error) {
	// Find the user by ID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update the user's name
	user.Name = req.Name

	// Save the changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Return the updated user profile
	response := user.ToResponse()
	return &response, nil
}
