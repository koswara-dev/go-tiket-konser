package service

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"math"

	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(req dto.UserQueryRequest) ([]dto.UserResponse, dto.PaginationMeta, error)
	GetUserByID(id uuid.UUID) (models.User, error)
	UpdateUser(id uuid.UUID, req *dto.UserUpdateRequest, updaterID uuid.UUID) (models.User, error)
	DeleteUser(id uuid.UUID, deleterID uuid.UUID) error
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

	var responses []dto.UserResponse = make([]dto.UserResponse, 0)
	for _, user := range users {
		var customerID uuid.UUID
		if user.Customer != nil {
			customerID = user.Customer.ID
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

func (s *userService) GetUserByID(id uuid.UUID) (models.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}
	return *user, nil
}

func (s *userService) UpdateUser(id uuid.UUID, req *dto.UserUpdateRequest, updaterID uuid.UUID) (models.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.Role = req.Role
	user.UpdatedBy = &updaterID

	if user.Customer != nil {
		user.Customer.Name = req.FullName
		user.Customer.Email = req.Email
		user.Customer.UpdatedBy = &updaterID
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *userService) DeleteUser(id uuid.UUID, deleterID uuid.UUID) error {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return models.ErrUserNotFound
	}
	user.DeletedBy = &deleterID
	if user.Customer != nil {
		user.Customer.DeletedBy = &deleterID
	}
	_ = s.userRepo.UpdateUser(user)
	return s.userRepo.DeleteUser(id)
}
