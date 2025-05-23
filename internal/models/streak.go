package models

import (
	"time"

	"gorm.io/gorm"
)

// HabitStreak represents a streak for a habit
type HabitStreak struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	HabitID           uint           `json:"habit_id" gorm:"not null;index"`
	TargetDays        int            `json:"target_days" gorm:"not null;check:target_days > 0"`
	CurrentStreak     int            `json:"current_streak" gorm:"default:0"`
	MaxStreakAchieved int            `json:"max_streak_achieved" gorm:"default:0"`
	StartDate         time.Time      `json:"start_date" gorm:"not null"`
	LastCheckInDate   *time.Time     `json:"last_check_in_date"`
	Status            string         `json:"status" gorm:"default:'active';check:status IN ('active', 'completed', 'failed')"`
	CompletedAt       *time.Time     `json:"completed_at"`
	FailedAt          *time.Time     `json:"failed_at"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Habit     Habit         `json:"-" gorm:"foreignKey:HabitID"`
	CheckIns  []HabitCheckIn `json:"check_ins,omitempty" gorm:"foreignKey:StreakID"`
}

// TableName specifies the table name for the HabitStreak model
func (HabitStreak) TableName() string {
	return "habit_streaks"
}

// HabitStreakResponse is the DTO for streak data sent to clients
type HabitStreakResponse struct {
	ID                uint       `json:"id"`
	HabitID           uint       `json:"habit_id"`
	TargetDays        int        `json:"target_days"`
	CurrentStreak     int        `json:"current_streak"`
	MaxStreakAchieved int        `json:"max_streak_achieved"`
	StartDate         time.Time  `json:"start_date"`
	LastCheckInDate   *time.Time `json:"last_check_in_date,omitempty"`
	Status            string     `json:"status"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	FailedAt          *time.Time `json:"failed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	
	// Optional related data
	CheckIns []HabitCheckInResponse `json:"check_ins,omitempty"`
}

// ToResponse converts a HabitStreak to a HabitStreakResponse
func (s *HabitStreak) ToResponse() HabitStreakResponse {
	return HabitStreakResponse{
		ID:                s.ID,
		HabitID:           s.HabitID,
		TargetDays:        s.TargetDays,
		CurrentStreak:     s.CurrentStreak,
		MaxStreakAchieved: s.MaxStreakAchieved,
		StartDate:         s.StartDate,
		LastCheckInDate:   s.LastCheckInDate,
		Status:            s.Status,
		CompletedAt:       s.CompletedAt,
		FailedAt:          s.FailedAt,
		CreatedAt:         s.CreatedAt,
	}
}

// HabitCheckIn represents a daily check-in for a habit streak
type HabitCheckIn struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	StreakID    uint           `json:"streak_id" gorm:"not null;index"`
	CheckInDate time.Time      `json:"check_in_date" gorm:"not null;index"`
	CheckedInAt time.Time      `json:"checked_in_at" gorm:"autoCreateTime"`
	Notes       string         `json:"notes"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Streak HabitStreak `json:"-" gorm:"foreignKey:StreakID"`
}

// TableName specifies the table name for the HabitCheckIn model
func (HabitCheckIn) TableName() string {
	return "habit_checkins"
}

// HabitCheckInResponse is the DTO for check-in data sent to clients
type HabitCheckInResponse struct {
	ID          uint      `json:"id"`
	StreakID    uint      `json:"streak_id"`
	CheckInDate time.Time `json:"check_in_date"`
	CheckedInAt time.Time `json:"checked_in_at"`
	Notes       string    `json:"notes,omitempty"`
}

// ToResponse converts a HabitCheckIn to a HabitCheckInResponse
func (c *HabitCheckIn) ToResponse() HabitCheckInResponse {
	return HabitCheckInResponse{
		ID:          c.ID,
		StreakID:    c.StreakID,
		CheckInDate: c.CheckInDate,
		CheckedInAt: c.CheckedInAt,
		Notes:       c.Notes,
	}
}
