package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// ModelsHandler handles model-related requests
type ModelsHandler struct {
	config *config.Config
}

// NewModelsHandler creates a new models handler
func NewModelsHandler(cfg *config.Config) *ModelsHandler {
	return &ModelsHandler{
		config: cfg,
	}
}

// ListGeminiModels lists all available Gemini models
func (h *ModelsHandler) ListGeminiModels(c *gin.Context) {
	ctx := context.Background()

	if h.config.GeminiAPIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GEMINI_API_KEY not set"})
		return
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(h.config.GeminiAPIKey))
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Gemini client"})
		return
	}
	defer client.Close()

	// List all available models
	iter := client.ListModels(ctx)

	type ModelInfo struct {
		Name             string   `json:"name"`
		DisplayName      string   `json:"display_name"`
		Description      string   `json:"description"`
		SupportedMethods []string `json:"supported_methods"`
		InputTokenLimit  int32    `json:"input_token_limit"`
		OutputTokenLimit int32    `json:"output_token_limit"`
	}

	var models []ModelInfo

	for {
		model, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error listing models: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list models"})
			return
		}

		// Check if model supports generateContent
		supportsGenerate := false
		for _, method := range model.SupportedGenerationMethods {
			if method == "generateContent" {
				supportsGenerate = true
				break
			}
		}

		// Only include models that support generateContent
		if supportsGenerate {
			models = append(models, ModelInfo{
				Name:             model.Name,
				DisplayName:      model.DisplayName,
				Description:      model.Description,
				SupportedMethods: model.SupportedGenerationMethods,
				InputTokenLimit:  model.InputTokenLimit,
				OutputTokenLimit: model.OutputTokenLimit,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
		"count":  len(models),
		"note":   "These models support generateContent method",
	})
}
