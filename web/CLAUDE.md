# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

**Development server:**
```bash
npm run dev
# Runs Next.js dev server with Turbopack on port 8000
```

**Production build:**
```bash
npm run build
npm start  # Runs production server on port 8000
```

**Linting and formatting:**
```bash
npm run lint          # Run ESLint
npm run format        # Format code with Prettier
npm run format:check  # Check formatting without writing
```

## Architecture Overview

### Tech Stack
- **Framework:** Next.js 15.4.7 with App Router (React 19)
- **Styling:** Tailwind CSS v4 with Geist font
- **Type Safety:** TypeScript with strict mode
- **Forms:** react-hook-form

### Route Groups and Layout Structure

The application uses Next.js route groups to organize pages with different layouts:

1. **`(non-auth)/`** - Public pages without authentication (e.g., login, register)
2. **`(site)/`** - Authenticated pages with full layout including navigation, overlays, and context providers
3. **`debug/`** - Development/testing pages
4. **`api/`** - API routes (e.g., SSE proxy at `/api/sse`)

The authenticated site layout (`app/(site)/layout.tsx`) wraps all authenticated pages in a nested provider hierarchy:
```
ProfileProvider → ProductProvider → SSEProvider → NotificationProvider →
ModalProvider → OverlayProvider → OverlayTop + OverlayContent + NotificationWidget
```

### Authentication Flow

**Middleware-based authentication** (`middleware.ts:5-69`):
- All routes except whitelisted paths (`/login`, `/register`, `/forgot-password`, `/reset-password`, `/`) require authentication
- Automatically attempts token refresh if `access_token` is missing but `refresh_token` exists
- Forwards client metadata headers (`x-user-ip`, `x-user-agent`, `x-user-host`, `x-user-origin`) to the gateway
- Redirects to `/login` on authentication failure

**Token management:**
- `access_token` and `refresh_token` stored in HTTP-only cookies
- Tokens automatically set by `fetchPOST()` when calling `/auth/login` or `/auth/refresh` endpoints
- Auth functions in `services/auth.ts` must only be called from server components/actions

### API Integration Pattern

**Server-side fetch utilities** (`utils/fetch.ts`):
- All API calls go through `GATEWAY_URL` environment variable
- `fetchGET()`, `fetchPOST()`, `fetchPUT()`, `fetchPATCH()` - JSON endpoints
- `fetchGETBytes()` - Binary data (e.g., images)
- `fetchPOSTFormData()` - Form uploads with multipart/form-data

**Response structure:**
```typescript
{
  success: boolean;
  message: string;
  data: T | null;
  errors: FormErrors | null;
}
```

**Important:** All fetch utilities are server actions (`"use server"`). They automatically:
- Forward client metadata headers from middleware
- Handle cookie-based authentication
- Set tokens on successful login/refresh

### Overlay System

**Centralized overlay management** (`hooks/useOverlay.tsx`):
- Overlays slide in from left or right sides of the screen
- Registration-based system with key-based lookup (e.g., `"left:navigation"`, `"right:notifications"`, `"right:profile"`)
- Actions: `openLeft(key)`, `closeLeft()`, `openRight(key)`, `closeRight()`
- Register new overlays in `registerOverlays()` function

### Component Organization

Components are organized by function in subdirectories:
- `components/actions/` - Interactive buttons and actions
- `components/cards/` - Card layouts and containers
- `components/feedback/` - Loading states, alerts, toasts
- `components/icons/` - SVG icon components
- `components/inputs/` - Form inputs and controls
- `components/layouts/` - Page layout components (OverlayTop, OverlayContent)
- `components/navigation/` - Navigation components
- `components/overlays/` - Overlay content (navigation, notifications, profile)

### Custom Hooks

- `useOverlay` - Control left/right sliding overlays
- `useModal` - Modal dialog state management
- `useNotification` - Toast/notification system
- `useSSE` - Server-Sent Events connection management
- `useProduct` - Access product list from context
- `useProfile` - Access user profile from context
- `useMediaQuery` - Responsive breakpoint detection

### Environment Variables

Required environment variables (see `.env.example`):
- `NODE_ENV` - Environment (development/production)
- `GATEWAY_URL` - Backend API gateway base URL
- `NOTIFICATION_URL` - SSE notification service URL (server-side only)

### Build Configuration

**Next.js config** (`next.config.ts`):
- Standalone output mode for Docker deployment
- Server actions with 5MB body size limit
- Console logs stripped in production (except warn/error)

**Path aliases:**
- `@/*` maps to project root (configured in `tsconfig.json`)

### SSE Proxy Pattern

The `/api/sse` route proxies Server-Sent Events from the private `NOTIFICATION_URL` service:
- Forwards cookies and origin headers to upstream
- Streams responses directly to client
- Handles connection errors and client disconnects
- Required because notification service is not publicly accessible

### Styling Guidelines

- Use native array filtering for conditional classes instead of external libraries
- Tailwind CSS with automatic class sorting via Prettier
- Geist font family applied globally
- No inline styles unless absolutely necessary
