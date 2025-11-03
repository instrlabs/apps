package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserSession represents a user's device session for binding tokens
// Each login creates a new session with device-specific binding
type UserSession struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         string             `json:"user_id" bson:"user_id"`                                    // User ID this session belongs to
	SessionID      string             `json:"session_id" bson:"session_id"`                              // Unique per session
	DeviceHash     string             `json:"-" bson:"device_hash"`                                      // Hash of IP + User-Agent
	IPAddress      string             `json:"ip_address" bson:"ip_address"`                              // Device IP for reference
	UserAgent      string             `json:"-" bson:"user_agent"`                                       // Device User-Agent (not exposed in JSON)
	RefreshToken   string             `json:"-" bson:"refresh_token"`                                    // Stored but not exposed
	IsActive       bool               `json:"is_active" bson:"is_active"`                                // Track if session is active
	LastActivityAt time.Time          `json:"last_activity_at" bson:"last_activity_at"`                  // Last request time
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`                              // Session creation time
	ExpiresAt      time.Time          `json:"-" bson:"expires_at"`                                       // Session expiry time
}

// GenerateDeviceHash creates a SHA256 hash of IP + User-Agent
// Used to bind tokens to specific devices
// If device info changes (IP or User-Agent), hash will be different
// This detects potential token theft or hijacking
func GenerateDeviceHash(ipAddress, userAgent string) string {
	h := sha256.New()
	h.Write([]byte(ipAddress + userAgent))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GenerateSessionID creates a unique random session ID
// Each device session gets a unique ID in base64 URL-safe format
// 16 bytes = 128 bits of randomness = ~16^2 = 256 combinations
func GenerateSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
