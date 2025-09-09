# Gateway Service

A Go (Fiber) API Gateway that fronts multiple backend services (Auth, Image, Payment). It handles CORS, logs inbound requests, authenticates requests via an HTTP-only AccessToken cookie, injects identity headers for downstream services, provides a consolidated Swagger UI, and proxies requests to registered services.

## Overview
- Framework: Fiber v2
- Responsibilities:
  - Inbound request logging
  - CORS policy enforcement
  - Access token cookie parsing and validation
  - Inject X-Authenticated, X-User-Id, X-User-Roles headers
  - Strip cookies before proxying
  - Reverse proxy by service prefix (e.g., /auth, /image, /payment)
  - Aggregate health and Swagger UI

## Features
- Cookie-based authentication parsing (no cookies forwarded downstream)
- JWT extraction and validation (see internal/middleware.go ExtractTokenInfo usage)
- Identity propagation via headers: X-Authenticated, X-User-Id, X-User-Roles, X-Gateway
- Reverse proxy with per-service prefix routing
- Health endpoint that pings downstream /health endpoints
- Single Swagger UI that lists multiple service OpenAPI JSON docs

## Routes
- GET /health
  - Returns overall status and per-service status by calling each configured service /health endpoint.
- GET /swagger
  - Serves a Swagger UI HTML that lists multiple APIs (Auth, Payment, Image, plus an example Petstore).
- Proxy routes (configured via prefixes):
  - /auth/* -> AUTH_SERVICE_URL
  - /payment/* -> PAYMENT_SERVICE_URL
  - /image/* -> IMAGE_SERVICE_URL

## Authentication Behavior
- Reads AccessToken from cookie on inbound request.
- Valid token -> sets headers:
  - X-Authenticated: "true"
  - X-User-Id: <user id>
  - X-User-Roles: comma-separated roles
- Invalid or absent token -> X-Authenticated: "false"
- If token is expired, responds 401 { "error": "EXPIRED_TOKEN" }
- Strips Cookie header before proxying to downstream services.
- Adds X-Gateway: "true" header on all forwarded requests.

## Configuration
Environment variables (see internal/config.go and internal/middleware.go):
- PORT: listen address (default: 3000)
- AUTH_SERVICE_URL: e.g., http://auth.localhost:8080
- PAYMENT_SERVICE_URL: e.g., http://payment.localhost:8082
- IMAGE_SERVICE_URL: e.g., http://image.localhost:8081
- CORS_ALLOWED_ORIGINS: comma-separated origins for browsers

Note: Services are registered with names, prefixes, and base URLs in the config. The gateway logs all proxy decisions.

## Running Locally
1. Create a .env containing PORT and service URLs (and CORS_ALLOWED_ORIGINS if needed).
2. Run: `go run main.go`
3. Check health: GET /health
4. Open Swagger UI: GET /swagger

## Logging
- JSON-structured logs via logrus.
- Logs inbound requests (method, path, ip, user-agent, duration).
- Logs proxy events and unmatched routes.

## Integration Notes for Downstream Services
- Do not rely on cookies; the gateway strips them.
- Trust the following headers for authentication state:
  - X-Authenticated: "true" or "false"
  - X-User-Id: user identifier
  - X-User-Roles: roles list
- Optionally verify X-Gateway === "true" if needed.
