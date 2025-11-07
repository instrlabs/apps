# Midtrans Payment Gateway Security Best Practices

## Overview

This document outlines security best practices for integrating Midtrans payment gateway into the InstrLabs microservices ecosystem. It covers authentication, data protection, key management, and compliance requirements.

## Security Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client App    │    │  Gateway Service│    │  Payment Service│
│                 │    │                 │    │                 │
│ - HTTPS/TLS 1.3 │◄──►│ - JWT Validation│◄──►│ - API Keys      │
│ - CSP Headers   │    │ - Rate Limiting │    │ - Encryption    │
│ - Token Storage │    │ - CORS          │    │ - Audit Logs    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Midtrans API  │    │    Database     │
                       │                 │    │                 │
                       │ - HTTPS/TLS 1.3 │    │ - Encrypted     │
                       │ - Basic Auth    │    │ - Access Control│
                       │ - Webhook Sigs  │    │ - Backups       │
                       └─────────────────┘    └─────────────────┘
```

## Authentication and Authorization

### 1. JWT Token Management

#### Token Structure
```go
type JWTPayload struct {
    UserID      string    `json:"user_id"`
    Email       string    `json:"email"`
    Role        string    `json:"role"`
    Permissions []string  `json:"permissions"`
    Iat         int64     `json:"iat"`
    Exp         int64     `json:"exp"`
    Iss         string    `json:"iss"`
    Aud         string    `json:"aud"`
}
```

#### Token Validation Middleware
```go
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Authorization header required",
                "errors": []fiber.Map{{
                    "field": "authorization",
                    "code":  "MISSING",
                    "message": "Authorization header is required",
                }},
            })
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid authorization format",
                "errors": []fiber.Map{{
                    "field": "authorization",
                    "code":  "INVALID_FORMAT",
                    "message": "Authorization header must be in format: Bearer <token>",
                }},
            })
        }

        token, err := jwt.ParseWithClaims(tokenString, &JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid or expired token",
                "errors": []fiber.Map{{
                    "field": "authorization",
                    "code":  "INVALID_TOKEN",
                    "message": "Token is invalid or expired",
                }},
            })
        }

        claims, ok := token.Claims.(*JWTPayload)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid token claims",
                "errors": []fiber.Map{{
                    "field": "authorization",
                    "code":  "INVALID_CLAIMS",
                    "message": "Token claims are invalid",
                }},
            })
        }

        // Store user context
        c.Locals("user_id", claims.UserID)
        c.Locals("user_email", claims.Email)
        c.Locals("user_role", claims.Role)
        c.Locals("permissions", claims.Permissions)

        return c.Next()
    }
}
```

### 2. Midtrans API Authentication

#### Secure Key Storage
```go
type MidtransConfig struct {
    ServerKey string `json:"server_key"`
    ClientKey string `json:"client_key"`
    Environment string `json:"environment"`
}

// Load from encrypted storage or environment variables
func loadMidtransConfig() (*MidtransConfig, error) {
    encryptedKey := os.Getenv("MIDTRANS_SERVER_KEY_ENCRYPTED")
    if encryptedKey == "" {
        return nil, errors.New("encrypted server key not found")
    }

    // Decrypt the key (implementation depends on your encryption method)
    serverKey, err := decrypt(encryptedKey, getEncryptionKey())
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt server key: %w", err)
    }

    return &MidtransConfig{
        ServerKey:   serverKey,
        ClientKey:   os.Getenv("MIDTRANS_CLIENT_KEY"),
        Environment: os.Getenv("MIDTRANS_ENVIRONMENT"),
    }, nil
}
```

#### API Client with Authentication
```go
type SecureMidtransClient struct {
    client      *http.Client
    serverKey   string
    environment string
    rateLimiter *rate.Limiter
}

func NewSecureMidtransClient(serverKey, environment string) *SecureMidtransClient {
    return &SecureMidtransClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    MinVersion: tls.VersionTLS12,
                    MaxVersion: tls.VersionTLS13,
                },
            },
        },
        serverKey:   serverKey,
        environment: environment,
        rateLimiter: rate.NewLimiter(rate.Limit(100), 10), // 100 requests per second
    }
}

