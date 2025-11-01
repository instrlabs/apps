package internal

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	initx "github.com/instrlabs/shared/init"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SessionRepository handles all database operations for user sessions
type SessionRepository struct {
	db         *initx.Mongo
	collection *mongo.Collection
}

// NewSessionRepository creates a new session repository instance
func NewSessionRepository(db *initx.Mongo) *SessionRepository {
	return &SessionRepository{
		db:         db,
		collection: db.DB.Collection("user_sessions"),
	}
}

// CreateSession creates a new session for a user with device binding
// Returns the created session or error
func (r *SessionRepository) CreateSession(userID, ipAddress, userAgent string) (*UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session := &UserSession{
		UserID:         userID,
		SessionID:      GenerateSessionID(),
		DeviceHash:     GenerateDeviceHash(ipAddress, userAgent),
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		IsActive:       true,
		LastActivityAt: time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
		ExpiresAt:      time.Now().UTC().AddDate(0, 0, 30), // 30 days
	}

	_, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		log.Errorf("CreateSession: Failed to create session: %v", err)
		return nil, err
	}

	log.Infof("CreateSession: Session created for user %s with IP %s", userID, ipAddress)
	return session, nil
}

// FindSessionByRefreshToken finds a session by refresh token and checks if active
// Returns nil if not found or expired
func (r *SessionRepository) FindSessionByRefreshToken(refreshToken string) (*UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session UserSession
	err := r.collection.FindOne(ctx, bson.M{
		"refresh_token": refreshToken,
		"is_active":     true,
	}).Decode(&session)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warn("FindSessionByRefreshToken: Session not found or inactive")
			return nil, nil
		}
		log.Errorf("FindSessionByRefreshToken: Failed to find session: %v", err)
		return nil, err
	}

	// Check if session has expired
	if time.Now().UTC().After(session.ExpiresAt) {
		log.Warn("FindSessionByRefreshToken: Session has expired")
		_ = r.DeactivateSession(session.SessionID)
		return nil, nil
	}

	return &session, nil
}

// FindSessionByID finds a specific session by ID and user ID
// Returns nil if not found
func (r *SessionRepository) FindSessionByID(sessionID, userID string) (*UserSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session UserSession
	err := r.collection.FindOne(ctx, bson.M{
		"session_id": sessionID,
		"user_id":    userID,
	}).Decode(&session)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnf("FindSessionByID: Session %s not found for user %s", sessionID, userID)
			return nil, nil
		}
		log.Errorf("FindSessionByID: Failed to find session: %v", err)
		return nil, err
	}

	return &session, nil
}

// ValidateSession checks if a session is valid for the given device
// Returns true only if session exists, is active, not expired, and device hash matches
func (r *SessionRepository) ValidateSession(sessionID, userID, deviceHash string) bool {
	session, err := r.FindSessionByID(sessionID, userID)
	if err != nil || session == nil {
		return false
	}

	if !session.IsActive {
		log.Warnf("ValidateSession: Session %s is inactive", sessionID)
		return false
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		log.Warnf("ValidateSession: Session %s has expired", sessionID)
		_ = r.DeactivateSession(sessionID)
		return false
	}

	if session.DeviceHash != deviceHash {
		log.Warnf("ValidateSession: Device hash mismatch for session %s - possible token theft", sessionID)
		// Deactivate session immediately due to potential token theft
		_ = r.DeactivateSession(sessionID)
		return false
	}

	return true
}

// UpdateSessionActivity updates the last activity timestamp for a session
// Called whenever the session is used
func (r *SessionRepository) UpdateSessionActivity(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"last_activity_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"session_id": sessionID}, update)
	if err != nil {
		log.Errorf("UpdateSessionActivity: Failed to update session: %v", err)
		return err
	}

	return nil
}

// DeactivateSession marks a session as inactive
// Called when user logs out from a device
func (r *SessionRepository) DeactivateSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"session_id": sessionID}, update)
	if err != nil {
		log.Errorf("DeactivateSession: Failed to deactivate session: %v", err)
		return err
	}

	log.Infof("DeactivateSession: Session %s deactivated", sessionID)
	return nil
}

// GetUserSessions returns all active, non-expired sessions for a user
func (r *SessionRepository) GetUserSessions(userID string) ([]UserSession, error) {
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

// ClearExpiredSessions removes expired sessions for a user
// Should be called periodically or after session operations
func (r *SessionRepository) ClearExpiredSessions(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteMany(ctx, bson.M{
		"user_id":    userID,
		"expires_at": bson.M{"$lt": time.Now().UTC()},
	})

	if err != nil {
		log.Errorf("ClearExpiredSessions: Failed to delete expired sessions: %v", err)
		return err
	}

	if result.DeletedCount > 0 {
		log.Infof("ClearExpiredSessions: Deleted %d expired sessions for user %s", result.DeletedCount, userID)
	}

	return nil
}

// ClearAllUserSessions deactivates all active sessions for a user
// Used for "logout all devices" functionality
func (r *SessionRepository) ClearAllUserSessions(userID string) error {
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

// UpdateSessionRefreshToken updates the refresh token for a session
// Called when issuing new tokens after refresh
func (r *SessionRepository) UpdateSessionRefreshToken(sessionID, newRefreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"refresh_token":    newRefreshToken,
			"last_activity_at": time.Now().UTC(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"session_id": sessionID}, update)
	if err != nil {
		log.Errorf("UpdateSessionRefreshToken: Failed to update refresh token: %v", err)
		return err
	}

	return nil
}
