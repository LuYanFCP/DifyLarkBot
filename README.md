# Dify Lark Bot

A production-ready Go application that integrates Dify AI with Lark (Feishu) Bot using event-driven messaging.

## Features

- ✅ Receives messages from Lark Bot when @mentioned
- ✅ Processes messages through Dify AI API
- ✅ Returns AI-generated responses to Lark with @mention
- ✅ **Asynchronous message processing** for better performance
- ✅ **Automatically @mentions the sender** in replies
- ✅ Event-driven architecture with proper error handling
- ✅ WebSocket long connection for real-time event receiving
- ✅ Command-line interface with flags and help
- ✅ Configuration file support (JSON)
- ✅ Graceful shutdown
- ✅ Concurrent message processing capability

## Project Structure

```
.
├── main.go              # Application entry point with enhanced features
├── cmd/                 # Command-line interface
├── config/              # Configuration management with file support
├── dify/               # Dify API client
├── lark/               # Lark Bot service
├── adapter/            # Type adapters
├── Dockerfile          # Container build file
├── .env.example        # Environment variables template
└── config.example.json # Configuration file template
```

## Quick Start

### Using Environment Variables

```bash
# Set required environment variables
export LARK_APP_ID=your_lark_app_id
export LARK_APP_SECRET=your_lark_app_secret
export LARK_VERIFICATION_TOKEN=your_lark_verification_token
export DIFY_API_KEY=your_dify_api_key

# Run the application
./dify_lark_bot
```

### Using Configuration File

```bash
# Create config file from template
cp config.example.json config.json
# Edit config.json with your settings

# Run with config file
./dify_lark_bot --config config.json
```

### Command Line Options

```bash
# Show help
./dify_lark_bot --help

# Show version
./dify_lark_bot --version

# Use config file
./dify_lark_bot --config config.json
./dify_lark_bot -c config.json
```

## Installation

```bash
# Clone repository
git clone <repository-url>
cd DifyLarkBot

# Install dependencies
go mod download

# Build application
go build -o dify_lark_bot .

# Run application
./dify_lark_bot --help
```

## Docker

```bash
# Build image
docker build -t dify-lark-bot .

# Run container
docker run -p 8080:8080 --env-file .env dify-lark-bot

# Or with config file
docker run -p 8080:8080 -v $(pwd)/config.json:/app/config.json dify-lark-bot --config config.json
```

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
LARK_APP_ID=your_lark_app_id
LARK_APP_SECRET=your_lark_app_secret
LARK_VERIFICATION_TOKEN=your_lark_verification_token
LARK_ENCRYPT_KEY=your_lark_encrypt_key
DIFY_API_KEY=your_dify_api_key
DIFY_BASE_URL=https://api.dify.ai
```

### Configuration File (JSON)

Copy `config.example.json` to `config.json`:

```json
{
  "lark": {
    "app_id": "your_lark_app_id",
    "app_secret": "your_lark_app_secret",
    "verification_token": "your_verification_token",
    "encrypt_key": "your_encrypt_key"
  },
  "dify": {
    "api_key": "your_dify_api_key",
    "base_url": "https://api.dify.ai",
    "timeout": 30
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  }
}
```

## Architecture

This application uses WebSocket long connection to receive events from Lark in real-time, eliminating the need for webhook configuration. The application connects directly to Lark's event stream and processes messages as they arrive.

### Key Features

- **Asynchronous Message Processing**: Messages are processed in background goroutines for better performance
- **@Mention Replies**: Bot automatically @mentions the sender when replying to messages
- **Real-time WebSocket Connection**: Direct event streaming without webhook setup
- **Concurrent Processing**: Multiple messages can be processed simultaneously

### Example Responses

```json
// GET /health
{
  "status": "ok",
  "app": "DifyLarkBot",
  "version": "1.0.0",
  "timestamp": 1234567890,
  "uptime": 3600.5
}

// GET /ready
{
  "status": "ready",
  "app": "DifyLarkBot",
  "version": "1.0.0"
}

// GET /info
{
  "app": "DifyLarkBot",
  "version": "1.0.0",
  "go_version": "1.21+",
  "build_time": "unknown",
  "git_commit": "unknown",
  "start_time": "2024-01-01T12:00:00Z",
  "uptime": 3600.5,
  "config": {
    "port": "8080",
    "dify_base_url": "https://api.dify.ai",
    "lark_app_id": "your_app_id"
  }
}
```

## Usage

1. **Configure Lark Bot**: The application uses WebSocket long connection, so no webhook URL configuration is needed
2. **Start the Application**: Run the bot with your configuration
3. **Chat with Bot**: When users @mention your bot in Lark, it will:
   - Receive the message via WebSocket
   - Process it asynchronously through Dify AI
   - Reply with an @mention to the original sender
4. **Monitor Logs**: Check console output for processing status and any errors

## Production Deployment

### Systemd Service

Create `/etc/systemd/system/dify-lark-bot.service`:

```ini
[Unit]
Description=Dify Lark Bot
After=network.target

[Service]
Type=simple
User=dify-lark-bot
WorkingDirectory=/opt/dify-lark-bot
ExecStart=/opt/dify-lark-bot/dify_lark_bot --config /etc/dify-lark-bot/config.json
Restart=always
RestartSec=10
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

### Docker Compose

```yaml
version: '3.8'
services:
  dify-lark-bot:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    env_file:
      - .env
    volumes:
      - ./config.json:/app/config.json
    command: ["./dify_lark_bot", "--config", "/app/config.json"]
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Dependencies

- **Gin Web Framework**: HTTP server and routing
- **Lark Suite OpenAPI SDK**: Lark Bot integration
- **Standard Go Libraries**: Core functionality

## Development

```bash
# Run in development mode
export GIN_MODE=debug
./dify_lark_bot

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
go vet ./...
```

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request