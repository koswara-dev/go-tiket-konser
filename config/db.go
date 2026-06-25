package config

import (
	"log"
	"time"

	"go-tiket-konser/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "postgres://postgres:secret45@localhost:5434/eticketdb?sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// auto migrate
	err = DB.AutoMigrate(&models.Concert{}, &models.TicketCategory{}, &models.Customer{},
		&models.Booking{}, &models.BookingDetail{}, &models.User{}, &models.BlacklistedToken{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("Database connected and migrated successfully")

	// seeder 5 data concert, 2 data customer, 2 ticket category
	var count int64
	DB.Model(&models.Concert{}).Count(&count)
	if count == 0 {
		initialConcerts := []models.Concert{
			{
				Title:       "Coldplay Music of the Spheres World Tour Jakarta",
				Description: "Konser perdana band asal Inggris, Coldplay, di Indonesia yang memukau ratusan ribu penonton dengan gelang Xyloband yang menyala warna-warni.",
				Date:        time.Date(2023, 11, 15, 20, 0, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "completed",
			},
			{
				Title:       "Blackpink [Born Pink] World Tour Jakarta",
				Description: "Konser megah dari girlgroup K-Pop fenomenal, Blackpink, yang berhasil meremajakan Jakarta menjadi lautan cahaya merah muda selama dua hari berturut-turut.",
				Date:        time.Date(2023, 3, 11, 19, 0, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "completed",
			},
			{
				Title:       "Metallica Live in Jakarta 2013",
				Description: "Konser sejarah kembalinya raja thrash metal dunia ke Indonesia setelah penantian 20 tahun, dihadiri oleh puluhan ribu pecinta musik cadas dari berbagai generasi.",
				Date:        time.Date(2013, 8, 25, 20, 0, 0, 0, time.UTC),
				Venue:       "Stadion Utama Gelora Bung Karno (SUGBK), Jakarta",
				Status:      "completed",
			},
			{
				Title:       "Bruno Mars Live in Jakarta 2026",
				Description: "Konser tur dunia dari solois legendaris Bruno Mars yang membawakan deretan lagu hitsnya dengan koreografi dan vokal yang sangat enerjik.",
				Date:        time.Date(2026, 6, 22, 20, 0, 0, 0, time.UTC),
				Venue:       "Jakarta International Stadium (JIS), Jakarta",
				Status:      "active",
			},
			{
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

	// seeder 2 customer
	var customerCount int64
	DB.Model(&models.Customer{}).Count(&customerCount)
	if customerCount == 0 {
		initialCustomers := []models.Customer{
			{
				Name:  "Budi Santoso",
				Email: "budi.santoso@gmail.com",
			},
			{
				Name:  "Ani Wijaya",
				Email: "ani.wijaya@gmail.com",
			},
		}

		if err := DB.Create(&initialCustomers).Error; err != nil {
			log.Println("failed to seed initial customers: ", err)
		} else {
			log.Println("Database successfully seeded with 2 initial customers")
		}
	}

	// seeder 2 ticket category
	var ticketCategoryCount int64
	DB.Model(&models.TicketCategory{}).Count(&ticketCategoryCount)
	if ticketCategoryCount == 0 {
		initialTicketCategories := []models.TicketCategory{
			{
				Name:           "Gold",
				Price:          1000000,
				ConcertID:      4,
				TotalQuota:     100,
				AvailableQuota: 50,
			},
			{
				Name:           "Silver",
				Price:          500000,
				ConcertID:      4,
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
}
