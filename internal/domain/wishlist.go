package domain

import (
	"github.com/google/uuid"
)

// Wishlist entity model
type Wishlist struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name               string    `gorm:"type:varchar(100);not null" json:"name"`
	TargetAmount       float64   `gorm:"type:numeric(15,2);not null" json:"target_amount"`
	SavedAmount        float64   `gorm:"type:numeric(15,2);default:0" json:"saved_amount"`
	AutoSavePercentage float64   `gorm:"type:numeric(5,2);default:0" json:"auto_save_percentage"`
	IconName           string    `gorm:"type:varchar(50);default:'star'" json:"icon_name"`
	LabelColor         string    `gorm:"type:varchar(10);default:'#2A9D8F'" json:"label_color"`
	IsCompleted        bool      `gorm:"type:boolean;default:false" json:"is_completed"`
}

func (Wishlist) TableName() string {
	return "wishlists"
}
