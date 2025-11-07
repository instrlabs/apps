# Midtrans Payment Gateway API Specification

## Overview

This document provides comprehensive API specifications for the Midtrans payment gateway integration in the InstrLabs payment service. It includes all endpoints, request/response formats, error codes, and implementation examples.

## Base URL

- **Sandbox**: `https://api.sandbox.midtrans.com/v2`
- **Production**: `https://api.midtrans.com/v2`
- **Payment Service**: `http://localhost:3002/api/v1` (development)

## Authentication

### Service-to-Service Authentication
- **Type**: HTTP Basic Authentication
- **Username**: Midtrans Server Key
- **Password**: Empty (Server Key acts as username)
- **Header**: `Authorization: Basic BASE64(SERVER_KEY:)`

### Client Authentication
- **Type**: Bearer Token (JWT)
- **Header**: `Authorization: Bearer <JWT_TOKEN>`
- **Validation**: Gateway service validates JWT before forwarding to payment service

## API Endpoints

### Payment Service Endpoints

#### 1. Process Payment

**Endpoint**: `POST /api/v1/payments/charge`

**Description**: Initiates a new payment transaction using the specified payment method.

**Request Headers**:
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body**:
```json
{
  "order_id": "ORDER-2024-001",
  "user_id": "user123",
  "payment_type": "bank_transfer",
  "gross_amount": 100000,
  "currency": "IDR",
  "customer_details": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+62812345678"
  },
  "payment_details": {
    "bank": "bca",
    "va_number": "1234567890"
  },
  "expires_at": "2024-12-31T23:59:59Z",
  "metadata": {
    "product_id": "prod123",
    "category": "electronics"
  }
}
```

**Response Examples**:

**Success Response (200)**:
```json
{
  "message": "Payment processed successfully",
  "errors": null,
  "data": {
    "order_id": "ORDER-2024-001",
    "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "payment_type": "bank_transfer",
    "status": "pending",
    "gross_amount": 100000,
    "currency": "IDR",
    "transaction_time": "2024-12-07T10:30:00Z",
    "va_numbers": [
      {
        "bank": "bca",
        "va_number": "1234567890"
      }
    ],
    "expires_at": "2024-12-07T10:30:00Z"
  }
}
```

**Error Response (400)**:
```json
{
  "message": "Invalid request parameters",
  "errors": [
    {
      "field": "order_id",
      "code": "REQUIRED",
      "message": "Order ID is required"
    }
  ],
  "data": null
}
```

#### 2. Get Payment Status

**Endpoint**: `GET /api/v1/payments/{order_id}/status`

**Description**: Retrieves the current status of a payment transaction.

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters**:
- `order_id` (string, required): Unique order identifier

**Response Examples**:

**Success Response (200)**:
```json
{
  "message": "Payment status retrieved successfully",
  "errors": null,
  "data": {
    "order_id": "ORDER-2024-001",
    "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "payment_type": "bank_transfer",
    "status": "settlement",
    "gross_amount": 100000,
    "currency": "IDR",
    "transaction_time": "2024-12-07T10:30:00Z",
    "payment_time": "2024-12-07T10:45:00Z",
    "fraud_status": "accept"
  }
}
```

**Not Found Response (404)**:
```json
{
  "message": "Payment not found",
  "errors": [
    {
      "field": "order_id",
      "code": "NOT_FOUND",
      "message": "Payment with order ID ORDER-2024-999 not found"
    }
  ],
  "data": null
}
```

#### 3. Cancel Payment

**Endpoint**: `POST /api/v1/payments/{order_id}/cancel`

**Description**: Cancels a pending payment transaction.

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Path Parameters**:
- `order_id` (string, required): Unique order identifier

**Response Examples**:

**Success Response (200)**:
```json
{
  "message": "Payment canceled successfully",
  "errors": null,
  "data": {
    "order_id": "ORDER-2024-001",
    "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "status": "cancel",
    "message": "Payment canceled successfully"
  }
}
```

#### 4. Get Payment Methods

**Endpoint**: `GET /api/v1/payments/methods`

**Description**: Retrieves available payment methods and their configurations.

**Request Headers**:
```
Authorization: Bearer <JWT_TOKEN>
```

**Query Parameters**:
- `enabled` (boolean, optional): Filter by enabled status

**Response Examples**:

