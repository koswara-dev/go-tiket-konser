package repository

import (
	"go-tiket-konser/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByID(id int) (models.Customer, error)
	FindByUserID(userID uint) (models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(id int) error
}

type customerRepo struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepo{db: db}
}

func (c *customerRepo) FindByID(id int) (models.Customer, error) {
	var customer models.Customer
	err := c.db.Where("id = ?", id).First(&customer).Error
	return customer, err
}

func (c *customerRepo) FindByUserID(userID uint) (models.Customer, error) {
	var customer models.Customer
	err := c.db.Where("user_id = ?", userID).First(&customer).Error
	return customer, err
}

func (c *customerRepo) GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	err := c.db.Find(&customers).Error
	return customers, err
}

func (c *customerRepo) UpdateCustomer(customer *models.Customer) error {
	return c.db.Save(customer).Error
}

func (c *customerRepo) DeleteCustomer(id int) error {
	var customer models.Customer
	if err := c.db.First(&customer, id).Error; err != nil {
		return err
	}
	return c.db.Delete(&customer).Error
}
