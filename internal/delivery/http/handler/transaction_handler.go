package handler

import (
	"net/http"
	"strconv"

	"backend_keping/internal/usecase"
	"backend_keping/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	txUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(txUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{txUsecase: txUsecase}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userID := userIDVal.(uuid.UUID)

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	transactions, err := h.txUsecase.GetUserTransactions(c.Request.Context(), userID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve transactions: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}
