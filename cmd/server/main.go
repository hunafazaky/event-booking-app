package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/handler"
	"github.com/hunafazaky/event-booking-app/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// .env is only used for local development — in production the real
	// environment variables are already set, so a missing file here isn't
	// fatal, it's expected. What DOES matter is config.Load() below, which
	// validates that whatever the source, the required values are present.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on process environment variables.")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	config.ConnectDB(cfg)

	server := gin.Default()
	server.SetTrustedProxies([]string{"localhost"})

	{
		api := server.Group("/api/events")
		api.GET("/", handler.GetEvents)
		api.GET("/:id", handler.GetEventById)
	}

	{
		api := server.Group("/api/auth")
		api.POST("/signup", handler.SignUpUser)
		api.POST("/signin", handler.SignInUser)
	}

	{
		protectedApi := server.Group("/api")
		protectedApi.Use(middleware.RequireAuth())
		protectedApi.GET("/auth/me", handler.GetAuthUser)
		protectedApi.POST("/events", handler.CreateEvent)
		protectedApi.PUT("/events/:id", handler.UpdateEventById)
		protectedApi.DELETE("/events/:id", handler.DeleteEventById)
		protectedApi.GET("/events/user", handler.GetEventsByUser)
		protectedApi.POST("/booking", handler.CreateBookingEvent)
		protectedApi.GET("/booking", handler.GetBooks)
		protectedApi.DELETE("/booking/:id", handler.DeleteBooking)
	}

	server.Run(":" + cfg.Port)
}
