package internal

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/instrlabs/shared/functionx"
)

func GenerateSixDigitPIN() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Errorf("GenerateSixDigitPIN: Failed to generate six digit PIN: %v", err)
		return ""
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func GenerateAccessToken(userID string) string {
	jwtSecret := functionx.GetEnvString("JWT_SECRET", "secret")
	now := time.Now().UTC()
	expirationTime := now.Add(1 * time.Hour)

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

func GenerateRefreshToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Errorf("GenerateRefreshToken: Failed to generate refresh token: %v", err)
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func GenerateUniqueUsername(email string) (string, error) {
	base := email
	if at := strings.Index(email, "@"); at != -1 {
		base = email[:at]
	}
	base = strings.ToLower(strings.TrimSpace(base))
	if base == "" {
		base = "user"
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	suffix := fmt.Sprintf("%04d", nBig.Int64())
	candidate := fmt.Sprintf("%s%s", base, suffix)

	return candidate, nil
}
