// Package models contains all the shared types used for API input, output, and validation.
package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Name         string         `json:"name" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserResponse is the DTO for user data sent to clients
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ConsistencyDataPoint struct {
	Date        time.Time `json:"date"`
	Percentage  float64   `json:"percentage"`
	CheckIns    int       `json:"check_ins"`
	TotalHabits int       `json:"total_habits"`
}

type StreakInsight struct {
	CurrentLongestStreak int     `json:"current_longest_streak"`
	BestStreakEver       int     `json:"best_streak_ever"`
	AverageStreakLength  float64 `json:"average_streak_length"`
	ActiveStreaksCount   int     `json:"active_streaks_count"`
}

type HabitPerformance struct {
	HabitID         uint       `json:"habit_id"`
	HabitName       string     `json:"habit_name"`
	ConsistencyRate float64    `json:"consistency_rate"`
	CurrentStreak   int        `json:"current_streak"`
	TotalCheckIns   int        `json:"total_check_ins"`
	LastCheckIn     *time.Time `json:"last_check_in,omitempty"`
}

type OverviewStats struct {
	TotalHabits        int     `json:"total_habits"`
	ActiveHabits       int     `json:"active_habits"`
	TotalCheckIns      int     `json:"total_check_ins"`
	TotalAchievements  int     `json:"total_achievements"`
	DaysSinceJoined    int     `json:"days_since_joined"`
	OverallConsistency float64 `json:"overall_consistency"`
	WeeklyConsistency  float64 `json:"weekly_consistency"`
	MonthlyConsistency float64 `json:"monthly_consistency"`
}

type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Overview           OverviewStats          `json:"overview"`
	StreakInsights     StreakInsight          `json:"streak_insights"`
	ConsistencyChart   []ConsistencyDataPoint `json:"consistency_chart"`
	TopHabits          []HabitPerformance     `json:"top_habits"`
	RecentAchievements []AchievementResponse  `json:"recent_achievements"`

	MostConsistentHabit *HabitPerformance    `json:"most_consistent_habit,omitempty"`
	ImprovementTrend    string               `json:"improvement_trend"` // "improving", "declining", "stable"
	NextMilestone       *AchievementResponse `json:"next_milestone,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}
