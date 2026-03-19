package main

import (
	"context"
	"daily-english-words/internal/config"
	"daily-english-words/internal/database"
	"daily-english-words/internal/handler"
	"daily-english-words/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Step 1: Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Step 2: Connect to database
	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	// Step 3: Initialize layers
	wordRepo := repository.NewWordRepository(pool)
	wordHandler := handler.NewWordHandler(wordRepo)

	// Step 4: Register routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/words/today", wordHandler.HandleTodayWords)
	mux.HandleFunc("/api/words/random", wordHandler.HandleRandomWords)

	// Serve frontend static files
	// The web/ directory contains index.html and static assets
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	// Step 5: Start server with graceful shutdown
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	// Listen for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func ()  {
		log.Printf("Server starting on http://localhost:%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Block until we receive a shutdown signal
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}