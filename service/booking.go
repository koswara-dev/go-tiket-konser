package service

import (
	"errors"
	"fmt"
	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingService interface {
	CreateBooking(req *dto.BookingRequest, userID uuid.UUID) (models.Booking, error)
	GetBookingByID(id uuid.UUID, userID uuid.UUID, role string) (models.Booking, error)
}

type bookingService struct {
	db           *gorm.DB
	bookingRepo  repository.BookingRepository
	customerRepo repository.CustomerRepository
	broker       *NotificationBroker
}

func NewBookingService(db *gorm.DB, bookingRepo repository.BookingRepository, customerRepo repository.CustomerRepository, broker *NotificationBroker) BookingService {
	return &bookingService{
		db:           db,
		bookingRepo:  bookingRepo,
		customerRepo: customerRepo,
		broker:       broker,
	}
}

func (s *bookingService) CreateBooking(req *dto.BookingRequest, userID uuid.UUID) (models.Booking, error) {
	var finalBooking models.Booking

	// Jalankan Transaksi Database Otomatis
	errTx := s.db.Transaction(func(tx *gorm.DB) error {
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
		bookingCode := fmt.Sprintf("TIX-%d-%s", time.Now().Unix(), uuid.New().String()[:8])

		// 2. Buat reservasi utama
		finalBooking = models.Booking{
			BaseModel: models.BaseModel{
				CreatedBy: &userID,
			},
			CustomerID:  customer.ID,
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
			errLock := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&category, "id = ?", item.TicketCategoryID).Error
			if errLock != nil {
				return fmt.Errorf("kategori tiket ID %s tidak ditemukan", item.TicketCategoryID.String())
			}

			// Validasi Quota
			if category.AvailableQuota < item.Quantity {
				return fmt.Errorf("kuota tiket '%s' tidak mencukupi (Tersisa: %d, Permintaan: %d)",
					category.Name, category.AvailableQuota, item.Quantity)
			}

			// Potong Quota Tiket
			category.AvailableQuota -= item.Quantity
			// Set UpdatedBy
			category.UpdatedBy = &userID
			if err := txTicketCategoryRepo.Update(&category); err != nil {
				return err
			}

			// Simpan Detail
			subtotal := category.Price * float64(item.Quantity)
			totalAmount += subtotal

			detail := models.BookingDetail{
				BaseModel: models.BaseModel{
					CreatedBy: &userID,
				},
				BookingID:        finalBooking.ID,
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

	if errTx == nil && s.broker != nil {
		// Send notification
		msg := fmt.Sprintf("Booking dengan kode %s berhasil dibuat. Total pembayaran: Rp %.2f", finalBooking.BookingCode, finalBooking.TotalAmount)
		_ = s.broker.SendNotification(customerIDToUserID(s.db, finalBooking.CustomerID), "Booking Berhasil", msg)
	}

	return finalBooking, errTx
}

func customerIDToUserID(db *gorm.DB, customerID uuid.UUID) string {
	var customer models.Customer
	if err := db.Select("user_id").First(&customer, "id = ?", customerID).Error; err == nil {
		return customer.UserID.String()
	}
	return ""
}

func (s *bookingService) GetBookingByID(id uuid.UUID, userID uuid.UUID, role string) (models.Booking, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return booking, models.ErrBookingNotFound
	}

	// MITIGASI IDOR:
	// Jika user aktif adalah customer, pastikan ID pembeli di DB cocok dengan CustomerID dari user tersebut
	if role == "customer" {
		customer, err := s.customerRepo.FindByUserID(userID)
		if err != nil || booking.CustomerID != customer.ID {
			// Mengembalikan 404 (ErrBookingNotFound) demi keamanan informasi
			return models.Booking{}, models.ErrBookingNotFound
		}
	}

	return booking, nil
}
