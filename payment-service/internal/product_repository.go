package internal

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	productCollectionName = "products"
)

type Product struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name"`
	Price       float64   `bson:"price" json:"price"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Image       string    `bson:"image,omitempty" json:"image,omitempty"`
	ProductType string    `bson:"productType,omitempty" json:"productType,omitempty"`
	Active      bool      `bson:"active" json:"active"`
	IsFree      bool      `bson:"isFree" json:"isFree"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

type ProductRepository struct {
	db         *MongoDB
	collection *mongo.Collection
}

func NewProductRepository(db *MongoDB) *ProductRepository {
	collection := db.GetCollection(productCollectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName("idx_products_name"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		fmt.Printf("Error creating index on products: %v\n", err)
	}

	return &ProductRepository{db: db, collection: collection}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *Product) error {
	if p.ID == "" {
		p.ID = primitive.NewObjectID().Hex()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	var p Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return &p, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id string, updateFields bson.M) error {
	if updateFields == nil {
		updateFields = bson.M{}
	}
	updateFields["updatedAt"] = time.Now()
	update := bson.M{"$set": updateFields}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// DeleteProduct removes a product by ID
func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// ListProducts returns all products (optionally filter by active)
func (r *ProductRepository) ListProducts(ctx context.Context, onlyActive bool) ([]*Product, error) {
	filter := bson.M{}
	if onlyActive {
		filter["active"] = true
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []*Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}
	return products, nil
}
