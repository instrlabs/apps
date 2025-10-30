package internal

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2/log"
)

func RefreshAccessToken(refreshToken string) (*http.Response, error) {
	authRefreshURL := os.Getenv("AUTH_SERVICE") + "/refresh"

	req, err := http.NewRequest("POST", authRefreshURL, bytes.NewReader([]byte{}))
	if err != nil {
		log.Errorf("RefreshAccessToken: Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-user-refresh", refreshToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("RefreshAccessToken: Failed to call auth service: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Warnf("RefreshAccessToken: Auth service returned status %d", resp.StatusCode)
		return nil, ErrTokenInvalid
	}

	// Consume body since we only need cookies from headers
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return resp, nil
}
