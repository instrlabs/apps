# Image Service

Image processing service for Instrlabs platform providing file upload, storage, processing pipeline, and instruction management with real-time notifications.

## Features

- **File Management**: Multi-part file upload with S3 storage integration
- **Image Processing**: Automated image processing pipeline with status tracking
- **Instruction Management**: Create and track processing instructions with detailed status
- **Product Integration**: Product-based file organization and management
- **Real-time Notifications**: NATS-based messaging for processing updates
- **Automated Cleanup**: Periodic cleanup of old and processed files
- **Storage Optimization**: S3-compatible storage with configurable endpoints

## Quick Start

```bash
# Local Setup
cp .env.example .env
go mod download
go run main.go

# Docker Setup
docker-compose up -d --build image-service
```

## Configuration

Refer to `.env.example` file for the list of required environment variables and their default values.

## API Endpoints

### Instructions Management

**Create Instruction**
```
POST /instructions
- Create new processing instruction
- Returns instruction ID for tracking
```

**Create Instruction Details**
```
POST /instructions/:id/details
- Add file details to instruction
- Upload files and associate with instruction
- Trigger processing pipeline
```

**Get Instruction**
```
GET /instructions/:id
- Retrieve instruction details by ID
- Include current processing status
```

**List Instructions**
```
GET /instructions
- Paginated list of user instructions
- Filter by status and date range
```

**Get Instruction Details**
```
GET /instructions/:id/details
- List all files associated with instruction
- Include processing status for each file
```

**Get Specific File Detail**
```
GET /instructions/:id/details/:detailId
- Retrieve specific file information
- Include processing metadata and status
```

**Download Processed File**
```
GET /instructions/:id/details/:detailId/file
- Download processed or original file
- Stream directly from S3 storage
```

### Product Management

**List Products**
```
GET /products
- Retrieve available products for current user
- Include product metadata and limits
```

### File Management

**List Uncleaned Files**
```
GET /files
- List files pending cleanup
- Administrative endpoint for file management
```

### Health Check

**Service Health**
```
GET /health
- Service health status and dependencies
- Database and storage connectivity check
```

## Processing Pipeline

### File Processing Stages

1. **Upload** (PENDING)
   - File received and stored in temporary location
   - Basic validation performed
   - Instruction created

2. **Processing** (PROCESSING)
   - Image processing algorithms applied
   - Format conversion and optimization
   - Metadata extraction and storage

3. **Completion** (DONE/FAILED)
   - Processing completed successfully or failed
   - Files moved to permanent storage
   - Notification sent via NATS

### Status Tracking

Each instruction detail tracks:
- Original file information
- Processing status (PENDING, PROCESSING, DONE, FAILED)
- Processed file URLs
- Error messages (if failed)
- Processing timestamps

## Project Structure

```
image-service/
├── main.go                       # Fiber app entry point
├── internal/
│   ├── config.go                 # Configuration management
│   ├── models.go                 # Data models (Instruction, Product, etc.)
│   ├── handlers/
│   │   ├── instruction_handler.go # Instruction CRUD operations
│   │   └── product_handler.go     # Product management
│   ├── repositories/
│   │   ├── instruction_repository.go    # Instruction DB operations
│   │   ├── instruction_detail_repository.go # File detail operations
│   │   └── product_repository.go         # Product DB operations
│   ├── services/
│   │   └── image_service.go       # Image processing logic
│   └── errors.go                 # Error definitions
├── static/
│   └── swagger.json              # API documentation
└── Dockerfile
```

## Storage Configuration

### S3 Integration

The service uses S3-compatible storage for file persistence:

```bash
# S3 Configuration
S3_ENDPOINT=https://s3.amazonaws.com
S3_REGION=us-east-1
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_BUCKET=instrlabs-images
S3_USE_SSL=true
```

**Storage Strategy:**
- Temporary storage during processing
- Permanent storage after processing completion
- Automatic cleanup of failed/pending files (30-minute intervals)

## Message Queue Integration

### NATS Configuration

```bash
# NATS Configuration
NATS_URI=nats://nats:4222
NATS_SUBJECT_IMAGE_REQUESTS=image.requests
NATS_SUBJECT_NOTIFICATIONS_SSE=notifications.sse
```

**Message Flow:**
1. Instruction created → Message published to `image.requests`
2. Processing service consumes message and processes files
3. Status updates sent to `notifications.sse`
4. Real-time updates delivered to users via SSE

## Security

### Authentication

- JWT-based authentication via shared middleware
- User identification through `X-User-ID` header
- Token validation and refresh handling
- Ownership validation for resource access

### Rate Limiting

- Default rate limits applied to all endpoints
- Higher limits for authenticated users
- File upload size restrictions

## Error Handling

### Standard Response Format

```json
{
  "message": "Operation description",
  "errors": null,
  "data": { ... }
}
```

### Error Responses

- **400 Bad Request**: Validation errors, malformed data
- **401 Unauthorized**: Authentication required or invalid token
- **403 Forbidden**: Resource access denied (ownership validation)
- **404 Not Found**: Resource does not exist
- **500 Internal Server Error**: Processing or system errors

## Monitoring and Observability

### Health Monitoring

- Service health endpoint with dependency checks
- Database connectivity monitoring
- Storage service availability checks
- NATS connection status

### Metrics

- Prometheus metrics integration
- Request rate and error tracking
- Processing pipeline metrics
- Storage usage statistics

### Logging

- Structured logging with request tracing
- Processing pipeline event logging
- Error and performance monitoring

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - Database
- [NATS Go Client](https://github.com/nats-io/nats.go) - Message queue
- [AWS SDK for Go](https://github.com/aws/aws-sdk-go) - S3 storage
- [Shared Init](github.com/instrlabs/shared/init) - Common utilities

## Development

### Local Development Setup

```bash
# Prerequisites
go 1.21+
MongoDB instance
S3-compatible storage (MinIO for local dev)
NATS server

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

### Docker Development

```bash
# Build image
docker build -t instrlabs/image-service .

# Run with docker-compose
docker-compose up image-service
```