package service

import (
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
)

type CustomerService interface {
	GetAllCustomers() ([]models.Customer, error)
	GetCustomerByID(id int) (models.Customer, error)
	UpdateCustomer(id int, req *dto.CustomerUpdateRequest) (models.Customer, error)
	DeleteCustomer(id int) error
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{customerRepo: customerRepo}
}

func (s *customerService) GetAllCustomers() ([]models.Customer, error) {
	return s.customerRepo.GetAllCustomers()
}

func (s *customerService) GetCustomerByID(id int) (models.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return models.Customer{}, models.ErrCustomerNotFound
	}
	return customer, nil
}

func (s *customerService) UpdateCustomer(id int, req *dto.CustomerUpdateRequest) (models.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return models.Customer{}, models.ErrCustomerNotFound
	}

	customer.Name = req.Name
	customer.Email = req.Email

	err = s.customerRepo.UpdateCustomer(&customer)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (s *customerService) DeleteCustomer(id int) error {
	_, err := s.customerRepo.FindByID(id)
	if err != nil {
		return models.ErrCustomerNotFound
	}
	return s.customerRepo.DeleteCustomer(id)
}
