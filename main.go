package main

import (
	"log"
	"net/http"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/aidilfitra08/simple-ai-agent/database"
	"github.com/aidilfitra08/simple-ai-agent/handlers"
	"github.com/aidilfitra08/simple-ai-agent/services"
	"github.com/gin-gonic/gin"
)

// Handler is exported for Vercel Serverless Functions runtime.
var Handler http.Handler

// buildHandler constructs the HTTP handler (Gin engine) with all dependencies.
func buildHandler() (http.Handler, error) {
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
		return nil, err
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		return nil, err
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

	return r, nil
}

func init() {
	h, err := buildHandler()
	if err != nil {
		log.Printf("Failed to initialize handler: %v", err)
		Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		})
		return
	}
	Handler = h
}

func main() {
	// Local/dev server run path
	cfg := config.Load()
	h, err := buildHandler()
	if err != nil {
		log.Fatalf("Startup failure: %v", err)
	}
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, h); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
