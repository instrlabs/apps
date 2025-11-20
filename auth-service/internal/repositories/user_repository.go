package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/instrlabs/auth-service/internal/models"
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	FindByGoogleID(googleID string) (*models.User, error)
	Update(user *models.User) error
	SetPinWithExpiry(email, hashedPin string) error
	ClearPin(userID string) error
	SetRegisteredAt(userID string) error
	AddRefreshToken(userID, token string) error
	RemoveRefreshToken(userID, token string) error
	ValidateRefreshToken(userID, token string) error
	ClearAllRefreshTokens(userID string) error
}

// UserRepository implements UserRepositoryInterface
type UserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	return &user, err
}

// FindByGoogleID finds a user by Google ID
func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	return &user, err
}

// Update updates a user in the database
func (r *UserRepository) Update(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now().UTC()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

// SetPinWithExpiry sets a PIN hash with expiry for a user
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
	return err
}

// ClearPin clears the PIN for a user
func (r *UserRepository) ClearPin(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
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

// SetRegisteredAt sets the registered timestamp for a user
func (r *UserRepository) SetRegisteredAt(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	t := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"registered_at": t,
			"updated_at":    t,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// AddRefreshToken adds a refresh token for a user
func (r *UserRepository) AddRefreshToken(userID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$addToSet": bson.M{
			"refresh_tokens": token,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// RemoveRefreshToken removes a specific refresh token for a user
func (r *UserRepository) RemoveRefreshToken(userID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$pull": bson.M{
			"refresh_tokens": token,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// ValidateRefreshToken validates that a refresh token exists for a user
func (r *UserRepository) ValidateRefreshToken(userID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{
		"_id":            objectID,
		"refresh_tokens": token,
	}).Decode(&user)

	return err
}

// ClearAllRefreshTokens removes all refresh tokens for a user
func (r *UserRepository) ClearAllRefreshTokens(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"refresh_tokens": []string{},
			"updated_at":     time.Now().UTC(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}
