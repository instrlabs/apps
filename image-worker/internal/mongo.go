package internal

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(cfg *Config) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")

	db := client.Database(cfg.MongoDB)

	return &MongoDB{Client: client, DB: db}
}

func (m *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}
