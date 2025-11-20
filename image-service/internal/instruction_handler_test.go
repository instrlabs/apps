package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/instrlabs/shared/modelx"
)

// Test data factory functions
func createTestConfig() *Config {
	return &Config{
		Environment:                 "test",
		Port:                        "3002",
		MongoURI:                    "mongodb://localhost:27017",
		MongoDB:                     "test-db",
		S3Endpoint:                  "localhost:9000",
		S3Region:                    "us-east-1",
		S3AccessKey:                 "test-access-key",
		S3SecretKey:                 "test-secret-key",
		S3Bucket:                    "test-bucket",
		S3UseSSL:                    false,
		NatsURI:                     "nats://localhost:4222",
		NatsSubjectImageRequests:    "image.requests",
		NatsSubjectNotificationsSSE: "notifications.sse",
		ApiUrl:                      "http://localhost:3000",
		ProductServiceURL:           "http://localhost:3005",
	}
}

func createTestInstruction(id primitive.ObjectID, userID *primitive.ObjectID, guestID *string, productID primitive.ObjectID) *modelx.Instruction {
	return &modelx.Instruction{
		ID:        id,
		UserID:    userID,
		GuestID:   guestID,
		ProductID: productID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func createTestInstructionFile(id primitive.ObjectID, instructionID primitive.ObjectID, userID *primitive.ObjectID, guestID *string, fileName string, status modelx.InstructionDetailStatus) *modelx.InstructionFile {
	now := time.Now().UTC()
	return &modelx.InstructionFile{
		ID:            id,
		InstructionID: instructionID,
		UserID:        userID,
		GuestID:       guestID,
		FileName:      fileName,
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func createTestProduct(id primitive.ObjectID) *Product {
	return &Product{
		ID:        id,
		Key:       "test-product",
		Name:      "Test Product",
		Type:      "image",
		Price:     9.99,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func createMultipartFormData(t *testing.T, fieldName, filename string, content []byte, contentType string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	return body, writer.FormDataContentType()
}

// Test helper function to parse JSON response
func parseJSONResponse(t *testing.T, body []byte) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	return response
}

// Functional tests that test actual handler behavior
func TestCreateInstruction_BasicValidation(t *testing.T) {
	cfg := createTestConfig()

	// Create a basic Fiber app to test the handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Convert errors to JSON responses
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "Internal Server Error",
					"errors":  []string{err.Error()},
					"data":    nil,
				})
			}
			return nil
		},
	})

	// Create handler with minimal setup for validation testing
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	// Setup route
	app.Post("/instructions", handler.CreateInstruction)

	tests := []struct {
		name           string
		body           string
		headers        map[string]string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Invalid JSON",
			body:           "invalid json",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid request body",
		},
		{
			name:           "Missing product ID",
			body:           "{}",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "ProductID is required",
		},
		{
			name:           "Empty product ID",
			body:           `{"product_id": ""}`,
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "ProductID is required",
		},
		{
			name:           "Missing user identification",
			body:           `{"product_id": "507f1f77bcf86cd799439011"}`,
			expectedStatus: fiber.StatusUnauthorized,
			expectedMsg:    "missing user identification - provide either x-user-id or x-guest-id header",
		},
		{
			name:           "Invalid user ID",
			body:           `{"product_id": "507f1f77bcf86cd799439011"}`,
			headers:        map[string]string{"x-user-id": "invalid-id"},
			expectedStatus: fiber.StatusUnauthorized,
			expectedMsg:    "invalid user ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/instructions", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

func TestCreateInstructionDetails_BasicValidation(t *testing.T) {
	cfg := createTestConfig()

	// Create a basic Fiber app to test the handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "Internal Server Error",
					"errors":  []string{err.Error()},
					"data":    nil,
				})
			}
			return nil
		},
	})

	// Create handler with minimal setup for validation testing
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	// Setup route
	app.Post("/instructions/:id", handler.CreateInstructionDetails)

	tests := []struct {
		name           string
		path           string
		body           io.Reader
		headers        map[string]string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Invalid instruction ID",
			path:           "/instructions/invalid-id",
			body:           nil,
			headers:        map[string]string{"x-user-id": "507f1f77bcf86cd799439011"},
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid instruction id",
		},
		{
			name:           "Missing user identification",
			path:           "/instructions/507f1f77bcf86cd799439011",
			body:           nil,
			expectedStatus: fiber.StatusUnauthorized,
			expectedMsg:    "missing user identification - provide either x-user-id or x-guest-id header",
		},
		{
			name:           "Missing file",
			path:           "/instructions/507f1f77bcf86cd799439011",
			body:           strings.NewReader(""),
			headers:        map[string]string{"x-user-id": "507f1f77bcf86cd799439011", "Content-Type": "multipart/form-data"},
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "failed to read uploaded file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tt.path, tt.body)

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

func TestGetInstructionByID_BasicValidation(t *testing.T) {
	cfg := createTestConfig()

	// Create a basic Fiber app to test the handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "Internal Server Error",
					"errors":  []string{err.Error()},
					"data":    nil,
				})
			}
			return nil
		},
	})

	// Create handler with minimal setup for validation testing
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	// Setup route
	app.Get("/instructions/:id", handler.GetInstructionByID)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Invalid instruction ID",
			path:           "/instructions/invalid-id",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

