package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// PaymentHandler handles HTTP requests for payments
type PaymentHandler struct {
	midtransService *MidtransService
	paymentRepo     *PaymentRepository
	natsService     *NatsService
	cfg             *Config
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(midtransService *MidtransService, paymentRepo *PaymentRepository, natsService *NatsService, cfg *Config) *PaymentHandler {
	return &PaymentHandler{
		midtransService: midtransService,
		paymentRepo:     paymentRepo,
		natsService:     natsService,
		cfg:             cfg,
	}
}

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	OrderID       string  `json:"orderId"`
	UserID        string  `json:"userId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"paymentMethod,omitempty"`
	Description   string  `json:"description,omitempty"`
	CustomerName  string  `json:"customerName,omitempty"`
	CustomerEmail string  `json:"customerEmail,omitempty"`
	CallbackURL   string  `json:"callbackUrl,omitempty"`
}

// CreatePayment handles payment creation requests
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.OrderID == "" || req.UserID == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: orderId, userId, amount",
		})
	}

	if req.Currency == "" {
		req.Currency = "IDR" // Default currency for Midtrans
	}

	// Check if payment already exists
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	existingPayment, err := h.paymentRepo.GetPaymentByOrderID(ctx, req.OrderID)
	if err != nil {
		log.Printf("Error checking existing payment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if existingPayment != nil {
		// Return existing payment information
		return c.JSON(fiber.Map{
			"id":            existingPayment.ID,
			"orderId":       existingPayment.OrderID,
			"userId":        existingPayment.UserID,
			"amount":        existingPayment.Amount,
			"currency":      existingPayment.Currency,
			"paymentMethod": existingPayment.PaymentMethod,
			"status":        existingPayment.Status,
			"redirectUrl":   existingPayment.RedirectURL,
		})
	}

	// Create payment with Midtrans
	paymentReq := &PaymentRequest{
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Description:   req.Description,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CallbackURL:   req.CallbackURL,
	}

	paymentResp, err := h.midtransService.CreatePayment(ctx, paymentReq)
	if err != nil {
		log.Printf("Error creating payment with Midtrans: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment",
		})
	}

	// Save payment to database
	payment := &Payment{
		ID:            paymentResp.ID,
		OrderID:       paymentResp.OrderID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentMethod: paymentResp.PaymentMethod,
		Status:        paymentResp.Status,
		RedirectURL:   paymentResp.RedirectURL,
	}

	if err := h.paymentRepo.CreatePayment(ctx, payment); err != nil {
		log.Printf("Error saving payment to database: %v", err)
		// Continue anyway, as the payment was created in Midtrans
	}

	// Publish payment event
	event := &PaymentEventMessage{
		ID:            paymentResp.ID,
		OrderID:       paymentResp.OrderID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentMethod: paymentResp.PaymentMethod,
		Status:        paymentResp.Status,
		RedirectURL:   paymentResp.RedirectURL,
		Timestamp:     time.Now(),
	}

	if err := h.natsService.PublishPaymentEvent(event); err != nil {
		log.Printf("Error publishing payment event: %v", err)
		// Continue anyway, as the payment was created
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":            paymentResp.ID,
		"orderId":       paymentResp.OrderID,
		"userId":        paymentResp.UserID,
		"amount":        paymentResp.Amount,
		"currency":      paymentResp.Currency,
		"paymentMethod": paymentResp.PaymentMethod,
		"status":        paymentResp.Status,
		"redirectUrl":   paymentResp.RedirectURL,
	})
}

// GetPaymentStatus handles payment status requests
func (h *PaymentHandler) GetPaymentStatus(c *fiber.Ctx) error {
	// Extract order ID from URL params
	orderID := c.Params("orderId")

	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Get payment status from Midtrans
	paymentResp, err := h.midtransService.GetPaymentStatus(ctx, orderID)
	if err != nil {
		log.Printf("Error getting payment status from Midtrans: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get payment status",
		})
	}

	// Update payment status in database
	err = h.paymentRepo.UpdatePaymentStatus(ctx, orderID, paymentResp.Status)
	if err != nil {
		log.Printf("Error updating payment status in database: %v", err)
		// Continue anyway, as we got the status from Midtrans
	}

	// Publish payment event
	event := &PaymentEventMessage{
		ID:            paymentResp.ID,
		OrderID:       paymentResp.OrderID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentMethod: paymentResp.PaymentMethod,
		Status:        paymentResp.Status,
		Timestamp:     time.Now(),
	}

	if err := h.natsService.PublishPaymentEvent(event); err != nil {
		log.Printf("Error publishing payment event: %v", err)
		// Continue anyway, as we got the status
	}

	// Return response
	return c.JSON(fiber.Map{
		"id":            paymentResp.ID,
		"orderId":       paymentResp.OrderID,
		"amount":        paymentResp.Amount,
		"status":        paymentResp.Status,
		"paymentMethod": paymentResp.PaymentMethod,
	})
}

