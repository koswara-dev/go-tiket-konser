package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-tiket-konser/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("failed to get database instance: ", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	// auto migrate
	err = DB.AutoMigrate(&models.User{}, &models.Customer{}, &models.Concert{}, &models.TicketCategory{},
		&models.Booking{}, &models.BookingDetail{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("Database connected and migrated successfully")

	// Initialize Redis and Mongo config
	InitRedis()
	InitMongo()

	// seeder 5 data concert, 2 data customer, 2 ticket category
	var count int64
	DB.Model(&models.Concert{}).Count(&count)
	if count == 0 {
		initialConcerts := []models.Concert{
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				},
				Title:       "Coldplay Music of the Spheres World Tour Jakarta",
				Description: "Konser perdana band asal Inggris, Coldplay, di Indonesia yang memukau ratusan ribu penonton dengan gelang Xyloband yang menyala warna-warni.",
				Date:        time.Date(2023, 11, 15, 20, 0, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "completed",
			},
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
				Title:       "Blackpink [Born Pink] World Tour Jakarta",
				Description: "Konser megah dari girlgroup K-Pop fenomenal, Blackpink, yang berhasil meremajakan Jakarta menjadi lautan cahaya merah muda selama dua hari berturut-turut.",
				Date:        time.Date(2023, 3, 11, 19, 0, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "completed",
			},
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				},
				Title:       "Metallica Live in Jakarta 2013",
				Description: "Konser sejarah kembalinya raja thrash metal dunia ke Indonesia setelah penantian 20 tahun, dihadiri oleh puluhan ribu pecinta musik cadas dari berbagai generasi.",
				Date:        time.Date(2013, 8, 25, 20, 0, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "completed",
			},
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				},
				Title:       "Bruno Mars Live in Jakarta 2026",
				Description: "Konser tur dunia dari solois legendaris Bruno Mars yang membawakan deretan lagu hitsnya dengan koreografi dan vokal yang sangat enerjik.",
				Date:        time.Date(2026, 6, 22, 20, 0, 0, 0, time.UTC),
				Venue:       "Jakarta International Stadium (JIS), Jakarta",
				Status:      "active",
			},
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				},
				Title:       "Pesta Rakyat Dewa 19 - 30 Tahun Berkarya",
				Description: "Konser selebrasi 3 dekade salah satu band rock terbesar di Indonesia, Dewa 19, yang memboyong 4 vokalis dan 5 drummer dalam satu panggung.",
				Date:        time.Date(2026, 6, 22, 19, 30, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "active",
			},
		}

		if err := DB.Create(&initialConcerts).Error; err != nil {
			log.Println("failed to seed initial concerts: ", err)
		} else {
			log.Println("Database successfully seeded with 5 initial concerts")
		}
	}

	// seeder 2 ticket category
	var ticketCategoryCount int64
	DB.Model(&models.TicketCategory{}).Count(&ticketCategoryCount)
	if ticketCategoryCount == 0 {
		initialTicketCategories := []models.TicketCategory{
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("33333333-3333-3333-3333-333333333331"),
				},
				Name:           "Gold",
				Price:          1000000,
				ConcertID:      uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				TotalQuota:     100,
				AvailableQuota: 50,
			},
			{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("33333333-3333-3333-3333-333333333332"),
				},
				Name:           "Silver",
				Price:          500000,
				ConcertID:      uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				TotalQuota:     200,
				AvailableQuota: 150,
			},
		}

		if err := DB.Create(&initialTicketCategories).Error; err != nil {
			log.Println("failed to seed initial ticket categories: ", err)
		} else {
			log.Println("Database successfully seeded with 2 initial ticket categories")
		}
	}

	// seeder users and customers
	var customerCount int64
	DB.Model(&models.Customer{}).Count(&customerCount)
	if customerCount == 0 {
		// Budi
		budiUserID := uuid.MustParse("11111111-1111-1111-1111-111111111110")
		budiCustID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		hashedPasswordBudi, _ := bcrypt.GenerateFromPassword([]byte("Indonesia"), bcrypt.DefaultCost)
		budiUser := models.User{
			BaseModel: models.BaseModel{
				ID: budiUserID,
			},
			FullName:   "Budi Santoso",
			Email:      "budi.santoso@gmail.com",
			Password:   string(hashedPasswordBudi),
			Role:       "customer",
			IsVerified: true,
			Customer: &models.Customer{
				BaseModel: models.BaseModel{
					ID: budiCustID,
				},
				UserID: budiUserID,
				Name:   "Budi Santoso",
				Email:  "budi.santoso@gmail.com",
			},
		}
		if err := DB.Create(&budiUser).Error; err != nil {
			log.Println("failed to seed Budi user/customer: ", err)
		}

		// Ani
		aniUserID := uuid.MustParse("22222222-2222-2222-2222-222222222220")
		aniCustID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
		hashedPasswordAni, _ := bcrypt.GenerateFromPassword([]byte("Indonesia"), bcrypt.DefaultCost)
		aniUser := models.User{
			BaseModel: models.BaseModel{
				ID: aniUserID,
			},
			FullName:   "Ani Wijaya",
			Email:      "ani.wijaya@gmail.com",
			Password:   string(hashedPasswordAni),
			Role:       "customer",
			IsVerified: true,
			Customer: &models.Customer{
				BaseModel: models.BaseModel{
					ID: aniCustID,
				},
				UserID: aniUserID,
				Name:   "Ani Wijaya",
				Email:  "ani.wijaya@gmail.com",
			},
		}
		if err := DB.Create(&aniUser).Error; err != nil {
			log.Println("failed to seed Ani user/customer: ", err)
		}

		log.Println("Database successfully seeded with 2 users and customers")
	}

	// seeder admin user
	var userCount int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Indonesia"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("failed to hash admin password: ", err)
		} else {
			adminUser := models.User{
				BaseModel: models.BaseModel{
					ID: uuid.MustParse("99999999-9999-9999-9999-999999999999"),
				},
				FullName:   "Admin Konser",
				Email:      "adminkonser@gmail.com",
				Password:   string(hashedPassword),
				Role:       "admin",
				IsVerified: true,
			}
			if err := DB.Create(&adminUser).Error; err != nil {
				log.Println("failed to seed admin user: ", err)
			} else {
				log.Println("Database successfully seeded with admin user")
			}
		}
	}
}
