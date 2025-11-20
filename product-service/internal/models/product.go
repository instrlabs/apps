package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the system
type Product struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key       string             `json:"key" bson:"key"`
	Name      string             `json:"name" bson:"name"`
	Desc      string             `json:"desc" bson:"desc"`
	Type      string             `json:"type" bson:"type"`
	Price     float64            `json:"price" bson:"price"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewProduct creates a new product with validation
func NewProduct(key, name, desc, productType string, price float64) *Product {
	now := time.Now()
	return &Product{
		Key:       key,
		Name:      name,
		Desc:      desc,
		Type:      productType,
		Price:     price,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the product data
func (p *Product) Validate() error {
	if p.Key == "" {
		return fmt.Errorf("product key is required")
	}

	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if p.Type == "" {
		return fmt.Errorf("product type is required")
	}

	if p.Price < 0 {
		return fmt.Errorf("product price cannot be negative")
	}

	// Validate key format (alphanumeric with dashes and underscores)
	if !isValidKey(p.Key) {
		return fmt.Errorf("product key contains invalid characters")
	}

	return nil
}

// isValidKey checks if the product key is valid (alphanumeric, dash, underscore)
func isValidKey(key string) bool {
	for _, char := range key {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	return true
}

// IsActive checks if the product is currently active
func (p *Product) IsActive() bool {
	return p.Active
}

// UpdateTimestamp updates the updated_at timestamp
func (p *Product) UpdateTimestamp() {
	p.UpdatedAt = time.Now()
}

// GetDisplayPrice returns the price formatted for display
func (p *Product) GetDisplayPrice() string {
	return fmt.Sprintf("%.2f", p.Price)
}

// IsValidType checks if the product type is valid
func (p *Product) IsValidType() bool {
	validTypes := []string{"digital", "physical", "service", "subscription"}
	for _, validType := range validTypes {
		if p.Type == validType {
			return true
		}
	}
	return false
}
