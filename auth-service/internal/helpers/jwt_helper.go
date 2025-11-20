package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/instrlabs/shared/functionx"
)

// GenerateAccessToken generates a JWT access token for the given user ID
func GenerateAccessToken(userID string, expiryHours int) string {
	jwtSecret := functionx.GetEnvString("JWT_SECRET", "secret")
	now := time.Now().UTC()
	expirationTime := now.Add(time.Duration(expiryHours) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Errorf("GenerateAccessToken: Failed to generate access token: %v", err)
		return ""
	}

	return tokenString
}

// GenerateRefreshToken generates a cryptographically secure refresh token
func GenerateRefreshToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Errorf("GenerateRefreshToken: Failed to generate refresh token: %v", err)
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
