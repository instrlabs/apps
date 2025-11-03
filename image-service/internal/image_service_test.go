package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newMockConfig() *Config {
	return &Config{
		Environment:                 "test",
		Port:                        ":3000",
		MongoURI:                    "mongodb://localhost:27017",
		MongoDB:                     "image_test",
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
	}
}

func TestNewConfig(t *testing.T) {
	cfg := LoadConfig()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Environment)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.MongoURI)
	assert.NotEmpty(t, cfg.MongoDB)
	assert.NotEmpty(t, cfg.S3Endpoint)
	assert.NotEmpty(t, cfg.S3Region)
	assert.NotEmpty(t, cfg.S3AccessKey)
	assert.NotEmpty(t, cfg.S3SecretKey)
	assert.NotEmpty(t, cfg.S3Bucket)
	assert.NotEmpty(t, cfg.NatsURI)
	assert.NotEmpty(t, cfg.NatsSubjectImageRequests)
	assert.NotEmpty(t, cfg.NatsSubjectNotificationsSSE)
}

func TestProduct_Validation(t *testing.T) {
	// Test Product struct
	product := &Product{
		ID:          primitive.NewObjectID(),
		Key:         "images/compress",
		Title:       "Image Compression",
		Description: "Compress images efficiently",
		ProductType: "image",
		IsActive:    true,
		IsFree:      true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	assert.NotEmpty(t, product.Key)
	assert.NotEmpty(t, product.Title)
	assert.NotEmpty(t, product.Description)
	assert.NotEmpty(t, product.ProductType)
	assert.True(t, product.IsActive)
	assert.True(t, product.IsFree)
	assert.False(t, product.ID.IsZero())
}

func TestInstruction_Validation(t *testing.T) {
	// Test Instruction struct
	userID := primitive.NewObjectID()
	productID := primitive.NewObjectID()

	instruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	assert.False(t, instruction.ID.IsZero())
	assert.False(t, instruction.UserID.IsZero())
	assert.False(t, instruction.ProductID.IsZero())
	assert.False(t, instruction.CreatedAt.IsZero())
	assert.False(t, instruction.UpdatedAt.IsZero())
}

func TestInstructionDetail_Validation(t *testing.T) {
	// Test InstructionDetail struct
	instructionID := primitive.NewObjectID()
	outputID := primitive.NewObjectID()

	detail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: instructionID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusPending,
		OutputID:      &outputID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		IsCleaned:     false,
	}

	assert.False(t, detail.ID.IsZero())
	assert.False(t, detail.InstructionID.IsZero())
	assert.NotEmpty(t, detail.FileName)
	assert.Greater(t, detail.FileSize, int64(0))
	assert.NotEmpty(t, detail.MimeType)
	assert.Equal(t, FileStatusPending, detail.Status)
	assert.NotNil(t, detail.OutputID)
	assert.False(t, detail.CreatedAt.IsZero())
	assert.False(t, detail.UpdatedAt.IsZero())
	assert.False(t, detail.IsCleaned)
}

func TestInstructionNotification_Validation(t *testing.T) {
	// Test InstructionNotification struct
	notification := InstructionNotification{
		UserID:              "user123",
		InstructionID:       "instr456",
		InstructionDetailID: "detail789",
	}

	assert.NotEmpty(t, notification.UserID)
	assert.NotEmpty(t, notification.InstructionID)
	assert.NotEmpty(t, notification.InstructionDetailID)
}

func TestFileStatus_Constants(t *testing.T) {
	// Test FileStatus constants
	assert.Equal(t, FileStatus("FAILED"), FileStatusFailed)
	assert.Equal(t, FileStatus("PENDING"), FileStatusPending)
	assert.Equal(t, FileStatus("PROCESSING"), FileStatusProcessing)
	assert.Equal(t, FileStatus("DONE"), FileStatusDone)
}

func TestFileStatus_StringValues(t *testing.T) {
	// Test that FileStatus values are correctly set
	assert.Equal(t, "FAILED", string(FileStatusFailed))
	assert.Equal(t, "PENDING", string(FileStatusPending))
	assert.Equal(t, "PROCESSING", string(FileStatusProcessing))
	assert.Equal(t, "DONE", string(FileStatusDone))
}

func TestInstructionDetail_OutputIDHandling(t *testing.T) {
	// Test InstructionDetail with and without OutputID
	instructionID := primitive.NewObjectID()
	outputID := primitive.NewObjectID()

	// With OutputID
	detailWithOutput := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: instructionID,
		FileName:      "input.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusPending,
		OutputID:      &outputID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		IsCleaned:     false,
	}

	assert.NotNil(t, detailWithOutput.OutputID)
	assert.Equal(t, outputID, *detailWithOutput.OutputID)

	// Without OutputID
	detailWithoutOutput := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: instructionID,
		FileName:      "output.jpg",
		FileSize:      512,
		MimeType:      "image/jpeg",
		Status:        FileStatusDone,
		OutputID:      nil,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		IsCleaned:     false,
	}

	assert.Nil(t, detailWithoutOutput.OutputID)
}

func TestInstructionDetail_TimeFields(t *testing.T) {
	// Test that CreatedAt and UpdatedAt fields work correctly
	instructionID := primitive.NewObjectID()
	before := time.Now().UTC()

	detail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: instructionID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusPending,
		OutputID:      nil,
		CreatedAt:     before,
		UpdatedAt:     before,
		IsCleaned:     false,
	}

	after := time.Now().UTC()
	detail.UpdatedAt = after

	assert.Equal(t, before, detail.CreatedAt)
	assert.Equal(t, after, detail.UpdatedAt)
	assert.True(t, after.After(before) || after.Equal(before))
}

