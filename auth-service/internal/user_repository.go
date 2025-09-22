package internal

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/histweety-labs/shared/init"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

func NewUserRepository(db *initx.Mongo) *UserRepository {
	return &UserRepository{
		db:         db,
		collection: db.DB.Collection("users"),
	}
}

func (r *UserRepository) Create(user *User) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser User
	err := r.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		log.Errorf("User already exists with email %s", user.Email)
		return nil
	}

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Failed to create user: %v", err)
		return nil
	}

	return user
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var user User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// FindByRefreshToken finds a user by refresh token
func (r *UserRepository) FindByRefreshToken(refreshToken string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateRefreshToken updates the refresh token for a user
func (r *UserRepository) UpdateRefreshToken(userID string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"refresh_token": refreshToken,
			"updated_at":    time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// FindByGoogleID finds a user by Google ID
func (r *UserRepository) FindByGoogleID(googleID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateGoogleID updates the Google ID for a user
func (r *UserRepository) UpdateGoogleID(userID string, googleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"google_id":  googleID,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// UpdateProfile updates the user's profile information (currently only name)
func (r *UserRepository) UpdateProfile(userID string, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// ClearRefreshToken removes the refresh token for a user
func (r *UserRepository) ClearRefreshToken(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"refresh_token": "",
			"updated_at":    time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// SetPinWithExpiry sets the user's PIN hash with an expiry (for OTP scenarios)
func (r *UserRepository) SetPinWithExpiry(email, hashedPin string, expiry time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"pin_hash":    hashedPin,
			"pin_expires": expiry,
			"updated_at":  time.Now().UTC(),
		},
	}
	res, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// ClearPin clears the PIN hash and expiry (e.g., after OTP is used)
func (r *UserRepository) ClearPin(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"pin_hash":    "",
			"pin_expires": time.Time{},
			"updated_at":  time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// UpdatePin updates the persistent pin hash for a user
func (r *UserRepository) UpdatePin(userID, hashedPin string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"pin_hash":   hashedPin,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// UpdateRegisteredAt sets the RegisteredAt timestamp for the user
func (r *UserRepository) UpdateRegisteredAt(userID string, t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"registered_at": t,
			"updated_at":    time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
