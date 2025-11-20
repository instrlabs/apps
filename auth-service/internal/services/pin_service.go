package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"github.com/instrlabs/auth-service/internal/helpers"
	"github.com/instrlabs/auth-service/internal/models"
	"github.com/instrlabs/auth-service/internal/repositories"
)

// PinService handles PIN generation and validation
type PinService struct {
	userRepo     repositories.UserRepositoryInterface
	emailService helpers.EmailSender
	pinEnabled   bool
}

// NewPinService creates a new PIN service
func NewPinService(userRepo repositories.UserRepositoryInterface, emailService helpers.EmailSender, pinEnabled bool) *PinService {
	return &PinService{
		userRepo:     userRepo,
		emailService: emailService,
		pinEnabled:   pinEnabled,
	}
}

// GenerateAndSendPIN generates and sends a PIN to the user's email
func (s *PinService) GenerateAndSendPIN(email string) error {
	// Find or create user
	_, err := s.findOrCreateUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to find/create user: %w", err)
	}

	// Generate and set PIN
	pin, err := s.generateAndSetPin(email)
	if err != nil {
		return fmt.Errorf("failed to set PIN: %w", err)
	}

	// Handle non-PIN enabled case
	if !s.pinEnabled {
		return nil
	}

	// Send email with PIN
	emailService := helpers.NewEmailService()
	emailService.SendPinEmail(email, pin)

	return nil
}

// ValidatePIN validates a PIN for a user
func (s *PinService) ValidatePIN(email, pin string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check PIN validity
	if !user.ComparePin(pin) {
		return nil, fmt.Errorf("invalid PIN")
	}

	// Check if PIN is expired
	if user.IsPinExpired() {
		return nil, fmt.Errorf("PIN has expired")
	}

	return user, nil
}

// findOrCreateUserByEmail finds a user by email or creates a new one
func (s *PinService) findOrCreateUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return user, nil
	}

	// User not found, create new one
	newUser := models.NewUser(email)
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// generateAndSetPin generates a PIN and sets it with expiry
func (s *PinService) generateAndSetPin(email string) (string, error) {
	pin := "000000"
	if s.pinEnabled {
		pin = helpers.GenerateSixDigitPIN()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	err = s.userRepo.SetPinWithExpiry(email, string(hash))
	if err != nil {
		return "", err
	}

	return pin, nil
}
