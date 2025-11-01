package internal

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenService struct {
	config *Config
}

func NewTokenService(config *Config) *TokenService {
	return &TokenService{
		config: config,
	}
}

// TokenClaims represents the JWT claims
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a JWT access token for a user
func (s *TokenService) GenerateAccessToken(user *User) (string, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(s.config.AccessTokenExpiry) * time.Hour)

	claims := TokenClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "auth-service-2",
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// GenerateRefreshToken generates a JWT refresh token for a user
func (s *TokenService) GenerateRefreshToken(user *User) (string, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(s.config.RefreshTokenExpiry) * time.Hour)

	claims := TokenClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "auth-service-2",
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// ValidateToken validates a JWT token and returns the claims
func (s *TokenService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a refresh token specifically
func (s *TokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// ParseUserID parses a user ID string to ObjectID
func ParseUserID(userIDStr string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(userIDStr)
}
