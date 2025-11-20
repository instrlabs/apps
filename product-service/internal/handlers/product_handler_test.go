package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockProductRepository mocks the ProductRepository for testing
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) List(productType string) ([]Product, error) {
	args := m.Called(productType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Product), args.Error(1)
}

func (m *MockProductRepository) FindByID(id primitive.ObjectID, productType string) (*Product, error) {
	args := m.Called(id, productType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func createTestProduct(id primitive.ObjectID, key, name, desc, productType string, price float64, active bool) Product {
	return Product{
		ID:        id,
		Key:       key,
		Name:      name,
		Desc:      desc,
		Type:      productType,
		Price:     price,
		Active:    active,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// Test setup function that creates a test app with routes
func setupTestApp(handler *ProductHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/api/v1/products", handler.ListProducts)
	app.Get("/api/v1/products/:id", handler.GetProductByID)

	return app
}

// Helper function to test and parse response
func testRequest(app *fiber.App, method, path string) (*http.Response, map[string]interface{}, error) {
	req := httptest.NewRequest(method, path, nil)
	resp, err := app.Test(req)
	if err != nil {
		return nil, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)

	return resp, response, err
}

// ListProducts Tests
func TestListProducts_Success(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	product1 := createTestProduct(
		primitive.NewObjectID(),
		"product1",
		"Test Product 1",
		"Description 1",
		"physical",
		29.99,
		true,
	)

	product2 := createTestProduct(
		primitive.NewObjectID(),
		"product2",
		"Test Product 2",
		"Description 2",
		"digital",
		19.99,
		true,
	)

	expectedProducts := []Product{product1, product2}

	repo.On("List", "").Return(expectedProducts, nil)

	t.Run("List all products", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "Success", response["message"])
		assert.Nil(t, response["errors"])
		assert.NotNil(t, response["data"])

		productsData := response["data"].(map[string]interface{})["products"].([]interface{})
		assert.Equal(t, 2, len(productsData))

		repo.AssertExpectations(t)
	})
}

func TestListProducts_WithFilter(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	product := createTestProduct(
		primitive.NewObjectID(),
		"physical-product",
		"Physical Product",
		"Physical description",
		"physical",
		29.99,
		true,
	)

	expectedProducts := []Product{product}

	repo.On("List", "physical").Return(expectedProducts, nil)

	t.Run("List products by type filter", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products?type=physical")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "Success", response["message"])
		assert.Nil(t, response["errors"])
		assert.NotNil(t, response["data"])

		productsData := response["data"].(map[string]interface{})["products"].([]interface{})
		assert.Equal(t, 1, len(productsData))
		assert.Equal(t, "Physical Product", productsData[0].(map[string]interface{})["name"])

		repo.AssertExpectations(t)
	})
}

func TestListProducts_EmptyResult(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	repo.On("List", "").Return([]Product{}, nil)

	t.Run("List products with no results", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "Success", response["message"])
		assert.Nil(t, response["errors"])
		assert.NotNil(t, response["data"])

		productsData := response["data"].(map[string]interface{})["products"].([]interface{})
		assert.Equal(t, 0, len(productsData))

		repo.AssertExpectations(t)
	})
}

func TestListProducts_DatabaseError(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	repo.On("List", "").Return(nil, errors.New("database connection failed"))

	t.Run("Database error handling", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "Failed to fetch products", response["message"])
		assert.NotNil(t, response["errors"])
		assert.Nil(t, response["data"])

		repo.AssertExpectations(t)
	})
}

func TestListProducts_ResponseFormat(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	product := createTestProduct(
		primitive.NewObjectID(),
		"format-test-product",
		"Format Test Product",
		"Format test description",
		"test",
		15.99,
		true,
	)

	expectedProducts := []Product{product}
	repo.On("List", "").Return(expectedProducts, nil)

	t.Run("Response format validation", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Check response structure
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "errors")
		assert.Contains(t, response, "data")

		// Check message
		assert.Equal(t, "Success", response["message"])

		// Check errors
		assert.Nil(t, response["errors"])

		// Check data structure
		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "products")

		// Check product structure
		products := data["products"].([]interface{})
		assert.Equal(t, 1, len(products))

		product := products[0].(map[string]interface{})
		assert.Contains(t, product, "id")
		assert.Contains(t, product, "key")
		assert.Contains(t, product, "name")
		assert.Contains(t, product, "desc")
		assert.Contains(t, product, "type")
		assert.Contains(t, product, "price")
		assert.Contains(t, product, "active")
		assert.Contains(t, product, "created_at")
		assert.Contains(t, product, "updated_at")

		// Check data types
		assert.IsType(t, "", product["id"])
		assert.IsType(t, "", product["key"])
		assert.IsType(t, "", product["name"])
		assert.IsType(t, "", product["desc"])
		assert.IsType(t, "", product["type"])
		assert.IsType(t, float64(0.0), product["price"])
		assert.IsType(t, false, product["active"])

		repo.AssertExpectations(t)
	})
}

