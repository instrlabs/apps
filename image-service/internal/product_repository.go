package internal

import (
	"context"
	"errors"
	"time"

	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func NewProductRepository(db *initx.Mongo) *ProductRepository {
	return &ProductRepository{
		db:         db,
		collection: db.DB.Collection("products"),
	}
}

func (r *ProductRepository) List() ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"product_type": "images",
		"is_active":    true,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// FindByKey returns an active 'images' product by its key or nil if not found.
func (r *ProductRepository) FindByKey(key string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"key":          key,
		"product_type": "images",
		"is_active":    true,
	}
	var p Product
	err := r.collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// FindByID returns an active 'images' product by its ObjectID or nil if not found.
func (r *ProductRepository) FindByID(id primitive.ObjectID) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":          id,
		"product_type": "images",
		"is_active":    true,
	}
	var p Product
	err := r.collection.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