func (c *SecureMidtransClient) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
    // Rate limiting
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limit exceeded: %w", err)
    }

    var bodyReader io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request: %w", err)
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Set headers
    req.Header.Set("Accept", "application/json")
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    // Basic authentication with server key
    auth := c.serverKey + ":"
    encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
    req.Header.Set("Authorization", "Basic "+encodedAuth)

    // Add request ID for traceability
    req.Header.Set("X-Request-ID", generateRequestID())

    return c.client.Do(req)
}
```

## Data Protection

### 1. Sensitive Data Handling

#### Card Data Tokenization
```go
// Never store raw card data
type PaymentRequest struct {
    OrderID       string          `json:"order_id" validate:"required"`
    TokenID       string          `json:"token_id" validate:"required"` // Midtrans token, not raw card
    GrossAmount   int64           `json:"gross_amount" validate:"required,min=1000"`
    CustomerID    string          `json:"customer_id" validate:"required"`
}

// Validate token before processing
func (s *PaymentService) validateToken(tokenID string) error {
    // Verify token with Midtrans API
    tokenInfo, err := s.midtransClient.GetTokenInfo(tokenID)
    if err != nil {
        return fmt.Errorf("invalid token: %w", err)
    }

    // Check token expiry
    if time.Now().After(tokenInfo.ExpiresAt) {
        return errors.New("token has expired")
    }

    return nil
}
```

#### Data Encryption at Rest
```go
// Encrypt sensitive fields before storing
type EncryptedPaymentData struct {
    OrderID       string `json:"order_id"`
    EncryptedData string `json:"encrypted_data"` // Encrypted JSON
    Nonce         string `json:"nonce"`          // Encryption nonce
    Tag           string `json:"tag"`            // Authentication tag
}

func encryptPaymentData(data interface{}, key []byte) (*EncryptedPaymentData, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal data: %w", err)
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nil, nonce, jsonData, nil)

    return &EncryptedPaymentData{
        EncryptedData: base64.StdEncoding.EncodeToString(ciphertext),
        Nonce:         base64.StdEncoding.EncodeToString(nonce),
        Tag:           "", // Included in ciphertext for GCM
    }, nil
}
```

### 2. Database Security

#### Connection Security
```go
func secureMongoConnection(uri string) (*mongo.Client, error) {
    clientOptions := options.Client().
        ApplyURI(uri).
        SetTLSConfig(&tls.Config{
            MinVersion: tls.VersionTLS12,
        }).
        SetAuth(options.Credential{
            AuthMechanism: "SCRAM-SHA-256",
        }).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second)

    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return client, nil
}
```

#### Field-Level Encryption
```go
type PaymentRecord struct {
    ID              primitive.ObjectID `bson:"_id,omitempty"`
    OrderID         string             `bson:"order_id"`
    UserID          string             `bson:"user_id"`
    TransactionID   string             `bson:"transaction_id"`

    // Encrypted fields
    EncryptedCard   *EncryptedField    `bson:"encrypted_card,omitempty"`

    // Non-sensitive fields
    Status          string             `bson:"status"`
    Amount          int64              `bson:"amount"`
    CreatedAt       time.Time          `bson:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at"`
}

type EncryptedField struct {
    Data string `bson:"data"`
    IV   string `bson:"iv"`
    Tag  string `bson:"tag"`
}
```

## Webhook Security

### 1. Signature Verification

#### SHA512 Signature Implementation
```go
func VerifyWebhookSignature(orderID, statusCode, grossAmount, receivedSignature, serverKey string) bool {
    // Construct the input string
    input := orderID + statusCode + grossAmount + serverKey

    // Calculate SHA512 hash
    hash := sha512.Sum512([]byte(input))
    calculatedSignature := fmt.Sprintf("%x", hash)

    // Use constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare([]byte(calculatedSignature), []byte(receivedSignature)) == 1
}

