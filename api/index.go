package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/aidilfitra08/simple-ai-agent/database"
	"github.com/aidilfitra08/simple-ai-agent/handlers"
	"github.com/aidilfitra08/simple-ai-agent/services"
	"github.com/gin-gonic/gin"
)

var (
	router http.Handler
	once   sync.Once
)

func initRouter() {
	cfg := config.Load()
	log.Printf("Using LLM Provider: %s", cfg.LLMProvider)

	// Gin mode
	switch cfg.AppEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("DB connection failed: %v", err)
		router = serviceUnavailable()
		return
	}

	// ❌ Do NOT auto-migrate on Vercel
	// Run migrations separately (CI / manual)

	llmService := services.NewLLMService(cfg)

	chatHandler := handlers.NewChatHandler(db, llmService)
	modelsHandler := handlers.NewModelsHandler(cfg)

	r := gin.New()
	r.Use(gin.Recovery())

	if cfg.AppEnv == "production" {
		_ = r.SetTrustedProxies([]string{"google.com"})
	} else if len(cfg.TrustedProxies) > 0 {
		_ = r.SetTrustedProxies(cfg.TrustedProxies)
	}

	r.POST("/chat", chatHandler.HandleChat)
	r.POST("/chat/stream", chatHandler.HandleChatStream)
	r.GET("/models/gemini", modelsHandler.ListGeminiModels)

	router = r
}

func serviceUnavailable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	})
}

// Handler is the Vercel Function entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initRouter)
	router.ServeHTTP(w, r)
}
