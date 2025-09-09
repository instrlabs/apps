package internal

import (
	"context"
	"errors"
	"time"

	initx "github.com/histweety-labs/shared/init"
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

func (r *ProductRepository) Create(p *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, p)
	return err
}

// FindByID returns a product by its hex string ID, using ObjectID for the _id filter.
func (r *ProductRepository) FindByID(id string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	var p Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &p, nil
}

// FindByKey finds a product by its key.
func (r *ProductRepository) FindByKey(key string) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p Product
	err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &p, nil
}

// Update sets provided fields and touches updatedAt (camelCase to match Product tags).
func (r *ProductRepository) Update(id string, updateFields bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}
	if updateFields == nil {
		updateFields = bson.M{}
	}
	updateFields["updatedAt"] = time.Now()

	update := bson.M{"$set": updateFields}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// Delete removes a product by ID.
func (r *ProductRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// List returns all products, optionally filtering by active=true.
func (r *ProductRepository) List(onlyActive bool) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if onlyActive {
		filter["active"] = true
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
