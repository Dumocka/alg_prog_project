package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"testing-platform/core/internal/api"
	"testing-platform/core/internal/config"
	"testing-platform/core/internal/repository"
	"testing-platform/core/internal/service"
)

func main() {
	cfg := config.Load()

	// Initialize repositories
	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize services
	svc := service.NewService(repo)

	// Initialize HTTP server
	server := api.NewServer(cfg, svc)

	// Start server
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
