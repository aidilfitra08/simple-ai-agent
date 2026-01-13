package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aidilfitra08/simple-ai-agent/models"
	"github.com/aidilfitra08/simple-ai-agent/services"
	"github.com/aidilfitra08/simple-ai-agent/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ChatHandler handles chat-related requests
type ChatHandler struct {
	db         *gorm.DB
	llmService *services.LLMService
}

// NewChatHandler creates a new chat handler
func NewChatHandler(db *gorm.DB, llmService *services.LLMService) *ChatHandler {
	return &ChatHandler{
		db:         db,
		llmService: llmService,
	}
}

// HandleChat processes chat requests
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user data
	// var user models.User
	// if err := h.db.First(&user, req.UserID).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 	return
	// }

	// Create system prompt
	systemPrompt := fmt.Sprintf(`
You are a personal AI. You will reply to:

%s

Note: If the user say I, me, myself, it means about them, not you. And sometimes you can be a philosopher too, to answer the user.
`, req.Prompt)

	// Generate response using LLM service
	finalResponse, err := h.llmService.GenerateResponse(systemPrompt)
	if err != nil {
		log.Printf("Error generating response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Final response: %s", finalResponse)

	// Clean response and return
	c.JSON(http.StatusOK, gin.H{
		"reply": utils.RemoveThinkBlocks(finalResponse),
	})
}

// HandleChatStream processes chat requests with Server-Sent Events streaming
func (h *ChatHandler) HandleChatStream(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user data
	// var user models.User
	// if err := h.db.First(&user, req.UserID).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 	return
	// }

	// Create system prompt
	systemPrompt := fmt.Sprintf(`
You are a personal AI. You will reply to:

%s

Note: If the user say I, me, myself, it means about them, not you. And sometimes you can be a philosopher too, to answer the user.
`, req.Prompt)

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no") // Disable proxy buffering

	// Channels for streaming
	streamChan := make(chan string)
	errChan := make(chan error)

	// Start streaming in goroutine
	go h.llmService.GenerateStreamResponse(c.Request.Context(), systemPrompt, streamChan, errChan)

	// Stream responses to client
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	fullResponse := ""

	for {
		select {
		case chunk, ok := <-streamChan:
			if !ok {
				// Stream finished successfully
				// Send final cleaned message
				cleanedResponse := utils.RemoveThinkBlocks(fullResponse)
				c.SSEvent("done", gin.H{"reply": cleanedResponse})
				flusher.Flush()
				return
			}

			fullResponse += chunk
			// Send chunk
			c.SSEvent("message", gin.H{"chunk": chunk})
			flusher.Flush()

		case err := <-errChan:
			if err != nil {
				log.Printf("Error streaming response: %v", err)
				c.SSEvent("error", gin.H{"error": err.Error()})
				flusher.Flush()
				return
			}

		case <-c.Request.Context().Done():
			log.Println("Client disconnected")
			return
		}
	}
}
