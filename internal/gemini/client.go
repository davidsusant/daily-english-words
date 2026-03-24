package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// The body we send to the Gemini REST API
type geminiRequest struct {
	Contents 			[]content 			`json:"contents"`
	GenerationConfig 	generationConfig 	`json:"generationConfig"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	Temperature 	float64 `json:"temperature"`
	ResponseMime 	string 	`json:"responseMimeType"`
}

// This is what Gemini sends back
type geminiResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content contentBlock `json:"content"`
}

type contentBlock struct {
	Parts []part `json:"parts"`
}

// Single vocabulary word returned by Gemini
type GeneratedWord struct {
	Word 			string 	`json:"word"`
	PartOfSpeech 	string 	`json:"part_of_speech"`
	Definition 		string 	`json:"definition"`
	Example 		string 	`json:"example"`
}

// Wraps the Gemini API
type Client struct {
	apiKey		string
	httpClient	*http.Client
	model		string
}

// Creates a Gemini client with the given API key
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		model: "gemini-2.5-flash",
	}
}

// Asks Gemini to produce vocabulary words
// Returns structured word data parsed from the JSON response
func (c *Client) GenerateWords(ctx context.Context, count int, difficulty string) ([]GeneratedWord, error) {
	prompt := buildPrompt(count, difficulty)

	// Build the request body
	reqBody := geminiRequest{
		Contents: []content{
			{
				Parts: []part{{Text: prompt}},
			},
		},
		GenerationConfig: generationConfig{
			Temperature: 0.8,
			ResponseMime: "application/json",
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Build the HTTP request
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.model, c.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	// Parse the Gemini response structure
	var geminiResp geminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("unmarshal gemini response: %w", err)
	}

	// Extract the text content from the first candidate
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	jsonText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Parse the JSON array from words
	var words []GeneratedWord
	if err := json.Unmarshal([]byte(jsonText), &words); err != nil {
		return nil, fmt.Errorf("unmarshal words JSON: %w (raw: %s)", err, jsonText)
	}

	return words, nil
}

// Creates the instruction for Gemini
func buildPrompt(count int, difficulty string) string {
	return fmt.Sprintf(`You are an English vocabulary teacher. Generate exactly %d English vocabulary words at the "%s" difficulty level.

	Rules:
	- Choose useful, real-world words that an English learner would benefit from knowing.
	- Each word must be distinct and not a common everyday word (no "hello", "good", "run", etc.).
	- The example sentence should clearly demonstrate the word's meaning in context.
	- Part of speech must be one of: noun, verb, adjective, adverb, preposition.

	Respond with a JSON array only. Each element must have these exact fields:
	- "word": the vocabulary word (lowercase)
	- "part_of_speech": one of noun, verb, adjective, adverb, preposition
	- "definition": a clear, concise definition (one sentence)
	- "example": an example sentence using the word

	Example format:
	[
		{
			"word": "meticulous",
			"part_of_speech": "adjective",
			"definition": "Showing great attention to detail; very careful and precise.",
			"example": "The meticulous engineer reviewed every line of the specification."
		}
	]

	Generate %d words now.`, count, difficulty, count)
}