package main

import (
	"log"
	"os"
	"strconv"

	deliveryHTTP "backend_keping/internal/delivery/http"
	"backend_keping/internal/delivery/http/handler"
	"backend_keping/internal/domain"
	"backend_keping/internal/repository"
	"backend_keping/internal/usecase"
	"backend_keping/pkg/ai"
	"backend_keping/pkg/logger"
	"backend_keping/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env if present
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system environment variables")
	}

	port := getEnv("PORT", "8080")

	dbURL := getEnv("DB_URL", os.Getenv("DATABASE_URL"))
	jwtSecret := getEnv("SUPABASE_JWT_SECRET", os.Getenv("JWT_SECRET"))
	ingestionSecret := getEnv("INGESTION_SECRET_KEY", "keping-super-secret-ingestion-key-2026")

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	geminiModel := getEnv("GEMINI_MODEL", "gemini-1.5-flash")

	cldCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cldAPIKey := os.Getenv("CLOUDINARY_API_KEY")
	cldAPISecret := os.Getenv("CLOUDINARY_API_SECRET")
	cldFolder := getEnv("CLOUDINARY_FOLDER", "receipts")

	rateLimitRPSStr := getEnv("RATE_LIMIT_RPS", "5")
	rateLimitRPS, _ := strconv.Atoi(rateLimitRPSStr)

	// Database Connection (Supabase PostgreSQL)
	var db *gorm.DB
	var err error
	if dbURL != "" {
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			logger.Error("Failed to connect to Supabase PostgreSQL: %v", err)
		} else {
			logger.Info("Successfully connected to Supabase PostgreSQL Database")
			if err := db.AutoMigrate(&domain.User{}, &domain.Transaction{}, &domain.Wishlist{}); err != nil {
				logger.Error("Failed to auto-migrate database schemas: %v", err)
			}
		}
	} else {
		logger.Info("DB_URL/DATABASE_URL not set. Running in dry-run mode.")
	}

	// Initialize Storage Service (Cloudinary)
	cldService, err := storage.NewCloudinaryService(cldCloudName, cldAPIKey, cldAPISecret, cldFolder)
	if err != nil {
		logger.Error("Failed to initialize Cloudinary storage: %v", err)
	} else {
		logger.Info("Successfully initialized Cloudinary Storage Service (Folder: %s)", cldFolder)
	}

	// Initialize AI Client (Google Gemini)
	geminiClient := ai.NewGeminiClient(geminiAPIKey, geminiModel)

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	wishRepo := repository.NewWishlistRepository(db)

	// Initialize Business Engine & Usecases
	regexParser := usecase.NewRegexParser()
	autoSaveCalc := usecase.NewAutoSaveCalculator()

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtSecret)
	ingestUsecase := usecase.NewIngestUsecase(txRepo, wishRepo, regexParser, autoSaveCalc, geminiClient, cldService)
	txUsecase := usecase.NewTransactionUsecase(txRepo)
	wishUsecase := usecase.NewWishlistUsecase(wishRepo)

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authUsecase)
	ingestHandler := handler.NewIngestHandler(ingestUsecase)
	txHandler := handler.NewTransactionHandler(txUsecase)
	wishHandler := handler.NewWishlistHandler(wishUsecase)

	// Setup Gin Engine
	r := gin.Default()

	deliveryHTTP.SetupRouter(deliveryHTTP.RouterConfig{
		Engine:             r,
		HealthHandler:      healthHandler,
		AuthHandler:        authHandler,
		IngestHandler:      ingestHandler,
		TxHandler:          txHandler,
		WishlistHandler:    wishHandler,
		SupabaseJWTSecret:  jwtSecret,
		IngestionSecretKey: ingestionSecret,
		RateLimitRPS:       rateLimitRPS,
	})

	logger.Info("Keping Backend Server starting on port :%s ...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
