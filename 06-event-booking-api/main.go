package main

import (
	"log"

	"example.com/event-booking-api/config"
	"example.com/event-booking-api/db"
	"example.com/event-booking-api/routes"
	"example.com/event-booking-api/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lpernett/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Could not load environment variables.")
	}

	config.Load()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterCustomValidators(v)
	}

	db.InitDB()
	defer db.DB.Close()

	server := gin.Default()
	routes.RegisterRoutes(server)

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
