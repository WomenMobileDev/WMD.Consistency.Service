package service

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
)

// AchievementService handles achievement-related business logic
type AchievementService struct {
	achievementRepo repository.AchievementRepository
	habitRepo       repository.HabitRepository
}

// NewAchievementService creates a new achievement service
func NewAchievementService(achievementRepo repository.AchievementRepository, habitRepo repository.HabitRepository) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		habitRepo:       habitRepo,
	}
}

// ListAchievements lists all achievements for a user
func (s *AchievementService) ListAchievements(ctx context.Context, userID uint) ([]models.AchievementResponse, error) {
	// Find all achievements for the user
	achievements, err := s.achievementRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get habit names for each achievement
	var responses []models.AchievementResponse
	for _, achievement := range achievements {
		response := achievement.ToResponse()

		// Get the habit name
		habit, err := s.habitRepo.FindByID(ctx, achievement.HabitID)
		if err == nil && habit != nil {
			response.HabitName = habit.Name
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetAchievement gets a specific achievement
func (s *AchievementService) GetAchievement(ctx context.Context, userID uint, achievementID uint) (*models.AchievementResponse, error) {
	// Find the achievement
	achievement, err := s.achievementRepo.FindByID(ctx, achievementID)
	if err != nil {
		return nil, err
	}
	if achievement == nil {
		return nil, errors.New("achievement not found")
	}

	// Verify the achievement belongs to the user
	if achievement.UserID != userID {
		return nil, errors.New("achievement not found")
	}

	// Get the habit name
	response := achievement.ToResponse()
	habit, err := s.habitRepo.FindByID(ctx, achievement.HabitID)
	if err == nil && habit != nil {
		response.HabitName = habit.Name
	}

	return &response, nil
}

// ListHabitAchievements lists all achievements for a habit
func (s *AchievementService) ListHabitAchievements(ctx context.Context, userID uint, habitID uint) ([]models.AchievementResponse, error) {
	// Verify the habit belongs to the user
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Find all achievements for the habit
	achievements, err := s.achievementRepo.FindByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	// Convert achievements to responses
	var responses []models.AchievementResponse
	for _, achievement := range achievements {
		response := achievement.ToResponse()
		response.HabitName = habit.Name
		responses = append(responses, response)
	}

	return responses, nil
}
