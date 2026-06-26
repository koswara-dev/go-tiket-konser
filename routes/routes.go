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
	r.Use(middleware.RateLimiter(100))

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

	// Inisialisasi layer Users & Customers
	userServiceInstance := service.NewUserService(userRepo)
	userHandlerInstance := handler.NewUserHandler(userServiceInstance)

	customerServiceInstance := service.NewCustomerService(customerRepo)
	customerHandlerInstance := handler.NewCustomerHandler(customerServiceInstance)

	// Group routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Concerts routes (Public GET, Admin for POST/PUT/DELETE)
		api.GET("/concerts", concertHandler.GetConcerts)
		api.GET("/concerts/:id", concertHandler.GetConcertByID)
		api.POST("/concerts", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.CreateConcert)
		api.PUT("/concerts/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.UpdateConcert)
		api.DELETE("/concerts/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.DeleteConcert)

		// Ticket Categories routes (Public GET, Admin for POST/PUT/DELETE)
		api.GET("/ticket-categories", ticketCategoryHandler.GetTicketCategories)
		api.GET("/ticket-categories/:id", ticketCategoryHandler.GetTicketCategoryByID)
		api.POST("/ticket-categories", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), ticketCategoryHandler.CreateTicketCategory)
		api.PUT("/ticket-categories/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), ticketCategoryHandler.UpdateTicketCategory)
		api.DELETE("/ticket-categories/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), ticketCategoryHandler.DeleteTicketCategory)

		// Booking routes (JWT Protected)
		api.POST("/bookings", middleware.JWTAuth(blacklistedTokenRepo), bookingHandler.CreateBooking)
		api.GET("/bookings/:id", middleware.JWTAuth(blacklistedTokenRepo), bookingHandler.GetBookingByID)

		// Users routes (Admin / Own-check)
		api.GET("/users", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), userHandlerInstance.GetAllUsers)
		api.GET("/users/:id", middleware.JWTAuth(blacklistedTokenRepo), userHandlerInstance.GetUserByID)
		api.PUT("/users/:id", middleware.JWTAuth(blacklistedTokenRepo), userHandlerInstance.UpdateUser)
		api.DELETE("/users/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), userHandlerInstance.DeleteUser)

		// Customers routes (Admin / Own-check)
		api.GET("/customers", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), customerHandlerInstance.GetAllCustomers)
		api.GET("/customers/:id", middleware.JWTAuth(blacklistedTokenRepo), customerHandlerInstance.GetCustomerByID)
		api.PUT("/customers/:id", middleware.JWTAuth(blacklistedTokenRepo), customerHandlerInstance.UpdateCustomer)
		api.DELETE("/customers/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), customerHandlerInstance.DeleteCustomer)

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
