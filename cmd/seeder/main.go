package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/config"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/database"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	setupLogging(cfg)
	log.Info().Msg("Starting Database Seeder")

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Connected to database")

	if err := db.AutoMigrate(); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database schemas")
	}
	log.Info().Msg("Database migrations completed")

	userRepo := repository.NewUserRepository(db.DB)
	habitRepo := repository.NewHabitRepository(db.DB)
	streakRepo := repository.NewStreakRepository(db.DB)
	checkInRepo := repository.NewCheckInRepository(db.DB)
	achievementRepo := repository.NewAchievementRepository(db.DB)

	ctx := context.Background()

	if err := clearExistingData(db.DB); err != nil {
		log.Fatal().Err(err).Msg("Failed to clear existing data")
	}
	log.Info().Msg("Existing data cleared")

	users, err := seedUsers(ctx, userRepo)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed users")
	}
	log.Info().Int("count", len(users)).Msg("Users seeded")

	habits, err := seedHabits(ctx, habitRepo, users)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed habits")
	}
	log.Info().Int("count", len(habits)).Msg("Habits seeded")

	err = seedComprehensiveData(ctx, habitRepo, streakRepo, checkInRepo, achievementRepo, habits)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to seed comprehensive data")
	}

	log.Info().Msg("Database seeding completed successfully!")
	log.Info().Msg("🎯 Rich analytics data created for the past 30 days")
	log.Info().Msg("📊 Multiple habits, streaks, check-ins, and achievements ready for testing")
}

