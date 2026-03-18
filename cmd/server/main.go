package main

import (
	"context"
	"daily-english-words/internal/config"
	"daily-english-words/internal/database"
	"daily-english-words/internal/handler"
	"daily-english-words/internal/repository"
	"log"
	"net/http"
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
}