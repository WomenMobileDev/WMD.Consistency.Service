package service

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/repository"
)

type UserService struct {
	userRepo        repository.UserRepository
	habitRepo       repository.HabitRepository
	streakRepo      repository.StreakRepository
	checkInRepo     repository.CheckInRepository
	achievementRepo repository.AchievementRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	habitRepo repository.HabitRepository,
	streakRepo repository.StreakRepository,
	checkInRepo repository.CheckInRepository,
	achievementRepo repository.AchievementRepository,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		habitRepo:       habitRepo,
		streakRepo:      streakRepo,
		checkInRepo:     checkInRepo,
		achievementRepo: achievementRepo,
	}
}

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*models.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	habits, err := s.habitRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	achievements, err := s.achievementRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	overview := s.calculateOverviewStats(user, habits, achievements)
	streakInsights := s.calculateStreakInsights(ctx, habits)
	consistencyChart := s.calculateConsistencyChart(ctx, habits, 30) // Last 30 days
	topHabits := s.calculateTopHabits(ctx, habits)
	recentAchievements := s.getRecentAchievements(achievements, 5) // Last 5 achievements

	mostConsistentHabit := s.findMostConsistentHabit(topHabits)
	improvementTrend := s.calculateImprovementTrend(consistencyChart)
	nextMilestone := s.predictNextMilestone(ctx, habits, achievements)

	profile := &models.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,

		Overview:           overview,
		StreakInsights:     streakInsights,
		ConsistencyChart:   consistencyChart,
		TopHabits:          topHabits,
		RecentAchievements: recentAchievements,

		MostConsistentHabit: mostConsistentHabit,
		ImprovementTrend:    improvementTrend,
		NextMilestone:       nextMilestone,
	}

	return profile, nil
}

func (s *UserService) calculateOverviewStats(user *models.User, habits []models.Habit, achievements []models.Achievement) models.OverviewStats {
	now := time.Now()
	daysSinceJoined := int(now.Sub(user.CreatedAt).Hours() / 24)

	activeHabits := 0
	for _, habit := range habits {
		if habit.IsActive {
			activeHabits++
		}
	}

	// Calculate total check-ins across all habits
	totalCheckIns := s.calculateTotalCheckIns(context.Background(), habits)

	// Calculate consistency percentages
	overallConsistency := s.calculateConsistencyForPeriod(context.Background(), habits, 0) // All time
	weeklyConsistency := s.calculateConsistencyForPeriod(context.Background(), habits, 7)
	monthlyConsistency := s.calculateConsistencyForPeriod(context.Background(), habits, 30)

	return models.OverviewStats{
		TotalHabits:        len(habits),
		ActiveHabits:       activeHabits,
		TotalCheckIns:      totalCheckIns,
		TotalAchievements:  len(achievements),
		DaysSinceJoined:    daysSinceJoined,
		OverallConsistency: overallConsistency,
		WeeklyConsistency:  weeklyConsistency,
		MonthlyConsistency: monthlyConsistency,
	}
}

func (s *UserService) calculateStreakInsights(ctx context.Context, habits []models.Habit) models.StreakInsight {
	var currentLongest, bestEver, totalStreaks, activeCount int
	var streakLengths []int

	for _, habit := range habits {
		streaks, err := s.streakRepo.FindByHabitID(ctx, habit.ID)
		if err != nil {
			continue
		}

		for _, streak := range streaks {
			if streak.Status == "active" {
				activeCount++
				if streak.CurrentStreak > currentLongest {
					currentLongest = streak.CurrentStreak
				}
			}

			if streak.MaxStreakAchieved > bestEver {
				bestEver = streak.MaxStreakAchieved
			}

			if streak.MaxStreakAchieved > 0 {
				streakLengths = append(streakLengths, streak.MaxStreakAchieved)
				totalStreaks++
			}
		}
	}

	var avgStreak float64
	if len(streakLengths) > 0 {
		sum := 0
		for _, length := range streakLengths {
			sum += length
		}
		avgStreak = float64(sum) / float64(len(streakLengths))
	}

	return models.StreakInsight{
		CurrentLongestStreak: currentLongest,
		BestStreakEver:       bestEver,
		AverageStreakLength:  math.Round(avgStreak*100) / 100, // Round to 2 decimal places
		ActiveStreaksCount:   activeCount,
	}
}

