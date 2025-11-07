# Midtrans Payment Gateway Technical Implementation Guide

## Overview

This guide provides detailed technical implementation instructions for integrating Midtrans payment gateway into the InstrLabs Go microservices architecture. It includes code examples, configuration templates, and step-by-step implementation procedures.

## Prerequisites

### Required Dependencies

```go
// go.mod
module github.com/instrlabs/payment-service

go 1.24

require (
    github.com/gofiber/fiber/v2 v2.52.0
    github.com/joho/godotenv v1.5.1
    go.mongodb.org/mongo-driver v1.17.0
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/go-playground/validator/v10 v10.19.0
)
```

### Environment Setup

Create `.env.example` file:

```bash
# Service Configuration
SERVICE_NAME=payment-service
PORT=3002
ENVIRONMENT=development

# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=payment_service
MONGO_TIMEOUT=10

# Midtrans Configuration
MIDTRANS_SERVER_KEY=SB-Mid-server-your-sandbox-key
MIDTRANS_CLIENT_KEY=SB-Mid-client-your-sandbox-key
MIDTRANS_ENVIRONMENT=sandbox

# Security
JWT_SECRET=your-super-secret-jwt-key
WEBHOOK_SECRET=your-webhook-signature-secret

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# Rate Limiting
RATE_LIMIT=100
RATE_WINDOW=60s

# Timeouts
READ_TIMEOUT=30
WRITE_TIMEOUT=30
IDLE_TIMEOUT=60

# External Services
AUTH_SERVICE_URL=http://localhost:3001
NOTIFICATION_SERVICE_URL=http://localhost:3004
```

## Project Structure

```
payment-service/
├── main.go                 # Entry point
├── go.mod                  # Dependencies
├── go.sum                  # Dependency lock
├── Dockerfile              # Multi-stage build
├── .dockerignore           # Build optimization
├── .env.example            # Environment template
├── internal/               # Private app code
│   ├── config.go           # Environment & configuration
│   ├── middleware.go       # HTTP middleware setup
│   ├── router.go           # Route definitions
│   ├── database.go         # Database connection
│   ├── errors.go           # Error definitions
│   ├── handlers/           # HTTP handlers
│   │   ├── payment.go      # Payment handlers
│   │   ├── webhook.go      # Webhook handlers
│   │   └── health.go       # Health check handlers
│   ├── models/             # Data structures
│   │   ├── payment.go      # Payment models
│   │   ├── customer.go     # Customer models
│   │   └── webhook.go      # Webhook models
│   ├── services/           # Business logic
│   │   ├── payment.go      # Payment service
│   │   ├── midtrans.go     # Midtrans integration
│   │   └── webhook.go      # Webhook service
│   └── repositories/       # Data access layer
│       ├── payment.go      # Payment repository
│       └── webhook.go      # Webhook repository
├── pkg/                    # Public reusable code
│   ├── midtrans/           # Midtrans client
│   │   ├── client.go       # API client
│   │   ├── types.go        # Type definitions
│   │   └── constants.go    # Constants
│   └── utils/              # Utility functions
│       ├── crypto.go       # Cryptographic utilities
│       └── validator.go    # Validation utilities
├── static/                 # Static assets
│   └── swagger.json        # API documentation
└── scripts/                # Build/deployment scripts
    └── build.sh            # Build script
```

## Configuration Implementation

### Configuration Structure (`internal/config.go`)

