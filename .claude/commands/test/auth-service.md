# Test Auth Service Features

Test the authentication service features according to feature workflows in `./auth-service/README.md`.

## Testing Scope

Test the following authentication flows:

1. **PIN Login Flow**
   - Send PIN to email
   - Login with PIN
   - Access protected endpoints
   - Logout
   - Verify access is denied after logout

2. **Google OAuth Flow**
   - Initiate Google login
   - Handle OAuth callback
   - Verify user creation
   - Access protected endpoints
   - Logout

3. **Token Management**
   - Verify token expiry times
   - Test token refresh
   - Verify cookie domain settings

## Service Architecture

- Services run in Docker containers
- Gateway service proxies requests to auth service
- Auth service endpoints are accessible at `$GATEWAY_URL/auth/*`
- Gateway handles CORS, CSRF protection, and authentication middleware

### Request Flow

**Production Flow (Browser → Traefik → Next.js → Gateway → Services):**
1. Browser sends request with standard `Origin` header
2. Traefik adds `x-forwarded-proto` and `x-forwarded-host` headers
3. Next.js middleware constructs `x-user-origin = x-forwarded-proto + "://" + x-forwarded-host`
4. Next.js forwards `x-user-origin` to gateway with all requests
5. Gateway validates `x-user-origin` for CSRF protection
6. Backend services use `x-user-origin` for origin-based logic

## Testing Requirements

### Environment
- Use the Gateway URL from `GATEWAY_URL` environment variable (default: `http://localhost:3000`)
- Auth service is proxied at `/auth` through the gateway
- Test with `PIN_ENABLED=false` (PIN should be "000000")

### CSRF Protection
- Gateway enforces CSRF protection when `CSRF_ENABLED=true`
- **Backend validates `x-user-origin` header** (NOT the standard `Origin` header)
- `x-user-origin` is automatically set by Next.js middleware in production
- Allowed origins are defined in `gateway-service/.env` under `ORIGINS_ALLOWED`
- GET and OPTIONS requests bypass CSRF checks

**For Testing Only:**
- When testing directly against the gateway (bypassing Next.js), you must manually set:
  - `x-user-origin: http://localhost:8000` (custom header that Next.js would normally set)

### Logging
- Auth service logs to Docker stdout (not file-based)
- View logs with: `docker logs instrlabs-auth-service --tail 50`
- Logs include timestamps and operation details

### Authentication
- Verify HTTP-only cookies are set correctly
- Test all endpoints for proper authentication
- Cookies should work across requests automatically

## Expected Behavior

### Token Management
- Access tokens expire after 1 hour (configurable via `TOKEN_EXPIRY_HOURS`)
- Refresh tokens expire after 30 days (configurable via `REFRESH_EXPIRY_HOURS`)
- Tokens are stored in HTTP-only cookies

### Cookie Configuration
- Cookie domain is `.localhost` in development
- In production, cookie domain is `.arthadede.com`
- Cookies are `SameSite=None` and `Secure` in production
- All authentication flows should set cookies automatically

### Endpoints
- Protected endpoints should work with valid cookies
- Protected endpoints should return 401 Unauthorized without valid cookies
- Logout should clear all authentication cookies
- After logout, cookies are expired (MaxAge: -1, Expires: Unix epoch)

## Test Implementation

Create comprehensive tests or manual test scripts to verify all features work as documented.

**IMPORTANT:** All test scripts should be created in `/tmp/` directory to keep them temporary and out of version control.

### Test Script Requirements
1. **Create scripts in `/tmp/`** - Test scripts should be temporary and not committed to repository
2. Use cookie files in `/tmp/` to persist authentication across requests
3. Verify HTTP status codes and response messages
4. Test the complete flow from authentication to logout
5. Verify protected endpoints are inaccessible after logout

### Example Test Flow

**Note:** These examples include `x-user-origin` header to simulate the production flow where Next.js middleware automatically sets this header from Traefik's forwarded headers.

### Why Two Origin Headers in Tests?

- **`Origin`** - Standard browser header (automatically sent by browsers)
- **`x-user-origin`** - Custom header set by Next.js middleware from Traefik's `x-forwarded-*` headers
- In production, Next.js sets `x-user-origin` automatically; **frontend code never sets it manually**
- In tests, we simulate both headers to bypass Next.js middleware and test the gateway directly
- **All backend services validate `x-user-origin` only** (NOT the standard `Origin` header)
