package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Achievement represents a user achievement or milestone
type Achievement struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"not null;index"`
	HabitID         uint           `json:"habit_id" gorm:"not null;index"`
	AchievementType string         `json:"achievement_type" gorm:"not null"` // 'streak_completed', 'max_streak', etc.
	TargetDays      int            `json:"target_days" gorm:"not null"`
	AchievedAt      time.Time      `json:"achieved_at" gorm:"autoCreateTime"`
	Metadata        datatypes.JSON `json:"metadata"` // Additional data like streak_id, etc.
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	User  User  `json:"-" gorm:"foreignKey:UserID"`
	Habit Habit `json:"-" gorm:"foreignKey:HabitID"`
}

// TableName specifies the table name for the Achievement model
func (Achievement) TableName() string {
	return "achievements"
}

// AchievementResponse is the DTO for achievement data sent to clients
type AchievementResponse struct {
	ID              uint           `json:"id"`
	UserID          uint           `json:"user_id"`
	HabitID         uint           `json:"habit_id"`
	AchievementType string         `json:"achievement_type"`
	TargetDays      int            `json:"target_days"`
	AchievedAt      time.Time      `json:"achieved_at"`
	Metadata        datatypes.JSON `json:"metadata,omitempty"`
	
	// Optional related data
	HabitName string `json:"habit_name,omitempty"`
}

// ToResponse converts an Achievement to an AchievementResponse
func (a *Achievement) ToResponse() AchievementResponse {
	return AchievementResponse{
		ID:              a.ID,
		UserID:          a.UserID,
		HabitID:         a.HabitID,
		AchievementType: a.AchievementType,
		TargetDays:      a.TargetDays,
		AchievedAt:      a.AchievedAt,
		Metadata:        a.Metadata,
	}
}
