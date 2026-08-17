package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionSource string
type TransactionType string

const (
	SourceNotification TransactionSource = "NOTIFICATION"
	SourceOCRScan      TransactionSource = "OCR_SCAN"
	SourceManual       TransactionSource = "MANUAL"

	TypeCredit TransactionType = "CREDIT" // Uang Masuk
	TypeDebit  TransactionType = "DEBIT"  // Uang Keluar
)

// Transaction entity model
type Transaction struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID         `gorm:"type:uuid;not null;index" json:"user_id"`
	Source      TransactionSource `gorm:"type:varchar(20);not null" json:"source"`
	Type        TransactionType   `gorm:"type:varchar(10);not null" json:"type"`
	Amount      float64           `gorm:"type:numeric(15,2);not null" json:"amount"`
	Description string            `gorm:"type:varchar(255)" json:"description"`
	RawText     string            `gorm:"type:text" json:"raw_text,omitempty"`
	ImageURL    string            `gorm:"type:text" json:"image_url,omitempty"`
	Timestamp   time.Time         `gorm:"autoCreateTime" json:"timestamp"`
}

func (Transaction) TableName() string {
	return "transactions"
}
