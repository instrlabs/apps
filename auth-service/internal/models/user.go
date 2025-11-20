package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity in the system
type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username      string             `json:"username" bson:"username"`
	Email         string             `json:"email" bson:"email"`
	PinHash       *string            `json:"-" bson:"pin_hash"`
	PinExpires    *time.Time         `json:"-" bson:"pin_expires"`
	GoogleID      *string            `json:"-" bson:"google_id"`
	RefreshTokens []string           `json:"-" bson:"refresh_tokens,omitempty"`
	RegisteredAt  *time.Time         `json:"registered_at" bson:"registered_at"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewUser creates a new user instance
func NewUser(email string) *User {
	now := time.Now().UTC()
	username, _ := GenerateUniqueUsername(email)
	return &User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewGoogleUser creates a new user with Google account
func NewGoogleUser(email, googleID string) *User {
	user := NewUser(email)
	user.GoogleID = &googleID
	return user
}

// ComparePin validates the provided PIN against stored hash
func (u *User) ComparePin(pin string) bool {
	if u.PinHash == nil || *u.PinHash == "" {
		return false
	}

	if u.PinExpires != nil && !u.PinExpires.IsZero() && time.Now().UTC().After(*u.PinExpires) {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(*u.PinHash), []byte(pin))
	return err == nil
}

// IsPinExpired checks if the PIN has expired
func (u *User) IsPinExpired() bool {
	if u.PinExpires == nil || u.PinExpires.IsZero() {
		return false
	}
	return time.Now().UTC().After(*u.PinExpires)
}

// HasRefreshToken checks if user has the specified refresh token
func (u *User) HasRefreshToken(token string) bool {
	for _, t := range u.RefreshTokens {
		if t == token {
			return true
		}
	}
	return false
}
