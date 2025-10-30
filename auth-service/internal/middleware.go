package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RefreshTokenIfNeeded() fiber.Handler {
	return func(c *fiber.Ctx) error {
		refreshToken := c.Get("x-user-refresh-token")
		if refreshToken == "" {
			return c.Next()
		}

		respRefresh, err := refreshAccessToken(refreshToken)
		if err == nil {
			setCookieHeaders := respRefresh.Header["Set-Cookie"]
			for _, setCookieHeader := range setCookieHeaders {
				c.Response().Header.Add("Set-Cookie", setCookieHeader)
			}

			_, _ = io.Copy(io.Discard, respRefresh.Body)
			_ = respRefresh.Body.Close()

			c.Request().Header.Set("x-user-id", respRefresh.Header.Get("x-user-id"))
		}

		return c.Next()
	}
}

func refreshAccessToken(refreshToken string) (*http.Response, error) {
	authServiceURL := os.Getenv("AUTH_SERVICE")
	refreshURL := strings.TrimSuffix(authServiceURL, "/") + "/refresh"

	req, err := http.NewRequest("POST", refreshURL, bytes.NewReader([]byte{}))
	if err != nil {
		log.Errorf("refreshAccessToken: Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-user-refresh", refreshToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("refreshAccessToken: Failed to call auth service: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		log.Warnf("refreshAccessToken: Auth service returned status %d", resp.StatusCode)
		return nil, fmt.Errorf(ErrInvalidToken)
	}

	return resp, nil
}