// HandleNotification handles payment notifications from Midtrans
func (h *PaymentHandler) HandleNotification(c *fiber.Ctx) error {
	// Read notification body
	var notificationBody map[string]interface{}
	if err := c.BodyParser(&notificationBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification body",
		})
	}

	// Convert notification body back to JSON
	notificationJSON, err := json.Marshal(notificationBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process notification",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	// Process notification
	paymentResp, err := h.midtransService.HandleNotification(ctx, notificationJSON)
	if err != nil {
		log.Printf("Error handling notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process notification",
		})
	}

	// Update payment status in database
	err = h.paymentRepo.UpdatePaymentStatus(ctx, paymentResp.OrderID, paymentResp.Status)
	if err != nil {
		log.Printf("Error updating payment status in database: %v", err)
		// Continue anyway, as we processed the notification
	}

	// Publish payment event
	event := &PaymentEventMessage{
		ID:            paymentResp.ID,
		OrderID:       paymentResp.OrderID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentMethod: paymentResp.PaymentMethod,
		Status:        paymentResp.Status,
		Timestamp:     time.Now(),
	}

	if err := h.natsService.PublishPaymentEvent(event); err != nil {
		log.Printf("Error publishing payment event: %v", err)
		// Continue anyway, as we processed the notification
	}

	// Return success response
	return c.JSON(fiber.Map{"status": "ok"})
}

// ProcessPaymentRequest processes payment requests from NATS
func (h *PaymentHandler) ProcessPaymentRequest(ctx context.Context, request *PaymentRequestMessage) (*PaymentEventMessage, error) {
	// Check if payment already exists
	existingPayment, err := h.paymentRepo.GetPaymentByOrderID(ctx, request.OrderID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing payment: %w", err)
	}

	if existingPayment != nil {
		// Return existing payment information
		return &PaymentEventMessage{
			ID:            existingPayment.ID,
			OrderID:       existingPayment.OrderID,
			UserID:        existingPayment.UserID,
			Amount:        existingPayment.Amount,
			Currency:      existingPayment.Currency,
			PaymentMethod: existingPayment.PaymentMethod,
			Status:        existingPayment.Status,
			RedirectURL:   existingPayment.RedirectURL,
			Timestamp:     time.Now(),
		}, nil
	}

	// Create payment with Midtrans
	paymentReq := &PaymentRequest{
		OrderID:       request.OrderID,
		UserID:        request.UserID,
		Amount:        request.Amount,
		Currency:      request.Currency,
		PaymentMethod: request.PaymentMethod,
		Description:   request.Description,
		CallbackURL:   request.CallbackURL,
	}

	paymentResp, err := h.midtransService.CreatePayment(ctx, paymentReq)
	if err != nil {
		return nil, fmt.Errorf("error creating payment with Midtrans: %w", err)
	}

	// Save payment to database
	payment := &Payment{
		ID:            paymentResp.ID,
		OrderID:       paymentResp.OrderID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentMethod: paymentResp.PaymentMethod,
		Status:        paymentResp.Status,
		RedirectURL:   paymentResp.RedirectURL,
	}

	if err := h.paymentRepo.CreatePayment(ctx, payment); err != nil {
		log.Printf("Error saving payment to database: %v", err)
		// Continue anyway, as the payment was created in Midtrans
	}

	// Return payment event
	return &PaymentEventMessage{
		ID:            paymentResp.ID,
		OrderID:       paymentResp.OrderID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentMethod: paymentResp.PaymentMethod,
		Status:        paymentResp.Status,
		RedirectURL:   paymentResp.RedirectURL,
		Timestamp:     time.Now(),
	}, nil
}
