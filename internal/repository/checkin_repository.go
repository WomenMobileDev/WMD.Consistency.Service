package repository

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GormCheckInRepository implements CheckInRepository using GORM
type GormCheckInRepository struct {
	db *gorm.DB
}

// NewCheckInRepository creates a new check-in repository
func NewCheckInRepository(db *gorm.DB) CheckInRepository {
	return &GormCheckInRepository{db: db}
}

// Create creates a new check-in
func (r *GormCheckInRepository) Create(ctx context.Context, checkIn *models.HabitCheckIn) error {
	result := r.db.WithContext(ctx).Create(checkIn)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create check-in")
		return result.Error
	}
	return nil
}

// FindByID finds a check-in by ID
func (r *GormCheckInRepository) FindByID(ctx context.Context, id uint) (*models.HabitCheckIn, error) {
	var checkIn models.HabitCheckIn
	result := r.db.WithContext(ctx).First(&checkIn, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to find check-in by ID")
		return nil, result.Error
	}
	return &checkIn, nil
}

// FindByStreakID finds all check-ins for a streak
func (r *GormCheckInRepository) FindByStreakID(ctx context.Context, streakID uint) ([]models.HabitCheckIn, error) {
	var checkIns []models.HabitCheckIn
	result := r.db.WithContext(ctx).Where("streak_id = ?", streakID).Order("check_in_date DESC").Find(&checkIns)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("streakID", streakID).Msg("Failed to find check-ins by streak ID")
		return nil, result.Error
	}
	return checkIns, nil
}

// FindByDate finds a check-in by date
func (r *GormCheckInRepository) FindByDate(ctx context.Context, streakID uint, date string) (*models.HabitCheckIn, error) {
	var checkIn models.HabitCheckIn
	result := r.db.WithContext(ctx).Where("streak_id = ? AND check_in_date = ?", streakID, date).First(&checkIn)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("streakID", streakID).Str("date", date).Msg("Failed to find check-in by date")
		return nil, result.Error
	}
	return &checkIn, nil
}

// FindLatestByStreakID finds the latest check-in for a streak
func (r *GormCheckInRepository) FindLatestByStreakID(ctx context.Context, streakID uint) (*models.HabitCheckIn, error) {
	var checkIn models.HabitCheckIn
	result := r.db.WithContext(ctx).Where("streak_id = ?", streakID).Order("check_in_date DESC").First(&checkIn)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("streakID", streakID).Msg("Failed to find latest check-in")
		return nil, result.Error
	}
	return &checkIn, nil
}

// Delete deletes a check-in
func (r *GormCheckInRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.HabitCheckIn{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to delete check-in")
		return result.Error
	}
	return nil
}
