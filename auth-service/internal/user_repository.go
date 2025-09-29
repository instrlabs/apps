package internal

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instr-labs/shared/init"
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

	// Ensure username is set and unique
	if strings.TrimSpace(user.Username) == "" {
		uname, genErr := r.generateUniqueUsername(ctx, user.Email)
		if genErr != nil {
			log.Errorf("Failed to generate unique username for %s: %v", user.Email, genErr)
			return nil
		}
		user.Username = uname
	}

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Failed to create user: %v", err)
		return nil
	}

	return user
}

func (r *UserRepository) FindByEmail(email string) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnf("User not found for email: %s", email)
			return nil
		}
		log.Errorf("Failed to find user by email: %v", err)
		return nil
	}

	return &user
}

func (r *UserRepository) FindByID(id string) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	objectID, _ := primitive.ObjectIDFromHex(id)
	err := r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnf("User not found for ID: %s", id)
			return nil
		}
		log.Errorf("Failed to find user by ID: %v", err)
		return nil
	}

	return &user
}

func (r *UserRepository) FindByRefreshToken(userId, refreshToken string) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	objectID, _ := primitive.ObjectIDFromHex(userId)
	err := r.collection.FindOne(ctx, bson.M{
		"_id":           objectID,
		"refresh_token": refreshToken,
	}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnf("Refresh token not found for user %s", userId)
			return nil
		}
		log.Errorf("Failed to find user by userId and refresh token: %v", err)
		return nil
	}

	return &user
}

func (r *UserRepository) UpdateRefreshToken(userID string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	update := bson.M{
		"$set": bson.M{
			"refresh_token": refreshToken,
			"updated_at":    time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Errorf("Failed to update refresh token for user %s: %v", userID, err)
		return err
	}
	return nil
}

func (r *UserRepository) ClearRefreshToken(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	update := bson.M{
		"$set": bson.M{
			"refresh_token": nil,
			"updated_at":    time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Errorf("Failed to clear refresh token for user %s: %v", userID, err)
		return err
	}
	return nil
}

func (r *UserRepository) FindByGoogleID(googleID string) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := r.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnf("Google ID not found for user: %s", googleID)
			return nil
		}

		log.Errorf("Failed to find user by Google ID: %v", err)
		return nil
	}

	return &user
}

func (r *UserRepository) UpdateGoogleID(userID string, googleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	update := bson.M{
		"$set": bson.M{
			"google_id":  googleID,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Errorf("Failed to update google ID for user %s: %v", userID, err)
		return err
	}
	return nil
}

func (r *UserRepository) SetPinWithExpiry(email, hashedPin string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expiry := time.Now().UTC().Add(10 * time.Minute)
	update := bson.M{
		"$set": bson.M{
			"pin_hash":    hashedPin,
			"pin_expires": expiry,
			"updated_at":  time.Now().UTC(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil {
		log.Errorf("Failed to set PIN for user %s: %v", email, err)
		return err
	}
	return nil
}

func (r *UserRepository) ClearPin(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"pin_hash":    nil,
			"pin_expires": nil,
			"updated_at":  time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *UserRepository) SetRegisteredAt(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().UTC()
	objectID, _ := primitive.ObjectIDFromHex(userID)
	update := bson.M{
		"$set": bson.M{
			"registered_at": t,
			"updated_at":    t,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Errorf("Failed to update registered_at for user %s: %v", userID, err)
		return err
	}
	return err
}

func (r *UserRepository) generateUniqueUsername(ctx context.Context, email string) (string, error) {
	base := email
	if at := strings.Index(email, "@"); at != -1 {
		base = email[:at]
	}
	base = strings.ToLower(strings.TrimSpace(base))
	if base == "" {
		base = "user"
	}

	for i := 0; i < 20; i++ { // up to 20 attempts
		nBig, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil {
			return "", err
		}
		suffix := fmt.Sprintf("%04d", nBig.Int64())
		candidate := fmt.Sprintf("%s%s", base, suffix)

		var tmp User
		err = r.collection.FindOne(ctx, bson.M{"username": candidate}).Decode(&tmp)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return candidate, nil
			}
			return "", err
		}
		// username exists, try again
	}
	return "", fmt.Errorf("unable to generate unique username")
}
