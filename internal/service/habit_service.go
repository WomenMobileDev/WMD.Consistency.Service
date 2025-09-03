package service

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
)

type HabitService struct {
	habitRepo  repository.HabitRepository
	streakRepo repository.StreakRepository
}

func NewHabitService(habitRepo repository.HabitRepository, streakRepo repository.StreakRepository) *HabitService {
	return &HabitService{
		habitRepo:  habitRepo,
		streakRepo: streakRepo,
	}
}

type CreateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" binding:"omitempty,len=7"` // Hex color code
	Icon        string `json:"icon"`
	TargetDays  int    `json:"target_days" binding:"required"`
	GoalUnit    string `json:"goal_unit" binding:"required,oneof=days weeks months"`
}

type UpdateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" binding:"omitempty,len=7"` // Hex color code
	Icon        string `json:"icon"`
	IsActive    *bool  `json:"is_active"`
}

func (s *HabitService) ListHabits(ctx context.Context, userID uint) ([]models.HabitResponse, error) {
	habits, err := s.habitRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []models.HabitResponse
	for _, habit := range habits {
		currentStreak, _ := s.streakRepo.FindActiveByHabitID(ctx, habit.ID)

		habitResponse := habit.ToResponseWithStreak(currentStreak)
		responses = append(responses, habitResponse)
	}

	return responses, nil
}

func (s *HabitService) CreateHabit(ctx context.Context, userID uint, req CreateHabitRequest) (*models.HabitResponse, error) {
	habit := models.Habit{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		IsActive:    true,
	}

	if err := s.habitRepo.Create(ctx, &habit); err != nil {
		return nil, err
	}

	// Calculate target days based on goal unit
	targetDays := req.TargetDays
	switch req.GoalUnit {
	case "days":
		targetDays = req.TargetDays
	case "weeks":
		targetDays = req.TargetDays * 7
	case "months":
		targetDays = req.TargetDays * 30
	}

	// Create streak service and start the streak
	streakService := &StreakService{
		habitRepo:  s.habitRepo,
		streakRepo: s.streakRepo,
	}
	streakReq := CreateStreakRequest{
		TargetDays: targetDays,
	}

	createdStreak, err := streakService.CreateStreak(ctx, userID, habit.ID, streakReq)
	if err != nil {
		return nil, err
	}

	// Convert streak response to streak model for the response
	var currentStreak *models.HabitStreak
	if createdStreak != nil {
		currentStreak = &models.HabitStreak{
			ID:                createdStreak.ID,
			HabitID:           createdStreak.HabitID,
			TargetDays:        createdStreak.TargetDays,
			CurrentStreak:     createdStreak.CurrentStreak,
			MaxStreakAchieved: createdStreak.MaxStreakAchieved,
			StartDate:         createdStreak.StartDate,
			Status:            createdStreak.Status,
			LastCheckInDate:   createdStreak.LastCheckInDate,
			CompletedAt:       createdStreak.CompletedAt,
			FailedAt:          createdStreak.FailedAt,
		}
	}

	response := habit.ToResponseWithStreak(currentStreak)
	return &response, nil
}

func (s *HabitService) GetHabit(ctx context.Context, userID uint, habitID uint) (*models.HabitResponse, error) {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}

	if habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	currentStreak, _ := s.streakRepo.FindActiveByHabitID(ctx, habit.ID)

	habitResponse := habit.ToResponseWithStreak(currentStreak)

	return &habitResponse, nil
}

func (s *HabitService) UpdateHabit(ctx context.Context, userID uint, habitID uint, req UpdateHabitRequest) (*models.HabitResponse, error) {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}

	if habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	habit.Name = req.Name
	habit.Description = req.Description
	habit.Color = req.Color
	habit.Icon = req.Icon
	if req.IsActive != nil {
		habit.IsActive = *req.IsActive
	}

	if err := s.habitRepo.Update(ctx, habit); err != nil {
		return nil, err
	}

	currentStreak, _ := s.streakRepo.FindActiveByHabitID(ctx, habit.ID)

	response := habit.ToResponseWithStreak(currentStreak)
	return &response, nil
}

func (s *HabitService) DeleteHabit(ctx context.Context, userID uint, habitID uint) error {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return err
	}
	if habit == nil {
		return errors.New("habit not found")
	}

	if habit.UserID != userID {
		return errors.New("habit not found")
	}

	return s.habitRepo.Delete(ctx, habitID)
}
