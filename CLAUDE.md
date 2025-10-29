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

- Each service requires its own `.env` file. 
- And root `.env` file for shared environment variables.
- See `.env.example` for example configuration.

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