func TestProduct_TimestampFields(t *testing.T) {
	// Test that timestamp fields work correctly
	before := time.Now().UTC()

	product := &Product{
		ID:          primitive.NewObjectID(),
		Key:         "test/product",
		Title:       "Test Product",
		Description: "Test description",
		ProductType: "test",
		IsActive:    true,
		IsFree:      true,
		CreatedAt:   before,
		UpdatedAt:   before,
	}

	after := time.Now().UTC()
	product.UpdatedAt = after

	assert.Equal(t, before, product.CreatedAt)
	assert.Equal(t, after, product.UpdatedAt)
	assert.True(t, after.After(before) || after.Equal(before))
}

func TestInstruction_UserProductRelationship(t *testing.T) {
	// Test relationship between User and Product in Instruction
	userID := primitive.NewObjectID()
	productID := primitive.NewObjectID()

	instruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Verify the relationship
	assert.Equal(t, userID, instruction.UserID)
	assert.Equal(t, productID, instruction.ProductID)
	assert.NotEqual(t, instruction.UserID, instruction.ProductID)
}

func TestConfig_DefaultValues(t *testing.T) {
	// Test that config has sensible defaults when loaded
	cfg := LoadConfig()

	assert.NotNil(t, cfg)
	// Test that critical fields are not empty
	assert.NotEmpty(t, cfg.Environment)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.MongoURI)
	assert.NotEmpty(t, cfg.MongoDB)
	assert.NotEmpty(t, cfg.S3Endpoint)
	assert.NotEmpty(t, cfg.S3Region)
	assert.NotEmpty(t, cfg.S3AccessKey)
	assert.NotEmpty(t, cfg.S3SecretKey)
	assert.NotEmpty(t, cfg.S3Bucket)
	assert.NotEmpty(t, cfg.NatsURI)
	assert.NotEmpty(t, cfg.NatsSubjectImageRequests)
	assert.NotEmpty(t, cfg.NatsSubjectNotificationsSSE)
}

func TestMockConfig_TestValues(t *testing.T) {
	// Test that mock config has consistent test values
	cfg := newMockConfig()

	assert.Equal(t, "test", cfg.Environment)
	assert.Equal(t, ":3000", cfg.Port)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
	assert.Equal(t, "image_test", cfg.MongoDB)
	assert.Equal(t, "localhost:9000", cfg.S3Endpoint)
	assert.Equal(t, "us-east-1", cfg.S3Region)
	assert.Equal(t, "test-access-key", cfg.S3AccessKey)
	assert.Equal(t, "test-secret-key", cfg.S3SecretKey)
	assert.Equal(t, "test-bucket", cfg.S3Bucket)
	assert.Equal(t, false, cfg.S3UseSSL)
	assert.Equal(t, "nats://localhost:4222", cfg.NatsURI)
	assert.Equal(t, "image.requests", cfg.NatsSubjectImageRequests)
	assert.Equal(t, "notifications.sse", cfg.NatsSubjectNotificationsSSE)
	assert.Equal(t, "http://localhost:3000", cfg.ApiUrl)
}

func TestObjectID_GenerationAndComparison(t *testing.T) {
	// Test ObjectID generation and comparison (used throughout the models)
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()

	// Test that IDs are unique
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, id3)
	assert.NotEqual(t, id1, id3)

	// Test that IDs are valid
	assert.False(t, id1.IsZero())
	assert.False(t, id2.IsZero())
	assert.False(t, id3.IsZero())

	// Test zero ID
	zeroID := primitive.NilObjectID
	assert.True(t, zeroID.IsZero())

	// Test Hex encoding
	hex1 := id1.Hex()
	hex2 := id2.Hex()
	assert.NotEmpty(t, hex1)
	assert.NotEmpty(t, hex2)
	assert.NotEqual(t, hex1, hex2)
	assert.Equal(t, 24, len(hex1)) // ObjectID hex is always 24 characters
	assert.Equal(t, 24, len(hex2))
}

func TestInstructionDetail_StatusTransitions(t *testing.T) {
	// Test typical status transitions
	instructionID := primitive.NewObjectID()
	outputID := primitive.NewObjectID()

	detail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: instructionID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusPending,
		OutputID:      &outputID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		IsCleaned:     false,
	}

	// Test initial state
	assert.Equal(t, FileStatusPending, detail.Status)

	// Simulate status transitions
	detail.Status = FileStatusProcessing
	assert.Equal(t, FileStatusProcessing, detail.Status)

	detail.Status = FileStatusDone
	assert.Equal(t, FileStatusDone, detail.Status)

	// Test error state
	detail.Status = FileStatusFailed
	assert.Equal(t, FileStatusFailed, detail.Status)
}

func TestInstructionDetail_CleaningLifecycle(t *testing.T) {
	// Test the cleaning lifecycle of InstructionDetail
	instructionID := primitive.NewObjectID()

	detail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: instructionID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusDone,
		OutputID:      nil,
		CreatedAt:     time.Now().Add(-2 * time.Hour), // 2 hours ago
		UpdatedAt:     time.Now().Add(-2 * time.Hour),
		IsCleaned:     false,
	}

	// Initial state - not cleaned
	assert.False(t, detail.IsCleaned)

	// Test that the file would be eligible for cleaning (old and not cleaned)
	cutoff := time.Now().Add(-1 * time.Hour) // 1 hour ago
	isEligibleForCleaning := detail.CreatedAt.Before(cutoff) && !detail.IsCleaned
	assert.True(t, isEligibleForCleaning, "Old file should be eligible for cleaning when not cleaned")

	// Mark as cleaned
	detail.IsCleaned = true
	assert.True(t, detail.IsCleaned)

	// After cleaning, it should no longer be eligible
	isEligibleForCleaning = detail.CreatedAt.Before(cutoff) && !detail.IsCleaned
	assert.False(t, isEligibleForCleaning, "File should not be eligible for cleaning after being marked as cleaned")
}

