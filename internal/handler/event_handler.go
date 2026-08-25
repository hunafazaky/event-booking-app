package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/response"
	"github.com/hunafazaky/event-booking-app/internal/service"
)

type EventHandler struct {
	service service.EventService
}

func NewEventHandler(service service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	dateTime, err := time.Parse(time.RFC3339, c.PostForm("date_time"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid datetime format")
		return
	}

	input := service.CreateEventInput{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Location:    c.PostForm("location"),
		DateTime:    dateTime,
		Image:       file,
		ImageName:   header.Filename,
	}

	event, err := h.service.Create(userID, input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "event created successfully", event)
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	events, meta, err := h.service.List(c.Query("search"), page, limit)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "data events retrieved", events, meta)
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameter value")
		return
	}

	event, err := h.service.GetByID(uint(eventID))
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data event retrieved", event)
}

func (h *EventHandler) GetMyEvents(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	events, err := h.service.GetByUser(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data event retrieved", events)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameter value")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	var (
		input service.UpdateEventInput
	)

	file, header, err := c.Request.FormFile("image")
	switch {
	case err == nil:
		defer file.Close()
		input.Image = file
		input.ImageName = header.Filename
	case errors.Is(err, http.ErrMissingFile):
		// no image sent — that's fine, Image/ImageName just stay zero-valued,
		// and EventService.Update's `if input.Image != nil` check skips
		// the upload entirely.
	default:
		response.Fail(c, http.StatusBadRequest, "failed to process image file")
		return
	}

	if c.PostForm("date_time") != "" {
		dateTime, err := time.Parse(time.RFC3339, c.PostForm("date_time"))
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid datetime format")
			return
		}
		input.DateTime = &dateTime
	}

	input.Name = c.PostForm("name")
	input.Description = c.PostForm("description")
	input.Location = c.PostForm("location")

	event, err := h.service.Update(userID, uint(eventID), input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data event updated", event)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameter value")
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	if err := h.service.Delete(userID, uint(eventID)); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data deleted successfully", nil)
}
