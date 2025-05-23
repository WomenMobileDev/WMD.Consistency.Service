package service

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
)

// HabitService handles habit-related business logic
type HabitService struct {
	habitRepo repository.HabitRepository
	streakRepo repository.StreakRepository
}

// NewHabitService creates a new habit service
func NewHabitService(habitRepo repository.HabitRepository, streakRepo repository.StreakRepository) *HabitService {
	return &HabitService{
		habitRepo: habitRepo,
		streakRepo: streakRepo,
	}
}

// CreateHabitRequest represents the request for creating a habit
type CreateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" binding:"omitempty,len=7"` // Hex color code
	Icon        string `json:"icon"`
}

// UpdateHabitRequest represents the request for updating a habit
type UpdateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" binding:"omitempty,len=7"` // Hex color code
	Icon        string `json:"icon"`
	IsActive    *bool  `json:"is_active"`
}

// ListHabits lists all habits for a user
func (s *HabitService) ListHabits(ctx context.Context, userID uint) ([]models.HabitResponse, error) {
	// Find all habits for the user
	habits, err := s.habitRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert habits to responses
	var responses []models.HabitResponse
	for _, habit := range habits {
		// Get the current streak for this habit
		currentStreak, _ := s.streakRepo.FindActiveByHabitID(ctx, habit.ID)

		habitResponse := habit.ToResponse()
		if currentStreak != nil {
			streakResponse := currentStreak.ToResponse()
			habitResponse.CurrentStreak = &streakResponse
		}

		responses = append(responses, habitResponse)
	}

	return responses, nil
}

// CreateHabit creates a new habit
func (s *HabitService) CreateHabit(ctx context.Context, userID uint, req CreateHabitRequest) (*models.HabitResponse, error) {
	// Create a new habit
	habit := models.Habit{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    true,
	}

	// Save the habit to the database
	if err := s.habitRepo.Create(ctx, &habit); err != nil {
		return nil, err
	}

	// Return the created habit
	response := habit.ToResponse()
	return &response, nil
}

// GetHabit gets a specific habit
func (s *HabitService) GetHabit(ctx context.Context, userID uint, habitID uint) (*models.HabitResponse, error) {
	// Find the habit
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}

	// Verify the habit belongs to the user
	if habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Get the current streak for this habit
	currentStreak, _ := s.streakRepo.FindActiveByHabitID(ctx, habit.ID)

	habitResponse := habit.ToResponse()
	if currentStreak != nil {
		streakResponse := currentStreak.ToResponse()
		habitResponse.CurrentStreak = &streakResponse
	}

	return &habitResponse, nil
}

// UpdateHabit updates a habit
func (s *HabitService) UpdateHabit(ctx context.Context, userID uint, habitID uint, req UpdateHabitRequest) (*models.HabitResponse, error) {
	// Find the habit
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}

	// Verify the habit belongs to the user
	if habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Update the habit
	habit.Name = req.Name
	habit.Description = req.Description
	habit.Color = req.Color
	habit.Icon = req.Icon
	if req.IsActive != nil {
		habit.IsActive = *req.IsActive
	}

	// Save the changes
	if err := s.habitRepo.Update(ctx, habit); err != nil {
		return nil, err
	}

	// Return the updated habit
	response := habit.ToResponse()
	return &response, nil
}

// DeleteHabit deletes a habit
func (s *HabitService) DeleteHabit(ctx context.Context, userID uint, habitID uint) error {
	// Find the habit
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return err
	}
	if habit == nil {
		return errors.New("habit not found")
	}

	// Verify the habit belongs to the user
	if habit.UserID != userID {
		return errors.New("habit not found")
	}

	// Delete the habit
	return s.habitRepo.Delete(ctx, habitID)
}
