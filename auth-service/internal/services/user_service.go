package services

import (
	"fmt"

	"github.com/instrlabs/auth-service/internal/models"
	"github.com/instrlabs/auth-service/internal/repositories"
)

// UserService handles user management operations
type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves user profile by ID
func (s *UserService) GetProfile(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// Create creates a new user
func (s *UserService) Create(email string) (*models.User, error) {
	// Check if user already exists
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Create new user
	newUser := models.NewUser(email)
	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// Update updates user information
func (s *UserService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

// FindByEmail finds a user by email
func (s *UserService) FindByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// FindByID finds a user by ID
func (s *UserService) FindByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
