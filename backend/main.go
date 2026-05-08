package main

import (
	"log"
	"os"
	"strings"
	"time"

	"ocs-room-booking/db"
	"ocs-room-booking/handlers"
	"ocs-room-booking/middleware"
	"ocs-room-booking/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found configured. Using environment variables.")
	}

	// Initialize the database
	db.Connect()
	services.EnsureAdminUser()
	if strings.EqualFold(os.Getenv("SEED_ROOMS"), "true") {
		services.SeedRooms()
	}

	r := gin.Default()

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	api.POST("/auth/login", handlers.Login)
	api.GET("/rooms", handlers.ListAvailableRooms)

	api.Use(middleware.AuthMiddleware())
	api.GET("/me", handlers.GetMe)
	api.POST("/bookings", handlers.CreateBooking)
	api.GET("/bookings", handlers.ListMyBookings)
	api.PATCH("/bookings/:id", handlers.UpdateBookingStatus)

	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly())
	admin.POST("/users", handlers.CreateUser)
	admin.GET("/users", handlers.ListUsers)
	admin.PATCH("/users/:id", handlers.UpdateUser)
	admin.DELETE("/users/:id", handlers.DeactivateUser)
	admin.POST("/rooms", handlers.CreateRoom)
	admin.GET("/rooms", handlers.ListRooms)
	admin.PATCH("/rooms/:id", handlers.UpdateRoom)
	admin.DELETE("/rooms/:id", handlers.DeleteRoom)
	admin.GET("/bookings", handlers.ListAllBookings)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
