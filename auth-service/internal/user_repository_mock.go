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

func (m *MockUserRepository) Create(user *User) *User {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*User)
}

func (m *MockUserRepository) FindByEmail(email string) *User {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*User)
}

func (m *MockUserRepository) FindByID(id string) *User {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*User)
}

func (m *MockUserRepository) FindByGoogleID(googleID string) *User {
	args := m.Called(googleID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*User)
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

func (m *MockUserRepository) ValidateRefreshToken(userID, token string) bool {
	args := m.Called(userID, token)
	return args.Bool(0)
}

func (m *MockUserRepository) ClearAllRefreshTokens(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