func (s *UserService) calculateConsistencyChart(ctx context.Context, habits []models.Habit, days int) []models.ConsistencyDataPoint {
	var dataPoints []models.ConsistencyDataPoint
	now := time.Now().UTC().Truncate(24 * time.Hour)

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		totalCheckIns := 0
		activeHabits := 0

		for _, habit := range habits {
			if habit.IsActive && habit.CreatedAt.Before(date.Add(24*time.Hour)) {
				activeHabits++

				streaks, err := s.streakRepo.FindByHabitID(ctx, habit.ID)
				if err != nil {
					continue
				}

				for _, streak := range streaks {
					if checkIn, err := s.checkInRepo.FindByDate(ctx, streak.ID, dateStr); err == nil && checkIn != nil {
						totalCheckIns++
						break
					}
				}
			}
		}

		var percentage float64
		if activeHabits > 0 {
			percentage = (float64(totalCheckIns) / float64(activeHabits)) * 100
		}

		dataPoints = append(dataPoints, models.ConsistencyDataPoint{
			Date:        date,
			Percentage:  math.Round(percentage*100) / 100,
			CheckIns:    totalCheckIns,
			TotalHabits: activeHabits,
		})
	}

	return dataPoints
}

func (s *UserService) calculateTopHabits(ctx context.Context, habits []models.Habit) []models.HabitPerformance {
	var performances []models.HabitPerformance

	for _, habit := range habits {
		performance := s.calculateHabitPerformance(ctx, habit)
		performances = append(performances, performance)
	}

	sort.Slice(performances, func(i, j int) bool {
		return performances[i].ConsistencyRate > performances[j].ConsistencyRate
	})

	if len(performances) > 5 {
		performances = performances[:5]
	}

	return performances
}

func (s *UserService) calculateHabitPerformance(ctx context.Context, habit models.Habit) models.HabitPerformance {
	streaks, err := s.streakRepo.FindByHabitID(ctx, habit.ID)
	if err != nil {
		return models.HabitPerformance{
			HabitID:   habit.ID,
			HabitName: habit.Name,
		}
	}

	var totalCheckIns, currentStreak int
	var lastCheckIn *time.Time

	for _, streak := range streaks {
		if streak.Status == "active" {
			currentStreak = streak.CurrentStreak
			if streak.LastCheckInDate != nil {
				lastCheckIn = streak.LastCheckInDate
			}
		}

		checkIns, err := s.checkInRepo.FindByStreakID(ctx, streak.ID)
		if err == nil {
			totalCheckIns += len(checkIns)
		}
	}

	daysSinceCreated := int(time.Since(habit.CreatedAt).Hours() / 24)
	if daysSinceCreated == 0 {
		daysSinceCreated = 1
	}

	consistencyRate := (float64(totalCheckIns) / float64(daysSinceCreated)) * 100
	if consistencyRate > 100 {
		consistencyRate = 100
	}

	return models.HabitPerformance{
		HabitID:         habit.ID,
		HabitName:       habit.Name,
		ConsistencyRate: math.Round(consistencyRate*100) / 100,
		CurrentStreak:   currentStreak,
		TotalCheckIns:   totalCheckIns,
		LastCheckIn:     lastCheckIn,
	}
}

func (s *UserService) calculateTotalCheckIns(ctx context.Context, habits []models.Habit) int {
	total := 0
	for _, habit := range habits {
		streaks, err := s.streakRepo.FindByHabitID(ctx, habit.ID)
		if err != nil {
			continue
		}
		for _, streak := range streaks {
			checkIns, err := s.checkInRepo.FindByStreakID(ctx, streak.ID)
			if err == nil {
				total += len(checkIns)
			}
		}
	}
	return total
}

