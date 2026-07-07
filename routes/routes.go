package routes

import (
	"go-tiket-konser/config"
	"go-tiket-konser/handler"
	"go-tiket-konser/middleware"
	"go-tiket-konser/repository"
	"go-tiket-konser/service"
	"os"

	_ "go-tiket-konser/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	// Initialize MongoDB audit log service and middleware
	auditService := service.NewAuditLogService(config.MongoDatabase)
	r.Use(middleware.AuditLogMiddleware(auditService))

	// Initialize SSE Notifications Broker and Handler
	notificationBroker := service.NewNotificationBroker(config.MongoDatabase)
	notificationHandler := handler.NewNotificationHandler(notificationBroker)

	// 1. inisiasi storage provider
	storageProvider := service.NewLocalStorageProvider("uploads", "http://localhost:8080")

	// 2. ekspose direktori upload
	r.Static("/uploads", "./uploads")

	// Global Middlewares
	r.Use(middleware.ApiKeyAuth())
	r.Use(middleware.RateLimiter(100))

	// Swagger Route
	if os.Getenv("APP_ENV") != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Initialize layers
	concertRepo := repository.NewConcertRepository(db)
	concertService := service.NewConcertService(concertRepo, notificationBroker)
	concertHandler := handler.NewConcertHandler(concertService, storageProvider)

	ticketCategoryRepo := repository.NewTicketCategoryRepository(db)
	ticketCategoryService := service.NewTicketCategoryService(ticketCategoryRepo, concertRepo)
	ticketCategoryHandler := handler.NewTicketCategoryHandler(ticketCategoryService)

	// Inisialisasi layer Booking
	customerRepo := repository.NewCustomerRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(db, bookingRepo, customerRepo, notificationBroker)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// inisialisasi layer authentication
	userRepo := repository.NewUserRepository(db)
	blacklistedTokenRepo := repository.NewBlacklistedTokenRepository(config.RedisClient)
	emailService := service.NewEmailService()
	authService := service.NewAuthService(userRepo, blacklistedTokenRepo, emailService)
	authHandler := handler.NewAuthHandler(authService)

	// Inisialisasi layer Users & Customers
	userServiceInstance := service.NewUserService(userRepo)
	userHandlerInstance := handler.NewUserHandler(userServiceInstance)

	customerServiceInstance := service.NewCustomerService(customerRepo)
	customerHandlerInstance := handler.NewCustomerHandler(customerServiceInstance)

	// Inisialisasi layer Chat
	chatHub := service.NewChatHub(config.MongoDatabase, db)
	chatHandler := handler.NewChatHandler(chatHub, blacklistedTokenRepo, db)

	// Group routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/verify-otp", authHandler.VerifyOTP)

		// Concerts routes (Public GET, Admin for POST/PUT/DELETE)
		api.GET("/concerts", concertHandler.GetConcerts)
		api.GET("/concerts/:id", concertHandler.GetConcertByID)
		api.POST("/concerts", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.CreateConcert)
		api.PUT("/concerts/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.UpdateConcert)
		api.DELETE("/concerts/:id", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.DeleteConcert)
		// upload thumbnail dan rules konser
		api.POST("/concerts/:id/thumbnail", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.UploadTumbnail)
		api.POST("/concerts/:id/rules", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), concertHandler.UploadRulesPDF)

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

		// Notifications (JWT Protected)
		api.GET("/notifications/stream", middleware.JWTAuth(blacklistedTokenRepo), notificationHandler.Stream)
		api.GET("/notifications", middleware.JWTAuth(blacklistedTokenRepo), notificationHandler.GetHistory)

		// Chat routes
		api.GET("/chat/ws", chatHandler.WS)
		api.GET("/chat/rooms", middleware.JWTAuth(blacklistedTokenRepo), middleware.RequireRole("admin"), chatHandler.GetRooms)
		api.GET("/chat/rooms/:roomId/messages", middleware.JWTAuth(blacklistedTokenRepo), chatHandler.GetMessages)

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