**Success Response (200)**:
```json
{
  "message": "Payment methods retrieved successfully",
  "errors": null,
  "data": {
    "payment_methods": [
      {
        "type": "bank_transfer",
        "provider": "bca",
        "enabled": true,
        "display_name": "BCA Virtual Account",
        "icon_url": "https://example.com/icons/bca.png",
        "min_amount": 10000,
        "max_amount": 100000000,
        "fee": {
          "type": "fixed",
          "amount": 4000
        }
      },
      {
        "type": "ewallet",
        "provider": "gopay",
        "enabled": true,
        "display_name": "GoPay",
        "icon_url": "https://example.com/icons/gopay.png",
        "min_amount": 1000,
        "max_amount": 10000000,
        "fee": {
          "type": "percentage",
          "amount": 2.0
        }
      }
    ]
  }
}
```

### Webhook Endpoints

#### 5. Payment Notification Webhook

**Endpoint**: `POST /api/v1/webhooks/payment-notification`

**Description**: Receives payment status updates from Midtrans.

**Request Headers**:
```
Content-Type: application/json
X-Midtrans-Signature: <SHA512_SIGNATURE>
```

**Request Body**:
```json
{
  "status_code": "200",
  "status_message": "Success, transaction is found",
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "order_id": "ORDER-2024-001",
  "gross_amount": "100000.00",
  "payment_type": "bank_transfer",
  "transaction_time": "2024-12-07T10:30:00Z",
  "transaction_status": "settlement",
  "signature_key": "a1b2c3d4e5f67890abcd1234567890a1b2c3d4e5f67890abcd1234567890a1b2c3d4e5f67890abcd1234567890a1b2c3d4e5f67890abcd1234567890a1b2c3d4e5f67890abcd1234567890",
  "fraud_status": "accept",
  "approval_code": "123456"
}
```

**Response Examples**:

**Success Response (200)**:
```json
{
  "message": "Webhook processed successfully",
  "errors": null,
  "data": {
    "status": "processed"
  }
}
```

### Admin Endpoints

#### 6. Get Transaction History

**Endpoint**: `GET /api/v1/admin/payments`

**Description**: Retrieves transaction history for administrative purposes.

**Request Headers**:
```
Authorization: Bearer <ADMIN_JWT_TOKEN>
```

**Query Parameters**:
- `user_id` (string, optional): Filter by user ID
- `status` (string, optional): Filter by transaction status
- `payment_type` (string, optional): Filter by payment type
- `start_date` (string, optional): Start date (ISO 8601)
- `end_date` (string, optional): End date (ISO 8601)
- `limit` (integer, optional): Number of records (default: 50)
- `offset` (integer, optional): Number of records to skip (default: 0)

**Response Examples**:

**Success Response (200)**:
```json
{
  "message": "Transaction history retrieved successfully",
  "errors": null,
  "data": {
    "transactions": [
      {
        "order_id": "ORDER-2024-001",
        "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "user_id": "user123",
        "payment_type": "bank_transfer",
        "status": "settlement",
        "gross_amount": 100000,
        "currency": "IDR",
        "created_at": "2024-12-07T10:30:00Z",
        "updated_at": "2024-12-07T10:45:00Z"
      }
    ],
    "pagination": {
      "total": 150,
      "limit": 50,
      "offset": 0,
      "has_more": true
    }
  }
}
```

## Payment Method Specifications

### 1. Credit Card Payment

**Payment Type**: `credit_card`

**Request Example**:
```json
{
  "order_id": "ORDER-2024-001",
  "payment_type": "credit_card",
  "gross_amount": 100000,
  "customer_details": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+62812345678"
  },
  "payment_details": {
    "card_number": "4811111111111114",
    "card_cvv": "123",
    "card_exp_month_year": "12/24",
    "secure": true,
    "save_token": false
  }
}
```

**Response Example**:
```json
{
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status_code": "200",
  "status_message": "Success, Credit Card transaction is successful",
  "transaction_status": "capture",
  "fraud_status": "accept",
  "approval_code": "123456"
}
```

### 2. Bank Transfer (Virtual Account)

**Payment Type**: `bank_transfer`

**Supported Banks**: `bca`, `bni`, `bri`, `mandiri`, `permata`, `cimb`, `danamon`

**Request Example**:
```json
{
  "order_id": "ORDER-2024-001",
  "payment_type": "bank_transfer",
  "gross_amount": 100000,
  "customer_details": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+62812345678"
  },
  "payment_details": {
    "bank": "bca",
    "va_number": "1234567890",
    "free_text": {
      "en": "Payment for ORDER-2024-001",
      "id": "Pembayaran untuk ORDER-2024-001"
    }
  }
}
```

**Response Example**:
```json
{
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status_code": "201",
  "status_message": "Success, Transaction is pending customer payment",
  "transaction_status": "pending",
  "va_numbers": [
    {
      "bank": "bca",
      "va_number": "1234567890"
    }
  ],
  "expiry_time": "2024-12-07T23:59:59Z"
}
```

### 3. E-Wallet (GoPay)

