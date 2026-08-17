package handler

import (
	"net/http"

	"backend_keping/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Database connection error")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Error(c, http.StatusServiceUnavailable, "Database ping failed")
		return
	}

	response.Success(c, http.StatusOK, "Keping Backend API & Database is Healthy", gin.H{
		"status":   "UP",
		"database": "CONNECTED",
	})
}
