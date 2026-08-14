package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/config"
	"github.com/hunafazaky/event-booking-app/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Err: ", err)
	}

	config.ConnectDB()

	server := gin.Default()
	server.SetTrustedProxies([]string{"localhost"})

	{
		api := server.Group("/api")
		api.GET("/events", handlers.GetEvents)
		api.GET("/events/:id", handlers.GetEventById)
		api.POST("/events", handlers.CreateEvent)
		api.PUT("/events/:id", handlers.UpdateEventById)
		api.DELETE("/events/:id", handlers.DeleteEventById)
	}

	server.Run(":8080")
}