```go
package internal

import (
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    // Service
    ServiceName string `env:"SERVICE_NAME,required"`
    Port        string `env:"PORT,default=3002"`
    Environment string `env:"ENVIRONMENT,default=development"`

    // Database
    MongoURI   string `env:"MONGODB_URI,required"`
    MongoDB    string `env:"MONGODB_DB,required"`
    MongoTimeout int `env:"MONGO_TIMEOUT,default=10"`

    // Midtrans
    MidtransServerKey string `env:"MIDTRANS_SERVER_KEY,required"`
    MidtransClientKey string `env:"MIDTRANS_CLIENT_KEY,required"`
    MidtransEnvironment string `env:"MIDTRANS_ENVIRONMENT,default=sandbox"`

    // Security
    JWTSecret     string `env:"JWT_SECRET,required"`
    WebhookSecret string `env:"WEBHOOK_SECRET,required"`

    // CORS
    Origins string `env:"CORS_ORIGINS,default=http://localhost:3000"`

    // Rate limiting
    RateLimit  int           `env:"RATE_LIMIT,default=100"`
    RateWindow time.Duration `env:"RATE_WINDOW,default=60s"`

    // Timeouts
    ReadTimeout  int `env:"READ_TIMEOUT,default=30"`
    WriteTimeout int `env:"WRITE_TIMEOUT,default=30"`
    IdleTimeout  int `env:"IDLE_TIMEOUT,default=60"`

    // External Services
    AuthServiceURL         string `env:"AUTH_SERVICE_URL,required"`
    NotificationServiceURL string `env:"NOTIFICATION_SERVICE_URL,required"`
}

func LoadConfig() (*Config, error) {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
        return nil, err
    }

    cfg := &Config{}

    // Parse environment variables
    cfg.ServiceName = os.Getenv("SERVICE_NAME")
    if cfg.ServiceName == "" {
        cfg.ServiceName = "payment-service"
    }

    cfg.Port = getEnv("PORT", "3002")
    cfg.Environment = getEnv("ENVIRONMENT", "development")

    cfg.MongoURI = os.Getenv("MONGODB_URI")
    cfg.MongoDB = getEnv("MONGODB_DB", "payment_service")
    cfg.MongoTimeout = getEnvInt("MONGO_TIMEOUT", 10)

    cfg.MidtransServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
    cfg.MidtransClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
    cfg.MidtransEnvironment = getEnv("MIDTRANS_ENVIRONMENT", "sandbox")

    cfg.JWTSecret = os.Getenv("JWT_SECRET")
    cfg.WebhookSecret = os.Getenv("WEBHOOK_SECRET")

    cfg.Origins = getEnv("CORS_ORIGINS", "http://localhost:3000")

    cfg.RateLimit = getEnvInt("RATE_LIMIT", 100)
    rateWindowSeconds := getEnvInt("RATE_WINDOW_SECONDS", 60)
    cfg.RateWindow = time.Duration(rateWindowSeconds) * time.Second

    cfg.ReadTimeout = getEnvInt("READ_TIMEOUT", 30)
    cfg.WriteTimeout = getEnvInt("WRITE_TIMEOUT", 30)
    cfg.IdleTimeout = getEnvInt("IDLE_TIMEOUT", 60)

    cfg.AuthServiceURL = os.Getenv("AUTH_SERVICE_URL")
    cfg.NotificationServiceURL = os.Getenv("NOTIFICATION_SERVICE_URL")

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

## Database Implementation

### Database Connection (`internal/database.go`)

```go
package internal

import (
    "context"
    "log"
    "time"

    "github.com/gofiber/fiber/v2/log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
    Client   *mongo.Client
    Database *mongo.Database
}

func NewDatabase(cfg *Config) (*Database, error) {
    ctx, cancel := context.WithTimeout(context.Background(),
        time.Duration(cfg.MongoTimeout)*time.Second)
    defer cancel()

    clientOpts := options.Client().
        ApplyURI(cfg.MongoURI).
        SetMaxPoolSize(100).
        SetMinPoolSize(10).
        SetMaxConnIdleTime(30 * time.Second)

    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        return nil, err
    }

    // Ping to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }

    log.Info("Database connected successfully")

    return &Database{
        Client:   client,
        Database: client.Database(cfg.MongoDB),
    }, nil
}

func (d *Database) Close() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return d.Client.Disconnect(ctx)
}

func (d *Database) CreateIndexes() error {
    ctx := context.Background()

    // Payments collection indexes
    paymentsCollection := d.Database.Collection("payments")

    // Unique index on order_id
    _, err := paymentsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys:    bson.D{{Key: "order_id", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
    if err != nil {
        return err
    }

    // Index on user_id for querying user payments
    _, err = paymentsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{Key: "user_id", Value: 1}},
    })
    if err != nil {
        return err
    }

    // Index on status for processing
    _, err = paymentsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{Key: "status", Value: 1}},
    })
    if err != nil {
        return err
    }

    // TTL index for expired transactions
    _, err = paymentsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys: bson.D{{Key: "expires_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(0),
    })

    return err
}
```

## Midtrans Client Implementation

### Midtrans API Client (`pkg/midtrans/client.go`)

```go
package midtrans

