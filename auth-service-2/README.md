# Auth Service 2

A modern authentication and authorization service with support for OAuth (Google) and PIN-based authentication.

## Features

### Authentication Methods

#### 1. **OAuth Flow (Google)**
High-level flow for Google OAuth authentication:

```
1. User clicks "Login with Google"
   ↓
2. GET /oauth/google
   - Generates OAuth state for CSRF protection
   - Stores state in cookie
   - Redirects to Google OAuth consent screen
   ↓
3. User authorizes on Google
   ↓
4. GET /oauth/google/callback?code=xxx&state=xxx
   - Validates state parameter
   - Exchanges code for Google access token
   - Fetches user info from Google
   - Creates/updates user in database
   - Generates JWT access & refresh tokens
   - Sets auth cookies
   - Redirects to web app
```

#### 2. **PIN Authentication Flow**
High-level flow for PIN-based authentication:

```
1. User enters email
   ↓
2. POST /auth/pin/request
   - Generates random PIN (default: 6 digits)
   - Hashes PIN using bcrypt
   - Creates/updates user with PIN hash
   - Sends PIN via email
   - PIN expires in 15 minutes (default)
   ↓
3. User receives email with PIN
   ↓
4. POST /auth/pin/verify
   - Validates PIN against hash
   - Checks PIN expiration
   - Marks user as verified
   - Generates JWT access & refresh tokens
   - Sets auth cookies
   - Clears PIN (one-time use)
```

### Token Management

#### **Refresh Token Flow**
```
POST /auth/refresh
- Validates refresh token
- Checks token expiration
- Generates new access token
- Rotates refresh token (security best practice)
- Updates tokens in database and cookies
```

#### **Logout Flow**
```
POST /auth/logout
- Revokes refresh token from database
- Clears auth cookies
```

### User Management

#### **Get Profile**
```
GET /auth/profile
- Returns authenticated user's profile
- Requires valid access token
```

## Architecture

### Directory Structure

```
auth-service-2/
├── internal/
│   ├── config.go           # Configuration loader
│   ├── models.go           # Data models
│   ├── user_repository.go  # Database operations
│   ├── token_service.go    # JWT token generation/validation
│   ├── email_service.go    # Email sending
│   ├── oauth_handler.go    # OAuth flow handlers
│   ├── pin_handler.go      # PIN authentication handlers
│   └── auth_handler.go     # Token management handlers
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
└── .env.example            # Environment configuration template
```

### Key Components

1. **Configuration (`config.go`)**
   - Loads environment variables
   - Provides typed configuration access

2. **Models (`models.go`)**
   - User model with OAuth and PIN fields
   - Request/Response DTOs
   - Auth response structures

3. **Repository (`user_repository.go`)**
   - MongoDB operations
   - User CRUD operations
   - Token management

4. **Services**
   - **TokenService**: JWT token generation and validation
   - **EmailService**: Email sending via SMTP

5. **Handlers**
   - **OAuthHandler**: Google OAuth flow
   - **PinHandler**: PIN authentication flow
   - **AuthHandler**: Token refresh, logout, profile

## API Endpoints

### OAuth Routes
- `GET /oauth/google` - Initiate Google OAuth login
- `GET /oauth/google/callback` - Handle Google OAuth callback

### PIN Authentication Routes
- `POST /auth/pin/request` - Request authentication PIN
- `POST /auth/pin/verify` - Verify PIN and get tokens

### Token Management Routes
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - Logout (revoke tokens)

### User Routes
- `GET /auth/profile` - Get authenticated user profile

## Security Features

1. **CSRF Protection**
   - OAuth state parameter validation
   - Secure cookie storage

2. **Token Security**
   - JWT with HMAC-SHA256 signing
   - Refresh token rotation
   - Token expiration validation

3. **PIN Security**
   - Bcrypt hashing
   - Time-based expiration
   - One-time use (cleared after verification)

4. **Cookie Security**
   - HTTPOnly cookies
   - Secure flag (production)
   - SameSite protection

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

### Required Configuration

1. **Database**: MongoDB connection string
2. **JWT Secret**: Strong random secret for token signing
3. **Google OAuth**: Client ID and Secret from Google Console
4. **SMTP**: Email server credentials for PIN delivery
5. **URLs**: API and Web app URLs for redirects

## Running the Service

```bash
# Install dependencies
go mod download

# Run the service
go run main.go
```

The service will start on the configured PORT (default: `:8002`).

## Development

### Adding New OAuth Providers

1. Create provider config in `config.go`
2. Add provider handler in `oauth_handler.go`
3. Register routes in `main.go`

### Extending User Model

1. Update `User` struct in `models.go`
2. Add repository methods in `user_repository.go`
3. Update handlers as needed

## Integration with Gateway

This service is designed to work behind an API gateway:

1. Gateway proxies requests to `/auth/*` → `auth-service-2`
2. Gateway validates JWT tokens from cookies
3. Gateway sets `x-user-origin` header for origin validation

## Best Practices

1. **Environment Variables**: Never commit `.env` files
2. **JWT Secret**: Use strong random secrets in production
3. **HTTPS**: Enable `COOKIE_SECURE=true` in production
4. **Token Rotation**: Refresh tokens are rotated on each use
5. **Error Handling**: Detailed errors in dev, generic in production

## License

Proprietary
