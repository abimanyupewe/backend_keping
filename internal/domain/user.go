package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user extension table in public schema
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}
