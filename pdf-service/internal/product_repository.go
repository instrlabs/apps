package internal

import (
	"context"
	"errors"
	"log"
	"time"

	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *initx.Mongo) *ProductRepository {
	return &ProductRepository{
		collection: db.DB.Collection("products"),
	}
}

func (r *ProductRepository) List(productType string) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"type": productType, "active": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list products of type %s: %v", productType, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("Failed to decode products: %v", err)
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindByID(id primitive.ObjectID, productType string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "type": productType, "active": true}

	var product Product
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Printf("Failed to find product by ID %s: %v", id.Hex(), err)
		return nil, err
	}

	return &product, nil
}
