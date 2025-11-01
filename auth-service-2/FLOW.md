# Authentication Flows - Visual Guide

This document provides detailed visual representations of the authentication flows in Auth Service 2.

## 1. OAuth Flow (Google)

```
┌─────────┐                                                         ┌─────────┐
│         │                                                         │         │
│  User   │                                                         │  Google │
│         │                                                         │         │
└────┬────┘                                                         └────┬────┘
     │                                                                   │
     │  1. Click "Login with Google"                                    │
     │ ──────────────────────────────────────────────────────►          │
     │                                                                   │
     │                           ┌──────────────────┐                   │
     │  2. GET /oauth/google     │                  │                   │
     │ ─────────────────────────►│  Auth Service 2  │                   │
     │                           │                  │                   │
     │                           │  - Generate      │                   │
     │                           │    OAuth state   │                   │
     │                           │  - Set cookie    │                   │
     │                           │  - Build OAuth   │                   │
     │                           │    URL           │                   │
     │                           └────────┬─────────┘                   │
     │                                    │                             │
     │  3. Redirect to Google OAuth       │                             │
     │ ◄──────────────────────────────────┘                             │
     │                                                                   │
     │  4. Authorize on Google (OAuth consent screen)                   │
     │ ─────────────────────────────────────────────────────────────────►
     │                                                                   │
     │  5. Google validates & generates auth code                       │
     │ ◄─────────────────────────────────────────────────────────────────
     │                                                                   │
     │  6. Redirect to callback with code                               │
     │     GET /oauth/google/callback?code=xxx&state=xxx                │
     │ ──────────────────────────────────────────────────────►          │
     │                           ┌──────────────────┐                   │
     │                           │  Auth Service 2  │                   │
     │                           │                  │                   │
     │                           │  - Verify state  │                   │
     │                           └────────┬─────────┘                   │
     │                                    │                             │
     │                                    │  7. Exchange code for token │
     │                                    │ ───────────────────────────►
     │                                    │                             │
     │                                    │  8. Return access token     │
     │                                    │ ◄───────────────────────────
     │                                    │                             │
     │                                    │  9. Get user info           │
     │                                    │ ───────────────────────────►
     │                                    │                             │
     │                                    │  10. Return user data       │
     │                                    │ ◄───────────────────────────
     │                                    │                             │
     │                           ┌────────┴─────────┐                   │
     │                           │  Auth Service 2  │                   │
     │                           │                  │                   │
     │                           │  - Create/Update │                   │
     │                           │    user          │                   │
     │                           │  - Generate JWT  │                   │
     │                           │    tokens        │                   │
     │                           │  - Set cookies   │                   │
     │                           └────────┬─────────┘                   │
     │                                    │                             │
     │  11. Redirect to web app with cookies                            │
     │ ◄──────────────────────────────────┘                             │
     │                                                                   │
     │  ✅ User is authenticated                                        │
     │                                                                   │
```

### OAuth Flow Steps

1. **User Initiates Login**: User clicks "Login with Google" button
2. **Request OAuth URL**: Frontend calls `GET /oauth/google`
3. **Generate State**: Service generates random state for CSRF protection
4. **Redirect to Google**: User is redirected to Google's OAuth consent page
5. **User Authorizes**: User grants permissions on Google
6. **Callback with Code**: Google redirects back with authorization code
7. **Exchange Code**: Service exchanges code for access token
8. **Fetch User Info**: Service fetches user profile from Google
9. **Create/Update User**: Service creates or updates user in database
10. **Generate Tokens**: Service generates JWT access and refresh tokens
11. **Set Cookies & Redirect**: Service sets auth cookies and redirects to web app

---

## 2. PIN Authentication Flow

```
┌─────────┐
│         │
│  User   │
│         │
└────┬────┘
     │
     │  1. Enter email address
     │ ──────────────────────────────────────────────────────►
     │                           ┌──────────────────┐
     │  2. POST /auth/pin/request│                  │
     │    { "email": "..." }     │  Auth Service 2  │
     │ ─────────────────────────►│                  │
     │                           │  - Generate PIN  │
     │                           │  - Hash PIN      │
     │                           │  - Store in DB   │
     │                           │  - Set expiry    │
     │                           └────────┬─────────┘
     │                                    │
     │                                    │  3. Send PIN email
     │                                    │ ───────────────────►
     │                                    │                     ┌──────────┐
     │                                    │                     │  Email   │
     │  4. Receive email with PIN         │                     │  Service │
     │ ◄──────────────────────────────────┼─────────────────────┤          │
     │                                    │                     └──────────┘
     │  ┌────────────────────────────┐    │
     │  │ Email Content:             │    │
     │  │                            │    │
     │  │ Your PIN: 123456           │    │
     │  │                            │    │
     │  │ Expires in 15 minutes      │    │
     │  └────────────────────────────┘    │
     │                                    │
     │  5. Enter PIN from email           │
     │ ──────────────────────────────────────────────────────►
     │                           ┌──────────────────┐
     │  6. POST /auth/pin/verify │                  │
     │    {                      │  Auth Service 2  │
     │      "email": "...",      │                  │
     │      "pin": "123456"      │  - Find user     │
     │    }                      │  - Verify PIN    │
     │ ─────────────────────────►│  - Check expiry  │
     │                           │  - Mark verified │
     │                           │  - Generate JWT  │
     │                           │  - Clear PIN     │
     │                           └────────┬─────────┘
     │                                    │
     │  7. Return tokens & set cookies    │
     │ ◄──────────────────────────────────┘
     │
     │  {
     │    "access_token": "eyJhbG...",
     │    "refresh_token": "eyJhbG...",
     │    "expires_at": "2025-11-01T12:00:00Z",
     │    "user": { ... }
     │  }
     │
     │  ✅ User is authenticated
     │
```

