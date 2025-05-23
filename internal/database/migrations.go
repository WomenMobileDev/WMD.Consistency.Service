package database

import (
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/rs/zerolog/log"
)

// RunMigrations runs all database migrations
func (db *Database) RunMigrations() error {
	log.Info().Msg("Running database migrations...")

	// List of models to migrate
	models := []interface{}{
		&models.User{},
		&models.Habit{},
		&models.HabitStreak{},
		&models.HabitCheckIn{},
		&models.Achievement{},
	}

	// Run migrations
	err := db.DB.AutoMigrate(models...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to run migrations")
		return err
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil
}