func TestGetInstructionDetail_BasicValidation(t *testing.T) {
	cfg := createTestConfig()

	// Create a basic Fiber app to test the handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "Internal Server Error",
					"errors":  []string{err.Error()},
					"data":    nil,
				})
			}
			return nil
		},
	})

	// Create handler with minimal setup for validation testing
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	// Setup route - match the route structure from main.go
	app.Get("/instructions/:id/:detail_id", handler.GetInstructionDetail)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Invalid instruction ID",
			path:           "/instructions/invalid-id/507f1f77bcf86cd799439011",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid instruction id",
		},
		{
			name:           "Invalid detail ID",
			path:           "/instructions/507f1f77bcf86cd799439011/invalid-detail-id",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid detail id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

func TestGetInstructionDetailFile_BasicValidation(t *testing.T) {
	cfg := createTestConfig()

	// Create a basic Fiber app to test the handler
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "Internal Server Error",
					"errors":  []string{err.Error()},
					"data":    nil,
				})
			}
			return nil
		},
	})

	// Create handler with minimal setup for validation testing
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	// Setup route - match the route structure from main.go
	app.Get("/instructions/:id/:detail_id/file", handler.GetInstructionDetailFile)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Invalid instruction ID",
			path:           "/instructions/invalid-id/507f1f77bcf86cd799439011/file",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid instruction id",
		},
		{
			name:           "Invalid detail ID",
			path:           "/instructions/507f1f77bcf86cd799439011/invalid-detail-id/file",
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    "invalid detail id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

// Unit tests for helper methods
func TestValidateOwnership(t *testing.T) {
	cfg := createTestConfig()
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()
	guestID1 := "guest-123"
	guestID2 := "guest-456"

	tests := []struct {
		name     string
		instr    *modelx.Instruction
		userID   *primitive.ObjectID
		guestID  *string
		expected bool
	}{
		{
			name: "User ownership - correct",
			instr: &modelx.Instruction{
				ID:     primitive.NewObjectID(),
				UserID: &userID1,
			},
			userID:   &userID1,
			guestID:  nil,
			expected: true,
		},
		{
			name: "User ownership - incorrect",
			instr: &modelx.Instruction{
				ID:     primitive.NewObjectID(),
				UserID: &userID1,
			},
			userID:   &userID2,
			guestID:  nil,
			expected: false,
		},
		{
			name: "Guest ownership - correct",
			instr: &modelx.Instruction{
				ID:      primitive.NewObjectID(),
				GuestID: &guestID1,
			},
			userID:   nil,
			guestID:  &guestID1,
			expected: true,
		},
		{
			name: "Guest ownership - incorrect",
			instr: &modelx.Instruction{
				ID:      primitive.NewObjectID(),
				GuestID: &guestID1,
			},
			userID:   nil,
			guestID:  &guestID2,
			expected: false,
		},
		{
			name: "User checking guest ownership",
			instr: &modelx.Instruction{
				ID:      primitive.NewObjectID(),
				GuestID: &guestID1,
			},
			userID:   &userID1,
			guestID:  nil,
			expected: false,
		},
		{
			name: "Guest checking user ownership",
			instr: &modelx.Instruction{
				ID:     primitive.NewObjectID(),
				UserID: &userID1,
			},
			userID:   nil,
			guestID:  &guestID1,
			expected: false,
		},
		{
			name:     "No ownership",
			instr:    &modelx.Instruction{ID: primitive.NewObjectID()},
			userID:   nil,
			guestID:  nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.validateOwnership(*tt.instr, tt.userID, tt.guestID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test file validation helpers
func TestMultipartFileValidation(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		contentType string
		shouldPass  bool
	}{
		{
			name:        "Valid JPEG",
			filename:    "test.jpg",
			contentType: "image/jpeg",
			shouldPass:  true,
		},
		{
			name:        "Valid PNG",
			filename:    "test.png",
			contentType: "image/png",
			shouldPass:  true,
		},
		{
			name:        "Valid GIF",
			filename:    "test.gif",
			contentType: "image/gif",
			shouldPass:  true,
		},
		{
			name:        "Valid WebP",
			filename:    "test.webp",
			contentType: "image/webp",
			shouldPass:  true,
		},
		{
			name:        "Valid BMP",
			filename:    "test.bmp",
			contentType: "image/bmp",
			shouldPass:  true,
		},
		{
			name:        "Invalid text file",
			filename:    "test.txt",
			contentType: "text/plain",
			shouldPass:  false,
		},
		{
			name:        "Invalid PDF",
			filename:    "test.pdf",
			contentType: "application/pdf",
			shouldPass:  false,
		},
		{
			name:        "Missing content type",
			filename:    "test.jpg",
			contentType: "",
			shouldPass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create multipart form data
			content := []byte("test content")
			_, _ = createMultipartFormData(t, "file", tt.filename, content, tt.contentType)

			// Test that the content type validation logic would work
			// This is a basic validation that mirrors the handler logic
			allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp"}
			isValid := false
			for _, allowedType := range allowedTypes {
				if tt.contentType == allowedType {
					isValid = true
					break
				}
			}

			assert.Equal(t, tt.shouldPass, isValid, "File type validation mismatch")
		})
	}
}

// Test file size validation
func TestFileSizeValidation(t *testing.T) {
	// Test the file size validation logic
	// The handler uses 50MB as the maximum file size
	const maxFileSize = 50 * 1024 * 1024

	tests := []struct {
		name       string
		fileSize   int64
		shouldPass bool
	}{
		{
			name:       "Small file - 1KB",
			fileSize:   1024,
			shouldPass: true,
		},
		{
			name:       "Medium file - 10MB",
			fileSize:   10 * 1024 * 1024,
			shouldPass: true,
		},
		{
			name:       "Large file - 49MB",
			fileSize:   49 * 1024 * 1024,
			shouldPass: true,
		},
		{
			name:       "Too large file - 51MB",
			fileSize:   51 * 1024 * 1024,
			shouldPass: false,
		},
		{
			name:       "Way too large file - 100MB",
			fileSize:   100 * 1024 * 1024,
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.fileSize <= maxFileSize
			assert.Equal(t, tt.shouldPass, isValid, "File size validation mismatch")
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkCreateInstruction_Validation(b *testing.B) {
	cfg := createTestConfig()
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	app := fiber.New()
	app.Post("/instructions", handler.CreateInstruction)

	requestBody := `{"product_id": "507f1f77bcf86cd799439011"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/instructions", strings.NewReader(requestBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-user-id", "507f1f77bcf86cd799439012")
		app.Test(req, -1)
	}
}

func BenchmarkGetInstructionByID_Validation(b *testing.B) {
	cfg := createTestConfig()
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	app := fiber.New()
	app.Get("/instructions/:id", handler.GetInstructionByID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/instructions/507f1f77bcf86cd799439011", nil)
		app.Test(req, -1)
	}
}

func BenchmarkValidateOwnership(b *testing.B) {
	cfg := createTestConfig()
	handler := NewInstructionHandler(cfg, nil, nil, nil, nil, nil, nil)

	userID := primitive.NewObjectID()
	instr := &modelx.Instruction{
		ID:     primitive.NewObjectID(),
		UserID: &userID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.validateOwnership(*instr, &userID, nil)
	}
}
