package handler

import (
	"net/http"

	"backend_keping/internal/domain"
	"backend_keping/internal/usecase"
	"backend_keping/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	wishUsecase usecase.WishlistUsecase
}

func NewWishlistHandler(wishUsecase usecase.WishlistUsecase) *WishlistHandler {
	return &WishlistHandler{wishUsecase: wishUsecase}
}

func (h *WishlistHandler) Create(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userID := userIDVal.(uuid.UUID)

	var wishlist domain.Wishlist
	if err := c.ShouldBindJSON(&wishlist); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if err := h.wishUsecase.CreateWishlist(c.Request.Context(), userID, &wishlist); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create wishlist: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Wishlist created successfully", wishlist)
}

func (h *WishlistHandler) GetWishlists(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userID := userIDVal.(uuid.UUID)

	wishlists, err := h.wishUsecase.GetUserWishlists(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve wishlists: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Wishlists retrieved successfully", wishlists)
}
