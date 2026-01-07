package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/aidilfitra08/simple-ai-agent/models"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// LLMService handles LLM interactions
type LLMService struct {
	config *config.Config
}

// NewLLMService creates a new LLM service
func NewLLMService(cfg *config.Config) *LLMService {
	return &LLMService{
		config: cfg,
	}
}

// GenerateResponse generates a response using the configured LLM provider
func (s *LLMService) GenerateResponse(prompt string) (string, error) {
	switch s.config.LLMProvider {
	case config.ProviderGemini:
		return s.callGemini(prompt)
	case config.ProviderLocal:
		return s.callOllama(prompt)
	default:
		return "", fmt.Errorf("invalid LLM provider: %s", s.config.LLMProvider)
	}
}

// callGemini sends prompt to Google Gemini API
func (s *LLMService) callGemini(prompt string) (string, error) {
	ctx := context.Background()

	if s.config.GeminiAPIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(s.config.GeminiAPIKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}
	defer client.Close()

	// Use the model name directly - the SDK handles the format
	model := client.GenerativeModel(s.config.GeminiModel)

	// Configure generation settings
	model.SetTemperature(0.7)
	model.SetTopP(0.95)
	model.SetTopK(40)
	model.SetMaxOutputTokens(2048)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	// Extract text from response
	var finalResponse string
	for _, part := range resp.Candidates[0].Content.Parts {
		finalResponse += fmt.Sprintf("%v", part)
	}

	return finalResponse, nil
}

// callOllama sends prompt to local Ollama instance
func (s *LLMService) callOllama(prompt string) (string, error) {
	ollamaReq := models.OllamaRequest{
		Model:       s.config.OllamaModel,
		Prompt:      prompt,
		Temperature: 0,
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(s.config.OllamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	// Ollama streams JSON lines, so we concatenate
	var finalResponse string
	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk models.OllamaResponse
		if err := decoder.Decode(&chunk); err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("failed to decode Ollama response: %w", err)
		}
		finalResponse += chunk.Response
	}

	return finalResponse, nil
}

// GenerateStreamResponse provides a simple streaming interface by chunking
// the full response from GenerateResponse and sending it over stream.
// It closes the stream channel when finished. Errors are sent on errChan.
func (s *LLMService) GenerateStreamResponse(ctx context.Context, prompt string, stream chan<- string, errChan chan<- error) {
	defer close(stream)

	resp, err := s.GenerateResponse(prompt)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		case errChan <- err:
		default:
		}
		return
	}

	for _, chunk := range chunkString(resp, 200) {
		select {
		case <-ctx.Done():
			return
		case stream <- chunk:
		}
		time.Sleep(15 * time.Millisecond)
	}
}

// chunkString splits a string into approximately n-rune chunks (UTF-8 safe).
func chunkString(s string, n int) []string {
	if n <= 0 {
		return []string{s}
	}
	rs := []rune(s)
	var out []string
	for i := 0; i < len(rs); i += n {
		end := i + n
		if end > len(rs) {
			end = len(rs)
		}
		out = append(out, string(rs[i:end]))
	}
	return out
}
