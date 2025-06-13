package models

import (
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name         string  `gorm:"not null" json:"name" validate:"required,min=2"`
	Description  string  `gorm:"not null" json:"description" validate:"required,min=2"`
	VideoURL     string  `gorm:"not null" json:"video_url" validate:"required,url"`
	ImageURL     string  `gorm:"not null" json:"image_url" validate:"required,url"`
	MuscleGroup  string  `gorm:"not null" json:"muscle_group" validate:"required,min=2"`
	Difficulty   string  `gorm:"not null" json:"difficulty" validate:"required,min=2"`
	Instructions string  `gorm:"not null" json:"instructions" validate:"required,min=2"`
	Sets         int     `gorm:"not null" json:"sets" validate:"required,min=1"`
	Reps         int     `gorm:"not null" json:"reps" validate:"required,min=1"`
	Weight       float64 `gorm:"not null" json:"weight" validate:"required,min=0"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"not null" json:"user_id"`
}