// Webhook middleware
func WebhookMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get signature from header
        signature := c.Get("X-Midtrans-Signature")
        if signature == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Missing signature",
                "errors": []fiber.Map{{
                    "code":    "MISSING_SIGNATURE",
                    "message": "X-Midtrans-Signature header is required",
                }},
            })
        }

        // Read body
        body := c.Body()
        if len(body) == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "message": "Empty request body",
                "errors": []fiber.Map{{
                    "code":    "EMPTY_BODY",
                    "message": "Request body cannot be empty",
                }},
            })
        }

        // Parse notification
        var notification midtrans.WebhookNotification
        if err := json.Unmarshal(body, &notification); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "message": "Invalid JSON",
                "errors": []fiber.Map{{
                    "code":    "INVALID_JSON",
                    "message": "Request body contains invalid JSON",
                }},
            })
        }

        // Verify signature
        serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
        if !VerifyWebhookSignature(
            notification.OrderID,
            notification.StatusCode,
            notification.GrossAmount,
            signature,
            serverKey,
        ) {
            // Log security incident
            logSecurityIncident("webhook_signature_verification_failed", map[string]interface{}{
                "order_id":     notification.OrderID,
                "remote_addr":  c.IP(),
                "user_agent":   c.Get("User-Agent"),
                "signature":    signature,
            })

            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Invalid signature",
                "errors": []fiber.Map{{
                    "code":    "INVALID_SIGNATURE",
                    "message": "Webhook signature verification failed",
                }},
            })
        }

        // Store verified notification in context
        c.Locals("webhook_notification", &notification)

        return c.Next()
    }
}
```

### 2. Replay Attack Prevention

```go
type WebhookNonceStore interface {
    SetUsed(ctx context.Context, orderID, nonce string, expiry time.Duration) error
    IsUsed(ctx context.Context, orderID, nonce string) (bool, error)
    Cleanup(ctx context.Context) error
}

func preventReplayAttack(store WebhookNonceStore) fiber.Handler {
    return func(c *fiber.Ctx) error {
        notification := c.Locals("webhook_notification").(*midtrans.WebhookNotification)

        // Create nonce from order_id + transaction_time
        nonce := fmt.Sprintf("%s_%s", notification.OrderID, notification.TransactionTime)

        // Check if this webhook was already processed
        used, err := store.IsUsed(c.Context(), notification.OrderID, nonce)
        if err != nil {
            log.Errorf("Failed to check webhook nonce: %v", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Internal server error",
            })
        }

        if used {
            log.Warnf("Duplicate webhook received: order_id=%s, nonce=%s", notification.OrderID, nonce)
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "message": "Duplicate webhook",
                "errors": []fiber.Map{{
                    "code":    "DUPLICATE_WEBHOOK",
                    "message": "This webhook has already been processed",
                }},
            })
        }

        // Mark this webhook as processed
        err = store.SetUsed(c.Context(), notification.OrderID, nonce, 24*time.Hour)
        if err != nil {
            log.Errorf("Failed to mark webhook as used: %v", err)
        }

        return c.Next()
    }
}
```

## Key Management

### 1. Environment Variable Encryption

#### Encrypted Configuration Storage
```go
type SecureConfig struct {
    encrypted map[string]string
    key       []byte
}

func NewSecureConfig(masterKey string) *SecureConfig {
    key := sha256.Sum256([]byte(masterKey))
    return &SecureConfig{
        encrypted: make(map[string]string),
        key:       key[:],
    }
}