**Payment Type**: `ewallet`

**Request Example**:
```json
{
  "order_id": "ORDER-2024-001",
  "payment_type": "ewallet",
  "gross_amount": 100000,
  "customer_details": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+62812345678"
  },
  "payment_details": {
    "gopay": {
      "payment_type": "generate-qr-code",
      "enable_callback": true,
      "callback_url": "https://yourdomain.com/api/v1/webhooks/gopay-callback"
    }
  }
}
```

**Response Example**:
```json
{
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status_code": "201",
  "status_message": "Success, Transaction is pending customer payment",
  "transaction_status": "pending",
  "actions": [
    {
      "name": "generate-qr-code",
      "method": "get",
      "url": "https://api.sandbox.midtrans.com/v2/gopay/qrcode/1234567890",
      "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
    }
  ]
}
```

### 4. Convenience Store

**Payment Type**: `cstore`

**Supported Stores**: `alfamart`, `indomaret`

**Request Example**:
```json
{
  "order_id": "ORDER-2024-001",
  "payment_type": "cstore",
  "gross_amount": 100000,
  "customer_details": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+62812345678"
  },
  "payment_details": {
    "store": "alfamart",
    "message": "Payment for ORDER-2024-001",
    "name": "John Doe",
    "phone": "+62812345678"
  }
}
```

**Response Example**:
```json
{
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status_code": "201",
  "status_message": "Success, Transaction is pending customer payment",
  "transaction_status": "pending",
  "payment_code": "8850820173905015",
  "store": "alfamart",
  "expiry_time": "2024-12-07T23:59:59Z"
}
```

### 5. Cardless Credit

**Payment Type**: `cardless_credit`

**Supported Providers**: `akulaku`, `kredivo`

**Request Example**:
```json
{
  "order_id": "ORDER-2024-001",
  "payment_type": "cardless_credit",
  "gross_amount": 100000,
  "customer_details": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+62812345678"
  },
  "payment_details": {
    "provider": "kredivo",
    "redirect_url": "https://yourdomain.com/payment/complete"
  }
}
```

**Response Example**:
```json
{
  "transaction_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status_code": "201",
  "status_message": "Success, Transaction is pending customer approval",
  "transaction_status": "pending",
  "redirect_url": "https://checkout.kredivo.com/kredivo/v2/pay?token=1234567890",
  "actions": [
    {
      "name": "kredivo-redirect",
      "method": "get",
      "url": "https://checkout.kredivo.com/kredivo/v2/pay?token=1234567890"
    }
  ]
}
```

## Transaction Status Codes

### HTTP Status Codes

| Status Code | Description | Example Scenarios |
|-------------|-------------|-------------------|
| 200 | Success | Transaction completed successfully |
| 201 | Created | Transaction created, pending customer action |
| 400 | Bad Request | Invalid parameters, validation errors |
| 401 | Unauthorized | Invalid authentication credentials |
| 404 | Not Found | Transaction or payment method not found |
| 409 | Conflict | Duplicate order ID |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | System or third-party service errors |

### Transaction Status Values

| Status | Description | Payment Methods |
|--------|-------------|-----------------|
| `capture` | Transaction successful (card payments) | Credit Card |
| `settlement` | Payment completed (non-card payments) | Bank Transfer, E-Wallet, C-Store, Cardless Credit |
| `pending` | Awaiting customer action | All payment methods |
| `deny` | Transaction rejected | All payment methods |
| `expire` | Transaction timeout | Bank Transfer, E-Wallet, C-Store, Cardless Credit |
| `cancel` | Transaction canceled by merchant | All payment methods |
| `refund` | Transaction refunded | All payment methods |

### Fraud Status Values

| Status | Description |
|--------|-------------|
| `accept` | Transaction approved by fraud detection |
| `challenge` | Transaction requires manual review |
| `deny` | Transaction rejected by fraud detection |

## Error Handling

### Error Response Format

```json
{
  "message": "Error description",
  "errors": [
    {
      "field": "field_name",
      "code": "ERROR_CODE",
      "message": "Specific error message for this field"
    }
  ],
  "data": null
}
```

### Common Error Codes

| Code | Description | Example |
|------|-------------|---------|
| `REQUIRED` | Required field is missing | `{"field": "order_id", "code": "REQUIRED", "message": "Order ID is required"}` |
| `INVALID_FORMAT` | Field format is invalid | `{"field": "email", "code": "INVALID_FORMAT", "message": "Invalid email format"}` |
| `MIN_VALUE` | Value below minimum | `{"field": "gross_amount", "code": "MIN_VALUE", "message": "Amount must be at least 1000"}` |
| `MAX_VALUE` | Value exceeds maximum | `{"field": "gross_amount", "code": "MAX_VALUE", "message": "Amount exceeds maximum limit"}` |
| `UNIQUE_VIOLATION` | Duplicate value | `{"field": "order_id", "code": "UNIQUE_VIOLATION", "message": "Order ID already exists"}` |
| `NOT_FOUND` | Resource not found | `{"field": "order_id", "code": "NOT_FOUND", "message": "Payment not found"}` |
| `INVALID_STATUS` | Invalid status transition | `{"field": "status", "code": "INVALID_STATUS", "message": "Cannot cancel completed payment"}` |