func setupLogging(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.Logger.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func seedUsers(ctx context.Context, repo repository.UserRepository) ([]*models.User, error) {
	createdAt := time.Now().UTC().AddDate(0, 0, -40)

	users := []*models.User{
		{
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: createdAt,
		},
		{
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			CreatedAt: createdAt,
		},
	}

	for _, user := range users {
		if err := user.SetPassword("password123"); err != nil {
			return nil, fmt.Errorf("failed to set password: %w", err)
		}
	}

	for _, user := range users {
		if err := repo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return users, nil
}

func seedHabits(ctx context.Context, repo repository.HabitRepository, users []*models.User) ([]*models.Habit, error) {
	now := time.Now().UTC()

	habitsData := []struct {
		name        string
		description string
		color       string
		icon        string
		daysAgo     int
		isActive    bool
	}{
		{"Morning Meditation", "Meditate for 10 minutes every morning", "#4287f5", "meditation", 30, true},
		{"Daily Exercise", "Exercise for 30 minutes every day", "#f54242", "exercise", 25, true},
		{"Reading Books", "Read for 20 minutes every day", "#42f5a7", "book", 20, true},
		{"Drink Water", "Drink 8 glasses of water daily", "#42d4f5", "water", 18, true},
		{"Learn Programming", "Code for 1 hour daily", "#f59342", "code", 15, true},
		{"Journaling", "Write in journal every evening", "#9342f5", "journal", 12, true},
		{"Healthy Eating", "Eat 5 servings of fruits/vegetables", "#42f542", "nutrition", 10, false}, // Deactivated
		{"Morning Walk", "Take a 20-minute walk", "#f542a7", "walk", 8, true},
	}

	var habits []*models.Habit
	for _, h := range habitsData {
		habit := &models.Habit{
			UserID:      users[0].ID,
			Name:        h.name,
			Description: h.description,
			Color:       h.color,
			Icon:        h.icon,
			IsActive:    h.isActive,
			CreatedAt:   now.AddDate(0, 0, -h.daysAgo),
		}

		if err := repo.Create(ctx, habit); err != nil {
			return nil, fmt.Errorf("failed to create habit: %w", err)
		}
		habits = append(habits, habit)
	}

	return habits, nil
}

func seedComprehensiveData(ctx context.Context, _ repository.HabitRepository, streakRepo repository.StreakRepository, checkInRepo repository.CheckInRepository, achievementRepo repository.AchievementRepository, habits []*models.Habit) error {
	now := time.Now().UTC()

	patterns := map[string][]float64{
		"Morning Meditation": {0.9, 0.8, 0.95, 0.7, 0.85, 0.9}, // Very consistent
		"Daily Exercise":     {0.7, 0.6, 0.8, 0.5, 0.7, 0.75},  // Moderately consistent
		"Reading Books":      {0.6, 0.7, 0.4, 0.8, 0.6, 0.5},   // Variable
		"Drink Water":        {0.8, 0.9, 0.85, 0.8, 0.9, 0.85}, // Very good
		"Learn Programming":  {0.5, 0.3, 0.6, 0.4, 0.7, 0.8},   // Improving over time
		"Journaling":         {0.4, 0.5, 0.3, 0.6, 0.5, 0.4},   // Struggling
		"Healthy Eating":     {0.6, 0.4, 0.2, 0.1, 0.0, 0.0},   // Gave up (deactivated)
		"Morning Walk":       {0.0, 0.0, 0.0, 0.8, 0.9, 0.95},  // Recently started, very good
	}

	for _, habit := range habits {
		pattern, exists := patterns[habit.Name]
		if !exists {
			continue
		}

		err := seedHabitData(ctx, streakRepo, checkInRepo, achievementRepo, habit, pattern, now)
		if err != nil {
			return fmt.Errorf("failed to seed data for habit %s: %w", habit.Name, err)
		}
	}

	log.Info().Msg("✅ Comprehensive dummy data created successfully")
	return nil
}

func seedHabitData(ctx context.Context, streakRepo repository.StreakRepository, checkInRepo repository.CheckInRepository, achievementRepo repository.AchievementRepository, habit *models.Habit, consistencyPattern []float64, now time.Time) error {

	daysActive := int(now.Sub(habit.CreatedAt).Hours() / 24)
	if daysActive > 30 {
		daysActive = 30
	}

	streakTargets := []int{3, 7, 14, 21, 30}
	currentStreakID := uint(0)
	currentStreak := 0
	achievementCount := 0

	segmentLength := 5
	patternIndex := 0

	for day := daysActive; day >= 0; day-- {
		date := now.AddDate(0, 0, -day)

		if date.Before(habit.CreatedAt) {
			continue
		}

		if patternIndex >= len(consistencyPattern) {
			patternIndex = len(consistencyPattern) - 1
		}

		segmentDay := (daysActive - day) % segmentLength
		if segmentDay == 0 && patternIndex < len(consistencyPattern)-1 {
			patternIndex++
		}

		consistency := consistencyPattern[patternIndex]

		actualConsistency := consistency + (rand.Float64()-0.5)*0.3
		if actualConsistency < 0 {
			actualConsistency = 0
		}
		if actualConsistency > 1 {
			actualConsistency = 1
		}

		shouldCheckIn := rand.Float64() < actualConsistency

		if shouldCheckIn {
			if currentStreakID == 0 {
				targetDays := streakTargets[rand.Intn(len(streakTargets))]
				streak := &models.HabitStreak{
					HabitID:           habit.ID,
					TargetDays:        targetDays,
					CurrentStreak:     0,
					MaxStreakAchieved: 0,
					StartDate:         date,
					Status:            "active",
				}

				if err := streakRepo.Create(ctx, streak); err != nil {
					return err
				}
				currentStreakID = streak.ID
				currentStreak = 0
			}

			checkIn := &models.HabitCheckIn{
				StreakID:    currentStreakID,
				CheckInDate: date,
				Notes:       generateCheckInNote(habit.Name, currentStreak+1),
			}

			if err := checkInRepo.Create(ctx, checkIn); err != nil {
				return err
			}

			currentStreak++

			streak, err := streakRepo.FindByID(ctx, currentStreakID)
			if err != nil {
				return err
			}

			streak.CurrentStreak = currentStreak
			streak.LastCheckInDate = &date
			if currentStreak > streak.MaxStreakAchieved {
				streak.MaxStreakAchieved = currentStreak
			}

			if currentStreak >= streak.TargetDays {
				streak.Status = "completed"
				streak.CompletedAt = &date

				achievement := &models.Achievement{
					UserID:          habit.UserID,
					HabitID:         habit.ID,
					AchievementType: "streak_completed",
					TargetDays:      streak.TargetDays,
					AchievedAt:      date,
					Metadata:        datatypes.JSON([]byte(fmt.Sprintf(`{"streak_id": %d, "habit_name": "%s"}`, streak.ID, habit.Name))),
				}

				if err := achievementRepo.Create(ctx, achievement); err != nil {
					return err
				}

				achievementCount++
				currentStreakID = 0
				currentStreak = 0
			}

			if err := streakRepo.Update(ctx, streak); err != nil {
				return err
			}

		} else {
			if currentStreakID != 0 && currentStreak > 0 {
				streak, err := streakRepo.FindByID(ctx, currentStreakID)
				if err != nil {
					return err
				}

				streak.Status = "failed"
				streak.FailedAt = &date

				if err := streakRepo.Update(ctx, streak); err != nil {
					return err
				}

				currentStreakID = 0
				currentStreak = 0
			}
		}
	}

	if achievementCount >= 3 {
		achievement := &models.Achievement{
			UserID:          habit.UserID,
			HabitID:         habit.ID,
			AchievementType: "streak_milestone",
			TargetDays:      3,
			AchievedAt:      now.AddDate(0, 0, -5),
			Metadata:        datatypes.JSON([]byte(fmt.Sprintf(`{"milestone": "3_streaks_completed", "habit_name": "%s"}`, habit.Name))),
		}

		if err := achievementRepo.Create(ctx, achievement); err != nil {
			return err
		}
	}

	return nil
}

func generateCheckInNote(habitName string, day int) string {
	notes := map[string][]string{
		"Morning Meditation": {
			"Feeling centered and calm", "Great session today", "Mind was wandering but pushed through",
			"10 minutes of peace", "Really needed this today", "Focused breathing session",
		},
		"Daily Exercise": {
			"30-min cardio workout", "Strength training complete", "Feeling energized!",
			"Pushed through the resistance", "Great gym session", "Home workout done",
		},
		"Reading Books": {
			"Finished chapter 3", "Learning so much", "Great insights today",
			"20 pages down", "Couldn't put it down", "New concepts absorbed",
		},
		"Drink Water": {
			"8 glasses complete!", "Staying hydrated", "Feeling refreshed",
			"Water intake on track", "Hydration goal met", "Body feels great",
		},
		"Learn Programming": {
			"Solved 3 algorithm problems", "Built a small feature", "Debugging session",
			"Learning new framework", "Code review complete", "Project progress made",
		},
		"Journaling": {
			"Reflected on today's events", "Gratitude practice", "Processing emotions",
			"Daily thoughts recorded", "Mindful writing session", "Self-reflection time",
		},
		"Morning Walk": {
			"Beautiful sunrise walk", "Fresh air and exercise", "20 minutes of movement",
			"Nature therapy", "Energizing start", "Peaceful morning stroll",
		},
	}

	habitNotes, exists := notes[habitName]
	if !exists {
		return fmt.Sprintf("Day %d completed!", day)
	}

	return habitNotes[rand.Intn(len(habitNotes))]
}

func clearExistingData(db *gorm.DB) error {
	log.Info().Msg("Clearing check-ins...")
	if err := db.Exec("DELETE FROM habit_checkins").Error; err != nil {
		return fmt.Errorf("failed to clear check-ins: %w", err)
	}

	log.Info().Msg("Clearing achievements...")
	if err := db.Exec("DELETE FROM achievements").Error; err != nil {
		return fmt.Errorf("failed to clear achievements: %w", err)
	}

	log.Info().Msg("Clearing streaks...")
	if err := db.Exec("DELETE FROM habit_streaks").Error; err != nil {
		return fmt.Errorf("failed to clear streaks: %w", err)
	}

	log.Info().Msg("Clearing habits...")
	if err := db.Exec("DELETE FROM habits").Error; err != nil {
		return fmt.Errorf("failed to clear habits: %w", err)
	}

	log.Info().Msg("Clearing users...")
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	tables := []string{"users", "habits", "habit_streaks", "habit_checkins", "achievements"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)).Error; err != nil {
			log.Warn().Err(err).Str("table", table).Msg("Failed to reset sequence (this is normal for new databases)")
		}
	}

	return nil
}
