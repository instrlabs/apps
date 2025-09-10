package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Product struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type apiResponse struct {
	Message string          `json:"message"`
	Errors  interface{}     `json:"errors"`
	Data    json.RawMessage `json:"data"`
}

type PaymentService struct {
	baseURL    string
	httpClient *http.Client
}

func NewPaymentService(cfg *Config) *PaymentService {
	return &PaymentService{
		baseURL:    cfg.PaymentServiceURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *PaymentService) GetProduct(c *fiber.Ctx, id string) *Product {
	url := fmt.Sprintf("%s/products/%s", s.baseURL, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("X-Authenticated", "true")
	req.Header.Set("X-User-Id", c.Locals("UserID").(string))
	req.Header.Set("X-User-Roles", c.Locals("Roles").(string))
	resp, _ := s.httpClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Errorf("payment-service returned status %d", resp.StatusCode)
		return nil
	}

	var envelope apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		log.Errorf("Failed to decode response: %v", err)
		return nil
	}

	var product Product
	if err := json.Unmarshal(envelope.Data, &product); err != nil {
		log.Errorf("Failed to unmarshal response: %v", err)
		return nil
	}

	return &product
}
