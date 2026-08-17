package usecase

import (
	"context"

	"backend_keping/internal/domain"
	"backend_keping/internal/repository"

	"github.com/google/uuid"
)

type WishlistUsecase interface {
	CreateWishlist(ctx context.Context, userID uuid.UUID, wishlist *domain.Wishlist) error
	GetUserWishlists(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error)
}

type wishlistUsecase struct {
	wishRepo repository.WishlistRepository
}

func NewWishlistUsecase(wishRepo repository.WishlistRepository) WishlistUsecase {
	return &wishlistUsecase{wishRepo: wishRepo}
}

func (u *wishlistUsecase) CreateWishlist(ctx context.Context, userID uuid.UUID, wishlist *domain.Wishlist) error {
	wishlist.UserID = userID
	if wishlist.IconName == "" {
		wishlist.IconName = "star"
	}
	if wishlist.LabelColor == "" {
		wishlist.LabelColor = "#2A9D8F"
	}
	return u.wishRepo.Create(ctx, wishlist)
}

func (u *wishlistUsecase) GetUserWishlists(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error) {
	return u.wishRepo.GetByUserID(ctx, userID)
}
