package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type BookingInput struct {
	Phone   string `json:"phone"`
	EventID uint   `json:"event_id"`
}

func CreateBookingEvent(c *gin.Context) {
	// 1. Ambil user_id dari context secara aman
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// 2. Cast tipe data dengan aman (mencegah panic/500 jika tipe di JWT adalah float64/int)
	var userID uint
	switch v := userIDVal.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
		return
	}

	var input BookingInput
	var booking model.Booking

	// 3. Validasi JSON Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}

	// 4. Cek apakah booking sudah ada
	bookingExist := config.DB.Where("user_id = ? AND event_id = ?", userID, input.EventID).First(&booking).Error
	if bookingExist == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Booking already exists"})
		return
	}

	// 5. Buat kode booking & simpan
	bookingCode := fmt.Sprintf("BK-%sE%dU%d", time.Now().Format("20060102"), input.EventID, userID)

	bookData := model.Booking{
		Phone:       input.Phone,
		EventID:     input.EventID,
		BookingCode: bookingCode,
		UserID:      userID,
	}

	if err := config.DB.Create(&bookData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking created successfully"})
}

func GetBooks(c *gin.Context) {
	var booking []model.Booking
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// 2. Cast tipe data dengan aman (mencegah panic/500 jika tipe di JWT adalah float64/int)
	var userID uint
	switch v := userIDVal.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
		return
	}

	err := config.DB.Preload("Event").Preload("Event.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("user_id = ?", userID).Find(&booking).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"booking": booking,
	})
}

func DeleteBooking(c *gin.Context) {
	var booking model.Booking
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// 2. Cast tipe data dengan aman (mencegah panic/500 jika tipe di JWT adalah float64/int)
	var userID uint
	switch v := userIDVal.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
		return
	}

	paramID := c.Param("id")
	err := config.DB.First(&booking, paramID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Booking not found"})
		return
	}

	if booking.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
		return
	}

	err = config.DB.Unscoped().Delete(&booking).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted"})
}
