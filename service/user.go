package service

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"math"
)

type UserService interface {
	GetAllUsers(req dto.UserQueryRequest) ([]dto.UserResponse, dto.PaginationMeta, error)
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

func (s *userService) GetAllUsers(req dto.UserQueryRequest) ([]dto.UserResponse, dto.PaginationMeta, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	users, totalData, err := s.userRepo.GetAllUsers(req.Search, limit, offset, req.Sort)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		var customerID uint
		if user.Customer != nil {
			customerID = uint(user.Customer.ID)
		}
		responses = append(responses, dto.UserResponse{
			ID:         user.ID,
			FullName:   user.FullName,
			Email:      user.Email,
			Role:       user.Role,
			CustomerID: customerID,
		})
	}

	totalPage := 0
	if totalData > 0 {
		totalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	}

	meta := dto.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return responses, meta, nil
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
