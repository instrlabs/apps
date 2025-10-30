# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.


## Services

**Gateway Service** (`gateway-service/`)
- Handle proxies requests to downstream services
- Handle Middlewares (CORS, Rate Limit, CSRF Protection, and Authentication)

**Auth Service** (`auth-service/`)
- Handle Authentication and Authorization
- Handle OAuth & JWT Flows

**Image Service** (`image-service/`)
- Handle Image processing requests (via NATS, API, Scheduling)

**Notification Service** (`notification-service/`)
- Handle Server-Sent Events (SSE) for real-time notifications to web clients

**Web Application** (`web/`)
- Next.js 15 with App Router and React 19

## Environment Configuration

Environment variables are separated into two categories for clarity:

### Public Exposed URLs (`*_URL`)
These are services accessible from the frontend and external clients:
- `WEB_URL` - Next.js web application (e.g., `http://localhost:3000`)
- `API_URL` - Backend API gateway (e.g., `http://localhost:3001`)
- `NOTIFICATION_URL` - SSE notification service (e.g., `http://localhost:3002`)

**Used by:** Web frontend and browser-based clients

### Internal Service URLs (`*_SERVICE`)
These are services that can only be accessed within the Docker network:
- `AUTH_SERVICE` - Authentication service (e.g., `http://auth-service:3000`)
- `IMAGE_SERVICE` - Image processing service (e.g., `http://image-service:3000`)
- `GATEWAY_SERVICE` - Backend API gateway service (e.g., `http://gateway-service:3000`)
- `NOTIFICATION_SERVICE` - SSE notification service (e.g., `http://notification-service:3000`)

**Used by:** Backend services (Gateway, Notification) for internal communication

### Configuration Files
- **Root `.env`** - Contains all public and internal URL definitions
- **Service `.env` files** - Reference root variables using `${VAR_NAME}` syntax
  - `gateway-service/.env` - Uses `AUTH_SERVICE` and `IMAGE_SERVICE`
  - `auth-service/.env` - Uses `API_URL` and `WEB_URL`
  - `notification-service/.env` - Uses only public `ORIGINS_ALLOWED`
  - `web/.env.local` - Uses `API_URL` and `NOTIFICATION_URL`

See `.env.example` for complete configuration template.

## Code Patterns

### Service Structure
```
<service>/
├── Dockerfile
├── go.mod
├── main.go          # Entry point with setup
├── internal/
│   ├── config.go    # Environment config loader
│   ├── *_handler.go # HTTP handlers (Fiber)
│   ├── *_repository.go # Database access
│   └── *.go         # Domain models
└── static/          # Static assets (if any)
```

### Origin Header Pattern (`x-user-origin`)

**IMPORTANT:** All backend services use `x-user-origin` for origin validation, NOT the standard `Origin` header.

**Key Points:**
- Frontend code NEVER manually sets `x-user-origin` - it's automatic via Next.js middleware
- All Go services check `c.Get("x-user-origin")` for origin validation
- For testing, manually set `x-user-origin` to simulate production flow
- Allowed origins configured in `gateway-service/.env` under `ORIGINS_ALLOWED`

## Guidelines
- Always get environment configuration from `.env`
- If Testing always refers to the .ai-guides directory in the service folder
- **Use `x-user-origin` header for all origin validation in backend services**
- **All test scripts must be created in `/tmp/` directory** - Keep test files temporary and out of version control