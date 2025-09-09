# Auth Service

A Go (Fiber) microservice providing user authentication and profile management. It issues JWT access tokens (via the Gateway cookie), manages refresh tokens, supports password reset via email, and offers Google OAuth login. Includes health check and Swagger JSON endpoint.

## Overview
- Framework: Fiber v2
- Storage: MongoDB
- Auth: JWT access token (validated by Gateway), HTTP-only cookies handled at the Gateway, refresh tokens persisted per user
- OAuth: Google
- Email: SMTP for password reset
- Observability: Request logging middleware

## Features
- User registration with hashed password (bcrypt)
- Login with JWT access token issuance (for Gateway) and refresh token persistence
- Logout (refresh token invalidation)
- Access token refresh via refresh token
- Forgot password and reset password flow (email + time-bound token)
- Get/Update user profile
- Change password (requires current password)
- Google OAuth login and callback
- Public Swagger JSON and health endpoints

## Endpoints
Public endpoints (no auth required by this service):
- GET /health → { status: "ok" }
- GET /swagger → serves static/swagger.json
- POST /register
- POST /login
- POST /refresh
- POST /forgot-password
- POST /reset-password
- GET /google (Google OAuth begin)
- GET /google/callback (Google OAuth callback)

Protected endpoints (require X-Authenticated: true from Gateway):
- GET /profile
- PUT /profile
- POST /change-password
- POST /logout

Response envelope (typical):
- message: string
- errors: null | string | object
- data: null | object

## Auth and Gateway Integration
This service relies on the Gateway to perform access-token parsing and to set the following headers when a user is authenticated:
- X-Authenticated: "true" | "false"
- X-User-Id: user ID (Mongo ObjectID hex)
- X-User-Roles: comma-separated roles (if any)

Middleware in this service treats these paths as public:
/health, /swagger, /login, /refresh, /register, /forgot-password, /reset-password, /google (and its callback).
Non-public paths require X-Authenticated: true; otherwise 401 Unauthorized is returned.

Access token lifetime is governed by TOKEN_EXPIRY_HOURS; expired/invalid access tokens are typically surfaced by the Gateway before requests reach this service.

## Password Reset Flow
- Client calls POST /forgot-password with email -> service sends email with reset token link
- Client calls POST /reset-password with token + new password, before ResetTokenExpires

## Google OAuth
- Start: GET /google -> redirects to Google
- Callback: GET /google/callback -> creates/links user and issues tokens

## Swagger
- GET /swagger returns static/swagger.json for API schema
- Gateway may aggregate multiple service swagger UIs under its /swagger page

## Configuration
Environment variables (see internal/config.go):
- ENVIRONMENT: deployment environment name
- PORT: service listen address, e.g. ":8080"
- MONGO_URI: MongoDB connection string
- MONGO_DB: MongoDB database name
- JWT_SECRET: secret for signing JWT access tokens
- TOKEN_EXPIRY_HOURS: int, access token expiry in hours
- SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD: SMTP settings for email
- EMAIL_FROM: From email address
- RESET_TOKEN_EXPIRY_HOURS: int, reset token expiry in hours
- GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL: Google OAuth credentials
- FE_RESET_PASSWORD: Frontend URL to handle password reset
- FE_OAUTH_REDIRECT: Frontend URL to redirect after OAuth
- CORS_ALLOWED_ORIGINS: allowed origins for CORS (defaults to http://web.localhost if empty)
- COOKIE_DOMAIN: cookie domain used by Gateway (defaults to .localhost if empty)

## Running Locally
1. Create a .env with the variables above (at minimum PORT, MONGO_URI, MONGO_DB, JWT_SECRET).
2. Run: go run main.go
3. Health: GET /health
4. Swagger JSON: GET /swagger

## Data Model (User)
Important fields (see internal/user.go):
- id (ObjectID)
- name, email
- password (bcrypt hash)
- google_id (optional)
- refresh_token (optional)
- reset_token, reset_token_expires (for password reset)
- created_at, updated_at

## Error Handling
- 400 for invalid input
- 401 for unauthorized (when Gateway indicates unauthenticated for protected routes)
- 404 when resource not found
- 500 for internal errors

## Notes
- Gateway removes cookies before proxying and communicates authentication via headers.
- This service trusts Gateway headers only for authentication and does not parse cookies directly.
