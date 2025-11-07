# Midtrans Payment Gateway Integration Architecture

## Overview

This document outlines the architectural approach for integrating Midtrans payment gateway into the InstrLabs microservices ecosystem. The integration follows Go microservices patterns using Fiber framework with a dedicated payment-service to handle all payment-related operations.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client App    │    │  Gateway Service│    │  Payment Service│
│                 │    │                 │    │                 │
│   - Frontend    │◄──►│   - Proxy       │◄──►│   - Core Logic  │
│   - Mobile      │    │   - Auth        │    │   - Midtrans API│
│   - Web         │    │   - Routing     │    │   - Validation  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Other Services │    │   Midtrans API  │
                       │                 │    │                 │
                       │ - Auth Service  │    │ - Charge API    │
                       │ - Image Service │    │ - Status API    │
                       │ - Notification  │    │ - Webhook       │
                       └─────────────────┘    └─────────────────┘
```

## Service Design

### Payment Service Responsibilities

1. **Payment Processing**
   - Handle all payment method integrations
   - Manage transaction lifecycle
   - Process refunds and cancellations

2. **API Management**
   - Midtrans API integration
   - Request/response transformation
   - Error handling and retry logic

3. **Data Management**
   - Transaction record keeping
   - Payment method tracking
   - Audit logging

4. **Webhook Handling**
   - Receive payment status updates
   - Validate webhook signatures
   - Trigger appropriate actions

### Database Schema

#### Payments Collection
```go
type Payment struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    OrderID         string            `bson:"order_id" json:"order_id"`
    TransactionID   string            `bson:"transaction_id" json:"transaction_id"`
    UserID          string            `bson:"user_id" json:"user_id"`
    PaymentType     string            `bson:"payment_type" json:"payment_type"`
    GrossAmount     int64             `bson:"gross_amount" json:"gross_amount"`
    Status          string            `bson:"status" json:"status"`
    Currency        string            `bson:"currency" json:"currency"`
    PaymentDetails  interface{}       `bson:"payment_details" json:"payment_details"`
    CustomerDetails CustomerDetails   `bson:"customer_details" json:"customer_details"`
    CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
    ExpiresAt       *time.Time        `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
    Metadata        map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type CustomerDetails struct {
    FirstName string `bson:"first_name" json:"first_name"`
    LastName  string `bson:"last_name" json:"last_name"`
    Email     string `bson:"email" json:"email"`
    Phone     string `bson:"phone" json:"phone"`
    BillingAddress *Address `bson:"billing_address,omitempty" json:"billing_address,omitempty"`
    ShippingAddress *Address `bson:"shipping_address,omitempty" json:"shipping_address,omitempty"`
}

type Address struct {
    Address     string `bson:"address" json:"address"`
    City        string `bson:"city" json:"city"`
    PostalCode  string `bson:"postal_code" json:"postal_code"`
    CountryCode string `bson:"country_code" json:"country_code"`
}
```

#### Payment Methods Collection
```go
type PaymentMethod struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Type        string            `bson:"type" json:"type"` // card, bank_transfer, ewallet, etc.
    Provider    string            `bson:"provider" json:"provider"` // gopay, kredivo, etc.
    Enabled     bool              `bson:"enabled" json:"enabled"`
    Config      interface{}       `bson:"config" json:"config"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}
```

## Inter-Service Communication

### API Gateway Integration

The payment service integrates with the gateway service through path-based routing:

```go
// Gateway service proxy configuration
app.All("/payments/*", proxy.Forward("http://payment-service:3002", proxy.Config{
    Timeout: 60 * time.Second,
    ModifyRequest: func(c *fiber.Ctx) error {
        c.Request().Header.Set("X-Forwarded-Host", c.Hostname())
        c.Request().Header.Set("X-User-ID", c.Locals("user_id").(string))
        return nil
    },
}))
```

### Service Dependencies

1. **Auth Service**
   - User authentication and authorization
   - User profile retrieval
   - Session management

2. **Notification Service**
   - Payment status notifications
   - Email/SMS confirmations
   - Real-time updates via SSE

## Payment Flow Architecture

### 1. Payment Initiation Flow

```
Client → Gateway → Auth Service → Payment Service → Midtrans API
   │        │          │                │               │
   │        │          │                │               │
   └────────┼──────────┼────────────────┼───────────────┘
            │          │                │
            ▼          ▼                ▼
        Validate   Authenticate    Process Payment
        Request     User           Request
```

### 2. Payment Completion Flow

```
Midtrans → Webhook → Payment Service → Notification Service → Client
    API        │           │                  │               │
               │           │                  │               │
               └───────────┼──────────────────┼───────────────┘
                           │                  │
                           ▼                  ▼
                     Update Status       Send Notification
