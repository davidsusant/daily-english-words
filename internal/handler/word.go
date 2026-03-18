package handler

import (
	"daily-english-words/internal/repository"
	"encoding/json"
	"log"
	"net/http"
)

// Holds the dependency on the word repository
type WordHandler struct {
	repo *repository.WordRepository
}

// Creates a handler with the given repository
func NewWordHandler(repo *repository.WordRepository) *WordHandler {
	return &WordHandler{repo: repo}
}

// Responds with today's assigned vocabulary words
// GET /api/words/today
func (h *WordHandler) HandleTodayWords(w http.ResponseWriter, r *http.Request) {
	// Only allow GET
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	words, err := h.repo.GetTodayWords(r.Context())
	if err != nil {
		log.Printf("ERROR GetTodayWords: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to get today's words")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"date": words[0].AssignedDate,
		"count": len(words),
		"words": words,
	})
}

// Responds with 5 random words for extra practice
// GET /api/words/random
func (h *WordHandler) HandleRandomWords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	words, err := h.repo.GetRandomWords(r.Context(), 5)
	if err != nil {
		log.Printf("ERROR GetRandomWords: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to get random words")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"count": len(words), 
		"words": words,
	})
}

// Sends a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

// Encodes data as JSON and writes it to the response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ERROR encoding JSON: %v", err)
	}
}