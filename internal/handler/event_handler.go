package handler

import (
	"net/http"
	"strconv"

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

	var input service.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.service.Create(userID, input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "event created successfully", event)
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameter value")
		return
	}

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

func (h *EventHandler) GetEventByUser(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	event, err := h.service.GetByID(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data event retrieved", event)
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

	var input service.UpdateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

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

	response.Success(c, http.StatusNoContent, "data deleted successfully", nil)
}
