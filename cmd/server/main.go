package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/handler"
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
		api.GET("/events", handler.GetEvents)
		api.GET("/events/:id", handler.GetEventById)
		api.POST("/events", handler.CreateEvent)
		api.PUT("/events/:id", handler.UpdateEventById)
		api.DELETE("/events/:id", handler.DeleteEventById)
	}

	server.Run(":8080")
}
