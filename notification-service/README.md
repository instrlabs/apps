# Notification Service

The Notification Service is responsible for handling real-time notifications between backend services and web clients. It acts as a bridge between the NATS messaging system and WebSocket connections from web clients.

## Architecture

The service follows a simple architecture:

1. **NATS Subscriber**: Listens for notifications on the NATS messaging system
2. **WebSocket Server**: Maintains connections with web clients
3. **Notification Handler**: Processes notifications from NATS and broadcasts them to connected WebSocket clients

```
┌────────────┐     ┌─────────────────────┐     ┌───────────┐
│            │     │                     │     │           │
│ Backend    │────▶│ NATS               │────▶│ Notification │
│ Services   │     │ (job.notifications) │     │ Service    │
│            │     │                     │     │           │
└────────────┘     └─────────────────────┘     └─────┬─────┘
                                                    │
                                                    │ WebSocket
                                                    ▼
                                              ┌───────────┐
                                              │           │
                                              │ Web       │
                                              │ Clients   │
                                              │           │
                                              └───────────┘
```

## Features

- Real-time notification delivery
- Automatic reconnection for WebSocket clients
- Scalable architecture
- Support for different notification types

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `PORT` | HTTP server port | `:3030` |
| `NATS_URL` | NATS server URL | `nats://localhost:4222` |
| `NATS_SUBJECT_JOB_NOTIFICATIONS` | NATS subject for job notifications | `job.notifications` |
| `WEBSOCKET_PATH` | WebSocket endpoint path | `/ws` |

## API

### WebSocket Endpoint

Connect to the WebSocket endpoint to receive real-time notifications:

```
ws://localhost:3030/ws
```

### Notification Format

Notifications are sent as JSON objects with the following format:

```json
{
  "id": "job-123",
  "status": "completed"
}
```

Possible status values:
- `pending`
- `processing`
- `completed`
- `failed`

## Integration with Web Client

The web client connects to the notification service using the WebSocket protocol. The connection is managed by the WebSocket service in the web application.

### Example Usage in React Component

```tsx
import { useNotifications } from '../hooks/useNotifications';

function NotificationDisplay() {
  const { notifications, connected } = useNotifications();
  
  return (
    <div>
      <div>Connection status: {connected ? 'Connected' : 'Disconnected'}</div>
      <h2>Recent Notifications</h2>
      <ul>
        {notifications.map((notification, index) => (
          <li key={index}>
            Job {notification.id}: {notification.status}
          </li>
        ))}
      </ul>
    </div>
  );
}
```

## Development

To run the service locally:

1. Make sure NATS is running
2. Set up environment variables or use defaults
3. Run the service:

```bash
go run main.go
```

## Deployment

The service is deployed as part of the Docker Compose setup. It is configured to connect to the NATS server and expose the WebSocket endpoint.