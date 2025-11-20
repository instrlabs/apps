package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/instrlabs/shared/modelx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock structs for PDF service testing
type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucket, objectName, reader, size, opts)
	return minio.UploadInfo{}, args.Error(0)
}

func (m *MockMinioClient) GetObject(ctx context.Context, bucket, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, bucket, objectName, opts)
	return args.Get(0).(*minio.Object), args.Error(1)
}

func (m *MockMinioClient) RemoveObject(ctx context.Context, bucket, objectName string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, bucket, objectName, opts)
	return args.Error(0)
}

type MockNATSConn struct {
	mock.Mock
}

func (m *MockNATSConn) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

func (m *MockNATSConn) Close() error {
	return nil
}

type MockInstructionRepository struct {
	mock.Mock
}

func (m *MockInstructionRepository) Create(instruction *modelx.Instruction) error {
	args := m.Called(instruction)
	return args.Error(0)
}

func (m *MockInstructionRepository) GetByID(id primitive.ObjectID, instruction modelx.Instruction) error {
	args := m.Called(id, instruction)
	return args.Error(0)
}

func (m *MockInstructionRepository) UpdateByID(id primitive.ObjectID, update bson.M) error {
	args := m.Called(id, update)
	return args.Error(0)
}

func (m *MockInstructionRepository) DeleteByID(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockInstructionDetailRepository struct {
	mock.Mock
}

func (m *MockInstructionDetailRepository) Create(detail *modelx.InstructionFile) error {
	args := m.Called(detail)
	return args.Error(0)
}

func (m *MockInstructionDetailRepository) GetByID(id primitive.ObjectID) (*modelx.InstructionFile, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*modelx.InstructionFile), args.Error(1)
}

func (m *MockInstructionDetailRepository) UpdateByID(id primitive.ObjectID, update bson.M) error {
	args := m.Called(id, update)
	return args.Error(0)
}