func (c *SecureConfig) SetEncrypted(key, value string) error {
    block, err := aes.NewCipher(c.key)
    if err != nil {
        return fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
    c.encrypted[key] = base64.StdEncoding.EncodeToString(ciphertext)

    return nil
}

func (c *SecureConfig) GetDecrypted(key string) (string, error) {
    encryptedValue, exists := c.encrypted[key]
    if !exists {
        return "", fmt.Errorf("key not found: %s", key)
    }

    data, err := base64.StdEncoding.DecodeString(encryptedValue)
    if err != nil {
        return "", fmt.Errorf("failed to decode base64: %w", err)
    }

    block, err := aes.NewCipher(c.key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt: %w", err)
    }

    return string(plaintext), nil
}
```

### 2. Key Rotation Strategy

```go
type KeyRotationManager struct {
    currentKeyID string
    keys         map[string][]byte
    keyStore     KeyStore
}

type KeyStore interface {
    StoreKey(ctx context.Context, keyID string, key []byte, expiresAt time.Time) error
    GetKey(ctx context.Context, keyID string) ([]byte, error)
    ListKeys(ctx context.Context) ([]KeyMetadata, error)
    RevokeKey(ctx context.Context, keyID string) error
}

type KeyMetadata struct {
    KeyID      string    `json:"key_id"`
    CreatedAt  time.Time `json:"created_at"`
    ExpiresAt  time.Time `json:"expires_at"`
    IsActive   bool      `json:"is_active"`
    KeyType    string    `json:"key_type"`
}

func (km *KeyRotationManager) RotateKeys(ctx context.Context) error {
    // Generate new key
    newKeyID := fmt.Sprintf("key_%d", time.Now().Unix())
    newKey := make([]byte, 32)
    if _, err := rand.Read(newKey); err != nil {
        return fmt.Errorf("failed to generate new key: %w", err)
    }

    // Store new key with expiration
    expiresAt := time.Now().Add(90 * 24 * time.Hour) // 90 days
    err := km.keyStore.StoreKey(ctx, newKeyID, newKey, expiresAt)
    if err != nil {
        return fmt.Errorf("failed to store new key: %w", err)
    }

    // Update current key
    if km.currentKeyID != "" {
        // Mark old key as inactive but keep for decryption
        // It will be cleaned up after all data encrypted with it expires
    }

    km.currentKeyID = newKeyID
    km.keys[newKeyID] = newKey

    log.Infof("Key rotated successfully: new_key_id=%s", newKeyID)
    return nil
}
```

## Input Validation and Sanitization

### 1. Request Validation

```go
type PaymentValidator struct {
    validator *validator.Validate
}

func NewPaymentValidator() *PaymentValidator {
    v := validator.New()

    // Register custom validation rules
    v.RegisterValidation("order_id", validateOrderID)
    v.RegisterValidation("currency", validateCurrency)
    v.RegisterValidation("phone", validatePhoneNumber)

    return &PaymentValidator{validator: v}
}

func validateOrderID(fl validator.FieldLevel) bool {
    orderID := fl.Field().String()

    // Order ID should be alphanumeric with optional dashes and underscores
    // Length between 3 and 50 characters
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_]{3,50}$`, orderID)
    return matched
}

func validateCurrency(fl validator.FieldLevel) bool {
    currency := fl.Field().String()

    // Supported currencies
    supportedCurrencies := map[string]bool{
        "IDR": true,
        "USD": true,
        "EUR": true,
        "SGD": true,
    }

    return supportedCurrencies[currency]
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
    phone := fl.Field().String()

    // Basic phone number validation for Indonesian format
    // Supports: +62XXX, 08XXX, (62) XXX
    matched, _ := regexp.MatchString(`^(\+62|62|\(62\))?[-]?0?8[1-9][0-9]{6,11}$`, phone)
    return matched
}
```

### 2. SQL/NoSQL Injection Prevention

```go
// Use parameterized queries for MongoDB
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*PaymentRecord, error) {
    filter := bson.M{"order_id": orderID} // Safe: uses Go types, not string concatenation

    var payment PaymentRecord
    err := r.collection.FindOne(ctx, filter).Decode(&payment)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("payment not found")
        }
        return nil, fmt.Errorf("database error: %w", err)
    }

    return &payment, nil
}

// Input sanitization
func sanitizeInput(input string) string {
    // Remove potentially dangerous characters
    sanitized := strings.ReplaceAll(input, "<", "&lt;")
    sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
    sanitized = strings.ReplaceAll(sanitized, "&", "&amp;")
    sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
    sanitized = strings.ReplaceAll(sanitized, "'", "&#x27;")

    // Trim whitespace
    sanitized = strings.TrimSpace(sanitized)

    return sanitized
}
```

## Rate Limiting and DDoS Protection

### 1. Multi-Level Rate Limiting

```go
type RateLimitConfig struct {
    Global    rate.Limit
    PerUser   rate.Limit
    PerIP     rate.Limit
    PerOrder  rate.Limit
}

type RateLimiter struct {
    globalLimiter    *rate.Limiter
    userLimiters     map[string]*rate.Limiter
    ipLimiters       map[string]*rate.Limiter
    orderLimiters    map[string]*rate.Limiter
    mu               sync.RWMutex
    config           RateLimitConfig
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
    return &RateLimiter{
        globalLimiter: rate.NewLimiter(config.Global, 100),
        userLimiters:  make(map[string]*rate.Limiter),
        ipLimiters:    make(map[string]*rate.Limiter),
        orderLimiters: make(map[string]*rate.Limiter),
        config:        config,
    }
}

func (rl *RateLimiter) CheckLimit(ctx context.Context, userID, ip, orderID string) error {
    // Check global limit
    if !rl.globalLimiter.Allow() {
        return errors.New("global rate limit exceeded")
    }

    // Check per-user limit
    if userID != "" {
        userLimiter := rl.getUserLimiter(userID)
        if !userLimiter.Allow() {
            return errors.New("user rate limit exceeded")
        }
    }

    // Check per-IP limit
    if ip != "" {
        ipLimiter := rl.getIPLimiter(ip)
        if !ipLimiter.Allow() {
            return errors.New("IP rate limit exceeded")
        }
    }

    // Check per-order limit (for repeated status checks)
    if orderID != "" {
        orderLimiter := rl.getOrderLimiter(orderID)
        if !orderLimiter.Allow() {
            return errors.New("order rate limit exceeded")
        }
    }

    return nil
}

func (rl *RateLimiter) getUserLimiter(userID string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if limiter, exists := rl.userLimiters[userID]; exists {
        return limiter
    }

    limiter := rate.NewLimiter(rl.config.PerUser, 10)
    rl.userLimiters[userID] = limiter
    return limiter
}
```

### 2. Request Size Limiting

```go
func RequestSizeLimitMiddleware(maxSize int64) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Check content-length header
        if contentLength := c.Get("Content-Length"); contentLength != "" {
            if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
                if size > maxSize {
                    return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
                        "message": "Request too large",
                        "errors": []fiber.Map{{
                            "code":    "REQUEST_TOO_LARGE",
                            "message": fmt.Sprintf("Request size %d exceeds maximum allowed size %d", size, maxSize),
                        }},
                    })
                }
            }
        }

        // Also limit the body reading
        c.Context().SetBodyStream(c.Request().BodyStream(), int(maxSize))

        return c.Next()
    }
}
```

## Logging and Monitoring

### 1. Security Event Logging

```go
type SecurityEvent struct {
    EventID     string                 `json:"event_id"`
    EventType   string                 `json:"event_type"`
    Timestamp   time.Time              `json:"timestamp"`
    UserID      string                 `json:"user_id,omitempty"`
    IPAddress   string                 `json:"ip_address"`
    UserAgent   string                 `json:"user_agent"`
    RequestID   string                 `json:"request_id"`
    Severity    string                 `json:"severity"`
    Description string                 `json:"description"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type SecurityLogger struct {
    logger *logrus.Logger
}

func NewSecurityLogger() *SecurityLogger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetLevel(logrus.InfoLevel)

    return &SecurityLogger{logger: logger}
}

