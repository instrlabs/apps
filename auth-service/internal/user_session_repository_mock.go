package internal

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserSessionRepository mocks the UserSessionRepositoryInterface
type MockUserSessionRepository struct {
	mock.Mock
}

// Ensure MockUserSessionRepository implements UserSessionRepositoryInterface
var _ UserSessionRepositoryInterface = (*MockUserSessionRepository)(nil)

func (m *MockUserSessionRepository) CreateUserSession(userID, ipAddress, userAgent string) (*UserSession, error) {
	args := m.Called(userID, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserSession), args.Error(1)
}

func (m *MockUserSessionRepository) FindUserSessionByRefreshToken(refreshToken string) (*UserSession, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserSession), args.Error(1)
}

func (m *MockUserSessionRepository) FindUserSessionByID(id primitive.ObjectID, userID string) (*UserSession, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserSession), args.Error(1)
}

func (m *MockUserSessionRepository) ValidateUserSession(id primitive.ObjectID, userID, deviceHash string) bool {
	args := m.Called(id, userID, deviceHash)
	return args.Bool(0)
}

func (m *MockUserSessionRepository) UpdateUserSessionActivity(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserSessionRepository) DeactivateUserSession(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserSessionRepository) GetUserSessions(userID string) ([]UserSession, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]UserSession), args.Error(1)
}

func (m *MockUserSessionRepository) ClearExpiredUserSessions(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserSessionRepository) ClearAllUserSessions(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserSessionRepository) UpdateUserSessionRefreshToken(id primitive.ObjectID, newRefreshToken string) error {
	args := m.Called(id, newRefreshToken)
	return args.Error(0)
}
