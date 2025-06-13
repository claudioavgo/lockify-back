package models

import (
	"time"

	"gorm.io/gorm"
)

type Habit struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name        string `gorm:"not null" json:"name" validate:"required,min=2"`
	Description string `gorm:"not null" json:"description" validate:"min=2"`
	DaysOfWeek  string `gorm:"not null" json:"days_of_week" validate:"required,min=1"`
	StartsAt    string `gorm:"not null" json:"starts_at" validate:"required,min=2"`
	EndsAt      string `gorm:"not null" json:"ends_at" validate:"min=2"`
	NoEndDate   bool   `gorm:"default:false" json:"no_end_date"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"not null" json:"user_id"`
}

type HabitDayCheckIn struct {
	ID uint `gorm:"primaryKey" json:"id"`

	HabitID uint  `gorm:"not null" json:"habit_id"`
	Habit   Habit `gorm:"foreignKey:HabitID" json:"habit"`

	Date time.Time `gorm:"not null" json:"date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"not null" json:"user_id"`
}
