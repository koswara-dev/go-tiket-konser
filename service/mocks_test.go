package service

import (
	"go-tiket-konser/models"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(search string, limit int, offset int, sort string) ([]models.User, int64, error) {
	args := m.Called(search, limit, offset, sort)
	var users []models.User
	if args.Get(0) != nil {
		users = args.Get(0).([]models.User)
	}
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockBlacklistedTokenRepository implements repository.BlacklistedTokenRepository
type MockBlacklistedTokenRepository struct {
	mock.Mock
}

func (m *MockBlacklistedTokenRepository) BlacklistToken(token string, expiresAt time.Time) error {
	args := m.Called(token, expiresAt)
	return args.Error(0)
}

func (m *MockBlacklistedTokenRepository) IsTokenBlacklisted(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

// MockEmailService implements EmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendOTP(to string, otp string) error {
	args := m.Called(to, otp)
	return args.Error(0)
}

// MockConcertRepository implements repository.ConcertRepository
type MockConcertRepository struct {
	mock.Mock
}

func (m *MockConcertRepository) Create(concert *models.Concert) error {
	args := m.Called(concert)
	return args.Error(0)
}

func (m *MockConcertRepository) FindAll() ([]models.Concert, error) {
	args := m.Called()
	var concerts []models.Concert
	if args.Get(0) != nil {
		concerts = args.Get(0).([]models.Concert)
	}
	return concerts, args.Error(1)
}

func (m *MockConcertRepository) FindAllPaginated(search string, limit int, offset int, sort string) ([]models.Concert, int64, error) {
	args := m.Called(search, limit, offset, sort)
	var concerts []models.Concert
	if args.Get(0) != nil {
		concerts = args.Get(0).([]models.Concert)
	}
	return concerts, args.Get(1).(int64), args.Error(2)
}

func (m *MockConcertRepository) FindByID(id uuid.UUID) (models.Concert, error) {
	args := m.Called(id)
	return args.Get(0).(models.Concert), args.Error(1)
}

func (m *MockConcertRepository) Update(concert *models.Concert) error {
	args := m.Called(concert)
	return args.Error(0)
}

func (m *MockConcertRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
