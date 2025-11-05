# Auth Service

## Quick Start

```bash
# Local Setup
cp .env.example .env
go mod download
go run main.go

# Docker Setup
docker-compose up -d --build auth-service
```

## Project Structure

```
auth-service/
├── main.go                    # Fiber app entry point
├── internal/
│   ├── config.go              # Config & env vars
│   ├── user.go                # User model
│   ├── user_handler.go        # User handlers
│   ├── user_repository.go     # User DB ops
│   ├── session.go             # Session model + device hashing
│   ├── session_repository.go  # Session DB ops
│   └── errors.go              # Error types
├── static/swagger.json        # API documentation
└── Dockerfile
```

