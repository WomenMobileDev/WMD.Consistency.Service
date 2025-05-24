package repository

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GormAchievementRepository implements AchievementRepository using GORM
type GormAchievementRepository struct {
	db *gorm.DB
}

// NewAchievementRepository creates a new achievement repository
func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &GormAchievementRepository{db: db}
}

// Create creates a new achievement
func (r *GormAchievementRepository) Create(ctx context.Context, achievement *models.Achievement) error {
	result := r.db.WithContext(ctx).Create(achievement)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create achievement")
		return result.Error
	}
	return nil
}

// FindByID finds an achievement by ID
func (r *GormAchievementRepository) FindByID(ctx context.Context, id uint) (*models.Achievement, error) {
	var achievement models.Achievement
	result := r.db.WithContext(ctx).First(&achievement, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to find achievement by ID")
		return nil, result.Error
	}
	return &achievement, nil
}

// FindByUserID finds all achievements for a user
func (r *GormAchievementRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Achievement, error) {
	var achievements []models.Achievement
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("achieved_at DESC").Find(&achievements)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("userID", userID).Msg("Failed to find achievements by user ID")
		return nil, result.Error
	}
	return achievements, nil
}

// FindByHabitID finds all achievements for a habit
func (r *GormAchievementRepository) FindByHabitID(ctx context.Context, habitID uint) ([]models.Achievement, error) {
	var achievements []models.Achievement
	result := r.db.WithContext(ctx).Where("habit_id = ?", habitID).Order("achieved_at DESC").Find(&achievements)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("habitID", habitID).Msg("Failed to find achievements by habit ID")
		return nil, result.Error
	}
	return achievements, nil
}

// Delete deletes an achievement
func (r *GormAchievementRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Achievement{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to delete achievement")
		return result.Error
	}
	return nil
}
