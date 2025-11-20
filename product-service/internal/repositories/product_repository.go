package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/instrlabs/product-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepositoryInterface defines the interface for the product repository
type ProductRepositoryInterface interface {
	List(productType string) ([]models.Product, error)
	FindByID(id primitive.ObjectID, productType string) (*models.Product, error)
	FindByKey(key string, productType string) (*models.Product, error)
}

type ProductRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		db:         db,
		collection: db.Collection("products"),
	}
}

func (r *ProductRepository) List(productType string) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"active": true}
	if productType != "" {
		filter["type"] = productType
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		log.Printf("Failed to decode products: %v", err)
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindByID(id primitive.ObjectID, productType string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "active": true}
	if productType != "" {
		filter["type"] = productType
	}

	var product models.Product
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

func (r *ProductRepository) FindByKey(key string, productType string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"key": key, "active": true}
	if productType != "" {
		filter["type"] = productType
	}

	var product models.Product
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		log.Printf("Failed to find product by key %s: %v", key, err)
		return nil, err
	}

	return &product, nil
}
