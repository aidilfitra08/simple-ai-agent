package main

import (
	"log"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/aidilfitra08/simple-ai-agent/database"
	"github.com/aidilfitra08/simple-ai-agent/handlers"
	"github.com/aidilfitra08/simple-ai-agent/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Using LLM Provider: %s", cfg.LLMProvider)

	// Configure Gin mode based on environment
	switch cfg.AppEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize services
	llmService := services.NewLLMService(cfg)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(db, llmService)
	modelsHandler := handlers.NewModelsHandler(cfg)

	// Setup router
	r := gin.Default()

	// Configure trusted proxies
	if cfg.AppEnv == "production" {
		r.SetTrustedProxies([]string{"google.com"})
	} else if len(cfg.TrustedProxies) > 0 {
		r.SetTrustedProxies(cfg.TrustedProxies)
	}

	// Register routes
	r.POST("/chat", chatHandler.HandleChat)
	r.POST("/chat/stream", chatHandler.HandleChatStream)
	r.GET("/models/gemini", modelsHandler.ListGeminiModels)

	// Start server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
