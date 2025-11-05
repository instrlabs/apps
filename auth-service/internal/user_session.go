package internal

import (
	"crypto/sha256"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID         string             `json:"user_id" bson:"user_id"`
	DeviceHash     string             `json:"-" bson:"device_hash"`
	IPAddress      string             `json:"ip_address" bson:"ip_address"`
	UserAgent      string             `json:"-" bson:"user_agent"`
	RefreshToken   string             `json:"-" bson:"refresh_token"`
	IsActive       bool               `json:"is_active" bson:"is_active"`
	LastActivityAt time.Time          `json:"last_activity_at" bson:"last_activity_at"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	ExpiresAt      time.Time          `json:"-" bson:"expires_at"`
}

func GenerateDeviceHash(ipAddress, userAgent string) string {
	h := sha256.New()
	h.Write([]byte(ipAddress + userAgent))
	return fmt.Sprintf("%x", h.Sum(nil))
}
