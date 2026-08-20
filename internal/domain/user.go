package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user extension table in public schema
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name         string    `gorm:"type:varchar(255)" json:"name"`
	PasswordHash string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
