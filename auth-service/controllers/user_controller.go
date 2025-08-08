package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/arthadede/auth-service/constants"
	"github.com/arthadede/auth-service/models"
	"github.com/arthadede/auth-service/repositories"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	userRepo    *repositories.UserRepository
	config      *constants.Config
	oauthConfig *oauth2.Config
}

func NewUserController(userRepo *repositories.UserRepository, config *constants.Config) *UserController {
	oauthConfig := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  config.GoogleRedirectUrl,
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

func (c *UserController) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (c *UserController) generateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})

	tokenString, err := token.SignedString([]byte(c.config.JWTSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

func (c *UserController) LoginUser(email, password string) (map[string]string, error) {
	user, err := c.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.ComparePassword(password) {
		return nil, errors.New("invalid email or password")
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

func (c *UserController) RefreshToken(refreshToken string) (map[string]string, error) {
	user, err := c.userRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, err := c.generateAccessToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := c.generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	err = c.userRepo.UpdateRefreshToken(user.ID.Hex(), newRefreshToken)
	if err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

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

	resetURL := fmt.Sprintf("%s?token=%s", c.config.FEResetPassword, resetToken)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link below to reset your password:\n\n%s\n\nIf you did not request a password reset, please ignore this email.", resetURL)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, email, subject, body)

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

func (c *UserController) VerifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(c.config.JWTSecret), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	user, err := c.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetEnvironment returns the current environment (production, development, etc.)
func (c *UserController) GetEnvironment() string {
	return c.config.Environment
}

// GetTokenExpiryHours returns the configured token expiry in hours
func (c *UserController) GetTokenExpiryHours() int {
	return c.config.TokenExpiryHours
}

// GetOAuthRedirectURL returns the frontend URL to redirect to after OAuth
func (c *UserController) GetOAuthRedirectURL() string {
	return c.config.FEOAuthRedirect
}
