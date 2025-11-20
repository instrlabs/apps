package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/instrlabs/auth-service/internal/models"
	"github.com/instrlabs/auth-service/internal/repositories"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleUserInfo represents Google user information
type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// OAuthService handles OAuth authentication flows
type OAuthService struct {
	userRepo           repositories.UserRepositoryInterface
	authService        *AuthService
	googleClientID     string
	googleClientSecret string
	googleRedirectUrl  string
	webUrl             string
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(
	userRepo repositories.UserRepositoryInterface,
	authService *AuthService,
	googleClientID, googleClientSecret, googleRedirectUrl, webUrl string,
) *OAuthService {
	return &OAuthService{
		userRepo:           userRepo,
		authService:        authService,
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		googleRedirectUrl:  googleRedirectUrl,
		webUrl:             webUrl,
	}
}

// InitiateGoogleLogin initiates Google OAuth login
func (s *OAuthService) InitiateGoogleLogin() (string, error) {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	conf := &oauth2.Config{
		ClientID:     s.googleClientID,
		ClientSecret: s.googleClientSecret,
		RedirectURL:  s.googleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return conf.AuthCodeURL(state), nil
}

// HandleGoogleCallback handles Google OAuth callback
func (s *OAuthService) HandleGoogleCallback(code string) (*TokenResponse, error) {
	// Exchange code for token
	token, err := s.exchangeGoogleCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get Google user info
	googleInfo, err := s.getGoogleUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	user, err := s.findOrCreateGoogleUser(googleInfo.ID, googleInfo.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find/create user: %w", err)
	}

	// Create tokens for user
	return s.authService.createTokensForUser(user.ID.Hex())
}

// BuildGoogleRedirectURL builds the redirect URL with tokens
func (s *OAuthService) BuildGoogleRedirectURL(accessToken, refreshToken string, expiresIn int) string {
	return fmt.Sprintf("%s?access_token=%s&refresh_token=%s&token_type=Bearer&expires_in=%s",
		s.webUrl,
		url.QueryEscape(accessToken),
		url.QueryEscape(refreshToken),
		strconv.Itoa(expiresIn),
	)
}

// exchangeGoogleCode exchanges authorization code for access token
func (s *OAuthService) exchangeGoogleCode(code string) (*oauth2.Token, error) {
	conf := &oauth2.Config{
		ClientID:     s.googleClientID,
		ClientSecret: s.googleClientSecret,
		RedirectURL:  s.googleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf.Exchange(context.Background(), code)
}

// getGoogleUserInfo gets user information from Google
func (s *OAuthService) getGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	conf := &oauth2.Config{
		ClientID:     s.googleClientID,
		ClientSecret: s.googleClientSecret,
		RedirectURL:  s.googleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleInfo GoogleUserInfo
	err = json.Unmarshal(data, &googleInfo)
	return &googleInfo, err
}

// findUserByGoogleID finds a user by Google ID
func (s *OAuthService) findUserByGoogleID(googleID string) (*models.User, error) {
	user, err := s.userRepo.FindByGoogleID(googleID)
	if err != nil {
		return nil, fmt.Errorf("user not found by Google ID")
	}
	return user, nil
}

// findUserByEmailAndLinkGoogle finds a user by email and links Google ID
func (s *OAuthService) findUserByEmailAndLinkGoogle(email, googleID string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found by email")
	}

	user.GoogleID = &googleID
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update Google ID: %w", err)
	}

	return user, nil
}

// createNewGoogleUser creates a new user from Google info
func (s *OAuthService) createNewGoogleUser(email, googleID string) (*models.User, error) {
	newUser := models.NewGoogleUser(email, googleID)
	err := s.userRepo.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google user: %w", err)
	}
	return newUser, nil
}

// findOrCreateGoogleUser finds or creates a user from Google OAuth info
func (s *OAuthService) findOrCreateGoogleUser(googleID, email string) (*models.User, error) {
	user, err := s.findUserByGoogleID(googleID)
	if err == nil {
		return user, nil
	}

	user, err = s.findUserByEmailAndLinkGoogle(email, googleID)
	if err == nil {
		return user, nil
	}

	return s.createNewGoogleUser(email, googleID)
}
