package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`

	// OAuth Fields
	GoogleID *string `json:"-" bson:"google_id,omitempty"`

	// PIN Authentication Fields
	PinHash    *string    `json:"-" bson:"pin_hash,omitempty"`
	PinExpires *time.Time `json:"-" bson:"pin_expires,omitempty"`

	// Token Management
	RefreshToken        *string    `json:"-" bson:"refresh_token,omitempty"`
	RefreshTokenExpires *time.Time `json:"-" bson:"refresh_token_expires,omitempty"`

	// Metadata
	IsVerified   bool       `json:"is_verified" bson:"is_verified"`
	RegisteredAt *time.Time `json:"registered_at,omitempty" bson:"registered_at,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" bson:"updated_at"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *User     `json:"user"`
}

// PinRequest represents a request to send a PIN
type PinRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PinVerifyRequest represents a request to verify a PIN
type PinVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Pin   string `json:"pin" validate:"required"`
}

// RefreshTokenRequest represents a request to refresh an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// OAuthState represents the state stored during OAuth flow
type OAuthState struct {
	State     string    `json:"state"`
	Provider  string    `json:"provider"`
	ExpiresAt time.Time `json:"expires_at"`
}
