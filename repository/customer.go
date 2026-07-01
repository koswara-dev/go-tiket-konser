package repository

import (
	"go-tiket-konser/models"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByID(id int) (models.Customer, error)
	FindByUserID(userID uint) (models.Customer, error)
	GetAllCustomers(search string, limit int, offset int, sort string) ([]models.Customer, int64, error)
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

func (c *customerRepo) GetAllCustomers(search string, limit int, offset int, sort string) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := c.db.Model(&models.Customer{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	case "email_asc":
		query = query.Order("email ASC")
	case "email_desc":
		query = query.Order("email DESC")
	default:
		query = query.Order("id DESC")
	}

	err := query.Limit(limit).Offset(offset).Find(&customers).Error
	return customers, total, err
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