func (sl *SecurityLogger) LogSecurityEvent(eventType, description string, severity string, metadata map[string]interface{}) {
    event := SecurityEvent{
        EventID:     generateEventID(),
        EventType:   eventType,
        Timestamp:   time.Now().UTC(),
        IPAddress:   metadata["remote_addr"].(string),
        UserAgent:   metadata["user_agent"].(string),
        RequestID:   metadata["request_id"].(string),
        Severity:    severity,
        Description: description,
        Metadata:    metadata,
    }

    fields := logrus.Fields{
        "event": event,
    }

    switch severity {
    case "critical":
        sl.logger.WithFields(fields).Error("Security event")
    case "high":
        sl.logger.WithFields(fields).Warn("Security event")
    case "medium":
        sl.logger.WithFields(fields).Info("Security event")
    default:
        sl.logger.WithFields(fields).Debug("Security event")
    }
}

// Example usage
func logSecurityIncident(eventType string, metadata map[string]interface{}) {
    logger := NewSecurityLogger()
    logger.LogSecurityEvent(
        eventType,
        "Security incident detected",
        "high",
        metadata,
    )
}
```

### 2. Monitoring and Alerting

```go
type SecurityMetrics struct {
    FailedAuths        prometheus.Counter
    RateLimitHits      prometheus.Counter
    WebhookFailures    prometheus.Counter
    SecurityEvents     prometheus.Counter
    ResponseTime       prometheus.Histogram
}

func NewSecurityMetrics() *SecurityMetrics {
    return &SecurityMetrics{
        FailedAuths: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "payment_service_failed_authentications_total",
            Help: "Total number of failed authentication attempts",
        }),
        RateLimitHits: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "payment_service_rate_limit_hits_total",
            Help: "Total number of rate limit violations",
        }),
        WebhookFailures: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "payment_service_webhook_failures_total",
            Help: "Total number of webhook processing failures",
        }),
        SecurityEvents: prometheus.NewCounterVec(prometheus.CounterOpts{
            Name: "payment_service_security_events_total",
            Help: "Total number of security events by type",
        }, []string{"event_type", "severity"}),
        ResponseTime: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name:    "payment_service_response_time_seconds",
            Help:    "Payment service response time in seconds",
            Buckets: prometheus.DefBuckets,
        }),
    }
}

