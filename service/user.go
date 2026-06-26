package service

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (models.User, error)
	UpdateUser(id uint, req *dto.UserUpdateRequest) (models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *userService) GetUserByID(id uint) (models.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}
	return *user, nil
}

func (s *userService) UpdateUser(id uint, req *dto.UserUpdateRequest) (models.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.Role = req.Role

	if user.Customer != nil {
		user.Customer.Name = req.FullName
		user.Customer.Email = req.Email
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.userRepo.GetUserById(id)
	if err != nil {
		return models.ErrUserNotFound
	}
	return s.userRepo.DeleteUser(id)
}
