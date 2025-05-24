package repository

import (
	"context"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

// HabitRepository defines the interface for habit data access
type HabitRepository interface {
	Create(ctx context.Context, habit *models.Habit) error
	FindByID(ctx context.Context, id uint) (*models.Habit, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) error
	Delete(ctx context.Context, id uint) error
}

// StreakRepository defines the interface for streak data access
type StreakRepository interface {
	Create(ctx context.Context, streak *models.HabitStreak) error
	FindByID(ctx context.Context, id uint) (*models.HabitStreak, error)
	FindByHabitID(ctx context.Context, habitID uint) ([]models.HabitStreak, error)
	FindActiveByHabitID(ctx context.Context, habitID uint) (*models.HabitStreak, error)
	Update(ctx context.Context, streak *models.HabitStreak) error
	Delete(ctx context.Context, id uint) error
}

// CheckInRepository defines the interface for check-in data access
type CheckInRepository interface {
	Create(ctx context.Context, checkIn *models.HabitCheckIn) error
	FindByID(ctx context.Context, id uint) (*models.HabitCheckIn, error)
	FindByStreakID(ctx context.Context, streakID uint) ([]models.HabitCheckIn, error)
	FindByDate(ctx context.Context, streakID uint, date string) (*models.HabitCheckIn, error)
	FindLatestByStreakID(ctx context.Context, streakID uint) (*models.HabitCheckIn, error)
	Delete(ctx context.Context, id uint) error
}

// AchievementRepository defines the interface for achievement data access
type AchievementRepository interface {
	Create(ctx context.Context, achievement *models.Achievement) error
	FindByID(ctx context.Context, id uint) (*models.Achievement, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Achievement, error)
	FindByHabitID(ctx context.Context, habitID uint) ([]models.Achievement, error)
	Delete(ctx context.Context, id uint) error
}