### PIN Flow Steps

1. **User Enters Email**: User provides email address
2. **Request PIN**: Frontend calls `POST /auth/pin/request`
3. **Generate & Store PIN**:
   - Service generates random 6-digit PIN
   - Hashes PIN using bcrypt
   - Stores hash in database with 15-minute expiry
   - Creates user if doesn't exist
4. **Send Email**: Service sends PIN via SMTP
5. **User Receives PIN**: User checks email for PIN code
6. **Verify PIN**: User submits PIN via `POST /auth/pin/verify`
7. **Validate & Authenticate**:
   - Service validates PIN hash
   - Checks expiration
   - Marks user as verified
   - Generates JWT tokens
   - Clears PIN (one-time use)
8. **Return Tokens**: Service returns tokens and sets auth cookies

---

## 3. Token Refresh Flow

```
┌─────────┐
│         │
│  User   │
│         │
└────┬────┘
     │
     │  Access token expired (401)
     │ ◄──────────────────────────────────────────────────────
     │                           ┌──────────────────┐
     │  1. POST /auth/refresh    │                  │
     │    {                      │  Auth Service 2  │
     │      "refresh_token":     │                  │
     │      "eyJhbG..."          │  - Validate      │
     │    }                      │    refresh token │
     │                           │  - Check expiry  │
     │ ─────────────────────────►│  - Verify user   │
     │                           │  - Generate new  │
     │                           │    tokens        │
     │                           │  - Rotate token  │
     │                           └────────┬─────────┘
     │                                    │
     │  2. Return new tokens              │
     │ ◄──────────────────────────────────┘
     │
     │  {
     │    "access_token": "eyJhbG...",  (new)
     │    "refresh_token": "eyJhbG...", (rotated)
     │    "expires_at": "2025-11-01T13:00:00Z",
     │    "user": { ... }
     │  }
     │
     │  ✅ User has new access token
     │
```

### Token Refresh Steps

1. **Expired Token**: Access token expires, API returns 401
2. **Request Refresh**: Frontend calls `POST /auth/refresh` with refresh token
3. **Validate Refresh Token**:
   - Service validates JWT signature
   - Checks token expiration
   - Verifies token exists in database
4. **Generate New Tokens**:
   - Creates new access token
   - Rotates refresh token (security best practice)
   - Updates database
5. **Return Tokens**: Service returns new tokens and updates cookies

---

## 4. Logout Flow

```
┌─────────┐
│         │
│  User   │
│         │
└────┬────┘
     │
     │  1. Click "Logout"
     │ ──────────────────────────────────────────────────────►
     │                           ┌──────────────────┐
     │  2. POST /auth/logout     │                  │
     │                           │  Auth Service 2  │
     │ ─────────────────────────►│                  │
     │                           │  - Get user ID   │
     │                           │    from token    │
     │                           │  - Clear refresh │
     │                           │    token in DB   │
     │                           │  - Clear cookies │
     │                           └────────┬─────────┘
     │                                    │
     │  3. Confirmation                   │
     │ ◄──────────────────────────────────┘
     │
     │  {
     │    "message": "Logged out successfully"
     │  }
     │
     │  ✅ User is logged out
     │
```

### Logout Steps

1. **User Clicks Logout**: User initiates logout
2. **Request Logout**: Frontend calls `POST /auth/logout`
3. **Revoke Tokens**:
   - Service extracts user ID from access token
   - Removes refresh token from database
   - Clears auth cookies
4. **Confirm Logout**: Service returns success message

---

## Security Considerations

### OAuth Flow
- **CSRF Protection**: Random state parameter validated on callback
- **Token Validation**: Google tokens validated before user creation
- **Email Verification**: Only verified Google emails accepted

### PIN Flow
- **Hashing**: PINs hashed with bcrypt before storage
- **Expiration**: PINs expire after 15 minutes
- **One-Time Use**: PINs cleared after successful verification
- **Rate Limiting**: Should be implemented at gateway level

### Token Management
- **JWT Signing**: Tokens signed with HMAC-SHA256
- **Token Rotation**: Refresh tokens rotated on each use
- **Secure Cookies**: HTTPOnly, Secure, SameSite attributes
- **Expiration**: Access tokens (1h), Refresh tokens (30 days)

### General
- **HTTPS**: All production traffic over HTTPS
- **Environment Secrets**: Sensitive config in environment variables
- **Database Security**: Passwords and tokens never stored in plain text
