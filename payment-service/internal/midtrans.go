package internal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService struct {
	cfg        *Config
	snapClient snap.Client
	coreClient coreapi.Client
}

// internal type representing subscription API response
type subscriptionAPIResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Amount   string `json:"amount"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	Schedule struct {
		Interval      int    `json:"interval"`
		IntervalUnit  string `json:"interval_unit"`
		NextExecution string `json:"next_execution_at"`
	} `json:"schedule"`
}

// SubscriptionRequest represents a Midtrans subscription creation request
type SubscriptionRequest struct {
	Name          string
	UserID        string
	Amount        float64
	Currency      string
	Token         string
	Interval      string // day, week, month
	IntervalCount int
	StartAt       string // RFC3339 or empty for immediate
	Description   string
	CustomerName  string
	CustomerEmail string
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
	Type          PaymentType
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
	Type          PaymentType
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

func (m *MidtransService) midtransBaseURL() string {
	if m.cfg.MidtransEnvironment == "production" {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

func (m *MidtransService) basicAuthHeader() string {
	basic := m.cfg.MidtransServerKey + ":"
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(basic))
}

// CreateSubscription creates a Midtrans subscription and returns a PaymentResponse
func (m *MidtransService) CreateSubscription(ctx context.Context, req *SubscriptionRequest) (*PaymentResponse, error) {
	payload := map[string]any{
		"name":         req.Name,
		"amount":       fmt.Sprintf("%0.0f", req.Amount),
		"currency":     req.Currency,
		"payment_type": "credit_card",
		"token":        req.Token,
		"schedule": map[string]any{
			"interval":      req.IntervalCount,
			"interval_unit": req.Interval,
		},
		"metadata": map[string]any{
			"user_id":     req.UserID,
			"description": req.Description,
		},
		"customer_details": map[string]any{
			"first_name": req.CustomerName,
			"email":      req.CustomerEmail,
		},
	}
	if req.StartAt != "" {
		// Midtrans expects start_time for subscription start
		if sch, ok := payload["schedule"].(map[string]any); ok {
			sch["start_time"] = req.StartAt
		}
	}

	data, _ := json.Marshal(payload)
	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodPost, m.midtransBaseURL()+"/v1/subscriptions", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	reqHTTP.Header.Set("Authorization", m.basicAuthHeader())
	reqHTTP.Header.Set("Content-Type", "application/json")

	respHTTP, err := http.DefaultClient.Do(reqHTTP)
	if err != nil {
		return nil, fmt.Errorf("failed to call Midtrans subscription API: %w", err)
	}
	defer respHTTP.Body.Close()

	if respHTTP.StatusCode < 200 || respHTTP.StatusCode >= 300 {
		b, _ := io.ReadAll(respHTTP.Body)
		return nil, fmt.Errorf("midtrans subscription API error: %s", string(b))
	}

	var subResp subscriptionAPIResponse
	if err := json.NewDecoder(respHTTP.Body).Decode(&subResp); err != nil {
		return nil, fmt.Errorf("failed to decode Midtrans subscription response: %w", err)
	}

	// Map subscription status
	var status PaymentStatus
	switch subResp.Status {
	case "active":
		status = PaymentStatusSuccess
	case "inactive", "paused":
		status = PaymentStatusCancelled
	default:
		status = PaymentStatusPending
	}

	amount, _ := strconv.ParseFloat(subResp.Amount, 64)
	return &PaymentResponse{
		ID:            subResp.ID,
		OrderID:       subResp.ID,
		UserID:        req.UserID,
		Amount:        amount,
		Currency:      req.Currency,
		Status:        status,
		Timestamp:     time.Now(),
		Type:          PaymentTypeSubscription,
		PaymentMethod: string(PaymentMethodCreditCard),
	}, nil
}

func (m *MidtransService) getSubscription(id string) (*subscriptionAPIResponse, error) {
	reqHTTP, err := http.NewRequest(http.MethodGet, m.midtransBaseURL()+"/v1/subscriptions/"+id, nil)
	if err != nil {
		return nil, err
	}
	reqHTTP.Header.Set("Authorization", m.basicAuthHeader())
	respHTTP, err := http.DefaultClient.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer respHTTP.Body.Close()
	if respHTTP.StatusCode < 200 || respHTTP.StatusCode >= 300 {
		b, _ := io.ReadAll(respHTTP.Body)
		return nil, fmt.Errorf("midtrans get subscription error: %s", string(b))
	}
	var subResp subscriptionAPIResponse
	if err := json.NewDecoder(respHTTP.Body).Decode(&subResp); err != nil {
		return nil, err
	}
	return &subResp, nil
}

// CreatePayment creates a new payment
func (m *MidtransService) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	// Generate a unique transaction ID if not provided
	if req.OrderID == "" {
		prefix := "ORDER-"
		if req.Type == PaymentTypeBalance {
			prefix = "TOPUP-"
		} else if req.Type == PaymentTypeProduct || req.Type == "" {
			prefix = "PROD-"
		}
		req.OrderID = fmt.Sprintf("%s%s", prefix, uuid.New().String())
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
			}
		case PaymentMethodGopay:
			snapReq.EnabledPayments = []snap.SnapPaymentType{snap.PaymentTypeGopay}
		case PaymentMethodShopeePay:
			snapReq.EnabledPayments = []snap.SnapPaymentType{snap.PaymentTypeShopeepay}
		case PaymentMethodQRIS:
			// QRIS support not available in current SDK version; ignoring
		}
	}

	// Set callback URL (finish) if provided
	if req.CallbackURL != "" {
		snapReq.Callbacks = &snap.Callbacks{Finish: req.CallbackURL}
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
		Type:          req.Type,
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
	amount, _ := strconv.ParseFloat(transactionStatusResp.GrossAmount, 64)
	resp := &PaymentResponse{
		ID:            transactionStatusResp.TransactionID,
		OrderID:       orderID,
		Amount:        amount,
		Status:        status,
		PaymentMethod: transactionStatusResp.PaymentType,
		Timestamp:     time.Now(),
	}

	log.Printf("Retrieved payment status for order %s: %s", orderID, status)
	return resp, nil
}

// HandleNotification handles payment notification from Midtrans
func (m *MidtransService) HandleNotification(ctx context.Context, notificationJSON []byte) (*PaymentResponse, error) {
	// Parse only the fields we need from notification JSON
	var notification map[string]any
	if err := json.Unmarshal(notificationJSON, &notification); err != nil {
		return nil, fmt.Errorf("failed to parse Midtrans notification: %w", err)
	}

	// Handle normal transaction notifications (Snap/Core)
	if orderIDRaw, ok := notification["order_id"]; ok {
		orderID, _ := orderIDRaw.(string)
		if orderID == "" {
			return nil, fmt.Errorf("notification missing order_id")
		}
		transactionStatusResp, err := m.coreClient.CheckTransaction(orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to check Midtrans transaction status: %w", err)
		}
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
		amount, _ := strconv.ParseFloat(transactionStatusResp.GrossAmount, 64)
		resp := &PaymentResponse{
			ID:            transactionStatusResp.TransactionID,
			OrderID:       orderID,
			Amount:        amount,
			Status:        status,
			PaymentMethod: transactionStatusResp.PaymentType,
			Timestamp:     time.Now(),
		}
		log.Printf("Processed payment notification for order %s: %s", orderID, status)
		return resp, nil
	}

	// Handle subscription notifications
	if subIDRaw, ok := notification["subscription_id"]; ok {
		subID, _ := subIDRaw.(string)
		if subID == "" {
			return nil, fmt.Errorf("notification missing subscription_id")
		}
		// Query subscription status
		subResp, err := m.getSubscription(subID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subscription status: %w", err)
		}
		// Map subscription status
		var status PaymentStatus
		switch subResp.Status {
		case "active":
			status = PaymentStatusSuccess
		case "inactive", "paused":
			status = PaymentStatusCancelled
		default:
			status = PaymentStatusPending
		}
		amount, _ := strconv.ParseFloat(subResp.Amount, 64)
		resp := &PaymentResponse{
			ID:            subResp.ID,
			OrderID:       subResp.ID,
			Amount:        amount,
			Status:        status,
			PaymentMethod: "credit_card",
			Timestamp:     time.Now(),
			Type:          PaymentTypeSubscription,
		}
		log.Printf("Processed subscription notification for subscription %s: %s", subResp.ID, status)
		return resp, nil
	}

	return nil, fmt.Errorf("notification missing order_id/subscription_id")
}
