package repository

import (
	"context"
	"daily-english-words/internal/model"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handles all word-related database operations
type WordRepository struct {
	pool *pgxpool.Pool
}

// Creates a new repository with the given connection pool
func NewWordRepository(pool *pgxpool.Pool) *WordRepository {
	return &WordRepository{pool: pool}
}

// Returns words assigned to today's date
// If no words are assigned yet, it auto-assigns 5 random unassigned words
func (r *WordRepository) GetTodayWords(ctx context.Context) ([]model.Word, error) {
	today := time.Now().Format("2006-01-02")

	// First, check if today already has assigned words
	words, err := r.getWordsByDate(ctx, today)
	if err != nil {
		return nil, fmt.Errorf("get words by date: %w", err)
	}

	// If we already have words for today, return them
	if len(words) > 0 {
		return words, nil
	}

	// Otherwise, assign 5 random unassigned words to today
	if err := r.assignWWordsToDate(ctx, today, 5); err != nil {
		return nil, fmt.Errorf("assign words: %w", err)
	}

	// Fetch the newly assigned words
	return r.getWordsByDate(ctx, today)
}

// Returns n random words regardless of assignment
// Useful for extra practice or quiz mode
func (r *WordRepository) GetRandomWords(ctx context.Context, n int) ([]model.Word, error) {
	query := `
		SELECT id, word, part_of_speech, definition, example, difficulty, assigned_date, created_at
		FROM words
		ORDER BY RANDOM()
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, n)
	if err != nil {
		return nil, fmt.Errorf("query random words: %w", err)
	}

	defer rows.Close()

	return scanWords(rows)
}

// Fetches all words assigned to a specific date
func (r *WordRepository) getWordsByDate(ctx context.Context, date string) ([]model.Word, error) {
	query := `
		SELECT id, word, part_of_speech, definition, example, difficulty, assigned_date, created_at
		FROM words
		WHERE assigned_date = $1
		ORDER BY difficulty ASC
	`

	rows, err := r.pool.Query(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("query words by date: %w", err)
	}

	defer rows.Close()

	return scanWords(rows)
}

// Helper that reads rows into a slice of Word structs
func scanWords(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]model.Word, error) {
	var words []model.Word

	for rows.Next() {
		var w model.Word
		err := rows.Scan(
			&w.ID,
			&w.Word,
			&w.PartOfSpeech,
			&w.Definition,
			&w.Example,
			&w.Difficulty,
			&w.AssignedDate,
			&w.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan word row: %w", err)
		}

		words = append(words, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return words, nil
}

// Picks n unassigned words and sets their assigned date
func (r *WordRepository) assignWWordsToDate(ctx context.Context, date string, n int) error {
	query := `
		UPDATE words
		SET assigned_date = $1
		WHERE id IN (
			SELECT id FROM words
			WHERE assigned_date IS NULL
			ORDER BY RANDOM()
			LIMIT $2
		)
	`

	_, err := r.pool.Exec(ctx, query, date, n)
	if err != nil {
		return fmt.Errorf("assign words to date: %w", err)
	}

	return  nil
}