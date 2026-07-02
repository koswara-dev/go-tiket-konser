package service

import (
	"errors"
	"go-tiket-konser/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistedTokenRepository)
	mockEmailService := new(MockEmailService)

	authServ := NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	email := "test@example.com"
	password := "securepassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	dummyUser := &models.User{
		Model:      gorm.Model{ID: 1},
		Email:      email,
		Password:   string(hashedPassword),
		Role:       "customer",
		FullName:   "Test User",
		IsVerified: true,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(dummyUser, nil)

	token, err := authServ.Login(email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistedTokenRepository)
	mockEmailService := new(MockEmailService)

	authServ := NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	email := "nonexistent@example.com"
	mockUserRepo.On("GetUserByEmail", email).Return((*models.User)(nil), errors.New("user not found"))

	token, err := authServ.Login(email, "password123")

	assert.Error(t, err)
	assert.Equal(t, "invalid email", err.Error())
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistedTokenRepository)
	mockEmailService := new(MockEmailService)

	authServ := NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	email := "test@example.com"
	password := "correctpassword"
	wrongPassword := "wrongpassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	dummyUser := &models.User{
		Model:      gorm.Model{ID: 1},
		Email:      email,
		Password:   string(hashedPassword),
		Role:       "customer",
		FullName:   "Test User",
		IsVerified: true,
	}

	mockUserRepo.On("GetUserByEmail", email).Return(dummyUser, nil)

	token, err := authServ.Login(email, wrongPassword)

	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, token)
	mockUserRepo.AssertExpectations(t)
}
