package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func CreateEvent(c *gin.Context) {
	userID, _ := c.Get("user_id")

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "error":   err.Error(),
			"message": "Failed to get image file.",
		})
		return
	}
	defer file.Close()

	// var event model.Event
	// if err := c.ShouldBindJSON(&event); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// event.UserID = userID.(int)

	fileName := header.Filename
	iKit := initImageKit()
	// uploadRes, err :=
	response, err := iKit.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload image.",
		})
		return
	}

	parseTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))

	event := model.Event{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Location:    c.PostForm("location"),
		DateTime:    parseTime,
		Image:       response.URL,
		ImageID:     response.FileID,
		UserID:      userID.(int),
	}

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
	userID, _ := c.Get("user_id")
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

	if event.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You are not authorized to update this event.",
		})
		return
	}

	// Bind the updated event data
	// var newEvent model.Event
	// err := c.ShouldBindJSON(&newEvent)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// "error":   err.Error(),
			"message": "Failed to get image file.",
		})
		return
	}

	fileName := header.Filename
	iKit := initImageKit()
	// uploadRes, err :=
	response, err := iKit.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File:     file,
		FileName: fileName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			// "error":   err.Error(),
			"message": "Failed to upload image.",
		})
		return
	}

	if event.ImageID != "" {
		err := iKit.Files.Delete(context.Background(), event.ImageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				// "error":   err.Error(),
				"message": "Failed to delete image.",
			})
			return
		}
		event.Image = response.URL
		event.ImageID = response.FileID
	}

	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}
	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}
	if dateTimeSTR := c.PostForm("date_time"); dateTimeSTR != "" {
		parseTime, err := time.Parse(time.RFC3339, dateTimeSTR)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				// "error":   err.Error(),
				"message": "Failed to parse date time.",
			})
			return
		}
		event.DateTime = parseTime
	}

	// Update the event with the new data
	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated.",
		"data":    event,
	})
}

func DeleteEventById(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var event model.Event
	paramID := c.Param("id")

	eventData := config.DB.First(&event, paramID).Error
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		return
	}

	if event.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You are not authorized to update this event.",
		})
		return
	}

	if event.ImageID != "" {
		ik := initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted.",
	})
}

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}
