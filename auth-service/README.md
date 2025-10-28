# Auth Service

Authentication and authorization service for Instrlabs platform.

## Features

- User authentication (PIN-based, OAuth)
- Google OAuth integration
- JWT token generation and validation
- Refresh token management
- Email verification
- User session management

## Prerequisites

- Go 1.23+
- MongoDB
- SMTP server (for email verification)

## Configuration

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

Required environment variables:
- `MONGO_URI` - MongoDB connection string
- `MONGO_DB` - Database name
- `JWT_SECRET` - Secret for signing JWT tokens
- `SMTP_*` - Email server configuration
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` - Google OAuth credentials
- `GATEWAY_URL` - Gateway service URL
- `WEB_URL` - Frontend application URL

## Running

```bash
# Install dependencies
go mod download

# Run the service
go run main.go
```

The service will start on the port specified in `PORT` environment variable (default: `:3000`).

## Feature Workflows

### 1. PIN-Based Authentication

**Send PIN**
1. User submits email address
2. System checks if user exists, creates new user if needed
3. Generates 6-digit PIN (or fixed `000000` if `PIN_ENABLED=false`)
4. Hashes PIN with bcrypt and stores with 10-minute expiry
5. Sends PIN via email (if `PIN_ENABLED=true`)

**Login with PIN**
1. User submits email and PIN
2. System validates credentials against stored hash
3. Clears PIN after successful validation
4. Generates JWT access token and refresh token
5. Sets secure HTTP-only cookies (`access_token`, `refresh_token`)
6. Sets `RegisteredAt` timestamp on first successful login

### 2. Google OAuth Flow

**Initiate OAuth**
1. User clicks "Login with Google"
2. System generates OAuth state parameter
3. Redirects to Google OAuth consent screen

**OAuth Callback**
1. Google redirects back with authorization code
2. System exchanges code for access token
3. Fetches user profile from Google API
4. Finds or creates user by Google ID or email
5. Generates JWT access token and refresh token
6. Sets secure HTTP-only cookies
7. Redirects to web application

### 3. Token Refresh

1. Client sends request with `x-user-refresh` header (refresh token)
2. System validates refresh token against database
3. Generates new access token and refresh token
4. Updates refresh token in database
5. Sets new cookies with updated tokens

### 4. User Profile

1. Client sends authenticated request (access token in cookie)
2. Gateway validates JWT and sets `userId` in request context
3. System retrieves user profile from database
4. Returns user data

### 5. Logout

1. Client sends authenticated logout request
2. System clears refresh token from database
3. Clears `access_token` and `refresh_token` cookies (sets to expired)

## API Endpoints

All endpoints are served through the Gateway service. Refer to the API documentation for available routes.

## Project Structure

```
auth-service/
├── main.go                      # Entry point
├── internal/
│   ├── config.go                # Environment configuration
│   ├── user.go                  # User domain model
│   ├── user_handler.go          # HTTP handlers
│   ├── user_repository.go       # Database access
│   └── errors.go                # Error definitions
└── Dockerfile                   # Container definition
```