// GetProductByID Tests
func TestGetProductByID_Success(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	productID := primitive.NewObjectID()
	product := createTestProduct(
		productID,
		"test-product",
		"Test Product",
		"Test description",
		"digital",
		29.99,
		true,
	)

	repo.On("FindByID", productID, "").Return(&product, nil)

	t.Run("Get product by ID successfully", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products/"+productID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "Success", response["message"])
		assert.Nil(t, response["errors"])
		assert.NotNil(t, response["data"])

		productData := response["data"].(map[string]interface{})["product"].(map[string]interface{})
		assert.Equal(t, "Test Product", productData["name"])
		assert.Equal(t, "test-product", productData["key"])
		assert.Equal(t, 29.99, productData["price"])

		repo.AssertExpectations(t)
	})
}

func TestGetProductByID_WithFilter(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	productID := primitive.NewObjectID()
	product := createTestProduct(
		productID,
		"physical-product",
		"Physical Product",
		"Physical description",
		"physical",
		49.99,
		true,
	)

	repo.On("FindByID", productID, "physical").Return(&product, nil)

	t.Run("Get product by ID with type filter", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products/"+productID.Hex()+"?type=physical")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "Success", response["message"])
		assert.NotNil(t, response["data"])

		productData := response["data"].(map[string]interface{})["product"].(map[string]interface{})
		assert.Equal(t, "Physical Product", productData["name"])
		assert.Equal(t, "physical", productData["type"])

		repo.AssertExpectations(t)
	})
}

func TestGetProductByID_NotFound(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	productID := primitive.NewObjectID()

	repo.On("FindByID", productID, "").Return(nil, nil)

	t.Run("Product not found", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products/"+productID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		assert.Equal(t, "Product not found", response["message"])
		assert.Equal(t, "Product with the specified ID not found", response["errors"])
		assert.Nil(t, response["data"])

		repo.AssertExpectations(t)
	})
}

func TestGetProductByID_DatabaseError(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	productID := primitive.NewObjectID()

	repo.On("FindByID", productID, "").Return(nil, errors.New("database connection failed"))

	t.Run("Database error handling", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products/"+productID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "Failed to fetch product", response["message"])
		assert.Contains(t, response["errors"], "database connection failed")
		assert.Nil(t, response["data"])

		repo.AssertExpectations(t)
	})
}

func TestGetProductByID_InvalidID(t *testing.T) {
	invalidIDs := []string{
		"invalid-id",
		"123",
		"zbcdef1234567890abcdef123456", // Invalid hex character
		"abcdef1234567890abcdef1",      // Too short
	}

	for _, invalidID := range invalidIDs {
		t.Run("Invalid ID format: "+invalidID, func(t *testing.T) {
			// We need a fresh mock for each test case
			repo := &MockProductRepository{}
			handler := NewProductHandler(repo)
			app := setupTestApp(handler)

			resp, response, err := testRequest(app, "GET", "/api/v1/products/"+invalidID)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			assert.Equal(t, "Invalid product ID format", response["message"])
			assert.Equal(t, "Product ID must be a valid ObjectID", response["errors"])
			assert.Nil(t, response["data"])
		})
	}
}

// Pagination Tests
func TestListProducts_WithPagination(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	// Create test data
	var expectedProducts []Product
	for i := 0; i < 25; i++ {
		product := createTestProduct(
			primitive.NewObjectID(),
			"product"+strconv.Itoa(i),
			"Product "+strconv.Itoa(i),
			"Description "+strconv.Itoa(i),
			"test",
			float64(10.00+float64(i)),
			true,
		)
		expectedProducts = append(expectedProducts, product)
	}

	repo.On("List", "").Return(expectedProducts, nil)

	t.Run("Pagination - first page", func(t *testing.T) {
		resp, response, err := testRequest(app, "GET", "/api/v1/products?page=1&limit=10")
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "Success", response["message"])
		assert.NotNil(t, response["data"])
		assert.NotNil(t, response["pagination"])

		// Check products
		productsData := response["data"].(map[string]interface{})["products"].([]interface{})
		assert.Equal(t, 10, len(productsData))

		// Check pagination metadata
		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["current_page"])
		assert.Equal(t, float64(10), pagination["per_page"])
		assert.Equal(t, float64(25), pagination["total"])
		assert.Equal(t, float64(3), pagination["total_pages"])
		assert.Equal(t, true, pagination["has_next"])
		assert.Equal(t, false, pagination["has_prev"])

		repo.AssertExpectations(t)
	})
}

