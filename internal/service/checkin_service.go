// Package service contains the business logic for habit tracking, including check-in and streak management.
package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
	"gorm.io/datatypes"
)

// CheckInService handles check-in-related business logic
type CheckInService struct {
	habitRepo       repository.HabitRepository
	streakRepo      repository.StreakRepository
	checkInRepo     repository.CheckInRepository
	achievementRepo repository.AchievementRepository
}

// NewCheckInService creates a new check-in service
func NewCheckInService(
	habitRepo repository.HabitRepository,
	streakRepo repository.StreakRepository,
	checkInRepo repository.CheckInRepository,
	achievementRepo repository.AchievementRepository,
) *CheckInService {
	return &CheckInService{
		habitRepo:       habitRepo,
		streakRepo:      streakRepo,
		checkInRepo:     checkInRepo,
		achievementRepo: achievementRepo,
	}
}

type CheckInRequest struct {
	Notes string `json:"notes"`
}

func (s *CheckInService) CheckIn(ctx context.Context, userID uint, habitID uint, req CheckInRequest) (*models.HabitCheckInResponse, error) {
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, &models.AppError{
			Code:    "habit_not_found",
			Message: "Habit not found",
		}
	}

	streak, err := s.streakRepo.FindActiveByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	if streak == nil {
		return nil, &models.AppError{
			Code:    "STREAK_NOT_FOUND",
			Message: "No active streak found. Please start a new streak first",
		}
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	todayStr := today.Format("2006-01-02")

	existingCheckIn, err := s.checkInRepo.FindByDate(ctx, streak.ID, todayStr)
	if err != nil {
		return nil, err
	}
	if existingCheckIn != nil {
		return nil, &models.AppError{
			Code:    "ALREADY_CHECKED_IN",
			Message: "Already checked in today",
		}
	}

	latestCheckIn, err := s.checkInRepo.FindLatestByStreakID(ctx, streak.ID)
	if err != nil {
		return nil, err
	}

	if latestCheckIn != nil {
		yesterday := today.AddDate(0, 0, -1)
		yesterdayStr := yesterday.Format("2006-01-02")
		latestDateStr := latestCheckIn.CheckInDate.Format("2006-01-02")

		if latestDateStr != yesterdayStr {
			streakID := streak.ID

			checkIns, err := s.checkInRepo.FindByStreakID(ctx, streakID)
			if err != nil {
				return nil, err
			}

			for _, checkIn := range checkIns {
				if err := s.checkInRepo.Delete(ctx, checkIn.ID); err != nil {
					return nil, err
				}
			}

			if err := s.streakRepo.Delete(ctx, streakID); err != nil {
				return nil, err
			}

			return nil, &models.AppError{
				Code:    "STREAK_BROKEN",
				Message: "Streak broken! You missed a day. The streak has been deleted. Please start a new streak.",
			}
		}
	}

	checkIn := models.HabitCheckIn{
		StreakID:    streak.ID,
		CheckInDate: today,
		Notes:       req.Notes,
	}

	if err := s.checkInRepo.Create(ctx, &checkIn); err != nil {
		return nil, err
	}

	streak.CurrentStreak++
	streak.LastCheckInDate = &today

	if streak.CurrentStreak >= streak.TargetDays {
		streak.Status = "completed"
		streak.CompletedAt = &today

		achievement := models.Achievement{
			UserID:          habit.UserID,
			HabitID:         habit.ID,
			AchievementType: "streak_completed",
			TargetDays:      streak.TargetDays,
			Metadata:        datatypes.JSON([]byte(`{"streak_id": ` + strconv.Itoa(int(streak.ID)) + `}`)),
		}

		if err := s.achievementRepo.Create(ctx, &achievement); err != nil {
			return nil, err
		}
	}

	if streak.CurrentStreak > streak.MaxStreakAchieved {
		streak.MaxStreakAchieved = streak.CurrentStreak
	}

	if err := s.streakRepo.Update(ctx, streak); err != nil {
		return nil, err
	}

	response := checkIn.ToResponse()
	return &response, nil
}

// ListCheckIns lists all check-ins for a habit
func (s *CheckInService) ListCheckIns(ctx context.Context, userID uint, habitID uint) ([]models.HabitCheckInResponse, error) {
	// Verify the habit belongs to the user
	habit, err := s.habitRepo.FindByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil || habit.UserID != userID {
		return nil, errors.New("habit not found")
	}

	// Get all streaks for this habit
	streaks, err := s.streakRepo.FindByHabitID(ctx, habitID)
	if err != nil {
		return nil, err
	}

	if len(streaks) == 0 {
		return []models.HabitCheckInResponse{}, nil
	}

	// Collect all streak IDs
	var streakIDs []uint
	for _, streak := range streaks {
		streakIDs = append(streakIDs, streak.ID)
	}

	// Get all check-ins for the streaks
	var allCheckIns []models.HabitCheckIn
	for _, streakID := range streakIDs {
		checkIns, err := s.checkInRepo.FindByStreakID(ctx, streakID)
		if err != nil {
			return nil, err
		}
		allCheckIns = append(allCheckIns, checkIns...)
	}

	// Convert check-ins to responses
	var responses []models.HabitCheckInResponse
	for _, checkIn := range allCheckIns {
		responses = append(responses, checkIn.ToResponse())
	}

	return responses, nil
}
