package http

import (
	"backend_keping/internal/delivery/http/handler"
	"backend_keping/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RouterConfig struct {
	Engine             *gin.Engine
	HealthHandler      *handler.HealthHandler
	AuthHandler        *handler.AuthHandler
	IngestHandler      *handler.IngestHandler
	TxHandler          *handler.TransactionHandler
	WishlistHandler    *handler.WishlistHandler
	SupabaseJWTSecret  string
	IngestionSecretKey string
	RateLimitRPS       int
}

func SetupRouter(cfg RouterConfig) {
	// Rate Limiter Middleware
	rps := cfg.RateLimitRPS
	if rps <= 0 {
		rps = 5
	}
	limiter := middleware.NewIPRateLimiter(rate.Limit(rps), rps*2)
	cfg.Engine.Use(middleware.RateLimiterMiddleware(limiter))

	// CORS / Recovery
	cfg.Engine.Use(gin.Recovery())

	// API v1 Group
	v1 := cfg.Engine.Group("/api/v1")
	{
		// Infrastructure / Health Check Endpoint (Public)
		v1.GET("/health", cfg.HealthHandler.Check)

		// Public Auth Endpoints (Register & Login)
		if cfg.AuthHandler != nil {
			authPublic := v1.Group("/auth")
			{
				authPublic.POST("/register", cfg.AuthHandler.Register)
				authPublic.POST("/login", cfg.AuthHandler.Login)
			}
		}

		// Notification Ingest Endpoint (Protected via INGESTION_SECRET_KEY & Supabase JWT)
		ingest := v1.Group("/ingest")
		ingest.Use(middleware.IngestionSecretMiddleware(cfg.IngestionSecretKey))
		ingest.Use(middleware.AuthMiddleware(cfg.SupabaseJWTSecret))
		{
			ingest.POST("/notification", cfg.IngestHandler.HandleNotification)
			ingest.POST("/scan", cfg.IngestHandler.HandleOCRScan)
		}

		// Authenticated Routes (Protected via Supabase JWT Auth)
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(cfg.SupabaseJWTSecret))
		{
			// Transactions
			auth.GET("/transactions", cfg.TxHandler.GetTransactions)

			// Wishlists
			auth.POST("/wishlists", cfg.WishlistHandler.Create)
			auth.GET("/wishlists", cfg.WishlistHandler.GetWishlists)
		}
	}
}
