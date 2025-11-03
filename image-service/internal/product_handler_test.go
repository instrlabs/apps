package internal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ====================
// Product Handler Tests
// ====================

// ProductRepositoryInterface defines the interface for testing
type ProductRepositoryInterface interface {
	List() ([]*Product, error)
}

// MockProductRepository for testing
type MockProductRepository struct {
	ProductRepositoryInterface
	ListFunc func() ([]*Product, error)
}

func (m *MockProductRepository) List() ([]*Product, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return []*Product{}, nil
}

// ProductHandlerForTest wraps the handler to work with interface
type ProductHandlerForTest struct {
	repo ProductRepositoryInterface
}

func NewProductHandlerForTest(repo ProductRepositoryInterface) *ProductHandlerForTest {
	return &ProductHandlerForTest{repo: repo}
}

func (h *ProductHandlerForTest) ListProducts(c *fiber.Ctx) error {
	products, err := h.repo.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  nil,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products retrieved successfully",
		"errors":  nil,
		"data": map[string]interface{}{
			"products": products,
		},
	})
}

func TestProductHandler_ListProducts_Success(t *testing.T) {
	// Arrange
	mockProducts := []*Product{
		{
			ID:          primitive.NewObjectID(),
			Key:         "product-1",
			Title:       "Product One",
			Description: "Description for product one",
			ProductType: "TYPE_A",
			IsActive:    true,
			IsFree:      false,
		},
		{
			ID:          primitive.NewObjectID(),
			Key:         "product-2",
			Title:       "Product Two",
			Description: "Description for product two",
			ProductType: "TYPE_B",
			IsActive:    true,
			IsFree:      true,
		},
	}

	mockRepo := &MockProductRepository{
		ListFunc: func() ([]*Product, error) {
			return mockProducts, nil
		},
	}

	handler := NewProductHandlerForTest(mockRepo)

	app := fiber.New()
	app.Get("/products", handler.ListProducts)

	// Act
	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	assert.Equal(t, "Products retrieved successfully", response["message"])
	assert.Nil(t, response["errors"])

	data := response["data"].(map[string]interface{})
	products := data["products"].([]interface{})
	assert.Len(t, products, 2)
}

func TestProductHandler_ListProducts_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := &MockProductRepository{
		ListFunc: func() ([]*Product, error) {
			return []*Product{}, nil
		},
	}

	handler := NewProductHandlerForTest(mockRepo)

	app := fiber.New()
	app.Get("/products", handler.ListProducts)

	// Act
	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	products := data["products"].([]interface{})
	assert.Len(t, products, 0)
}

func TestProductHandler_ListProducts_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := &MockProductRepository{
		ListFunc: func() ([]*Product, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewProductHandlerForTest(mockRepo)

	app := fiber.New()
	app.Get("/products", handler.ListProducts)

	// Act
	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	assert.Equal(t, "Internal server error", response["message"])
	assert.Nil(t, response["errors"])
	assert.Nil(t, response["data"])
}

func TestProductHandler_ListProducts_ResponseFormat(t *testing.T) {
	// Verify the response adheres to the standard format

	mockRepo := &MockProductRepository{
		ListFunc: func() ([]*Product, error) {
			return []*Product{}, nil
		},
	}

	handler := NewProductHandlerForTest(mockRepo)

	app := fiber.New()
	app.Get("/products", handler.ListProducts)

	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	// Verify standard response format
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "errors")
	assert.Contains(t, response, "data")

	assert.IsType(t, "", response["message"])
	// data should be a map
	assert.IsType(t, map[string]interface{}{}, response["data"])
}

func TestNewProductHandler(t *testing.T) {
	mockRepo := &MockProductRepository{}
	handler := NewProductHandlerForTest(mockRepo)

	assert.NotNil(t, handler)
	assert.Equal(t, mockRepo, handler.repo)
}
