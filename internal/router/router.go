// internal/router/router.go
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/handler"
	"github.com/hunafazaky/event-booking-app/internal/middleware"

	"github.com/hunafazaky/event-booking-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// scalarReferenceHTML is a static page embedding Scalar's API Reference
// component via CDN. It only needs a URL to fetch the OpenAPI spec from —
// /openapi.json below, which is the same generated spec Swagger UI reads,
// just served as raw JSON instead of wrapped in a UI.
//
// data-configuration is Scalar's own settings object — theme:"default" and
// layout:"classic" together give the plainest, least "app-like" look
// Scalar offers: a single scrolling document instead of a three-pane
// dashboard. This block is intentionally generic — copy it as-is into any
// other project's docs page for the same clean baseline.
const scalarReferenceHTML = `<!doctype html>
<html>
<head>
	<title>Event Booking API Reference</title>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
	<script
		id="api-reference"
		data-url="/openapi.json"
		data-configuration='{"theme":"elysiajs","darkMode":true}'
	></script>
	<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

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
		api.GET("", eventHandler.GetEvents)
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

	// --- API documentation ---
	// Both UIs read the SAME generated spec (docs.SwaggerInfo) — main.go
	// overrides its Host once at startup if cfg.PublicHost is set, so
	// neither UI needs to know about deployment vs local on its own.
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(docs.SwaggerInfo.ReadDoc()))
	})

	server.GET("/reference", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(scalarReferenceHTML))
	})
}
