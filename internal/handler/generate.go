package handler

import (
	"daily-english-words/internal/gemini"
	"daily-english-words/internal/model"
	"daily-english-words/internal/repository"
	"encoding/json"
	"log"
	"net/http"
)

// Handles word generation via Gemini
type GenerateHandler struct {
	repo 	*repository.WordRepository
	gemini 	*gemini.Client
}

// Expected JSON body from the frontend
type generateRequest struct {
	Count 		int 	`json:"count"`
	Difficulty 	string 	`json:"difficulty"`
}

// Creates a handler with repo and Gemini client
func NewGenerateHandler(repo *repository.WordRepository, gemini *gemini.Client) *GenerateHandler {
	return &GenerateHandler{
		repo: repo,
		gemini: gemini,
	}
}

func (h *GenerateHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse request body
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Validate inputs
	if req.Count < 1 || req.Count > 20 {
		writeError(w, http.StatusBadRequest, "count must be between 1 and 20")
		return
	}

	validDifficulties := map[string]bool{"easy": true, "medium": true, "hard": true}
	if !validDifficulties[req.Difficulty] {
		writeError(w, http.StatusBadRequest, "difficulty must be easy, medium, or hard")
		return
	}

	// Call Gemini API to generate words
	log.Printf("Generating %d %s words via Gemini...", req.Count, req.Difficulty)
	generated, err := h.gemini.GenerateWords(r.Context(), req.Count, req.Difficulty)
	if err != nil {
		log.Printf("ERROR Gemini generate: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to generate words from Gemini")
		return
	}

	// Convert Gemini output to model.Word slice
	var words []model.Word
	for _, g := range generated {
		words = append(words, model.Word{
			Word: 			g.Word,
			PartOfSpeech: 	g.PartOfSpeech,
			Definition: 	g.Definition,
			Example: 		g.Example,
			Difficulty: 	req.Difficulty,
		})
	}

	// Insert into database (duplicates are silently skipped)
	inserted, err := h.repo.InsertWords(r.Context(), words)
	if err != nil {
		log.Printf("ERROR inserting words: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to save generated words")
		return
	}

	log.Printf("Generated %d words, inserted %d new (skipped %d duplicate)",
		len(generated), len(inserted), len(generated) - len(inserted))

	// Respond with the newly inserted words
	writeJSON(w, http.StatusCreated, map[string]any{
		"message": 		formatMessage(len(generated), len(inserted)),
		"generated": 	len(generated),
		"inserted": 	len(inserted),
		"duplicates": 	len(generated) - len(inserted),
		"words": 		inserted,
	})
	
}

func formatMessage(generated, inserted int) string {
	if inserted == generated {
		return "All words generated and saved succefully"
	}
	if inserted == 0 {
		return "All generated words already exist in the database"
	}
	return "Words generated - some duplucates were skipped"
}