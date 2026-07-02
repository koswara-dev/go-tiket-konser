package main

import (
	"go-tiket-konser/config"
	"go-tiket-konser/routes"
	"go-tiket-konser/utils/logger"
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

// @title           Go Ticket Concert API
// @version         1.0
// @description     This is a ticket concert API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apiKey  ApiKeyAuth
// @in                          header
// @name                        x-api-key

// @securityDefinitions.apiKey  BearerAuth
// @in                          header
// @name                        Authorization

func main() {

	// Initialize logger
	logger.InitLogger()
	logger.Log.Info("Server is starting...")

	// Initialize database, migrations, and seeding
	config.InitDB()

	registerCustomValidators()

	// Setup routes and inject database connection
	r := routes.SetupRouter(config.DB)

	logger.Log.Info("Server is running on port :8080")
	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("Failed to start server: ", err)
	}
}
