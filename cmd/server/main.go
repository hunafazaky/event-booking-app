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
		api := server.Group("/api/events")
		api.GET("/", handler.GetEvents)
		api.GET("/:id", handler.GetEventById)
		api.POST("", handler.CreateEvent)
		api.PUT("/:id", handler.UpdateEventById)
		api.DELETE("/:id", handler.DeleteEventById)
	}

	{
		api := server.Group("/api/auth")
		api.POST("/signup", handler.SignUpUser)
	}

	server.Run(":8080")
}