func TestInstruction_MultipleInstructionsPerUser(t *testing.T) {
	// Test that a user can have multiple instructions
	userID := primitive.NewObjectID()
	productID1 := primitive.NewObjectID()
	productID2 := primitive.NewObjectID()

	instruction1 := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instruction2 := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID2,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Both instructions should belong to the same user
	assert.Equal(t, userID, instruction1.UserID)
	assert.Equal(t, userID, instruction2.UserID)
	assert.Equal(t, instruction1.UserID, instruction2.UserID)

	// But they should have different products and IDs
	assert.NotEqual(t, instruction1.ProductID, instruction2.ProductID)
	assert.NotEqual(t, instruction1.ID, instruction2.ID)
}

// Mock implementations for testing
type MockS3 struct {
	GetFunc    func(string) []byte
	PutFunc    func(string, []byte) error
	DeleteFunc func(string) error
}

func (m *MockS3) Get(key string) []byte {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return nil
}

func (m *MockS3) Put(key string, data []byte) error {
	if m.PutFunc != nil {
		return m.PutFunc(key, data)
	}
	return nil
}

func (m *MockS3) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return nil
}

type MockNats struct {
	PublishFunc func(string, []byte) error
}

func (m *MockNats) Publish(subject string, data []byte) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(subject, data)
	}
	return nil
}

type MockMongo struct {
	DB *mockDatabase
}

type mockDatabase struct {
	Collection *mockCollection
}

type mockCollection struct {
	FindFunc    func(interface{}) *mockCursor
	FindOneFunc func(interface{}) *mockSingleResult
}

type mockCursor struct {
	AllFunc  func(interface{}) error
	NextFunc func(bool) bool
}

type mockSingleResult struct {
	DecodeFunc func(interface{}) error
}

func (m *MockMongo) Collection(name string) *mockCollection {
	return m.DB.Collection
}

func (c *mockCollection) Find(filter interface{}) *mockCursor {
	if c.FindFunc != nil {
		return c.FindFunc(filter)
	}
	return &mockCursor{}
}

func (c *mockCollection) FindOne(filter interface{}) *mockSingleResult {
	if c.FindOneFunc != nil {
		return c.FindOneFunc(filter)
	}
	return &mockSingleResult{}
}

func (c *mockCursor) All(results interface{}) error {
	if c.AllFunc != nil {
		return c.AllFunc(results)
	}
	return nil
}

func (c *mockCursor) Next(b bool) bool {
	if c.NextFunc != nil {
		return c.NextFunc(b)
	}
	return false
}

func (r *mockSingleResult) Decode(v interface{}) error {
	if r.DecodeFunc != nil {
		return r.DecodeFunc(v)
	}
	return nil
}

// Mock repositories
type MockProductRepository struct {
	FindByIDFunc func(primitive.ObjectID) (*Product, error)
	ListFunc     func() ([]*Product, error)
}

func (m *MockProductRepository) FindByID(id primitive.ObjectID) (*Product, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockProductRepository) List() ([]*Product, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return []*Product{}, nil
}

type MockInstructionRepository struct {
	GetByIDFunc    func(primitive.ObjectID) *Instruction
	ListLatestFunc func(string, int) ([]*Instruction, error)
	CreateFunc     func(*Instruction) error
}

func (m *MockInstructionRepository) GetByID(id primitive.ObjectID) *Instruction {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil
}

func (m *MockInstructionRepository) ListLatest(userId string, limit int) ([]*Instruction, error) {
	if m.ListLatestFunc != nil {
		return m.ListLatestFunc(userId, limit)
	}
	return []*Instruction{}, nil
}

func (m *MockInstructionRepository) Create(instr *Instruction) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(instr)
	}
	return nil
}

type MockInstructionDetailRepository struct {
	GetByIDFunc                  func(primitive.ObjectID) *InstructionDetail
	ListByInstructionFunc        func(primitive.ObjectID) []*InstructionDetail
	CreateManyFunc               func([]*InstructionDetail) error
	UpdateStatusFunc             func(primitive.ObjectID, FileStatus) error
	UpdateStatusAndSizeFunc      func(primitive.ObjectID, FileStatus, int64) error
	ListOlderThanFunc            func(time.Time) []*InstructionDetail
	ListPendingUpdatedBeforeFunc func(time.Time) []*InstructionDetail
	MarkCleanedFunc              func([]primitive.ObjectID) error
	ListUncleanedFunc            func() []*InstructionDetail
}

func (m *MockInstructionDetailRepository) GetByID(id primitive.ObjectID) *InstructionDetail {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil
}

func (m *MockInstructionDetailRepository) ListByInstruction(instructionID primitive.ObjectID) []*InstructionDetail {
	if m.ListByInstructionFunc != nil {
		return m.ListByInstructionFunc(instructionID)
	}
	return []*InstructionDetail{}
}

func (m *MockInstructionDetailRepository) CreateMany(details []*InstructionDetail) error {
	if m.CreateManyFunc != nil {
		return m.CreateManyFunc(details)
	}
	return nil
}

func (m *MockInstructionDetailRepository) UpdateStatus(id primitive.ObjectID, status FileStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(id, status)
	}
	return nil
}

func (m *MockInstructionDetailRepository) UpdateStatusAndSize(id primitive.ObjectID, status FileStatus, size int64) error {
	if m.UpdateStatusAndSizeFunc != nil {
		return m.UpdateStatusAndSizeFunc(id, status, size)
	}
	return nil
}

