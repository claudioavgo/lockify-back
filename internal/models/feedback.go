package models

import (
	"time"

	"gorm.io/gorm"
)

type Feedback struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Rating    int            `gorm:"not null" json:"rating" validate:"required,min=1,max=5"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"not null" json:"user_id"`
}
