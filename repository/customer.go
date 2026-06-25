package repository

import (
	"go-tiket-konser/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByID(id int) (models.Customer, error)
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
