package repository

import (
	"context"
	"errors"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

// Create creates a new user
func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create user")
		return result.Error
	}
	return nil
}

// FindByID finds a user by ID
func (r *GormUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to find user by ID")
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(result.Error).Str("email", email).Msg("Failed to find user by email")
		return nil, result.Error
	}
	return &user, nil
}

// Update updates a user
func (r *GormUserRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", user.ID).Msg("Failed to update user")
		return result.Error
	}
	return nil
}

// Delete deletes a user
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("id", id).Msg("Failed to delete user")
		return result.Error
	}
	return nil
}
