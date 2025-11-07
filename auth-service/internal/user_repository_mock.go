package internal

import (
	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks the UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepositoryInterface
var _ UserRepositoryInterface = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string, user *User) error {
	args := m.Called(email, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id string, user *User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByGoogleID(googleID string, user *User) error {
	args := m.Called(googleID, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateGoogleID(userID string, googleID string) error {
	args := m.Called(userID, googleID)
	return args.Error(0)
}

func (m *MockUserRepository) SetPinWithExpiry(email, hashedPin string) error {
	args := m.Called(email, hashedPin)
	return args.Error(0)
}

func (m *MockUserRepository) ClearPin(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) SetRegisteredAt(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) AddRefreshToken(userID, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveRefreshToken(userID, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) ValidateRefreshToken(userID, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) ClearAllRefreshTokens(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