// Alert on suspicious activity
func (sm *SecurityMetrics) CheckAnomalies() {
    // Check for unusual authentication failures
    if sm.FailedAuths.Get() > 100 {
        sendAlert("high_authentication_failure_rate", map[string]interface{}{
            "failed_auths": sm.FailedAuths.Get(),
            "threshold":    100,
        })
    }

    // Check for rate limit violations
    if sm.RateLimitHits.Get() > 50 {
        sendAlert("high_rate_limit_violations", map[string]interface{}{
            "rate_limit_hits": sm.RateLimitHits.Get(),
            "threshold":       50,
        })
    }
}
```

## Compliance Requirements

### 1. PCI DSS Compliance

```go
// PCI DSS requirements implementation

// 1. Never store raw card data
type PCICompliantPaymentRequest struct {
    TokenID     string `json:"token_id"`      // Midtrans token, not PAN
    OrderID     string `json:"order_id"`
    Amount      int64  `json:"amount"`
    ExpiryMonth string `json:"expiry_month"` // Only if absolutely necessary
    ExpiryYear  string `json:"expiry_year"`  // Only if absolutely necessary
}

// 2. Implement strong access controls
func (s *PaymentService) checkPCIAccess(userID string) error {
    // Check if user has PCI compliance training
    hasTraining, err := s.userRepo.HasPCITraining(userID)
    if err != nil {
        return fmt.Errorf("failed to check PCI training: %w", err)
    }

    if !hasTraining {
        return errors.New("user does not have PCI compliance training")
    }

    // Check role-based access
    userRole, err := s.userRepo.GetUserRole(userID)
    if err != nil {
        return fmt.Errorf("failed to get user role: %w", err)
    }

    allowedRoles := []string{"payment_admin", "payment_operator"}
    if !contains(allowedRoles, userRole) {
        return errors.New("user does not have required role for payment operations")
    }

    return nil
}

// 3. Encrypt all data in transit
func configureTLS() *tls.Config {
    return &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }
}
```

### 2. Data Privacy and GDPR

```go
// GDPR compliance implementation

type PersonalData struct {
    UserID      string    `json:"user_id"`
    Email       string    `json:"email"`
    Phone       string    `json:"phone"`
    FullName    string    `json:"full_name"`
    Address     string    `json:"address"`
    CreatedAt   time.Time `json:"created_at"`
    LastUpdated time.Time `json:"last_updated"`
}

// Right to be forgotten
func (s *PaymentService) DeleteUserData(ctx context.Context, userID string) error {
    // Anonymize payment records instead of deleting (for audit purposes)
    filter := bson.M{"user_id": userID}
    update := bson.M{
        "$set": bson.M{
            "user_id":          "ANONYMIZED_" + generateHash(userID),
            "customer_details": bson.M{},
            "email":           "",
            "phone":           "",
            "anonymized_at":   time.Now(),
        },
    }

    _, err := s.paymentCollection.UpdateMany(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to anonymize user data: %w", err)
    }

    log.Infof("User data anonymized: user_id=%s", userID)
    return nil
}

// Data export functionality
func (s *PaymentService) ExportUserData(ctx context.Context, userID string) (*UserDataExport, error) {
    // Collect all user-related data
    payments, err := s.paymentRepository.GetByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user payments: %w", err)
    }

    userProfile, err := s.userRepo.GetProfile(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user profile: %w", err)
    }

    export := &UserDataExport{
        UserID:        userID,
        ExportDate:    time.Now(),
        Payments:      payments,
        Profile:       userProfile,
        ExportFormat:  "JSON",
    }

    // Log the export for compliance
    log.Infof("User data exported: user_id=%s, export_id=%s", userID, export.ExportID)

    return export, nil
}

