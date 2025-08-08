package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model (adjust fields as needed)
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string
    Interests string
}

// ChatRequest from frontend
type ChatRequest struct {
    UserID uint   `json:"user_id"`
    Prompt string `json:"prompt"`
}

// Ollama request format
type OllamaRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Temperature   float64 `json:"temperature"`
    TopP         float64 `json:"top_p"`
    RepeatPenalty float64 `json:"repeat_penalty"`
}

// Ollama streaming response chunk
type OllamaResponse struct {
    Model    string `json:"model"`
    CreatedAt string `json:"created_at"`
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

func main() {
		err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using system environment variables")
    }

    // Get app environment and configure Gin mode
    appEnv := os.Getenv("APP_ENV")
    if appEnv == "" {
        appEnv = "development"
    }
    
    switch appEnv {
    case "production":
        gin.SetMode(gin.ReleaseMode)
    case "test":
        gin.SetMode(gin.TestMode)
    default:
        gin.SetMode(gin.DebugMode)
    }
		
    // DB connection
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // Auto-migrate
    db.AutoMigrate(&User{})

    r := gin.Default()

		// Configure trusted proxies based on environment
		// Configure trusted proxies from environment
    trustedProxiesEnv := os.Getenv("TRUSTED_PROXIES")
    var trustedProxies []string

    if trustedProxiesEnv != "" {
        // Split by comma and trim spaces
        for _, proxy := range strings.Split(trustedProxiesEnv, ",") {
            proxy = strings.TrimSpace(proxy)
            if proxy != "" {
                trustedProxies = append(trustedProxies, proxy)
            }
        }
    }

    if appEnv == "production" {
        r.SetTrustedProxies([]string{"google.com"}) // Only google.com in production
    } else {
        r.SetTrustedProxies(trustedProxies) // More permissive for dev
    }

    // r.SetTrustedProxies([]string{"localhost"}) // Trust only localhost for development

    r.POST("/chat", func(c *gin.Context) {
        var req ChatRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Fetch user data
        var user User
        if err := db.First(&user, req.UserID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        // Combine user data into system prompt
        // systemPrompt := fmt.Sprintf(
        //     "You are a personal AI assistant for %s. User's interests: %s.\nQuestion: %s. If the user use porn word, reply with 'I don't know what you are talking about'.",
        //     user.Name, user.Interests, req.Prompt,
        // )
        systemPrompt := fmt.Sprintf(`
            You are a personal AI. You will reply to:

            %s

            Note: If the user say I, me, myself, it means about them, not you. And sometimes you can be a philosopher too, to answer the user.
            `, req.Prompt)

        // Call Ollama
        ollamaReq := OllamaRequest{
            Model:  "qwen3:0.6b",
            Prompt: systemPrompt,
			Temperature:  0,
        }

        jsonData, _ := json.Marshal(ollamaReq)
        resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer resp.Body.Close()

        // Ollama streams JSON lines, so we concatenate
        var finalResponse string
        decoder := json.NewDecoder(resp.Body)
        for {
            var chunk OllamaResponse
            if err := decoder.Decode(&chunk); err == io.EOF {
                break
            } else if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            finalResponse += chunk.Response
        }
        log.Printf("Final response: %s", finalResponse)
        c.JSON(http.StatusOK, gin.H{
            "reply": removeThinkBlocks(finalResponse),
        })
    })

		// Get port from environment or default to 8080
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    r.Run(":" + port)
}

func removeThinkBlocks(input string) string {
	// Regex to match <think>...</think> including newlines
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	// Replace with empty string

    // Replace with empty string
    result := re.ReplaceAllString(input, "")
    
    // Remove only leading newlines (at the beginning of the string)
    result = strings.TrimLeft(result, "\n")
	return result
}