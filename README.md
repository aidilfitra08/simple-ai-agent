# Simple AI Agent

A Go-based REST API that serves as a personal AI assistant with support for both Google Gemini and local Ollama LLMs. The AI agent provides philosophical responses and maintains conversation context.

## Features

- 🤖 **Dual LLM Support**: Switch between Google Gemini and local Ollama models via environment variable
- 🎭 **Philosophical Personality**: AI responds with philosophical insights
- 🔧 **Environment-Based Configuration**: Supports development, test, and production modes
- 🗄️ **PostgreSQL Integration**: User data storage with GORM
- 🛡️ **Security**: Configurable trusted proxies based on environment
- 🧠 **Smart Response Processing**: Removes AI thinking blocks from responses
- ⚡ **Gin Framework**: Fast HTTP web framework
- ☁️ **Vercel Compatible**: Ready for serverless deployment on Vercel

## Prerequisites

- Go 1.23 or higher
- PostgreSQL database
- **For Local LLM**: Ollama installed and running locally with model downloaded
- **For Gemini**: Google Gemini API key

### Option 1: Local LLM Setup (Ollama)

```bash
# Install Ollama from https://ollama.ai/
# Pull the model
ollama pull qwen3:0.6b
```

### Option 2: Google Gemini Setup

1. Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Set `LLM_PROVIDER=gemini` in your `.env` file
3. Add your API key to `GEMINI_API_KEY`

## Installation

1. **Clone the repository**

   ```bash
   git clone <your-repo-url>
   cd simple-ai-agent
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Create environment file**

   ```bash
   cp .env.example .env
   ```

4. **Configure environment variables** (see Configuration section)

5. **Run the application**
   ```bash
   go run main.go
   ```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Application Environment
APP_ENV=development
APP_PORT=8080

# LLM Provider Configuration
# Options: "local" or "gemini"
LLM_PROVIDER=local

# Gemini API Configuration (required if LLM_PROVIDER=gemini)
GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-1.5-flash-latest

# Ollama Configuration (required if LLM_PROVIDER=local)
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=qwen3:0.6b

# Database Configuration
DB_PROVIDER=postgres # options: postgres | supabase
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=simple_ai_agent
DB_PORT=5432

# Supabase (if using DB_PROVIDER=supabase)
# Prefer a full connection string; example:
# SUPABASE_DB_URL=postgres://USER:PASSWORD@HOST:PORT/DB?sslmode=require
# Alternatively set the DB_* fields above; sslmode will default to require
SUPABASE_DB_URL=

# Security Configuration - Trusted Proxies (comma-separated)
TRUSTED_PROXIES=localhost,127.0.0.1
```

### LLM Provider Configuration

#### Using Local Ollama (LLM_PROVIDER=local)

- Set `LLM_PROVIDER=local`
- Configure `OLLAMA_URL` (default: `http://localhost:11434`)
- Configure `OLLAMA_MODEL` (default: `qwen3:0.6b`)
- Ensure Ollama is running: `ollama serve`

#### Using Google Gemini (LLM_PROVIDER=gemini)

- Set `LLM_PROVIDER=gemini`
- Set `GEMINI_API_KEY` to your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
- Configure `GEMINI_MODEL` (default: `gemini-1.5-flash-latest`)
  - Available models: `gemini-1.5-flash-latest`, `gemini-1.5-pro-latest`, `gemini-pro`, `gemini-1.0-pro`
  - Note: Use the `-latest` suffix for the most up-to-date model versions

### Database Provider (Postgres vs Supabase)

You can switch database backends via `DB_PROVIDER`.

#### Using local/managed Postgres (`DB_PROVIDER=postgres`)

- Set `DB_PROVIDER=postgres` (default)
- Configure `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- Local development uses `sslmode=disable` by default

#### Using Supabase (`DB_PROVIDER=supabase`)

- Set `DB_PROVIDER=supabase`
- Preferred: set `SUPABASE_DB_URL` to the connection string from Supabase (includes sslmode)
- Alternative: set `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`; `sslmode` will default to `require`

Examples:

```env
# Local Postgres
DB_PROVIDER=postgres
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=simple_ai_agent
DB_PORT=5432

