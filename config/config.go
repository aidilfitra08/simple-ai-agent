package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// LLMProvider represents the type of LLM provider
type LLMProvider string

const (
	ProviderLocal  LLMProvider = "local"
	ProviderGemini LLMProvider = "gemini"
)

// DBProvider represents the database provider type
type DBProvider string

const (
	DBPostgres DBProvider = "postgres"
	DBSupabase DBProvider = "supabase"
)

// Config holds all application configuration
type Config struct {
	// Application settings
	AppEnv         string
	AppPort        string
	TrustedProxies []string

	// LLM settings
	LLMProvider LLMProvider

	// Gemini settings
	GeminiAPIKey string
	GeminiModel  string

	// Ollama settings
	OllamaURL   string
	OllamaModel string

	// Database settings
	DBProvider DBProvider
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	// For Supabase or URL-based connections
	SupabaseDBURL string // e.g., postgres://... from Supabase
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	cfg := &Config{
		// Application
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		// LLM Provider
		LLMProvider: LLMProvider(getEnv("LLM_PROVIDER", "local")),

		// Gemini
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		GeminiModel:  getEnv("GEMINI_MODEL", "gemini-1.5-flash-latest"),
		OllamaURL:    getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:  getEnv("OLLAMA_MODEL", "qwen3:0.6b"),

		// Database
		DBProvider:    DBProvider(getEnv("DB_PROVIDER", string(DBPostgres))),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "simple_ai_agent"),
		DBPort:        getEnv("DB_PORT", "5432"),
		SupabaseDBURL: getEnv("SUPABASE_DB_URL", getEnv("DB_URL", "")),
	}

	// Parse trusted proxies
	trustedProxiesEnv := os.Getenv("TRUSTED_PROXIES")
	if trustedProxiesEnv != "" {
		for _, proxy := range strings.Split(trustedProxiesEnv, ",") {
			proxy = strings.TrimSpace(proxy)
			if proxy != "" {
				cfg.TrustedProxies = append(cfg.TrustedProxies, proxy)
			}
		}
	}

	return cfg
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if c.DBProvider == DBSupabase {
		// Prefer full URL if provided (Supabase connection string)
		if c.SupabaseDBURL != "" {
			return c.SupabaseDBURL
		}
		// Fall back to host params with SSL enabled if using Supabase fields only
		return "host=" + c.DBHost +
			" user=" + c.DBUser +
			" password=" + c.DBPassword +
			" dbname=" + c.DBName +
			" port=" + c.DBPort +
			" sslmode=require"
	}
	// Default: local Postgres style DSN
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable"
}
