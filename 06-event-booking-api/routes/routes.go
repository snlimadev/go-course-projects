package routes

import (
	"example.com/event-booking-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	registerDocsRoutes(server)
	registerPublicRoutes(server)
	registerPrivateRoutes(server)
}

func registerDocsRoutes(server *gin.Engine) {
	server.StaticFile("/docs", "./docs/swagger-ui.html")
	server.StaticFile("/docs/swagger.json", "./docs/swagger.json")
	server.StaticFile("/docs/favicon.png", "./docs/favicon.png")
}

func registerPublicRoutes(server *gin.Engine) {
	// Authentication
	server.POST("/auth/login", login)
	server.POST("/auth/refresh", refresh)
	server.POST("/auth/logout", logout)

	// Users
	server.POST("/users", createUser)

	// Events
	server.GET("/events", getEvents)
	server.GET("/events/:id", middlewares.OptionalAuthMiddleware(), getEventDetails)

	// Bookings
	server.GET("/events/:id/bookings", getBookings)
}

func registerPrivateRoutes(server *gin.Engine) {
	protected := server.Group("/")

	protected.Use(middlewares.AuthMiddleware())
	{
		// Users
		protected.DELETE("/users/:id", deleteUser)

		// Events
		protected.POST("/events", createEvent)
		protected.PUT("/events/:id", updateEvent)
		protected.DELETE("/events/:id", deleteEvent)

		// Bookings
		protected.POST("/events/:id/bookings", createBooking)
		protected.DELETE("/events/:id/bookings/:bookingId", deleteBooking)
	}
}
