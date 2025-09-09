# Shared Init Package

This module provides reusable initialization helpers for all services: MongoDB, S3 (MinIO), NATS, and Swagger route for Fiber.

Package path: `github.com/arthadede/shared/init` (package name `initx`).

Contents:
- Config (FromEnv): loads configuration from environment variables with sensible defaults
- NewMongo / (*Mongo).Close: MongoDB client and database handle
- NewS3 / Put / Get: MinIO S3 client helpers
- NewNats / (*Nats).Close: NATS connection
- SetupServiceSwagger: register GET /swagger serving a JSON file

Environment variables recognized:
- ENVIRONMENT, PORT
- MONGO_URI, MONGO_DB
- S3_ENDPOINT, S3_REGION, S3_ACCESS_KEY, S3_SECRET_KEY, S3_BUCKET, S3_USE_SSL
- NATS_URL, NATS_SUBJECT_REQUESTS, NATS_SUBJECT_NOTIFICATIONS

Example usage in a service (Fiber):

```go
import (
    initx "github.com/arthadede/shared/init"
    "github.com/gofiber/fiber/v2"
)

func main() {
    cfg := initx.FromEnv()

    mongo := initx.NewMongo(cfg)
    defer mongo.Close()

    s3, err := initx.NewS3(cfg)
    if err != nil { panic(err) }

    natsConn, err := initx.NewNats(cfg.NatsURL)
    if err != nil { panic(err) }
    defer natsConn.Close()

    app := fiber.New()
    initx.SetupServiceSwagger(app, "./static/swagger.json")

    // ... your routes
    app.Listen(cfg.Port)
}
```
