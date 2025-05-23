package repository

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GormHabitRepository implements HabitRepository using GORM
type GormHabitRepository struct {
	db *gorm.DB
}

// NewHabitRepository creates a new habit repository
func NewHabitRepository(db *gorm.DB) HabitRepository {
	return &GormHabitRepository{db: db}
}

// Create creates a new habit
func (r *GormHabitRepository) Create(ctx context.Context, habit *models.Habit) error {
	result := r.db.WithContext(ctx).Create(habit)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create habit")
		return result.Error
	}
	return nil
}

// FindByID finds a habit by ID
func (r *GormHabitRepository) FindByID(ctx context.Context, id uint) (*models.Habit, error) {
	var habit models.Habit
	result := r.db.WithContext(ctx).First(&habit, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to find habit by ID")
		return nil, result.Error
	}
	return &habit, nil
}

// FindByUserID finds all habits for a user
func (r *GormHabitRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Habit, error) {
	var habits []models.Habit
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&habits)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("userID", userID).Msg("Failed to find habits by user ID")
		return nil, result.Error
	}
	return habits, nil
}

// Update updates a habit
func (r *GormHabitRepository) Update(ctx context.Context, habit *models.Habit) error {
	result := r.db.WithContext(ctx).Save(habit)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", habit.ID).Msg("Failed to update habit")
		return result.Error
	}
	return nil
}

// Delete deletes a habit
func (r *GormHabitRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Habit{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to delete habit")
		return result.Error
	}
	return nil
}
