package repository

import (
	"go-tiket-konser/models"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *models.Booking) error
	CreateDetail(detail *models.BookingDetail) error
	FindByID(id int) (models.Booking, error)
	Update(booking *models.Booking) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) CreateDetail(detail *models.BookingDetail) error {
	return r.db.Create(detail).Error
}

func (r *bookingRepository) FindByID(id int) (models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Customer").
		Preload("Details.TicketCategory.Concert").
		First(&booking, id).Error
	return booking, err
}

func (r *bookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}
