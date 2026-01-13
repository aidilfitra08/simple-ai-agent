package models

// ChatRequest represents a chat request from the frontend
type ChatRequest struct {
	// UserID uint   `json:"user_id"`
	Prompt string `json:"prompt"`
}

// OllamaRequest represents the request format for Ollama API
type OllamaRequest struct {
	Model         string  `json:"model"`
	Prompt        string  `json:"prompt"`
	Temperature   float64 `json:"temperature"`
	TopP          float64 `json:"top_p"`
	RepeatPenalty float64 `json:"repeat_penalty"`
}

// OllamaResponse represents a streaming response chunk from Ollama
type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}
