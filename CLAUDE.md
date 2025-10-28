# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based application with a Next.js frontend and four Go backend services communicating via NATS message broker. All services use Fiber web framework and share common infrastructure through `github.com/instrlabs/shared`.

### Services

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
