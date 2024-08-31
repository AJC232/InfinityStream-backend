package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username  string    `gorm:"size:100;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	Email     string    `gorm:"size:100;not null" json:"email" validate:"email" unique:"true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserResponse represents a user response
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
