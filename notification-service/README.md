# Notification Service

Real-time notification service for Instrlabs platform providing Server-Sent Events (SSE) for instant message delivery with NATS message bus integration.

## Features

- **Real-time Notifications**: Server-Sent Events (SSE) for instant message delivery
- **NATS Integration**: Message bus subscription for receiving notifications from other services
- **User-Specific Streams**: Individual notification streams per user with isolation
- **Connection Management**: Automatic connection tracking and cleanup
- **Keep-Alive System**: Ping mechanism to maintain connection health
- **JWT Authentication**: Secure user identification and authorization
- **Rate Limiting**: Protection against connection spam and abuse

## Quick Start

```bash
# Local Setup
cp .env.example .env
go mod download
go run main.go

# Docker Setup
docker-compose up -d --build notification-service
```

## Configuration

Refer to `.env.example` file for the list of required environment variables and their default values.

## API Endpoints

### Server-Sent Events

**Connect to Notification Stream**
```
GET /sse
- Establish SSE connection for real-time notifications
- Requires JWT authentication via headers
- Returns text/event-stream content type
- User-specific notification delivery
```

**Required Headers:**
- `Authorization: Bearer <jwt_token>` - JWT authentication token
- `X-User-ID: <user_id>` - User identification (automatically added by auth middleware)

**SSE Event Types:**

1. **Connection Event**
```
event: connected
data: {"message": "Connected to notification stream", "timestamp": "2024-01-01T00:00:00Z"}
```

2. **Ping Event** (every 30 seconds)
```
event: ping
data: {"timestamp": "2024-01-01T00:00:00Z"}
```

3. **Notification Event**
```
event: message
data: {"id": "notif_123", "type": "instruction", "title": "Processing Complete", "body": "Your image processing has finished successfully", "timestamp": "2024-01-01T00:00:00Z"}
```

### Health Check

**Service Health**
```
GET /health
- Service health status and dependencies
- NATS connection status
- Active connection count
```

## Project Structure

```
notification-service/
├── main.go                    # Fiber app entry point
├── internal/
│   ├── config.go              # Configuration management
│   ├── middleware.go          # CORS, rate limiting, logging
│   ├── sse_service.go         # SSE connection management
│   ├── sse_handler.go         # SSE HTTP handlers
│   └── models.go              # Data models and types
├── static/
│   └── swagger.json           # API documentation
└── Dockerfile
```

## Connection Management

### SSE Connection Lifecycle

1. **Connection Establishment**
   - JWT token validation
   - User identification and authentication
   - SSE stream initialization
   - Connection event sent to client

2. **Active Connection**
   - Real-time message delivery from NATS
   - Periodic ping messages (30-second intervals)
   - Connection health monitoring
   - Automatic reconnection handling

3. **Connection Termination**
   - Graceful disconnection handling
   - Resource cleanup
   - Connection pool management

### User Isolation

- Each user receives only their own notifications
- Connections tracked by user ID
- Automatic cleanup of orphaned connections
- Support for multiple concurrent connections per user

## Message Queue Integration

### NATS Configuration

```bash
# NATS Configuration
NATS_URI=nats://nats:4222
NATS_SUBJECT_NOTIFICATIONS_SSE=notifications.sse
```

**Message Flow:**
1. Services publish notifications to `notifications.sse`
2. Notification service subscribes and filters messages by user
3. Messages delivered to appropriate SSE connections
4. Real-time delivery to connected clients

