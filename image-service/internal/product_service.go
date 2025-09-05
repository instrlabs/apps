package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type Product struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type productResponse struct {
	Message string          `json:"message"`
	Errors  interface{}     `json:"errors"`
	Data    json.RawMessage `json:"data"`
}

type ProductService struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductService() *ProductService {
	base := os.Getenv("PAYMENT_SERVICE_URL")
	return &ProductService{
		baseURL:    base,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *ProductService) GetProduct(id string) *Product {
	url := fmt.Sprintf("%s/products/%s", s.baseURL, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, _ := s.httpClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Errorf("payment-service returned status %d", resp.StatusCode)
		return nil
	}

	var envelope productResponse
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
