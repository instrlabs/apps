## POST /login - Complete Flow

### Request Format
```json
{
  "email": "string (required)",
  "pin": "string (required)"
}
```

### Success Response Flow (HTTP 200)
**Condition**: All validations pass, user exists, PIN matches, session created successfully

```json
{
  "message": "Login successful",
  "errors": null,
  "data": {
    "access_token": "JWT with user_id + session_id claims",
    "refresh_token": "base64 random token (32 bytes)",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### Error Response Flows

#### 1. Invalid Request Body (HTTP 400)
**Trigger**: Malformed JSON or unparsable request body
**Code Location**: user_handler.go:79-86
```json
{
  "message": "Invalid request body",
  "errors": null,
  "data": null
}
```

#### 2. Missing Email Field (HTTP 400)
**Trigger**: `input.Email == ""` (empty string or missing field)
**Code Location**: user_handler.go:88-95
```json
{
  "message": "Email is required",
  "errors": null,
  "data": null
}
```

#### 3. Missing PIN Field (HTTP 400)
**Trigger**: `input.Pin == ""` (empty string or missing field)
**Code Location**: user_handler.go:97-104
```json
{
  "message": "Pin is required",
  "errors": null,
  "data": null
}
```

#### 4. User Not Found (HTTP 400)
**Trigger**: `user == nil || user.ID.IsZero()` (no user with provided email)
**Code Location**: user_handler.go:108-114
```json
{
  "message": "Invalid email or pin",
  "errors": null,
  "data": null
}
```

#### 5. Invalid PIN (HTTP 400)
**Trigger**: `!user.ComparePin(input.Pin)` (PIN doesn't match hash or expired)
**Code Location**: user_handler.go:116-123
```json
{
  "message": "Invalid email or pin",
  "errors": null,
  "data": null
}
```

