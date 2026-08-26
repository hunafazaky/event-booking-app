// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/handler"
	"github.com/hunafazaky/event-booking-app/internal/middleware"
)

// Setup registers every route on server. It takes the already-constructed
// handlers and config as parameters — it builds NOTHING itself, it only
// wires paths to methods. Construction (repos → services → handlers)
// stays entirely in main.go; this file's only job is routing.
func Setup(
	server *gin.Engine,
	cfg *config.Config,
	userHandler *handler.UserHandler,
	eventHandler *handler.EventHandler,
	bookingHandler *handler.BookingHandler,
) {
	{
		api := server.Group("/api/auth")
		api.POST("/signup", userHandler.SignUp)
		api.POST("/signin", userHandler.SignIn)
	}

	{
		api := server.Group("/api")
		api.GET("/events", eventHandler.GetEvents)
		api.GET("/events/:id", eventHandler.GetEventByID)
	}

	{
		api := server.Group("/api")
		api.Use(middleware.RequireAuth(cfg.JWTSecret))
		api.GET("/auth/me", userHandler.GetMe)
		api.POST("/events", eventHandler.CreateEvent)
		api.PUT("/events/:id", eventHandler.UpdateEvent)
		api.DELETE("/events/:id", eventHandler.DeleteEvent)
		api.GET("/events/mine", eventHandler.GetEventsMine)
		api.POST("/bookings", bookingHandler.CreateBooking)
		api.GET("/bookings", bookingHandler.GetBooks)
		api.DELETE("/bookings/:id", bookingHandler.DeleteBooking)
	}
}
