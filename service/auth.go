package service

import (
	"errors"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecretKey = []byte("juara-coding-super-secret-key-2026-batch-1")

type JWTClaims struct {
	UserID     uint   `json:"user_id"`
	CustomerID uint   `json:"customer_id,omitempty"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(email, password, fullName string) (*models.User, error)
	Login(email, password string) (string, error)
	Logout(tokenString string) error
	GetProfile(id uint) (*models.User, error)
}

type authService struct {
	userRepo      repository.UserRepository
	blacklistRepo repository.BlacklistedTokenRepository
}

// New Auth Service
func NewAuthService(userRepo repository.UserRepository, blacklistRepo repository.BlacklistedTokenRepository) *authService {
	return &authService{userRepo: userRepo, blacklistRepo: blacklistRepo}
}

// Register
func (s *authService) Register(email, password, fullName string) (*models.User, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := "customer"

	var customer *models.Customer
	if role == "customer" {
		customer = &models.Customer{
			Name:  fullName,
			Email: email,
		}
	}

	// Create the new user
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
		Role:     role,
		Customer: customer,
	}

	// Save the user to the database
	err = s.userRepo.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Login
func (s *authService) Login(email, password string) (string, error) {
	// Check if email already exists
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email")
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate the JWT token
	var customerID uint
	if user.Customer != nil {
		customerID = uint(user.Customer.ID)
	}

	claims := JWTClaims{
		UserID:     user.ID,
		CustomerID: customerID,
		Email:      user.Email,
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey)
}

func (s *authService) Logout(tokenString string) error {
	return s.blacklistRepo.BlacklistToken(tokenString, time.Now().Add(24*time.Hour))
}

func (s *authService) GetProfile(id uint) (*models.User, error) {
	return s.userRepo.GetUserById(id)
}

// email verification

// change password
