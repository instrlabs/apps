package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/smtp"
	"time"

	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/models"
	"github.com/arthadede/auth-service/repositories"
	"github.com/golang-jwt/jwt/v5"
)

// UserController handles business logic for user operations
type UserController struct {
	userRepo    *repositories.UserRepository
	config      *constants.Config
	oauthConfig *oauth2.Config
}

func NewUserController(userRepo *repositories.UserRepository, config *constants.Config) *UserController {
	redirectUrl := "https://" + config.Hostname + config.Port + config.GoogleRedirectPath

	oauthConfig := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &UserController{
		userRepo:    userRepo,
		config:      config,
		oauthConfig: oauthConfig,
	}
}

// RegisterUser handles the business logic for user registration
func (c *UserController) RegisterUser(email, password string) (*models.User, error) {
	user, err := models.NewUser(email, password)
	if err != nil {
		return nil, err
	}

	err = c.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// generateRefreshToken generates a random refresh token
func (c *UserController) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// generateAccessToken generates a JWT access token for a user
func (c *UserController) generateAccessToken(userID string) (string, error) {
	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(c.config.TokenExpiryHours)).Unix(), // Token expiration from config
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(c.config.JWTSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

// LoginUser handles the business logic for user login
func (c *UserController) LoginUser(email, password string) (map[string]string, error) {
	// Find user by email
	user, err := c.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !user.ComparePassword(password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate access token
	accessToken, err := c.generateAccessToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := c.generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Save refresh token to user
	err = c.userRepo.UpdateRefreshToken(user.ID.Hex(), refreshToken)
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

// RefreshToken handles the business logic for refreshing an access token
func (c *UserController) RefreshToken(refreshToken string) (map[string]string, error) {
	// Find user by refresh token
	user, err := c.userRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := c.generateAccessToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := c.generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Save new refresh token to user
	err = c.userRepo.UpdateRefreshToken(user.ID.Hex(), newRefreshToken)
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

// generateResetToken generates a random reset token
func (c *UserController) generateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (c *UserController) sendResetEmail(email, resetToken string) error {
	if c.config.Environment == "development" {
		log.Printf("Password reset token for %s: %s", email, resetToken)
		return nil
	}

	from := c.config.EmailFrom
	to := []string{email}

	// Construct email body
	resetURL := fmt.Sprintf("%s?token=%s", c.config.FEResetPassword, resetToken)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link below to reset your password:\n\n%s\n\nIf you did not request a password reset, please ignore this email.", resetURL)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, email, subject, body)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", c.config.SMTPUsername, c.config.SMTPPassword, c.config.SMTPHost)
	err := smtp.SendMail(c.config.SMTPHost+":"+c.config.SMTPPort, auth, from, to, []byte(message))
	return err
}

func (c *UserController) RequestPasswordReset(email string) error {
	_, err := c.userRepo.FindByEmail(email)
	if err != nil {
		return nil
	}

	resetToken, err := c.generateResetToken()
	if err != nil {
		return errors.New("failed to generate reset token")
	}

	expiry := time.Now().Add(time.Hour * time.Duration(c.config.ResetTokenExpiryHours))

	err = c.userRepo.SetResetToken(email, resetToken, expiry)
	if err != nil {
		return errors.New("failed to save reset token")
	}

	err = c.sendResetEmail(email, resetToken)
	if err != nil {
		return errors.New("failed to send reset email")
	}

	return nil
}

func (c *UserController) ResetPassword(resetToken, newPassword string) error {
	user, err := c.userRepo.FindByResetToken(resetToken)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	err = c.userRepo.UpdatePassword(user.ID.Hex(), string(hashedPassword))
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (c *UserController) GetGoogleAuthURL() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	return c.oauthConfig.AuthCodeURL(state)
}

func (c *UserController) HandleGoogleCallback(code string) (map[string]string, error) {
	token, err := c.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("failed to exchange code for token")
	}

	client := c.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("failed to get user info from Google")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response from Google")
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, errors.New("failed to parse user info from Google")
	}

	user, err := c.userRepo.FindByGoogleID(userInfo.ID)
	if err != nil {
		user, err = c.userRepo.FindByEmail(userInfo.Email)
		if err != nil {
			user = models.NewGoogleUser(userInfo.Email, userInfo.ID)
			err = c.userRepo.Create(user)
			if err != nil {
				return nil, errors.New("failed to create user")
			}
		} else {
			err = c.userRepo.UpdateGoogleID(user.ID.Hex(), userInfo.ID)
			if err != nil {
				return nil, errors.New("failed to update user with Google ID")
			}
		}
	}

	accessToken, err := c.generateAccessToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	refreshToken, err := c.generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	err = c.userRepo.UpdateRefreshToken(user.ID.Hex(), refreshToken)
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}
