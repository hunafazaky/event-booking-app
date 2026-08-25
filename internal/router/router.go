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
		api := server.Group("/api/events")
		api.GET("/", eventHandler.GetEvents)
		api.GET("/:id", eventHandler.GetEventByID)
	}

	{
		api := server.Group("/api/auth")
		api.POST("/signup", userHandler.SignUp)
		api.POST("/signin", userHandler.SignIn)
	}

	{
		protectedApi := server.Group("/api")
		protectedApi.Use(middleware.RequireAuth(cfg.JWTSecret))
		protectedApi.GET("/auth/me", userHandler.GetMe)
		protectedApi.POST("/events", eventHandler.CreateEvent)
		protectedApi.PUT("/events/:id", eventHandler.UpdateEvent)
		protectedApi.DELETE("/events/:id", eventHandler.DeleteEvent)
		protectedApi.GET("/events/mine", eventHandler.GetEventsMine)
		protectedApi.POST("/bookings", bookingHandler.CreateBooking)
		protectedApi.GET("/bookings", bookingHandler.GetBooks)
		protectedApi.DELETE("/bookings/:id", bookingHandler.DeleteBooking)
	}
}
