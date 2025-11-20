package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/instrlabs/product-service/internal/models"
	"github.com/instrlabs/product-service/internal/repositories"
)

// ProductServiceInterface defines the interface for product business logic
type ProductServiceInterface interface {
	ListProducts(productType string, page, limit int) (*ProductListResult, error)
	GetProductByID(id string, productType string) (*models.Product, error)
}

// ProductListResult contains the result of a product listing with pagination
type ProductListResult struct {
	Products   []models.Product   `json:"products"`
	Pagination PaginationMetadata `json:"pagination"`
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

// productService implements ProductServiceInterface
type productService struct {
	repo repositories.ProductRepositoryInterface
}

// NewProductService creates a new ProductService
func NewProductService(repo repositories.ProductRepositoryInterface) ProductServiceInterface {
	return &productService{
		repo: repo,
	}
}

// ListProducts retrieves a paginated list of products
func (s *productService) ListProducts(productType string, page, limit int) (*ProductListResult, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Get products from repository
	products, err := s.repo.List(productType)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Apply pagination
	total := len(products)
	start := (page - 1) * limit
	end := start + limit

	var paginatedProducts []models.Product
	if start >= total {
		paginatedProducts = []models.Product{}
	} else if end > total {
		paginatedProducts = products[start:total]
	} else {
		paginatedProducts = products[start:end]
	}

	// Calculate pagination metadata
	totalPages := (total + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := PaginationMetadata{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	return &ProductListResult{
		Products:   paginatedProducts,
		Pagination: pagination,
	}, nil
}

// GetProductByID retrieves a single product by ID
func (s *productService) GetProductByID(id string, productType string) (*models.Product, error) {
	// Validate product ID format
	if id == "" {
		return nil, fmt.Errorf("product ID is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	// Get product from repository
	product, err := s.repo.FindByID(objectID, productType)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

// ValidateProductType validates if the product type is supported
func (s *productService) ValidateProductType(productType string) error {
	if productType == "" {
		return nil // Empty type means all types
	}

	supportedTypes := []string{"digital", "physical", "service", "subscription"}
	for _, supportedType := range supportedTypes {
		if productType == supportedType {
			return nil
		}
	}

	return fmt.Errorf("unsupported product type: %s", productType)
}
