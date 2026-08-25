package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/handler"
	"github.com/hunafazaky/event-booking-app/internal/repository"
	"github.com/hunafazaky/event-booking-app/internal/router"
	"github.com/hunafazaky/event-booking-app/internal/service"
	"github.com/joho/godotenv"

	// documentation
	_ "github.com/hunafazaky/event-booking-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Event Booking API
// @version         1.0
// @description     API for browsing events and managing bookings.
// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and your JWT.
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

	db := config.ConnectDB(cfg)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Services — note the extra dependencies beyond just their repo
	uploader := service.NewImageKitUploader(cfg.ImageKitPrivateKey)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	eventService := service.NewEventService(eventRepo, uploader)
	bookingService := service.NewBookingService(bookingRepo, eventRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	eventHandler := handler.NewEventHandler(eventService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	server := gin.Default()
	server.SetTrustedProxies([]string{"localhost"})

	router.Setup(server, cfg, userHandler, eventHandler, bookingHandler)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Run(":" + cfg.Port)
}
