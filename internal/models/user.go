package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name     string `gorm:"not null" json:"name" validate:"required,min=2"`
	Email    string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"password" validate:"required,min=6"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserResponse struct {
	ID uint `json:"id"`

	Name  string `json:"name"`
	Email string `json:"email"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