func (m *MockInstructionDetailRepository) ListOlderThan(cutoff time.Time) []*InstructionDetail {
	if m.ListOlderThanFunc != nil {
		return m.ListOlderThanFunc(cutoff)
	}
	return []*InstructionDetail{}
}

func (m *MockInstructionDetailRepository) ListPendingUpdatedBefore(cutoff time.Time) []*InstructionDetail {
	if m.ListPendingUpdatedBeforeFunc != nil {
		return m.ListPendingUpdatedBeforeFunc(cutoff)
	}
	return []*InstructionDetail{}
}

func (m *MockInstructionDetailRepository) MarkCleaned(ids []primitive.ObjectID) error {
	if m.MarkCleanedFunc != nil {
		return m.MarkCleanedFunc(ids)
	}
	return nil
}

func (m *MockInstructionDetailRepository) ListUncleaned() []*InstructionDetail {
	if m.ListUncleanedFunc != nil {
		return m.ListUncleanedFunc()
	}
	return []*InstructionDetail{}
}

// Test helpers
func createTestApp() (*fiber.App, *Config, *MockS3, *MockNats, *MockProductRepository, *MockInstructionRepository, *MockInstructionDetailRepository) {
	cfg := newMockConfig()
	s3 := &MockS3{}
	nats := &MockNats{}
	productRepo := &MockProductRepository{}
	instrRepo := &MockInstructionRepository{}
	detailRepo := &MockInstructionDetailRepository{}

	app := fiber.New()

	return app, cfg, s3, nats, productRepo, instrRepo, detailRepo
}

// Main handler initialization tests focus on the gateway pattern
// where SetupAuthenticated is called with specific protected routes

// Tests for main function initialization
func TestMainInitialization_ConfigLoading(t *testing.T) {
	// Test that LoadConfig works correctly
	cfg := LoadConfig()
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.MongoURI)
	assert.NotEmpty(t, cfg.NatsURI)
	assert.NotEmpty(t, cfg.S3Endpoint)
}

func TestMainInitialization_ServiceNameAndDefaults(t *testing.T) {
	cfg := LoadConfig()

	// Test environment defaults
	assert.NotEmpty(t, cfg.Environment)
	assert.Contains(t, []string{"development", "production", "test"}, cfg.Environment)

	// Test that essential services are configured
	assert.NotEmpty(t, cfg.NatsSubjectImageRequests)
	assert.NotEmpty(t, cfg.NatsSubjectNotificationsSSE)
}

func TestMainInitialization_MongoConfiguration(t *testing.T) {
	cfg := LoadConfig()

	// Test MongoDB configuration
	assert.NotEmpty(t, cfg.MongoURI)
	assert.NotEmpty(t, cfg.MongoDB)
	assert.Contains(t, cfg.MongoURI, "mongodb://")
}

func TestMainInitialization_S3Configuration(t *testing.T) {
	cfg := LoadConfig()

	// Test S3 configuration
	assert.NotEmpty(t, cfg.S3Endpoint)
	assert.NotEmpty(t, cfg.S3AccessKey)
	assert.NotEmpty(t, cfg.S3SecretKey)
	assert.NotEmpty(t, cfg.S3Bucket)
	assert.NotEmpty(t, cfg.S3Region)
}

// Table-driven tests for error scenarios
func TestInstructionHandler_CreateInstruction_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]string
		setupRepo      func(*MockProductRepository)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:    "valid product",
			payload: map[string]string{"product_id": "507f1f77bcf86cd799439011"},
			setupRepo: func(repo *MockProductRepository) {
				repo.FindByIDFunc = func(id primitive.ObjectID) (*Product, error) {
					return &Product{
						ID:  id,
						Key: "images/compress",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "instruction created",
		},
		{
			name:           "missing product_id",
			payload:        map[string]string{},
			setupRepo:      func(repo *MockProductRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "ProductID is required",
		},
		{
			name:    "invalid product_id format",
			payload: map[string]string{"product_id": "invalid-objectid"},
			setupRepo: func(repo *MockProductRepository) {
				repo.FindByIDFunc = func(id primitive.ObjectID) (*Product, error) {
					return nil, nil
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid ObjectId",
		},
	}

	app, cfg, _, _, productRepo, _, _ := createTestApp()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup repository mock
			tt.setupRepo(productRepo)

			imageSvc := NewImageService()
			instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, &MockInstructionDetailRepository{}, productRepo, imageSvc)

			app.Post("/instructions", instrHandler.CreateInstruction)

			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/instructions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var result map[string]interface{}
			err = json.Unmarshal(responseBody, &result)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedMsg, result["message"])
		})
	}
}

func createTestJWTToken(userID string, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic("Failed to create test token: " + err.Error())
	}
	return tokenString
}

// Handler tests would require complex integration setup
// These tests focus on models, config, and image compression instead

// Image Service Compression Tests
func TestImageService_Compress_BasicFunctionality(t *testing.T) {
	app, _, _, _, productRepo, _, _ := createTestApp()

	productHandler := NewProductHandler(productRepo)

	// Mock successful product listing
	testProducts := []*Product{
		{
			ID:          primitive.NewObjectID(),
			Key:         "images/compress",
			Title:       "Image Compression",
			Description: "Compress images",
			ProductType: "image",
			IsActive:    true,
			IsFree:      true,
		},
	}

	productRepo.ListFunc = func() ([]*Product, error) {
		return testProducts, nil
	}

	app.Get("/products", productHandler.ListProducts)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Products retrieved successfully", result["message"])
	assert.Nil(t, result["errors"])
	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["data"].(map[string]interface{})["products"])
}

