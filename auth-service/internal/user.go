package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name,omitempty"`
	Email        string             `json:"email" bson:"email,unique"`
	PinHash      string             `json:"-" bson:"pin_hash,omitempty"`
	PinExpires   time.Time          `json:"-" bson:"pin_expires,omitempty"`
	GoogleID     string             `json:"-" bson:"google_id,omitempty"`
	RefreshToken string             `json:"-" bson:"refresh_token,omitempty"`
	RegisteredAt time.Time          `json:"registered_at" bson:"registered_at,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewUser(email string) *User {
	now := time.Now().UTC()
	return &User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewGoogleUser(name, email, googleID string) *User {
	user := NewUser(email)
	user.Name = name
	user.GoogleID = googleID
	return user
}

func (u *User) ComparePin(pin string) bool {
	if u.PinHash == "" {
		return false
	}
	if !u.PinExpires.IsZero() && time.Now().UTC().After(u.PinExpires) {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.PinHash), []byte(pin))
	return err == nil
}
