# Auth Service

Authentication and authorization service for Instrlabs platform.

## Features

- User authentication (PIN-based, OAuth)
- Google OAuth integration
- JWT token generation and validation
- Refresh token management
- Email verification
- User session management
- Session-Based Token Binding with device validation (IP + User-Agent hash)
- Multiple concurrent sessions per user with per-device revocation

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

### 6. Session-Based Token Binding with Device Validation

**Overview**
This feature prevents token theft by binding JWT tokens to specific users AND devices. Each login creates a unique session tied to the device (identified by IP + User-Agent hash). If a token is used from a different device, the session is immediately deactivated.

**Device Hash Mechanism**
- Device Hash = `SHA256(IP + User-Agent)`
- IP and User-Agent extracted from `x-user-ip` and `x-user-agent` headers
- Headers set by `SetupAuthenticated()` middleware as fiber.Ctx locals
- Each session stores the device hash for validation

**Session Lifecycle**

**Login (PIN or OAuth)**
1. User authenticates successfully
2. System generates unique SessionID
3. Device hash calculated from current IP + User-Agent
4. Session record created in MongoDB with:
   - UserID, SessionID, DeviceHash, IPAddress, UserAgent
   - RefreshToken, IsActive=true, ExpiresAt (7 days default)
5. JWT access token generated with `sessionId` claim
6. Refresh token (long-lived) stored in session document

**Token Refresh**
1. Client sends refresh token with current request headers
2. System finds session by refresh token
3. Validates session is active and not expired
4. **Device validation**: Current device hash compared with stored hash
   - If mismatch: Session deactivated immediately, refresh rejected
   - If match: New access token issued, activity timestamp updated
5. New access token includes session_id claim

**Device Mismatch Detection**
- Occurs during token refresh
- Indicates token may have been compromised or stolen
- Immediate action: Session deactivated to prevent further use
- User must re-authenticate from original device

**Multiple Concurrent Sessions**
- Each device login creates independent session
- Sessions don't interfere with each other
- User can have different sessions on phone, laptop, tablet simultaneously
- Each session can be revoked independently

**Device Management Endpoints**
- `GET /devices` - Lists all active sessions for authenticated user
- `POST /devices/:sessionId/revoke` - Deactivates specific device session
- `POST /devices/revoke-all` - Deactivates all sessions for user

**Data Storage (MongoDB)**
```
users_sessions collection:
{
  _id: ObjectID,
  userId: string,
  sessionId: string,           # Unique per session
  deviceHash: string,          # SHA256(IP + User-Agent)
  ipAddress: string,
  userAgent: string,
  refreshToken: string,        # Hashed refresh token
  isActive: boolean,
  lastActivityAt: timestamp,
  createdAt: timestamp,
  expiresAt: timestamp
}
```

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
│   ├── session.go               # Session domain model and device hashing
│   ├── session_repository.go    # Session database operations
│   └── errors.go                # Error definitions
└── Dockerfile                   # Container definition
```