**Message Format:**
```json
{
  "user_id": "user_123",
  "id": "notif_456",
  "type": "instruction",
  "title": "Processing Started",
  "body": "Your image is being processed",
  "data": {"instruction_id": "instr_789"},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Security

### Authentication

- **JWT Token Validation**: Bearer token required for all connections
- **User Identification**: X-User-ID header automatically populated
- **Token Refresh**: Automatic token refresh for long-lived connections
- **CORS Protection**: Configurable origin validation

### Rate Limiting

- Connection rate limiting to prevent abuse
- Request throttling for SSE connections
- Configurable rate limits per endpoint

**Rate Limiting Configuration:**
```bash
# Default limits (configurable via shared middleware)
RATE_LIMIT=100      # requests per window
RATE_WINDOW=60s     # time window
```

## Error Handling

### Standard Response Format

```json
{
  "message": "Error description",
  "errors": ["Detailed error information"],
  "data": null
}
```

### SSE Error Handling

**Connection Errors:**
- **400 Bad Request**: Missing headers or invalid parameters
- **401 Unauthorized**: Invalid or expired JWT token
- **403 Forbidden**: User access denied
- **500 Internal Server Error**: Service or NATS connection issues

**Error Event Format:**
```
event: error
data: {"error": "Authentication failed", "code": 401}
```

## Monitoring and Observability

### Health Monitoring

- Service health endpoint with dependency checks
- NATS connection status monitoring
- Active connection count tracking
- Memory and resource usage monitoring

### Metrics

- Prometheus metrics integration
- Connection establishment and termination rates
- Message delivery success rates
- Error rates by type

### Logging

- Structured logging for connection events
- Message delivery logging
- Error and performance monitoring
- User activity tracking

## Client Integration

### JavaScript Client Example

```javascript
// Connect to notification stream
const eventSource = new EventSource('/sse', {
  headers: {
    'Authorization': `Bearer ${jwtToken}`,
    'X-User-ID': userId
  }
});

eventSource.addEventListener('connected', (event) => {
  console.log('Connected to notifications:', JSON.parse(event.data));
});

eventSource.addEventListener('message', (event) => {
  const notification = JSON.parse(event.data);
  console.log('Received notification:', notification);
  // Handle notification display
});

eventSource.addEventListener('ping', (event) => {
  console.log('Connection alive:', JSON.parse(event.data));
});

eventSource.addEventListener('error', (event) => {
  console.error('SSE error:', event);
  // Implement reconnection logic
});

eventSource.onerror = (event) => {
  console.error('Connection lost:', event);
  // Implement reconnection logic
};
```

### React Hook Example

```javascript
function useNotifications(userId, jwtToken) {
  const [notifications, setNotifications] = useState([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const eventSource = new EventSource('/sse', {
      headers: {
        'Authorization': `Bearer ${jwtToken}`,
        'X-User-ID': userId
      }
    });

    eventSource.addEventListener('connected', () => setConnected(true));

    eventSource.addEventListener('message', (event) => {
      const notification = JSON.parse(event.data);
      setNotifications(prev => [notification, ...prev]);
    });

    eventSource.onerror = () => {
      setConnected(false);
    };

    return () => {
      eventSource.close();
    };
  }, [userId, jwtToken]);

  return { notifications, connected };
}
```

## Performance Considerations

### Scalability

- **Connection Pooling**: Efficient connection management
- **Memory Management**: Automatic cleanup of inactive connections
- **Message Filtering**: User-based filtering at source
- **Resource Limits**: Configurable connection limits per user

### Optimization

- **Binary Message Support**: Efficient message encoding
- **Compression**: Optional response compression
- **Connection Reuse**: Persistent connections where possible
- **Lazy Loading**: On-demand connection establishment

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [NATS Go Client](https://github.com/nats-io/nats.go) - Message queue
- [Shared Init](github.com/instrlabs/shared/init) - Common utilities

## Development

### Local Development Setup

```bash
# Prerequisites
go 1.21+
NATS server instance

# Setup
cp .env.example .env
# Edit .env with local configuration
go mod download
go run main.go
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Testing SSE Connections

```bash
# Test SSE connection with curl
curl -H "Authorization: Bearer <token>" \
     -H "X-User-ID: test-user" \
     http://localhost:3000/sse
```

### Docker Development

```bash
# Build image
docker build -t instrlabs/notification-service .

# Run with docker-compose
docker-compose up notification-service
```

## Troubleshooting

### Common Issues

1. **Connection Failed**
   - Check JWT token validity
   - Verify X-User-ID header presence
   - Ensure NATS server is accessible

2. **No Messages Received**
   - Verify NATS subscription status
   - Check message format compliance
   - Ensure user ID matches in messages

3. **Frequent Disconnections**
   - Check network stability
   - Verify keep-alive settings
   - Monitor server resources

### Debug Logging

Enable debug logging for troubleshooting:

```bash
# Set environment
ENVIRONMENT=development
LOG_LEVEL=debug
```

This will provide detailed connection and message logging.