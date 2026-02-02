# Fly-Go RESTful API

A Golang RESTful API project with MongoDB integration and dual logging system.

## Project Structure

```
fly-go/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   └── mongodb.go       # MongoDB connection
│   ├── handlers/
│   │   └── health.go        # API handlers
│   ├── middleware/
│   │   └── middleware.go    # Middleware (CORS, Logger, Recovery)
│   ├── models/
│   │   └── models.go        # Data models
│   ├── routes/
│   │   └── routes.go        # Route definitions
│   └── utils/
│       └── response.go      # Response utilities
├── pkg/
│   └── logger/
│       └── logger.go        # Logger configuration
├── logs/                    # Log files directory
├── config.yaml              # Configuration file
├── .env.example             # Environment variables example
└── go.mod                   # Go module file
```

## Features

- RESTful API with Gin framework
- MongoDB database integration
- Dual logging system:
  - Console: INFO level
  - File: ERROR level (stored in logs/error.log)
- CORS middleware
- Request logging middleware
- Panic recovery middleware
- Graceful shutdown
- Configuration management with Viper

## Configuration

Edit `config.yaml` to configure your application:

```yaml
server:
  port: "8080"
  mode: "debug"  # debug, release, test

database:
  uri: "mongodb://localhost:27017"
  database: "flygo"
```

Or use environment variables (see `.env.example`):

```bash
SERVER_PORT=8080
SERVER_MODE=debug
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=flygo
```

## Installation

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run cmd/api/main.go
```

3. Build the application:
```bash
go build -o bin/api cmd/api/main.go
```

## API Endpoints

### Health Check
```
GET /api/v1/health
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "message": "Service is running"
  }
}
```

## Logging

- Console logs: INFO level and above
- File logs: ERROR level only (logs/error.log)
- Log rotation: 100MB max size, 3 backups, 30 days retention

## Dependencies

- github.com/gin-gonic/gin - Web framework
- go.mongodb.org/mongo-driver - MongoDB driver
- go.uber.org/zap - Logging
- github.com/spf13/viper - Configuration management
- gopkg.in/natefinch/lumberjack.v2 - Log rotation
