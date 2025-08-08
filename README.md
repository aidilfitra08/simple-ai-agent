# Simple AI Agent

A Go-based REST API that serves as a personal AI assistant using Ollama for natural language processing. The AI agent provides philosophical responses and maintains conversation context.

## Features

- 🤖 **AI-Powered Conversations**: Uses Ollama's `qwen3:0.6b` model for generating responses
- 🎭 **Philosophical Personality**: AI responds with philosophical insights
- 🔧 **Environment-Based Configuration**: Supports development, test, and production modes
- 🗄️ **PostgreSQL Integration**: User data storage with GORM
- 🛡️ **Security**: Configurable trusted proxies based on environment
- 🧠 **Smart Response Processing**: Removes AI thinking blocks from responses
- ⚡ **Gin Framework**: Fast HTTP web framework

## Prerequisites

- Go 1.19 or higher
- PostgreSQL database
- Ollama installed and running locally
- Ollama `qwen3:0.6b` model downloaded

### Install Ollama Model

```bash
ollama pull qwen3:0.6b
```

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
# Database Configuration
DB_HOST=localhost
DB_USER=your-db-password
DB_PASSWORD=your-password
DB_NAME=your-db-name
DB_PORT=5432

# Application Configuration
APP_PORT=8080
APP_ENV=development

# Security Configuration
TRUSTED_PROXIES=127.0.0.1,::1,localhost
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

## How It Works

1. **User sends a message** via POST request to `/chat`
2. **User lookup**: The system fetches user data from PostgreSQL
3. **Prompt construction**: Creates a philosophical system prompt
4. **AI processing**: Sends the prompt to Ollama's `qwen3:0.6b` model
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

1. **Ollama Connection Error**

   - Ensure Ollama is running: `ollama serve`
   - Verify the model is installed: `ollama list`

2. **Database Connection Error**

   - Check PostgreSQL is running
   - Verify database credentials in `.env`
   - Ensure database exists

3. **Port Already in Use**
   - Change `APP_PORT` in `.env` file
   - Kill the process using the port: `netstat -ano | findstr :8080`

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

- [Ollama](https://ollama.ai/) for providing the AI model infrastructure
- [Gin Framework](https://gin-gonic.com/) for the excellent HTTP framework
- [GORM](https://gorm.io/) for the powerful ORM capabilities
