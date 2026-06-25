package routes

import (
	"go-tiket-konser/handler"
	"go-tiket-konser/middleware"
	"go-tiket-konser/repository"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Global Middlewares
	r.Use(middleware.ApiKeyAuth())
	r.Use(middleware.RateLimiter(5))

	// Initialize layers
	concertRepo := repository.NewConcertRepository(db)
	concertService := service.NewConcertService(concertRepo)
	concertHandler := handler.NewConcertHandler(concertService)

	ticketCategoryRepo := repository.NewTicketCategoryRepository(db)
	ticketCategoryService := service.NewTicketCategoryService(ticketCategoryRepo, concertRepo)
	ticketCategoryHandler := handler.NewTicketCategoryHandler(ticketCategoryService)

	// Inisialisasi layer Booking
	customerRepo := repository.NewCustomerRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(db, bookingRepo, customerRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// inisialisasi layer authentication
	userRepo := repository.NewUserRepository(db)
	blacklistedTokenRepo := repository.NewBlacklistedTokenRepository(db)
	authService := service.NewAuthService(userRepo, blacklistedTokenRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Group routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Concerts routes
		api.POST("/concerts", concertHandler.CreateConcert)
		api.GET("/concerts", concertHandler.GetConcerts)
		api.GET("/concerts/:id", concertHandler.GetConcertByID)
		api.PUT("/concerts/:id", concertHandler.UpdateConcert)
		api.DELETE("/concerts/:id", concertHandler.DeleteConcert)

		// Ticket Categories routes
		api.POST("/ticket-categories", ticketCategoryHandler.CreateTicketCategory)
		api.GET("/ticket-categories", ticketCategoryHandler.GetTicketCategories)
		api.GET("/ticket-categories/:id", ticketCategoryHandler.GetTicketCategoryByID)
		api.PUT("/ticket-categories/:id", ticketCategoryHandler.UpdateTicketCategory)
		api.DELETE("/ticket-categories/:id", ticketCategoryHandler.DeleteTicketCategory)

		// Booking routes
		api.POST("/bookings", bookingHandler.CreateBooking)
		api.GET("/bookings/:id", bookingHandler.GetBookingByID)

		// Protected routes (JWT)
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(blacklistedTokenRepo))
		{
			protected.POST("/logout", authHandler.Logout)
			protected.GET("/profile", authHandler.GetProfile)
		}
	}

	return r
}