```

## Security Architecture

### Authentication Flow

1. **Client Authentication**
   - JWT token validation via gateway service
   - User context propagation to payment service

2. **Midtrans Authentication**
   - Server Key stored in environment variables
   - Basic authentication for API calls
   - Request signature verification for webhooks

### Data Security

1. **Sensitive Data Handling**
   - Card tokenization via Midtrans.js
   - No raw card data stored in database
   - Encrypted configuration storage

2. **API Security**
   - Rate limiting per user/IP
   - Request size limits
   - Input validation and sanitization

## Configuration Management

### Environment Variables

```bash
# Service Configuration
SERVICE_NAME=payment-service
PORT=3002
ENVIRONMENT=development

# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=payment_service

# Midtrans Configuration
MIDTRANS_SERVER_KEY=your-server-key
MIDTRANS_CLIENT_KEY=your-client-key
MIDTRANS_ENVIRONMENT=sandbox # sandbox or production

# Security
JWT_SECRET=your-jwt-secret
WEBHOOK_SECRET=your-webhook-secret

# External Services
AUTH_SERVICE_URL=http://auth-service:3001
NOTIFICATION_SERVICE_URL=http://notification-service:3004
```

## Error Handling Strategy

### Error Categories

1. **Client Errors (4xx)**
   - Invalid request parameters
   - Authentication failures
   - Insufficient permissions

2. **Server Errors (5xx)**
   - Midtrans API failures
   - Database connection issues
   - Internal service errors

3. **Business Logic Errors**
   - Payment method not available
   - Transaction timeout
   - Invalid transaction status

### Error Response Format

```go
type ErrorResponse struct {
    Message string      `json:"message"`
    Errors  []ErrorItem `json:"errors,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type ErrorItem struct {
    Field   string `json:"field,omitempty"`
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

## Monitoring and Observability

### Metrics Collection

1. **Transaction Metrics**
   - Payment success/failure rates
   - Payment method distribution
   - Transaction processing times

2. **System Metrics**
   - API response times
   - Error rates by endpoint
   - Database query performance

### Logging Strategy

Following InstrLabs logging standards:

```go
// Payment processing logs
log.Infof("Payment initiated: order_id=%s, amount=%d, payment_type=%s", orderID, amount, paymentType)
log.Errorf("Payment failed: order_id=%s, error=%v", orderID, err)
log.Warnf("Payment webhook validation failed: order_id=%s, reason=%s", orderID, reason)
```

## Deployment Architecture

### Container Configuration

```dockerfile
# Multi-stage build for Go payment service
FROM golang:1.24-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment-service .

FROM alpine:latest AS runtime
RUN apk --no-cache add ca-certificates tzdata dumb-init
WORKDIR /app
COPY --from=builder /go/src/app/payment-service .
EXPOSE 3002
CMD ["dumb-init", "./payment-service"]
```

### Service Discovery

Services register with the gateway through environment-based configuration:

```yaml
# docker-compose.yml
services:
  payment-service:
    build: ./payment-service
    environment:
      - SERVICE_NAME=payment-service
      - PORT=3002
      - MONGODB_URI=mongodb://mongodb:27017
      - MIDTRANS_SERVER_KEY=${MIDTRANS_SERVER_KEY}
    networks:
      - instrlabs-network
    depends_on:
      - mongodb
```

## Scaling Considerations

### Horizontal Scaling

1. **Stateless Design**
   - Payment service designed for horizontal scaling
   - Session data stored in MongoDB
   - No in-memory state dependencies

2. **Database Scaling**
   - MongoDB replica sets for high availability
   - Read replicas for payment status queries
   - Indexing strategy for performance

### Performance Optimization

1. **Caching Strategy**
   - Payment method configuration caching
   - Transaction status caching (TTL-based)
   - Midtrans API response caching where appropriate

2. **Async Processing**
   - Webhook processing via background workers
   - Notification queuing for non-critical updates
   - Batch processing for refunds

## Integration Testing Strategy

### Test Environments

1. **Sandbox Environment**
   - Full Midtrans sandbox API integration
   - Mock payment flows for all payment methods
   - Comprehensive webhook testing

2. **Mock Environment**
   - Local development with mocked Midtrans API
   - Unit tests for business logic
   - Integration tests with other services

### Test Scenarios

1. **Happy Path Testing**
   - Successful payment flows for all methods
   - Proper status updates and notifications
   - Accurate financial record keeping

2. **Error Scenario Testing**
   - Network failures and timeouts
   - Invalid payment details
   - Midtrans service unavailability

This architecture provides a robust, scalable, and maintainable foundation for integrating Midtrans payment gateway into the InstrLabs microservices ecosystem.