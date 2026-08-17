package repository

import (
	"context"

	"backend_keping/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WishlistRepository interface {
	Create(ctx context.Context, wishlist *domain.Wishlist) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error)
	GetActiveByUserIDWithLock(ctx context.Context, tx *gorm.DB, userID uuid.UUID) ([]domain.Wishlist, error)
	UpdateSavedAmount(ctx context.Context, tx *gorm.DB, wishlistID uuid.UUID, addedAmount float64) error
	GetDB() *gorm.DB
}

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *wishlistRepository) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	return r.db.WithContext(ctx).Create(wishlist).Error
}

func (r *wishlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error) {
	var wishlists []domain.Wishlist
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&wishlists).Error
	return wishlists, err
}

func (r *wishlistRepository) GetActiveByUserIDWithLock(ctx context.Context, tx *gorm.DB, userID uuid.UUID) ([]domain.Wishlist, error) {
	var wishlists []domain.Wishlist
	db := r.db
	if tx != nil {
		db = tx
	}

	// Clauses Locking (FOR UPDATE) to prevent race conditions during auto-save calculation
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND is_completed = ?", userID, false).
		Find(&wishlists).Error

	return wishlists, err
}

func (r *wishlistRepository) UpdateSavedAmount(ctx context.Context, tx *gorm.DB, wishlistID uuid.UUID, addedAmount float64) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	var wishlist domain.Wishlist
	if err := db.WithContext(ctx).First(&wishlist, "id = ?", wishlistID).Error; err != nil {
		return err
	}

	newSaved := wishlist.SavedAmount + addedAmount
	isCompleted := newSaved >= wishlist.TargetAmount

	return db.WithContext(ctx).Model(&wishlist).Updates(map[string]interface{}{
		"saved_amount": newSaved,
		"is_completed": isCompleted,
	}).Error
}
