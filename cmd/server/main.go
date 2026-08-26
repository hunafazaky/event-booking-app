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

	// Named (not blank) import now — main.go reads docs.SwaggerInfo
	// directly below to override Host once deployed. The generated
	// docs.go's init() side-effect (registering the spec) still runs the
	// same way on import either way.
	"github.com/hunafazaky/event-booking-app/docs"
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

	// The @host annotation above is what generates the spec's DEFAULT
	// host — correct for local dev (localhost:8080), wrong once deployed.
	// Overriding it here, once at startup, means both doc UIs (Swagger UI
	// and Scalar, wired in router.go) automatically point "try it out"
	// requests at the real deployed domain instead of localhost.
	if cfg.PublicHost != "" {
		docs.SwaggerInfo.Host = cfg.PublicHost
		docs.SwaggerInfo.Schemes = []string{"https"}
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

	server.Run(":" + cfg.Port)
}
