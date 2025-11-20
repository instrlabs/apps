package services

import (
	"fmt"

	"github.com/instrlabs/auth-service/internal/helpers"
	"github.com/instrlabs/auth-service/internal/models"
	"github.com/instrlabs/auth-service/internal/repositories"
)

// TokenResponse represents the response for login/refresh token operations
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthService handles authentication operations
type AuthService struct {
	userRepo         repositories.UserRepositoryInterface
	tokenExpiryHours int
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repositories.UserRepositoryInterface, tokenExpiryHours int) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		tokenExpiryHours: tokenExpiryHours,
	}
}

// Login authenticates a user with email and PIN
func (s *AuthService) Login(email, pin string) (*TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Validate PIN
	if !user.ComparePin(pin) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Handle post-login tasks
	s.handlePostLoginTasks(user)

	// Create tokens
	return s.createTokensForUser(user.ID.Hex())
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(userID, refreshToken string) (*TokenResponse, error) {
	// Validate refresh token
	if err := s.userRepo.ValidateRefreshToken(userID, refreshToken); err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Rotate tokens (remove old, add new)
	return s.rotateRefreshTokens(user.ID.Hex(), refreshToken)
}

// Logout handles user logout by removing refresh tokens
func (s *AuthService) Logout(userID, refreshToken string) error {
	if refreshToken != "" {
		// Logout specific token
		return s.userRepo.RemoveRefreshToken(userID, refreshToken)
	}

	// Logout all tokens
	return s.userRepo.ClearAllRefreshTokens(userID)
}

// handlePostLoginTasks handles tasks after successful login
func (s *AuthService) handlePostLoginTasks(user *models.User) {
	// Clear PIN after successful login
	_ = s.userRepo.ClearPin(user.ID.Hex())

	// Set registered at if not already set
	if user.RegisteredAt == nil {
		_ = s.userRepo.SetRegisteredAt(user.ID.Hex())
	}
}

// createTokensForUser creates access and refresh tokens for a user
func (s *AuthService) createTokensForUser(userID string) (*TokenResponse, error) {
	accessToken := helpers.GenerateAccessToken(userID, s.tokenExpiryHours)
	refreshToken := helpers.GenerateRefreshToken()

	err := s.userRepo.AddRefreshToken(userID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokenExpiryHours * 3600,
	}, nil
}

// rotateRefreshTokens rotates a refresh token (removes old, adds new)
func (s *AuthService) rotateRefreshTokens(userID, oldRefreshToken string) (*TokenResponse, error) {
	// Remove old refresh token
	if err := s.userRepo.RemoveRefreshToken(userID, oldRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to remove old refresh token: %w", err)
	}

	// Create new tokens
	return s.createTokensForUser(userID)
}