# Supabase via URL
DB_PROVIDER=supabase
SUPABASE_DB_URL=postgres://USER:PASSWORD@HOST:6543/postgres?sslmode=require
```

### Environment Modes

- **`development`**: Debug mode with verbose logging
- **`production`**: Release mode with minimal logging and strict security
- **`test`**: Test mode for running tests

## API Endpoints

### Chat with AI

**POST** `/chat`

Send a message to the AI assistant.

**Request Body:**

```json
{
  "user_id": 1,
  "prompt": "What is the meaning of life?"
}
```

**Response:**

```json
{
  "reply": "The meaning of life, from a philosophical perspective, is a question that has puzzled humanity for centuries..."
}
```

## Database Schema

### Users Table

| Column    | Type   | Description      |
| --------- | ------ | ---------------- |
| id        | uint   | Primary key      |
| name      | string | User's name      |
| interests | string | User's interests |

## Project Structure

```
simple-ai-agent/
├── main.go              # Main application file
├── .env                 # Environment variables
├── .env.example         # Example environment file
├── go.mod              # Go module file
├── go.sum              # Go dependencies
└── README.md           # This file
```

## Dependencies

- **[Gin](https://github.com/gin-gonic/gin)**: HTTP web framework
- **[GORM](https://gorm.io/)**: ORM library for Go
- **[PostgreSQL Driver](https://github.com/go-gorm/postgres)**: PostgreSQL driver for GORM
- **[godotenv](https://github.com/joho/godotenv)**: Load environment variables from .env file
- **[Google Generative AI Go SDK](https://github.com/google/generative-ai-go)**: Google Gemini API client

## Deployment

### Deploy to Vercel

1. **Install Vercel CLI**

   ```bash
   npm i -g vercel
   ```

2. **Set up environment variables in Vercel**

   ```bash
   vercel env add LLM_PROVIDER
   vercel env add GEMINI_API_KEY
   vercel env add GEMINI_MODEL
   vercel env add DB_PROVIDER
   vercel env add SUPABASE_DB_URL
   vercel env add DB_HOST
   vercel env add DB_USER
   vercel env add DB_PASSWORD
   vercel env add DB_NAME
   vercel env add DB_PORT
   ```

3. **Deploy**

   ```bash
   vercel
   ```

4. **Note for Vercel**:
   - Use `LLM_PROVIDER=gemini` for Vercel deployments (Ollama local server won't be accessible)
   - Make sure your PostgreSQL database is accessible from Vercel (consider using services like Neon, Supabase, or Railway)
   - When using Supabase, set `DB_PROVIDER=supabase` and `SUPABASE_DB_URL` in Vercel env
   - The `vercel.json` file is already configured for Go deployment

## How It Works

1. **User sends a message** via POST request to `/chat`
2. **User lookup**: The system fetches user data from PostgreSQL
3. **Prompt construction**: Creates a philosophical system prompt
4. **LLM Routing**: Based on `LLM_PROVIDER` environment variable:
   - **Gemini**: Sends to Google Gemini API using the configured model
   - **Local**: Sends to local Ollama instance
5. **Response processing**: Removes AI thinking blocks and formats the response
6. **Return response**: Sends the philosophical answer back to the user

## AI Behavior

The AI assistant is configured to:

- Respond with philosophical insights
- Understand when users refer to themselves ("I", "me", "myself")
- Remove internal thinking processes from responses
- Maintain a personal, thoughtful conversation style

## Security Features

- **Trusted Proxies**: Configurable based on environment
- **Environment-based Security**: Different security levels for dev/prod
- **Input Validation**: Request validation using Gin's binding

## Development

### Running in Development Mode

```bash
# Set environment to development
echo "APP_ENV=development" >> .env

# Run the application
go run main.go
```

### Running in Production Mode

```bash
# Set environment to production
echo "APP_ENV=production" >> .env

# Run the application
go run main.go
```

## Troubleshooting

### Common Issues

1. **Ollama Connection Error** (when using LLM_PROVIDER=local)

   - Ensure Ollama is running: `ollama serve`
   - Verify the model is installed: `ollama list`
   - Check `OLLAMA_URL` is correct in `.env`

2. **Gemini API Error** (when using LLM_PROVIDER=gemini)

   - Verify your API key is valid
   - Check you have API quota available
   - Ensure the model name is correct (`gemini-1.5-flash`, `gemini-1.5-pro`, etc.)

3. **Database Connection Error**

   - Check PostgreSQL is running
   - Verify database credentials in `.env`
   - Ensure database exists

4. **Port Already in Use**

   - Change `APP_PORT` in `.env` file
   - Kill the process using the port: `netstat -ano | findstr :8080`

5. **Vercel Deployment Issues**
   - Ensure you're using `LLM_PROVIDER=gemini` (local Ollama won't work on Vercel)
   - Verify all environment variables are set in Vercel dashboard
   - Check that your database is accessible from Vercel's network

### Logs

The application provides detailed logging including:

- Environment configuration
- Database connections
- AI responses (in development mode)
- Trusted proxy configuration

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Google Gemini](https://ai.google.dev/) for providing the Gemini API
- [Ollama](https://ollama.ai/) for providing the local AI model infrastructure
- [Gin Framework](https://gin-gonic.com/) for the excellent HTTP framework
- [GORM](https://gorm.io/) for the powerful ORM capabilities
