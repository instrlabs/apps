package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService struct {
	cfg        *Config
	snapClient snap.Client
	coreClient coreapi.Client
}

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodGopay        PaymentMethod = "gopay"
	PaymentMethodShopeePay    PaymentMethod = "shopeepay"
	PaymentMethodQRIS         PaymentMethod = "qris"
)

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID       string
	UserID        string
	Amount        float64
	Currency      string
	PaymentMethod string
	Description   string
	CustomerName  string
	CustomerEmail string
	CallbackURL   string
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID            string
	OrderID       string
	UserID        string
	Amount        float64
	Currency      string
	PaymentMethod string
	Status        PaymentStatus
	RedirectURL   string
	Timestamp     time.Time
}

func NewMidtransService(cfg *Config) *MidtransService {
	// Set Midtrans environment
	var environment midtrans.EnvironmentType
	if cfg.MidtransEnvironment == "production" {
		environment = midtrans.Production
	} else {
		environment = midtrans.Sandbox
	}

	// Initialize Snap client
	snapClient := snap.Client{}
	snapClient.New(cfg.MidtransServerKey, environment)

	// Initialize Core API client
	coreClient := coreapi.Client{}
	coreClient.New(cfg.MidtransServerKey, environment)

	return &MidtransService{
		cfg:        cfg,
		snapClient: snapClient,
		coreClient: coreClient,
	}
}

// CreatePayment creates a new payment
func (m *MidtransService) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	// Generate a unique transaction ID if not provided
	if req.OrderID == "" {
		req.OrderID = fmt.Sprintf("ORDER-%s", uuid.New().String())
	}

	// Create transaction details
	transactionDetails := midtrans.TransactionDetails{
		OrderID:  req.OrderID,
		GrossAmt: int64(req.Amount),
	}

	// Create customer details if available
	customerDetails := &midtrans.CustomerDetails{
		FName: req.CustomerName,
		Email: req.CustomerEmail,
	}

	// Create Snap request
	snapReq := &snap.Request{
		TransactionDetails: transactionDetails,
		CustomerDetail:     customerDetails,
	}

	// Set payment method if specified
	if req.PaymentMethod != "" {
		switch PaymentMethod(req.PaymentMethod) {
		case PaymentMethodCreditCard:
			snapReq.EnabledPayments = []snap.SnapPaymentType{snap.PaymentTypeCreditCard}
		case PaymentMethodBankTransfer:
			snapReq.EnabledPayments = []snap.SnapPaymentType{
				snap.PaymentTypeBCAVA,
				snap.PaymentTypeBNIVA,
				snap.PaymentTypePermataVA,
				snap.PaymentTypeMandiriVA,
			}
		case PaymentMethodGopay:
			snapReq.EnabledPayments = []snap.SnapPaymentType{snap.PaymentTypeGopay}
		case PaymentMethodShopeePay:
			snapReq.EnabledPayments = []snap.SnapPaymentType{snap.PaymentTypeShopeepay}
		case PaymentMethodQRIS:
			snapReq.EnabledPayments = []snap.SnapPaymentType{snap.PaymentTypeQRIS}
		}
	}

	// Set callback URLs if provided
	if req.CallbackURL != "" {
		snapReq.CallbackURL = req.CallbackURL
	}

	// Create transaction
	snapResp, err := m.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	// Create payment response
	resp := &PaymentResponse{
		ID:            uuid.New().String(),
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		Status:        PaymentStatusPending,
		RedirectURL:   snapResp.RedirectURL,
		Timestamp:     time.Now(),
	}

	log.Printf("Created payment for order %s with redirect URL: %s", resp.OrderID, resp.RedirectURL)
	return resp, nil
}

// GetPaymentStatus gets the status of a payment
func (m *MidtransService) GetPaymentStatus(ctx context.Context, orderID string) (*PaymentResponse, error) {
	// Get transaction status from Midtrans
	transactionStatusResp, err := m.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check Midtrans transaction status: %w", err)
	}

	// Map Midtrans transaction status to our payment status
	var status PaymentStatus
	switch transactionStatusResp.TransactionStatus {
	case "capture", "settlement":
		status = PaymentStatusSuccess
	case "pending":
		status = PaymentStatusPending
	case "deny", "failure":
		status = PaymentStatusFailed
	case "cancel":
		status = PaymentStatusCancelled
	case "expire":
		status = PaymentStatusExpired
	case "refund":
		status = PaymentStatusRefunded
	default:
		status = PaymentStatusPending
	}

	// Create payment response
	resp := &PaymentResponse{
		ID:            transactionStatusResp.TransactionID,
		OrderID:       orderID,
		Amount:        float64(transactionStatusResp.GrossAmount),
		Status:        status,
		PaymentMethod: transactionStatusResp.PaymentType,
		Timestamp:     time.Now(),
	}

	log.Printf("Retrieved payment status for order %s: %s", orderID, status)
	return resp, nil
}

// HandleNotification handles payment notification from Midtrans
func (m *MidtransService) HandleNotification(ctx context.Context, notificationJSON []byte) (*PaymentResponse, error) {
	// Parse notification
	notification, err := coreapi.ParseWebhookJSON(notificationJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Midtrans notification: %w", err)
	}

	// Get transaction status from Midtrans
	transactionStatusResp, err := m.coreClient.CheckTransaction(notification.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check Midtrans transaction status: %w", err)
	}

	// Map Midtrans transaction status to our payment status
	var status PaymentStatus
	switch transactionStatusResp.TransactionStatus {
	case "capture", "settlement":
		status = PaymentStatusSuccess
	case "pending":
		status = PaymentStatusPending
	case "deny", "failure":
		status = PaymentStatusFailed
	case "cancel":
		status = PaymentStatusCancelled
	case "expire":
		status = PaymentStatusExpired
	case "refund":
		status = PaymentStatusRefunded
	default:
		status = PaymentStatusPending
	}

	// Create payment response
	resp := &PaymentResponse{
		ID:            transactionStatusResp.TransactionID,
		OrderID:       notification.OrderID,
		Amount:        float64(transactionStatusResp.GrossAmount),
		Status:        status,
		PaymentMethod: transactionStatusResp.PaymentType,
		Timestamp:     time.Now(),
	}

	log.Printf("Processed payment notification for order %s: %s", notification.OrderID, status)
	return resp, nil
}