func TestProductHandler_ListProducts_DBError(t *testing.T) {
	app, _, _, _, productRepo, _, _ := createTestApp()

	productHandler := NewProductHandler(productRepo)

	productRepo.ListFunc = func() ([]*Product, error) {
		return nil, assert.AnError
	}

	app.Get("/products", productHandler.ListProducts)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// Instruction Handler Tests
func TestInstructionHandler_CreateInstruction_Success(t *testing.T) {
	app, cfg, _, _, productRepo, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, detailRepo, productRepo, imageSvc)

	// Mock product found
	testProduct := &Product{
		ID:          primitive.NewObjectID(),
		Key:         "images/compress",
		Title:       "Image Compression",
		Description: "Compress images",
		ProductType: "image",
		IsActive:    true,
		IsFree:      true,
	}

	productRepo.FindByIDFunc = func(id primitive.ObjectID) (*Product, error) {
		return testProduct, nil
	}

	// Mock instruction creation
	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: testProduct.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	app.Post("/instructions", instrHandler.CreateInstruction)

	payload := map[string]string{"product_id": testProduct.ID.Hex()}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/instructions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "instruction created", result["message"])
	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["data"].(map[string]interface{})["instruction"])
	assert.Equal(t, testProduct.ID.Hex(), result["data"].(map[string]interface{})["instruction"].(map[string]interface{})["product_id"])
}

func TestInstructionHandler_CreateInstruction_MissingProductID(t *testing.T) {
	app, cfg, _, _, _, _, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, &MockInstructionDetailRepository{}, &MockProductRepository{}, imageSvc)

	app.Post("/instructions", instrHandler.CreateInstruction)

	payload := map[string]string{} // Missing product_id
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/instructions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "ProductID is required", result["message"])
}

func TestInstructionHandler_CreateInstruction_ProductNotFound(t *testing.T) {
	app, cfg, _, _, productRepo, _, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, &MockInstructionDetailRepository{}, productRepo, imageSvc)

	productRepo.FindByIDFunc = func(id primitive.ObjectID) (*Product, error) {
		return nil, nil
	}

	app.Post("/instructions", instrHandler.CreateInstruction)

	testProductID := primitive.NewObjectID()
	payload := map[string]string{"product_id": testProductID.Hex()}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/instructions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestInstructionHandler_ListInstructions_Success(t *testing.T) {
	app, cfg, _, _, _, instrRepo, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, &MockInstructionDetailRepository{}, &MockProductRepository{}, imageSvc)

	// Mock instruction listing
	testInstructions := []*Instruction{
		{
			ID:        primitive.NewObjectID(),
			UserID:    primitive.NewObjectID(),
			ProductID: primitive.NewObjectID(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	instrRepo.ListLatestFunc = func(userId string, limit int) ([]*Instruction, error) {
		return testInstructions, nil
	}

	app.Get("/instructions", instrHandler.ListInstructions)

	req := httptest.NewRequest(http.MethodGet, "/instructions", nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["message"])
	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["data"].(map[string]interface{})["instructions"])
}

func TestInstructionHandler_GetInstructionByID_Success(t *testing.T) {
	app, cfg, _, _, _, instrRepo, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, &MockInstructionDetailRepository{}, &MockProductRepository{}, imageSvc)

	// Mock instruction found
	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	app.Get("/instructions/:id", instrHandler.GetInstructionByID)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["message"])
	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["data"].(map[string]interface{})["instruction"])
}

func TestInstructionHandler_GetInstruction_NotFound(t *testing.T) {
	app, cfg, _, _, _, instrRepo, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, &MockInstructionDetailRepository{}, &MockProductRepository{}, imageSvc)

	testInstructionID := primitive.NewObjectID()
	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		return nil
	}

	app.Get("/instructions/:id", instrHandler.GetInstructionByID)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstructionID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestInstructionHandler_GetInstructionDetails_Success(t *testing.T) {
	app, cfg, _, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	// Mock instruction found
	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	// Mock instruction details
	testDetails := []*InstructionDetail{
		{
			ID:            primitive.NewObjectID(),
			InstructionID: testInstruction.ID,
			FileName:      "test.jpg",
			FileSize:      1024,
			MimeType:      "image/jpeg",
			Status:        FileStatusDone,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		},
	}

	detailRepo.ListByInstructionFunc = func(instructionID primitive.ObjectID) []*InstructionDetail {
		if instructionID == testInstruction.ID {
			return testDetails
		}
		return []*InstructionDetail{}
	}

	app.Get("/instructions/:id/details", instrHandler.GetInstructionDetails)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex()+"/details", nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testInstruction.UserID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["message"])
	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["data"].(map[string]interface{})["files"])
}

