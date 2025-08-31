# Payment Service

A microservice for handling payments integrated with Midtrans payment gateway, built with the Fiber web framework.

## Features

- Create payments with Midtrans (supports top-up, product purchase, and product subscription)
- Check payment status
- Handle payment notifications (webhooks)
- Process payment requests via NATS
- Store payment data in MongoDB
- Fast and efficient API using Fiber web framework

## API Endpoints

### Create Payment (Top-up or Product)

```
POST /payments
```

Request body:
```json
{
  "orderId": "",
  "userId": "USER-123",
  "amount": 100000,
  "currency": "IDR",
  "paymentMethod": "gopay",
  "description": "Payment for order #123",
  "customerName": "John Doe",
  "customerEmail": "john@example.com",
  "callbackUrl": "https://example.com/callback",
  "type": "product"
}
```

Notes:
- Set `type` to "product" (default) or "topup".
- If `orderId` is empty, the service will generate one with a prefix according to `type` (TOPUP-, PROD-).
- Default `type` is `product` if not provided.

Response:
```json
{
  "id": "PAY-123",
  "orderId": "ORDER-123",
  "userId": "USER-123",
  "amount": 100000,
  "currency": "IDR",
  "paymentMethod": "gopay",
  "status": "pending",
  "redirectUrl": "https://app.midtrans.com/snap/v2/vtweb/..."
}
```

### Create Subscription

```
POST /payments/subscriptions
```

Request body:
```json
{
  "name": "Pro Plan Monthly",
  "userId": "USER-123",
  "amount": 50000,
  "currency": "IDR",
  "token": "snap_credit_card_token",
  "interval": "month",
  "intervalCount": 1,
  "startAt": "", 
  "description": "Subscription for Pro Plan",
  "customerName": "John Doe",
  "customerEmail": "john@example.com"
}
```

Notes:
- Midtrans subscription API requires a credit card token.
- `startAt` can be empty to start immediately or an RFC3339 timestamp.

Response:
```json
{
  "id": "sub_xxx",
  "orderId": "sub_xxx",
  "userId": "USER-123",
  "amount": 50000,
  "currency": "IDR",
  "status": "pending",
  "type": "subscription"
}
```

### Get Payment Status

```
GET /payments/{orderId}
```

Response:
```json
{
  "id": "PAY-123",
  "orderId": "ORDER-123",
  "amount": 100000,
  "status": "success",
  "paymentMethod": "gopay"
}
```

### Handle Notification

```
POST /payments/notification
```

This endpoint is used by Midtrans to send payment notifications.

## NATS Integration

The service subscribes to the `payment.requests` subject to process payment requests from other services.

Payment request message format:
```json
{
  "orderId": "ORDER-123",
  "userId": "USER-123",
  "amount": 100000,
  "currency": "IDR",
  "paymentMethod": "gopay",
  "description": "Payment for order #123",
  "callbackUrl": "https://example.com/callback",
  "type": "product"
}
```

Notes:
- Set `type` to "product" or "topup". Subscription is created via HTTP at `/payments/subscriptions`. 

The service publishes payment events to the `payment.events` subject.

Payment event message format:
```json
{
  "id": "PAY-123",
  "orderId": "ORDER-123",
  "userId": "USER-123",
  "amount": 100000,
  "currency": "IDR",
  "paymentMethod": "gopay",
  "status": "pending",
  "redirectUrl": "https://app.midtrans.com/snap/v2/vtweb/...",
  "timestamp": "2025-08-29T18:11:00Z",
  "type": "product"
}
```

## Configuration

The service is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| ENVIRONMENT | Environment (development, production) | development |
| PORT | HTTP server port | :3040 |
| MONGO_URI | MongoDB connection URI | mongodb://localhost:27017 |
| MONGO_DB | MongoDB database name | payment_service |
| NATS_URL | NATS server URL | nats://localhost:4222 |
| NATS_SUBJECT_PAYMENT_EVENTS | NATS subject for payment events | payment.events |
| NATS_SUBJECT_PAYMENT_REQUESTS | NATS subject for payment requests | payment.requests |
| MIDTRANS_SERVER_KEY | Midtrans server key | - |
| MIDTRANS_CLIENT_KEY | Midtrans client key | - |
| MIDTRANS_ENVIRONMENT | Midtrans environment (sandbox, production) | sandbox |
| MIDTRANS_NOTIFICATION_URL | URL for Midtrans notifications | - |
| CORS_ALLOWED_ORIGINS | CORS allowed origins | http://web.localhost |

## Getting Started

1. Set up environment variables in `.env.local`
2. Run the service using Docker Compose:
   ```
   docker-compose up payment-service
   ```

## Development

This service uses the [Fiber](https://gofiber.io/) web framework, which is an Express-inspired web framework for Go that focuses on performance and minimal memory allocation.

To run the service locally:

1. Install dependencies:
   ```
   go mod tidy
   ```

2. Run the service:
   ```
   go run main.go
   ```

Fiber provides several benefits:
- Fast HTTP routing
- Low memory footprint
- Built-in middleware support
- Express-like API for easy development