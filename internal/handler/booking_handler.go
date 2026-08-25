package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/response"
	"github.com/hunafazaky/event-booking-app/internal/service"
)

type BookingInput struct {
	Phone   string `json:"phone"`
	EventID uint   `json:"event_id"`
}

type BookingHandler struct {
	service service.BookingService
}

func NewBookingHandler(service service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	var input service.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := h.service.Create(userID, input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "booking created successfully", booking)
}

func (h *BookingHandler) GetBooks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	bookings, err := h.service.GetByUser(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "booking retrieved successfully", bookings)
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameter value")
		return
	}

	if err := h.service.Delete(userID, uint(bookingID)); err != nil {
		response.FromError(c, err)
	}

	response.Success(c, http.StatusNoContent, "data deleted successfully", nil)
}
