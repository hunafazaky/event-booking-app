package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/model"
)

func CreateEvent(c *gin.Context) {
	var event model.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	event.UserID = 1

	config.DB.Create(&event)

	c.JSON(http.StatusCreated, gin.H{
		"message": "New Event Created",
		"data":    event,
	})
}

// GET /events - Retrieve all events
func GetEvents(c *gin.Context) {
	var events []model.Event

	config.DB.Find(&events)
	c.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved.",
		"data":    events,
	})
}

// GET /events - Retrieve all events
func GetEventById(c *gin.Context) {
	var event model.Event
	paramID := c.Param("id")

	var eventData = config.DB.First(&event, paramID).Error

	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event retrieved.",
		"data":    event,
	})
}

func UpdateEventById(c *gin.Context) {
	// Get the event by ID
	var event model.Event
	paramID := c.Param("id")

	// Find the event by ID
	eventData := config.DB.First(&event, paramID).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		return
	}

	// Bind the updated event data
	var newEvent model.Event
	err := c.ShouldBindJSON(&newEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update the event with the new data
	config.DB.Model(&event).Updates(newEvent)
	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated.",
		"data":    event,
	})
}

func DeleteEventById(c *gin.Context) {
	var event model.Event
	paramID := c.Param("id")

	eventData := config.DB.First(&event, paramID).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted.",
	})
}
