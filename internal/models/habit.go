package models

import (
	"time"

	"gorm.io/gorm"
)

// Habit represents a habit that a user wants to track
type Habit struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Color       string         `json:"color" gorm:"size:7"` // Hex color code
	Icon        string         `json:"icon" gorm:"size:50"` // Icon identifier
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User    User          `json:"-" gorm:"foreignKey:UserID"`
	Streaks []HabitStreak `json:"streaks,omitempty" gorm:"foreignKey:HabitID"`
}

// TableName specifies the table name for the Habit model
func (Habit) TableName() string {
	return "habits"
}

// HabitResponse is the DTO for habit data sent to clients
type HabitResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Optional related data
	CurrentStreak *HabitStreakResponse `json:"current_streak,omitempty"`
}

// ToResponse converts a Habit to a HabitResponse
func (h *Habit) ToResponse() HabitResponse {
	return HabitResponse{
		ID:          h.ID,
		UserID:      h.UserID,
		Name:        h.Name,
		Description: h.Description,
		Color:       h.Color,
		Icon:        h.Icon,
		IsActive:    h.IsActive,
		Status:      "inactive", // Default to inactive, will be updated by service
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
	}
}

func (h *Habit) ToResponseWithStreak(currentStreak *HabitStreak) HabitResponse {
	response := h.ToResponse()

	if currentStreak != nil && currentStreak.Status == "active" && currentStreak.CurrentStreak > 0 {
		response.Status = "active"
		streakResponse := currentStreak.ToResponse()
		response.CurrentStreak = &streakResponse
	} else {
		response.Status = "inactive"
		if currentStreak != nil {
			streakResponse := currentStreak.ToResponse()
			response.CurrentStreak = &streakResponse
		}
	}

	return response
}
