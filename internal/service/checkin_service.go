package service

import (
	"context"
	"errors"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
	"gorm.io/datatypes"
	"strconv"
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

// CheckInRequest represents the request for checking in
type CheckInRequest struct {
	Notes string `json:"notes"`
}

// CheckIn checks in for a habit
func (s *CheckInService) CheckIn(ctx context.Context, userID uint, habitID uint, req CheckInRequest) (*models.HabitCheckInResponse, error) {
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
		return nil, errors.New("no active streak found. Please start a new streak first")
	}

	// Get today's date in UTC
	today := time.Now().UTC().Truncate(24 * time.Hour)
	todayStr := today.Format("2006-01-02")

	// Check if already checked in today
	existingCheckIn, err := s.checkInRepo.FindByDate(ctx, streak.ID, todayStr)
	if err != nil {
		return nil, err
	}
	if existingCheckIn != nil {
		return nil, errors.New("already checked in today")
	}

	// Check if the last check-in was yesterday or if this is the first check-in
	latestCheckIn, err := s.checkInRepo.FindLatestByStreakID(ctx, streak.ID)
	if err != nil {
		return nil, err
	}

	// If there's a last check-in, verify it was yesterday
	if latestCheckIn != nil {
		yesterday := today.AddDate(0, 0, -1)
		yesterdayStr := yesterday.Format("2006-01-02")
		latestDateStr := latestCheckIn.CheckInDate.Format("2006-01-02")

		if latestDateStr != yesterdayStr {
			// If the last check-in wasn't yesterday, the streak is broken
			// Instead of updating the streak as failed, delete it
			streakID := streak.ID
			
			// Delete all check-ins for this streak first (to maintain referential integrity)
			checkIns, err := s.checkInRepo.FindByStreakID(ctx, streakID)
			if err != nil {
				return nil, err
			}
			
			// Delete each check-in by its ID
			for _, checkIn := range checkIns {
				if err := s.checkInRepo.Delete(ctx, checkIn.ID); err != nil {
					return nil, err
				}
			}
			
			// Now delete the streak itself by its ID
			if err := s.streakRepo.Delete(ctx, streakID); err != nil {
				return nil, err
			}

			return nil, errors.New("streak broken! You missed a day. The streak has been deleted. Please start a new streak")
		}
	}

	// Create a new check-in
	checkIn := models.HabitCheckIn{
		StreakID:    streak.ID,
		CheckInDate: today,
		Notes:       req.Notes,
	}

	// Begin transaction
	// In a real application, you would use a transaction here
	// For simplicity, we'll just perform the operations sequentially

	// Save the check-in
	if err := s.checkInRepo.Create(ctx, &checkIn); err != nil {
		return nil, err
	}

	// Update the streak
	streak.CurrentStreak++
	streak.LastCheckInDate = &today

	// Check if the streak has reached its target
	if streak.CurrentStreak >= streak.TargetDays {
		streak.Status = "completed"
		streak.CompletedAt = &today

		// Create an achievement
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

	// Update max streak if current streak is higher
	if streak.CurrentStreak > streak.MaxStreakAchieved {
		streak.MaxStreakAchieved = streak.CurrentStreak
	}

	// Save the streak
	if err := s.streakRepo.Update(ctx, streak); err != nil {
		return nil, err
	}

	// Return the check-in
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