func TestListProducts_Pagination_InvalidParameters(t *testing.T) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	// Create test data
	var expectedProducts []Product
	for i := 0; i < 10; i++ {
		product := createTestProduct(
			primitive.NewObjectID(),
			"product"+strconv.Itoa(i),
			"Product "+strconv.Itoa(i),
			"Description "+strconv.Itoa(i),
			"test",
			float64(10.00+float64(i)),
			true,
		)
		expectedProducts = append(expectedProducts, product)
	}

	repo.On("List", "").Return(expectedProducts, nil)

	testCases := []struct {
		name          string
		url           string
		expectedPage  int
		expectedLimit int
		expectedItems int
	}{
		{
			name:          "Negative page defaults to 1",
			url:           "/api/v1/products?page=-1&limit=10",
			expectedPage:  1,
			expectedLimit: 10,
			expectedItems: 10,
		},
		{
			name:          "Zero page defaults to 1",
			url:           "/api/v1/products?page=0&limit=10",
			expectedPage:  1,
			expectedLimit: 10,
			expectedItems: 10,
		},
		{
			name:          "Negative limit defaults to 50",
			url:           "/api/v1/products?page=1&limit=-5",
			expectedPage:  1,
			expectedLimit: 50,
			expectedItems: 10,
		},
		{
			name:          "Limit over 100 defaults to 50",
			url:           "/api/v1/products?page=1&limit=150",
			expectedPage:  1,
			expectedLimit: 50,
			expectedItems: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, response, err := testRequest(app, "GET", tc.url)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)

			// Check products
			productsData := response["data"].(map[string]interface{})["products"].([]interface{})
			assert.Equal(t, tc.expectedItems, len(productsData))

			// Check pagination metadata
			pagination := response["pagination"].(map[string]interface{})
			assert.Equal(t, float64(tc.expectedPage), pagination["current_page"])
			assert.Equal(t, float64(tc.expectedLimit), pagination["per_page"])
		})
	}

	repo.AssertExpectations(t)
}

// Benchmark Tests
func BenchmarkListProducts_Success(b *testing.B) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	// Create test data
	var expectedProducts []Product
	for i := 0; i < 100; i++ {
		product := createTestProduct(
			primitive.NewObjectID(),
			"product"+strconv.Itoa(i),
			"Product "+strconv.Itoa(i),
			"Description "+strconv.Itoa(i),
			"type"+strconv.Itoa(i%5),
			float64(10.00+float64(i)),
			true,
		)
		expectedProducts = append(expectedProducts, product)
	}

	repo.On("List", "").Return(expectedProducts, nil)

	req := httptest.NewRequest("GET", "/api/v1/products", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}

	repo.AssertExpectations(b)
}

func BenchmarkListProducts_WithFilter(b *testing.B) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	var expectedProducts []Product
	for i := 0; i < 50; i++ {
		product := createTestProduct(
			primitive.NewObjectID(),
			"digital-product"+strconv.Itoa(i),
			"Digital Product "+strconv.Itoa(i),
			"Digital description "+strconv.Itoa(i),
			"digital",
			float64(20.00+float64(i)),
			true,
		)
		expectedProducts = append(expectedProducts, product)
	}

	repo.On("List", "digital").Return(expectedProducts, nil)

	req := httptest.NewRequest("GET", "/api/v1/products?type=digital", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}

	repo.AssertExpectations(b)
}

func BenchmarkGetProductByID_Success(b *testing.B) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	productID := primitive.NewObjectID()
	product := createTestProduct(
		productID,
		"benchmark-product",
		"Benchmark Product",
		"Benchmark description",
		"test",
		99.99,
		true,
	)

	repo.On("FindByID", productID, "").Return(&product, nil)

	req := httptest.NewRequest("GET", "/api/v1/products/"+productID.Hex(), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}

	repo.AssertExpectations(b)
}

func BenchmarkGetProductByID_WithFilter(b *testing.B) {
	repo := &MockProductRepository{}
	handler := NewProductHandler(repo)
	app := setupTestApp(handler)

	productID := primitive.NewObjectID()
	product := createTestProduct(
		productID,
		"filtered-product",
		"Filtered Product",
		"Filtered description",
		"physical",
		149.99,
		true,
	)

	repo.On("FindByID", productID, "physical").Return(&product, nil)

	req := httptest.NewRequest("GET", "/api/v1/products/"+productID.Hex()+"?type=physical", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req)
	}

	repo.AssertExpectations(b)
}
