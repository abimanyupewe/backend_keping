package usecase

import (
	"context"

	"backend_keping/internal/domain"
	"backend_keping/internal/repository"

	"github.com/google/uuid"
)

type TransactionUsecase interface {
	GetUserTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Transaction, error)
}

type transactionUsecase struct {
	txRepo repository.TransactionRepository
}

func NewTransactionUsecase(txRepo repository.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{txRepo: txRepo}
}

func (u *transactionUsecase) GetUserTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Transaction, error) {
	return u.txRepo.GetByUserID(ctx, userID, limit)
}
