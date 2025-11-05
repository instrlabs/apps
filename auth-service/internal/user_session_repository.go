package internal

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserSessionRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserSessionRepository(db *mongo.Database) *UserSessionRepository {
	return &UserSessionRepository{
		db:         db,
		collection: db.Collection("user_sessions"),
	}
}

func (r *UserSessionRepository) CreateUserSession(userID, ipAddress, userAgent string) (*UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := &UserSession{
		ID:             primitive.NewObjectID(),
		UserID:         userID,
		DeviceHash:     GenerateDeviceHash(ipAddress, userAgent),
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		IsActive:       true,
		LastActivityAt: time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
		ExpiresAt:      time.Now().UTC().AddDate(0, 0, 30),
	}

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		log.Errorf("CreateUserSession: Failed to create session: %v", err)
		return nil, err
	}

	log.Infof("CreateUserSession: Session created for user %s with IP %s", userID, ipAddress)
	return session, nil
}

func (r *UserSessionRepository) FindUserSessionByRefreshToken(refreshToken string) (*UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session UserSession
	err := r.collection.FindOne(ctx, bson.M{
		"refresh_token": refreshToken,
		"is_active":     true,
	}).Decode(&session)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warn("FindUserSessionByRefreshToken: Session not found or inactive")
			return nil, nil
		}
		log.Errorf("FindUserSessionByRefreshToken: Failed to find session: %v", err)
		return nil, err
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		log.Warn("FindUserSessionByRefreshToken: Session has expired")
		_ = r.DeactivateUserSession(session.ID)
		return nil, nil
	}

	return &session, nil
}

func (r *UserSessionRepository) FindUserSessionByID(id primitive.ObjectID, userID string) (*UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session UserSession
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     id,
		"user_id": userID,
	}).Decode(&session)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnf("FindUserSessionByID: Session %s not found for user %s", id, userID)
			return nil, nil
		}
		log.Errorf("FindUserSessionByID: Failed to find session: %v", err)
		return nil, err
	}

	return &session, nil
}

func (r *UserSessionRepository) ValidateUserSession(id primitive.ObjectID, userID, deviceHash string) bool {
	session, err := r.FindUserSessionByID(id, userID)
	if err != nil || session == nil {
		return false
	}

	if !session.IsActive {
		log.Warnf("ValidateUserSession: Session %s is inactive", id)
		return false
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		log.Warnf("ValidateUserSession: Session %s has expired", id)
		_ = r.DeactivateUserSession(id)
		return false
	}

	if session.DeviceHash != deviceHash {
		log.Warnf("ValidateUserSession: Device hash mismatch for session %s - possible token theft", id)
		_ = r.DeactivateUserSession(id)
		return false
	}

	return true
}

func (r *UserSessionRepository) UpdateUserSessionActivity(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"last_activity_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Errorf("UpdateUserSessionActivity: Failed to update session: %v", err)
		return err
	}

	return nil
}

func (r *UserSessionRepository) DeactivateUserSession(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Errorf("DeactivateUserSession: Failed to deactivate session: %v", err)
		return err
	}

	log.Infof("DeactivateUserSession: Session %s deactivated", id)
	return nil
}

func (r *UserSessionRepository) GetUserSessions(userID string) ([]UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":    userID,
		"is_active":  true,
		"expires_at": bson.M{"$gt": time.Now().UTC()},
	})

	if err != nil {
		log.Errorf("GetUserSessions: Failed to find sessions: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []UserSession
	if err = cursor.All(ctx, &sessions); err != nil {
		log.Errorf("GetUserSessions: Failed to decode sessions: %v", err)
		return nil, err
	}

	return sessions, nil
}

func (r *UserSessionRepository) ClearExpiredUserSessions(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteMany(ctx, bson.M{
		"user_id":    userID,
		"expires_at": bson.M{"$lt": time.Now().UTC()},
	})

	if err != nil {
		log.Errorf("ClearExpiredUserSessions: Failed to delete expired sessions: %v", err)
		return err
	}

	if result.DeletedCount > 0 {
		log.Infof("ClearExpiredUserSessions: Deleted %d expired sessions for user %s", result.DeletedCount, userID)
	}

	return nil
}

func (r *UserSessionRepository) ClearAllUserSessions(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}

	result, err := r.collection.UpdateMany(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		log.Errorf("ClearAllUserSessions: Failed to deactivate sessions: %v", err)
		return err
	}

	log.Infof("ClearAllUserSessions: Deactivated %d sessions for user %s", result.ModifiedCount, userID)
	return nil
}

func (r *UserSessionRepository) UpdateUserSessionRefreshToken(id primitive.ObjectID, newRefreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"refresh_token":    newRefreshToken,
			"last_activity_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		log.Errorf("UpdateUserSessionRefreshToken: Failed to update refresh token: %v", err)
		return err
	}

	return nil
}
