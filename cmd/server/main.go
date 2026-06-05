package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/config"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/database"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/handler"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/middleware"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/repository"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	ctx := context.Background()
	pool, err := database.NewPostgresPool(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	orgRepo := repository.NewOrganizationRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	apiKeyRepo := repository.NewAPIKeyRepository(pool)
	usageRepo := repository.NewUsageRepository(pool)

	// Initialize services
	authService := service.NewAuthService(orgRepo, userRepo, apiKeyRepo, cfg.Auth.JWTSecret, cfg.Auth.KeyExpireDays)
	proxyService := service.NewProxyService(
		usageRepo,
		cfg.AIProviders.OpenAI.BaseURL,
		cfg.AIProviders.OpenAI.APIKey,
		cfg.AIProviders.Anthropic.BaseURL,
		cfg.AIProviders.Anthropic.APIKey,
	)
	usageService := service.NewUsageService(usageRepo, userRepo)
	billingService := service.NewBillingService(usageRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	proxyHandler := handler.NewProxyHandler(proxyService)
	usageHandler := handler.NewUsageHandler(usageService, billingService)
	orgHandler := handler.NewOrgHandler(orgRepo, userRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.JWTSecret)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Setup router
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected auth routes
		authProtected := v1.Group("/auth")
		authProtected.Use(authMiddleware.JWTAuth())
		{
			authProtected.POST("/keys", authHandler.GenerateAPIKey)
			authProtected.DELETE("/keys/:key_id", authHandler.RevokeAPIKey)
			authProtected.GET("/keys", authHandler.ListAPIKeys)
		}

		// Organization routes (protected)
		org := v1.Group("/org")
		org.Use(authMiddleware.JWTAuth())
		{
			org.GET("", orgHandler.GetOrg)
			org.PUT("", orgHandler.UpdateOrg)
			org.POST("/members", orgHandler.AddMember)
			org.DELETE("/members/:user_id", orgHandler.RemoveMember)
			org.GET("/members", orgHandler.ListMembers)
		}

		// Proxy routes (protected, with rate limiting)
		proxy := v1.Group("/proxy")
		proxy.Use(authMiddleware.JWTAuth())
		proxy.Use(rateLimiter.Limit(func(c *gin.Context) string {
			if userID, exists := c.Get("user_id"); exists {
				return userID.(string)
			}
			return ""
		}))
		{
			proxy.POST("/openai/*path", proxyHandler.ProxyOpenAI)
			proxy.POST("/anthropic/*path", proxyHandler.ProxyAnthropic)
		}

		// Usage routes (protected)
		usage := v1.Group("/usage")
		usage.Use(authMiddleware.JWTAuth())
		{
			usage.GET("", usageHandler.GetCurrentUsage)
			usage.GET("/history", usageHandler.GetUsageHistory)
		}

		// Billing routes (protected)
		billing := v1.Group("/billing")
		billing.Use(authMiddleware.JWTAuth())
		{
			billing.GET("", usageHandler.GetBilling)
			billing.GET("/export", usageHandler.ExportBilling)
		}
	}

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