func (s *UserService) calculateConsistencyForPeriod(ctx context.Context, habits []models.Habit, days int) float64 {
	if len(habits) == 0 {
		return 0
	}

	now := time.Now().UTC()
	var totalPossible, totalActual int

	for _, habit := range habits {
		if !habit.IsActive {
			continue
		}

		var periodStart time.Time
		if days > 0 {
			periodStart = now.AddDate(0, 0, -days)
		} else {
			periodStart = habit.CreatedAt
		}

		if habit.CreatedAt.After(periodStart) {
			periodStart = habit.CreatedAt
		}

		daysInPeriod := int(now.Sub(periodStart).Hours() / 24)
		if daysInPeriod <= 0 {
			continue
		}

		totalPossible += daysInPeriod

		streaks, err := s.streakRepo.FindByHabitID(ctx, habit.ID)
		if err != nil {
			continue
		}

		for _, streak := range streaks {
			checkIns, err := s.checkInRepo.FindByStreakID(ctx, streak.ID)
			if err != nil {
				continue
			}

			for _, checkIn := range checkIns {
				if checkIn.CheckInDate.After(periodStart) && checkIn.CheckInDate.Before(now) {
					totalActual++
				}
			}
		}
	}

	if totalPossible == 0 {
		return 0
	}

	consistency := (float64(totalActual) / float64(totalPossible)) * 100
	return math.Round(consistency*100) / 100
}

func (s *UserService) getRecentAchievements(achievements []models.Achievement, limit int) []models.AchievementResponse {
	sort.Slice(achievements, func(i, j int) bool {
		return achievements[i].AchievedAt.After(achievements[j].AchievedAt)
	})

	var recent []models.AchievementResponse
	count := limit
	if len(achievements) < limit {
		count = len(achievements)
	}

	for i := 0; i < count; i++ {
		recent = append(recent, achievements[i].ToResponse())
	}

	return recent
}

func (s *UserService) findMostConsistentHabit(habits []models.HabitPerformance) *models.HabitPerformance {
	if len(habits) == 0 {
		return nil
	}

	return &habits[0]
}

func (s *UserService) calculateImprovementTrend(chartData []models.ConsistencyDataPoint) string {
	if len(chartData) < 14 {
		return "insufficient_data"
	}

	lastWeekAvg := s.calculateAverageConsistency(chartData[len(chartData)-7:])
	prevWeekAvg := s.calculateAverageConsistency(chartData[len(chartData)-14 : len(chartData)-7])

	diff := lastWeekAvg - prevWeekAvg
	if diff > 5 {
		return "improving"
	} else if diff < -5 {
		return "declining"
	}
	return "stable"
}

func (s *UserService) calculateAverageConsistency(dataPoints []models.ConsistencyDataPoint) float64 {
	if len(dataPoints) == 0 {
		return 0
	}

	sum := 0.0
	for _, point := range dataPoints {
		sum += point.Percentage
	}
	return sum / float64(len(dataPoints))
}

func (s *UserService) predictNextMilestone(ctx context.Context, habits []models.Habit, _ []models.Achievement) *models.AchievementResponse {
	for _, habit := range habits {
		streaks, err := s.streakRepo.FindByHabitID(ctx, habit.ID)
		if err != nil {
			continue
		}

		for _, streak := range streaks {
			if streak.Status == "active" {
				if streak.CurrentStreak > 0 && streak.CurrentStreak < streak.TargetDays {
					remaining := streak.TargetDays - streak.CurrentStreak
					if remaining <= 3 {
						milestone := &models.AchievementResponse{
							HabitID:         habit.ID,
							AchievementType: "streak_completion_prediction",
							TargetDays:      streak.TargetDays,
							HabitName:       habit.Name,
						}
						return milestone
					}
				}
			}
		}
	}

	return nil
}

// UpdateProfile updates the user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*models.UserResponse, error) {
	// Find the user by ID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update the user's name
	user.Name = req.Name

	// Save the changes
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Return the updated user profile
	response := user.ToResponse()
	return &response, nil
}
