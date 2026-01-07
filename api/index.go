package handler

import (
	"log"
	"net/http"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/aidilfitra08/simple-ai-agent/database"
	"github.com/aidilfitra08/simple-ai-agent/handlers"
	"github.com/aidilfitra08/simple-ai-agent/services"
	"github.com/gin-gonic/gin"
)

var router http.Handler

func init() {
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
		log.Printf("Failed to connect to database: %v", err)
		router = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		})
		return
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		router = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		})
		return
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

	router = r
}

// Handler is the Vercel serverless function entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
