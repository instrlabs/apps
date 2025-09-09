package initx

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	MongoURI string
	MongoDB  string
}

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongo(cfg *MongoConfig) *Mongo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	db := client.Database(cfg.MongoDB)
	return &Mongo{Client: client, DB: db}
}

func (m *Mongo) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := m.Client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect from MongoDB: %v", err)
	}
}
