package handler

import (
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/config"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"gorm.io/gorm"
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
		UserID:      userID.(uint),
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

	// config.DB.Find(&events)
	//

	// Inisiasi data query
	query := config.DB.Model(&model.Event{})

	// tangkap flter
	search := c.Query("search")

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// hitung total data
	var totalRows int64
	query.Count(&totalRows)

	// tangkap parameter
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "6")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 6
	}

	// hitung offset
	offset := (page - 1) * limit

	// hitung total page
	totalPage := int(math.Ceil(float64(totalRows) / float64(limit)))

	// eksekusi
	if err := query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve events.",
			"error":   err.Error(),
		})
		return
	}

	// if page != "" {
	// 	query = query.Offset((pageNumber - 1) * limitNumber)
	// }

	// if limit != "" {
	// 	query = query.Limit(limitNumber)
	// }

	query.Find(&events)
	c.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved.",
		"data":    events,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total_rows": totalRows,
			"total_page": totalPage,
		},
	})
}

// GET /events - Retrieve all events
func GetEventById(c *gin.Context) {
	var event model.Event
	paramID := c.Param("id")

	var eventData = config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Preload("Booking").Preload("Booking.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).First(&event, paramID).Error

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

func GetEventsByUser(c *gin.Context) {
	var events []model.Event
	userID, _ := c.Get("user_id")

	err := config.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("user_id", userID).Find(&events).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Events retrieved.",
		"data":    events,
	})
}

func UpdateEventById(c *gin.Context) {
	// 1. Ambil user_id dari context dengan aman (mencegah panic)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Cast tipe data secara aman ke uint (menyesuaikan tipe ID di GORM)
	var currentUserID uint
	switch v := userIDVal.(type) {
	case uint:
		currentUserID = v
	case int:
		currentUserID = uint(v)
	case float64: // Banyak library JWT mengurai number sebagai float64
		currentUserID = uint(v)
	}

	paramID := c.Param("id")
	var event model.Event

	// 2. Cari Event di Database
	if err := config.DB.First(&event, paramID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Event not found.",
		})
		return
	}

	// 3. Cek Hak Akses (Otorisasi)
	if event.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You are not authorized to update this event.",
		})
		return
	}

	// 4. Proses Upload Gambar (Opsional)
	file, header, err := c.Request.FormFile("image") // Shortcut Gin
	if err != nil {

		if !errors.Is(err, http.ErrMissingFile) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to process image file",
			})
			return
		}
	} else {
		iKit := initImageKit()

		// Upload gambar baru ke ImageKit
		response, err := iKit.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File:     file,
			FileName: header.Filename,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to upload image.",
			})
			return
		}

		// Jika sebelumnya ada gambar lama, hapus dari ImageKit (Opsional/Best effort)
		if event.ImageID != "" {
			_ = iKit.Files.Delete(context.Background(), event.ImageID)
			// Catatan: Jika hapus gambar lama gagal, proses TETAP dilanjutkan
		}

		// PERBAIKAN BUG: Selalu update URL & FileID baru di luar blok 'if' hapus
		event.Image = response.URL
		event.ImageID = response.FileID
	}

	// 5. Update Field Teks
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
				"message": "Failed to parse date time. Format must be RFC3339.",
			})
			return
		}
		event.DateTime = parseTime
	}

	// 6. Simpan Perubahan ke Database & Cek Error
	if err := config.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update event in database.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully.",
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

	if event.UserID != userID.(uint) {
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
