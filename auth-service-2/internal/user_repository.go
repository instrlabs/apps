package internal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByGoogleID finds a user by Google ID
func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now().UTC()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

// UpdatePIN updates the PIN hash and expiry for a user
func (r *UserRepository) UpdatePIN(ctx context.Context, email, pinHash string, expiresAt time.Time) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.M{
			"$set": bson.M{
				"pin_hash":    pinHash,
				"pin_expires": expiresAt,
				"updated_at":  time.Now().UTC(),
			},
		},
	)
	return err
}

// UpdateRefreshToken updates the refresh token for a user
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, userID primitive.ObjectID, refreshToken string, expiresAt time.Time) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"refresh_token":         refreshToken,
				"refresh_token_expires": expiresAt,
				"updated_at":            time.Now().UTC(),
			},
		},
	)
	return err
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now().UTC()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"last_login_at": now,
				"updated_at":    now,
			},
		},
	)
	return err
}

// ClearRefreshToken clears the refresh token (for logout)
func (r *UserRepository) ClearRefreshToken(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$unset": bson.M{
				"refresh_token":         "",
				"refresh_token_expires": "",
			},
			"$set": bson.M{
				"updated_at": time.Now().UTC(),
			},
		},
	)
	return err
}
