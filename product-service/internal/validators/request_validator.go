package validators

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RequestValidator handles input validation for HTTP requests
type RequestValidator struct{}

// NewRequestValidator creates a new RequestValidator
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// ProductListRequest represents the validated parameters for product listing
type ProductListRequest struct {
	Type  string `json:"type"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// ValidateProductListRequest validates and parses product list request parameters
func (v *RequestValidator) ValidateProductListRequest(c *fiber.Ctx) (*ProductListRequest, error) {
	req := &ProductListRequest{
		Type:  c.Query("type", ""),
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 50),
	}

	// Validate product type if provided
	if req.Type != "" {
		if err := v.validateProductType(req.Type); err != nil {
			return nil, err
		}
	}

	// Validate pagination parameters
	if err := v.validatePagination(req.Page, req.Limit); err != nil {
		return nil, err
	}

	return req, nil
}

// validatePagination validates pagination parameters
func (v *RequestValidator) validatePagination(page, limit int) error {
	if page < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Page must be greater than 0")
	}

	if limit < 1 || limit > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Limit must be between 1 and 100")
	}

	return nil
}

// validateProductType validates the product type
func (v *RequestValidator) validateProductType(productType string) error {
	validTypes := []string{"digital", "physical", "service", "subscription"}
	for _, validType := range validTypes {
		if productType == validType {
			return nil
		}
	}

	return fiber.NewError(fiber.StatusBadRequest,
		fmt.Sprintf("Invalid product type '%s'. Valid types are: %v", productType, validTypes))
}

// ValidateProductID validates product ID parameter and returns ObjectID
func (v *RequestValidator) ValidateProductID(c *fiber.Ctx) (primitive.ObjectID, error) {
	productID := c.Params("id")
	if productID == "" {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusBadRequest, "Product ID is required")
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusBadRequest,
			"Product ID must be a valid ObjectID format")
	}

	return objectID, nil
}

// ValidateGetProductRequest validates the get product request parameters
func (v *RequestValidator) ValidateGetProductRequest(c *fiber.Ctx) (*GetProductRequest, error) {
	objectID, err := v.ValidateProductID(c)
	if err != nil {
		return nil, err
	}

	productType := c.Query("type", "")

	// Validate product type if provided
	if productType != "" {
		if err := v.validateProductType(productType); err != nil {
			return nil, err
		}
	}

	return &GetProductRequest{
		ID:   objectID,
		Type: productType,
	}, nil
}

// GetProductRequest represents the validated parameters for getting a single product
type GetProductRequest struct {
	ID   primitive.ObjectID `json:"id"`
	Type string             `json:"type"`
}

// ValidationErrorMessage represents a structured validation error
type ValidationErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationError creates a structured validation error response
func (v *RequestValidator) ValidationError(message string, fieldErrors []ValidationErrorMessage) *fiber.Error {
	return fiber.NewError(fiber.StatusBadRequest, message)
}

// ValidateSearchQuery validates search-related query parameters
func (v *RequestValidator) ValidateSearchQuery(c *fiber.Ctx) (*SearchRequest, error) {
	req := &SearchRequest{
		Query:    c.Query("q", ""),
		Type:     c.Query("type", ""),
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 50),
		MinPrice: c.QueryFloat("min_price", -1),
		MaxPrice: c.QueryFloat("max_price", -1),
	}

	// Validate query length
	if len(req.Query) > 100 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Search query cannot exceed 100 characters")
	}

	// Validate product type if provided
	if req.Type != "" {
		if err := v.validateProductType(req.Type); err != nil {
			return nil, err
		}
	}

	// Validate pagination
	if err := v.validatePagination(req.Page, req.Limit); err != nil {
		return nil, err
	}

	// Validate price range
	if err := v.validatePriceRange(req.MinPrice, req.MaxPrice); err != nil {
		return nil, err
	}

	return req, nil
}

// SearchRequest represents validated search parameters
type SearchRequest struct {
	Query    string  `json:"query"`
	Type     string  `json:"type"`
	Page     int     `json:"page"`
	Limit    int     `json:"limit"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

// validatePriceRange validates the minimum and maximum price parameters
func (v *RequestValidator) validatePriceRange(minPrice, maxPrice float64) error {
	if minPrice < 0 && minPrice != -1 {
		return fiber.NewError(fiber.StatusBadRequest, "Minimum price cannot be negative")
	}

	if maxPrice < 0 && maxPrice != -1 {
		return fiber.NewError(fiber.StatusBadRequest, "Maximum price cannot be negative")
	}

	if minPrice != -1 && maxPrice != -1 && minPrice > maxPrice {
		return fiber.NewError(fiber.StatusBadRequest, "Minimum price cannot be greater than maximum price")
	}

	return nil
}
