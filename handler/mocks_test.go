package handler

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockAuthService implements service.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(email, password, fullName string) (*models.User, error) {
	args := m.Called(email, password, fullName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) Logout(tokenString string) error {
	args := m.Called(tokenString)
	return args.Error(0)
}

func (m *MockAuthService) GetProfile(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) VerifyOTP(email, otp string) error {
	args := m.Called(email, otp)
	return args.Error(0)
}

// MockConcertService implements service.ConcertService
type MockConcertService struct {
	mock.Mock
}

func (m *MockConcertService) CreateConcert(concert *models.Concert) error {
	args := m.Called(concert)
	return args.Error(0)
}

func (m *MockConcertService) GetAllConcerts(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error) {
	args := m.Called(req)
	var concerts []dto.ConcertResponse
	if args.Get(0) != nil {
		concerts = args.Get(0).([]dto.ConcertResponse)
	}
	return concerts, args.Get(1).(dto.PaginationMeta), args.Error(2)
}

func (m *MockConcertService) GetConcertByID(id uuid.UUID) (models.Concert, error) {
	args := m.Called(id)
	return args.Get(0).(models.Concert), args.Error(1)
}

func (m *MockConcertService) UpdateConcert(concert *models.Concert) error {
	args := m.Called(concert)
	return args.Error(0)
}

func (m *MockConcertService) DeleteConcert(id uuid.UUID, deleterID uuid.UUID) error {
	args := m.Called(id, deleterID)
	return args.Error(0)
}

// MockStorageProvider implements service.StorageProvider
type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error) {
	args := m.Called(file, c)
	return args.String(0), args.Error(1)
}
