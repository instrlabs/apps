package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductClient struct {
	baseURL string
	client  *http.Client
}

type Product struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key       string             `json:"key" bson:"key"`
	Name      string             `json:"name" bson:"name"`
	Type      string             `json:"type" bson:"type"`
	Price     float64            `json:"price" bson:"price"`
	Active    bool               `json:"active" bson:"active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (pc *ProductClient) FindByID(id primitive.ObjectID, productType string) (*Product, error) {
	url := fmt.Sprintf("%s/?type=%s", pc.baseURL, productType)
	log.Printf("Fetching product: id=%s type=%s url=%s", id.Hex(), productType, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, err
	}

	start := time.Now()
	resp, err := pc.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Failed to call product service: duration=%s error=%v", duration.String(), err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Product service response: statusCode=%d duration=%s", resp.StatusCode, duration.String())

	if resp.StatusCode != http.StatusOK {
		log.Printf("Product service returned error: status %d", resp.StatusCode)
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var result struct {
		Message string      `json:"message"`
		Errors  interface{} `json:"errors"`
		Data    []Product   `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		return nil, err
	}

	log.Printf("Products received: count=%d", len(result.Data))

	for _, p := range result.Data {
		if p.ID == id {
			log.Println("Product found")
			return &p, nil
		}
	}

	log.Println("Product not found")
	return nil, nil
}
