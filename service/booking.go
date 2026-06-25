package service

import (
	"errors"
	"fmt"
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingService interface {
	CreateBooking(req *dto.BookingRequest) (models.Booking, error)
	GetBookingByID(id int) (models.Booking, error)
}

type bookingService struct {
	db           *gorm.DB
	bookingRepo  repository.BookingRepository
	customerRepo repository.CustomerRepository
}

func NewBookingService(db *gorm.DB, bookingRepo repository.BookingRepository, customerRepo repository.CustomerRepository) BookingService {
	return &bookingService{
		db:           db,
		bookingRepo:  bookingRepo,
		customerRepo: customerRepo,
	}
}

func (s *bookingService) CreateBooking(req *dto.BookingRequest) (models.Booking, error) {
	var finalBooking models.Booking

	// Jalankan Transaksi Database Otomatis
	errTx := s.db.Transaction(func(tx *gorm.DB) error {
		// Instansiasi repository khusus di dalam scope transaksi (Menggunakan tx)
		txCustomerRepo := repository.NewCustomerRepository(tx)
		txTicketCategoryRepo := repository.NewTicketCategoryRepository(tx)
		txBookingRepo := repository.NewBookingRepository(tx)

		// 1. Validasi Customer
		customer, err := txCustomerRepo.FindByID(req.CustomerID)
		if err != nil {
			return errors.New("customer tidak ditemukan di database")
		}

		var totalAmount float64
		var details []models.BookingDetail
		bookingCode := fmt.Sprintf("TIX-%d-%d", time.Now().Unix(), time.Now().UnixNano()%1000)

		// 2. Buat reservasi utama
		finalBooking = models.Booking{
			CustomerID:  uint(customer.ID),
			BookingCode: bookingCode,
			TotalAmount: 0,
		}
		if err := txBookingRepo.Create(&finalBooking); err != nil {
			return err
		}

		// 3. Proses setiap item tiket
		for _, item := range req.BookingDetails {
			var category models.TicketCategory

			// PROTEKSI RACE CONDITION: Kunci baris kategori tiket menggunakan FOR UPDATE
			errLock := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&category, item.TicketCategoryID).Error
			if errLock != nil {
				return fmt.Errorf("kategori tiket ID %d tidak ditemukan", item.TicketCategoryID)
			}

			// Validasi Quota
			if category.AvailableQuota < item.Quantity {
				return fmt.Errorf("kuota tiket '%s' tidak mencukupi (Tersisa: %d, Permintaan: %d)",
					category.Name, category.AvailableQuota, item.Quantity)
			}

			// Potong Quota Tiket
			category.AvailableQuota -= item.Quantity
			if err := txTicketCategoryRepo.Update(&category); err != nil {
				return err
			}

			// Simpan Detail
			subtotal := category.Price * float64(item.Quantity)
			totalAmount += subtotal

			detail := models.BookingDetail{
				BookingID:        int(finalBooking.ID),
				TicketCategoryID: category.ID,
				Quantity:         item.Quantity,
				SubTotal:         subtotal,
			}
			if err := txBookingRepo.CreateDetail(&detail); err != nil {
				return err
			}

			details = append(details, detail)
		}

		// 4. Update total_amount akhir pada tabel bookings
		finalBooking.TotalAmount = totalAmount
		finalBooking.Details = details
		if err := txBookingRepo.Update(&finalBooking); err != nil {
			return err
		}

		return nil // Commit otomatis jika tanpa error
	})

	return finalBooking, errTx
}

func (s *bookingService) GetBookingByID(id int) (models.Booking, error) {
	return s.bookingRepo.FindByID(id)
}