import (
    "bytes"
    "context"
    "crypto/sha512"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    serverKey   string
    clientKey   string
    environment string
    httpClient  *http.Client
}

func NewClient(serverKey, clientKey, environment string) *Client {
    return &Client{
        serverKey:   serverKey,
        clientKey:   clientKey,
        environment: environment,
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

func (c *Client) getBaseURL() string {
    if c.environment == "production" {
        return "https://api.midtrans.com/v2"
    }
    return "https://api.sandbox.midtrans.com/v2"
}

func (c *Client) createAuthHeader() string {
    auth := c.serverKey + ":"
    return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) Charge(ctx context.Context, request *ChargeRequest) (*ChargeResponse, error) {
    url := c.getBaseURL() + "/charge"

    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", c.createAuthHeader())

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    var chargeResp ChargeResponse
    if err := json.Unmarshal(body, &chargeResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    // Handle HTTP errors
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return &chargeResp, fmt.Errorf("midtrans API error: status=%d, body=%s", resp.StatusCode, string(body))
    }

    return &chargeResp, nil
}

func (c *Client) GetStatus(ctx context.Context, orderID string) (*StatusResponse, error) {
    url := c.getBaseURL() + "/" + orderID + "/status"

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", c.createAuthHeader())

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    var statusResp StatusResponse
    if err := json.Unmarshal(body, &statusResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return &statusResp, fmt.Errorf("midtrans API error: status=%d, body=%s", resp.StatusCode, string(body))
    }

    return &statusResp, nil
}

func (c *Client) Cancel(ctx context.Context, orderID string) (*CancelResponse, error) {
    url := c.getBaseURL() + "/" + orderID + "/cancel"

    req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("Authorization", c.createAuthHeader())

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    var cancelResp CancelResponse
    if err := json.Unmarshal(body, &cancelResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return &cancelResp, fmt.Errorf("midtrans API error: status=%d, body=%s", resp.StatusCode, string(body))
    }

    return &cancelResp, nil
}

func (c *Client) Refund(ctx context.Context, orderID string, request *RefundRequest) (*RefundResponse, error) {
    url := c.getBaseURL() + "/" + orderID + "/refund"

    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", c.createAuthHeader())

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    var refundResp RefundResponse
    if err := json.Unmarshal(body, &refundResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return &refundResp, fmt.Errorf("midtrans API error: status=%d, body=%s", resp.StatusCode, string(body))
    }

    return &refundResp, nil
}

func VerifyWebhookSignature(orderID, statusCode, grossAmount, signatureKey, serverKey string) bool {
    input := orderID + statusCode + grossAmount + serverKey
    hash := sha512.Sum512([]byte(input))
    calculatedSignature := fmt.Sprintf("%x", hash)
    return calculatedSignature == signatureKey
}
```

### Type Definitions (`pkg/midtrans/types.go`)

```go
package midtrans

import "time"

// Core Request/Response Types
type ChargeRequest struct {
    PaymentType       string                 `json:"payment_type"`
    TransactionDetails TransactionDetails    `json:"transaction_details"`
    CustomerDetails   CustomerDetails       `json:"customer_details"`
    ItemDetails       []ItemDetails         `json:"item_details,omitempty"`
    CustomField       map[string]interface{} `json:"custom_field,omitempty"`

    // Payment method specific fields
    CreditCard      *CreditCardDetails      `json:"credit_card,omitempty"`
    BankTransfer    *BankTransferDetails    `json:"bank_transfer,omitempty"`
    EWallet         *EWalletDetails         `json:"ewallet,omitempty"`
    CStore          *CStoreDetails          `json:"cstore,omitempty"`
    CardlessCredit  *CardlessCreditDetails  `json:"cardless_credit,omitempty"`
}

type ChargeResponse struct {
    StatusCode        string                 `json:"status_code"`
    StatusMessage     string                 `json:"status_message"`
    TransactionID     string                 `json:"transaction_id"`
    OrderID           string                 `json:"order_id"`
    GrossAmount       string                 `json:"gross_amount"`
    PaymentType       string                 `json:"payment_type"`
    TransactionTime   string                 `json:"transaction_time"`
    TransactionStatus string                 `json:"transaction_status"`
    FraudStatus       string                 `json:"fraud_status,omitempty"`
    ApprovalCode      string                 `json:"approval_code,omitempty"`

    // Payment method specific response fields
    RedirectURL       string                 `json:"redirect_url,omitempty"`
    VaNumbers         []VANumber             `json:"va_numbers,omitempty"`
    QRCode            string                 `json:"qr_code,omitempty"`
    DeepLink          string                 `json:"deeplink,omitempty"`
    PaymentCode       string                 `json:"payment_code,omitempty"`
    Store             string                 `json:"store,omitempty"`
    MerchantID        string                 `json:"merchant_id,omitempty"`

    // Additional data
    Actions           []Action               `json:"actions,omitempty"`
    CustomFields      map[string]interface{} `json:"custom_fields,omitempty"`
}

type TransactionDetails struct {
    OrderID     string `json:"order_id"`
    GrossAmount int64  `json:"gross_amount"`
}

type CustomerDetails struct {
    FirstName      string        `json:"first_name"`
    LastName       string        `json:"last_name"`
    Email          string        `json:"email"`
    Phone          string        `json:"phone"`
    BillingAddress *Address      `json:"billing_address,omitempty"`
    ShippingAddress *Address     `json:"shipping_address,omitempty"`
}

type Address struct {
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
    Address       string `json:"address"`
    City          string `json:"city"`
    PostalCode    string `json:"postal_code"`
    Phone         string `json:"phone"`
    CountryCode   string `json:"country_code"`
}

type ItemDetails struct {
    ID       string  `json:"id"`
    Price    int64   `json:"price"`
    Quantity int     `json:"quantity"`
    Name     string  `json:"name"`
    Category string  `json:"category,omitempty"`
    Brand    string  `json:"brand,omitempty"`
    Merchant string  `json:"merchant,omitempty"`
}

// Credit Card Details
type CreditCardDetails struct {
    CardNumber  string `json:"card_number"`
    CardCVV     string `json:"card_cvv"`
    CardExpire  string `json:"card_exp_month_year"`
    TokenID     string `json:"token_id,omitempty"`
    SaveToken   bool   `json:"save_token,omitempty"`
    ThreeDSecure bool  `json:"secure"`
}

// Bank Transfer Details
type BankTransferDetails struct {
    Bank            string `json:"bank"`
    VaNumber        string `json:"va_number,omitempty"`
    FreeText        map[string]string `json:"free_text,omitempty"`
    BillKey         string `json:"bill_key,omitempty"`
    BillerCode      string `json:"biller_code,omitempty"`
}

type VANumber struct {
    Bank     string `json:"bank"`
    VANumber string `json:"va_number"`
}

// E-Wallet Details
type EWalletDetails struct {
    Gopay *EWalletGopay `json:"gopay,omitempty"`
}

type EWalletGopay struct {
    EnableCallback   bool   `json:"enable_callback"`
    CallbackURL      string `json:"callback_url,omitempty"`
    Acquirer         string `json:"acquirer,omitempty"`
    PaymentType      string `json:"payment_type"`
    TransactionUUID  string `json:"transaction_uuid,omitempty"`
}

// Convenience Store Details
type CStoreDetails struct {
    Store         string `json:"store"`
    Message       string `json:"message,omitempty"`
    Name          string `json:"name,omitempty"`
    Phone         string `json:"phone,omitempty"`
    Email         string `json:"email,omitempty"`
}

// Cardless Credit Details
type CardlessCreditDetails struct {
    Provider    string `json:"provider"`
    CancelURL   string `json:"cancel_url,omitempty"`
    ReturnURL   string `json:"return_url,omitempty"`
    RedirectURL string `json:"redirect_url,omitempty"`
}

// Action Details
type Action struct {
    Name       string `json:"name"`
    Method     string `json:"method"`
    URL        string `json:"url"`
    Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// Status Response
type StatusResponse struct {
    StatusCode        string                 `json:"status_code"`
    StatusMessage     string                 `json:"status_message"`
    TransactionID     string                 `json:"transaction_id"`
    OrderID           string                 `json:"order_id"`
    GrossAmount       string                 `json:"gross_amount"`
    PaymentType       string                 `json:"payment_type"`
    TransactionTime   string                 `json:"transaction_time"`
    TransactionStatus string                 `json:"transaction_status"`
    FraudStatus       string                 `json:"fraud_status,omitempty"`
    SignatureKey      string                 `json:"signature_key"`
    ApprovalCode      string                 `json:"approval_code,omitempty"`
    PaymentCode       string                 `json:"payment_code,omitempty"`
    Store             string                 `json:"store,omitempty"`
    MerchantID        string                 `json:"merchant_id,omitempty"`
    CustomFields      map[string]interface{} `json:"custom_fields,omitempty"`
}

// Cancel Response
type CancelResponse struct {
    StatusCode    string `json:"status_code"`
    StatusMessage string `json:"status_message"`
    TransactionID string `json:"transaction_id"`
    OrderID       string `json:"order_id"`
    GrossAmount   string `json:"gross_amount"`
    PaymentType   string `json:"payment_type"`
}

// Refund Request
type RefundRequest struct {
    RefundKeys    []string `json:"refund_keys,omitempty"`
    Amount        int64    `json:"amount,omitempty"`
    Reason        string   `json:"reason,omitempty"`
}

// Refund Response
type RefundResponse struct {
    StatusCode    string `json:"status_code"`
    StatusMessage string `json:"status_message"`
    TransactionID string `json:"transaction_id"`
    OrderID       string `json:"order_id"`
    GrossAmount   string `json:"gross_amount"`
    RefundAmount  string `json:"refund_amount"`
    RefundKeys    []string `json:"refund_keys"`
}

// Webhook Notification
type WebhookNotification struct {
    StatusCode        string                 `json:"status_code"`
    StatusMessage     string                 `json:"status_message"`
    TransactionID     string                 `json:"transaction_id"`
    OrderID           string                 `json:"order_id"`
    GrossAmount       string                 `json:"gross_amount"`
    PaymentType       string                 `json:"payment_type"`
    TransactionTime   string                 `json:"transaction_time"`
    TransactionStatus string                 `json:"transaction_status"`
    FraudStatus       string                 `json:"fraud_status,omitempty"`
    SignatureKey      string                 `json:"signature_key"`
    ApprovalCode      string                 `json:"approval_code,omitempty"`
    PaymentCode       string                 `json:"payment_code,omitempty"`
    Store             string                 `json:"store,omitempty"`
    MerchantID        string                 `json:"merchant_id,omitempty"`
    CustomFields      map[string]interface{} `json:"custom_fields,omitempty"`
}
```

## Service Layer Implementation

### Payment Service (`internal/services/payment.go`)

```go
package services

import (
    "context"
    "fmt"
    "time"

    "github.com/instrlabs/payment-service/internal/models"
    "github.com/instrlabs/payment-service/pkg/midtrans"
    "github.com/gofiber/fiber/v2/log"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
    midtransClient *midtrans.Client
    paymentRepo    PaymentRepository
    webhookService *WebhookService
}

func NewPaymentService(midtransClient *midtrans.Client, paymentRepo PaymentRepository, webhookService *WebhookService) *PaymentService {
    return &PaymentService{
        midtransClient: midtransClient,
        paymentRepo:    paymentRepo,
        webhookService: webhookService,
    }
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
    log.Infof("Payment: Processing payment for order_id=%s, amount=%d, payment_type=%s",
        req.OrderID, req.GrossAmount, req.PaymentType)

    // Check if payment already exists
    existingPayment, err := s.paymentRepo.GetByOrderID(ctx, req.OrderID)
    if err == nil && existingPayment != nil {
        log.Warnf("Payment: Payment already exists for order_id=%s", req.OrderID)
        return nil, fmt.Errorf("payment already exists for order ID: %s", req.OrderID)
    }

    // Create payment record
    payment := &models.Payment{
        OrderID:         req.OrderID,
        UserID:          req.UserID,
        PaymentType:     req.PaymentType,
        GrossAmount:     req.GrossAmount,
        Status:          "pending",
        Currency:        req.Currency,
        CustomerDetails: req.CustomerDetails,
        PaymentDetails:  req.PaymentDetails,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
        ExpiresAt:       req.ExpiresAt,
        Metadata:        req.Metadata,
    }

    // Save initial payment record
    err = s.paymentRepo.Create(ctx, payment)
    if err != nil {
        log.Errorf("Payment: Failed to create payment record: %v", err)
        return nil, fmt.Errorf("failed to create payment record: %w", err)
    }

    // Build Midtrans charge request
    chargeReq := s.buildChargeRequest(payment)

    // Call Midtrans API
    chargeResp, err := s.midtransClient.Charge(ctx, chargeReq)
    if err != nil {
        log.Errorf("Payment: Midtrans charge failed: %v", err)
        // Update payment status to failed
        payment.Status = "failed"
        payment.UpdatedAt = time.Now()
        s.paymentRepo.Update(ctx, payment)
        return nil, fmt.Errorf("payment processing failed: %w", err)
    }

    // Update payment with response
    s.updatePaymentFromResponse(payment, chargeResp)
    err = s.paymentRepo.Update(ctx, payment)
    if err != nil {
        log.Errorf("Payment: Failed to update payment: %v", err)
        return nil, fmt.Errorf("failed to update payment: %w", err)
    }

    // Build response
    response := &ProcessPaymentResponse{
        OrderID:           payment.OrderID,
        TransactionID:     payment.TransactionID,
        PaymentType:       payment.PaymentType,
        Status:            payment.Status,
        GrossAmount:       payment.GrossAmount,
        Currency:          payment.Currency,
        TransactionTime:   payment.UpdatedAt,
        RedirectURL:       chargeResp.RedirectURL,
        PaymentCode:       chargeResp.PaymentCode,
        Store:             chargeResp.Store,
        VANumbers:         chargeResp.VaNumbers,
        QRCode:            chargeResp.QRCode,
        DeepLink:          chargeResp.DeepLink,
        Actions:           chargeResp.Actions,
    }

    log.Infof("Payment: Payment processed successfully: order_id=%s, transaction_id=%s, status=%s",
        payment.OrderID, payment.TransactionID, payment.Status)

    return response, nil
}

func (s *PaymentService) GetPaymentStatus(ctx context.Context, orderID string) (*PaymentStatusResponse, error) {
    log.Infof("Payment: Getting payment status for order_id=%s", orderID)

    // Get payment from database
    payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
    if err != nil {
        log.Errorf("Payment: Failed to get payment: %v", err)
        return nil, fmt.Errorf("payment not found: %w", err)
    }

    // If payment is still pending, check status from Midtrans
    if payment.Status == "pending" {
        statusResp, err := s.midtransClient.GetStatus(ctx, orderID)
        if err != nil {
            log.Errorf("Payment: Failed to get status from Midtrans: %v", err)
            return nil, fmt.Errorf("failed to get payment status: %w", err)
        }

        // Update payment if status changed
        if payment.TransactionStatus != statusResp.TransactionStatus {
            s.updatePaymentFromStatusResponse(payment, statusResp)
            err = s.paymentRepo.Update(ctx, payment)
            if err != nil {
                log.Errorf("Payment: Failed to update payment status: %v", err)
            }
        }

        return &PaymentStatusResponse{
            OrderID:           payment.OrderID,
            TransactionID:     payment.TransactionID,
            PaymentType:       payment.PaymentType,
            Status:            statusResp.TransactionStatus,
            GrossAmount:       payment.GrossAmount,
            Currency:          payment.Currency,
            TransactionTime:   payment.UpdatedAt,
            FraudStatus:       statusResp.FraudStatus,
            PaymentCode:       statusResp.PaymentCode,
            Store:             statusResp.Store,
        }, nil
    }

    return &PaymentStatusResponse{
        OrderID:           payment.OrderID,
        TransactionID:     payment.TransactionID,
        PaymentType:       payment.PaymentType,
        Status:            payment.TransactionStatus,
        GrossAmount:       payment.GrossAmount,
        Currency:          payment.Currency,
        TransactionTime:   payment.UpdatedAt,
        FraudStatus:       payment.FraudStatus,
        PaymentCode:       payment.PaymentCode,
        Store:             payment.Store,
    }, nil
}

func (s *PaymentService) CancelPayment(ctx context.Context, orderID string) (*CancelPaymentResponse, error) {
    log.Infof("Payment: Canceling payment for order_id=%s", orderID)

    // Get payment from database
    payment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
    if err != nil {
        log.Errorf("Payment: Failed to get payment: %v", err)
        return nil, fmt.Errorf("payment not found: %w", err)
    }

    // Check if payment can be canceled
    if payment.Status != "pending" {
        return nil, fmt.Errorf("payment cannot be canceled, current status: %s", payment.Status)
    }

    // Call Midtrans cancel API
    cancelResp, err := s.midtransClient.Cancel(ctx, orderID)
    if err != nil {
        log.Errorf("Payment: Failed to cancel payment: %v", err)
        return nil, fmt.Errorf("payment cancellation failed: %w", err)
    }

    // Update payment status
    payment.Status = "cancel"
    payment.TransactionStatus = "cancel"
    payment.UpdatedAt = time.Now()
    err = s.paymentRepo.Update(ctx, payment)
    if err != nil {
        log.Errorf("Payment: Failed to update payment: %v", err)
    }

    response := &CancelPaymentResponse{
        OrderID:       payment.OrderID,
        TransactionID: payment.TransactionID,
        Status:        payment.Status,
        Message:       "Payment canceled successfully",
    }

    log.Infof("Payment: Payment canceled successfully: order_id=%s", orderID)
    return response, nil
}

func (s *PaymentService) ProcessWebhook(ctx context.Context, notification *midtrans.WebhookNotification) error {
    log.Infof("Webhook: Processing webhook for order_id=%s, status=%s",
        notification.OrderID, notification.TransactionStatus)

    // Get payment from database
    payment, err := s.paymentRepo.GetByOrderID(ctx, notification.OrderID)
    if err != nil {
        log.Errorf("Webhook: Payment not found: %v", err)
        return fmt.Errorf("payment not found: %w", err)
    }

    // Update payment with webhook data
    s.updatePaymentFromWebhook(payment, notification)
    payment.UpdatedAt = time.Now()

    err = s.paymentRepo.Update(ctx, payment)
    if err != nil {
        log.Errorf("Webhook: Failed to update payment: %v", err)
        return fmt.Errorf("failed to update payment: %w", err)
    }

    // Trigger notification if needed
    if notification.TransactionStatus == "settlement" ||
       notification.TransactionStatus == "capture" {
        // Send payment success notification
        err = s.webhookService.SendPaymentNotification(ctx, payment, "success")
        if err != nil {
            log.Errorf("Webhook: Failed to send notification: %v", err)
        }
    }

    log.Infof("Webhook: Webhook processed successfully: order_id=%s, new_status=%s",
        notification.OrderID, notification.TransactionStatus)
    return nil
}

// Helper methods
func (s *PaymentService) buildChargeRequest(payment *models.Payment) *midtrans.ChargeRequest {
    return &midtrans.ChargeRequest{
        PaymentType:       payment.PaymentType,
        TransactionDetails: midtrans.TransactionDetails{
            OrderID:     payment.OrderID,
            GrossAmount: payment.GrossAmount,
        },
        CustomerDetails: payment.CustomerDetails,
        PaymentDetails:  payment.PaymentDetails,
    }
}

func (s *PaymentService) updatePaymentFromResponse(payment *models.Payment, resp *midtrans.ChargeResponse) {
    payment.TransactionID = resp.TransactionID
    payment.TransactionStatus = resp.TransactionStatus
    payment.FraudStatus = resp.FraudStatus
    payment.RedirectURL = resp.RedirectURL
    payment.PaymentCode = resp.PaymentCode
    payment.Store = resp.Store
    payment.VANumbers = resp.VANumbers
    payment.QRCode = resp.QRCode
    payment.DeepLink = resp.DeepLink
    payment.Actions = resp.Actions
    payment.ApprovalCode = resp.ApprovalCode
    payment.UpdatedAt = time.Now()
}

func (s *PaymentService) updatePaymentFromStatusResponse(payment *models.Payment, resp *midtrans.StatusResponse) {
    payment.TransactionStatus = resp.TransactionStatus
    payment.FraudStatus = resp.FraudStatus
    payment.PaymentCode = resp.PaymentCode
    payment.Store = resp.Store
    payment.ApprovalCode = resp.ApprovalCode
    payment.UpdatedAt = time.Now()
}

func (s *PaymentService) updatePaymentFromWebhook(payment *models.Payment, notification *midtrans.WebhookNotification) {
    payment.TransactionStatus = notification.TransactionStatus
    payment.FraudStatus = notification.FraudStatus
    payment.PaymentCode = notification.PaymentCode
    payment.Store = notification.Store
    payment.ApprovalCode = notification.ApprovalCode
    payment.SignatureKey = notification.SignatureKey
}

// Request/Response Types
type ProcessPaymentRequest struct {
    OrderID         string                     `json:"order_id" validate:"required"`
    UserID          string                     `json:"user_id" validate:"required"`
    PaymentType     string                     `json:"payment_type" validate:"required"`
    GrossAmount     int64                      `json:"gross_amount" validate:"required,min=1000"`
    Currency        string                     `json:"currency" validate:"required"`
    CustomerDetails midtrans.CustomerDetails   `json:"customer_details" validate:"required"`
    PaymentDetails  interface{}                `json:"payment_details" validate:"required"`
    ExpiresAt       *time.Time                 `json:"expires_at,omitempty"`
    Metadata        map[string]interface{}     `json:"metadata,omitempty"`
}

type ProcessPaymentResponse struct {
    OrderID         string                   `json:"order_id"`
    TransactionID   string                   `json:"transaction_id"`
    PaymentType     string                   `json:"payment_type"`
    Status          string                   `json:"status"`
    GrossAmount     int64                    `json:"gross_amount"`
    Currency        string                   `json:"currency"`
    TransactionTime time.Time                `json:"transaction_time"`
    RedirectURL     string                   `json:"redirect_url,omitempty"`
    PaymentCode     string                   `json:"payment_code,omitempty"`
    Store           string                   `json:"store,omitempty"`
    VANumbers       []midtrans.VANumber      `json:"va_numbers,omitempty"`
    QRCode          string                   `json:"qr_code,omitempty"`
    DeepLink        string                   `json:"deeplink,omitempty"`
    Actions         []midtrans.Action        `json:"actions,omitempty"`
}

type PaymentStatusResponse struct {
    OrderID         string    `json:"order_id"`
    TransactionID   string    `json:"transaction_id"`
    PaymentType     string    `json:"payment_type"`
    Status          string    `json:"status"`
    GrossAmount     int64     `json:"gross_amount"`
    Currency        string    `json:"currency"`
    TransactionTime time.Time `json:"transaction_time"`
    FraudStatus     string    `json:"fraud_status,omitempty"`
    PaymentCode     string    `json:"payment_code,omitempty"`
    Store           string    `json:"store,omitempty"`
}

type CancelPaymentResponse struct {
    OrderID       string `json:"order_id"`
    TransactionID string `json:"transaction_id"`
    Status        string `json:"status"`
    Message       string `json:"message"`
}

// Repository interface
type PaymentRepository interface {
    Create(ctx context.Context, payment *models.Payment) error
    GetByOrderID(ctx context.Context, orderID string) (*models.Payment, error)
    GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Payment, error)
    Update(ctx context.Context, payment *models.Payment) error
    Delete(ctx context.Context, orderID string) error
}
```

This technical implementation guide provides the foundational code structure for integrating Midtrans payment gateway. The remaining files (handlers, repository layer, models, etc.) follow the same patterns established in the InstrLabs architecture.

The next steps would be to:
1. Complete the handler implementations
2. Set up the repository layer with MongoDB
3. Create the webhook service
4. Implement middleware and routing
5. Add comprehensive error handling and logging
6. Set up monitoring and health checks

Would you like me to continue with any specific component or move on to the next document in the plan?