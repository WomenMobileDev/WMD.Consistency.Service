package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Configure logging
	setupLogging(cfg)
	log.Info().Msg("Starting Database Seeder")

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Connected to database")

	// Auto-migrate database schemas
	if err := db.AutoMigrate(); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database schemas")
	}
	log.Info().Msg("Database migrations completed")

	// Create repositories
	userRepo := repository.NewUserRepository(db.DB)
	habitRepo := repository.NewHabitRepository(db.DB)
	streakRepo := repository.NewStreakRepository(db.DB)
	checkInRepo := repository.NewCheckInRepository(db.DB)
	achievementRepo := repository.NewAchievementRepository(db.DB)

	// Create context
	ctx := context.Background()

	// Seed users
	users, err := seedUsers(ctx, userRepo)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed users")
	}
	log.Info().Int("count", len(users)).Msg("Users seeded")

	// Seed habits
	habits, err := seedHabits(ctx, habitRepo, users)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed habits")
	}
	log.Info().Int("count", len(habits)).Msg("Habits seeded")

	// Seed streaks
	streaks, err := seedStreaks(ctx, streakRepo, habits)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed streaks")
	}
	log.Info().Int("count", len(streaks)).Msg("Streaks seeded")

	// Seed check-ins
	checkIns, err := seedCheckIns(ctx, checkInRepo, streaks)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed check-ins")
	}
	log.Info().Int("count", len(checkIns)).Msg("Check-ins seeded")

	// Seed achievements
	achievements, err := seedAchievements(ctx, achievementRepo, habits)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed achievements")
	}
	log.Info().Int("count", len(achievements)).Msg("Achievements seeded")

	log.Info().Msg("Database seeding completed")
}

func setupLogging(cfg *config.Config) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure pretty logging for development
	if cfg.Logger.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Set default time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func seedUsers(ctx context.Context, repo repository.UserRepository) ([]*models.User, error) {
	users := []*models.User{
		{
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			Name:  "Jane Smith",
			Email: "jane@example.com",
		},
	}

	// Set passwords
	for _, user := range users {
		if err := user.SetPassword("password123"); err != nil {
			return nil, fmt.Errorf("failed to set password: %w", err)
		}
	}

	// Save users to database
	for _, user := range users {
		if err := repo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return users, nil
}

func seedHabits(ctx context.Context, repo repository.HabitRepository, users []*models.User) ([]*models.Habit, error) {
	habits := []*models.Habit{
		{
			UserID:      users[0].ID,
			Name:        "Morning Meditation",
			Description: "Meditate for 10 minutes every morning",
			Color:       "#4287f5",
			Icon:        "meditation",
			IsActive:    true,
		},
		{
			UserID:      users[0].ID,
			Name:        "Daily Exercise",
			Description: "Exercise for 30 minutes every day",
			Color:       "#f54242",
			Icon:        "exercise",
			IsActive:    true,
		},
		{
			UserID:      users[1].ID,
			Name:        "Read a Book",
			Description: "Read for 20 minutes every day",
			Color:       "#42f5a7",
			Icon:        "book",
			IsActive:    true,
		},
	}

	// Save habits to database
	for _, habit := range habits {
		if err := repo.Create(ctx, habit); err != nil {
			return nil, fmt.Errorf("failed to create habit: %w", err)
		}
	}

	return habits, nil
}

func seedStreaks(ctx context.Context, repo repository.StreakRepository, habits []*models.Habit) ([]*models.HabitStreak, error) {
	now := time.Now().UTC()
	streaks := []*models.HabitStreak{
		{
			HabitID:          habits[0].ID,
			TargetDays:       7,
			CurrentStreak:    3,
			MaxStreakAchieved: 3,
			StartDate:        now.AddDate(0, 0, -3),
			LastCheckInDate:  &now,
			Status:           "active",
		},
		{
			HabitID:          habits[1].ID,
			TargetDays:       14,
			CurrentStreak:    0,
			MaxStreakAchieved: 0,
			StartDate:        now,
			Status:           "active",
		},
		{
			HabitID:          habits[2].ID,
			TargetDays:       21,
			CurrentStreak:    5,
			MaxStreakAchieved: 5,
			StartDate:        now.AddDate(0, 0, -5),
			LastCheckInDate:  &now,
			Status:           "active",
		},
	}

	// Save streaks to database
	for _, streak := range streaks {
		if err := repo.Create(ctx, streak); err != nil {
			return nil, fmt.Errorf("failed to create streak: %w", err)
		}
	}

	return streaks, nil
}

func seedCheckIns(ctx context.Context, repo repository.CheckInRepository, streaks []*models.HabitStreak) ([]*models.HabitCheckIn, error) {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	checkIns := []*models.HabitCheckIn{}

	// Create check-ins for the first streak (3 days)
	for i := 0; i < 3; i++ {
		checkIns = append(checkIns, &models.HabitCheckIn{
			StreakID:    streaks[0].ID,
			CheckInDate: now.AddDate(0, 0, -i),
			Notes:       fmt.Sprintf("Day %d of meditation", i+1),
		})
	}

	// Create check-ins for the third streak (5 days)
	for i := 0; i < 5; i++ {
		checkIns = append(checkIns, &models.HabitCheckIn{
			StreakID:    streaks[2].ID,
			CheckInDate: now.AddDate(0, 0, -i),
			Notes:       fmt.Sprintf("Day %d of reading", i+1),
		})
	}

	// Save check-ins to database
	for _, checkIn := range checkIns {
		if err := repo.Create(ctx, checkIn); err != nil {
			return nil, fmt.Errorf("failed to create check-in: %w", err)
		}
	}

	return checkIns, nil
}

func seedAchievements(ctx context.Context, repo repository.AchievementRepository, habits []*models.Habit) ([]*models.Achievement, error) {
	now := time.Now().UTC()
	achievements := []*models.Achievement{
		{
			UserID:          habits[0].UserID,
			HabitID:         habits[0].ID,
			AchievementType: "streak_milestone",
			TargetDays:      3,
			AchievedAt:      now.AddDate(0, 0, -1),
			Metadata:        datatypes.JSON([]byte(`{"streak_id": 1}`)),
		},
		{
			UserID:          habits[2].UserID,
			HabitID:         habits[2].ID,
			AchievementType: "streak_milestone",
			TargetDays:      5,
			AchievedAt:      now,
			Metadata:        datatypes.JSON([]byte(`{"streak_id": 3}`)),
		},
	}

	// Save achievements to database
	for _, achievement := range achievements {
		if err := repo.Create(ctx, achievement); err != nil {
			return nil, fmt.Errorf("failed to create achievement: %w", err)
		}
	}

	return achievements, nil
}
