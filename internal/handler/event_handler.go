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

// CreateEvent godoc
// @Summary Create New Event
// @Description Creates a new event for the authenticated user.
// @Security BearerAuth
// @Tags Events
// @Accept mpfd
// @Produce json
// @Param name formData string true "Event name"
// @Param description formData string true "Event description"
// @Param location formData string true "Event location"
// @Param datetime formData string true "RFC3339 datetime"
// @Param image formData file true "Event image"
// @Success 201 {object} response.Envelope{data=dto.EventResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /events [post]
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

	dateTime, err := time.Parse(time.RFC3339, c.PostForm("datetime"))
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

// GetEvents godoc
// @Summary Get All Events
// @Description Returns a paginated list of events, optionally filtered by search term.
// @Tags Events
// @Produce json
// @Param search query string false "Search by name or description"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Results per page (default 6)"
// @Success 200 {object} response.Envelope{data=[]dto.EventResponse,meta=dto.EventListMeta}
// @Failure 500 {object} response.Envelope
// @Router /events [get]
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

// GetEventByID godoc
// @Summary Get an Event by ID
// @Description Returns an event by ID, including its bookings.
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} response.Envelope{data=dto.EventDetailResponse}
// @Failure 404 {object} response.Envelope
// @Router /events/{id} [get]
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

// GetEventsMine godoc
// @Summary Get Current User Events
// @Description Returns events created by the authenticated user.
// @Security BearerAuth
// @Tags Events
// @Produce json
// @Success 200 {object} response.Envelope{data=[]dto.EventResponse}
// @Failure 401 {object} response.Envelope
// @Router /events/mine [get]
func (h *EventHandler) GetEventsMine(c *gin.Context) {
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

// UpdateEvent godoc
// @Summary Update an Event
// @Description Updates an event owned by the authenticated user.
// @Security BearerAuth
// @Tags Events
// @Accept mpfd
// @Produce json
// @Param id path int true "Event ID"
// @Param name formData string false "Event name"
// @Param description formData string false "Event description"
// @Param location formData string false "Event location"
// @Param datetime formData string false "RFC3339 datetime"
// @Param image formData file false "Event image"
// @Success 200 {object} response.Envelope{data=dto.EventResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /events/{id} [put]
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

	if c.PostForm("datetime") != "" {
		dateTime, err := time.Parse(time.RFC3339, c.PostForm("datetime"))
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

// DeleteEvent godoc
// @Summary Delete an Event
// @Description Deletes an event owned by the authenticated user.
// @Security BearerAuth
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Router /events/{id} [delete]
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
