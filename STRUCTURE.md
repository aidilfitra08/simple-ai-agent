# Project Structure

This document describes the organization of the Simple AI Agent codebase.

## Directory Structure

```
simple-ai-agent/
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   └── config.go          # Environment variables and app settings
├── models/                 # Data models
│   ├── user.go            # User database model
│   └── request.go         # API request/response models
├── database/               # Database operations
│   └── database.go        # Connection and migration logic
├── services/               # Business logic layer
│   └── llm.go             # LLM service (Gemini/Ollama integration)
├── handlers/               # HTTP request handlers
│   └── chat.go            # Chat endpoint handler
├── utils/                  # Utility functions
│   └── text.go            # Text processing utilities
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── vercel.json            # Vercel deployment configuration
├── .env.example           # Environment variables template
└── README.md              # Project documentation
```

## Package Descriptions

### `config/`

Handles all application configuration including:

- Environment variable loading
- LLM provider configuration (Gemini/Ollama)
- Database connection settings
- Application settings (port, environment, trusted proxies)

**Key types:**

- `Config`: Main configuration struct
- `LLMProvider`: Enum for LLM provider types

### `models/`

Defines data structures used throughout the application:

- `User`: Database user model
- `ChatRequest`: API request structure
- `OllamaRequest/Response`: Ollama API structures

### `database/`

Manages database connections and migrations:

- `Connect()`: Establishes PostgreSQL connection
- `Migrate()`: Runs GORM auto-migrations

### `services/`

Contains business logic and external service integrations:

- `LLMService`: Handles all LLM operations
  - `GenerateResponse()`: Routes requests to appropriate provider
  - `callGemini()`: Google Gemini API integration
  - `callOllama()`: Local Ollama integration

### `handlers/`

HTTP request handlers for API endpoints:

- `ChatHandler`: Processes chat requests
  - Validates input
  - Fetches user data
  - Generates LLM responses
  - Returns formatted results

### `utils/`

Shared utility functions:

- `RemoveThinkBlocks()`: Cleans AI thinking blocks from responses

## Benefits of This Structure

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Maintainability**: Easy to locate and modify specific functionality
3. **Testability**: Packages can be tested independently
4. **Scalability**: Easy to add new features or extend existing ones
5. **Clean Dependencies**: Clear dependency flow from main → handlers → services → models

## Dependency Flow

```
main.go
  ├─> config (loads settings)
  ├─> database (establishes connection)
  ├─> services (business logic)
  └─> handlers (HTTP layer)
        ├─> models (data structures)
        ├─> services (LLM operations)
        └─> utils (helper functions)
```

## Adding New Features

### Adding a New API Endpoint

1. Define request/response models in `models/`
2. Create handler function in `handlers/`
3. Register route in `main.go`

### Adding a New LLM Provider

1. Add provider constant in `config/config.go`
2. Add configuration fields in `Config` struct
3. Implement provider method in `services/llm.go`
4. Add case in `GenerateResponse()` switch

### Adding a New Database Model

1. Define struct in `models/`
2. Add to `Migrate()` in `database/database.go`

## Best Practices

- Keep handlers thin - move business logic to services
- Use dependency injection (pass dependencies to constructors)
- Handle errors gracefully with descriptive messages
- Use proper logging for debugging
- Keep models simple and focused on data structure