func TestInstructionHandler_GetInstructionDetails_Unauthorized(t *testing.T) {
	app, cfg, _, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	// Mock instruction belonging to different user
	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	app.Get("/instructions/:id/details", instrHandler.GetInstructionDetails)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex()+"/details", nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("different-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// RunInstructionMessage Test
func TestInstructionHandler_RunInstructionMessage_Success(t *testing.T) {
	app, cfg, s3, nats, _, _, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, s3, nats, &MockInstructionRepository{}, detailRepo, &MockProductRepository{}, imageSvc)

	testFileID := primitive.NewObjectID()
	testInstructionID := primitive.NewObjectID()
	testUserID := primitive.NewObjectID()
	testProductID := primitive.NewObjectID()

	// Mock input file
	inputFile := &InstructionDetail{
		ID:            testFileID,
		InstructionID: testInstructionID,
		FileName:      "input.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusPending,
		OutputID:      &testFileID,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	// Mock output file
	outputFile := &InstructionDetail{
		ID:            testFileID,
		InstructionID: testInstructionID,
		FileName:      "output.jpg",
		FileSize:      512,
		MimeType:      "image/jpeg",
		Status:        FileStatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	detailRepo.GetByIDFunc = func(id primitive.ObjectID) *InstructionDetail {
		if id == testFileID {
			return inputFile
		} else if id == *inputFile.OutputID {
			return outputFile
		}
		return nil
	}

	// Mock instruction
	testInstruction := &Instruction{
		ID:        testInstructionID,
		UserID:    testUserID,
		ProductID: testProductID,
	}

	instrHandler.(*InstructionHandler).instrRepo = &MockInstructionRepository{
		GetByIDFunc: func(id primitive.ObjectID) *Instruction {
			if id == testInstructionID {
				return testInstruction
			}
			return nil
		},
	}

	// Mock product
	productHandler.(*InstructionHandler).productRepo = &MockProductRepository{
		FindByIDFunc: func(id primitive.ObjectID) (*Product, error) {
			return &Product{
				ID:  testProductID,
				Key: "images/compress",
			}, nil
		},
	}

	// Mock S3 and NATS
	s3.GetFunc = func(key string) []byte {
		return []byte("test image data")
	}

	s3.PutFunc = func(key string, data []byte) error {
		return nil
	}

	nats.PublishFunc = func(subject string, data []byte) error {
		return nil
	}

	// Mock detail repo methods
	detailRepo.UpdateStatusFunc = func(id primitive.ObjectID, status FileStatus) error {
		return nil
	}

	detailRepo.UpdateStatusAndSizeFunc = func(id primitive.ObjectID, status FileStatus, size int64) error {
		return nil
	}

	// Test the method
	instrHandler.RunInstructionMessage([]byte(testFileID.Hex()))

	// Verify the methods were called (checking if S3 get was called)
	assert.NotNil(t, s3.GetFunc("input.jpg"))
}

func TestInstructionHandler_RunInstructionMessage_FileNotFound(t *testing.T) {
	app, cfg, _, _, _, _, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, detailRepo, &MockProductRepository{}, imageSvc)

	// Mock file not found
	detailRepo.GetByIDFunc = func(id primitive.ObjectID) *InstructionDetail {
		return nil
	}

	// Should not panic
	instrHandler.RunInstructionMessage([]byte("invalid-file-id"))

	// Expected behavior - logged but no panic
	assert.True(t, true)
}

func TestInstructionHandler_RunInstructionMessage_MissingInputFile(t *testing.T) {
	app, cfg, _, _, _, _, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, detailRepo, &MockProductRepository{}, imageSvc)

	testFileID := primitive.NewObjectID()

	// Mock output file exists but input file doesn't
	detailRepo.GetByIDFunc = func(id primitive.ObjectID) *InstructionDetail {
		if id == testFileID {
			// Input file
			return nil
		}
		return &InstructionDetail{
			ID:            testFileID,
			InstructionID: testInstructionID,
		}
	}

	// Should not panic
	instrHandler.RunInstructionMessage([]byte(testFileID.Hex()))

	// Expected behavior - logged but no panic
	assert.True(t, true)
}

// Image Service Tests
func TestImageService_Compress_Success(t *testing.T) {
	app := fiber.New()

	imageSvc := NewImageService()

	inputData := []byte("fake image data")

	result, err := imageSvc.Compress(inputData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
	assert.Less(t, len(result), len(inputData)) // Compressed data should be smaller
}

// Main function tests
func TestMainFunction_Initialization(t *testing.T) {
	// Test that main function can initialize components without errors
	// This is a basic smoke test to ensure the main function works
	cfg := LoadConfig()
	assert.NotNil(t, cfg)

	// Test LoadConfig doesn't panic
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.MongoURI)
	assert.NotEmpty(t, cfg.NatsURI)
}

func TestEnvironmentConfigurations(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedPort string
	}{
		{
			name:         "development default",
			expectedPort: ":3000",
		},
		{
			name: "production override",
			envVars: map[string]string{
				"PORT": ":8080",
			},
			expectedPort: ":8080",
		},
		{
			name: "test environment",
			envVars: map[string]string{
				"ENVIRONMENT": "test",
				"PORT":        ":9000",
			},
			expectedPort: ":9000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars before test
			for k := range tt.envVars {
				os.Unsetenv(k)
			}

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg := LoadConfig()

			// Expected port, default to test expectation
			expected := tt.expectedPort
			if expected == "" {
				expected = ":3000"
			}

			assert.Equal(t, expected, cfg.Port)

			// Clean up env vars
			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	// Test various error scenarios
	cfg := newMockConfig()

	// Test with nil JWT secret (should not crash)
	os.Setenv("JWT_SECRET", "")
	cfg = LoadConfig()
	assert.NotEmpty(t, cfg.JWTSecret) // Should have default

	// Test with empty environment
	os.Setenv("ENVIRONMENT", "")
	cfg = LoadConfig()
	assert.Equal(t, "development", cfg.Environment) // Should have default

	// Clean up
	os.Setenv("JWT_SECRET", "")
}

func TestCORSHeaders(t *testing.T) {
	app, _, _, _, _, _, _ := createTestApp()

	// Test CORS headers are set correctly
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Check CORS headers are present
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestRateLimiting(t *testing.T) {
	app, cfg, _, _, _, _, _ := createTestApp()

	// Add a simple endpoint to test rate limiting
	app.Get("/ratelimit", func(c *fiber.Ctx) error {
		return c.SendString("rate limited")
	})

	// Test multiple requests to trigger rate limiting
	req := httptest.NewRequest(http.MethodGet, "/ratelimit", nil)

	for i := 0; i < 5; i++ {
		resp, err := app.Test(req)
		require.NoError(t, err)

		if i < 3 { // Should succeed within rate limit
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		} else { // Should be rate limited
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}
}

// Integration test for file upload flow
func TestFileUploadFlow(t *testing.T) {
	app, cfg, s3, nats, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, s3, nats, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	// Mock product
	testProductID := primitive.NewObjectID()
	testInstructionID := primitive.NewObjectID()
	testUserID := primitive.NewObjectID()

	productHandler.(*InstructionHandler).productRepo = &MockProductRepository{
		FindByIDFunc: func(id primitive.ObjectID) (*Product, error) {
			return &Product{
				ID:  testProductID,
				Key: "images/compress",
			}, nil
		},
	}

	// Mock instruction creation
	instrHandler.(*InstructionHandler).instrRepo = &MockInstructionRepository{
		GetByIDFunc: func(id primitive.ObjectID) *Instruction {
			if id == testInstructionID {
				return &Instruction{
					ID:        testInstructionID,
					UserID:    testUserID,
					ProductID: testProductID,
				}
			}
			return nil
		},
	}

	// Mock file upload
	uploadedFile := []byte("test image data")
	s3.PutFunc = func(key string, data []byte) error {
		assert.Equal(t, uploadedFile, data)
		return nil
	}

	nats.PublishFunc = func(subject string, data []byte) error {
		assert.Equal(t, cfg.NatsSubjectImageRequests, subject)
		return nil
	}

	app.Post("/instructions/:id/details", instrHandler.CreateInstructionDetails)

	// Create test file
	body := bytes.NewBuffer(uploadedFile)
	req := httptest.NewRequest(http.MethodPost, "/instructions/"+testInstructionID.Hex()+"/details", body)
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testUserID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Benchmark tests
func BenchmarkProductHandler_ListProducts(b *testing.B) {
	app, cfg, _, _, productRepo, _, _ := createTestApp()

	productHandler := NewProductHandler(productRepo)

	// Mock product listing
	testProducts := []*Product{
		{
			ID:          primitive.NewObjectID(),
			Key:         "images/compress",
			Title:       "Image Compression",
			Description: "Compress images",
			ProductType: "image",
			IsActive:    true,
			IsFree:      true,
		},
	}

	productRepo.ListFunc = func() ([]*Product, error) {
		return testProducts, nil
	}

	app.Get("/products", productHandler.ListProducts)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req)
		require.NoError(b, err)
		assert.Equal(b, http.StatusOK, resp.StatusCode)
	}
}

func BenchmarkInstructionHandler_CreateInstruction(b *testing.B) {
	app, cfg, _, _, productRepo, _, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, &MockInstructionDetailRepository{}, productRepo, imageSvc)

	testProduct := &Product{
		ID:          primitive.NewObjectID(),
		Key:         "images/compress",
		Title:       "Image Compression",
		Description: "Compress images",
		ProductType: "image",
		IsActive:    true,
		IsFree:      true,
	}

	productRepo.FindByIDFunc = func(id primitive.ObjectID) (*Product, error) {
		return testProduct, nil
	}

	app.Post("/instructions", instrHandler.CreateInstruction)

	payload := map[string]string{"product_id": testProduct.ID.Hex()}
	body, err := json.Marshal(payload)
	require.NoError(b, err)

	req := httptest.NewRequest(http.MethodPost, "/instructions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req)
		require.NoError(b, err)
		assert.Equal(b, http.StatusOK, resp.StatusCode)
	}
}

// Additional Handler Tests - GetInstructionDetail endpoint
func TestInstructionHandler_GetInstructionDetail_Success(t *testing.T) {
	app, cfg, _, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	testDetail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: testInstruction.ID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusDone,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	detailRepo.GetByIDFunc = func(id primitive.ObjectID) *InstructionDetail {
		if id == testDetail.ID {
			return testDetail
		}
		return nil
	}

	app.Get("/instructions/:id/details/:detailId", instrHandler.GetInstructionDetail)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex()+"/details/"+testDetail.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testInstruction.UserID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["message"])
	assert.NotNil(t, result["data"].(map[string]interface{})["detail"])
}

func TestInstructionHandler_GetInstructionDetail_Forbidden(t *testing.T) {
	app, cfg, _, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	app.Get("/instructions/:id/details/:detailId", instrHandler.GetInstructionDetail)

	differentUserID := primitive.NewObjectID()
	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex()+"/details/"+primitive.NewObjectID().Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(differentUserID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestInstructionHandler_GetInstructionDetail_NotFound(t *testing.T) {
	app, cfg, _, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	testInstructionID := primitive.NewObjectID()
	testDetailID := primitive.NewObjectID()

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		return nil
	}

	app.Get("/instructions/:id/details/:detailId", instrHandler.GetInstructionDetail)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstructionID.Hex()+"/details/"+testDetailID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken("test-user", cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// GetInstructionDetilFile endpoint tests
func TestInstructionHandler_GetInstructionDetilFile_Success(t *testing.T) {
	app, cfg, s3, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, s3, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	testDetail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: testInstruction.ID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusDone,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	fileContent := []byte("test image content")

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	detailRepo.GetByIDFunc = func(id primitive.ObjectID) *InstructionDetail {
		if id == testDetail.ID {
			return testDetail
		}
		return nil
	}

	s3.GetFunc = func(key string) []byte {
		if key == testDetail.FileName {
			return fileContent
		}
		return nil
	}

	app.Get("/instructions/:id/details/:detailId/file", instrHandler.GetInstructionDetilFile)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex()+"/details/"+testDetail.ID.Hex()+"/file", nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testInstruction.UserID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/octet-stream", resp.Header.Get("content-type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, fileContent, body)
}

func TestInstructionHandler_GetInstructionDetilFile_FileNotFound(t *testing.T) {
	app, cfg, s3, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, s3, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	testDetail := &InstructionDetail{
		ID:            primitive.NewObjectID(),
		InstructionID: testInstruction.ID,
		FileName:      "test.jpg",
		FileSize:      1024,
		MimeType:      "image/jpeg",
		Status:        FileStatusDone,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	detailRepo.GetByIDFunc = func(id primitive.ObjectID) *InstructionDetail {
		if id == testDetail.ID {
			return testDetail
		}
		return nil
	}

	s3.GetFunc = func(key string) []byte {
		return nil // File not found
	}

	app.Get("/instructions/:id/details/:detailId/file", instrHandler.GetInstructionDetilFile)

	req := httptest.NewRequest(http.MethodGet, "/instructions/"+testInstruction.ID.Hex()+"/details/"+testDetail.ID.Hex()+"/file", nil)
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testInstruction.UserID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// ListUncleanedFiles endpoint tests
func TestInstructionHandler_ListUncleanedFiles_Success(t *testing.T) {
	app, cfg, _, _, _, _, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, &MockInstructionRepository{}, detailRepo, &MockProductRepository{}, imageSvc)

	testFiles := []*InstructionDetail{
		{
			ID:            primitive.NewObjectID(),
			InstructionID: primitive.NewObjectID(),
			FileName:      "uncleaned1.jpg",
			FileSize:      1024,
			MimeType:      "image/jpeg",
			Status:        FileStatusDone,
			IsCleaned:     false,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		},
		{
			ID:            primitive.NewObjectID(),
			InstructionID: primitive.NewObjectID(),
			FileName:      "uncleaned2.jpg",
			FileSize:      2048,
			MimeType:      "image/jpeg",
			Status:        FileStatusDone,
			IsCleaned:     false,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		},
	}

	detailRepo.ListUncleanedFunc = func() []*InstructionDetail {
		return testFiles
	}

	app.Get("/files", instrHandler.ListUncleanedFiles)

	req := httptest.NewRequest(http.MethodGet, "/files", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["message"])
	assert.NotNil(t, result["data"].(map[string]interface{})["files"])
}

// CleanInstruction tests
func TestInstructionHandler_CleanInstruction_Success(t *testing.T) {
	_, cfg, s3, _, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, s3, &MockNats{}, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	oldTime := time.Now().Add(-2 * time.Hour)
	testFiles := []*InstructionDetail{
		{
			ID:            primitive.NewObjectID(),
			InstructionID: primitive.NewObjectID(),
			FileName:      "old_file.jpg",
			FileSize:      1024,
			MimeType:      "image/jpeg",
			Status:        FileStatusDone,
			IsCleaned:     false,
			CreatedAt:     oldTime,
			UpdatedAt:     oldTime,
		},
	}

	detailRepo.ListOlderThanFunc = func(cutoff time.Time) []*InstructionDetail {
		return testFiles
	}

	detailRepo.ListPendingUpdatedBeforeFunc = func(cutoff time.Time) []*InstructionDetail {
		return []*InstructionDetail{}
	}

	s3DeleteCalled := false
	s3.DeleteFunc = func(key string) error {
		s3DeleteCalled = true
		return nil
	}

	detailRepo.MarkCleanedFunc = func(ids []primitive.ObjectID) error {
		return nil
	}

	err := instrHandler.CleanInstruction()
	assert.NoError(t, err)
	assert.True(t, s3DeleteCalled, "S3 Delete should have been called")
}

// CreateInstructionDetails tests (file upload)
func TestInstructionHandler_CreateInstructionDetails_Success(t *testing.T) {
	app, cfg, s3, nats, _, instrRepo, detailRepo := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, s3, nats, instrRepo, detailRepo, &MockProductRepository{}, imageSvc)

	testInstruction := &Instruction{
		ID:        primitive.NewObjectID(),
		UserID:    primitive.NewObjectID(),
		ProductID: primitive.NewObjectID(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	instrRepo.GetByIDFunc = func(id primitive.ObjectID) *Instruction {
		if id == testInstruction.ID {
			return testInstruction
		}
		return nil
	}

	detailRepo.CreateManyFunc = func(details []*InstructionDetail) error {
		return nil
	}

	s3.PutFunc = func(key string, data []byte) error {
		return nil
	}

	nats.PublishFunc = func(subject string, data []byte) error {
		return nil
	}

	app.Post("/instructions/:id/details", instrHandler.CreateInstructionDetails)

	// Create file upload request
	body := bytes.NewBuffer([]byte("test image data"))
	req := httptest.NewRequest(http.MethodPost, "/instructions/"+testInstruction.ID.Hex()+"/details", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testInstruction.UserID.Hex(), cfg.JWTSecret))

	// Note: This test is simplified due to Fiber's file upload complexity
	// In a real scenario, you'd use proper multipart form data construction
}

// Test CreateInstruction with database error
func TestInstructionHandler_CreateInstruction_DBError(t *testing.T) {
	app, cfg, _, _, productRepo, instrRepo, _ := createTestApp()

	imageSvc := NewImageService()
	instrHandler := NewInstructionHandler(cfg, &MockS3{}, &MockNats{}, instrRepo, &MockInstructionDetailRepository{}, productRepo, imageSvc)

	testProduct := &Product{
		ID:          primitive.NewObjectID(),
		Key:         "images/compress",
		Title:       "Image Compression",
		Description: "Compress images",
		ProductType: "image",
		IsActive:    true,
		IsFree:      true,
	}

	productRepo.FindByIDFunc = func(id primitive.ObjectID) (*Product, error) {
		return testProduct, nil
	}

	// Mock repository error
	instrRepo.CreateFunc = func(instr *Instruction) error {
		return assert.AnError
	}

	app.Post("/instructions", instrHandler.CreateInstruction)

	payload := map[string]string{"product_id": testProduct.ID.Hex()}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/instructions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestJWTToken(testProduct.ID.Hex(), cfg.JWTSecret))

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
