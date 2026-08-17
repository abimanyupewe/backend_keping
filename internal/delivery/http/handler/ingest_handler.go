package handler

import (
	"io"
	"net/http"

	"backend_keping/internal/domain"
	"backend_keping/internal/usecase"
	"backend_keping/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IngestHandler struct {
	ingestUsecase usecase.IngestUsecase
}

func NewIngestHandler(ingestUsecase usecase.IngestUsecase) *IngestHandler {
	return &IngestHandler{ingestUsecase: ingestUsecase}
}

func (h *IngestHandler) HandleNotification(c *gin.Context) {
	var payload domain.NotificationIngestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userID := userIDVal.(uuid.UUID)

	tx, err := h.ingestUsecase.ProcessNotification(c.Request.Context(), userID, payload)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process notification: "+err.Error())
		return
	}

	if tx == nil {
		response.Success(c, http.StatusOK, "Notification processed (no transaction recorded)", nil)
		return
	}

	response.Success(c, http.StatusCreated, "Notification ingested successfully", tx)
}

func (h *IngestHandler) HandleOCRScan(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User context missing")
		return
	}
	userID := userIDVal.(uuid.UUID)

	fileHeader, err := c.FormFile("receipt")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "File 'receipt' form field is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to open uploaded file")
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read uploaded image bytes")
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	tx, err := h.ingestUsecase.ProcessOCRScan(c.Request.Context(), userID, imageBytes, mimeType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to scan receipt via Gemini AI: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Receipt scanned and transaction saved successfully", tx)
}
