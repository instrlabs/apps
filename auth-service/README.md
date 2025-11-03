# Auth Service

Authentication service for Instrlabs platform with session-based token binding and device validation.

## Features

- PIN-based and Google OAuth authentication
- JWT token generation and validation
- Session management with device binding (IP + User-Agent hash)
- Multiple concurrent sessions with per-device revocation
- Email verification via SMTP

## Quick Start

```bash
# Local Setup
cp .env.example .env
go mod download
go run main.go

# Docker Setup
docker-compose up -d --build auth-service
```

## Configuration

Refer to `.env.example` file for the list of required environment variables and their default values.

## Authentication Flows

### PIN Authentication

```
POST /auth/send-pin
- Generate 6-digit PIN (or 000000 if PIN_ENABLED=false)
- Hash with bcrypt, 10-min expiry
- Send via email

POST /auth/login
- Validate email + PIN
- Create session with device hash (SHA256(IP + User-Agent))
- Return JWT access token + refresh token in response body
```

### Google OAuth

```
GET /auth/google
- Redirect to Google OAuth consent

GET /auth/google/callback
- Exchange code for access token
- Fetch user profile
- Create/update user and session
- Redirect to frontend with tokens as URL parameters
```

### Session Management

**Device Binding**
- Each session tied to device via `SHA256(IP + User-Agent)`
- Device hash validated on token refresh
- Mismatch detection deactivates session immediately

**Token Refresh**
```
POST /auth/refresh
Body: {"refresh_token": "<refresh_token>"}
- Validate refresh token from request body
- Check device hash match
- Issue new access and refresh tokens
- Return tokens in response body
- Update session activity
```

**Session Schema**
```go
type Session struct {
    UserID         string
    SessionID      string    // Unique per device
    DeviceHash     string    // SHA256(IP + User-Agent)
    IPAddress      string
    UserAgent      string
    RefreshToken   string    // Hashed
    IsActive       bool
    LastActivityAt time.Time
    ExpiresAt      time.Time // 7 days default
}
```

**Device Management**
```
GET    /devices              - List all active sessions
POST   /devices/:id/revoke   - Revoke specific session
POST   /devices/revoke-all   - Revoke all sessions
DELETE /auth/logout          - Clear current session
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

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - Database
- [JWT Go](https://github.com/golang-jwt/jwt) - Token handling
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing
