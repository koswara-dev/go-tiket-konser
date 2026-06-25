package main

import (
	"go-tiket-konser/config"
	"go-tiket-konser/routes"
	"log"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// custom validators
func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// future date jika format tanggal benar dan waktunya minimal 7 hari dari hari ini
		_ = v.RegisterValidation("future_date", func(fl validator.FieldLevel) bool {
			dateStr, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return false
			}
			return date.After(time.Now().AddDate(0, 0, 7))
		})
	}
}

func main() {
	// Initialize database, migrations, and seeding
	config.InitDB()

	registerCustomValidators()

	// Setup routes and inject database connection
	r := routes.SetupRouter(config.DB)

	log.Println("Server is running on port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
