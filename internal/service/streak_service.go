package service

import (
	"context"
	"errors"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
)

// StreakService handles streak-related business logic
type StreakService struct {
	habitRepo   repository.HabitRepository
	streakRepo  repository.StreakRepository
	checkInRepo repository.CheckInRepository
}

// NewStreakService creates a new streak service
func NewStreakService(habitRepo repository.HabitRepository, streakRepo repository.StreakRepository, checkInRepo repository.CheckInRepository) *StreakService {
	return &StreakService{
		habitRepo:   habitRepo,
		streakRepo:  streakRepo,
		checkInRepo: checkInRepo,
	}
}

// CreateStreakRequest represents the request for creating a streak
type CreateStreakRequest struct {
	TargetDays int `json:"target_days" binding:"required,gt=0"`
}

// CreateStreak creates a new streak for a habit
func (s *StreakService) CreateStreak(ctx context.Context, userID uint, habitID uint, req CreateStreakRequest) (*models.HabitStreakResponse, error) {
	// Verify the habit belongs to the user
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Check if there's already an active streak
	activeStreak, err := s.streakRepo.FindActiveByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if activeStreak != nil {
		return nil, errors.New("an active streak already exists for this habit")
	}

	// Check if the target days is valid
	if req.TargetDays <= 0 {
		return nil, errors.New("target days must be greater than 0")
	}

	// Create a new streak
	now := time.Now().UTC()
	streak := models.HabitStreak{
		HabitID:       habitID,
		TargetDays:    req.TargetDays,
		CurrentStreak: 0,
		StartDate:     now,
		Status:        "active",
	}

	// Save the streak to the database
	if err := s.streakRepo.Create(ctx, &streak); err != nil {
		return nil, err
	}

	// Return the created streak
	response := streak.ToResponse()
	return &response, nil
}

// ListStreaks lists all streaks for a habit
func (s *StreakService) ListStreaks(ctx context.Context, userID uint, habitID uint) ([]models.HabitStreakResponse, error) {
	// Verify the habit belongs to the user
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Find all streaks for the habit
	streaks, err := s.streakRepo.FindByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	// Convert streaks to responses
	var responses []models.HabitStreakResponse
	for _, streak := range streaks {
		responses = append(responses, streak.ToResponse())
	}

	return responses, nil
}

// GetCurrentStreak gets the current active streak for a habit
func (s *StreakService) GetCurrentStreak(ctx context.Context, userID uint, habitID uint) (*models.HabitStreakResponse, error) {
	// Verify the habit belongs to the user
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Find the active streak for the habit
	streak, err := s.streakRepo.FindActiveByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if streak == nil {
		return nil, errors.New("no active streak found")
	}

	// Get the check-ins for this streak
	checkIns, err := s.checkInRepo.FindByStreakID(ctx, streak.ID)
	if err != nil {
		return nil, err
	}

	// Convert check-ins to responses
	var checkInResponses []models.HabitCheckInResponse
	for _, checkIn := range checkIns {
		checkInResponses = append(checkInResponses, checkIn.ToResponse())
	}

	// Create the response
	response := streak.ToResponse()
	response.CheckIns = checkInResponses

	return &response, nil
}
