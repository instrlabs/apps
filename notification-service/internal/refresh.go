package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2/log"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	authURL := os.Getenv("AUTH_SERVICE")
	if authURL == "" {
		authURL = "http://auth-service:3000"
	}
	authURL = authURL + "/auth/refresh"

	req, err := http.NewRequest("POST", authURL, bytes.NewReader([]byte{}))
	if err != nil {
		log.Errorf("RefreshAccessToken: Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("Cookie", "refresh_token="+refreshToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("RefreshAccessToken: Failed to call auth service: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warnf("RefreshAccessToken: Auth service returned status %d", resp.StatusCode)
		return nil, ErrTokenInvalid
	}

	var tokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		log.Errorf("RefreshAccessToken: Failed to decode response: %v", err)
		return nil, err
	}

	return &tokens, nil
}
