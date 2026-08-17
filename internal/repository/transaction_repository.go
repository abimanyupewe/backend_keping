package repository

import (
	"context"

	"backend_keping/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *domain.Transaction) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *domain.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&transactions).Error
	return transactions, err
}
