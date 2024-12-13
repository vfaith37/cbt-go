package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/cbt-platform/internal/api/router"
	"github.com/yourusername/cbt-platform/internal/config"
	"github.com/yourusername/cbt-platform/internal/repository"
	"github.com/yourusername/cbt-platform/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := repository.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	examRepo := repository.NewExamRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	examService := service.NewExamService(examRepo)

	// Initialize health checker
	healthChecker := health.NewChecker(db, redisClient)

	// Initialize security components
	waf := security.NewWAF()
	rateLimiter := security.NewRateLimiter(redisClient, 100, time.Minute)

	// Initialize sync manager
	syncManager := sync.NewSyncManager(db, redisClient)

	// Initialize offline store
	offlineStore, err := offline.NewOfflineStore("./offline-data")
	if err != nil {
		log.Fatalf("Failed to initialize offline store: %v", err)
	}
	defer offlineStore.Close()

	// Update router setup to include new middleware
	r := router.Setup(cfg, authService, examService, waf, rateLimiter)

	// Add health check endpoint
	r.Get("/health", func(c *fiber.Ctx) error {
		status := healthChecker.GetStatus()
		if !status.Healthy {
			return c.Status(fiber.StatusServiceUnavailable).JSON(status)
		}
		return c.JSON(status)
	})

	// Setup router
	r := router.Setup(cfg, authService, examService)

	// Start server
	go func() {
		if err := r.Listen(cfg.Server.Address); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}
}
