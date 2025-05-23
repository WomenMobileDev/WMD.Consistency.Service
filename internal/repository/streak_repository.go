package repository

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GormStreakRepository implements StreakRepository using GORM
type GormStreakRepository struct {
	db *gorm.DB
}

// NewStreakRepository creates a new streak repository
func NewStreakRepository(db *gorm.DB) StreakRepository {
	return &GormStreakRepository{db: db}
}

// Create creates a new streak
func (r *GormStreakRepository) Create(ctx context.Context, streak *models.HabitStreak) error {
	result := r.db.WithContext(ctx).Create(streak)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create streak")
		return result.Error
	}
	return nil
}

// FindByID finds a streak by ID
func (r *GormStreakRepository) FindByID(ctx context.Context, id uint) (*models.HabitStreak, error) {
	var streak models.HabitStreak
	result := r.db.WithContext(ctx).First(&streak, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to find streak by ID")
		return nil, result.Error
	}
	return &streak, nil
}

// FindByHabitID finds all streaks for a habit
func (r *GormStreakRepository) FindByHabitID(ctx context.Context, habitID uint) ([]models.HabitStreak, error) {
	var streaks []models.HabitStreak
	result := r.db.WithContext(ctx).Where("habit_id = ?", habitID).Order("created_at DESC").Find(&streaks)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("habitID", habitID).Msg("Failed to find streaks by habit ID")
		return nil, result.Error
	}
	return streaks, nil
}

// FindActiveByHabitID finds the active streak for a habit
func (r *GormStreakRepository) FindActiveByHabitID(ctx context.Context, habitID uint) (*models.HabitStreak, error) {
	var streak models.HabitStreak
	result := r.db.WithContext(ctx).Where("habit_id = ? AND status = 'active'", habitID).First(&streak)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("habitID", habitID).Msg("Failed to find active streak by habit ID")
		return nil, result.Error
	}
	return &streak, nil
}

// Update updates a streak
func (r *GormStreakRepository) Update(ctx context.Context, streak *models.HabitStreak) error {
	result := r.db.WithContext(ctx).Save(streak)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", streak.ID).Msg("Failed to update streak")
		return result.Error
	}
	return nil
}

// Delete deletes a streak
func (r *GormStreakRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.HabitStreak{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to delete streak")
		return result.Error
	}
	return nil
}
