package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/models"
)

func main() {
	server := gin.Default()
	server.SetTrustedProxies([]string{"localhost"})

	{
		api := server.Group("/api")
		api.GET("/events", getEvents)
		api.POST("/events", createEvent)
	}

	server.Run(":8080")
}

// GET /events - Retrieve all events
func getEvents(c *gin.Context) {
	events := models.GetAllEvents()

	c.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved.",
		"data":    events,
	})
}

// POST /events - Create a new event
func createEvent(c *gin.Context) {
	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data.",
			"error":   err.Error(),
		})
		return
	}

	event.Id = 1
	event.UserId = 1
	event.DateTime = time.Now()

	event.Save()

	c.JSON(http.StatusCreated, gin.H{
		"message": "New Event Created",
		"data":    event,
	})
}