type UserDataExport struct {
    ExportID    string                `json:"export_id"`
    UserID      string                `json:"user_id"`
    ExportDate  time.Time             `json:"export_date"`
    Payments    []PaymentRecord       `json:"payments"`
    Profile     UserProfile           `json:"profile"`
    ExportFormat string               `json:"export_format"`
}
```

## Security Testing

### 1. Penetration Testing Checklist

```
- [ ] Authentication bypass attempts
- [ ] SQL/NoSQL injection testing
- [ ] XSS vulnerability scanning
- [ ] CSRF token validation
- [ ] Rate limiting effectiveness
- [ ] Webhook signature verification
- [ ] Session management security
- [ ] Input validation bypass
- [ ] File upload security
- [ ] HTTPS/TLS configuration
- [ ] Security headers presence
- [ ] Error information disclosure
- [ ] Authentication token handling
- [ ] API endpoint enumeration
- [ ] Business logic vulnerabilities
```

### 2. Security Testing Implementation

```go
// Integration tests for security
func TestWebhookSignatureVerification(t *testing.T) {
    // Test cases for signature verification
    testCases := []struct {
        name           string
        orderID        string
        statusCode     string
        grossAmount    string
        signature      string
        serverKey      string
        expectedResult bool
    }{
        {
            name:           "Valid signature",
            orderID:        "ORDER-123",
            statusCode:     "200",
            grossAmount:    "10000.00",
            signature:      "valid_signature_hash",
            serverKey:      "test_key",
            expectedResult: true,
        },
        {
            name:           "Invalid signature",
            orderID:        "ORDER-123",
            statusCode:     "200",
            grossAmount:    "10000.00",
            signature:      "invalid_signature_hash",
            serverKey:      "test_key",
            expectedResult: false,
        },
        {
            name:           "Modified amount",
            orderID:        "ORDER-123",
            statusCode:     "200",
            grossAmount:    "20000.00", // Modified
            signature:      "valid_signature_hash", // Original signature
            serverKey:      "test_key",
            expectedResult: false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := VerifyWebhookSignature(
                tc.orderID,
                tc.statusCode,
                tc.grossAmount,
                tc.signature,
                tc.serverKey,
            )

            if result != tc.expectedResult {
                t.Errorf("Expected %v, got %v", tc.expectedResult, result)
            }
        })
    }
}

// Test for rate limiting
func TestRateLimiting(t *testing.T) {
    limiter := NewRateLimiter(RateLimitConfig{
        Global:   rate.Limit(10),
        PerUser:  rate.Limit(5),
        PerIP:    rate.Limit(3),
        PerOrder: rate.Limit(2),
    })

    ctx := context.Background()
    userID := "test_user"
    ip := "192.168.1.1"
    orderID := "ORDER-123"

    // Test per-user rate limiting
    for i := 0; i < 6; i++ {
        err := limiter.CheckLimit(ctx, userID, ip, orderID)
        if i < 5 {
            if err != nil {
                t.Errorf("Request %d should pass: %v", i, err)
            }
        } else {
            if err == nil {
                t.Error("Request 6 should be rate limited")
            }
        }
    }
}
```

## Security Checklist

### Pre-Deployment Security Checklist

```
Authentication & Authorization
- [ ] JWT token validation implemented
- [ ] Role-based access control configured
- [ ] API key encryption implemented
- [ ] Session timeout configured
- [ ] Multi-factor authentication (if applicable)

Data Protection
- [ ] Sensitive data encrypted at rest
- [ ] TLS 1.2+ configured for all communications
- [ ] Input validation implemented
- [ ] SQL/NoSQL injection prevention
- [ ] Data anonymization for deleted records

API Security
- [ ] Rate limiting configured
- [ ] Request size limits implemented
- [ ] CORS properly configured
- [ ] Security headers added
- [ ] API versioning implemented

Webhook Security
- [ ] Signature verification implemented
- [ ] Replay attack prevention
- [ ] IP whitelisting (if applicable)
- [ ] Timeout and retry logic

Monitoring & Logging
- [ ] Security event logging implemented
- [ ] Alert thresholds configured
- [ ] Audit trail maintained
- [ ] Performance monitoring enabled

Compliance
- [ ] PCI DSS requirements met
- [ ] GDPR compliance implemented
- [ ] Data retention policy defined
- [ ] Privacy policy updated

Testing
- [ ] Security testing completed
- [ ] Penetration testing performed
- [ ] Vulnerability scan completed
- [ ] Load testing with security scenarios
```

This security guide provides comprehensive coverage of security considerations for integrating Midtrans payment gateway while maintaining compliance with industry standards and best practices.