func (m *MockInstructionDetailRepository) DeleteByID(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockProductClient struct {
	mock.Mock
}

func (m *MockProductClient) FindByID(id primitive.ObjectID, service string) (*modelx.Product, error) {
	args := m.Called(id, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*modelx.Product), args.Error(1)
}

type MockPDFService struct {
	mock.Mock
}

func (m *MockPDFService) ProcessPDF(inputData []byte, fileName string) ([]byte, error) {
	args := m.Called(inputData, fileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func setupTestApp() (*fiber.App, *InstructionHandler, *Config) {
	app := fiber.New()

	// Create test config
	cfg := &Config{
		ServiceName:            "pdf-service-test",
		Port:                   "3003",
		Environment:            "test",
		S3Bucket:               "test-bucket",
		S3Endpoint:             "localhost:9000",
		S3AccessKey:            "test-access-key",
		S3SecretKey:            "test-secret-key",
		S3Region:               "us-east-1",
		NatsSubjectPdfRequests: "pdf.process.requests",
	}

	return app, NewInstructionHandler(cfg, &MockMinioClient{}, &MockNATSConn{}, &MockInstructionRepository{}, &MockInstructionDetailRepository{}, &MockProductClient{}, &MockPDFService{}), cfg
}

func createTestInstruction(id primitive.ObjectID, userID *primitive.ObjectID, guestID *string, productID primitive.ObjectID) *modelx.Instruction {
	return &modelx.Instruction{
		ID:        id,
		UserID:    userID,
		GuestID:   guestID,
		ProductID: productID,
		Status:    modelx.InstructionStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func createTestInstructionFile(id primitive.ObjectID, instructionID primitive.ObjectID, status modelx.InstructionDetailStatus, key string) *modelx.InstructionFile {
	return &modelx.InstructionFile{
		ID:            id,
		InstructionID: instructionID,
		Type:          "input",
		FileName:      "test.pdf",
		FileKey:       key,
		MimeType:      "application/pdf",
		Status:        status,
		Size:          1024,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
}

func createTestProduct(id primitive.ObjectID) *modelx.Product {
	return &modelx.Product{
		ID:        id,
		Name:      "Test Product",
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func createTestPDF(t *testing.T, filename string, content []byte) *multipart.Writer {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write(content)

	writer.Close()
	return writer
}

func createTestContextWithHeaders(app *fiber.App, method, path string, headers map[string]string, body io.Reader) *fiber.Ctx {
	req := httptest.NewRequest(method, path, body)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")

	return app.Test(req, -1)
}

// CreateInstruction Tests
func TestCreateInstruction_Success(t *testing.T) {
	app, handler, _ := setupTestApp()

	productID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	instructionID := primitive.NewObjectID()

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("Create", mock.AnythingOfType("*modelx.Instruction")).Return(nil)

	productClient := handler.productClient.(*MockProductClient).(*MockProductClient)
	productClient.On("FindByID", productID, "pdf").Return(createTestProduct(productID), nil)

	t.Run("Valid instruction creation", func(t *testing.T) {
		ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions", map[string]string{}, bytes.NewBufferString(`{"product_id": "`+productID.Hex()+`"}`))
		ctx.Locals("headers", map[string]string{"x-user-id": userID.Hex()})

		err := handler.CreateInstruction(ctx)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

		var response map[string]interface{}
		json.Unmarshal(ctx.Response().Body(), &response)
		assert.Equal(t, "instruction creation successful", response["message"])
		assert.NotNil(t, response["data"])

		instrRepo.AssertExpectations(t)
		productClient.AssertExpectations(t)
	})
}

func TestCreateInstruction_InvalidBody(t *testing.T) {
	app, handler, _ := setupTestApp()

	t.Run("Invalid JSON", func(t *testing.T) {
		ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions", map[string]string{}, strings.NewReader("invalid json"))
		err := handler.CreateInstruction(ctx)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
	})

	t.Run("Missing product ID", func(t *testing.T) {
		ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions", map[string]string{}, bytes.NewBufferString(`{}`))
		err := handler.CreateInstruction(ctx)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
	})
}

func TestCreateInstruction_ProductNotFound(t *testing.T) {
	app, handler, _ := setupTestApp()

	productID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	productClient := handler.productClient.(*MockProductClient).(*MockProductClient)
	productClient.On("FindByID", productID, "pdf").Return(nil, errors.New("product not found"))

	ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions", map[string]string{"x-user-id": userID.Hex()}, bytes.NewBufferString(`{"product_id": "`+productID.Hex()+`"}`))

	err := handler.CreateInstruction(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "product not found", response["message"])

	productClient.AssertExpectations(t)
}

func TestCreateInstruction_MissingHeaders(t *testing.T) {
	app, handler, _ := setupTestApp()

	productID := primitive.NewObjectID()
	productClient := handler.productClient.(*MockProductClient).(*MockProductClient)
	productClient.On("FindByID", productID, "pdf").Return(createTestProduct(productID), nil)

	ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions", map[string]string{}, bytes.NewBufferString(`{"product_id": "`+productID.Hex()+`"}`))

	err := handler.CreateInstruction(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "missing user identification - provide either x-user-id or x-guest-id header", response["message"])

	productClient.AssertExpectations(t)
}

// CreateInstructionDetails Tests
func TestCreateInstructionDetails_Success(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	detailID := primitive.NewObjectID()

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("Create", mock.AnythingOfType("*modelx.InstructionFile")).Return(nil)

	s3Client := handler.s3.(*MockMinioClient).(*MockMinioClient)
	s3Client.On("PutObject", mock.Anything, "test-bucket", mock.AnythingOfType("string"), mock.Anything, mock.Anything, minio.PutObjectOptions{}).Return(minio.UploadInfo{}, nil)

	pdfSvc := handler.pdfSvc.(*MockPDFService).(*MockPDFService)
	pdfSvc.On("ProcessPDF", mock.AnythingOfType("[]byte"), "test.pdf").Return([]byte("processed-pdf"), nil)

	natsConn := handler.nats.(*MockNATSConn).(*MockNATSConn)
	natsConn.On("Publish", "pdf.process.requests", mock.AnythingOfType("[]byte")).Return(nil)

	t.Run("Valid PDF upload", func(t *testing.T) {
		testPDF := createTestPDF(t, "test.pdf", []byte("fake PDF content"))
		ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

		err := handler.CreateInstructionDetails(ctx)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

		var response map[string]interface{}
		json.Unmarshal(ctx.Response().Body(), &response)
		assert.Equal(t, "instruction details uploaded successfully", response["message"])
		assert.NotNil(t, response["data"])

		instrRepo.AssertExpectations(t)
		detailRepo.AssertExpectations(t)
		s3Client.AssertExpectations(t)
		pdfSvc.AssertExpectations(t)
		natsConn.AssertExpectations(t)
	})
}

func TestCreateInstructionDetails_InvalidInstructionID(t *testing.T) {
	app, handler, _ := setupTestApp()

	ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions/invalid-id/details", map[string]string{"x-user-id": "507f1f77bcf86cd799439011"}, strings.NewReader(""))
	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
}

func TestCreateInstructionDetails_InstructionNotFound(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(errors.New("instruction not found"))

	ctx := createTestContextWithHeaders(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "instruction not found", response["message"])

	instrRepo.AssertExpectations(t)
}

func TestCreateInstructionDetails_Unauthorized(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID1 := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	userID2 := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")

	instruction := createTestInstruction(instructionID, &userID2, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	ctx := createTestContextWithHeaders(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID1.Hex()}, strings.NewReader(""))

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "user not authorized to create instruction details", response["message"])

	instrRepo.AssertExpectations(t)
}

func TestCreateInstructionDetails_MissingFile(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	ctx := createTestContextWithHeaders(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "failed to read uploaded file", response["message"])

	instrRepo.AssertExpectations(t)
}

func TestCreateInstructionDetails_FileTooLarge(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	// Create a PDF that's larger than 50MB
	largePDF := make([]byte, 51*1024*1024)
	testPDF := createTestPDF(t, "large.pdf", largePDF)

	ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "file size exceeds maximum limit of 50MB", response["message"])

	instrRepo.AssertExpectations(t)
}

func TestCreateInstructionDetails_InvalidFileType(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	testPDF := createTestPDF(t, "test.txt", []byte("text content"))

	// Override the Content-Type to simulate an invalid file type
	testPDF.SetBoundary("invalid-boundary")

	ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "invalid file type", response["message"])

	instrRepo.AssertExpectations(t)
}

func TestCreateInstructionDetails_S3UploadFailure(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("Create", mock.AnythingOfType("*modelx.InstructionFile")).Return(nil)

	s3Client := handler.s3.(*MockMinioClient).(*MockMinioClient)
	s3Client.On("PutObject", mock.Anything, "test-bucket", mock.AnythingOfType("string"), mock.Anything, mock.Anything, minio.PutObjectOptions{}).Return(minio.UploadInfo{}, errors.New("S3 upload failed"))

	pdfSvc := handler.pdfSvc.(*MockPDFService).(*MockPDFService)
	pdfSvc.On("ProcessPDF", mock.AnythingOfType("[]byte"), "test.pdf").Return([]byte("processed-pdf"), nil)

	testPDF := createTestPDF(t, "test.pdf", []byte("fake PDF content"))

	ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "failed to upload file to S3", response["message"])

	instrRepo.AssertExpectations(t)
	detailRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
	pdfSvc.AssertExpectations(t)
}

func TestCreateInstructionDetails_NATSFailure(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("Create", mock.AnythingOfType("*modelx.InstructionFile")).Return(nil)

	s3Client := handler.s3.(*MockMinioClient).(*MockMinioClient)
	s3Client.On("PutObject", mock.Anything, "test-bucket", mock.AnythingOfType("string"), mock.Anything, mock.Anything, minio.PutObjectOptions{}).Return(minio.UploadInfo{}, nil)

	pdfSvc := handler.pdfSvc.(*MockPDFService).(*MockPDFService)
	pdfSvc.On("ProcessPDF", mock.AnythingOfType("[]byte"), "test.pdf").Return([]byte("processed-pdf"), nil)

	natsConn := handler.nats.(*MockNATSConn).(*MockNATSConn)
	natsConn.On("Publish", "pdf.process.requests", mock.AnythingOfType("[]byte")).Return(errors.New("NATS publish failed"))

	testPDF := createTestPDF(t, "test.pdf", []byte("fake PDF content"))

	ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "failed to send message to NATS", response["message"])

	instrRepo.AssertExpectations(t)
	detailRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
	pdfSvc.AssertExpectations(t)
	natsConn.AssertExpectations(t)
}

// GetInstructionByID Tests
func TestGetInstructionByID_Success(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	ctx := createTestContextWithHeaders(app, "GET", fmt.Sprintf("/api/v1/instructions/%s", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.GetInstructionByID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "instruction retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])

	instrRepo.AssertExpectations(t)
}

func TestGetInstructionByID_NotFound(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(errors.New("instruction not found"))

	ctx := createTestContextWithHeaders(app, "GET", fmt.Sprintf("/api/v1/instructions/%s", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.GetInstructionByID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "instruction not found", response["message"])

	instrRepo.AssertExpectations(t)
}

// GetInstructionDetail Tests
func TestGetInstructionDetail_Success(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	detailID := primitive.NewObjectID()

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())
	detail := createTestInstructionFile(detailID, instructionID, modelx.InstructionStatusSuccess, "test-key")

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("GetByID", detailID).Return(detail, nil)

	ctx := createTestContextWithHeaders(app, "GET", fmt.Sprintf("/api/v1/instructions/%s/details/%s", instructionID.Hex(), detailID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.GetInstructionDetail(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "instruction detail retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])

	instrRepo.AssertExpectations(t)
	detailRepo.AssertExpectations(t)
}

func TestGetInstructionDetail_Unauthorized(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID1 := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	userID2 := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	detailID := primitive.NewObjectID()

	instruction := createTestInstruction(instructionID, &userID2, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	ctx := createTestContextWithHeaders(app, "GET", fmt.Sprintf("/api/v1/instructions/%s/details/%s", instructionID.Hex(), detailID.Hex()), map[string]string{"x-user-id": userID1.Hex()}, strings.NewReader(""))

	err := handler.GetInstructionDetail(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "user not authorized to access instruction detail", response["message"])

	instrRepo.AssertExpectations(t)
}

// GetInstructionDetailFile Tests
func TestGetInstructionDetailFile_Success(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	detailID := primitive.NewObjectID()
	detail := createTestInstructionFile(detailID, instructionID, modelx.InstructionStatusSuccess, "test-key")

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("GetByID", detailID).Return(detail, nil)

	s3Client := handler.s3.(*MockMinioClient).(*MockMinioClient)
	s3Client.On("GetObject", mock.Anything, "test-bucket", "test-key", minio.GetObjectOptions{}).Return(&minio.Object{
		ObjectInfo: minio.ObjectInfo{Key: "test-key"},
	}, nil)

	ctx := createTestContextWithHeaders(app, "GET", fmt.Sprintf("/api/v1/instructions/%s/details/%s/file", instructionID.Hex(), detailID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.GetInstructionDetailFile(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

	detailRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
}

func TestGetInstructionDetailFile_NotFound(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	detailID := primitive.NewObjectID()

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("GetByID", detailID).Return(nil, errors.New("detail not found"))

	ctx := createTestContextWithHeaders(app, "GET", fmt.Sprintf("/api/v1/instructions/%s/details/%s/file", instructionID.Hex(), detailID.Hex()), map[string]string{"x-user-id": userID.Hex()}, strings.NewReader(""))

	err := handler.GetInstructionDetailFile(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "instruction detail not found", response["message"])

	detailRepo.AssertExpectations(t)
}

// Helper functions
func createTestContextWithFile(app *fiber.App, method, path string, headers map[string]string, fileData *multipart.Writer) *fiber.Ctx {
	body := &bytes.Buffer{}
	body.Write(fileData.FormDataBoundary())
	body.Write(fileData.Close())
	body.Write(fileData.Buffer().Bytes())

	req := httptest.NewRequest(method, path, body)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", fileData.FormDataContentType())

	return app.Test(req, -1)
}

// Mock utility functions
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GetMimeTypeFromName(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// PDF processing specific tests
func TestCreateInstructionDetails_CorruptedPDF(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("Create", mock.AnythingOfType("*modelx.InstructionFile")).Return(nil)

	s3Client := handler.s3.(*MockMinioClient).(*MockMinioClient)
	s3Client.On("PutObject", mock.Anything, "test-bucket", mock.AnythingOfType("string"), mock.Anything, mock.Anything, minio.PutObjectOptions{}).Return(minio.UploadInfo{}, nil)

	pdfSvc := handler.pdfSvc.(*MockPDFService).(*MockPDFService)
	pdfSvc.On("ProcessPDF", mock.AnythingOfType("[]byte"), "test.pdf").Return(nil, errors.New("PDF corrupted"))

	testPDF := createTestPDF(t, "test.pdf", []byte("corrupted PDF content"))

	ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "failed to process PDF", response["message"])

	instrRepo.AssertExpectations(t)
	detailRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
	pdfSvc.AssertExpectations(t)
}

func TestCreateInstructionDetails_PasswordProtectedPDF(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	detailRepo := handler.detailRepo.(*MockInstructionDetailRepository).(*MockInstructionDetailRepository)
	detailRepo.On("Create", mock.AnythingOfType("*modelx.InstructionFile")).Return(nil)

	s3Client := handler.s3.(*MockMinioClient).(*MockMinioClient)
	s3Client.On("PutObject", mock.Anything, "test-bucket", mock.AnythingOfType("string"), mock.Anything, mock.Anything, minio.PutObjectOptions{}).Return(minio.UploadInfo{}, nil)

	pdfSvc := handler.pdfSvc.(*MockPDFService).(*MockPDFService)
	pdfSvc.On("ProcessPDF", mock.AnythingOfType("[]byte"), "protected.pdf").Return(nil, errors.New("PDF is password protected"))

	testPDF := createTestPDF(t, "protected.pdf", []byte("password protected PDF content"))

	ctx := createTestContextWithFile(app, "POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), map[string]string{"x-user-id": userID.Hex()}, testPDF)

	err := handler.CreateInstructionDetails(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, ctx.Response().StatusCode())

	var response map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &response)
	assert.Equal(t, "failed to process PDF", response["message"])

	instrRepo.AssertExpectations(t)
	detailRepo.AssertExpectations(t)
	s3Client.AssertExpectations(t)
	pdfSvc.AssertExpectations(t)
}

// PDF-specific validation tests
func TestCreateInstructionDetails_PDFSpecificValidation(t *testing.T) {
	app, handler, _ := setupTestApp()

	instructionID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")

	instruction := createTestInstruction(instructionID, &userID, nil, primitive.NewObjectID())

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("GetByID", instructionID, mock.AnythingOfType("modelx.Instruction")).Return(nil)

	// Test PDF-specific allowed types
	testCases := []struct {
		name        string
		filename    string
		contentType string
		shouldPass  bool
	}{
		{"Valid PDF", "document.pdf", "application/pdf", true},
		{"Valid PDF with uppercase extension", "DOCUMENT.PDF", "application/pdf", true},
		{"PDF with spaces", "my document.pdf", "application/pdf", true},
		{"Invalid PDF content type", "test.pdf", "text/plain", false},
		{"Image pretending to be PDF", "fake.pdf", "image/jpeg", false},
		{"Word document", "document.docx", "application/pdf", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testPDF := createTestPDF(t, tc.filename, []byte("test content"))

			// Create context with custom content type to simulate the actual upload
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", tc.filename)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			part.Write([]byte("test content"))
			writer.Close()

			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/instructions/%s/details", instructionID.Hex()), body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("x-user-id", userID.Hex())
			ctx := app.Test(req, -1)

			if tc.shouldPass {
				assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())
			} else {
				assert.Equal(t, fiber.StatusBadRequest, ctx.Response().StatusCode())
			}
		})
	}
}

// Ensure PDF service specific configuration is being used
func TestProductClientIntegration_PDFService(t *testing.T) {
	app, handler, _ := setupTestApp()

	productID := primitive.NewObjectID()
	userID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	instructionID := primitive.NewObjectID()

	instrRepo := handler.instrRepo.(*MockInstructionRepository).(*MockInstructionRepository)
	instrRepo.On("Create", mock.AnythingOfType("*modelx.Instruction")).Return(nil)

	productClient := handler.productClient.(*MockProductClient).(*MockProductClient)
	productClient.On("FindByID", productID, "pdf").Return(createTestProduct(productID), nil)

	ctx := createTestContextWithHeaders(app, "POST", "/api/v1/instructions", map[string]string{"x-user-id": userID.Hex()}, bytes.NewBufferString(`{"product_id": "`+productID.Hex()+`"}`))

	err := handler.CreateInstruction(ctx)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, ctx.Response().StatusCode())

	productClient.AssertExpectations(t)
	// Verify that "pdf" was passed as the service type to FindByID
}
