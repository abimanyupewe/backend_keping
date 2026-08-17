package usecase

import (
	"context"

	"backend_keping/internal/domain"
	"backend_keping/internal/repository"
	"backend_keping/pkg/ai"
	"backend_keping/pkg/logger"
	"backend_keping/pkg/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IngestUsecase interface {
	ProcessNotification(ctx context.Context, userID uuid.UUID, payload domain.NotificationIngestPayload) (*domain.Transaction, error)
	ProcessOCRScan(ctx context.Context, userID uuid.UUID, imageBytes []byte, mimeType string) (*domain.Transaction, error)
}

type ingestUsecase struct {
	txRepo     repository.TransactionRepository
	wishRepo   repository.WishlistRepository
	parser     RegexParser
	calculator AutoSaveCalculator
	aiClient   ai.GeminiClient
	cldService storage.CloudinaryService
}

func NewIngestUsecase(
	txRepo repository.TransactionRepository,
	wishRepo repository.WishlistRepository,
	parser RegexParser,
	calculator AutoSaveCalculator,
	aiClient ai.GeminiClient,
	cldService storage.CloudinaryService,
) IngestUsecase {
	return &ingestUsecase{
		txRepo:     txRepo,
		wishRepo:   wishRepo,
		parser:     parser,
		calculator: calculator,
		aiClient:   aiClient,
		cldService: cldService,
	}
}

func (u *ingestUsecase) ProcessNotification(ctx context.Context, userID uuid.UUID, payload domain.NotificationIngestPayload) (*domain.Transaction, error) {
	parsed := u.parser.Parse(payload.PackageName, payload.Title, payload.Content)
	if !parsed.IsValid {
		logger.Info("Ingest notification ignored (not a valid transaction): %s", payload.Content)
		return nil, nil
	}

	txModel := &domain.Transaction{
		UserID:      userID,
		Source:      domain.SourceNotification,
		Type:        parsed.Type,
		Amount:      parsed.Amount,
		Description: parsed.Description,
		RawText:     payload.Content,
	}

	// Begin DB Transaction for atomic transaction save + wishlist auto-save updates
	db := u.wishRepo.GetDB()
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.txRepo.Create(ctx, tx, txModel); err != nil {
			return err
		}

		// Trigger auto-save allocation if transaction is CREDIT (Income)
		if parsed.Type == domain.TypeCredit {
			u.triggerAutoSave(ctx, tx, userID, parsed.Amount)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return txModel, nil
}

func (u *ingestUsecase) ProcessOCRScan(ctx context.Context, userID uuid.UUID, imageBytes []byte, mimeType string) (*domain.Transaction, error) {
	// 1. Scan image via Gemini AI
	ocrRes, err := u.aiClient.ScanReceipt(ctx, imageBytes, mimeType)
	if err != nil {
		return nil, err
	}

	// 2. Upload image to Cloudinary storage if configured
	imageURL, errUpload := u.cldService.UploadReceipt(ctx, imageBytes, uuid.New().String())
	if errUpload != nil {
		logger.Error("Failed to upload receipt to Cloudinary: %v", errUpload)
	}

	txType := domain.TransactionType(ocrRes.Type)
	if txType != domain.TypeCredit && txType != domain.TypeDebit {
		txType = domain.TypeDebit
	}

	txModel := &domain.Transaction{
		UserID:      userID,
		Source:      domain.SourceOCRScan,
		Type:        txType,
		Amount:      ocrRes.Amount,
		Description: ocrRes.Description,
		ImageURL:    imageURL,
	}

	db := u.wishRepo.GetDB()
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.txRepo.Create(ctx, tx, txModel); err != nil {
			return err
		}

		if txType == domain.TypeCredit {
			u.triggerAutoSave(ctx, tx, userID, ocrRes.Amount)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return txModel, nil
}

func (u *ingestUsecase) triggerAutoSave(ctx context.Context, tx *gorm.DB, userID uuid.UUID, incomeAmount float64) {
	wishlists, err := u.wishRepo.GetActiveByUserIDWithLock(ctx, tx, userID)
	if err != nil {
		logger.Error("Failed to fetch active wishlists for auto-save: %v", err)
		return
	}

	for _, w := range wishlists {
		allocation := u.calculator.CalculateAllocation(incomeAmount, w)
		if allocation > 0 {
			if err := u.wishRepo.UpdateSavedAmount(ctx, tx, w.ID, allocation); err != nil {
				logger.Error("Failed to update saved amount for wishlist %s: %v", w.ID, err)
			} else {
				logger.Info("Auto-allocated Rp%.2f to wishlist %s (%s)", allocation, w.Name, w.ID)
			}
		}
	}
}