### Midtrans API Error Codes

| Status Code | Message | Description |
|-------------|---------|-------------|
| 400 | Bad Request | Invalid request parameters |
| 401 | Access Denied | Invalid API credentials |
| 404 | Not Found | Transaction not found |
| 410 | Gone | Transaction expired |
| 411 | Duplicate Order ID | Order ID already exists |
| 412 | Invalid Payment Type | Payment type not enabled |
| 413 | Amount Out of Range | Amount exceeds limits |
| 414 | Merchant Inactive | Merchant account inactive |
| 500 | Internal Server Error | Midtrans service error |

## Webhook Security

### Signature Verification

Midtrans sends a SHA512 signature in the `X-Midtrans-Signature` header. To verify:

1. Concatenate: `order_id + status_code + gross_amount + server_key`
2. Calculate SHA512 hash
3. Compare with received signature

**Example Verification Code**:
```go
func VerifyWebhookSignature(orderID, statusCode, grossAmount, signatureKey, serverKey string) bool {
    input := orderID + statusCode + grossAmount + serverKey
    hash := sha512.Sum512([]byte(input))
    calculatedSignature := fmt.Sprintf("%x", hash)
    return calculatedSignature == signatureKey
}
```

## Rate Limiting

### API Rate Limits

| Endpoint | Rate Limit | Time Window |
|----------|------------|-------------|
| `/payments/charge` | 100 requests | 1 minute |
| `/payments/{order_id}/status` | 200 requests | 1 minute |
| `/payments/{order_id}/cancel` | 50 requests | 1 minute |
| `/webhooks/payment-notification` | 1000 requests | 1 minute |

### Response Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1702005600
```

## Testing Examples

### cURL Examples

**Process Credit Card Payment**:
```bash
curl -X POST http://localhost:3002/api/v1/payments/charge \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "order_id": "TEST-ORDER-001",
    "user_id": "test-user",
    "payment_type": "credit_card",
    "gross_amount": 100000,
    "currency": "IDR",
    "customer_details": {
      "first_name": "Test",
      "last_name": "User",
      "email": "test@example.com",
      "phone": "+62812345678"
    },
    "payment_details": {
      "card_number": "4811111111111114",
      "card_cvv": "123",
      "card_exp_month_year": "12/24",
      "secure": true
    }
  }'
```

**Get Payment Status**:
```bash
curl -X GET http://localhost:3002/api/v1/payments/TEST-ORDER-001/status \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

**Cancel Payment**:
```bash
curl -X POST http://localhost:3002/api/v1/payments/TEST-ORDER-001/cancel \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### JavaScript Examples

**Process Payment with Fetch API**:
```javascript
const processPayment = async (paymentData) => {
  try {
    const response = await fetch('/api/v1/payments/charge', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify(paymentData)
    });

    const result = await response.json();

    if (response.ok) {
      console.log('Payment processed:', result.data);
      return result.data;
    } else {
      console.error('Payment failed:', result.errors);
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Network error:', error);
    throw error;
  }
};
```

**Check Payment Status**:
```javascript
const checkPaymentStatus = async (orderId) => {
  try {
    const response = await fetch(`/api/v1/payments/${orderId}/status`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${getAuthToken()}`
      }
    });

    const result = await response.json();

    if (response.ok) {
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Error checking status:', error);
    throw error;
  }
};
```

## Implementation Checklist

### Prerequisites
- [ ] Midtrans account (Sandbox and Production)
- [ ] Server Key and Client Key
- [ ] Webhook URL configured in Midtrans dashboard
- [ ] SSL certificate for production environment

### Integration Steps
- [ ] Set up payment service infrastructure
- [ ] Implement core payment endpoints
- [ ] Add webhook handling with signature verification
- [ ] Implement error handling and retry logic
- [ ] Add logging and monitoring
- [ ] Configure rate limiting
- [ ] Set up database indexes
- [ ] Test with sandbox environment
- [ ] Perform load testing
- [ ] Deploy to production
- [ ] Configure production webhooks
- [ ] Monitor and optimize

This API specification provides comprehensive guidance for implementing Midtrans payment gateway integration with detailed examples and best practices.