package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"-" bson:"password"`
	GoogleID          string             `json:"-" bson:"google_id,omitempty"`
	RefreshToken      string             `json:"-" bson:"refresh_token,omitempty"`
	ResetToken        string             `json:"-" bson:"reset_token,omitempty"`
	ResetTokenExpires time.Time          `json:"-" bson:"reset_token_expires,omitempty"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewUser creates a new user with the given name, email and password
func NewUser(name, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NewGoogleUser creates a new user from a Google profile
func NewGoogleUser(name, email, googleID string) *User {
	now := time.Now().UTC()
	return &User{
		Name:      name,
		Email:     email,
		GoogleID:  googleID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ComparePassword compares the given password with the user's password
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
