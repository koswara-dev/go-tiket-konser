package service

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"math"

	"github.com/google/uuid"
)

type CustomerService interface {
	GetAllCustomers(req dto.CustomerQueryRequest) ([]dto.CustomerResponse, dto.PaginationMeta, error)
	GetCustomerByID(id uuid.UUID) (models.Customer, error)
	UpdateCustomer(id uuid.UUID, req *dto.CustomerUpdateRequest, updaterID uuid.UUID) (models.Customer, error)
	DeleteCustomer(id uuid.UUID, deleterID uuid.UUID) error
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{customerRepo: customerRepo}
}

func (s *customerService) GetAllCustomers(req dto.CustomerQueryRequest) ([]dto.CustomerResponse, dto.PaginationMeta, error) {
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

	customers, totalData, err := s.customerRepo.GetAllCustomers(req.Search, limit, offset, req.Sort)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	var responses []dto.CustomerResponse = make([]dto.CustomerResponse, 0)
	for _, cust := range customers {
		responses = append(responses, dto.CustomerResponse{
			ID:        cust.ID,
			UserID:    cust.UserID,
			Name:      cust.Name,
			Email:     cust.Email,
			CreatedAt: cust.CreatedAt,
			UpdatedAt: cust.UpdatedAt,
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

func (s *customerService) GetCustomerByID(id uuid.UUID) (models.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return models.Customer{}, models.ErrCustomerNotFound
	}
	return customer, nil
}

func (s *customerService) UpdateCustomer(id uuid.UUID, req *dto.CustomerUpdateRequest, updaterID uuid.UUID) (models.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return models.Customer{}, models.ErrCustomerNotFound
	}

	customer.Name = req.Name
	customer.Email = req.Email
	customer.UpdatedBy = &updaterID

	err = s.customerRepo.UpdateCustomer(&customer)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (s *customerService) DeleteCustomer(id uuid.UUID, deleterID uuid.UUID) error {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return models.ErrCustomerNotFound
	}
	customer.DeletedBy = &deleterID
	_ = s.customerRepo.UpdateCustomer(&customer)
	return s.customerRepo.DeleteCustomer(id)
}
